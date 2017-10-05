// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kv_test

import (
	"testing"
	"net/http"
	"github.com/dlmc/ids/global"
	"github.com/dlmc/golight/ghttp"
	"strings"
)


func tOsetTest(t *testing.T, rootUrl, storetype string) {

	t.Run("StorePost", func(t *testing.T) {
		tests := []struct {
			url string
			want string
			code int
		}{
			{"/kvstore?storename=mystore&keytype=string&storetype="+storetype+"1", 
			 `{"code":400,"message":"storetype [`+storetype+"1"+`] not in [oset, set]","data":"/kvstore?storename=mystore\u0026keytype=string\u0026storetype=`+storetype+"1"+`"}`+"\n",
			http.StatusBadRequest},
			{"/kvstore?storename=mystore&keytype=string&store1type="+storetype, 
			`{"code":400,"message":"storetype [] not in [oset, set]","data":"/kvstore?storename=mystore\u0026keytype=string\u0026store1type=`+storetype+`"}`+"\n",
			http.StatusBadRequest},
			{"/kvstore?storename=mystore&keytype=string1&storetype="+storetype, 
			`{"code":400,"message":"keytype [string1] not in [string, integer, float]","data":"/kvstore?storename=mystore\u0026keytype=string1\u0026storetype=`+storetype+`"}`+"\n",
			http.StatusBadRequest},
			{"/kvstore?storename=mystore&key1type=string&storetype="+storetype, 
			`{"code":400,"message":"keytype [] not in [string, integer, float]","data":"/kvstore?storename=mystore\u0026key1type=string\u0026storetype=`+storetype+`"}`+"\n",
			http.StatusBadRequest},
			{"/kvstore?store1name=mystore&keytype=string&storetype="+storetype, 
			`{"code":400,"message":"storename empty","data":"/kvstore?store1name=mystore\u0026keytype=string\u0026storetype=`+storetype+`"}`+"\n",
			http.StatusBadRequest},
			{"/kvstore?storename=&keytype=string&storetype="+storetype, 
			`{"code":400,"message":"storename empty","data":"/kvstore?storename=\u0026keytype=string\u0026storetype=`+storetype+`"}`+"\n",
			http.StatusBadRequest},
		}
		for _, test := range tests {
			res,err := http.Post(rootUrl+test.url, "text/html; charset=utf-8", nil)
			tResult(t, res, err, test.want, test.code)
		}		
	})
	t.Run("StoreDelete", func(t *testing.T) {
		tests := []struct {
			method string
			url string
			want string
			code int
		}{
			{"POST", "/kvstore?storename=mystore&keytype=string&storetype="+storetype, 
			 `{"code":201,"message":"success"}`+"\n",
			http.StatusCreated},
			{"DELETE", "/kvstore?store1name=mystore", 
			 `{"code":400,"message":"storename empty","data":"/kvstore?store1name=mystore"}`+"\n",
			http.StatusBadRequest},
			{"DELETE", "/kvstore?storename=my1store", 
			 `{"code":400,"message":"store does not exists - my1store","data":"/kvstore?storename=my1store"}`+"\n",
			http.StatusBadRequest},
			{"DELETE", "/kvstore?storename=mystore", 
			 `{"code":200,"message":"success"}`+"\n",
			http.StatusOK},
		}
		for _, test := range tests {
			res,err := ghttp.ClientHttp(test.method,  rootUrl+test.url, "text/html; charset=utf-8", nil)
			tResult(t, res, err, test.want, test.code)
		}
		
	})
	t.Run("StoreClear", func(t *testing.T) {
		tests := []struct {
			method string
			url string
			want string
			code int
		}{
			{"POST", "/kvstore?storename=mystore&keytype=string&storetype="+storetype, 
			 `{"code":201,"message":"success"}`+"\n",
			http.StatusCreated},
			{"PUT", "/kvstore?storename=mystore&action=cl1ear", 
			 `{"code":400,"message":"cl1ear is not in [clear]","data":"/kvstore?storename=mystore\u0026action=cl1ear"}`+"\n",
			http.StatusBadRequest},
			{"PUT", "/kvstore?storename=mystore&acti1on=clear", 
			 `{"code":400,"message":"action empty - [clear]","data":"/kvstore?storename=mystore\u0026acti1on=clear"}`+"\n",
			http.StatusBadRequest},
			{"PUT", "/kvstore?storename=mystore1&action=clear", 
			 `{"code":400,"message":"store does not exists - mystore1","data":"/kvstore?storename=mystore1\u0026action=clear"}`+"\n",
			http.StatusBadRequest},
			{"PUT", "/kvstore?store1name=mystore&action=clear", 
			 `{"code":400,"message":"storename empty","data":"/kvstore?store1name=mystore\u0026action=clear"}`+"\n",
			http.StatusBadRequest},
			{"PUT", "/kvstore?storename=mystore&action=clear", 
			 `{"code":200,"message":"success"}`+"\n",
			http.StatusOK},
			{"DELETE", "/kvstore?storename=mystore", 
			 `{"code":200,"message":"success"}`+"\n",
			http.StatusOK},
		}
		for _, test := range tests {
			res,err := ghttp.ClientHttp(test.method,  rootUrl+test.url, "text/html; charset=utf-8", nil)
			tResult(t, res, err, test.want, test.code)
		}
		
	})
	t.Run("Data", func(t *testing.T) {
		url := rootUrl+"/kvstore?storename=mystore&keytype=string&storetype="+storetype
		res,err := http.Post(url, "text/html; charset=utf-8", nil)
		want := `{"code":201,"message":"success"}`+"\n"
		tResult(t, res, err, want, http.StatusCreated)
		t.Run("DataPost", func(t *testing.T) {
			smsg := `{"code":201,"message":"success"}`+"\n"
			tests := []struct {
				k, v string
				code int
				want string
			}{
				{"1k", "a", http.StatusCreated, smsg},
				{"2k", "b", http.StatusCreated, smsg},
				{"3k", "c", http.StatusCreated, smsg},
				{"4k", "d", http.StatusCreated, smsg},
				{"4k", "g", http.StatusBadRequest, `{"code":400,"message":"key exists - 4k","data":"/kvdata?storename=mystore\u0026k=4k\u0026v=g"}`+"\n"},
			}
			for _, test := range tests {
				url := rootUrl+"/kvdata?storename=mystore&k="+test.k+"&v="+test.v
				res,err := http.Post(url, "text/html; charset=utf-8", nil)
				tResult(t, res, err, test.want, test.code)
			}
		})
		t.Run("DataPut", func(t *testing.T) {
			smsg := `{"code":200,"message":"success"}`+"\n"
			tests := []struct {
				k, v string
				code int
				want string
			}{
				{"1k", "g", http.StatusOK, smsg},
				{"2k", "f", http.StatusOK, smsg},
				{"3k", "e", http.StatusOK, smsg},
				{"4k", "d", http.StatusOK, smsg},
				{"8k", "g", http.StatusBadRequest, `{"code":400,"message":"key does not exist - 8k","data":"/kvdata?storename=mystore\u0026k=8k\u0026v=g"}`+"\n"},
			}
			for _, test := range tests {
				url := rootUrl+"/kvdata?storename=mystore&k="+test.k+"&v="+test.v
				res,err := ghttp.Put(url, "text/html; charset=utf-8", nil)
				tResult(t, res, err, test.want, test.code)
			}
		})
		t.Run("DataDelete", func(t *testing.T) {
			smsg := `{"code":200,"message":"success"}`+"\n"
			tests := []struct {
				k string
				code int
				want string
			}{
				{"1k", http.StatusOK, smsg},
				{"2k", http.StatusOK, smsg},
				{"8k", http.StatusBadRequest, `{"code":400,"message":"key does not exist - 8k","data":"/kvdata?storename=mystore\u0026k=8k"}`+"\n"},
			}
			for _, test := range tests {
				url := rootUrl+"/kvdata?storename=mystore&k="+test.k
				res,err := ghttp.Del(url, "text/html; charset=utf-8", nil)
				tResult(t, res, err, test.want, test.code)
			}
		})
		t.Run("DataGet", func(t *testing.T) {
			total := 100
			want = `{"code":201,"message":"success"}`+"\n"
			sCode := http.StatusCreated
			t.Run("PostDataAsQueryParameter", func(t *testing.T) {
				for i:=0; i<total/2; i++ {
					k := global.Int2StrPadZero(i, 10)
					url := rootUrl+"/kvdata?storename=mystore&k="+ k +"&v="+ k
					res,err := http.Post(url, "text/html; charset=utf-8", nil)
					tResult(t, res, err, want, sCode)
				}
			})
			t.Run("PostDataInBody", func(t *testing.T) {
				for i:=total/2; i<total; i++ {
					k := global.Int2StrPadZero(i, 10)
					url := rootUrl+"/kvdata?storename=mystore&k="+ k
					res,err := http.Post(url, "text/html; charset=utf-8", strings.NewReader(k))
					tResult(t, res, err, want, sCode)
				}
			})
			sCode = http.StatusOK
			t.Run("GetDataBack", func(t *testing.T) {
				for i:=0; i<total; i++ {
					k := global.Int2StrPadZero(i, 10)
					url := rootUrl+"/kvdata?storename=mystore&k="+ k
					res,err := http.Get(url)
					want = `{"code":200,"message":"success","data":"`+k+`"}`+"\n"
					tResult(t, res, err, want, sCode)
				}
			})
		})
		t.Run("DeleteStore", func(t *testing.T) {
			want := `{"code":200,"message":"success"}`+"\n"
			res,err := ghttp.Del(rootUrl+"/kvstore?storename=mystore", "text/html; charset=utf-8", nil)
			tResult(t, res, err, want, http.StatusOK)
		})
		t.Run("Scale", func(t *testing.T) {
			tCreatePerformance(t, rootUrl, storetype)
		})
	})
}
