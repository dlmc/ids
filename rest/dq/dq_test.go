// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dq_test

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"github.com/dlmc/ids/rest"
	"github.com/dlmc/ids/global"
	"github.com/dlmc/golight/ghttp"
	"github.com/dlmc/golight/decorator/logging"
	"os"
	"io/ioutil"
	"strings"
)


//Process test results
func tResult(t *testing.T, res *http.Response, err error, wantBody string, wantStateCode int) {
	if err != nil {
		t.Errorf("tResult failed", err)
	}
	got, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Errorf("tResult failed", err)
	}
	if wantBody != string(got) {
		t.Errorf("tResult body, got: %s, want: %s", got, wantBody)
	}
	if res.StatusCode != wantStateCode {
		t.Errorf("tResult StatusCode, got: %v, want: %v", res.StatusCode, wantStateCode)
	} 
}

func TestDq(t *testing.T) {
	l := logging.NewLogger(os.Stdout).Level(logging.LogError).With()
	ts := httptest.NewServer(rest.NewServerMux("/", l))
	url := ts.URL
	defer ts.Close()
	t.Run("dq store", func(t *testing.T) {
		tDqStoreTest(t, url)
	})
	t.Run("dq data", func(t *testing.T) {
		tDqDataTest(t, url)
	})
	t.Run("dq scale1000Elementsx2.6k", func(t *testing.T) {
		tDqScaleTest(t, 1000, 100, url)
	})

}
func tDqStoreTest(t *testing.T, rootUrl string) {

	t.Run("StorePost", func(t *testing.T) {
		tests := []struct {
			url string
			want string
			code int
		}{
			{"/dqstore?n=mystore&t=string", 
			 `{"code":201,"message":"success"}`+"\n",
			http.StatusCreated},
			{"/dqstore?n=intstore&t=int", 
			 `{"code":201,"message":"success"}`+"\n",
			http.StatusCreated},
			{"/dqstore?n=floatstore&t=float", 
			 `{"code":201,"message":"success"}`+"\n",
			http.StatusCreated},
			{"/dqstore?store1name=mystore&t=string", 
			`{"code":400,"message":"n empty","data":"/dqstore?store1name=mystore\u0026t=string"}`+"\n",
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
			{"DELETE", "/dqstore?store1name=mystore", 
			 `{"code":400,"message":"n empty","data":"/dqstore?store1name=mystore"}`+"\n",
			http.StatusBadRequest},
			{"DELETE", "/dqstore?n=my1store", 
			 `{"code":400,"message":"store not exist - my1store","data":"/dqstore?n=my1store"}`+"\n",
			http.StatusBadRequest},
			{"DELETE", "/dqstore?n=mystore", 
			 `{"code":200,"message":"success"}`+"\n",
			http.StatusOK},
			{"DELETE", "/dqstore?n=intstore", 
			 `{"code":200,"message":"success"}`+"\n",
			http.StatusOK},
			{"DELETE", "/dqstore?n=floatstore", 
			 `{"code":200,"message":"success"}`+"\n",
			http.StatusOK},
		}
		for _, test := range tests {
			res,err := ghttp.Del(rootUrl+test.url, "text/html; charset=utf-8", nil)
			tResult(t, res, err, test.want, test.code)
		}
		
	})
	t.Run("StorePut", func(t *testing.T) {
		tests := []struct {
			method string
			url string
			want string
			code int
		}{
			{"POST", "/dqstore?n=mystore&t=string", 
			 `{"code":201,"message":"success"}`+"\n",
			http.StatusCreated},
			{"PUT", "/dqstore?n=mystore&a=cl1ear", 
			 `{"code":400,"message":"cl1ear is not in [clear]","data":"/dqstore?n=mystore\u0026a=cl1ear"}`+"\n",
			http.StatusBadRequest},
			{"PUT", "/dqstore?n=mystore&acti1on=clear", 
			 `{"code":400,"message":"a empty - [clear]","data":"/dqstore?n=mystore\u0026acti1on=clear"}`+"\n",
			http.StatusBadRequest},
			{"PUT", "/dqstore?n=mystore1&a=clear", 
			 `{"code":400,"message":"store not exist - mystore1","data":"/dqstore?n=mystore1\u0026a=clear"}`+"\n",
			http.StatusBadRequest},
			{"PUT", "/dqstore?store1name=mystore&a=clear", 
			 `{"code":400,"message":"n empty","data":"/dqstore?store1name=mystore\u0026a=clear"}`+"\n",
			http.StatusBadRequest},
			{"PUT", "/dqstore?n=mystore&a=clear", 
			 `{"code":200,"message":"success"}`+"\n",
			http.StatusOK},
			{"DELETE", "/dqstore?n=mystore", 
			 `{"code":200,"message":"success"}`+"\n",
			http.StatusOK},
		}
		for _, test := range tests {
			res,err := ghttp.ClientHttp(test.method,  rootUrl+test.url, "text/html; charset=utf-8", nil)
			tResult(t, res, err, test.want, test.code)
		}		
	})
	t.Run("StoreGet", func(t *testing.T) {
		tests := []struct {
			method string
			url string
			want string
			code int
		}{
			{"POST", "/dqstore?n=mystore&t=string", 
			 `{"code":201,"message":"success"}`+"\n",
			http.StatusCreated},
			{"GET", "/dqstore?n=mystore&a=size1", 
			 `{"code":400,"message":"size1 is not in [size]","data":"/dqstore?n=mystore\u0026a=size1"}`+"\n",
			http.StatusBadRequest},
			{"GET", "/dqstore?n=mystore&acti1on=size", 
			 `{"code":400,"message":"a emtpy - [size]","data":"/dqstore?n=mystore\u0026acti1on=size"}`+"\n",
			http.StatusBadRequest},
			{"GET", "/dqstore?n=mystore1&a=size", 
			 `{"code":400,"message":"store not exist - mystore1","data":"/dqstore?n=mystore1\u0026a=size"}`+"\n",
			http.StatusBadRequest},
			{"GET", "/dqstore?store1name=mystore&a=size", 
			 `{"code":400,"message":"n empty","data":"/dqstore?store1name=mystore\u0026a=size"}`+"\n",
			http.StatusBadRequest},
			{"GET", "/dqstore?n=mystore&a=size", 
			 `{"code":200,"message":"success","data":0}`+"\n",
			http.StatusOK},
			{"DELETE", "/dqstore?n=mystore", 
			 `{"code":200,"message":"success"}`+"\n",
			http.StatusOK},
		}
		for _, test := range tests {
			res,err := ghttp.ClientHttp(test.method,  rootUrl+test.url, "text/html; charset=utf-8", nil)
			tResult(t, res, err, test.want, test.code)
		}
	})
}


func tDqDataTest(t *testing.T, rootUrl string) {

	t.Run("StorePost", func(t *testing.T) {
		tests := []struct {
			method string
			url string
			want string
			code int
		}{
			{"POST", "/dqstore?n=mystore&t=string", 
			 `{"code":201,"message":"success"}`+"\n",
			http.StatusCreated},
			{"POST", "/dqstore?n=intstore&t=int", 
			 `{"code":201,"message":"success"}`+"\n",
			http.StatusCreated},
			{"POST", "/dqstore?n=floatstore&t=float", 
			 `{"code":201,"message":"success"}`+"\n",
			http.StatusCreated},
		}
		for _, test := range tests {
			res,err := ghttp.ClientHttp(test.method,  rootUrl+test.url, "text/html; charset=utf-8", nil)
			tResult(t, res, err, test.want, test.code)
		}
	})
	t.Run("DataPostPutString", func(t *testing.T) {
		tests := []struct {
			method string
			url string
			want string
			code int
		}{
			{"POST", "/dqdata?n=mystore&a=f&v=a", 
			 `{"code":201,"message":"success"}`+"\n",
			http.StatusCreated},
			{"POST", "/dqdata?n=mystore&a=f&v=b", 
			 `{"code":201,"message":"success"}`+"\n",
			http.StatusCreated},
			{"POST", "/dqdata?n=mystore&a=b&v=c", 
			 `{"code":201,"message":"success"}`+"\n",
			http.StatusCreated},
			{"POST", "/dqdata?n=mystore&a=b&v=d", 
			 `{"code":201,"message":"success"}`+"\n",
			http.StatusCreated},
			{"PUT", "/dqdata?n=mystore&a=f", 
			 `{"code":200,"message":"success","data":"b"}`+"\n",
			http.StatusOK},
			{"PUT", "/dqdata?n=mystore&a=f", 
			 `{"code":200,"message":"success","data":"a"}`+"\n",
			http.StatusOK},
			{"PUT", "/dqdata?n=mystore&a=b", 
			 `{"code":200,"message":"success","data":"d"}`+"\n",
			http.StatusOK},
			{"PUT", "/dqdata?n=mystore&a=b", 
			 `{"code":200,"message":"success","data":"c"}`+"\n",
			http.StatusOK},
			{"PUT", "/dqdata?n=mystore&a=f", 
			 `{"code":400,"message":"store empty - mystore","data":"/dqdata?n=mystore\u0026a=f"}`+"\n",
			http.StatusBadRequest},
			{"PUT", "/dqdata?n=mystore&a=b", 
			 `{"code":400,"message":"store empty - mystore","data":"/dqdata?n=mystore\u0026a=b"}`+"\n",
			http.StatusBadRequest},
		}
		for _, test := range tests {
			res,err := ghttp.ClientHttp(test.method,  rootUrl+test.url, "text/html; charset=utf-8", nil)
			tResult(t, res, err, test.want, test.code)
		}
	})
	t.Run("DataPostPutInt", func(t *testing.T) {
		tests := []struct {
			method string
			url string
			want string
			code int
		}{
			{"POST", "/dqdata?n=intstore&a=f&v=1", 
			 `{"code":201,"message":"success"}`+"\n",
			http.StatusCreated},
			{"POST", "/dqdata?n=intstore&a=f&v=2", 
			 `{"code":201,"message":"success"}`+"\n",
			http.StatusCreated},
			{"POST", "/dqdata?n=intstore&a=b&v=0", 
			 `{"code":201,"message":"success"}`+"\n",
			http.StatusCreated},
			{"POST", "/dqdata?n=intstore&a=b&v=-1", 
			 `{"code":201,"message":"success"}`+"\n",
			http.StatusCreated},
			{"PUT", "/dqdata?n=intstore&a=f", 
			 `{"code":200,"message":"success","data":2}`+"\n",
			http.StatusOK},
			{"PUT", "/dqdata?n=intstore&a=f", 
			 `{"code":200,"message":"success","data":1}`+"\n",
			http.StatusOK},
			{"PUT", "/dqdata?n=intstore&a=b", 
			 `{"code":200,"message":"success","data":-1}`+"\n",
			http.StatusOK},
			{"PUT", "/dqdata?n=intstore&a=b", 
			 `{"code":200,"message":"success","data":0}`+"\n",
			http.StatusOK},
		}
		for _, test := range tests {
			res,err := ghttp.ClientHttp(test.method,  rootUrl+test.url, "text/html; charset=utf-8", nil)
			tResult(t, res, err, test.want, test.code)
		}
	})	
	t.Run("DataPostPutFloat", func(t *testing.T) {
		tests := []struct {
			method string
			url string
			want string
			code int
		}{
			{"POST", "/dqdata?n=floatstore&a=f&v=1.1", 
			 `{"code":201,"message":"success"}`+"\n",
			http.StatusCreated},
			{"POST", "/dqdata?n=floatstore&a=f&v=2.2", 
			 `{"code":201,"message":"success"}`+"\n",
			http.StatusCreated},
			{"POST", "/dqdata?n=floatstore&a=b&v=0.5", 
			 `{"code":201,"message":"success"}`+"\n",
			http.StatusCreated},
			{"POST", "/dqdata?n=floatstore&a=b&v=-1.5", 
			 `{"code":201,"message":"success"}`+"\n",
			http.StatusCreated},
			{"PUT", "/dqdata?n=floatstore&a=f", 
			 `{"code":200,"message":"success","data":2.2}`+"\n",
			http.StatusOK},
			{"PUT", "/dqdata?n=floatstore&a=f", 
			 `{"code":200,"message":"success","data":1.1}`+"\n",
			http.StatusOK},
			{"PUT", "/dqdata?n=floatstore&a=b", 
			 `{"code":200,"message":"success","data":-1.5}`+"\n",
			http.StatusOK},
			{"PUT", "/dqdata?n=floatstore&a=b", 
			 `{"code":200,"message":"success","data":0.5}`+"\n",
			http.StatusOK},
		}
		for _, test := range tests {
			res,err := ghttp.ClientHttp(test.method,  rootUrl+test.url, "text/html; charset=utf-8", nil)
			tResult(t, res, err, test.want, test.code)
		}
	})	
}

func tDqScaleTest(t *testing.T, totalKeys, repeats int, rootUrl string) {
	key := "abcdefghijklmnopqrstuvwxyz"
	str := strings.Repeat(key,repeats)
	t.Run("DataWithValueInQuery", func(t *testing.T) {
		tDqScaleTestRun(t, totalKeys, str, rootUrl, true)
	})
	t.Run("DataWithValueInBody", func(t *testing.T) {
		tDqScaleTestRun(t, totalKeys, str, rootUrl, false)
	})
}

func tDqScaleTestRun(t *testing.T, totalKeys int, str, rootUrl string, inQuery bool) {
	t.Run("PushFrontPopBack", func(t *testing.T) {
		want := `{"code":201,"message":"success"}`+"\n"
		for i:=0; i<totalKeys; i++ {
			item := str+global.Int2StrPadZero(i, 10)
			if inQuery {
				url := rootUrl+"/dqdata?n=mystore&a=f&v="+ item
				res,err := http.Post(url, "text/html; charset=utf-8", nil)				
				tResult(t, res, err, want, http.StatusCreated)
			} else {
				url := rootUrl+"/dqdata?n=mystore&a=f"
				res,err := http.Post(url, "text/html; charset=utf-8", strings.NewReader(item))
				tResult(t, res, err, want, http.StatusCreated)
			}
		}
		want = `{"code":200,"message":"success","data":"`
		for i:=0; i<totalKeys; i++ {
			item := str+global.Int2StrPadZero(i, 10)
			url := rootUrl+"/dqdata?n=mystore&a=b"
			res,err := ghttp.Put(url, "text/html; charset=utf-8", nil)
			tResult(t, res, err, want+item+`"}`+"\n", http.StatusOK)
		}
	})	
	t.Run("PushBackPopFront", func(t *testing.T) {
		want := `{"code":201,"message":"success"}`+"\n"
		for i:=0; i<totalKeys; i++ {
			item := str+global.Int2StrPadZero(i, 10)
			if inQuery {
				url := rootUrl+"/dqdata?n=mystore&a=b&v="+ item
				res,err := http.Post(url, "text/html; charset=utf-8", nil)
				tResult(t, res, err, want, http.StatusCreated)
			} else {
				url := rootUrl+"/dqdata?n=mystore&a=b"
				res,err := http.Post(url, "text/html; charset=utf-8", strings.NewReader(item))
				tResult(t, res, err, want, http.StatusCreated)
			}
		}
		want = `{"code":200,"message":"success","data":"`
		for i:=0; i<totalKeys; i++ {
			item := str+global.Int2StrPadZero(i, 10)
			url := rootUrl+"/dqdata?n=mystore&a=f"
			res,err := ghttp.Put(url, "text/html; charset=utf-8", nil)
			tResult(t, res, err, want+item+`"}`+"\n", http.StatusOK)
		}
	})	
	t.Run("PushFrontBackPopFrontBack", func(t *testing.T) {
		want := `{"code":201,"message":"success"}`+"\n"
		var url string
		for i:=0; i<totalKeys; i++ {
			item := str+global.Int2StrPadZero(i, 10)
			if inQuery {
				url = rootUrl+"/dqdata?n=mystore&a=f&v="+ item
				res,err := http.Post(url, "text/html; charset=utf-8", nil)
				tResult(t, res, err, want, http.StatusCreated)
				url = rootUrl+"/dqdata?n=mystore&a=b&v="+ item
				res,err = http.Post(url, "text/html; charset=utf-8", nil)
				tResult(t, res, err, want, http.StatusCreated)			
			} else {
				url = rootUrl+"/dqdata?n=mystore&a=f"
				res,err := http.Post(url, "text/html; charset=utf-8", strings.NewReader(item))
				tResult(t, res, err, want, http.StatusCreated)
				url = rootUrl+"/dqdata?n=mystore&a=b"
				res,err = http.Post(url, "text/html; charset=utf-8", strings.NewReader(item))
				tResult(t, res, err, want, http.StatusCreated)							
			}
		}
		want = `{"code":200,"message":"success","data":"`
		for i:=0; i<totalKeys; i++ {
			item := str+global.Int2StrPadZero(totalKeys-i-1, 10)
			url = rootUrl+"/dqdata?n=mystore&a=f"
			res,err := ghttp.Put(url, "text/html; charset=utf-8", nil)
			tResult(t, res, err, want+item+`"}`+"\n", http.StatusOK)
			url = rootUrl+"/dqdata?n=mystore&a=b"
			res,err = ghttp.Put(url, "text/html; charset=utf-8", nil)
			tResult(t, res, err, want+item+`"}`+"\n", http.StatusOK)
		}
	})	
}

