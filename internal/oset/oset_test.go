

// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package oset_test

import (
	"testing"
	"fmt"
	"reflect"
	"time"
	"math/rand"
	"strings"
	"github.com/dlmc/ids/oset"
	"github.com/dlmc/ids/global"
)

func ExampleInt2StrPadZero() {
	fmt.Println(global.Int2StrPadZero(123, 5))
	fmt.Println(global.Int2StrPadZero(1234, 3))
	//Output: 
	//00123
	//1234
}


func ExampleOrderedSet() {
		ost := oset.New("first", global.KTStr)
		//ds.Clear() // empty
		fmt.Println(ost.Size())
		// Output: 
		// 0
}

func tRandomRangeRequest(t *testing.T, ost *oset.OSet, nRequest, nReqRange, nTotalKeys int, key, str string) {
		for i:=0; i < nRequest; i++ {
			k := rand.Intn(nTotalKeys-nReqRange+1)
			skey := global.Int2StrPadZero(k,10)+key
			ekey := global.Int2StrPadZero(k+nReqRange,10)+key
			
			if v, cnt := ost.RangeGet(skey, ekey, nReqRange, true); cnt==nReqRange {
				for idx, d := range v {
					if exp := global.Int2StrPadZero(k+idx,10)+str; d != exp {
						t.Errorf("Got %v expected %v", d, exp)
					}
				}
			} else {
				t.Errorf("Got %v expect %v\n", cnt, nReqRange)			
			}
		}	
}

func tReadRequest(t *testing.T, ost *oset.OSet, nRequest, nTotalKeys int, key, str string) {
	for i:=0; i < nRequest; i++ {
		k := rand.Intn(nTotalKeys)
		kpad := global.Int2StrPadZero(k,10)
		if v, found := ost.Read(kpad+key); found==true {
			if v != kpad+str {
				t.Errorf("Got %v expected %v\n", v, kpad+str)
			}
		} else {
			t.Errorf("Got not found for key %v\n", v, kpad+key)			
		}
	}
}

func tRunTest(t *testing.T, nTotalKeys int, a []int, key, str string) {
	ost := oset.New("oset", global.KTStr)
	t.Run("Create elements", func(t *testing.T) {
		a = a[:nTotalKeys]
		for _, v := range a {
			kpad := global.Int2StrPadZero(v,10)
			ost.Create(kpad+key, kpad+str)
		}
	})
	t.Run("10000 random keys request", func(t *testing.T) {
		tReadRequest(t, ost, 10000, nTotalKeys, key, str)
	})
	t.Run("10000 random range request with 10 elements", func(t *testing.T) {
		tRandomRangeRequest(t, ost, 10000, 10, nTotalKeys, key, str)
	})
}


func TestCreate(t* testing.T) {
	t.Run("CreatePositiveCase", func(t *testing.T) {
		ost := oset.New("oset", global.KTStr)
		ost.Create("5", "e")
		ost.Create("6", "f")
		ost.Create("7", "g")
		ost.Create("3", "c")
		ost.Create("4", "d")
		ost.Create("1", "a")
		ost.Create("2", "b")
		
		if actualValue := ost.Size(); actualValue != 7 {
			t.Errorf("Got %v expected %v", actualValue, 7)
		}
		if actualValue, expectedValue := fmt.Sprintf("%s%s%s%s%s%s%s", ost.Keys()...), "1234567"; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
		if actualValue, expectedValue := fmt.Sprintf("%s%s%s%s%s%s%s", ost.Values()...), "abcdefg"; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	
		tests1 := [][]interface{}{
			{"1", "a", true},
			{"2", "b", true},
			{"3", "c", true},
			{"4", "d", true},
			{"5", "e", true},
			{"6", "f", true},
			{"7", "g", true},
			{"8", nil, false},
		}
	
		for _, test := range tests1 {
			actualValue, actualFound := ost.Read(test[0])
			if actualValue != test[1] || actualFound != test[2] {
				t.Errorf("Got %v expected %v", actualValue, test[1])
			}
		}
	})
	t.Run("CreateNegativeCase", func(t *testing.T) {
		ost := oset.New("oset", global.KTStr)
		ost.Create("2", "b")
		r1 := ost.Create("1", "x")
		r2 := ost.Create("1", "a") 
		
		if r1 != true || r2 != false {
			t.Errorf("Got r1 %v r2 %v", r1, r2)
		}
		if actualValue, expectedValue := fmt.Sprintf("%s%s", ost.Values()...), "xb"; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	})
}


