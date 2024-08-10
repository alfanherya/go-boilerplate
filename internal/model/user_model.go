package model

import "time"

type LoginUserResponse struct {
	User  LoginUser `json:"user"`
	Token string    `json:"token"`
}

type LoginUser struct {
	ID        string     `json:"id"`
	Username  string     `json:"username"`
	Roles     string     `json:"role"`
	CreatedAt *time.Time `json:"createdAt"`
	IsAdmin   bool       `json:"isAdmin"`
}

type LoginUserRequest struct {
	Username string `json:"username" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type VerifyUserRequest struct {
	ID string `json:"id" validate:"required"`
}
