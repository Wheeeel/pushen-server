package handler

import (
	"github.com/Wheeeel/pushen-server/model"
	"github.com/kataras/iris"
	"github.com/pkg/errors"
)

func AuthHandler(ctx iris.Context) {
	// white list
	switch ctx.Path() {
	case "/users/signup", "/users/signin", "devices/bind":
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

func deviceAuth(tokenStr string) (token model.AuthToken, err error) {
	token, err = model.AuthTokenByToken(tokenStr)
	if err != nil {
		err = errors.Wrap(err, "device auth error")
		return
	}
	return
}
