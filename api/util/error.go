package util

import "github.com/kataras/iris"

func WriteError(ctx iris.Context, status int, error, msg string) {
	ctx.StatusCode(status)
	ctx.Values().Set("error", error)
	ctx.Application().Logger().Info(msg)
}

func WriteErrorf(ctx iris.Context, status int, error, format string, msg ...interface{}) {
	ctx.StatusCode(status)
	ctx.Values().Set("error", error)
	ctx.Application().Logger().Infof(format, msg...)
}
