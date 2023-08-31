package tinyCache

// 选择相应的节点
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// 获取缓存值
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
