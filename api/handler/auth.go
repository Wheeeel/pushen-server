package handler

import (
	"net/http"

	"github.com/Wheeeel/pushen-server/api/util"
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
		switch ctx.Method() {
		case http.MethodPost:
			tokenStr := ctx.GetHeader("X-Auth-Token")
			if tokenStr == "" {
				util.WriteError(ctx, iris.StatusForbidden, "auth error", "auth handler error: token empty")
				return
			}
			token, err := model.AuthTokenByToken(model.DefaultDB, tokenStr)
			if err != nil {
				util.WriteErrorf(ctx, iris.StatusForbidden, "auth error", "auth handler error: %v", err.Error())
				return
			}

			ctx.Values().Set("userID", token.UserID)
		case http.MethodGet:
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
			util.WriteErrorf(ctx, iris.StatusForbidden, "auth error", "auth handler error: %v", err.Error())
			return
		}
		user, err := model.UserByID(model.DefaultDB, token.UserID)
		if err != nil {
			util.WriteErrorf(ctx, iris.StatusForbidden, "auth error", "auth handler error: %v", err.Error())
			return
		}

		ctx.Values().Set("email", user.Email)
		ctx.Next()
		return
	}

	// check cookie
	s := session.Start(ctx)
	if auth, _ := s.GetBoolean(authenticated); !auth {
		util.WriteError(ctx, iris.StatusForbidden, "auth error", "auth handler error: header auth flag empty")
		return
	}

	var email string
	if email = s.GetString("email"); email == "" {
		util.WriteError(ctx, iris.StatusForbidden, "auth error", "auth handler error: email empty")
		return
	}

	// set email
	ctx.Values().Set("email", email)
	ctx.Next()
}
