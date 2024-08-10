package converter

import (
	"umami-go/internal/entity"
	"umami-go/internal/model"
)

func UserLoginToResponse(user *entity.User, token string) *model.LoginUserResponse {
	return &model.LoginUserResponse{
		User: model.LoginUser{
			ID:        user.ID,
			Roles:     user.Role,
			IsAdmin:   user.Role == "admin",
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		},
		Token: token,
	}
}

func UserVerifyToResponse(user *entity.User) *model.LoginUser {
	return &model.LoginUser{
		ID:        user.ID,
		Username:  user.Username,
		IsAdmin:   user.Role == "admin",
		CreatedAt: user.CreatedAt,
		Roles:     user.Role,
	}
}
