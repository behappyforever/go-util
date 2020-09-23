package util

import (
	"container/list"
	"sync"
)

type Key string

type Value interface{}

type LruI interface {
	Put(k Key, v Value) Value
	Get(k Key) Value
	Rem(k Key) Value
}

type lru struct {
	*list.List
	mp    map[Key]*list.Element
	limit int
	sync.Mutex
}

type node struct {
	k Key
	v Value
}

// thread-safe lru cache, all op are O(1)
func NewLruCache(limit int, lazyInit bool) LruI {
	if limit <= 0 {
		return nil
	}

	cap := 0
	if !lazyInit {
		cap = limit
	}

	return &lru{
		limit: limit,
		mp:    make(map[Key]*list.Element, cap),
		List:  list.New(),
	}
}

func (l *lru) Put(k Key, v Value) Value {
	if l == nil {
		return nil
	}

	var old Value

	l.Lock()
	defer l.Unlock()

	if e, ok := l.mp[k]; ok {
		old = l.Remove(e).(node).v
	}

	if l.Len()+1 > l.limit {
		e := l.Back()
		k := e.Value.(node).k
		l.Remove(e)
		delete(l.mp, k)
	}
	l.mp[k] = l.PushFront(node{
		k: k,
		v: v,
	})

	return old
}

func (l *lru) Get(k Key) Value {
	if l == nil {
		return nil
	}

	l.Lock()
	defer l.Unlock()

	if e, ok := l.mp[k]; ok {
		l.MoveToFront(e)
		return e.Value.(node).v
	}

	return nil
}

func (l *lru) Rem(k Key) Value {
	if l == nil {
		return nil
	}

	l.Lock()
	defer l.Unlock()

	if e, ok := l.mp[k]; ok {
		delete(l.mp, k)
		return l.Remove(e).(node).v
	}

	return nil
}
