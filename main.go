// Copyright 2017 The Golight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package main

import (
	"github.com/dlmc/golight/decorator/logging"
	"github.com/dlmc/ids/rest"
	"flag"
	"net/http"
	"os"
)

var addr = flag.String("addr", ":8081", "Server address to run - ip:port")
var path = flag.String("path", "/", "Root path to service")


// addr - ":8081"
// path - "/"
func startService(addr, path string, lc log.Context) error {
	s := &http.Server{
		Addr:           addr,
		Handler:        newServerMux(path, l),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:		10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}	
	return s.ListenAndServe()
}


func main() {
	flag.Parse()	
	l := logging.NewContext(os.Stdout)
	
	e := startService(*addr, *path, l)

	//e := http.ListenAndServe(*addr, rest.NewServerMux(*path, l))
	l.Logger().Fatal().Msg(e.Error())
}