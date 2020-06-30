package model

import "errors"

var ErrNotFound = errors.New("the requested data was not found in the store")

type Store interface {
	VideoStreamStore
	BuffStore
}