func TestUpdate(t* testing.T) {
	t.Run("UpdatePositiveCase", func(t *testing.T) {
		ost := oset.New("oset", global.KTStr)
		ost.Create("1", "a")
		ost.Create("2", "b")
		ost.Update("1", "x")
		
		if actualValue := ost.Size(); actualValue != 2 {
			t.Errorf("Got %v expected %v", actualValue, 2)
		}
		if actualValue, expectedValue := fmt.Sprintf("%s%s", ost.Keys()...), "12"; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
		if actualValue, expectedValue := fmt.Sprintf("%s%s", ost.Values()...), "xb"; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	
		tests1 := [][]interface{}{
			{"1", "x", true},
			{"2", "b", true},
			{"8", nil, false},
		}
	
		for _, test := range tests1 {
			actualValue, actualFound := ost.Read(test[0])
			if actualValue != test[1] || actualFound != test[2] {
				t.Errorf("Got %v expected %v", actualValue, test[1])
			}
		}
	})
	t.Run("UpdateNegativeCase", func(t *testing.T) {
		ost := oset.New("oset", global.KTStr)
		ost.Create("2", "b")
		ost.Create("1", "a")
		r1 := ost.Update("1", "x") 
		r2 := ost.Update("3", "c") 
		
		if r1 != true || r2 != false {
			t.Errorf("Got r1 %v r2 %v", r1, r2)
		}
		if actualValue, expectedValue := fmt.Sprintf("%s%s", ost.Values()...), "xb"; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	})
}


func TestDelete(t* testing.T) {
	t.Run("DeleteCase", func(t *testing.T) {
		ost := oset.New("oset", global.KTStr)
		ost.Create("1", "a")
		ost.Create("2", "b")
		ost.Create("3", "c")
	
		tests1 := [][]interface{}{
			{"1", true},
			{"4", false},
			{"2", true},
			{"3", true},
			{"8", false},
		}
	
		for _, test := range tests1 {
			deleted := ost.Delete(test[0])
			if deleted != test[1] {
				t.Errorf("Got %v expected %v", deleted, test[1])
			}
		}
	})	
}

func TestClear(t* testing.T) {
	t.Run("ClearCase", func(t *testing.T) {
		ost := oset.New("oset", global.KTStr)
		ost.Create("1", "a")
		ost.Create("2", "b")
		ost.Create("3", "c")
		if actualValue := ost.Size(); actualValue != 3 {
			t.Errorf("Got %v expected %v", actualValue, 3)
		}		
		ost.Clear()
		if actualValue := ost.Size(); actualValue != 0 {
			t.Errorf("Got %v expected %v", actualValue, 0)
		}
	})	
}


func TestRangeGet(t* testing.T) {
	t.Run("RangeGetPositiveCase", func(t *testing.T) {
		ost := oset.New("oset", global.KTStr)
		ost.Create("5", "e")
		ost.Create("6", "f")
		ost.Create("7", "g")
		ost.Create("3", "c")
		ost.Create("4", "d")
		ost.Create("1", "a")
		ost.Create("2", "b")
			
		tests1 := [][]interface{}{
			{"2", "5", 10, true,  4, "b","c","d","e"},
			{"2", "5", 10, false, 4, "e","d","c","b"},
			{"5", "a", 10, true,  3, "e","f","g"},
			{"5", "a", 10, false, 3, "g","f","e"},
			{"0", "4", 3,  true,  3, "a","b","c"},
			{"0", "4", 3,  false, 3, "d","c","b"},			
			{"0", "4", 0,  false, 0},			
		}
	
		for _, test := range tests1 {
			aValues, count := ost.RangeGet(test[0], test[1], test[2].(int), test[3].(bool))
			if count != test[4] {
				t.Errorf("Got %v expected %v", count, test[4])
			}
			if !reflect.DeepEqual(aValues,test[5:]) {
				t.Errorf("Got %v expected %v", aValues, test[5:])
			}
		}
	})
}

func TestRangeGet1000000ElementsSmallValue(t *testing.T) {
	//key := "abcdefghijklmnopqrstuvwxyz"
	//str := strings.Repeat(key,100)
	rand.Seed(time.Now().UTC().UnixNano())
	totalKeys := 1000000
	a := rand.Perm(totalKeys)
	tRunTest(t, totalKeys, a, "key", "value")
}

func TestCreatePerformance1000Elements(t *testing.T) {
	key := "abcdefghijklmnopqrstuvwxyz"
	str := strings.Repeat(key,100)
	rand.Seed(time.Now().UTC().UnixNano())
	totalKeys := 1000
	a := rand.Perm(totalKeys)
	tRunTest(t, totalKeys, a, key, str)
}

func TestCreatePerformance10000Elements(t *testing.T) {
	key := "abcdefghijklmnopqrstuvwxyz"
	str := strings.Repeat(key,100)
	rand.Seed(time.Now().UTC().UnixNano())
	totalKeys := 10000
	a := rand.Perm(totalKeys)
	tRunTest(t, totalKeys, a, key, str)
}

func TestCreatePerformance100000(t *testing.T) {
	key := "abcdefghijklmnopqrstuvwxyz"
	str := strings.Repeat(key,100)
	rand.Seed(time.Now().UTC().UnixNano())
	totalKeys := 100000
	a := rand.Perm(totalKeys)
	tRunTest(t, totalKeys, a, key, str)
}

func TestCreatePerformance1000000Elements(t *testing.T) {
	key := "abcdefghijklmnopqrstuvwxyz"
	str := strings.Repeat(key,100)
	rand.Seed(time.Now().UTC().UnixNano())
	totalKeys := 1000000
	a := rand.Perm(totalKeys)
	tRunTest(t, totalKeys, a, key, str)
}


