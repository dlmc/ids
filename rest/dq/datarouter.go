// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dq

import (
	"github.com/dlmc/golight/ghttp"
	"net/http"
	"errors"
)

var DataRouter = ghttp.Router{
	"POST":&dataPost{},
	"PUT":&dataPut{},
}

// Push item v to the front or back of the queue
// POST /data?storename=sname&v=value&action=f|b
// POST /data?storename=sname&action=f|b and v in body
type dataPost struct {}
func (d *dataPost) ServeHTTPWithCtx(c ghttp.Ctx, h *ghttp.Http) ghttp.Ctx {
	reqUri := h.R.RequestURI
	log := h.Log.Str("Post", reqUri).Logger()
	sn, action, v, err := decodeRequestNAV(h, strDataPutActionEmpty)
	if err == nil {
		if st, found := store.GetStore(sn); !found {
			err = errors.New(strStoreNotExist + sn)
		} else {
			switch action {
			case "f":
				st.PushFront(v)
			case "b":
				st.PushBack(v)
			default:
				err = errors.New(action + strDataPutActionNotExist)
			}
		}
	}

	finishHandling(err,  h.Resp, log, reqUri, http.StatusCreated)
	return c
}


// Pop and return the item from the front or back of the queue
// PUT ?storename=sname&action=f|b
// Shall we use DELETE instead or both?
type dataPut struct {}
func (d *dataPut) ServeHTTPWithCtx(c ghttp.Ctx, h *ghttp.Http) ghttp.Ctx {
	reqUri := h.R.RequestURI
	resp := h.Resp
	log := h.Log.Str("Put", reqUri).Logger()
	sn, action, err := decodeRequestNA(h, strDataPutActionEmpty)
	if err == nil {
		if st, found := store.GetStore(sn); !found {
			err = errors.New(strStoreNotExist + sn)
		} else {
			switch action {
			case "f":
				if resp.Data, found = st.PopFront(); !found {
					err = errors.New(strStoreEmpty + sn)
				} 
			case "b":
				if resp.Data, found = st.PopBack(); !found {
					err = errors.New(strStoreEmpty + sn)
				} 
			default:
				err = errors.New(action + strDataPutActionNotExist)
			}
		}
	}

	finishHandling(err,  h.Resp, log, reqUri, http.StatusOK)
	return c
}

