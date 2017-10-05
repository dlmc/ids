// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package set_test

import (
//	"testing"
	"fmt"
	"github.com/dlmc/ids/global"
	"github.com/dlmc/ids/set"
)


func ExampleOrderedSet() {
		st := set.New("first", global.KTStr)
		//ds.Clear() // empty
		fmt.Println(st.Size())
		// Output: 
		// 0
}
