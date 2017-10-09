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
	"encoding/json"
	"io"
	"io/ioutil"
	"github.com/dlmc/ids/internal/oset"
	"github.com/dlmc/ids/global"
	"github.com/dlmc/golight/ghttp"
)


func tRangeTest(t *testing.T, rootUrl string) {

	t.Run("CreateStore", func(t *testing.T) {
		url := rootUrl+"/kvstore?n=mystore&st=oset&t=str"
		res,err := http.Post(url, "text/html; charset=utf-8", nil)
		want := `{"code":201,"message":"success"}`+"\n"
		tResult(t, res, err, want, http.StatusCreated)
	})
	t.Run("Range", func(t *testing.T) {
		total := 100
		want := `{"code":201,"message":"success"}`+"\n"
		sCode := http.StatusCreated
		t.Run("PostData", func(t *testing.T) {
			for i:=0; i<total; i++ {
				k := global.Int2StrPadZero(i, 10)
				url := rootUrl+"/kvdata?n=mystore&k="+ k +"&v="+ k
				res,err := http.Post(url, "text/html; charset=utf-8", nil)
				tResult(t, res, err, want, sCode)
			}
		})
		sCode = http.StatusOK
		t.Run("RangeGetDataBack", func(t *testing.T) {
			for i:=0; i<5; i++ {
				k := global.Int2StrPadZero(i, 10)+"~"+global.Int2StrPadZero(i+10, 10)
				url := rootUrl+"/kvdata/range?n=mystore&range="+ k +"&limit=10&ascending=true"
				res,err := http.Get(url)
				var data []string
				for j:=i; j<i+10; j++ {
					data = append(data, global.Int2StrPadZero(j, 10))
				}
				wd, _ := json.Marshal(data)
				want = `{"code":200,"message":"success","data":{"Count":10,"Values":`+string(wd)+`}}`+"\n"
				tResult(t, res, err, want, sCode)
			}
		})
	})
	t.Run("Scale", func(t *testing.T) {
		tRangePerformance(t, rootUrl)
	})
	t.Run("DeleteStore", func(t *testing.T) {
		want := `{"code":200,"message":"success"}`+"\n"
		res,err := ghttp.Del(rootUrl+"/kvstore?n=mystore", "text/html; charset=utf-8", nil)
		tResult(t, res, err, want, http.StatusOK)
	})
}

func tRandomRangeRequest(t *testing.T, ost *oset.OSet, nRequest, nReqRange, nTotalKeys int, key, str string) {
		for i:=0; i < nRequest; i++ {
			k := rand.Intn(nTotalKeys-nReqRange)
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

func tRandomRangeRequestHttpDiscard(t *testing.T, nRequest, nReqRange, nTotalKeys int, key, str, rootUrl string) {
	for i:=0; i < nRequest; i++ {
		k := rand.Intn(nTotalKeys-nReqRange)
		key := global.Int2StrPadZero(k,10)+key+"~"+global.Int2StrPadZero(k+nReqRange,10)+key
		u := rootUrl+"/kvdata/range?n=mystore&range="+ key +"&limit=10&ascending=true"
		resp,err := http.Get(u)
		if err != nil {
			t.Errorf("tResult failed", err)
		}
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}	
}
func tRandomRangeRequestHttpCheckResp(t *testing.T, nRequest, nReqRange, nTotalKeys int, key, str, rootUrl string) {
	code := http.StatusOK
	for i:=0; i < nRequest; i++ {
		k := rand.Intn(nTotalKeys-nReqRange)
		key := global.Int2StrPadZero(k,10)+key+"~"+global.Int2StrPadZero(k+nReqRange,10)+key
		u := rootUrl+"/kvdata/range?n=mystore&range="+ key +"&limit=10&ascending=true"
		resp,err := http.Get(u)
		var data []string
		for j:=k; j<k+nReqRange; j++ {
			data = append(data, global.Int2StrPadZero(j, 10)+str)
		}
		wd, _ := json.Marshal(data)
		want := `{"code":200,"message":"success","data":{"Count":10,"Values":`+string(wd)+`}}`+"\n"
		tResult(t, resp, err, want, code)
	}	
}


func tRunRangeTest(t *testing.T, nTotalKeys int, a []int, key, str string) {
	ost, _ := oset.New("oset", "str")
	t.Run("Create elements", func(t *testing.T) {
		a = a[:nTotalKeys]
		for _, v := range a {
			vpad := global.Int2StrPadZero(v,10)
			k := vpad+key
			value := vpad+str
			ost.Create(k, value)
		}
	})
	t.Run("1000 random range request with 10 elements", func(t *testing.T) {
		tRandomRangeRequest(t, ost, 1000, 10, nTotalKeys, key, str)
	})
}

func tRunRangeTestHttp(t *testing.T, nTotalKeys int, a []int, key, str, rootUrl string) {
	want := `{"code":201,"message":"success"}`+"\n"
	code := http.StatusCreated

	t.Run("Clear Store", func(t *testing.T) {
		res,err := ghttp.Put(rootUrl +"/kvstore?n=mystore&a=clear", "", nil)
		tResult(t, res, err, `{"code":200,"message":"success"}`+"\n", http.StatusOK)
	})
	
	t.Run("Create elements", func(t *testing.T) {
		a = a[:nTotalKeys]
		for _, v := range a {
			kpad := global.Int2StrPadZero(v,10)
			k := kpad+key
			value := kpad+str
			u := rootUrl +"/kvdata?n=mystore&k="+k+"&v="+value
			res,err := http.Post(u, "text/html; charset=utf-8", nil)
			tResult(t, res, err, want, code)
		}
	})
	t.Run("1000 random range request with 10 elements", func(t *testing.T) {
		t.Run("Check Response", func(t *testing.T) {
			tRandomRangeRequestHttpCheckResp(t, 1000, 10, nTotalKeys, key, str, rootUrl)
		})
		t.Run("Discard Response", func(t *testing.T) {
			tRandomRangeRequestHttpDiscard(t, 1000, 10, nTotalKeys, key, str, rootUrl)
		})
	})
}




func tRangePerformance20000Elements(t *testing.T, totalKeys, repeats int, rootUrl string) {
	key := "abcdefghijklmnopqrstuvwxyz"
	str := strings.Repeat(key,repeats)
	rand.Seed(time.Now().UTC().UnixNano())
	a := rand.Perm(totalKeys)
	t.Run("DataGetFunctionCall", func(t *testing.T) {
		tRunRangeTest(t, totalKeys, a, key, str)
	})
	t.Run("DataGetHttp", func(t *testing.T) {
		tRunRangeTestHttp(t, totalKeys, a, key, str, rootUrl)
	})
}



func tRangePerformance(t *testing.T, rootUrl string) {	
	t.Run("20000x26 Elements", func(t *testing.T) {
		tRangePerformance20000Elements(t, 20000, 1, rootUrl)
	})
	t.Run("20000x130 Elements", func(t *testing.T) {
		tRangePerformance20000Elements(t, 20000, 5, rootUrl)
	})
	t.Run("20000x1.3k Elements", func(t *testing.T) {
		tRangePerformance20000Elements(t, 20000, 50, rootUrl)
	})
	t.Run("20000x2.6k Elements", func(t *testing.T) {
		tRangePerformance20000Elements(t, 20000, 100, rootUrl)
	})
}





