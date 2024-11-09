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
	ttl    int64
	lastId int64
	data   cache
	sync.RWMutex
}

type cache map[string]cacheItem

type cacheItem struct {
	id  int64
	val []byte
	exp int64
}

// New is a constructor for Cache.
func New(log *slog.Logger, TTL int64) *Cache {
	data := make(cache)

	return &Cache{
		log:  log,
		ttl:  TTL,
		data: data,
	}
}

// Clear checks whether entries are expired and deletes them if true.
func (c *Cache) Clear(ctx context.Context) {
	c.Lock()
	defer c.Unlock()

	now := time.Now().Unix()
	for key, item := range c.data {
		if now > item.exp {
			delete(c.data, key)
		}
	}
}

// Get returns value by key if latter exists else ErrKeyNotExists.
func (c *Cache) Get(
	ctx context.Context,
	key string,
) ([]byte, error) {
	c.RLock()
	defer c.RUnlock()

	item, ok := c.data[key]

	if !ok || time.Now().Unix() > item.exp {
		delete(c.data, key)
		return nil, cachepkg.ErrKeyNotExists
	}

	return item.val, nil
}

// Set saves given value by given key and sets expiration time.
func (c *Cache) Set(
	ctx context.Context,
	key string,
	value []byte,
) error {
	c.Lock()
	defer c.Unlock()

	c.lastId++
	c.data[key] = cacheItem{
		id:  c.lastId,
		val: value,
		exp: time.Now().Unix() + c.ttl,
	}

	return nil
}
