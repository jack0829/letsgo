package wechat

type AccessTokenStorage interface {
	SetAccessToken(token *AccessToken) error
	GetAccessToken(appID string) *AccessToken
}

type JsApiTicketStorage interface {
	SetJsApiTicket(ticket *JsApiTicket) error
	GetJsApiTicket(appID string) *JsApiTicket
}

type SessionStorage interface {
	SetSession(session *Session) error
	GetSession(appID, openID string) *Session
	DeleteSession(appID, openID string) error
}

type OAuthAccessTokenStorage interface {
	SetOAuthAccessToken(t *OAuthAccessToken) error
	GetOAuthAccessToken(appID, openID string) *OAuthAccessToken
}

type Storage interface {
	AccessTokenStorage
	JsApiTicketStorage
	SessionStorage
	OAuthAccessTokenStorage
}

type storage struct {
	accessToken      AccessTokenStorage
	jsApiTicket      JsApiTicketStorage
	session          SessionStorage
	oauthAccessToken OAuthAccessTokenStorage
}
