// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package set

import (
	"github.com/dlmc/ids/global"
)

// Set implements the IDataStore interface
// Set is a set with k/v elements
// Access to the set is protected for concurrent access
type Set struct {
	//sync.RWMutex
	m map[interface{}]interface{}
	KeyType global.Type
	Name string
}

// New creates a Set instance with the specified KeyType
func New(name, ktype string) (*Set, error) {
	t, err := global.MapQueryType(ktype)
	if err == nil {
		return &Set{Name:name, KeyType:t, m:make(map[interface{}]interface{})}, nil
	}
	return nil, err
}

// Create creates the key/value pair in the set
// Return true if key does not exist, and key/value pair created
// Return false if key already exists, and no change made to the set
func (d *Set) Create(key, value string) bool {
	//may get ride of the lock by
	m := d.m
	//d.Lock()
	//defer d.Unlock()
	
	k, ok := global.ParseInput(key, d.KeyType)
	if ok {
		if _, ok = m[k]; !ok {
			m[k]=value
			ok = true
		} else {
			ok = false
		}
	}
	return ok
}

// Read gets the value of the specified key
// Return value if found is true or nil if found false
func (d *Set) Read(key string) (value interface{}, found bool){
	//d.RLock()
	//defer d.RUnlock()
	m := d.m

	k, ok := global.ParseInput(key, d.KeyType)
	if ok {
		value, found = m[k]
	}
	return
}

// Update updates the value of an existing key in the set
// Return true if key exists, and value updates sucessfully
// Return false if key does not exist, and no change made to the set
func (d *Set) Update(key, value string) bool {
	//d.Lock()
	//defer d.Unlock()
	m := d.m

	k, ok := global.ParseInput(key, d.KeyType)
	if ok {
		if _, ok = m[k]; ok {
			m[k]=value
		}
	}
	return ok
}

// Delete deletes an existing key in the set
// Return true if key exists, and deletes sucessfully
// Return false if key does not exist, and no change made to the set
func (d *Set) Delete(key string) bool {
	//d.Lock()
	//defer d.Unlock()	
	m := d.m
	k, ok := global.ParseInput(key, d.KeyType)
	if ok {
		if _, ok = m[k]; ok {
			delete(m, k)
			return ok
		}
	}
	return false
}

// Clear wapes out the whole set and make it empty
func (d *Set) Clear() {
	//d.Lock()
	//defer d.Unlock()
	d.m = make(map[interface{}]interface{})
}

// Size returns the number of elements in the set
func (d *Set) Size() int {
	//d.RLock()
	//defer d.RUnlock()
	m := d.m
	return len(m)
}

// Values returns the list of all the values in the set
func (d *Set) Values() []interface{} {
	//d.RLock()
	//defer d.RUnlock()
	m := d.m
	values := make([]interface{}, len(m))
	i := 0
	for _, value := range m {
		values[i] = value
		i++
	}
	return values
}

// Keys returns the list of all the keys in the set
func (d *Set) Keys() []interface{} {
	//d.RLock()
	//defer d.RUnlock()
	m := d.m
	values := make([]interface{}, len(m))
	i := 0
	for key, _ := range m {
		values[i] = key
		i++
	}
	return values
}

// Not implemented
func (d *Set) RangeGet(startkey, endkey string, limit int, ascending bool) (values []interface{}, count int) {
	return nil, 0
}