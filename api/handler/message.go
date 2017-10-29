package handler

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Wheeeel/pushen-server/api/request"
	"github.com/Wheeeel/pushen-server/model"
	"github.com/Wheeeel/pushen-server/util"
	"github.com/go-playground/validator"
	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
)

func ReceiveMessageHandler(ctx iris.Context) {
	var mc request.MessageCreate
	err := ctx.ReadJSON(&mc)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Values().Set("error", err.Error())
		return
	}

	validate := validator.New()
	err = validate.Struct(mc)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Values().Set("error", err.Error())
		return
	}

	var msg model.Message
	err = util.CopyStruct(&msg, mc)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.Values().Set("error", err.Error())
		return
	}

	tx := model.DefaultDB.Begin()
	err = model.MessageCreate(tx, &msg)
	if err != nil {
		tx.Rollback()
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.Values().Set("error", err.Error())
		return
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.Values().Set("error", err.Error())
		return
	}

	var resp Response
	resp.Code = 200
	resp.Msg = "create message success"
	ctx.JSON(resp)
}

func SendMessageHandler(c websocket.Connection) {
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	for range ticker.C {
		ms, err := model.MessageByCreateTimestamp(model.DefaultDB)
		if err != nil {
			log.Println(err)
			continue
		}
		if len(ms) == 0 {
			continue
		}

		for _, m := range ms {
			b, err := json.Marshal(m)
			if err != nil {
				log.Println(err)
				continue
			}
			c.EmitMessage(b)

			tx := model.DefaultDB.Begin()
			m.Status = model.MessageStatusSendt
			err = model.MessageUpdateStatus(tx, &m)
			if err != nil {
				tx.Rollback()
				log.Println(err)
				continue
			}
			if err = tx.Commit().Error; err != nil {
				tx.Rollback()
				log.Println(err)
				continue
			}
		}
	}
}
