// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kv

import (
	"github.com/dlmc/ids/kvstore"
	"github.com/dlmc/ids/global"
	"github.com/dlmc/golight/ghttp"
	zlog "github.com/rs/zerolog"
	"net/http"
	"errors"
	"io/ioutil"
)

var store = kvstore.New()


func decodeRequestNSK(h *ghttp.Http) (st, n, t string, err error) {
	query := h.Query
	if n = query.Get(global.QueryName); n == "" {
		err = errors.New(global.StrNameEmpty)
	} else if st = query.Get(global.QueryStoreType); st == "" {
		err = errors.New(global.StrStoreTypeEmpty)
	} else if t = query.Get(global.QueryType); t == "" {
		err = errors.New(global.StrTypeEmpty)
	}
	return
}

func decodeRequestNKV(h *ghttp.Http) (sn, k, v string, err error) {
	query := h.Query
	if sn = query.Get(global.QueryName); sn == "" {
		err = errors.New(global.StrNameEmpty)
	} else if k = query.Get(global.QueryKey); k == "" {
		err = errors.New(global.StrKeyEmpty)
	} else if v = query.Get(global.QueryValue); v == "" {
		if b, e := ioutil.ReadAll(h.R.Body); e == nil {
			v = string(b)
		} else {
			err = errors.New("Error reading body: " + e.Error())
        }
	}
	return
}

func decodeRequestNK(h *ghttp.Http) (sn, k string, err error) {
	query := h.Query
	if sn = query.Get(global.QueryName); sn == "" {
		err = errors.New(global.StrNameEmpty)
	} else if k = query.Get(global.QueryKey); k == "" {
		err = errors.New(global.StrKeyEmpty)
	}
	return
}

func decodeRequestNA(h *ghttp.Http, emptyMsg string)  (sn, action string, err error) {
	query := h.Query
	if sn = query.Get(global.QueryName); sn == "" {
		err = errors.New(global.StrNameEmpty)
	} else if action = query.Get(global.QueryAction); action == "" {
		err = errors.New(emptyMsg)
	}
	return
}


func  decodeRequestN(h *ghttp.Http) (sn string, err error) {
	query := h.Query
	if sn = query.Get(global.QueryName); sn == "" {
		err = errors.New(global.StrNameEmpty)
	}
	return
}


func finishHandling(err error, resp  *ghttp.Response, log zlog.Logger, reqUrl string, codeOK int) {
	if err == nil {
		resp.Code = codeOK
		resp.Message = ghttp.SuccessMessage
		log.Info().Msg(resp.Message)
	} else {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		resp.Data = reqUrl
		log.Warn().Msg(resp.Message)
	}
}
