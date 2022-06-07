package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu sync.Mutex

	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if item, exist := l.items[key]; exist {
		l.queue.MoveToFront(item)
		item.Value = cacheItem{key: key, value: value}

		return true
	}

	l.items[key] = l.queue.PushFront(cacheItem{key: key, value: value})

	if l.queue.Len() > l.capacity {
		tail := l.queue.Back()

		l.queue.Remove(tail)
		delete(l.items, tail.Value.(cacheItem).key)
	}

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if item, exist := l.items[key]; exist {
		l.queue.MoveToFront(item)

		return item.Value.(cacheItem).value, true
	}

	return nil, false
}

func (l *lruCache) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}
