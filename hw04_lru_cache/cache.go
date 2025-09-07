package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	value    interface{}
	listItem *ListItem
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*cacheItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*cacheItem, capacity),
	}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	if element, ok := l.items[key]; ok {
		element.value = value
		l.queue.MoveToFront(element.listItem)
		return true
	}

	newItem := l.queue.PushFront(key)
	l.items[key] = &cacheItem{
		value:    value,
		listItem: newItem,
	}

	if l.queue.Len() > l.capacity {
		back := l.queue.Back()
		if back != nil {
			deleteKey := back.Value.(Key)
			l.queue.Remove(back)
			delete(l.items, deleteKey)
		}
	}

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	if element, ok := l.items[key]; ok {
		l.queue.MoveToFront(element.listItem)
		return element.value, true
	}

	return nil, false
}

func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*cacheItem, l.capacity)
}
