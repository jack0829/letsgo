package restful

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Request struct {
	Gin *gin.Context
}

func (r *Request) GetUID() int {
	return r.Gin.GetInt("uid")
}

func (r *Request) GetEID() int {
	return r.Gin.GetInt("eid")
}

func (r *Request) GetSvc() []string {
	return r.Gin.GetStringSlice("svc")
}

func (r *Request) GetNick() string {
	return r.Gin.GetString("nick")
}

func (r *Request) GetRoles() []int {
	if v, ok := r.Gin.Get("roles"); ok {
		if r, ok := v.([]int); ok {
			return r
		}
	}
	return nil
}

func (r *Request) GetJWT() string {
	return r.Gin.GetString("jwt")
}

func (r *Request) GetGUID() string {
	return r.Gin.GetString(GUIDName)
}

func (r *Request) BindJSON(v any) error {
	return r.Gin.ShouldBindWith(v, binding.JSON)
}

func (r *Request) BindQuery(v any) error {
	return r.Gin.ShouldBindWith(v, binding.Query)
}
