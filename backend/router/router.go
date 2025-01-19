package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rejxcy/colorgame/controllers"
	"github.com/rejxcy/colorgame/controllers/game"
)

func middlewareCors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//请求方法
		method := ctx.Request.Method

		// 允许任何源
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Headers", "*")
		//服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
		// 跨域关键设置 让浏览器可以解析
		ctx.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusOK)
			return
		}

		// 处理请求
		ctx.Next()
	}
}

func Routers(engine *gin.Engine) {
	engine.Use(middlewareCors())
	ctx := &controllers.Context{}

	v1 := engine.Group("/api")

	{
		c := game.New(ctx)
		r := v1.Group("/game")
		r.GET("/ws", c.HandleWebSocket)
	}
}
