// Copyright 2012 The Golanger Authors. All rights reserved.
// Use of this source code is governed by a GPLv3
// license that can be found in the LICENSE file.

// Compile The Go Source By play.golang.org.
package controllers

import (
    "fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	PLAY = "http://play.golang.org/p/%s"
	COMPILE = "http://play.golang.org/compile"
	FMT = "http://play.golang.org/fmt"
	SHARE = "http://play.golang.org/share"
)

type PagePlay struct {
	*App
}

type Compiled struct {
	Compile_errors string
	Output         string
}

func init() {
	Page.RegisterController("play/", &PagePlay{Page})
}

func (p *PagePlay) Index() {
	if p.Request.Method != "GET" {
		return
	}
    id := p.GET["p"]
    if id == "" {
        return
    }
    fmt.Println(id)
    play := fmt.Sprintf(PLAY, id)

    resp, _ := http.Get(play)
	defer resp.Body.Close()

	buf, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(buf))
}

func (p *PagePlay) Compile() {
	if p.Request.Method != "POST" {
		return
	}

	data := url.Values{"body": {strings.TrimSpace(p.POST["body"])}}
	resp, err := http.PostForm(COMPILE, data)

	defer resp.Body.Close()
	if err != nil {
		m := Compiled{"Error communicating with remote server.", ""}
		ret, _ := json.Marshal(m)
		p.ResponseWriter.Write(ret)
		return
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	p.ResponseWriter.Write(buf)
}

func (p *PagePlay) Fmt() {
	if p.Request.Method != "POST" {
		return
	}

	data := url.Values{"body": {strings.TrimSpace(p.POST["body"])}}
	resp, err := http.PostForm(FMT, data)

	defer resp.Body.Close()
	if err != nil {
		m := Compiled{"Error communicating with remote server.", ""}
		ret, _ := json.Marshal(m)
		p.ResponseWriter.Write(ret)
		return
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	p.ResponseWriter.Write(buf)
}

func (p *PagePlay) Share() {
	if p.Request.Method != "POST" {
		return
	}

    // FIXED
    data := url.Values{"body":{strings.TrimSpace(p.POST["body"])}}
	resp, err := http.PostForm(SHARE, data)

	defer resp.Body.Close()
	if err != nil {
		m := Compiled{"Error communicating with remote server.", ""}
		ret, _ := json.Marshal(m)
		p.ResponseWriter.Write(ret)
		return
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	p.ResponseWriter.Write(buf)
}
