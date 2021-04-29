package core

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"net/http"
	"net/url"
)

type HttpGetter struct {
	baseUrl string
}

func (h *HttpGetter) Get(in *Request, out *Response) ([]byte, error) {
	apiUrl := fmt.Sprintf(
		"%v%v/%v",
		h.baseUrl,
		url.QueryEscape(in.GetGroup()),
		url.QueryEscape(in.GetKey()),
		)
	resp, err := http.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned error: %v", resp.Status)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}
	if err = proto.Unmarshal(bytes, out); err != nil {
		return bytes, fmt.Errorf("decoding response error: %v", err)
	}
	return bytes, nil
}

var _ PeerGetter = (*HttpGetter)(nil)

/// 实现 PeerPicker

func (p *HttpPool) Set(peers ...string)  {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = NewMap(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*HttpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &HttpGetter{
			baseUrl: peer + p.basePath,
		}
	}
}

func (p *HttpPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

var _ PeerPicker = (*HttpPool)(nil)
