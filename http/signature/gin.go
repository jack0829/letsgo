package signature

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GinMiddleWare 中间件
func GinMiddleWare(s *Signature) gin.HandlerFunc {

	return func(g *gin.Context) {

		if err := s.Check(g.Request); err != nil {
			g.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		g.Next()
	}
}
