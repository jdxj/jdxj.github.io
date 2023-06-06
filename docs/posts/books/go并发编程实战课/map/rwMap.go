package test

import "sync"

type RWMap[K comparable, V any] struct {
	sync.RWMutex
	m map[K]V
}

func NewRWMap[K comparable, V any](n int) *RWMap[K, V] {
	return &RWMap[K, V]{
		m: make(map[K]V),
	}
}

func (m *RWMap[K, V]) Get(k K) (V, bool) {
	m.RLock()
	defer m.RUnlock()
	v, existed := m.m[k]
	return v, existed
}

func (m *RWMap[K, V]) Set(k K, v V) {
	m.Lock()
	defer m.Unlock()
	m.m[k] = v
}

func (m *RWMap[K, V]) Delete(k K) {
	m.Lock()
	defer m.Unlock()
	delete(m.m, k)
}

func (m *RWMap[K, V]) Len() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.m)
}

func (m *RWMap[K, V]) Each(f func(k K, v V) bool) {
	m.RLock()
	defer m.RUnlock()

	for k, v := range m.m {
		if !f(k, v) {
			return
		}
	}
}
