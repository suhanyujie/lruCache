package core

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/_lruCache/"
	defaultReplicas = 50
)

type HttpPool struct {
	self string
	basePath string
	mu sync.Mutex
	peers *Map
	httpGetters map[string]*HttpGetter
}

func NewHttpPool(self string) *HttpPool {
	return &HttpPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// todo
func (_this *HttpPool) Log(format string, v ...interface{})  {
	log.Printf("[Server %s] %s", _this.self, fmt.Sprintf(format, v...))
}

// 实现 ServeHTTP
func (_this *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	if !strings.HasPrefix(r.URL.Path, _this.basePath) {
		panic("HttpPool serving unexpected path: " + r.URL.Path)
	}
	_this.Log("%s %s", r.Method, r.URL.Path)

	// /{basePath}/{groupName}/{key}
	parts := strings.SplitN(r.URL.Path[len(_this.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	key := parts[1]
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: " + groupName, http.StatusNotFound)
		return
	}
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 封装 body，构造请求
	body, err := proto.Marshal(&Response{Value: view.ByteSlice()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(body)
}
