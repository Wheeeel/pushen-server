package handler

import (
	"github.com/Wheeeel/pushen-server/api/request"
	"github.com/Wheeeel/pushen-server/model"
	"github.com/Wheeeel/pushen-server/util"
	"github.com/go-playground/validator"
	"github.com/kataras/iris"
)

var (
	useToken  = "PUSHEN-USE-TOKEN"
	tokenName = "PUSHEN-TOKEN"
)

func UserLoginHandler(ctx iris.Context) {
	sess := session.Start(ctx)

	var ul request.UserLogin
	err := ctx.ReadJSON(&ul)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Values().Set("error", err.Error())
		return
	}

	validate := validator.New()
	err = validate.Struct(ul)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Values().Set("error", err.Error())
		return
	}

	ok, err := model.UserValidate(ul.Email, ul.Password)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Values().Set("error", err.Error())
		return
	}

	if !ok {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.Values().Set("error", "permission denied")
		return
	}

	// set cookie
	sess.Set(authenticated, true)
	sess.Set("email", ul.Email)

	var resp Response
	resp.Code = iris.StatusOK
	resp.Msg = "success"
	resp.Data = map[string]interface{}{
		"cookie": sess.ID(),
	}
	ctx.JSON(resp)
}

func UserLogoutHandler(ctx iris.Context) {
}

func UserCreateHandler(ctx iris.Context) {
	var uc request.UserCreate
	err := ctx.ReadJSON(&uc)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Values().Set("error", err.Error())
		return
	}

	validate := validator.New()
	err = validate.Struct(uc)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Values().Set("error", err.Error())
		return
	}

	var user model.User
	err = util.CopyStruct(&user, uc)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.Values().Set("error", err.Error())
		return
	}

	err = model.UserCreate(&user)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.Values().Set("error", err.Error())
		return
	}

	var resp Response
	resp.Code = 200
	resp.Msg = "create user success"
	ctx.JSON(resp)
}

func UserInfoHandler(ctx iris.Context) {
	var email string
	sess := session.Start(ctx)
	ei := sess.Get("email")
	if ei == nil {
		var resp Response
		resp.Code = 403
		resp.Msg = "please sign in"
		ctx.JSON(resp)
		return
	}
	var ok bool
	email, ok = ei.(string)
	if !ok {
		var resp Response
		resp.Code = 403
		resp.Msg = "please sign in"
		ctx.JSON(resp)
		return
	}

	user, err := model.UserByEmail(email)
	if err != nil {
		var resp Response
		resp.Code = 500
		resp.Msg = err.Error()
		ctx.JSON(resp)
		return
	}

	var resp Response
	resp.Code = 500
	resp.Msg = "ok"
	resp.Data = user
	ctx.JSON(resp)
	return
}

func UserTokenHandler(ctx iris.Context) {
}
