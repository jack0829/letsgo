package rbac

import (
	"github.com/gin-gonic/gin"
	"github.com/jack0829/letsgo/restful"
	"net/http"
)

var (
	errorNotAcceptable = "权限不足"
)

const (
	RoleSuperAdmin      = 101  // 超级管理员
	RoleEnterpriseAdmin = 1001 // 企业管理员
)

func basic(ctx *gin.Context, roles []int) (role map[int]struct{}, finish bool) {

	if len(roles) == 0 {
		finish = true
		ctx.Next()
		return
	}

	req := &restful.Request{Gin: ctx}
	userRoles := req.GetRoles()
	l := len(userRoles)
	if l == 0 {
		finish = true
		ctx.AbortWithStatusJSON(
			http.StatusNotAcceptable,
			restful.Error(http.StatusNotAcceptable, errorNotAcceptable),
		)
		return
	}

	role = make(map[int]struct{})
	for _, r := range userRoles {
		role[r] = struct{}{}
	}
	return
}

// Any 满足任意一个角色即可
func Any(roles ...int) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		role, finish := basic(ctx, roles)
		if finish {
			return
		}

		for _, r := range roles {
			if _, ok := role[r]; ok {
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatusJSON(
			http.StatusNotAcceptable,
			restful.Error(http.StatusNotAcceptable, errorNotAcceptable),
		)
	}
}

// All 满足所有角色才通过
func All(roles ...int) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		role, finish := basic(ctx, roles)
		if finish {
			return
		}

		for _, r := range roles {
			if _, ok := role[r]; !ok {
				ctx.AbortWithStatusJSON(
					http.StatusNotAcceptable,
					restful.Error(http.StatusNotAcceptable, errorNotAcceptable),
				)
				return
			}
		}

		ctx.Next()
	}
}

// Not 满足任意一个角色就不通过
func Not(roles ...int) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		role, finish := basic(ctx, roles)
		if finish {
			return
		}

		for _, r := range roles {
			if _, ok := role[r]; ok {
				ctx.AbortWithStatusJSON(
					http.StatusNotAcceptable,
					restful.Error(http.StatusNotAcceptable, errorNotAcceptable),
				)
				return
			}
		}

		ctx.Next()
	}
}
