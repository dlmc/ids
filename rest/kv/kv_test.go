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


func tOsetTest(t *testing.T, rootUrl, st string) {

	t.Run("StorePost", func(t *testing.T) {
		tests := []struct {
			url string
			want string
			code int
		}{
			{"/kvstore?n=mystore&t=str&st="+st+"1", 
			 `{"code":400,"message":"`+st+"1"+` isn't in [set, oset]","data":"/kvstore?n=mystore\u0026t=str\u0026st=`+st+"1"+`"}`+"\n",
			http.StatusBadRequest},
			{"/kvstore?n=mystore&t=str&store1type="+st, 
			`{"code":400,"message":"st empty - [set, oset]","data":"/kvstore?n=mystore\u0026t=str\u0026store1type=`+st+`"}`+"\n",
			http.StatusBadRequest},
			{"/kvstore?n=mystore&t=str1&st="+st, 
			`{"code":400,"message":"str1 isn't in [int, float, str]","data":"/kvstore?n=mystore\u0026t=str1\u0026st=`+st+`"}`+"\n",
			http.StatusBadRequest},
			{"/kvstore?n=mystore&key1type=str&st="+st, 
			`{"code":400,"message":"t empty - [int, float, str]","data":"/kvstore?n=mystore\u0026key1type=str\u0026st=`+st+`"}`+"\n",
			http.StatusBadRequest},
			{"/kvstore?store1name=mystore&t=str&st="+st, 
			`{"code":400,"message":"n empty","data":"/kvstore?store1name=mystore\u0026t=str\u0026st=`+st+`"}`+"\n",
			http.StatusBadRequest},
			{"/kvstore?n=&t=str&st="+st, 
			`{"code":400,"message":"n empty","data":"/kvstore?n=\u0026t=str\u0026st=`+st+`"}`+"\n",
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
			{"POST", "/kvstore?n=mystore&t=str&st="+st, 
			 `{"code":201,"message":"success"}`+"\n",
			http.StatusCreated},
			{"DELETE", "/kvstore?store1name=mystore", 
			 `{"code":400,"message":"n empty","data":"/kvstore?store1name=mystore"}`+"\n",
			http.StatusBadRequest},
			{"DELETE", "/kvstore?n=my1store", 
			 `{"code":400,"message":"store not exist - my1store","data":"/kvstore?n=my1store"}`+"\n",
			http.StatusBadRequest},
			{"DELETE", "/kvstore?n=mystore", 
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
			{"POST", "/kvstore?n=mystore&t=str&st="+st, 
			 `{"code":201,"message":"success"}`+"\n",
			http.StatusCreated},
			{"PUT", "/kvstore?n=mystore&a=cl1ear", 
			 `{"code":400,"message":"cl1ear isn't in [clear]","data":"/kvstore?n=mystore\u0026a=cl1ear"}`+"\n",
			http.StatusBadRequest},
			{"PUT", "/kvstore?n=mystore&acti1on=clear", 
			 `{"code":400,"message":"a empty - [clear]","data":"/kvstore?n=mystore\u0026acti1on=clear"}`+"\n",
			http.StatusBadRequest},
			{"PUT", "/kvstore?n=mystore1&a=clear", 
			 `{"code":400,"message":"store not exist - mystore1","data":"/kvstore?n=mystore1\u0026a=clear"}`+"\n",
			http.StatusBadRequest},
			{"PUT", "/kvstore?store1name=mystore&a=clear", 
			 `{"code":400,"message":"n empty","data":"/kvstore?store1name=mystore\u0026a=clear"}`+"\n",
			http.StatusBadRequest},
			{"PUT", "/kvstore?n=mystore&a=clear", 
			 `{"code":200,"message":"success"}`+"\n",
			http.StatusOK},
			{"DELETE", "/kvstore?n=mystore", 
			 `{"code":200,"message":"success"}`+"\n",
			http.StatusOK},
		}
		for _, test := range tests {
			res,err := ghttp.ClientHttp(test.method,  rootUrl+test.url, "text/html; charset=utf-8", nil)
			tResult(t, res, err, test.want, test.code)
		}
		
	})
	t.Run("Data", func(t *testing.T) {
		url := rootUrl+"/kvstore?n=mystore&t=str&st="+st
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
				{"4k", "g", http.StatusBadRequest, `{"code":400,"message":"key exists - 4k","data":"/kvdata?n=mystore\u0026k=4k\u0026v=g"}`+"\n"},
			}
			for _, test := range tests {
				url := rootUrl+"/kvdata?n=mystore&k="+test.k+"&v="+test.v
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
				{"8k", "g", http.StatusBadRequest, `{"code":400,"message":"key not exist - 8k","data":"/kvdata?n=mystore\u0026k=8k\u0026v=g"}`+"\n"},
			}
			for _, test := range tests {
				url := rootUrl+"/kvdata?n=mystore&k="+test.k+"&v="+test.v
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
				{"8k", http.StatusBadRequest, `{"code":400,"message":"key not exist - 8k","data":"/kvdata?n=mystore\u0026k=8k"}`+"\n"},
			}
			for _, test := range tests {
				url := rootUrl+"/kvdata?n=mystore&k="+test.k
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
					url := rootUrl+"/kvdata?n=mystore&k="+ k +"&v="+ k
					res,err := http.Post(url, "text/html; charset=utf-8", nil)
					tResult(t, res, err, want, sCode)
				}
			})
			t.Run("PostDataInBody", func(t *testing.T) {
				for i:=total/2; i<total; i++ {
					k := global.Int2StrPadZero(i, 10)
					url := rootUrl+"/kvdata?n=mystore&k="+ k
					res,err := http.Post(url, "text/html; charset=utf-8", strings.NewReader(k))
					tResult(t, res, err, want, sCode)
				}
			})
			sCode = http.StatusOK
			t.Run("GetDataBack", func(t *testing.T) {
				for i:=0; i<total; i++ {
					k := global.Int2StrPadZero(i, 10)
					url := rootUrl+"/kvdata?n=mystore&k="+ k
					res,err := http.Get(url)
					want = `{"code":200,"message":"success","data":"`+k+`"}`+"\n"
					tResult(t, res, err, want, sCode)
				}
			})
		})
		t.Run("DeleteStore", func(t *testing.T) {
			want := `{"code":200,"message":"success"}`+"\n"
			res,err := ghttp.Del(rootUrl+"/kvstore?n=mystore", "text/html; charset=utf-8", nil)
			tResult(t, res, err, want, http.StatusOK)
		})
		t.Run("Scale", func(t *testing.T) {
			//tCreatePerformance(t, rootUrl, st)
		})
	})
}
