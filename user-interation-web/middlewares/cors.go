// cors 解决跨域
// 原理：浏览器在发起负责请求时（例如有自定义的header），首先会发起预检请求询问服务端是否允许访问，我们这里通过设置header的方式通过响应预检请求告诉浏览器完全可以访问。
package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*") // 允许所有域（生产请改成具体域名）
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, AccessToken, X-CSRF-Token, Token, X-Requested-With, x-token")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")

		// 处理预检请求 OPTIONS
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}
