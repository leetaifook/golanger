package session

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type SessionManager struct {
	CookieName    string
	mutex         *sync.RWMutex
	sessions      map[string][2]map[string]interface{}
	expires       int
	timerDuration time.Duration
}

func New(cookieName string, expires int, timerDuration time.Duration) *SessionManager {
	if cookieName == "" {
		cookieName = "GoLangerSession"
	}

	if expires <= 0 {
		expires = 3600
	}

	if timerDuration <= 0 {
		timerDuration, _ = time.ParseDuration("24h")
	}

	s := &SessionManager{
		CookieName:    cookieName,
		mutex:         &sync.RWMutex{},
		sessions:      map[string][2]map[string]interface{}{},
		expires:       expires,
		timerDuration: timerDuration,
	}

	time.AfterFunc(s.timerDuration, func() { s.GC() })

	return s
}

func (s *SessionManager) Get(rw http.ResponseWriter, req *http.Request) map[string]interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	var sessionSign string
	if c, err := req.Cookie(s.CookieName); err == nil {
		sessionSign = c.Value
		if sessionValue, ok := s.sessions[sessionSign]; ok {
			return sessionValue[1]
		}
	}

	s.mutex.RUnlock()
	s.mutex.Lock()
	sessionSign = s.new(rw)
	s.mutex.Unlock()
	s.mutex.RLock()

	return s.sessions[sessionSign][1]
}

func (s *SessionManager) Len() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return int64(len(s.sessions))
}

func (s *SessionManager) new(rw http.ResponseWriter) string {
	timeNano := time.Now().UnixNano()
	sessionSign := s.sessionSign()
	s.sessions[sessionSign] = [2]map[string]interface{}{
		map[string]interface{}{
			"create": timeNano,
		},
		map[string]interface{}{},
	}

	bCookie := &http.Cookie{
		Name:  s.CookieName,
		Value: url.QueryEscape(sessionSign),
		Path:  "/",
	}

	http.SetCookie(rw, bCookie)

	return sessionSign
}

func (s *SessionManager) Clear(sessionSign string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.sessions, sessionSign)
}

func (s *SessionManager) GC() {
	s.mutex.Lock()
	for sessionSign, _ := range s.sessions {
		if (s.sessions[sessionSign][0]["create"].(int64) + int64(s.expires)) <= time.Now().Unix() {
			delete(s.sessions, sessionSign)
		}
	}

	s.mutex.Unlock()
	time.AfterFunc(s.timerDuration, func() { s.GC() })
}

func (s *SessionManager) sessionSign() string {
	var n int = 24
	b := make([]byte, n)
	io.ReadFull(rand.Reader, b)

	//return length:32
	return base64.URLEncoding.EncodeToString(b)
}
