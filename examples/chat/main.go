// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"fmt"
	"io"
	"path/filepath"
)

var addr = flag.String("addr", ":"+os.Getenv("PORT"), "http service address")
var cert = flag.String("cert", "./algo.crt", "certificate to be used")
var key = flag.String("key", "./algo.key", "certificate key to be used")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if _, err := os.Stat("./home.html"); err == nil {
		http.ServeFile(w, r, "./home.html")
	} else {
		http.ServeFile(w, r, "examples/chat/home.html")
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World!")

	var files []string

    root := "./"
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        files = append(files, path)
        return nil
    })
    if err != nil {
        panic(err)
    }
    for _, file := range files {
		io.WriteString(w, file)
		io.WriteString(w, "\n")
    }
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", hello)
	http.HandleFunc("/home", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	var httpErr error
	if _, err := os.Stat(*cert); err == nil {
		fmt.Println(*cert, " found. Switching to https")
		if httpErr = http.ListenAndServeTLS(*addr, *cert, *key, nil); httpErr != nil {
			log.Fatal("The process exited with https error: ", httpErr.Error())
		}
	} else {
		fmt.Println("No cert, using http")
		httpErr = http.ListenAndServe(*addr, nil)
		if httpErr != nil {
			log.Fatal("The process exited with http error: ", httpErr.Error())
		}
	}
	/*flag.Parse()
	http.HandleFunc("/", hello)
	fmt.Println("listen on ", *addr)
	http.ListenAndServe(*addr, nil)*/
}
