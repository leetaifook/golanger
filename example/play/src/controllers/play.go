// Copyright 2012 The Golanger Authors. All rights reserved.
// Use of this source code is governed by a GPLv3
// license that can be found in the LICENSE file.

// Compile The Go Source By play.golang.org.
package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	COMPILE = "http://play.golang.org/compile"
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
