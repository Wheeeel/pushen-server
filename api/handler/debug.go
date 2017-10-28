package handler

import (
	"github.com/kataras/iris"
)

func DebugHandler(ctx iris.Context) {
	ctx.ResponseWriter().Header().Set("Set-Cookie", "GUI=helloworld; Path=/; Domain=tony6.com; Expires=Fri, 27 Oct 2017 16:44:44 GMT; Max-Age=1799; HttpOnly")
	var resp Response
	resp.Code = 200
	resp.Msg = "debug"
	ctx.JSON(resp)
}
