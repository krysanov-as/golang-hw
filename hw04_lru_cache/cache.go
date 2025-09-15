package hw04lrucache

import "sync"

type Key string

type cacheItem struct {
	key   Key
	value interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       sync.Mutex
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
	if element, ok := l.items[key]; ok {
		ci := cacheItem{key: key, value: value}
		element.Value = ci
		l.queue.MoveToFront(element)
		return true
	}

	if l.queue.Len() == l.capacity {
		back := l.queue.Back()
		val, _ := back.Value.(cacheItem)
		delete(l.items, val.key)
		l.queue.Remove(back)
	}

	ci := l.queue.PushFront(cacheItem{key: key, value: value})
	l.items[key] = ci

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if element, ok := l.items[key]; ok {
		l.queue.MoveToFront(element)
		val, _ := element.Value.(cacheItem)
		return val.value, true
	}

	return nil, false
}

func (l *lruCache) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}
