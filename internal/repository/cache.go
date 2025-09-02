package repository

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
)

type CacheItem struct {
	Value      interface{}
	Expiration time.Time
}

type CacheRepository struct {
	mu    sync.RWMutex
	items map[string]CacheItem
}

func NewCacheRepository() *CacheRepository {
	cache := &CacheRepository{
		items: make(map[string]CacheItem),
	}

	go cache.cleanup()

	return cache
}

func (c *CacheRepository) Set(key string, value interface{}, expiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheItem{
		Value:      value,
		Expiration: time.Now().Add(expiration),
	}

	return nil
}

func (c *CacheRepository) Get(key string, dest interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return errors.New("key not found")
	}

	if time.Now().After(item.Expiration) {
		delete(c.items, key)
		return errors.New("key expired")
	}

	jsonBytes, err := json.Marshal(item.Value)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonBytes, dest)
}

func (c *CacheRepository) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
	return nil
}

func (c *CacheRepository) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]CacheItem)
	return nil
}

func (c *CacheRepository) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.Expiration) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}
