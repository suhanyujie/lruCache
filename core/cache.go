package core

import "sync"

/// 基于 lru 封装缓存

// 缓存值 b。之所以使用 []byte 是为了可以存储任意类型的值。
type ByteView struct {
	b []byte
}

func (v ByteView) Len() int {
	return len(v.b)
}

func (v ByteView) ByteSlice() []byte {
	newV := make([]byte, v.Len())
	copy(newV, v.b)
	return newV
}

func (v ByteView) String() string {
	return string(v.b)
}

type cache struct {
	mu sync.Mutex
	lru *Lru
	cacheBytes int64
}

func NewCache(maxBytes int64) *cache {
	return &cache{
		mu:         sync.Mutex{},
		lru:        NewLru(maxBytes, nil),
		cacheBytes: maxBytes,
	}
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = NewLru(c.cacheBytes, nil)
	}
	c.lru.AddOrUpdate(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = NewLru(c.cacheBytes, nil)
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}

func cloneBytes(b []byte) []byte {
	newByte := make([]byte, len(b))
	copy(newByte, b)
	return newByte
}

