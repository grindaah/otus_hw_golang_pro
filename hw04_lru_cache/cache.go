package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	sync.Mutex
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

func (lru *lruCache) Set(key Key, value interface{}) bool {
	lru.Lock()
	defer lru.Unlock()
	if li, ok := lru.items[key]; ok {
		li.Value = cacheItem{key: key, value: value}
		lru.queue.MoveToFront(li)
		return true
	} else {
		if lru.queue.Len() == lru.capacity {
			// remove last elem
			lastKey := lru.queue.Back().Value.(cacheItem).key
			lastElem := lru.items[lastKey]
			delete(lru.items, lastKey)
			lru.queue.Remove(lastElem)

			// append new
			li := lru.queue.PushFront(cacheItem{key: key, value: value})
			lru.items[key] = li
		} else {
			li := lru.queue.PushFront(cacheItem{key: key, value: value})
			lru.items[key] = li
		}
		return false
	}
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	lru.Lock()
	defer lru.Unlock()

	if li, ok := lru.items[key]; ok {
		lru.queue.MoveToFront(li)
		return li.Value.(cacheItem).value, true
	}
	return nil, false
}

func (lru *lruCache) Clear() {
	lru.Lock()
	defer lru.Unlock()

	lru.queue = NewList()
	lru.items = make(map[Key]*ListItem, lru.capacity)
}
