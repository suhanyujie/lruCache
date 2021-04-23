package core

import (
	"hash/crc32"
	"sort"
	"strconv"
)

/// 一致性哈希算法
/// 主要通过计算将一个节点映射为多个哈希值，对其排序，然后放到环上，使多个虚拟节点可以首尾相接。
/// 当需要通过接口获取一个 key 的缓存值时，通过 key 从上述的节点结构中获取对应的节点，从而进一步获取缓存值。

type Hash func(data []byte) uint32

type Map struct {
	hash Hash
	// 虚拟节点倍数。如倍数为 5，在新增1个真实节点时，需要映射为 5 个虚拟节点。
	replicas int
	// 哈希环。存放的是虚拟节点的哈希值
	keys []int  // 有序的
	// 虚拟节点哈希值和真实节点名称的映射
	// 一个真实节点对应多个虚拟节点
	hashMap map[int]string
}

func NewMap(replicas int, fn Hash) *Map {
	m := &Map{
		hash:     fn,
		replicas: replicas,
		keys:     nil,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// 新增真实节点。
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// 获取节点
func (m *Map) Get(key string) string {
	if key == "" {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	// 如果 idx 为 len(m.keys)，因为是环状结构，所以表示是 `m.keys[0]` 的元素
	keysIndex := idx % len(m.keys)
	return m.hashMap[m.keys[keysIndex]]
}