package core

import (
	"container/list"
)

/// LRU 实现
/// 参考 https://geektutu.com/post/geecache-day1.html

type Lru struct {
	maxBytes int64
	nbytes int64
	ll *list.List
	cacheMap map[string]*list.Element
	OnEvicted func(key string, value Value)
}

type Value interface {
	Len() int
}

type entry struct {
	key string
	value Value
}

func NewLru(maxBytes int64, onEvicted func(string, Value)) *Lru {
	return &Lru{
		maxBytes:  maxBytes,
		nbytes:    0,
		ll:        list.New(),
		cacheMap:  make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (lru *Lru) Get(key string) (value Value, ok bool) {
	if ele, ok := lru.cacheMap[key]; ok {
		lru.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// 删除链表中的元素，删除 map 中的键值对，重新计算 nbytes 值。
func (lru *Lru) Del(key string) (ok bool) {
	if ele, ok := lru.cacheMap[key]; ok {
		lru.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(lru.cacheMap, kv.key)
		lru.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if lru.OnEvicted != nil {
			lru.OnEvicted(kv.key, kv.value)
		}
	}
	ok = true
	return
}

func (lru *Lru) RemoveOldest() {
	ele := lru.ll.Back()
	if ele != nil {
		lru.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(lru.cacheMap, kv.key)
		lru.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if lru.OnEvicted != nil {
			lru.OnEvicted(kv.key, kv.value)
		}
	}
}

func (lru *Lru) AddOrUpdate(key string, value Value) {
	if ele, ok := lru.cacheMap[key]; ok {
		lru.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		// 键不变，但值可能发生变化，所以需要更新大小。
		lru.nbytes += int64(value.Len() - kv.value.Len())
		kv.value = value
	} else {
		ele := lru.ll.PushFront(&entry{
			key, value,
		})
		lru.cacheMap[key] = ele
		lru.nbytes += int64(len(key) + value.Len())
	}
	for lru.maxBytes != 0 && lru.nbytes > lru.maxBytes {
		lru.RemoveOldest()
	}
}

func (lru *Lru) Len() int {
	return lru.ll.Len()
}

