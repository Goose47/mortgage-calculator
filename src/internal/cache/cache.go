// Package cache contains errors and types common to every cache provider.
package cache

import "errors"

// ErrKeyNotExists represents error when key is not present in cache.
var ErrKeyNotExists = errors.New("key is not found")

// Entry describes data stored in cache.
type Entry struct {
	Key string
	Val []byte
	ID  int64
}
