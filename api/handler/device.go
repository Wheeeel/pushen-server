package handler

import (
	"github.com/Wheeeel/pushen-server/model"
	"github.com/kataras/iris"
)

func DeviceListHandler(ctx iris.Context) {
	email := ctx.Values().GetString("email")
	if email == "" {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.Values().Set("error", "auth error")
		return
	}

	user, err := model.UserByEmail(model.DefaultDB, email)
	if err != nil {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.Values().Set("error", "auth error")
		return
	}

	devices, err := model.DevicesByUserID(model.DefaultDB, user.ID)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.Values().Set("error", err.Error())
		return
	}

	var resp Response
	resp.Code = 200
	resp.Msg = "ok"
	resp.Data = devices
	ctx.JSON(resp)
}
