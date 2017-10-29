package handler

import (
	"github.com/Wheeeel/pushen-server/api/request"
	apiutil "github.com/Wheeeel/pushen-server/api/util"
	"github.com/Wheeeel/pushen-server/model"
	"github.com/Wheeeel/pushen-server/util"
	"github.com/go-playground/validator"
	"github.com/kataras/iris"
	"github.com/satori/go.uuid"
)

func DeviceBindTokenRenewHandler(ctx iris.Context) {
	var dbt request.DeviceBindToken
	err := ctx.ReadJSON(&dbt)
	if err != nil {
		apiutil.WriteErrorf(ctx, iris.StatusBadRequest, err.Error(), "device bind token renew error: %v", err.Error())
		return
	}

	validate := validator.New()
	err = validate.Struct(dbt)
	if err != nil {
		apiutil.WriteErrorf(ctx, iris.StatusBadRequest, err.Error(), "device bind token renew error: %v", err.Error())
		return
	}

	var bindToken model.BindToken
	err = util.CopyStruct(&bindToken, dbt)
	if err != nil {
		apiutil.WriteErrorf(ctx, iris.StatusInternalServerError, err.Error(), "device bind token renew error: %v", err.Error())
		return
	}

	bindToken, err = model.BindTokenByToken(model.DefaultDB, bindToken.Token)
	if err != nil {
		apiutil.WriteErrorf(ctx, iris.StatusBadRequest, err.Error(), "device bind token renew error: %v", err.Error())
		return
	}

	if bindToken.Status != model.BindStatusBinded {
		apiutil.WriteError(ctx, iris.StatusBadRequest, "token has been binded", "device bind token renew error")
		return
	}

	var resp Response
	resp.Code = 200
	resp.Msg = "bind success"
	ctx.JSON(resp)
}

// DeviceBindTokenHandler generate device bind token
func DeviceBindTokenHandler(ctx iris.Context) {
	email := ctx.Values().GetString("email")
	user, err := model.UserByEmail(model.DefaultDB, email)
	if err != nil {
		apiutil.WriteErrorf(ctx, iris.StatusBadRequest, err.Error(), "device bind token error: %v", err.Error())
		return
	}

	var bindToken model.BindToken
	bindToken.Status = model.BindStatusNotBinded
	bindToken.Token = uuid.NewV4().String()
	bindToken.UserID = user.ID
	tx := model.DefaultDB.Begin()
	err = model.BindTokenCreate(tx, &bindToken)
	if err != nil {
		tx.Rollback()
		apiutil.WriteErrorf(ctx, iris.StatusInternalServerError, err.Error(), "device bind token error: %v", err.Error())
		return
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		apiutil.WriteErrorf(ctx, iris.StatusInternalServerError, err.Error(), "device bind token error: %v", err.Error())
		return
	}

	var resp Response
	resp.Code = 200
	resp.Msg = "generate device bind token success"
	resp.Data = map[string]interface{}{
		"Token": bindToken.Token,
	}
	ctx.JSON(resp)
}
