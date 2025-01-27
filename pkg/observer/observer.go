/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package observer

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"github.com/trustbloc/edge-core/pkg/log"
	"github.com/trustbloc/sidetree-core-go/pkg/api/operation"
	"github.com/trustbloc/sidetree-core-go/pkg/api/protocol"
	txnapi "github.com/trustbloc/sidetree-core-go/pkg/api/txn"

	"github.com/trustbloc/orb/pkg/activitypub/vocab"
	"github.com/trustbloc/orb/pkg/anchor/graph"
	anchorinfo "github.com/trustbloc/orb/pkg/anchor/info"
	"github.com/trustbloc/orb/pkg/anchor/util"
	"github.com/trustbloc/orb/pkg/errors"
	"github.com/trustbloc/orb/pkg/hashlink"
)

var logger = log.New("orb-observer")

// AnchorGraph interface to access anchors.
type AnchorGraph interface {
	Read(cid string) (*verifiable.Credential, error)
	GetDidAnchors(cid, suffix string) ([]graph.Anchor, error)
}

// OperationStore interface to access operation store.
type OperationStore interface {
	Put(ops []*operation.AnchoredOperation) error
}

// OperationFilter filters out operations before they are persisted.
type OperationFilter interface {
	Filter(uniqueSuffix string, ops []*operation.AnchoredOperation) ([]*operation.AnchoredOperation, error)
}

type didAnchors interface {
	PutBulk(dids []string, cid string) error
}

// Publisher publishes anchors and DIDs to a message queue for processing.
type Publisher interface {
	PublishAnchor(anchor *anchorinfo.AnchorInfo) error
	PublishDID(did string) error
}

type pubSub interface {
	Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error)
	Publish(topic string, messages ...*message.Message) error
	Close() error
}

type metricsProvider interface {
	ProcessAnchorTime(value time.Duration)
	ProcessDIDTime(value time.Duration)
}

// Outbox defines an ActivityPub outbox.
type Outbox interface {
	Post(activity *vocab.ActivityType) (*url.URL, error)
}

type outboxProvider func() Outbox

// Option is an option for observer.
type Option func(opts *Observer)

// WithDiscoveryDomain sets optional discovery domain hint (used for did equivalent ids).
func WithDiscoveryDomain(domain string) Option {
	return func(opts *Observer) {
		opts.discoveryDomain = domain
	}
}

// Providers contains all of the providers required by the TxnProcessor.
type Providers struct {
	ProtocolClientProvider protocol.ClientProvider
	AnchorGraph
	DidAnchors didAnchors
	PubSub     pubSub
	Metrics    metricsProvider
	Outbox     outboxProvider
}

// Observer receives transactions over a channel and processes them by storing them to an operation store.
type Observer struct {
	*Providers

	pubSub          *PubSub
	discoveryDomain string
}

// New returns a new observer.
func New(providers *Providers, opts ...Option) (*Observer, error) {
	o := &Observer{
		Providers: providers,
	}

	ps, err := NewPubSub(providers.PubSub, o.handleAnchor, o.processDID)
	if err != nil {
		return nil, err
	}

	o.pubSub = ps

	// apply options
	for _, opt := range opts {
		opt(o)
	}

	return o, nil
}

// Start starts observer routines.
func (o *Observer) Start() {
	o.pubSub.Start()
}

// Stop stops the observer.
func (o *Observer) Stop() {
	o.pubSub.Stop()
}

// Publisher returns the publisher that adds anchors and DIDs to a message queue for processing.
func (o *Observer) Publisher() Publisher {
	return o.pubSub
}

func (o *Observer) handleAnchor(anchor *anchorinfo.AnchorInfo) error {
	logger.Debugf("observing anchor - hashlink [%s], local hashlink [%s], attributedTo [%s]",
		anchor.Hashlink, anchor.Hashlink, anchor.AttributedTo)

	startTime := time.Now()

	defer o.Metrics.ProcessAnchorTime(time.Since(startTime))

	anchorInfo, err := o.AnchorGraph.Read(anchor.Hashlink)
	if err != nil {
		logger.Warnf("Failed to get anchor[%s] node from anchor graph: %s", anchor.Hashlink, err.Error())

		return err
	}

	logger.Debugf("successfully read anchor[%s] from anchor graph", anchor.Hashlink)

	if err := o.processAnchor(anchor, anchorInfo); err != nil {
		logger.Warnf(err.Error())

		return err
	}

	return nil
}

func (o *Observer) processDID(did string) error {
	logger.Debugf("processing out-of-system did[%s]", did)

	startTime := time.Now()

	defer func() {
		o.Metrics.ProcessDIDTime(time.Since(startTime))
	}()

	cidWithHint, suffix, err := getDidParts(did)
	if err != nil {
		logger.Warnf("process did failed for did[%s]: %s", did, err.Error())

		return err
	}

	anchors, err := o.AnchorGraph.GetDidAnchors(cidWithHint, suffix)
	if err != nil {
		logger.Warnf("process did failed for did[%s]: %s", did, err.Error())

		return err
	}

	logger.Debugf("got %d anchors for out-of-system did[%s]", len(anchors), did)

	for _, anchor := range anchors {
		logger.Debugf("processing anchor[%s] for out-of-system did[%s]", anchor.CID, did)

		if err := o.processAnchor(&anchorinfo.AnchorInfo{Hashlink: anchor.CID},
			anchor.Info, suffix); err != nil {
			if errors.IsTransient(err) {
				// Return an error so that the message is redelivered and retried.
				return fmt.Errorf("process out-of-system anchor [%s]: %w", anchor.CID, err)
			}

			logger.Warnf("ignoring anchor[%s] for did[%s]", anchor.CID, did, err.Error())

			continue
		}
	}

	return nil
}

