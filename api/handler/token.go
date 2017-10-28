package handler

import (
	"github.com/Wheeeel/pushen-server/api/request"
	"github.com/Wheeeel/pushen-server/model"
	"github.com/Wheeeel/pushen-server/util"
	"github.com/go-playground/validator"
	"github.com/kataras/iris"
	"github.com/satori/go.uuid"
)

func BindAuthTokenHandler(ctx iris.Context) {
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

	if token.Status == model.BindStatusBinded {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Values().Set("error", "token has been binded.")
		return
	}

	token.Status = model.BindStatusBinded
	err = model.BindTokenUpdateStatus(&bindToken)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.Values().Set("error", err.Error())
		return
	}

	// generate auth token
	var at model.AuthToken
	at.Token = uuid.NewV4().String()
	at.UserID = token.UserID
	err = model.AuthTokenCreate(&at)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.Values().Set("error", err.Error())
		return
	}

	user, err := model.UserByID(token.UserID)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Values().Set("error", err.Error())
		return
	}

	// generate device
	var device model.Device
	device.Type = "phone"
	device.UUID = uuid.NewV4().String()
	device.UserID = user.ID
	//device.Status = 1
	err = model.DeviceCreate(&device)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.Values().Set("error", err.Error())
		return
	}

	var resp Response
	resp.Code = 200
	resp.Msg = "bind success"
	resp.Data = map[string]interface{}{
		"Token":  at.Token,
		"Device": device.UUID,
	}
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
