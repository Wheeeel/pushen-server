package request

type UserCreate struct {
	Email              string    `json:"email" validate:"required,email"`
	Name               string    `json:"name" validate:"required"`
	Password           string    `json:"password" validate:"required"`
	AvatarURL          string    `json:"avatarUrl" validate:"omitempty,url"`
}

type UserLogin struct {
	Email              string    `json:"email" validate:"required,email"`
	Password           string    `json:"password" validate:"required"`
}
