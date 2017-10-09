// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dq

import (
	"github.com/dlmc/golight/ghttp"
	"github.com/dlmc/ids/global"
	"net/http"
	"errors"
)

var DataRouter = ghttp.Router{
	"POST":&dataPost{},
	"PUT":&dataPut{},
}

// Push item v to the front or back of the queue
// POST /data?n=sname&v=value&a=f|b
// POST /data?n=sname&a=f|b and v in body
type dataPost struct {}
func (d *dataPost) ServeHTTPWithCtx(c ghttp.Ctx, h *ghttp.Http) ghttp.Ctx {
	reqUri := h.R.RequestURI
	log := h.Log.Str("Post", reqUri).Logger()
	n, a, v, err := decodeRequestNAV(h, global.StrDqDataPutActionEmpty)
	var ok bool = true
	if err == nil {
		if st, found := store.GetStore(n); !found {
			err = errors.New(global.StrStoreNotExist + n)
		} else {
			switch a {
			case "f":
				ok = st.PushFront(v)
			case "b":
				ok = st.PushBack(v)
			default:
				err = errors.New(a + global.StrDqDataPutActionNotExist)
			}
		}
	}

	if !ok {
		err = errors.New(global.StrDataParseError + v)
	}
	finishHandling(err,  h.Resp, log, reqUri, http.StatusCreated)
	return c
}


// Pop and return the item from the front or back of the queue
// PUT ?n=sname&a=f|b
// Shall we use DELETE instead or both?
type dataPut struct {}
func (d *dataPut) ServeHTTPWithCtx(c ghttp.Ctx, h *ghttp.Http) ghttp.Ctx {
	reqUri := h.R.RequestURI
	resp := h.Resp
	log := h.Log.Str("Put", reqUri).Logger()
	n, a, err := decodeRequestNA(h, global.StrDqDataPutActionEmpty)
	if err == nil {
		if st, found := store.GetStore(n); !found {
			err = errors.New(global.StrStoreNotExist + n)
		} else {
			switch a {
			case global.QueryActionQueFront:
				if resp.Data, found = st.PopFront(); !found {
					err = errors.New(global.StrStoreEmpty + n)
				} 
			case global.QueryActionQueBack:
				if resp.Data, found = st.PopBack(); !found {
					err = errors.New(global.StrStoreEmpty + n)
				} 
			default:
				err = errors.New(a + global.StrDqDataPutActionNotExist)
			}
		}
	}

	finishHandling(err,  h.Resp, log, reqUri, http.StatusOK)
	return c
}

