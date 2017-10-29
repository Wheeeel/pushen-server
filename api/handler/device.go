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

func DeviceListHandler(ctx iris.Context) {
	email := ctx.Values().GetString("email")
	if email == "" {
		apiutil.WriteError(ctx, iris.StatusForbidden, "auth error", "device list handler error: email empty")
		return
	}

	user, err := model.UserByEmail(model.DefaultDB, email)
	if err != nil {
		apiutil.WriteErrorf(ctx, iris.StatusForbidden, "auth error", "device list handler error: %v", err.Error())
		return
	}

	devices, err := model.DevicesByUserID(model.DefaultDB, user.ID)
	if err != nil {
		apiutil.WriteErrorf(ctx, iris.StatusInternalServerError, err.Error(), "device list handler error: %v", err.Error())
		return
	}

	var resp Response
	resp.Code = 200
	resp.Msg = "ok"
	resp.Data = devices
	ctx.JSON(resp)
}

// DeviceBindHandler bind device using device bind token
func DeviceBindHandler(ctx iris.Context) {
	var dbt request.DeviceBindToken
	err := ctx.ReadJSON(&dbt)
	if err != nil {
		apiutil.WriteErrorf(ctx, iris.StatusBadRequest, err.Error(), "device bind handler error: %v", err.Error())
		return
	}

	validate := validator.New()
	err = validate.Struct(dbt)
	if err != nil {
		apiutil.WriteErrorf(ctx, iris.StatusBadRequest, err.Error(), "device bind handler error: %v", err.Error())
		return
	}

	var bindToken model.BindToken
	err = util.CopyStruct(&bindToken, dbt)
	if err != nil {
		apiutil.WriteErrorf(ctx, iris.StatusInternalServerError, err.Error(), "device bind handler error: %v", err.Error())
		return
	}

	bindToken, err = model.BindTokenByToken(model.DefaultDB, bindToken.Token)
	if err != nil {
		apiutil.WriteErrorf(ctx, iris.StatusBadRequest, err.Error(), "device bind handler error: %v", err.Error())
		return
	}

	if bindToken.Status == model.BindStatusBinded {
		apiutil.WriteError(ctx, iris.StatusBadRequest, "token has been binded", "device bind handler error: status error")
		return
	}

	tx := model.DefaultDB.Begin()
	bindToken.Status = model.BindStatusBinded
	err = model.BindTokenUpdateStatus(tx, &bindToken)
	if err != nil {
		tx.Rollback()
		apiutil.WriteErrorf(ctx, iris.StatusInternalServerError, err.Error(), "device bind handler error: %v", err.Error())
		return
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		apiutil.WriteErrorf(ctx, iris.StatusInternalServerError, err.Error(), "device bind handler error: %v", err.Error())
		return
	}

	// generate auth token
	var at model.AuthToken
	at.Token = uuid.NewV4().String()
	at.UserID = bindToken.UserID
	tx = model.DefaultDB.Begin()
	err = model.AuthTokenCreate(tx, &at)
	if err != nil {
		tx.Rollback()
		apiutil.WriteErrorf(ctx, iris.StatusInternalServerError, err.Error(), "device bind handler error: %v", err.Error())
		return
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		apiutil.WriteErrorf(ctx, iris.StatusInternalServerError, err.Error(), "device bind handler error: %v", err.Error())
		return
	}

	user, err := model.UserByID(model.DefaultDB, bindToken.UserID)
	if err != nil {
		apiutil.WriteErrorf(ctx, iris.StatusBadRequest, err.Error(), "device bind handler error: %v", err.Error())
		return
	}

	// generate device
	var device model.Device
	device.Type = "phone"
	device.UUID = uuid.NewV4().String()
	device.UserID = user.ID
	device.Status = model.DeviceStatusBinded
	tx = model.DefaultDB.Begin()
	err = model.DeviceCreate(tx, &device)
	if err != nil {
		tx.Rollback()
		apiutil.WriteErrorf(ctx, iris.StatusInternalServerError, err.Error(), "device bind handler error: %v", err.Error())
		return
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		apiutil.WriteErrorf(ctx, iris.StatusInternalServerError, err.Error(), "device bind handler error: %v", err.Error())
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

func DeviceUnbindHandler(ctx iris.Context) {
	var du request.DeviceUnbind
	err := ctx.ReadJSON(&du)
	if err != nil {
		apiutil.WriteErrorf(ctx, iris.StatusBadRequest, err.Error(), "device unbind handler error: %v", err.Error())
		return
	}

	validate := validator.New()
	err = validate.Struct(du)
	if err != nil {
		apiutil.WriteErrorf(ctx, iris.StatusBadRequest, err.Error(), "device unbind handler error: %v", err.Error())
		return
	}

	device, err := model.DeviceByUUID(model.DefaultDB, du.Device)
	if err != nil {
		apiutil.WriteErrorf(ctx, iris.StatusBadRequest, err.Error(), "device unbind handler error: %v", err.Error())
		return
	}

	device.Status = model.DeviceStatusUnbinded
	tx := model.DefaultDB.Begin()
	err = model.DeviceUpdateStatus(tx, &device)
	if err != nil {
		tx.Rollback()
		apiutil.WriteErrorf(ctx, iris.StatusInternalServerError, err.Error(), "device unbind handler error: %v", err.Error())
		return
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		apiutil.WriteErrorf(ctx, iris.StatusInternalServerError, err.Error(), "device unbind handler error: %v", err.Error())
		return
	}

	var resp Response
	resp.Code = 200
	resp.Msg = "device unbind success"
	ctx.JSON(resp)
}
