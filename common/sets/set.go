package sets

import (
	"golang.org/x/exp/maps"
	"sync"
)

type Operator uint8

const (
	NoAction Operator = iota
	Delete
	Overwrite
)

type Set[K comparable, V any] struct {
	data map[K]V
	x    sync.RWMutex
}

func (s *Set[K, V]) init() {
	if s.data == nil {
		s.data = make(map[K]V)
	}
}

func (s *Set[K, V]) Set(k K, v V) (replaced bool) {

	s.x.Lock()
	defer s.x.Unlock()
	s.init()

	_, replaced = s.data[k]
	s.data[k] = v
	return
}

func (s *Set[K, V]) SetFn(k K, fn func(k K, v V, exist bool) (nv V, op Operator)) (nv V, op Operator) {

	s.x.Lock()
	defer s.x.Unlock()
	s.init()

	v, ok := s.data[k]
	nv, op = fn(k, v, ok)
	switch op {
	case Overwrite:
		s.data[k] = nv
	case Delete:
		delete(s.data, k)
	default:
	}
	return
}

func (s *Set[K, V]) SetX(k K, v V) (old V, ok bool) {

	s.x.Lock()
	defer s.x.Unlock()
	s.init()

	if old, ok = s.data[k]; ok {
		s.data[k] = v
	}

	return
}

func (s *Set[K, V]) SetNX(k K, v V) (ok bool) {

	s.x.Lock()
	defer s.x.Unlock()
	s.init()

	if _, x := s.data[k]; !x {
		s.data[k] = v
		ok = true
	}

	return
}

func (s *Set[K, V]) Get(k K) (v V, ok bool) {

	s.x.RLock()
	defer s.x.RUnlock()

	if s.data == nil {
		return
	}

	v, ok = s.data[k]
	return
}

func (s *Set[K, V]) Delete(k K) (old V, ok bool) {

	s.x.Lock()
	defer s.x.Unlock()

	if s.data == nil {
		return
	}

	old, ok = s.data[k]
	delete(s.data, k)
	return
}

func (s *Set[K, V]) Clear() (cnt int) {

	s.x.Lock()
	defer s.x.Unlock()

	if s.data == nil {
		return
	}

	cnt = len(s.data)
	maps.Clear(s.data)
	return
}

func (s *Set[K, V]) Each(handler func(k K, v V)) {

	var data map[K]V
	s.x.RLock()
	if len(s.data) > 0 {
		data = maps.Clone(s.data)
	}
	s.x.RUnlock()

	if len(data) < 1 {
		return
	}

	for k, v := range data {
		handler(k, v)
	}
}

func (s *Set[K, V]) EachForWrite(handler func(k K, v V)) {
	s.x.Lock()
	defer s.x.Unlock()
	for k, v := range s.data {
		handler(k, v)
	}
}
