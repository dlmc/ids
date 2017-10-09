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

var DataRouter = ghttp.Router{
	"GET":&dataGet{}, 
	"POST":&dataPost{}, 
	"PUT":&dataPut{},
	"DELETE":&dataDelete{},
}


// GET /data?n=sname&k=key
type dataGet struct {}
func (s *dataGet) ServeHTTPWithCtx(c ghttp.Ctx, h *ghttp.Http) ghttp.Ctx {
	reqUri := h.R.RequestURI
	resp := h.Resp
	log := h.Log.Str("Get", reqUri).Logger()
	n, k, err := decodeRequestNK(h)
	var ok bool

	if err == nil {
		if st, found := store.GetStore(n); !found {
			err = errors.New(global.StrStoreNotExist + n)
		} else if resp.Data, ok = st.Read(k); !ok {
			err = errors.New(global.StrKeyNotExist + k)
		}
	}
	finishHandling(err,  resp, log, reqUri, http.StatusOK)
	return c
}

// Create a new k/v pair in the store
// POST /data?n=sname&k=key&v=value 
// POST /data?n=sname&k=key and v in the body
type dataPost struct {}
func (d *dataPost) ServeHTTPWithCtx(c ghttp.Ctx, h *ghttp.Http) ghttp.Ctx {
	reqUri := h.R.RequestURI
	log := h.Log.Str("Post", reqUri).Logger()
	n, k, v, err := decodeRequestNKV(h)

	if err == nil {
		if st, found := store.GetStore(n); !found {
			err = errors.New(global.StrStoreNotExist + n)
		} else if ok := st.Create(k, v); !ok {
			err = errors.New(global.StrKeyExist + k)
		}
	}

	finishHandling(err, h.Resp, log, reqUri, http.StatusCreated)
	return c
}

// Replace an existing k/v pair
// PUT /data?n=sname&k=key&v=value
// PUT /data?n=sname&k=key and v in the body
type dataPut struct {}
func (d *dataPut) ServeHTTPWithCtx(c ghttp.Ctx, h *ghttp.Http) ghttp.Ctx {
	reqUri := h.R.RequestURI
	log := h.Log.Str("Put", reqUri).Logger()
	n, k, v, err := decodeRequestNKV(h)

	if err == nil {
		if st, found := store.GetStore(n); !found {
			err = errors.New(global.StrStoreNotExist + n)
		} else if ok := st.Update(k, v); !ok {
			err = errors.New(global.StrKeyNotExist + k)
		}
	}

	finishHandling(err,  h.Resp, log, reqUri, http.StatusOK)
	return c
}

// Delete a k/v pair
// DELETE /data?n=sname&k=key
type dataDelete struct {}
func (d *dataDelete) ServeHTTPWithCtx(c ghttp.Ctx, h *ghttp.Http) ghttp.Ctx {
	reqUri := h.R.RequestURI
	log := h.Log.Str("Delete", reqUri).Logger()
	n, k, err := decodeRequestNK(h)

	if err == nil {
		if st, found := store.GetStore(n); !found {
			err = errors.New(global.StrStoreNotExist + n)
		} else if ok := st.Delete(k); !ok {
			err = errors.New(global.StrKeyNotExist + k)
		}
	}
	finishHandling(err,  h.Resp, log, reqUri, http.StatusOK)
	return c
}
