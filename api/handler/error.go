package handler

import "github.com/kataras/iris"

func ErrorForbidden(ctx iris.Context) {
	var resp Response
	resp.Code = iris.StatusForbidden

	errMessage := ctx.Values().GetString("error")
	if errMessage != "" {
		resp.Msg = errMessage
	}

	ctx.JSON(resp)
}

func ErrorInternal(ctx iris.Context) {
	var resp Response
	resp.Code = iris.StatusInternalServerError

	errMessage := ctx.Values().GetString("error")
	if errMessage != "" {
		resp.Msg = errMessage
	}

	ctx.JSON(resp)
}

func ErrorServiceUnavailable(ctx iris.Context) {
	var resp Response
	resp.Code = iris.StatusServiceUnavailable

	errMessage := ctx.Values().GetString("error")
	if errMessage != "" {
		resp.Msg = errMessage
	}

	ctx.JSON(resp)
}

func ErrorNotFound(ctx iris.Context) {
	var resp Response
	resp.Code = iris.StatusNotFound
	ctx.JSON(resp)
}

func ErrorBadRequest(ctx iris.Context) {
	var resp Response
	resp.Code = iris.StatusBadRequest

	errMessage := ctx.Values().GetString("error")
	if errMessage != "" {
		resp.Msg = errMessage
	}

	ctx.JSON(resp)
}
