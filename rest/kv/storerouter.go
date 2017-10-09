// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kv

import (
	"github.com/dlmc/ids/global"
	"github.com/dlmc/golight/ghttp"
	"net/http"
	"errors"
)

var StoreRouter = ghttp.Router{
	"GET":&storeGet{}, 
	"POST":&storePost{}, 
	"PUT":&storePut{},
	"DELETE":&storeDelete{},
}

// GET /store?n=sname&a=size|keys|values 
type storeGet struct {}
func (s *storeGet) ServeHTTPWithCtx(c ghttp.Ctx, h *ghttp.Http) ghttp.Ctx {
	reqUri := h.R.RequestURI
	resp := h.Resp
	log := h.Log.Str("Get",reqUri).Logger()
	n, a, err := decodeRequestNA(h, global.StrKvStoreGetActionEmpty)
	if err == nil {
		if st, found := store.GetStore(n); !found {
			err = errors.New(global.StrStoreNotExist + n)
		} else {
			switch a {
			case "size":
				resp.Data = st.Size()
			case "keys":
				resp.Data = st.Keys()
			case "values":
				resp.Data = st.Values()
			default:
				err = errors.New(a + global.StrKvStoreGetActionNotExist)
			}
		}
	}
	finishHandling(err,  resp, log, reqUri, http.StatusOK)
	return c
}

// POST /store?n=sname&k=string|int|float&st=oset  
type storePost struct {}
func (s *storePost) ServeHTTPWithCtx(c ghttp.Ctx, h *ghttp.Http) ghttp.Ctx {
	reqUri := h.R.RequestURI
	log := h.Log.Str("Post", reqUri).Logger()
	st, n, t, err := decodeRequestNSK(h)
	if err == nil {
		err = store.CreateStore(n, st, t)
	}
	finishHandling(err,  h.Resp, log, reqUri, http.StatusCreated)
	return c
}

// PUT /store?storename=sname&action=clear 
type storePut struct {}
func (s *storePut) ServeHTTPWithCtx(c ghttp.Ctx, h *ghttp.Http) ghttp.Ctx {
	reqUri := h.R.RequestURI
	log := h.Log.Str("Put", reqUri).Logger()
	n, a, err := decodeRequestNA(h, global.StrKvStorePutActionEmpty)
	if err == nil {
		if a != "clear" {
			err = errors.New(a + global.StrKvStorePutActionNotExit)
		} else if st, found := store.GetStore(n); found {
			st.Clear()
		} else {
			err = errors.New(global.StrStoreNotExist + n)
		}
	}
	finishHandling(err,  h.Resp, log, reqUri, http.StatusOK)
	return c
}

// DELETE /store?storename=sname
type storeDelete struct {}
func (s *storeDelete) ServeHTTPWithCtx(c ghttp.Ctx, h *ghttp.Http) ghttp.Ctx {
	reqUri := h.R.RequestURI
	log := h.Log.Str("Delete",reqUri).Logger()
	n, err := decodeRequestN(h)
	if err == nil {		
		if _, found := store.GetStore(n); found {
			store.DeleteStore(n)
		} else {
			err = errors.New(global.StrStoreNotExist + n)
		}
	}
	finishHandling(err, h.Resp, log, reqUri, http.StatusOK)
	return c
}

