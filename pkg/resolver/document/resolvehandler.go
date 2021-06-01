/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package document

import (
	"fmt"
	"strings"

	"github.com/trustbloc/edge-core/pkg/log"
	"github.com/trustbloc/sidetree-core-go/pkg/document"
	"github.com/trustbloc/sidetree-core-go/pkg/restapi/dochandler"
)

var logger = log.New("orb-resolver")

const (
	delimiter         = ":"
	minOrbSuffixParts = 2
)

// ResolveHandler resolves generic documents.
type ResolveHandler struct {
	coreResolver dochandler.Resolver
	discovery    discovery

	namespace           string
	aliases             []string
	unpublishedDIDLabel string
}

// did discovery service.
type discovery interface {
	RequestDiscovery(id string) error
}

// NewResolveHandler returns a new document resolve handler.
func NewResolveHandler(namespace string, aliases []string, unpublishedDIDLabel string, resolver dochandler.Resolver, discovery discovery) *ResolveHandler { //nolint:lll
	return &ResolveHandler{
		namespace:           namespace,
		aliases:             aliases,
		unpublishedDIDLabel: unpublishedDIDLabel,
		coreResolver:        resolver,
		discovery:           discovery,
	}
}

// ResolveDocument resolves a document.
func (r *ResolveHandler) ResolveDocument(id string) (*document.ResolutionResult, error) {
	response, err := r.coreResolver.ResolveDocument(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			// this can only happen if short form resolution returned not found
			r.requestDiscovery(id)
		}

		return nil, err
	}

	return response, nil
}

func (r *ResolveHandler) requestDiscovery(id string) {
	// it only makes sense to request discovery if did has cid
	orbSuffix, err := r.getOrbSuffix(id)
	if err != nil {
		// not proper orb suffix - nothing to do
		logger.Debugf("get orb suffix from id[%s] error: %s", id, err.Error())

		return
	}

	if strings.HasPrefix(orbSuffix, r.unpublishedDIDLabel) {
		// we cannot request discovery for unpublished DIDs - nothing to do
		return
	}

	logger.Infof("requesting discovery for orb suffix[%s]", orbSuffix)

	err = r.discovery.RequestDiscovery(orbSuffix)
	if err != nil {
		logger.Warnf("error while requesting discovery for orb suffix[%s]: %s", orbSuffix, err.Error())
	}
}

func (r *ResolveHandler) getNamespace(shortFormDID string) (string, error) {
	if strings.HasPrefix(shortFormDID, r.namespace+delimiter) {
		return r.namespace, nil
	}

	// check aliases
	for _, ns := range r.aliases {
		if strings.HasPrefix(shortFormDID, ns+delimiter) {
			return ns, nil
		}
	}

	return "", fmt.Errorf("did must start with configured namespace[%s] or aliases%v", r.namespace, r.aliases)
}

// getOrbSuffix fetches unique portion of ID which is string after namespace.
// Valid Orb suffix has two parts cas hint + cid:suffix.
func (r *ResolveHandler) getOrbSuffix(shortFormDID string) (string, error) {
	namespace, err := r.getNamespace(shortFormDID)
	if err != nil {
		return "", err
	}

	orbSuffix := shortFormDID[len(namespace+delimiter):]

	parts := strings.Split(orbSuffix, delimiter)
	if len(parts) < minOrbSuffixParts {
		return "", fmt.Errorf("invalid number of parts for orb suffix[%s]", orbSuffix)
	}

	return orbSuffix, nil
}
