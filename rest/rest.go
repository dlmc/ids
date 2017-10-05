// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rest


import (
	"net/http"
	"strings"
	"github.com/dlmc/ids/rest/kv"
	"github.com/dlmc/ids/rest/dq"
	"github.com/dlmc/golight/decorator"
	"github.com/dlmc/golight/decorator/logging"
	"github.com/dlmc/golight/decorator/respond"
	log "github.com/rs/zerolog"
)


// path - "/"
func NewServerMux(path string, lc log.Context) *http.ServeMux {
	mux := http.NewServeMux()
	
	lc = lc.Str("service","IDS")
	
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	
	rd := respond.CreateDecor()
	ld := logging.CreateDecor(lc)
	
	mux.Handle(path + "kvdata",       decorator.DecorateRouter(kv.DataRouter,  rd, ld))
	mux.Handle(path + "kvdata/range", decorator.DecorateRouter(kv.RangeRouter, rd, ld))
	mux.Handle(path + "kvstore",      decorator.DecorateRouter(kv.StoreRouter, rd, ld))
	mux.Handle(path + "dqdata",       decorator.DecorateRouter(dq.DataRouter,  rd, ld))
	mux.Handle(path + "dqstore",      decorator.DecorateRouter(dq.StoreRouter, rd, ld))
	return mux
}
