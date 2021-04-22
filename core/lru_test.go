package lru

import (
	"encoding/json"
	"errors"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestLruGet(t *testing.T) {
	lru := NewLru(1000, nil)
	lru.AddOrUpdate("name1", String("liudahua"))
	lru.AddOrUpdate("name2", String("wangchuanguang"))
	lru.AddOrUpdate("name3", String("xujingzhong"))
	lru.Get("name2")
	lru.Get("name3")
	lru.RemoveOldest()
	// name1 被删掉了，所以不应该存在
	if _, ok := lru.Get("name1"); ok {
		t.Error(errors.New("error 1"))
	}
	t.Log("end...")
}

func TestLruOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := NewLru(1000, callback)
	lru.AddOrUpdate("name1", String("liudahua"))
	lru.AddOrUpdate("name2", String("wangchuanguang"))
	lru.AddOrUpdate("name3", String("xujingzhong"))
	lru.RemoveOldest()
	json1, _ := json.Marshal(keys)
	t.Log(string(json1))
}
