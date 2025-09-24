package storage

import (
	"encoding/json"
	"fmt"
	"github.com/jack0829/letsgo/common/fs"
	"github.com/jack0829/letsgo/wechat"
	"os"
	"path/filepath"
	"strings"
)

const (
	filenameAccessToken = "access_token.json"
	filenameJsApiTicket = "jsapi_ticket.json"
)

type file struct {
	baseDir string
}

func File(baseDir string) *file {
	return &file{
		baseDir: strings.TrimRight(baseDir, "/"),
	}
}

func (s *file) write(path string, v any) error {

	dir := filepath.Dir(path)
	if err := fs.MustDir(dir); err != nil {
		return err
	}

	fp, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fp.Close()

	return json.NewEncoder(fp).Encode(v)
}

func (s *file) read(path string, v any) error {

	fp, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fp.Close()

	return json.NewDecoder(fp).Decode(&v)
}

func (s *file) SetAccessToken(token *wechat.AccessToken) error {
	return s.write(
		fmt.Sprintf("%s/%s/%s", s.baseDir, token.AppID, filenameAccessToken),
		token,
	)
}

func (s *file) GetAccessToken(appID string) *wechat.AccessToken {
	var tk wechat.AccessToken
	if err := s.read(
		fmt.Sprintf("%s/%s/%s", s.baseDir, appID, filenameAccessToken),
		&tk,
	); err != nil {
		return nil
	}
	return &tk
}

func (s *file) SetJsApiTicket(ticket *wechat.JsApiTicket) error {
	return s.write(
		fmt.Sprintf("%s/%s/%s", s.baseDir, ticket.AppID, filenameJsApiTicket),
		ticket,
	)
}

func (s *file) GetJsApiTicket(appID string) *wechat.JsApiTicket {
	var tk wechat.JsApiTicket
	if err := s.read(
		fmt.Sprintf("%s/%s/%s", s.baseDir, appID, filenameJsApiTicket),
		&tk,
	); err != nil {
		return nil
	}
	return &tk
}

func (s *file) GetSession(appID, openID string) *wechat.Session {
	var ws wechat.Session
	if err := s.read(
		fmt.Sprintf("%s/%s/session/%s.json", s.baseDir, appID, openID),
		&ws,
	); err != nil {
		return nil
	}
	return &ws
}

func (s *file) SetSession(ws *wechat.Session) error {
	return s.write(
		fmt.Sprintf("%s/%s/session/%s.json", s.baseDir, ws.AppID, ws.OpenID),
		ws,
	)
}

func (s *file) DeleteSession(appID, openID string) error {
	os.Remove(
		fmt.Sprintf("%s/%s/session/%s.json", s.baseDir, appID, openID),
	)
	return nil
}

func (s *file) GetOAuthAccessToken(appID, openID string) *wechat.OAuthAccessToken {
	var tk wechat.OAuthAccessToken
	if err := s.read(
		fmt.Sprintf("%s/%s/oauth-access-token/%s.json", s.baseDir, appID, openID),
		&tk,
	); err != nil {
		return nil
	}
	return &tk
}

func (s *file) SetOAuthAccessToken(tk *wechat.OAuthAccessToken) error {
	return s.write(
		fmt.Sprintf("%s/%s/oauth-access-token/%s.json", s.baseDir, tk.AppID, tk.OpenID),
		tk,
	)
}
