/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package unpublished

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/aries-framework-go/spi/storage"
	"github.com/trustbloc/edge-core/pkg/log"
	"github.com/trustbloc/sidetree-core-go/pkg/api/operation"

	orberrors "github.com/trustbloc/orb/pkg/errors"
)

const nameSpace = "unpublished-operation"

var logger = log.New("unpublished-operation-store")

// New returns new instance of unpublished operation store.
func New(provider storage.Provider) (*Store, error) {
	store, err := provider.OpenStore(nameSpace)
	if err != nil {
		return nil, fmt.Errorf("failed to open unpublished operation store: %w", err)
	}

	return &Store{
		store: store,
	}, nil
}

// Store implements storage for unpublished operation.
type Store struct {
	store storage.Store
}

// Put saves an unpublished operation. If it it already exists an error will be returned.
func (s *Store) Put(op *operation.AnchoredOperation) error {
	if op.UniqueSuffix == "" {
		return fmt.Errorf("failed to save unpublished operation: suffix is empty")
	}

	_, err := s.Get(op.UniqueSuffix)
	if err == nil {
		return fmt.Errorf("pending operation found for suffix[%s], please re-submit your operation request at later time", op.UniqueSuffix) //nolint:lll
	}

	if !errors.Is(err, storage.ErrDataNotFound) {
		return fmt.Errorf("unable to check for pending operations for suffix[%s], please re-submit your operation request at later time: %w", op.UniqueSuffix, err) //nolint:lll
	}

	opBytes, err := json.Marshal(op)
	if err != nil {
		return fmt.Errorf("failed to marshal unpublished operation: %w", err)
	}

	logger.Debugf("storing unpublished '%s' operation for suffix[%s]: %s", op.Type, op.UniqueSuffix, string(opBytes))

	if e := s.store.Put(op.UniqueSuffix, opBytes); e != nil {
		return fmt.Errorf("failed to put unpublished operation for suffix[%s]: %w", op.UniqueSuffix, e)
	}

	return nil
}

// Get retrieves unpublished operation by suffix.
func (s *Store) Get(suffix string) (*operation.AnchoredOperation, error) {
	opBytes, err := s.store.Get(suffix)
	if err != nil {
		return nil, err
	}

	var op operation.AnchoredOperation

	err = json.Unmarshal(opBytes, &op)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal unpublished operation for suffix[%s]: %w", suffix, err)
	}

	logger.Debugf("retrieved unpublished '%s' operation for suffix[%s]: %s", op.Type, suffix, string(opBytes))

	return &op, nil
}

// Delete will delete unpublished operation for suffix.
func (s *Store) Delete(suffix string) error {
	if err := s.store.Delete(suffix); err != nil {
		return fmt.Errorf("failed to delete unpublished operation for suffix[%s]: %w", suffix, err)
	}

	return nil
}

// DeleteAll deletes all operations for suffixes.
func (s *Store) DeleteAll(suffixes []string) error {
	if len(suffixes) == 0 {
		return nil
	}

	operations := make([]storage.Operation, len(suffixes))

	for i, k := range suffixes {
		operations[i] = storage.Operation{Key: k}
	}

	err := s.store.Batch(operations)
	if err != nil {
		return orberrors.NewTransient(fmt.Errorf("failed to delete unpublished operations: %w", err))
	}

	logger.Debugf("deleted unpublished operations for %d suffixes", len(suffixes))

	return nil
}
