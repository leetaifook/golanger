package web

import (
	"mime"
	"net/http"
	"net/url"
	"sync"
)

type Base struct {
	rmutex sync.RWMutex
	mutex  sync.Mutex
	header map[string][]string
}

func (b *Base) Init(w http.ResponseWriter, r *http.Request) *Base {

	return b
}

func (b *Base) AddHeader(k, v string) {
	if _, ok := b.header[k]; ok {
		b.header[k] = append(b.header[k], v)
	} else {
		b.header[k] = []string{v}
	}
}

func (b *Base) DelHeader(k string) {
	delete(b.header, k)
}

func (b *Base) getHttpGet(r *http.Request) map[string]string {
	g := map[string]string{}
	q := r.URL.Query()
	for key, _ := range q {
		g[key] = q.Get(key)
	}

	return g
}

func (b *Base) getHttpPost(r *http.Request, MAX_FORM_SIZE int64) map[string]string {
	ct := r.Header.Get("Content-Type")
	ct, _, _ = mime.ParseMediaType(ct)
	if ct == "multipart/form-data" {
		r.ParseMultipartForm(MAX_FORM_SIZE)
	} else {
		r.ParseForm()
	}

	p := map[string]string{}
	for key, _ := range r.Form {
		p[key] = r.FormValue(key)
	}

	return p
}

func (b *Base) getHttpCookie(r *http.Request) map[string]string {
	c := map[string]string{}
	for _, ck := range r.Cookies() {
		c[ck.Name], _ = url.QueryUnescape(ck.Value)
	}

	return c
}
