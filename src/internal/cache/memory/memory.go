// Package memory provides cache implementation and api that is safe to use in concurrent applications.
package memory

import (
	"context"
	"log/slog"
	cachepkg "mortgage-calculator/src/internal/cache"
	"sync"
	"time"
)

// Cache stores cached data.
type Cache struct {
	log    *slog.Logger
	data   cache
	ttl    int64
	lastID int64
	mu     sync.RWMutex
}

type cache map[string]cacheItem

type cacheItem struct {
	val []byte
	id  int64
	exp int64
}

// New is a constructor for Cache.
func New(log *slog.Logger, ttl int64) *Cache {
	data := make(cache)

	return &Cache{
		log:  log,
		ttl:  ttl,
		data: data,
	}
}

// Clear checks whether entries are expired and deletes them if true.
func (c *Cache) Clear(_ context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().Unix()
	for key, item := range c.data {
		if now > item.exp {
			delete(c.data, key)
		}
	}
}

// Get returns value by key if latter exists else ErrKeyNotExists.
func (c *Cache) Get(
	_ context.Context,
	key string,
) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.data[key]

	if !ok || time.Now().Unix() > item.exp {
		delete(c.data, key)
		return nil, cachepkg.ErrKeyNotExists
	}

	return item.val, nil
}

// Set saves given value by given key and sets expiration time.
func (c *Cache) Set(
	_ context.Context,
	key string,
	value []byte,
) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.lastID++
	c.data[key] = cacheItem{
		id:  c.lastID,
		val: value,
		exp: time.Now().Unix() + c.ttl,
	}

	return nil
}

// List returns all active cache entries.
func (c *Cache) List(_ context.Context) ([]*cachepkg.Entry, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	res := make([]*cachepkg.Entry, 0)
	for k, v := range c.data {
		res = append(res, &cachepkg.Entry{
			ID:  v.id,
			Key: k,
			Val: v.val,
		})
	}

	return res, nil
}
