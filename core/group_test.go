package lru

import (
	"fmt"
	"log"
	"testing"
)

var (
	db = map[string]string{
		"Liudehua": "001-Liudehua",
		"Sunwenyu": "002-Sunwenyu",
		"Zhangsan": "003-Zhangsan",
	}
)

func TestGroup1(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	lc := NewGroup("nameInfo", 1<<10, GetterFunc(func(key string) ([]byte, error) {
		log.Println("[lc] search by key" )
		if v, ok := db[key]; ok {
			if _, ok2 := loadCounts[key]; !ok2 {
				loadCounts[key] = 0
			}
			loadCounts[key] += 1
			return []byte(v), nil
		}
		return nil, fmt.Errorf("key %s not exist", key)
	}))

	for k, v := range db {
		if view, err := lc.Get(k); err != nil || view.String() != v {
			t.Fatalf("failed to get data for key %s", k)
		}
		if _, err := lc.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache miss for key %s", k)
		}
	}

	if view, err := lc.Get("unknown"); err == nil {
		t.Fatalf("the value of unknown should be empty, but got value: %s", view)
	}
}
