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

const (
	strStoreNameEmpty = "storename empty"
	strStoreExist = "store exists - "
	strStoreNotExist = "store does not exists - "
	strStoreGetActionEmpty = "action emtpy - [size, keys, values]"
	strStoreGetActionNotExist = " is not in [size, keys, values]"
	strStorePutActionEmpty = "action empty - [clear]"
	strStorePutActionNotExit = " is not in [clear]"
	strKeyEmpty = "key empty"
	strKeyExist = "key exists - "
	strKeyNotExist = "key does not exist - "
)

func decodeRequestNSK(h *ghttp.Http) (sn string, st global.StoreType, kt global.KeyType, err error) {
	query := h.Query
	if sn = query.Get("storename"); sn == "" {
		err = errors.New(strStoreNameEmpty)
	} else if st, err = global.GetStoreType(query.Get("storetype")); err == nil {
		kt, err = global.GetKeyType(query.Get("keytype"))
	}
	return
}

func decodeRequestNKV(h *ghttp.Http) (sn, k, v string, err error) {
	query := h.Query
	if sn = query.Get("storename"); sn == "" {
		err = errors.New(strStoreNameEmpty)
	} else if k = query.Get("k"); k == "" {
		err = errors.New(strKeyEmpty)
	} else if v = query.Get("v"); v == "" {
		if b, e := ioutil.ReadAll(h.R.Body); e == nil {
			v = string(b)
		} else {
			err = errors.New("Error reading body: " + err.Error())
        }
	}
	return
}

func decodeRequestNK(h *ghttp.Http) (sn, k string, err error) {
	query := h.Query
	if sn = query.Get("storename"); sn == "" {
		err = errors.New(strStoreNameEmpty)
	} else if k = query.Get("k"); k == "" {
		err = errors.New(strKeyEmpty)
	}
	return
}

func decodeRequestNA(h *ghttp.Http, emptyMsg string)  (sn, action string, err error) {
	query := h.Query
	if sn = query.Get("storename"); sn == "" {
		err = errors.New(strStoreNameEmpty)
	} else if action = query.Get("action"); action == "" {
		err = errors.New(emptyMsg)
	}
	return
}


func  decodeRequestN(h *ghttp.Http) (sn string, err error) {
	query := h.Query
	if sn = query.Get("storename"); sn == "" {
		err = errors.New(strStoreNameEmpty)
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
