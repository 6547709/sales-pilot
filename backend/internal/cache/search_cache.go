package cache

import (
	"sync"
	"time"
)

// SearchCache 产品搜索简单 TTL 缓存
type SearchCache struct {
	mu    sync.RWMutex
	items map[string]entry
	ttl   time.Duration
}

type entry struct {
	body      []byte
	expiresAt time.Time
}

// NewSearchCache 默认 30 秒 TTL
func NewSearchCache(ttl time.Duration) *SearchCache {
	if ttl <= 0 {
		ttl = 30 * time.Second
	}
	return &SearchCache{items: make(map[string]entry), ttl: ttl}
}

// Get 命中返回 body 与 true
func (c *SearchCache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	e, ok := c.items[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(e.expiresAt) {
		return nil, false
	}
	return e.body, true
}

// Set 写入缓存
func (c *SearchCache) Set(key string, body []byte) {
	c.mu.Lock()
	c.items[key] = entry{body: append([]byte(nil), body...), expiresAt: time.Now().Add(c.ttl)}
	c.mu.Unlock()
}
