package restful

import (
	"github.com/gin-gonic/gin"
	"github.com/jack0829/letsgo/jwt"
	"net/http"
	"strings"
)

const (
	errorNotLogin   = "请先登录"
	errorInvalidUid = "用户ID异常"
	errorInvalidEid = "企业ID异常"
)

const (
	defaultAuthHeaderName = "Authorization"
)

type authOpt struct {
	HeaderName string // 默认 Authorization，优先
	OmitCookie string // 非空时，尝试从 Cookie 获取，默认不尝试
	Required   bool   // 是否必须，默认为否
}

type AuthOption func(opt *authOpt)

func AuthOmitCookie(cookie string) AuthOption {
	return func(opt *authOpt) {
		opt.OmitCookie = cookie
	}
}

func AuthRequired(opt *authOpt) {
	opt.Required = true
}

func (a *authOpt) unauthorized(ctx *gin.Context, message string) {
	if a.Required {
		ctx.AbortWithStatusJSON(
			http.StatusUnauthorized,
			Error(http.StatusUnauthorized, message),
		)
	} else {
		ctx.Next()
	}
}

func (a *authOpt) ginHandler(ctx *gin.Context) {

	if a.HeaderName == "" {
		a.HeaderName = defaultAuthHeaderName
	}

	var token string

	if header := ctx.GetHeader(a.HeaderName); header != "" {
		// 没传 或 不符合规则：未登录
		auth := strings.SplitN(header, " ", 2)
		if len(auth) != 2 || auth[0] != "Bearer" {
			a.unauthorized(ctx, errorNotLogin)
			return
		}
		token = auth[1]
	}

	if token == "" {
		if c := a.OmitCookie; c != "" {
			token, _ = ctx.Cookie(c)
		}
	}

	if token == "" {
		a.unauthorized(ctx, errorNotLogin)
		return
	}

	if err := jwt.Auth.Decode(token); err != nil {
		a.unauthorized(ctx, errorNotLogin)
		return
	}

	data := jwt.Auth.Data()
	encoder := jwt.Auth.GetEncoder()
	uid, err := data.GetUid(encoder)
	if err != nil {
		a.unauthorized(ctx, errorInvalidUid)
		return
	}
	eid, err := data.GetEid(encoder)
	if err != nil {
		a.unauthorized(ctx, errorInvalidEid)
		return
	}

	ctx.Set("uid", uid)
	ctx.Set("eid", eid)
	ctx.Set("nick", data.Nick)
	ctx.Set("avatar", data.Avatar)
	ctx.Set("svc", data.Svc)
	ctx.Set("roles", data.Roles)
	ctx.Set("jwt", token)
	ctx.Next()
}

func Auth(ops ...AuthOption) gin.HandlerFunc {

	var opt authOpt
	for _, op := range ops {
		op(&opt)
	}

	return opt.ginHandler
}
