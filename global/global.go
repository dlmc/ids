
// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package global

import (
	"strconv"
	"strings"
	"unsafe"
	"errors"
)

const (
	StrStoreEmpty = "store empty - "
	StrStoreExist = "store exists - "
	StrStoreNotExist = "store not exist - "
	StrStoreTypeEmpty = "st empty - [set, oset]"
	StrStoreTypeNotExist = " isn't in [set, oset]"
	StrNameEmpty = "n empty"
	StrTypeEmpty = "t empty - [int, float, str]"
	StrKeyEmpty = "k empty"
	StrDataParseError = "pasring error - "
	StrKeyExist = "key exists - "
	StrKeyNotExist = "key not exist - "
	StrTypeNotExist = " isn't in [int, float, str]"
	
	StrDqStoreTypeNotExist = " isn't in [int, float, string]"
	StrDqStoreGetActionEmpty = "a emtpy - [size]"
	StrDqStoreGetActionNotExist = " isn't in [size]"
	StrDqStorePutActionEmpty = "a empty - [clear]"
	StrDqStorePutActionNotExit = " isn't in [clear]"
	StrDqDataPutActionEmpty = "a empty - [f,b]"
	StrDqDataPutActionNotExist = " isn't in - [f,b]"

	StrKvStoreGetActionEmpty = "a emtpy - [size, keys, values]"
	StrKvStoreGetActionNotExist = " isn't in [size, keys, values]"
	StrKvStorePutActionEmpty = "a empty - [clear]"
	StrKvStorePutActionNotExit = " isn't in [clear]"
)

const (
	QueryStoreType = "st"
	QueryStoreTypeSet	= "set"
	QueryStoreTypeOSet	= "oset"
	QueryType = "t"
	QueryTypeInt		= "int"
	QueryTypeStr		= "str"
	QueryTypeFloat	= "float"
	QueryName = "n"
	QueryKey = "k"
	QueryValue = "v"
	QueryAction = "a"
	QueryActionQueFront = "f"
	QueryActionQueBack = "b"
	QueryActionSize = "size"
	QueryActionClear = "clear"
)

type Type uint8
const (
	TypeInt64 = iota
	TypeFloat64
	TypeStr
)

func MapQueryType(s string) (t Type, err error) {
	switch s {
	case QueryTypeInt:
		t = TypeInt64
	case QueryTypeFloat:
		t = TypeFloat64
	case QueryTypeStr:
		t = TypeStr
	default:
		err = errors.New(s + StrTypeNotExist)		
	}
	return
}


func ParseInput(s string, t Type) (r interface{}, ok bool) {
	switch t {
	case TypeInt64:
		r, ok = ParseIntBytes(*(*[]byte)(unsafe.Pointer(&s)))
	case TypeFloat64:
		if v, err := strconv.ParseFloat(s, 64); err==nil  {
			r, ok = v, true
		} else {
			r, ok = nil, false
		}
	case TypeStr:
		r, ok = s, true
	default:
		r, ok = nil, false
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

func ParseFloat(s string) (float64, bool) {
	if v, err := strconv.ParseFloat(s, 64); err==nil  {
		return v, true
	}
	return 0, false
}

func ParseInt(s string) (int64, bool) {	
	return ParseIntBytes(*(*[]byte)(unsafe.Pointer(&s)))
}
// About 3x faster then strconv.ParseInt because does not check for range error and support only base 10, which is enough for JSON
// credit goes to github.com/buger/jsonparser
func ParseIntBytes(bytes []byte) (v int64, ok bool) {
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
	Create(key, value string) bool
	Read(key string) (value interface{}, found bool)
	Update(key, value string) bool
	Delete(Key string) bool
	Clear()
	Size() int
	Values() []interface{}
	Keys() []interface{} 
	RangeGet(startkey, endkey string, limit int, ascending bool) (values []interface{}, count int)
	// Todo: add block CRUD operations: 
	// BlockCreate, BLockRead, BlockUpdate, BlockDelete???
	// BlockCreateUpdate ???
}

// IDataStore defines the interface that all the data store implementation shall comply
// For the data store that can not support perticular interface method, create a dummy
// Supported Store: Deque
type IDequeStore interface {
	PushFront(string) bool
	PushBack(string) bool
	PopFront() (interface{}, bool)
	PopBack() (interface{}, bool) 
	Size() int
	Clear()
}
