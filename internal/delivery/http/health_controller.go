package http

import (
	"umami-go/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type HealthController struct {
	Log     *logrus.Logger
	UseCase *usecase.HealthUseCase
}

func NewHealthController(useCase *usecase.HealthUseCase, logger *logrus.Logger) *HealthController {
	return &HealthController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *HealthController) All(ctx *fiber.Ctx) error {
	response, err := c.UseCase.All(ctx.UserContext())
	if err != nil {
		c.Log.Warnf("Failed to get health : %+v", err)
		return err
	}

	return ctx.JSON(response)
}

func (c *HealthController) CheckDB(ctx *fiber.Ctx) error {
	response, err := c.UseCase.CheckDB(ctx.UserContext())
	if err != nil {
		c.Log.Warnf("Failed to check db : %+v", err)
		return err
	}
	return ctx.JSON(response)
}

func (c *HealthController) CheckRedis(ctx *fiber.Ctx) error {
	response, err := c.UseCase.CheckRedis(ctx.UserContext())
	if err != nil {
		c.Log.Warnf("Failed to check redis : %+v", err)
		return err
	}
	return ctx.JSON(response)
}
