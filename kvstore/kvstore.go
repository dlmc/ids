// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kvstore


import (
	"github.com/dlmc/ids/internal/set"
	"github.com/dlmc/ids/internal/oset"
	"github.com/dlmc/ids/global"
	"errors"
)


type KvStore map[string]global.IKvStore

func New() KvStore {
	return KvStore{}
}

func (s KvStore) CreateStore(name, stype, ktype string) error {
	var err error
	if _, ok := s[name];  ok {
		err = errors.New(global.StrStoreExist + name)
	} else {
		switch stype {
		case global.QueryStoreTypeSet:
			if st, e := set.New(name, ktype); e == nil {
				s[name] = st
			} else {
				err = e
			}
		case global.QueryStoreTypeOSet:
			if st, e := oset.New(name, ktype); e == nil {
				s[name] = st
			} else {
				err = e
			}
		default:
			err = errors.New(stype + global.StrStoreTypeNotExist)
		}		
	}
	return err
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
