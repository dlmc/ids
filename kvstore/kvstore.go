// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kvstore


import (
	"github.com/dlmc/ids/internal/set"
	"github.com/dlmc/ids/internal/oset"
	"github.com/dlmc/ids/global"
)


type KvStore map[string]global.IKvStore

func New() KvStore {
	return KvStore{}
}

// KeyType is ignored here, all uses string key
func (s KvStore) CreateStore(name string, st global.StoreType, kt global.KeyType) bool{
	if _, ok := s[name];  ok {
		return false
	}
	switch st {
	case global.STSet:
		s[name] = set.New(name, global.KTStr)
	case global.STOSet:
		fallthrough 
	default:
		s[name] = oset.New(name, global.KTStr)
	}
	return true
}

func (s KvStore) DeleteStore(name string) {
	t := s[name]
	delete(s, name)	 
	t.Clear()
}

func (s KvStore) GetStore(name string) (store global.IKvStore, found bool) {
	store, found = s[name]
	return
}
