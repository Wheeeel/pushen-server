package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/Wheeeel/pushen-server/api/request"
	apiutil "github.com/Wheeeel/pushen-server/api/util"
	"github.com/Wheeeel/pushen-server/model"
	"github.com/Wheeeel/pushen-server/util"
	"github.com/go-playground/validator"
	"github.com/gorilla/websocket"
	"github.com/kataras/iris"
	"github.com/pkg/errors"
)

var upgrader websocket.Upgrader

func init() {
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
}

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

func SendMessageHandler(ctx iris.Context) {
	conn, err := upgrader.Upgrade(ctx.ResponseWriter(), ctx.Request(), nil)
	if err != nil {
		log.Printf("send message handler upgrade error: %v", err)
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	for range ticker.C {
		err := sendMessage(conn)
		if err != nil {
			log.Printf("write message error: %v", err)
			break
		}
	}
}

func sendMessage(c *websocket.Conn) (err error) {
	ms, err := model.MessageOrderByCreateTimestamp(model.DefaultDB)
	if err != nil {
		err = errors.Wrap(err, "send message error")
		return
	}
	if len(ms) == 0 {
		return
	}

	for _, m := range ms {
		err = c.WriteJSON(m)
		if err != nil {
			log.Println(err)
			break
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
	return
}
