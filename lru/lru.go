package lru

import "container/list"

// Cache map存储的是key和value的映射，list是用来调整最近是否使用过，通过list还存储了key，方便删除
type Cache struct {
	maxBytes  int64
	nBytes    int64
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}
type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		// 移动到队头
		c.ll.MoveToFront(ele)
		res := ele.Value.(*entry)
		return res.value, ok
	}
	return
}

// RemoveOldEle 删除最久未被使用的
func (c *Cache) RemoveOldEle() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		// 删除后当前的bytes要减小
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		// 调用回调函数
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes -= int64(kv.value.Len()) - int64(value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{
			key:   key,
			value: value,
		})
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}

	// 如果超过内存了，就要删
	for c.maxBytes != 0 && c.nBytes > c.maxBytes {
		c.RemoveOldEle()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
