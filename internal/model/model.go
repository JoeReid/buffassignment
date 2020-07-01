package model

import "errors"

// ErrNotFound should be returned by store implementations when they
// couldn't find the requested data
var ErrNotFound = errors.New("the requested data was not found in the store")

// Store defines all the actions needed to implement the full storage layer
// This could be implemented by:
//   - A relational database (for production)
//   - A mock implementation (for testing)
//   - An RPC backend (for unforeseen future developments)
//
// Genericising the storage actions in this way makes the code considerably
// easier to re-factor with respect to storage sub-systems, should they need to change
type Store interface {
	VideoStreamStore
	BuffStore
}
