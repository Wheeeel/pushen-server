package handler

import (
	"github.com/kataras/iris"
)

func CorsHandler(ctx iris.Context) {
	ctx.Header("Access-Control-Request-Headers", "*")
	ctx.Next()
}
