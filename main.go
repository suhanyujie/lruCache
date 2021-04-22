package main

import (
	"fmt"
	"github.com/suhanyujie/lruCache/core"
	"log"
	"net/http"
)

var (
	db = map[string]string{
		"Stu1": "zhangSan",
		"Stu2": "zhangSan2",
		"Stu3": "zhangSan3",
	}
)

func main() {
	core.NewGroup("stuTest", 1 << 10, core.GetterFunc(func (key string) ([]byte, error) {
		log.Println("[cacheDb] search key", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("key %s not exist. ", key)
	}))
	addr := ":9999"
	peers := core.NewHttpPool(addr)
	log.Printf("lruCache is running at %s. \n", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}


