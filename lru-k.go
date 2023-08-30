package tinyCache

import "container/list"

type historyCounter struct {
	key   string
	count int
}

type Lru_kCache struct {
	k            int
	rest         int
	historyRest  int
	historyL     *list.List
	historyCache map[string]*list.Element
	ll           *list.List
	Cache        map[string]*list.Element
}
