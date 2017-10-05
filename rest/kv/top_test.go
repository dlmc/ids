// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kv_test

import (
	"testing"
	"net/http/httptest"
	"github.com/dlmc/ids/rest"
	"github.com/dlmc/golight/ghttp"
	"github.com/dlmc/golight/decorator/logging"
	"os"
	"encoding/json"
	"fmt"
)



func ExampleJsonMarshal() {
	msg1, _ := json.Marshal(&ghttp.Response{202, "OK","stringing"})
	fmt.Println(string(msg1))
	// Output: 
	// {"code":202,"message":"OK","data":"stringing"}

	msg2, _ := json.Marshal(&ghttp.Response{202,"OK",nil})
	fmt.Println(string(msg2))
	// {"code":202,"message":"OK"}

	msg3, _ := json.Marshal(&ghttp.Response{})
	fmt.Println(string(msg3))
	// {}

	msgA := ghttp.Response{
		Code:202, 
		Message:"OK", 
		Data: ghttp.Response{
			Code:202, 
			Message:"OK", 
			Data: "fe",
		}}
	msg1a, _ := json.Marshal(&msgA)
	fmt.Println(string(msg1a))
	// {"code":202,"message":"OK","data":{"code":202,"message":"OK","data":"fe"}}
}

func TestRestData(t *testing.T) {
	l := logging.NewLogger(os.Stdout).Level(logging.LogError).With()
	ts := httptest.NewServer(rest.NewServerMux("/", l))
	url := ts.URL
	defer ts.Close()
	t.Run("set", func(t *testing.T) {
		tOsetTest(t, url, "set")
	})
	t.Run("oset", func(t *testing.T) {
		tOsetTest(t, url, "oset")
	})
	t.Run("osetRange", func(t *testing.T) {
		tRangeTest(t, url)
	})

}

