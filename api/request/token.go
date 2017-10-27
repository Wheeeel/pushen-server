package request

type DeviceAuthToken struct {
	Token string `json:"token" validate:"required"`
}

type DeviceBindToken struct {
	Token string `json:"token" validate:"required"`
}
