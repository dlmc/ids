// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dq

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

// GET ?n=sname&a=size
type storeGet struct {}
func (s *storeGet) ServeHTTPWithCtx(c ghttp.Ctx, h *ghttp.Http) ghttp.Ctx {
	reqUri := h.R.RequestURI
	resp := h.Resp
	log := h.Log.Str("Get",reqUri).Logger()
	n, a, err := decodeRequestNA(h, global.StrDqStoreGetActionEmpty)
	if err == nil {
		if a != "size" {
			err = errors.New(a + global.StrDqStoreGetActionNotExist)
		} else if st, found := store.GetStore(n); found {
			resp.Data = st.Size()
		} else {
			err = errors.New(global.StrStoreNotExist + n)
		}
	}
	finishHandling(err,  resp, log, reqUri, http.StatusOK)
	return c
}

// POST ?n=sname&t=int|float|string
type storePost struct {}
func (s *storePost) ServeHTTPWithCtx(c ghttp.Ctx, h *ghttp.Http) ghttp.Ctx {
	reqUri := h.R.RequestURI
	log := h.Log.Str("Post", reqUri).Logger()
	n, t, err := decodeRequestNT(h)
	if err == nil {
		err = store.CreateStore(n, t)
	}
	finishHandling(err,  h.Resp, log, reqUri, http.StatusCreated)
	return c
}

// PUT ?n=sname&a=clear
type storePut struct {}
func (s *storePut) ServeHTTPWithCtx(c ghttp.Ctx, h *ghttp.Http) ghttp.Ctx {
	reqUri := h.R.RequestURI
	log := h.Log.Str("Put", reqUri).Logger()
	n, a, err := decodeRequestNA(h, global.StrDqStorePutActionEmpty)
	if err == nil {
		if a != "clear" {
			err = errors.New(a + global.StrDqStorePutActionNotExit)
		} else if st, found := store.GetStore(n); found {
			st.Clear()
		} else {
			err = errors.New(global.StrStoreNotExist + n)
		}
	}
	finishHandling(err,  h.Resp, log, reqUri, http.StatusOK)
	return c
}

// DELETE ?n=sname
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

