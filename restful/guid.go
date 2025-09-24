package restful

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const GUIDName = "_guid"

func GUID(domain string) gin.HandlerFunc {
	return func(g *gin.Context) {

		var v string
		if v, _ = g.Cookie(GUIDName); v == "" {
			v = uuid.New().String()
			g.SetCookie(
				GUIDName,
				v,
				0,
				"",
				domain,
				false,
				true,
			)
		}
		g.Set(GUIDName, v)
		g.Next()
	}
}
