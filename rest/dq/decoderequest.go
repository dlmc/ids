// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dq

import (
	"github.com/dlmc/ids/global"
	"github.com/dlmc/ids/dqstore"
	"github.com/dlmc/golight/ghttp"
	zlog "github.com/rs/zerolog"
	"net/http"
	"errors"
	"io/ioutil"
)

var store = dqstore.New()

func decodeRequestN(h *ghttp.Http) (sn string, err error) {
	query := h.Query
	if sn = query.Get(global.QueryName); sn == "" {
		err = errors.New(global.StrNameEmpty)
	}
	return
}

func decodeRequestNT(h *ghttp.Http) (n, t string, err error) {
	query := h.Query
	if n = query.Get(global.QueryName); n == "" {
		err = errors.New(global.StrNameEmpty)
	} else if t = query.Get(global.QueryType); t == "" {
		err = errors.New(global.StrTypeEmpty)
	}
	return
}
func decodeRequestNA(h *ghttp.Http, emptyMsg string)  (n, a string, err error) {
	query := h.Query
	if n = query.Get(global.QueryName); n == "" {
		err = errors.New(global.StrNameEmpty)
	} else if a = query.Get(global.QueryAction); a == "" {
		err = errors.New(emptyMsg)
	}
	return
}

func decodeRequestNAV(h *ghttp.Http, emptyMsg string)  (n, a, v string, err error) {
	query := h.Query
	if n = query.Get(global.QueryName); n == "" {
		err = errors.New(global.StrNameEmpty)
	} else if a = query.Get(global.QueryAction); a == "" {
		err = errors.New(emptyMsg)
	} else if v = query.Get(global.QueryValue); v == "" {
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
