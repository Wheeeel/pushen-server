package api

import (
	"fmt"

	"github.com/Wheeeel/pushen-server/api/handler"
	"github.com/Wheeeel/pushen-server/config"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/handlerconv"
	"github.com/kataras/iris/websocket"
	corscore "github.com/rs/cors"
)

type Server struct {
	app  *iris.Application
	Addr string
}

func New(addr string) (srv *Server, err error) {
	srv = &Server{
		app:  iris.New(),
		Addr: addr,
	}

	ws := websocket.New(websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	})

	ws.OnConnection(handler.SendMessageHandler)
	wsHandler := func(ctx iris.Context) {
		ws.Handler()(ctx)
	}

	srv.app.Use(func(ctx iris.Context) {
		fmt.Printf("%s %s \n", ctx.Method(), ctx.Request().URL.String())
		ctx.Next()
	})

	// cors
	ch := corscore.AllowAll().ServeHTTP
	h := handlerconv.FromStdWithNext(ch)
	srv.app.Use(h)
	srv.app.Use(handler.CorsHandler)

	// filter
	srv.app.Use(handler.AuthHandler)

	// error handler
	srv.app.OnErrorCode(iris.StatusForbidden, handler.ErrorForbidden)
	srv.app.OnErrorCode(iris.StatusNotFound, handler.ErrorNotFound)
	srv.app.OnErrorCode(iris.StatusBadRequest, handler.ErrorBadRequest)

	// router
	// user related
	srv.app.Post("/users/signup", handler.UserCreateHandler)
	srv.app.Post("/users/signin", handler.UserLoginHandler)
	srv.app.Post("/users/logout", handler.UserLogoutHandler)
	srv.app.Get("/me", handler.UserInfoHandler)
	srv.app.Post("/messages", handler.ReceiveMessageHandler)
	srv.app.Get("/messages", wsHandler)
	srv.app.Post("/devices/bind", handler.BindAuthTokenHandler)
	srv.app.Get("/token", handler.DeviceBindTokenHandler)

	srv.app.Use(func (ctx iris.Context) {
		fmt.Println("post")
		ctx.Next()
	})

	return
}

func (srv *Server) Run() (err error) {
	return srv.app.Run(iris.Addr(config.DefaultAppConfig.Addr))
}
