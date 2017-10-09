// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dqstore


import (
	"github.com/dlmc/ids/internal/deque"
	"github.com/dlmc/ids/global"
 	"sync"
 	"errors"
)

type DequeType uint8

type dequeue struct {
	dq *deque.Deque
	sync.Mutex
	KeyType global.Type
}


func NewDeque(dtype string) (*dequeue, error) {
	t, err := global.MapQueryType(dtype)
	return &dequeue{dq:deque.New(), KeyType:t}, err
}


func (d *dequeue) PushFront(i string) bool {
	d.Lock()
	defer d.Unlock()
	v, ok := global.ParseInput(i, d.KeyType)
	if ok {
		d.dq.PushFront(v)	
	}
	return ok
}
func (d *dequeue) PushBack(i string) bool {
	d.Lock()
	defer d.Unlock()
	v, ok := global.ParseInput(i, d.KeyType)
	if ok {
		d.dq.PushBack(v)	
	}
	return ok
}
func (d *dequeue) PopFront() (interface{}, bool) {
	d.Lock()
	defer d.Unlock()
	return d.dq.PopFront()
}
func (d *dequeue) PopBack() (interface{}, bool) {
	d.Lock()
	defer d.Unlock()
	return d.dq.PopBack()	
}
func (d *dequeue) Size() int {
	d.Lock()
	defer d.Unlock()
	return d.dq.Len()		
}
func (d *dequeue) Clear() {
	d.Lock()
	defer d.Unlock()
	d.dq.Clear()
}


type DequeStore map[string]*dequeue

func New() DequeStore {
	return DequeStore{}
}

func (s DequeStore) CreateStore(name, dtype string) error {
	if _, ok := s[name];  ok {
		return errors.New(global.StrStoreExist + name)
	}
	st, err := NewDeque(dtype)
	if err == nil {
		s[name] = st
	}
	return err
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
