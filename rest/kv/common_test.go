// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kv_test

import (
	"testing"
	"net/http"
	"io/ioutil"
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
