package handler

import (
	"github.com/Wheeeel/pushen-server/model"
	"github.com/kataras/iris"
)

func AuthHandler(ctx iris.Context) {
	// white list
	switch ctx.Path() {
	case "/users/signup", "/users/signin", "/devices/bind", "/debug":
		ctx.Next()
		return
	case "/messages":
		token := ctx.GetHeader("X-Auth-Token")
		if token == "" {
			ctx.StatusCode(iris.StatusForbidden)
			ctx.Values().Set("error", "auth error")
			return
		}
		ctx.Next()
		return
	case "/me":
		tokenStr := ctx.GetHeader("X-Auth-Token")
		// use cookie
		if tokenStr == "" {
			break
		}

		// use token
		token, err := model.AuthTokenByToken(model.DefaultDB, tokenStr)
		if err != nil {
			ctx.StatusCode(iris.StatusForbidden)
			ctx.Values().Set("error", "auth error")
			return
		}

		user, err := model.UserByID(model.DefaultDB, token.UserID)
		if err != nil {
			ctx.StatusCode(iris.StatusForbidden)
			ctx.Values().Set("error", "auth error")
			return
		}

		ctx.Values().Set("email", user.Email)
		ctx.Next()
		return
	}

	// check whether use cookie or token
	if isUseToken := ctx.GetHeader(useToken); isUseToken == "1" || isUseToken == "true" {
		// check token
		// TODO
	}

	// check cookie
	s := session.Start(ctx)
	if auth, _ := s.GetBoolean(authenticated); !auth {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.Values().Set("error", "auth error")
		return
	}

	var email string
	if email = s.GetString("email"); email == "" {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.Values().Set("error", "auth error")
		return
	}

	// set email
	ctx.Values().Set("email", email)

	//session.ShiftExpiration(ctx)
	ctx.Next()
}
