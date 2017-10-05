// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dqstore


import (
	"github.com/dlmc/ids/internal/deque"
	"github.com/dlmc/ids/global"
	"sync"
)

type dequeue struct {
	*deque.Deque
	sync.Mutex
}

func NewDeque() *dequeue {
	return &dequeue{Deque:deque.New()}
}

func (d *dequeue) Pushfront(i interface{}) {
	d.Lock()
	defer d.Unlock()
	d.Deque.PushFront(i)
}
func (d *dequeue) Pushback(i interface{}) {
	d.Lock()
	defer d.Unlock()
	d.Deque.PushBack(i)	
}
func (d *dequeue) Popfront() (interface{}, bool) {
	d.Lock()
	defer d.Unlock()
	return d.Deque.PopFront()	
}
func (d *dequeue) Popback() (interface{}, bool) {
	d.Lock()
	defer d.Unlock()
	return d.Deque.PopBack()	
}
func (d *dequeue) Size() int {
	d.Lock()
	defer d.Unlock()
	return d.Deque.Len()		
}
func (d *dequeue) Clear() {
	d.Lock()
	defer d.Unlock()
	d.Deque.Clear()
}


type DequeStore map[string]*dequeue

func New() DequeStore {
	return DequeStore{}
}

func (s DequeStore) CreateStore(name string) bool{
	if _, ok := s[name];  ok {
		return false
	}
	s[name] = NewDeque()
	return true
}

func (s DequeStore) DeleteStore(name string) {
	t := s[name]
	delete(s, name)	 
	t.Clear()
}

func (s DequeStore) GetStore(name string) (store global.IDequeStore, found bool) {
	store, found = s[name]
	return
}
