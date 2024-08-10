package http

import (
	"umami-go/internal/delivery/http/middleware"
	"umami-go/internal/model"
	"umami-go/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log     *logrus.Logger
	UseCase *usecase.UserUseCase
}

func NewUserController(useCase *usecase.UserUseCase, logger *logrus.Logger) *UserController {
	return &UserController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Login(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to login user : %+v", err)
		return err
	}

	return ctx.JSON(response)
}

func (c *UserController) Verify(ctx *fiber.Ctx) error {
	userID := middleware.GetUserID(ctx)

	response, err := c.UseCase.Verify(ctx.UserContext(), userID)
	if err != nil {
		c.Log.Warnf("Failed to verify user : %+v", err)
		return err
	}

	return ctx.JSON(response)
}
