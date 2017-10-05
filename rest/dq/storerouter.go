// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dq

import (
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

// GET ?storename=sname&action=size
type storeGet struct {}
func (s *storeGet) ServeHTTPWithCtx(c ghttp.Ctx, h *ghttp.Http) ghttp.Ctx {
	reqUri := h.R.RequestURI
	resp := h.Resp
	log := h.Log.Str("Get",reqUri).Logger()
	sn, action, err := decodeRequestNA(h, strStoreGetActionEmpty)
	if err == nil {
		if action != "size" {
			err = errors.New(action + strStoreGetActionNotExist)
		} else if st, found := store.GetStore(sn); found {
			resp.Data = st.Size()
		} else {
			err = errors.New(strStoreNotExist + sn)
		}
	}
	finishHandling(err,  resp, log, reqUri, http.StatusOK)
	return c
}

// POST ?storename=sname  
type storePost struct {}
func (s *storePost) ServeHTTPWithCtx(c ghttp.Ctx, h *ghttp.Http) ghttp.Ctx {
	reqUri := h.R.RequestURI
	log := h.Log.Str("Post", reqUri).Logger()
	sn, err := decodeRequestN(h)
	if err == nil {
		if ok := store.CreateStore(sn); !ok {
			err = errors.New(strStoreExist + sn)
		}
	}
	finishHandling(err,  h.Resp, log, reqUri, http.StatusCreated)
	return c
}

// PUT ?storename=sname&action=clear
type storePut struct {}
func (s *storePut) ServeHTTPWithCtx(c ghttp.Ctx, h *ghttp.Http) ghttp.Ctx {
	reqUri := h.R.RequestURI
	log := h.Log.Str("Put", reqUri).Logger()
	sn, action, err := decodeRequestNA(h, strStorePutActionEmpty)
	if err == nil {
		if action != "clear" {
			err = errors.New(action + strStorePutActionNotExit)
		} else if st, found := store.GetStore(sn); found {
			st.Clear()
		} else {
			err = errors.New(strStoreNotExist + sn)
		}
	}
	finishHandling(err,  h.Resp, log, reqUri, http.StatusOK)
	return c
}

// DELETE ?storename=sname
type storeDelete struct {}
func (s *storeDelete) ServeHTTPWithCtx(c ghttp.Ctx, h *ghttp.Http) ghttp.Ctx {
	reqUri := h.R.RequestURI
	log := h.Log.Str("Delete",reqUri).Logger()
	sn, err := decodeRequestN(h)
	if err == nil {		
		if _, found := store.GetStore(sn); found {
			store.DeleteStore(sn)
		} else {
			err = errors.New(strStoreNotExist + sn)
		}
	}
	finishHandling(err, h.Resp, log, reqUri, http.StatusOK)
	return c
}

