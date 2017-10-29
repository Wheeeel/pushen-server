package handler

import (
	"math/rand"
	"time"

	"github.com/Wheeeel/pushen-server/api/request"
	apiutil "github.com/Wheeeel/pushen-server/api/util"
	"github.com/Wheeeel/pushen-server/model"
	"github.com/Wheeeel/pushen-server/util"
	"github.com/go-playground/validator"
	"github.com/kataras/iris"
)

var (
	useToken = "PUSHEN-USE-TOKEN"
	//tokenName = "PUSHEN-TOKEN"

	avatar = []string{
		"http://a.hiphotos.baidu.com/image/pic/item/00e93901213fb80e139458333cd12f2eb8389488.jpg",
		"http://a.hiphotos.baidu.com/image/pic/item/342ac65c103853437bb76ac59913b07ecb8088de.jpg",
		"http://s0.hao123img.com/res/img/moe/0707QX_333G.jpg",
		"http://s0.hao123img.com/res/img/moe/0707QX_333G.jpg",
		"http://fdfs.xmcdn.com/group5/M01/33/4B/wKgDtVORjieivT6jAAO-Rdq2VaU918_mobile_large.jpg",
	}
)

func UserLoginHandler(ctx iris.Context) {
	s := session.Start(ctx)

	var ul request.UserLogin
	err := ctx.ReadJSON(&ul)
	if err != nil {
		apiutil.WriteErrorf(ctx, iris.StatusBadRequest, err.Error(), "user login error: %v", err.Error())
		return
	}

	validate := validator.New()
	err = validate.Struct(ul)
	if err != nil {
		apiutil.WriteErrorf(ctx, iris.StatusBadRequest, err.Error(), "user login error: %v", err.Error())
		return
	}

	ok, err := model.UserValidate(model.DefaultDB, ul.Email, ul.Password)
	if err != nil {
		apiutil.WriteErrorf(ctx, iris.StatusBadRequest, err.Error(), "user login error: %v", err.Error())
		return
	}

	if !ok {
		apiutil.WriteError(ctx, iris.StatusForbidden, "permission denied", "user login error")
		return
	}

	// set cookie
	s.Set(authenticated, true)
	s.Set("email", ul.Email)

	var resp Response
	resp.Code = iris.StatusOK
	resp.Msg = "success"
	ctx.JSON(resp)
}

func UserLogoutHandler(ctx iris.Context) {
	s := session.Start(ctx)
	s.Clear()

	var resp Response
	resp.Code = iris.StatusOK
	resp.Msg = "success"
	ctx.JSON(resp)
}

func UserCreateHandler(ctx iris.Context) {
	var uc request.UserCreate
	err := ctx.ReadJSON(&uc)
	if err != nil {
		apiutil.WriteErrorf(ctx, iris.StatusBadRequest, err.Error(), "user create error: %v", err.Error())
		return
	}

	validate := validator.New()
	err = validate.Struct(uc)
	if err != nil {
		apiutil.WriteErrorf(ctx, iris.StatusBadRequest, err.Error(), "user create error: %v", err.Error())
		return
	}

	var user model.User
	err = util.CopyStruct(&user, uc)
	if err != nil {
		apiutil.WriteErrorf(ctx, iris.StatusInternalServerError, err.Error(), "user create error: %v", err.Error())
		return
	}

	rand.Seed(time.Now().Unix())
	user.AvatarURL = avatar[rand.Intn(5)]
	tx := model.DefaultDB.Begin()
	err = model.UserCreate(tx, &user)
	if err != nil {
		tx.Rollback()
		apiutil.WriteErrorf(ctx, iris.StatusInternalServerError, err.Error(), "user create error: %v", err.Error())
		return
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		apiutil.WriteErrorf(ctx, iris.StatusInternalServerError, err.Error(), "user create error: %v", err.Error())
		return
	}

	var resp Response
	resp.Code = 200
	resp.Msg = "create user success"
	ctx.JSON(resp)
}

func UserInfoHandler(ctx iris.Context) {
	email := ctx.Values().GetString("email")
	if email == "" {
		apiutil.WriteError(ctx, iris.StatusBadRequest, "auth error", "user info error")
		return
	}

	user, err := model.UserByEmail(model.DefaultDB, email)
	if err != nil {
		apiutil.WriteErrorf(ctx, iris.StatusInternalServerError, err.Error(), "user info error: %v", err.Error())
		return
	}

	var resp Response
	resp.Code = 200
	resp.Msg = "ok"
	resp.Data = user
	ctx.JSON(resp)
}
