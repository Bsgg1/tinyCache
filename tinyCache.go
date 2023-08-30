package tinyCache

import (
	"fmt"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	getter    Getter //缓存未命中时的回调
	name      string
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		getter:    getter,
		name:      name,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	if value, ok := g.mainCache.Get(key); ok {
		return value, nil
	}
	// 去其他地方找key
	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	return g.GetLocally(key)
}

func (g *Group) GetLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: bytes}
	g.Add(key, value)
	return value, nil
}

func (g *Group) Add(key string, value ByteView) {
	g.mainCache.Add(key, value)
}
