package response

import "github.com/Wheeeel/pushen-server/model"

type Message struct {
	model.Message

	DeviceUUID string `json:"deviceUUID"`
}
