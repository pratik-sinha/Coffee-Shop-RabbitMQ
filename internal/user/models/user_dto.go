package models

type RegisterReq struct {
	Email    string `json:"email" validate:"required"`
	UserName string `json:"user_name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginReq struct {
	UserName string `json:"user_name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LogoutReq struct {
	SessionId string
}

type LoginRes struct {
	Token string `json:"token"`
}
