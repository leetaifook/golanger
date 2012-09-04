package controllers

import (
	"fmt"
	"time"
)

type PageIndex struct {
	*App
}

func init() {
	Page.RegisterController("index/", &PageIndex{Page})
}

func (p *PageIndex) Index() {
	p.SESSION["String"] = "String"
	p.SESSION["string"] = "string"
	p.SESSION["Int"] = 1
	p.SESSION["Map"] = map[string]string{
		"a": "b",
		"b": "c",
	}
	p.SESSION["Int64"] = time.Now().UnixNano()
}

func (p *PageIndex) TestPage() {
	p.Document.Title = "测试页面"
	p.ResponseWriter.Write([]byte(fmt.Sprintf("%v", p.SESSION["String"])))
	p.ResponseWriter.Write([]byte(fmt.Sprintf("%v", p.SESSION["string"])))
	p.ResponseWriter.Write([]byte(fmt.Sprintf("%v", p.SESSION["Int"])))
	p.ResponseWriter.Write([]byte(fmt.Sprintf("%v", p.SESSION["Map"])))
	p.ResponseWriter.Write([]byte(fmt.Sprintf("%v", p.SESSION["Int64"])))
}
