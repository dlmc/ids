// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package set

import (
	"github.com/dlmc/ids/global"
	//"sync"
	//"fmt"
)

// Set implements the IDataStore interface
// Set is a set with k/v elements
// Access to the set is protected for concurrent access
type Set struct {
	m map[interface{}]interface{}
	Keytype global.KeyType
	//sync.RWMutex
	Name string
}

// New creates a Set instance with the specified KeyType
func New(name string, kt global.KeyType) *Set {
	//fmt.Println("New Set Store")
	return &Set{Name:name, Keytype:kt, m:make(map[interface{}]interface{})}
}

// Create creates the key/value pair in the set
// Return true if key does not exist, and key/value pair created
// Return false if key already exists, and no change made to the set
func (d *Set) Create(key, value interface{}) bool{
	//may get ride of the lock by
	m := d.m
	//d.Lock()
	//defer d.Unlock()
	var found bool
	if _, found = m[key]; !found {
		m[key]=value
	}
	return !found
}

// Read gets the value of the specified key
// Return value if found is true or nil if found false
func (d *Set) Read(key interface{}) (value interface{}, found bool){
	//d.RLock()
	//defer d.RUnlock()
	m := d.m
	value, found = m[key]
	return
}

// Update updates the value of an existing key in the set
// Return true if key exists, and value updates sucessfully
// Return false if key does not exist, and no change made to the set
func (d *Set) Update(key, value interface{}) bool {
	//d.Lock()
	//defer d.Unlock()
	m := d.m
	var found bool
	if _, found = m[key]; found {
		m[key]=value
	}
	return found
}

// Delete deletes an existing key in the set
// Return true if key exists, and deletes sucessfully
// Return false if key does not exist, and no change made to the set
func (d *Set) Delete(key interface{}) bool {
	//d.Lock()
	//defer d.Unlock()	
	m := d.m
	var found bool
	if _, found = m[key]; found {
		delete(m, key)
	}
	return found
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
func (d *Set) RangeGet(startkey, endkey interface{}, limit int, ascending bool) (values []interface{}, count int) {
	return nil, 0
}