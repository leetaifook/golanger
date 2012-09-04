package web

import (
	"golanger/session"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Base struct {
	GET            map[string]string
	POST           map[string]string
	COOKIE         map[string]string
	SESSION        map[string]interface{}
	MAX_FORM_SIZE  int64
	SupportSession bool
	Session        *session.SessionManager
	SessionName    string
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Cookie         []*http.Cookie
}

func (b *Base) Init() *Base {
	b.GET = func() map[string]string {
		g := map[string]string{}
		q := b.Request.URL.Query()
		for key, _ := range q {
			g[key] = q.Get(key)
		}

		return g
	}()

	b.POST = func() map[string]string {
		ct := b.Request.Header.Get("Content-Type")
		ct, _, _ = mime.ParseMediaType(ct)
		if ct == "multipart/form-data" {
			b.Request.ParseMultipartForm(b.MAX_FORM_SIZE)
		} else {
			b.Request.ParseForm()
		}

		p := map[string]string{}
		for key, _ := range b.Request.Form {
			p[key] = b.Request.FormValue(key)
		}

		return p
	}()

	b.COOKIE = func() map[string]string {
		c := map[string]string{}
		for _, ck := range b.Cookie {
			c[ck.Name], _ = url.QueryUnescape(ck.Value)
		}

		return c
	}()

	if b.SupportSession {
		b.SESSION = b.Session.Get(b.ResponseWriter, b.Request)
		b.SessionName = b.Session.CookieName
	} else {
		b.SESSION = map[string]interface{}{}
	}

	return b
}

/*
cookie[0] => name string
cookie[1] => value string
cookie[2] => expires string
cookie[3] => path string
cookie[4] => domain string
*/
func (b *Base) SetCookie(args ...interface{}) {
	if len(args) < 2 {
		return
	}

	const LEN = 5
	var cookie = [LEN]interface{}{}

	for k, v := range args {
		if k >= LEN {
			break
		}

		cookie[k] = v
	}

	var (
		name    string
		value   string
		expires int
		path    string
		domain  string
	)

	if v, ok := cookie[0].(string); ok {
		name = v
	} else {
		return
	}

	if v, ok := cookie[1].(string); ok {
		value = v
	} else {
		return
	}

	if v, ok := cookie[2].(int); ok {
		expires = v
	}

	if v, ok := cookie[3].(string); ok {
		path = v
	}

	if v, ok := cookie[4].(string); ok {
		domain = v
	}

	bCookie := &http.Cookie{
		Name:   name,
		Value:  url.QueryEscape(value),
		Path:   path,
		Domain: domain,
	}

	if expires > 0 {
		d, _ := time.ParseDuration(strconv.Itoa(expires) + "s")
		bCookie.Expires = time.Now().Add(d)
	}

	http.SetCookie(b.ResponseWriter, bCookie)

	if expires > 0 {
		b.COOKIE[bCookie.Name] = bCookie.Value
	} else {
		delete(b.COOKIE, bCookie.Name)
	}
}