func getDidParts(did string) (cid, suffix string, err error) {
	const delimiter = ":"

	pos := strings.LastIndex(did, delimiter)
	if pos == -1 {
		return "", "", fmt.Errorf("invalid number of parts for did[%s]", did)
	}

	return did[0:pos], did[pos+1:], nil
}

//nolint:funlen
func (o *Observer) processAnchor(anchor *anchorinfo.AnchorInfo, info *verifiable.Credential, suffixes ...string) error {
	logger.Debugf("processing anchor[%s] from [%s], suffixes: %s", anchor.Hashlink, anchor.AttributedTo, suffixes)

	anchorPayload, err := util.GetAnchorSubject(info)
	if err != nil {
		return fmt.Errorf("failed to extract anchor payload from anchor[%s]: %w", anchor.Hashlink, err)
	}

	pc, err := o.ProtocolClientProvider.ForNamespace(anchorPayload.Namespace)
	if err != nil {
		return fmt.Errorf("failed to get protocol client for namespace [%s]: %w", anchorPayload.Namespace, err)
	}

	v, err := pc.Get(anchorPayload.Version)
	if err != nil {
		return fmt.Errorf("failed to get protocol version for transaction time [%d]: %w",
			anchorPayload.Version, err)
	}

	ad := &util.AnchorData{OperationCount: anchorPayload.OperationCount, CoreIndexFileURI: anchorPayload.CoreIndex}

	canonicalID, err := hashlink.GetResourceHashFromHashLink(anchor.Hashlink)
	if err != nil {
		return fmt.Errorf("failed to get canonical ID from hl[%s]: %w", anchor.Hashlink, err)
	}

	equivalentRefs := []string{anchor.Hashlink}
	if o.discoveryDomain != "" {
		// only makes sense to have discovery domain with webcas (may change with ipfs gateway requirements)
		equivalentRefs = append(equivalentRefs, "https:"+o.discoveryDomain+":"+canonicalID)
	}

	sidetreeTxn := txnapi.SidetreeTxn{
		TransactionTime:      uint64(info.Issued.Unix()),
		AnchorString:         ad.GetAnchorString(),
		Namespace:            anchorPayload.Namespace,
		ProtocolGenesisTime:  anchorPayload.Version,
		CanonicalReference:   canonicalID,
		EquivalentReferences: equivalentRefs,
	}

	logger.Debugf("processing anchor[%s], core index[%s]", anchor.Hashlink, anchorPayload.CoreIndex)

	err = v.TransactionProcessor().Process(sidetreeTxn, suffixes...)
	if err != nil {
		return fmt.Errorf("failed to processAnchors core index[%s]: %w", anchorPayload.CoreIndex, err)
	}

	// update global did/anchor references
	acSuffixes := getKeys(anchorPayload.PreviousAnchors)

	err = o.DidAnchors.PutBulk(acSuffixes, anchor.Hashlink)
	if err != nil {
		return fmt.Errorf("failed updating did anchor references for anchor credential[%s]: %w", anchor.Hashlink, err)
	}

	logger.Infof("Successfully processed %d DIDs in anchor[%s], core index[%s]",
		len(anchorPayload.PreviousAnchors), anchor.Hashlink, anchorPayload.CoreIndex)

	// Post a 'Like' activity to the originator of the anchor credential.
	err = o.postLikeActivity(anchor)
	if err != nil {
		// This is not a critical error. We have already processed the anchor, so we don't want
		// to trigger a retry by returning a transient error. Just log a warning.
		logger.Warnf("A 'Like' activity could not be posted to the outbox: %s", err)
	}

	return nil
}

func (o *Observer) postLikeActivity(anchor *anchorinfo.AnchorInfo) error {
	if anchor.AttributedTo == "" {
		logger.Debugf("Not posting 'Like' activity since no attributedTo ID was specified for anchor [%s]",
			anchor.Hashlink)

		return nil
	}

	logger.Debugf("Posting 'Like' activity to [%s] for anchor [%s]",
		anchor.AttributedTo, anchor.Hashlink)

	refURL, err := url.Parse(anchor.Hashlink)
	if err != nil {
		return fmt.Errorf("parse hash link [%s]: %w", anchor.Hashlink, err)
	}

	attributedTo, err := url.Parse(anchor.AttributedTo)
	if err != nil {
		return fmt.Errorf("parse origin [%s]: %w", anchor.AttributedTo, err)
	}

	publishedTime := time.Now()

	var result *vocab.ObjectProperty

	if anchor.LocalHashlink != "" {
		u, e := url.Parse(anchor.LocalHashlink)
		if e != nil {
			return fmt.Errorf("parse local hashlink [%s]: %w", anchor.LocalHashlink, e)
		}

		result = vocab.NewObjectProperty(vocab.WithAnchorReference(
			vocab.NewAnchorReferenceWithOpts(
				vocab.WithURL(u),
			)),
		)
	}

	like := vocab.NewLikeActivity(
		vocab.NewObjectProperty(vocab.WithAnchorReference(
			vocab.NewAnchorReferenceWithOpts(
				vocab.WithURL(refURL),
			)),
		),
		vocab.WithTo(attributedTo, vocab.PublicIRI),
		vocab.WithPublishedTime(&publishedTime),
		vocab.WithResult(result),
	)

	_, err = o.Outbox().Post(like)
	if err != nil {
		return fmt.Errorf("post like: %w", err)
	}

	return nil
}

func getKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}
