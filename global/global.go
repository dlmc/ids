
// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package global

import (
	"strconv"
	"strings"
	"errors"
)


type KeyType uint8
const (
	KTStr = iota					// using string key
	KTInt						// using int key
	KTFloat64					// using float64 key
)
func GetKeyType(str string) (kt KeyType, err error) {
	switch str {
	case "string":
		kt = KTStr
	case "integer":
		kt = KTInt
	case "float":
		kt = KTFloat64
	default:
		err = errors.New("keytype [" + str + "] not in [string, integer, float]")
	}
	return
}

type StoreType uint8
const (
	STOSet = iota				// using OSet Store 
	STSet						// using Set Store
	STDeque						// useing Deque Store
)
func GetStoreType(str string) (st StoreType, err error) {
	switch str {
	case "oset":
		st = STOSet
	case "set":
		st = STSet
	default:
		err = errors.New("storetype [" + str + "] not in [oset, set]")
	}
	return
}

// Int2StrPadZero converts integer to string 
// Adding prefix "0" only if the length of the converted string is less than digits number
func Int2StrPadZero(i, digits int) string {
	s := strconv.Itoa(i)
	n := 1+digits-strings.Count(s,"")
	if n > 0 {
		return strings.Repeat("0", n) + s
	}
	return s
}


// About 3x faster then strconv.ParseInt because does not check for range error and support only base 10, which is enough for JSON
// credit goes to github.com/buger/jsonparser
func ParseInt(bytes []byte) (v int64, ok bool) {
	if len(bytes) == 0 {
		return 0, false
	}

	var neg bool = false
	if bytes[0] == '-' {
		neg = true
		bytes = bytes[1:]
	}

	for _, c := range bytes {
		if c >= '0' && c <= '9' {
			v = (10 * v) + int64(c-'0')
		} else {
			return 0, false
		}
	}

	if neg {
		return -v, true
	} else {
		return v, true
	}
}


// IKvStore defines the interface that all the data store implementation shall comply
// For the data store that can not support perticular interface method, create a dummy
// supported Store: OSet Set
type IKvStore interface {
	Create(key, value interface{}) bool
	Read(key interface{}) (value interface{}, found bool)
	Update(key, value interface{}) bool
	Delete(Key interface{}) bool
	Clear()
	Size() int
	Values() []interface{}
	Keys() []interface{} 
	RangeGet(startkey, endkey interface{}, limit int, ascending bool) (values []interface{}, count int)
	// Todo: add block CRUD operations: 
	// BlockCreate, BLockRead, BlockUpdate, BlockDelete???
	// BlockCreateUpdate ???
}

// IDataStore defines the interface that all the data store implementation shall comply
// For the data store that can not support perticular interface method, create a dummy
// Supported Store: Deque
type IDequeStore interface {
	PushFront(interface{})
	PushBack(interface{})
	PopFront() (interface{}, bool)
	PopBack() (interface{}, bool) 
	Size() int
	Clear()
}
