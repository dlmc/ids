// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kv

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

// GET /store?storename=sname&action=size|keys|values 
type storeGet struct {}
func (s *storeGet) ServeHTTPWithCtx(c ghttp.Ctx, h *ghttp.Http) ghttp.Ctx {
	reqUri := h.R.RequestURI
	resp := h.Resp
	log := h.Log.Str("Get",reqUri).Logger()
	sn, action, err := decodeRequestNA(h, strStoreGetActionEmpty)
	if err == nil {
		if st, found := store.GetStore(sn); !found {
			err = errors.New(strStoreNotExist + sn)
		} else {
			switch action {
			case "size":
				resp.Data = st.Size()
			case "keys":
				resp.Data = st.Keys()
			case "values":
				resp.Data = st.Values()
			default:
				err = errors.New(action + strStoreGetActionNotExist)
			}
		}
	}
	finishHandling(err,  resp, log, reqUri, http.StatusOK)
	return c
}

// POST /store?storename=sname&keytype=string|int|float&storetype=oset  
type storePost struct {}
func (s *storePost) ServeHTTPWithCtx(c ghttp.Ctx, h *ghttp.Http) ghttp.Ctx {
	reqUri := h.R.RequestURI
	log := h.Log.Str("Post", reqUri).Logger()
	sn, st, kt, err := decodeRequestNSK(h)
	if err == nil {
		//kt ignored, all uses string key for now
		if ok := store.CreateStore(sn, st, kt); !ok {
			err = errors.New(strStoreExist + sn)
		}
	}
	finishHandling(err,  h.Resp, log, reqUri, http.StatusCreated)
	return c
}

// PUT /store?storename=sname&action=clear 
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

// DELETE /store?storename=sname
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

