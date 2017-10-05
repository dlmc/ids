// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kv

import (
	. "github.com/dlmc/ids/global"
	"github.com/dlmc/golight/ghttp"
	"net/http"
	"strings"
	"errors"
)


var RangeRouter = ghttp.Router{
	"GET":&rangeGet{}, 
//	"DELETE":&rangeDelete{},
}

// GET /data/range?storename=sname&range=skey~ekey&limit=number&ascending=true
type rangeGet struct {
	sn, sk, ek string
	ascending bool
	limit int64
	Count int
	Values []interface{}
}

func (s *rangeGet) decodeRequest(h *ghttp.Http) (err error) {
	query := h.Query
	var ok bool
	if s.sn = query.Get("storename"); s.sn == "" {
		err = errors.New(strStoreNameEmpty)
	} else if qrng := query.Get("range"); qrng == "" {
		err = errors.New("either k or range is required")
	} else if rng := strings.Split(qrng, "~"); rng[0] == "" || rng[1] == "" {
		err = errors.New("range should be in [startkey~endkey]")
	} else if limit := query.Get("limit"); limit == "" {
		err = errors.New("limit [] empty")
	} else if s.limit, ok = ParseInt([]byte(limit)); !ok {
		err = errors.New("limit not integer")			
	} else {
		if ascending := query.Get("ascending"); ascending == "true" {
			s.ascending = true
		} else {
			s.ascending = false
		}
		s.sk=rng[0]
		s.ek=rng[1]
	}
	return
}
func (s *rangeGet) ServeHTTPWithCtx(c ghttp.Ctx, h *ghttp.Http) ghttp.Ctx {
	reqUri := h.R.RequestURI
	resp := h.Resp
	log := h.Log.Str("Get", reqUri).Logger()
	err := s.decodeRequest(h)

	if err == nil {
		if st, found := store.GetStore(s.sn); !found {
			err = errors.New(strStoreNotExist + s.sn)
		} else {
			s.Values, s.Count = st.RangeGet(s.sk, s.ek, int(s.limit), s.ascending)
			resp.Data = s
		} 
	}
	finishHandling(err,  resp, log, reqUri, http.StatusOK)
	return c
}
