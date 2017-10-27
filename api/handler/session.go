package handler

import (
	"time"

	"github.com/kataras/iris/sessions"
)

var (
	cookieName    = "PUSHEN"
	session       = sessions.New(sessions.Config{Cookie: cookieName, Expires: time.Minute * 30})
	authenticated = "authenticated"
)
