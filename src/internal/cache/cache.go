// Package cache contains errors common to every cache provider.
package cache

import "errors"

// ErrKeyNotExists represents error when key is not present in cache.
var ErrKeyNotExists = errors.New("key is not found")
