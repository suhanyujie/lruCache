package main

import (
	"flag"
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

var (
	port int
	api bool
	apiAddr = "http://localhost:9001"
)

func main() {
	flag.IntVar(&port, "port", 8081, "lruCache server port")
	flag.BoolVar(&api, "api", false, "start a api server or not")
	flag.Parse()

	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}
	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}
	lru := createGroup()
	if api {
		go startApiServer(apiAddr, lru)
	}

	startSomeLruCacheServer(addrMap[port], addrs, lru)
}

func startSomeLruCacheServer(addr string, addrs []string, g *core.Group) {
	peers := core.NewHttpPool(addr)
	peers.Set(addrs...)
	g.RegisterPeers(peers)
	log.Println("lruCache is running at: ", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func createGroup() *core.Group {
	return core.NewGroup("someCache", 1 << 10, core.GetterFunc(func(key string) ([]byte, error) {
		log.Println("[lruCacheDb] search key: ", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist. ", key)
	}))
}

// 创建与用户交互的 api 服务器
func startApiServer(addr string, lru *core.Group) {
	http.Handle("/api", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		view, err := lru.Get(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(view.ByteSlice())
	}))
	log.Printf("frontend api server is running at: %s", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}




