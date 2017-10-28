package handler

import (
	"github.com/kataras/iris"
)

func CorsHandler(ctx iris.Context) {
	ctx.ResponseWriter().Header().Set("Access-Control-Request-Headers", "*")
	ctx.ResponseWriter().Header().Set("Access-Control-Allow-Origin", ctx.Request().Header.Get("Origin"))
	ctx.Next()
}
