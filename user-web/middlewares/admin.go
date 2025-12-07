// 鉴权是否为管理员中间件
package middlewares

import (
	"mx-shop-api/user-web/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func IsAdmin() gin.HandlerFunc {
	s := zap.S()
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		currentUser := claims.(*models.CustomClaims)

		s.Infof("currentUser AuthorityId %d", currentUser.AuthorityId)
		// 不是管理员
		if currentUser.AuthorityId == 1 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg": "无权限",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
