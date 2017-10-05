// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kvstore_test

import (	
	"testing"
	"fmt"
/*	"reflect"
	"strings"
	"strconv"
*/	"github.com/dlmc/ids/global"
	dstore "github.com/dlmc/ids/datastore"
)


func ExampleDatastoreOSet() {
		ds := dstore.New()
		if ok := ds.CreateStore("first", global.STOSet, global.KTStr); !ok {
			fmt.Println("CreateStore failed")
		}
		
		oset, ok := ds.GetStore("first")
		//oset.Clear() // empty
		oset.Create("5", "e")
		oset.Create("6", "f")
		oset.Create("7", "g")

		fmt.Println(ok)
		// Output: 
		// true
		fmt.Println(oset.Empty())
		// false		
		fmt.Println(oset.Size())
		// 3
		fmt.Println(oset.Keys())
		// [5 6 7]
		fmt.Println(oset.Values())
		// [e f g]
		fmt.Println(oset.Read("6"))		
		// f true
}
func ExampleDatastoreSet() {
		ds := dstore.New()
		if ok := ds.CreateStore("first", global.STSet, global.KTStr); !ok {
			fmt.Println("CreateStore failed")
		}
		
		oset, ok := ds.GetStore("first")
		//oset.Clear() // empty
		oset.Create("5", "e")
		oset.Create("6", "f")
		oset.Create("7", "g")

		fmt.Println(ok)
		// Output: 
		// true
		fmt.Println(oset.Empty())
		// false		
		fmt.Println(oset.Size())
		// 3
		fmt.Println(oset.Keys())
		// [5 6 7]
		fmt.Println(oset.Read("5"))
		// e true
		fmt.Println(oset.Read("6"))		
		// f true
		fmt.Println(oset.Read("7"))		
		// g true
		fmt.Println(oset.Read("8"))		
		// nil false
}
func TestCreate(t* testing.T) {
}
