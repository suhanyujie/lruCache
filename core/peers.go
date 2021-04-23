package core

type PeerPicker interface {
	// 通过 key 获取访问节点
	PickPeer(key string) (peer PeerGetter, ok bool)
}

type PeerGetter interface {
	// 用于从节点（group）查找缓存值
	Get(group string, key string) ([]byte, error)
}