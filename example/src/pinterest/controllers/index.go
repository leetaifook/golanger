package controllers

import (
	. "../models"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"golanger/utils"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type PageIndex struct {
	*App
}

func init() {
	Page.RegisterController("index/", &PageIndex{Page})
}

func (p *PageIndex) Index() {
	body := utils.M{}
	body["images"], _ = GetImagesLists()
	p.Body = body
}

func (p *PageIndex) Upload() {
	if p.Request.Method == "POST" {
		buf := new(bytes.Buffer)
		tnow := time.Now()
		binary.Write(buf, binary.LittleEndian, tnow.UnixNano())
		fileName := strings.TrimRight(base64.URLEncoding.EncodeToString(buf.Bytes()), "=")
		filePath := p.UploadDirectory + "images/"
		os.MkdirAll(filePath, 0777)
		file, fileHeader, err := p.Request.FormFile("file")
		if err != nil {
			fmt.Println(err)
			return
		}

		fileExt := strings.ToLower(path.Ext(fileHeader.Filename))
		fileContent, _ := ioutil.ReadAll(file)
		ioutil.WriteFile(filePath+fileName+fileExt, fileContent, 0777)

		go SaveImages(Images{
			Name:       fileName,
			Ext:        fileExt,
			Path:       filePath[1:],
			CreateTime: time.Now().Unix(),
		})

		http.Redirect(p.ResponseWriter, p.Request, "/", http.StatusFound)
	}
}

func (p *PageIndex) Page() {
	if p.Request.Method == "GET" {
		if pg, ok := p.GET["page"]; ok {
			ipg, _ := strconv.Atoi(pg)
			body := utils.M{}
			body["images"], _ = GetImagesLists(ipg)
			p.Body = body
		}
	}

}
