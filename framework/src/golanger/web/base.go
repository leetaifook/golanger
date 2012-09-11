package web

import (
	"golanger/session"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type Base struct {
	rmutex         sync.RWMutex
	mutex          sync.Mutex
	GET            map[string]string
	POST           map[string]string
	COOKIE         map[string]string
	SESSION        map[string]interface{}
	MAX_FORM_SIZE  int64
	SupportSession bool
	Session        *session.SessionManager
	Cookie         []*http.Cookie
}

func (b *Base) Init(w http.ResponseWriter, r *http.Request) *Base {
	b.GET = func() map[string]string {
		g := map[string]string{}
		q := r.URL.Query()
		for key, _ := range q {
			g[key] = q.Get(key)
		}

		return g
	}()

	b.POST = func() map[string]string {
		ct := r.Header.Get("Content-Type")
		ct, _, _ = mime.ParseMediaType(ct)
		if ct == "multipart/form-data" {
			r.ParseMultipartForm(b.MAX_FORM_SIZE)
		} else {
			r.ParseForm()
		}

		p := map[string]string{}
		for key, _ := range r.Form {
			p[key] = r.FormValue(key)
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
		b.mutex.Lock()
		b.SESSION = b.Session.Get(w, r)
		b.mutex.Unlock()
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
func (b *Base) SetCookie(w http.ResponseWriter, args ...interface{}) {
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

	http.SetCookie(w, bCookie)

	if expires > 0 {
		b.COOKIE[bCookie.Name] = bCookie.Value
	} else {
		delete(b.COOKIE, bCookie.Name)
	}
}
