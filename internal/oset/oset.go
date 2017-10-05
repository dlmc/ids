// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oset

import (
	"github.com/dlmc/ids/global"
	avl "github.com/dlmc/ids/internal/avl"
	"sync"
//	"fmt"
)

// OSet implements the IDataStore interface
// OSet is a set with elements sorted by its keys
// Access to the set is protected for concurrent access
type OSet struct {
	*avl.Tree
	sync.RWMutex
	Name string
}

// New creates a OSet instance with the specified KeyType
func New(name string, kt global.KeyType) *OSet {
	//fmt.Println("New OSet Store")
	var tree *avl.Tree
	switch kt {
	case global.KTInt:
		tree = avl.NewWithIntComparator()
	case global.KTFloat64:
		tree = avl.NewWithFloat64Comparator()
	case global.KTStr:
		fallthrough 
	default:
		tree = avl.NewWithStringComparator()
	}
	return &OSet{Name:name, Tree:tree}
}

// Create creates the key/value pair in the set
// Return true if key does not exist, and key/value pair created
// Return false if key already exists, and no change made to the set
func (d *OSet) Create(key interface{}, value interface{}) bool{
	d.Lock()
	defer d.Unlock()
	return d.Tree.Put(key, value)
}

// Read gets the value of the specified key
// Return value if found is true or nil if found false
func (d *OSet) Read(key interface{}) (value interface{}, found bool){
	d.RLock()
	defer d.RUnlock()
	return d.Tree.Get(key)
}

// Update updates the value of an existing key in the set
// Return true if key exists, and value updates sucessfully
// Return false if key does not exist, and no change made to the set
func (d *OSet) Update(key interface{}, value interface{}) bool {
	d.Lock()
	defer d.Unlock()
	if d.Tree.Remove(key) {
		return d.Tree.Put(key, value)
	}
	return false
}

// Delete deletes an existing key in the set
// Return true if key exists, and deletes sucessfully
// Return false if key does not exist, and no change made to the set
func (d *OSet) Delete(key interface{}) bool {
	d.Lock()
	defer d.Unlock()
	return d.Tree.Remove(key)
}

// Clear wapes out the whole set and make it empty
func (d *OSet) Clear() {
	d.Lock()
	defer d.Unlock()
	d.Tree.Clear()
}

// Size returns the number of elements in the set
func (d *OSet) Size() int {
	d.RLock()
	defer d.RUnlock()
	return d.Tree.Size()
}

// Values returns the list of all the values in the set in ascending sorted order
func (d *OSet) Values() []interface{} {
	d.RLock()
	defer d.RUnlock()
	return d.Tree.Values()
}

// Keys returns the list of all the keys in the set in ascending sorted order
func (d *OSet) Keys() []interface{} {
	d.RLock()
	defer d.RUnlock()
	return d.Tree.Keys()
}

// Ceiling node is defined as the smallest node that is larger than or equal to the given node.
// Floor node is defined as the largest node that is smaller than or equal to the given node.
// Floor(key interface{}) (floor *Node, found bool) 
func (d *OSet) RangeGet(startkey, endkey interface{}, limit int, ascending bool) (values []interface{}, count int) {
	d.RLock()
	defer d.RUnlock()
	values = make([]interface{}, limit) //incase limit is 0
	count = 0
	if ascending {
		if n, found := d.Tree.Ceiling(startkey); found {
			for (n != nil) && (count < limit) && (d.Tree.Comparator(n.Key, endkey)<=0) {
				values[count] = n.Value
				count++
				n = n.Next()
			}
			if count != limit {
				values = values[:count]
			}
			//return values, count
			return
		}
	} else {
		if n, found := d.Tree.Floor(endkey); found {
			for (n != nil) && (count < limit) && (d.Tree.Comparator(n.Key, startkey)>=0) {
				values[count] = n.Value
				count++
				n = n.Prev()
			}
			if count != limit {
				values = values[:count]
			}
			return
		}
	}
	return nil, 0
}