// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kv_test

import (
	"testing"
	"time"
	"math/rand"
	"strings"
	"net/http"
	"github.com/dlmc/ids/kvstore"
	"github.com/dlmc/ids/global"
	"github.com/dlmc/golight/ghttp"
)

func tReadRequest(t *testing.T, ost global.IKvStore, nRequest, nTotalKeys int, key, str string) {
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
func tReadRequestHttp(t *testing.T, nRequest, nTotalKeys int, key, str, rootUrl string) {
	code := http.StatusOK
	for i:=0; i < nRequest; i++ {
		k := rand.Intn(nTotalKeys)
		kpad := global.Int2StrPadZero(k,10)
		
		u := rootUrl+"/kvdata?n=mystore&k="+ kpad+key
		res,err := http.Get(u)
		want := `{"code":200,"message":"success","data":"`+kpad+str+`"}`+"\n"
		tResult(t, res, err, want, code)

	}
}

func tRunTest(t *testing.T, nTotalKeys int, a []int, key, str, storetype string) {
	store := kvstore.New()
	if err:= store.CreateStore("mystore", storetype, "str"); err != nil {
		t.Errorf("creating storetype: %s failed", storetype)			
	} else {
		if kvs, found := store.GetStore("mystore"); !found {
			t.Errorf("can't find store %s", "mystore")			
		} else {
			t.Run("Create elements", func(t *testing.T) {
				a = a[:nTotalKeys]
				for _, v := range a {
					vpad := global.Int2StrPadZero(v,10)
					k := vpad+key
					value := vpad+str
					kvs.Create(k, value)
				}
			})
			t.Run("1000 random keys request", func(t *testing.T) {
				tReadRequest(t, kvs, 1000, nTotalKeys, key, str)
			})
		}
	}
}

func tRunTestHttp(t *testing.T, nTotalKeys int, a []int, key, str, rootUrl, storetype string, inQuery bool) {
	want := `{"code":201,"message":"success"}`+"\n"
	code := http.StatusCreated

	t.Run("Create Store", func(t *testing.T) {
		res,err := http.Post(rootUrl +"/kvstore?n=mystore&t=string&st="+storetype, "text/html; charset=utf-8", nil)
		tResult(t, res, err, want, code)
	})
	
	t.Run("Create elements", func(t *testing.T) {
		a = a[:nTotalKeys]
		for _, v := range a {
			kpad := global.Int2StrPadZero(v,10)
			k := kpad+key
			value := kpad+str
			if inQuery {
				u := rootUrl +"/kvdata?n=mystore&k="+k+"&v="+value
				res,err := http.Post(u, "text/html; charset=utf-8", nil)
				tResult(t, res, err, want, code)
			} else {
				u := rootUrl +"/kvdata?n=mystore&k="+k
				res,err := http.Post(u, "text/html; charset=utf-8", strings.NewReader(value))
				tResult(t, res, err, want, code)				
			}
		}
	})
	t.Run("1000 random keys request", func(t *testing.T) {
		tReadRequestHttp(t, 1000, nTotalKeys, key, str, rootUrl)
	})
	t.Run("DeleteStore", func(t *testing.T) {
		want := `{"code":200,"message":"success"}`+"\n"
		res,err := ghttp.Del(rootUrl+"/kvstore?n=mystore", "text/html; charset=utf-8", nil)
		tResult(t, res, err, want, http.StatusOK)
	})

}



func tCreatePerformance20000Elements(t *testing.T, totalKeys, repeats int, rootUrl, storetype string) {
	key := "abcdefghijklmnopqrstuvwxyz"
	str := strings.Repeat(key,repeats)
	rand.Seed(time.Now().UTC().UnixNano())
	a := rand.Perm(totalKeys)
	t.Run("DataGetFunctionCall", func(t *testing.T) {
		tRunTest(t, totalKeys, a, key, str, storetype)
	})
	t.Run("DataInHttpQuery", func(t *testing.T) {
		tRunTestHttp(t, totalKeys, a, key, str, rootUrl, storetype, true)
	})
	t.Run("DataInHttpBody", func(t *testing.T) {
		tRunTestHttp(t, totalKeys, a, key, str, rootUrl, storetype, false)
	})
}

func tCreatePerformance(t *testing.T, rootUrl, storetype string) {	
	t.Run("20000x26 Elements", func(t *testing.T) {
		tCreatePerformance20000Elements(t, 20000, 1, rootUrl, storetype)
	})
	t.Run("20000x130 Elements", func(t *testing.T) {
		tCreatePerformance20000Elements(t, 20000, 5, rootUrl, storetype)
	})
	t.Run("20000x1.3k Elements", func(t *testing.T) {
		tCreatePerformance20000Elements(t, 20000, 50, rootUrl, storetype)
	})
	t.Run("20000x2.6k Elements", func(t *testing.T) {
		tCreatePerformance20000Elements(t, 20000, 100, rootUrl, storetype)
	})
}

