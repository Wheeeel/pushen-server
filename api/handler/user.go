package handler

import (
	"math/rand"
	"time"

	"github.com/Wheeeel/pushen-server/api/request"
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

	rand.Seed(time.Now().Unix())
	user.AvatarURL = avatar[rand.Intn(5)]
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
	email := ctx.Values().GetString("email")
	if email == "" {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.Values().Set("error", "auth error")
		return
	}

	user, err := model.UserByEmail(email)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.Values().Set("error", err.Error())
		return
	}

	var resp Response
	resp.Code = 200
	resp.Msg = "ok"
	resp.Data = user
	ctx.JSON(resp)
}
