package request

type MessageCreate struct {
	AppName  string `json:"appName" validate:"required"`
	AppIcon  string `json:"appIcon"`
	Body     string `json:"messageBody" validate:"required"`
	DeviceID string `json:"deviceID" validate:"required"`
}
