package storage

import (
	"context"
	"fmt"
	REDIS "github.com/go-redis/redis/v8"
	"github.com/jack0829/letsgo/wechat"
	jsoniter "github.com/json-iterator/go"
	"time"
)

const (
	redisKeyAccessToken      = "AccessToken"
	redisKeyJsApiTicket      = "JsApi:Ticket"
	redisKeySession          = "Session:"
	redisKeyOAuthAccessToken = "OAuthAccessToken:"
)

type redis struct {
	ctx context.Context
	c   REDIS.Cmdable
}

func Redis(
	ctx context.Context,
	cmd REDIS.Cmdable,
) *redis {
	return &redis{
		ctx: ctx,
		c:   cmd,
	}
}

func (s *redis) key(appID, obj string) string {
	return fmt.Sprintf("Wechat:%s:%s", appID, obj)
}

func (s *redis) SetAccessToken(token *wechat.AccessToken) error {

	if token == nil {
		return nil
	}

	if token.Expired() {
		return fmt.Errorf("cannot save expired token")
	}

	v, _ := jsoniter.MarshalToString(token)
	return s.c.Set(
		s.ctx,
		s.key(token.AppID, redisKeyAccessToken),
		v,
		time.Until(token.ExpireAt),
	).Err()
}

func (s *redis) GetAccessToken(appID string) *wechat.AccessToken {

	cmd := s.c.Get(s.ctx, s.key(appID, redisKeyAccessToken))
	if cmd.Err() != nil {
		return nil
	}

	var tk wechat.AccessToken
	if jsoniter.UnmarshalFromString(cmd.Val(), &tk) != nil {
		return nil
	}

	return &tk
}

func (s *redis) SetJsApiTicket(ticket *wechat.JsApiTicket) error {

	if ticket == nil {
		return nil
	}

	if ticket.Expired() {
		return fmt.Errorf("cannot save expired ticket")
	}

	v, _ := jsoniter.MarshalToString(ticket)
	return s.c.Set(
		s.ctx,
		s.key(ticket.AppID, redisKeyJsApiTicket),
		v,
		time.Until(ticket.ExpireAt),
	).Err()
}

func (s *redis) GetJsApiTicket(appID string) *wechat.JsApiTicket {

	cmd := s.c.Get(s.ctx, s.key(appID, redisKeyJsApiTicket))
	if cmd.Err() != nil {
		return nil
	}

	var tk wechat.JsApiTicket
	if jsoniter.UnmarshalFromString(cmd.Val(), &tk) != nil {
		return nil
	}

	return &tk
}

func (s *redis) SetSession(ws *wechat.Session) error {

	if ws == nil {
		return nil
	}

	v, _ := jsoniter.MarshalToString(ws)
	return s.c.Set(
		s.ctx,
		s.key(ws.AppID, redisKeySession+ws.OpenID),
		v,
		time.Hour*24*7,
	).Err()
}

func (s *redis) GetSession(appID, openID string) *wechat.Session {

	cmd := s.c.Get(s.ctx, s.key(appID, redisKeySession+openID))
	if cmd.Err() != nil {
		return nil
	}

	var ws wechat.Session
	if jsoniter.UnmarshalFromString(cmd.Val(), &ws) != nil {
		return nil
	}

	return &ws
}

func (s *redis) DeleteSession(appID, openID string) error {
	_ = s.c.Del(s.ctx, s.key(appID, redisKeySession+openID)).Err()
	return nil
}

func (s *redis) SetOAuthAccessToken(tk *wechat.OAuthAccessToken) error {

	if tk == nil {
		return nil
	}

	v, _ := jsoniter.MarshalToString(tk)
	return s.c.Set(
		s.ctx,
		s.key(tk.AppID, redisKeyOAuthAccessToken+tk.OpenID),
		v,
		time.Hour*24*20, // 强制20天过期重新授权，但实际 refreshToken 是 30 天过期
	).Err()
}

func (s *redis) GetOAuthAccessToken(appID, openID string) *wechat.OAuthAccessToken {

	cmd := s.c.Get(s.ctx, s.key(appID, redisKeyOAuthAccessToken+openID))
	if cmd.Err() != nil {
		return nil
	}

	var tk wechat.OAuthAccessToken
	if jsoniter.UnmarshalFromString(cmd.Val(), &tk) != nil {
		return nil
	}

	return &tk
}
