package storage

import (
	"github.com/jack0829/letsgo/common/sets"
	"github.com/jack0829/letsgo/wechat"
	"gopkg.in/yaml.v3"
	"io"
	"sync"
)

var (
	defaultMemory memory
	Memory        = &defaultMemory
)

type memory struct {
	once             sync.Once
	initialized      bool
	accessToken      *sets.Set[string, *wechat.AccessToken]
	jsApiTicket      *sets.Set[string, *wechat.JsApiTicket]
	session          *sets.Set[string, *wechat.Session]
	oauthAccessToken *sets.Set[string, *wechat.OAuthAccessToken]
}

func (s *memory) initialize() {
	s.once.Do(func() {
		s.accessToken = &sets.Set[string, *wechat.AccessToken]{}
		s.jsApiTicket = &sets.Set[string, *wechat.JsApiTicket]{}
		s.session = &sets.Set[string, *wechat.Session]{}
		s.oauthAccessToken = &sets.Set[string, *wechat.OAuthAccessToken]{}
		s.initialized = true
	})
}

func (s *memory) SetAccessToken(token *wechat.AccessToken) error {
	s.initialize()
	s.accessToken.Set(token.AppID, token)
	return nil
}

func (s *memory) GetAccessToken(appID string) *wechat.AccessToken {
	s.initialize()
	if tk, ok := s.accessToken.Get(appID); ok {
		return tk
	}
	return nil
}

func (s *memory) SetJsApiTicket(ticket *wechat.JsApiTicket) error {
	s.initialize()
	s.jsApiTicket.Set(ticket.AppID, ticket)
	return nil
}

func (s *memory) GetJsApiTicket(appID string) *wechat.JsApiTicket {
	s.initialize()
	if tk, ok := s.jsApiTicket.Get(appID); ok {
		return tk
	}
	return nil
}

func (s *memory) SetSession(ws *wechat.Session) error {
	s.initialize()
	key := ws.AppID + ":" + ws.OpenID
	s.session.Set(key, ws)
	return nil
}

func (s *memory) GetSession(appID, openID string) *wechat.Session {
	s.initialize()
	key := appID + ":" + openID
	if sess, ok := s.session.Get(key); ok {
		return sess
	}
	return nil
}

func (s *memory) DeleteSession(appID, openID string) error {
	s.initialize()
	key := appID + ":" + openID
	s.session.Delete(key)
	return nil
}

func (s *memory) SetOAuthAccessToken(tk *wechat.OAuthAccessToken) error {
	s.initialize()
	key := tk.AppID + ":" + tk.OpenID
	s.oauthAccessToken.Set(key, tk)
	return nil
}

func (s *memory) GetOAuthAccessToken(appID, openID string) *wechat.OAuthAccessToken {
	s.initialize()
	key := appID + ":" + openID
	if tk, ok := s.oauthAccessToken.Get(key); ok {
		return tk
	}
	return nil
}

func (s *memory) dumpTo(w io.Writer) {

	var dump struct {
		AccessToken      []*wechat.AccessToken
		JsApiTicket      []*wechat.JsApiTicket
		Session          []*wechat.Session
		OAuthAccessToken []*wechat.OAuthAccessToken
	}

	s.accessToken.Each(func(_ string, tk *wechat.AccessToken) {
		dump.AccessToken = append(dump.AccessToken, tk)
	})

	s.jsApiTicket.Each(func(_ string, tk *wechat.JsApiTicket) {
		dump.JsApiTicket = append(dump.JsApiTicket, tk)
	})

	s.session.Each(func(_ string, s *wechat.Session) {
		dump.Session = append(dump.Session, s)
	})

	s.oauthAccessToken.Each(func(_ string, tk *wechat.OAuthAccessToken) {
		dump.OAuthAccessToken = append(dump.OAuthAccessToken, tk)
	})

	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	enc.Encode(dump)
}

func DumpMemory(w io.Writer) {
	if defaultMemory.initialized {
		defaultMemory.dumpTo(w)
	}
}
