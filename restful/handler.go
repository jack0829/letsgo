package restful

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type JsonHandler func(req *Request) *Response

func JSON(handler JsonHandler) gin.HandlerFunc {

	return func(g *gin.Context) {

		if resp := handler(&Request{
			Gin: g,
		}); resp != nil {
			g.PureJSON(http.StatusOK, resp)
		}

	}
}

type NopHandler func(req *Request)

func Nop(handler NopHandler) gin.HandlerFunc {
	return func(g *gin.Context) {
		handler(&Request{
			Gin: g,
		})
	}
}
