package handler

import (
	"github.com/Wheeeel/pushen-server/api/request"
	"github.com/Wheeeel/pushen-server/model"
	"github.com/Wheeeel/pushen-server/util"
	"github.com/go-playground/validator"
	"github.com/kataras/iris"
	"github.com/satori/go.uuid"
)

func DeviceListHandler(ctx iris.Context) {
	email := ctx.Values().GetString("email")
	if email == "" {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.Values().Set("error", "auth error")
		return
	}

	user, err := model.UserByEmail(email)
	if err != nil {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.Values().Set("error", "auth error")
		return
	}

	devices, err := model.DevicesByUserID(user.ID)
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

func BindAuthTokenRenewHandler(ctx iris.Context) {
	var dbt request.DeviceBindToken
	err := ctx.ReadJSON(&dbt)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Values().Set("error", err.Error())
		return
	}

	validate := validator.New()
	err = validate.Struct(dbt)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Values().Set("error", err.Error())
		return
	}

	var bindToken model.BindToken
	err = util.CopyStruct(&bindToken, dbt)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.Values().Set("error", err.Error())
		return
	}

	token, err := model.BindTokenByToken(bindToken.Token)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Values().Set("error", err.Error())
		return
	}

	if token.Status != model.BindStatusBinded {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Values().Set("error", "token has been binded.")
		return
	}

	var resp Response
	resp.Code = 200
	resp.Msg = "bind success"
	resp.Data = map[string]interface{}{
	}
	ctx.JSON(resp)
}

func DeviceBindTokenHandler(ctx iris.Context) {
	user, err := model.UserByEmail(ctx.Values().GetString("email"))
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Values().Set("error", err.Error())
		return
	}

	var bindToken model.BindToken
	bindToken.Status = model.BindStatusNotBinded
	bindToken.Token = uuid.NewV4().String()
	bindToken.UserID = user.ID
	err = model.BindTokenCreate(&bindToken)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.Values().Set("error", err.Error())
		return
	}

	var resp Response
	resp.Code = 200
	resp.Msg = "create message success"
	resp.Data = map[string]interface{}{
		"Token": bindToken.Token,
	}
	ctx.JSON(resp)
}
