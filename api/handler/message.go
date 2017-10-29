package handler

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Wheeeel/pushen-server/api/request"
	apiutil "github.com/Wheeeel/pushen-server/api/util"
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
		apiutil.WriteErrorf(ctx, iris.StatusBadRequest, err.Error(), "receive message error: %v", err.Error())
		return
	}

	validate := validator.New()
	err = validate.Struct(mc)
	if err != nil {
		apiutil.WriteErrorf(ctx, iris.StatusBadRequest, err.Error(), "receive message error: %v", err.Error())
		return
	}

	device, err := model.DeviceByUUID(model.DefaultDB, mc.DeviceID)
	if err != nil {
		apiutil.WriteErrorf(ctx, iris.StatusBadRequest, err.Error(), "receive message error: %v", err.Error())
		return
	}

	var msg model.Message
	err = util.CopyStruct(&msg, mc)
	if err != nil {
		apiutil.WriteErrorf(ctx, iris.StatusInternalServerError, err.Error(), "receive message error: %v", err.Error())
		return
	}

	tx := model.DefaultDB.Begin()
	msg.DeviceID = device.ID
	err = model.MessageCreate(tx, &msg)
	if err != nil {
		tx.Rollback()
		apiutil.WriteErrorf(ctx, iris.StatusInternalServerError, err.Error(), "receive message error: %v", err.Error())
		return
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		apiutil.WriteErrorf(ctx, iris.StatusInternalServerError, err.Error(), "receive message error: %v", err.Error())
		return
	}

	var resp Response
	resp.Code = 200
	resp.Msg = "receive message success"
	ctx.JSON(resp)
}

func SendMessageHandler(c websocket.Connection) {
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	for range ticker.C {
		ms, err := model.MessageOrderByCreateTimestamp(model.DefaultDB)
		if err != nil {
			//c.Context()..Application().Logger().Infof("receive message error: %v", err)
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
			err = c.EmitMessage(b)
			if err != nil {
				log.Println(err)
				continue
			}

			tx := model.DefaultDB.Begin()
			m.Status = model.MessageStatusSent
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
