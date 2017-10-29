package request

type DeviceUnbind struct {
	Device string `json:"device" validate:"required"`
}
