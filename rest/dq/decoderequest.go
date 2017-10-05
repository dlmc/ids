// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dq

import (
	"github.com/dlmc/ids/dqstore"
	"github.com/dlmc/golight/ghttp"
	zlog "github.com/rs/zerolog"
	"net/http"
	"errors"
	"io/ioutil"
)

var store = dqstore.New()

const (
	strStoreEmpty = "store empty - "
	strStoreNameEmpty = "storename empty"
	strStoreExist = "store exists - "
	strStoreNotExist = "store does not exists - "
	strStoreGetActionEmpty = "action emtpy - [size]"
	strStoreGetActionNotExist = " is not in [size]"
	strStorePutActionEmpty = "action empty - [clear]"
	strStorePutActionNotExit = " is not in [clear]"
	strDataPutActionEmpty = "action empty - [f,b]"
	strDataPutActionNotExist = " is not in - [f,b]"
)

func decodeRequestN(h *ghttp.Http) (sn string, err error) {
	query := h.Query
	if sn = query.Get("storename"); sn == "" {
		err = errors.New(strStoreNameEmpty)
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

func decodeRequestNAV(h *ghttp.Http, emptyMsg string)  (sn, action, v string, err error) {
	query := h.Query
	if sn = query.Get("storename"); sn == "" {
		err = errors.New(strStoreNameEmpty)
	} else if action = query.Get("action"); action == "" {
		err = errors.New(emptyMsg)
	} else if v = query.Get("v"); v == "" {
		if b, e := ioutil.ReadAll(h.R.Body); e == nil {
			v = string(b)
		} else {
			err = errors.New("Error reading body: " + err.Error())
        }
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
