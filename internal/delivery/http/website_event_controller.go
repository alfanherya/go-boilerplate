package http

import (
	"umami-go/internal/model"
	"umami-go/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type WebsiteEventController struct {
	Log     *logrus.Logger
	UseCase *usecase.WebsiteEventUseCase
}

func NewWebsiteEventController(useCase *usecase.WebsiteEventUseCase, logger *logrus.Logger) *WebsiteEventController {
	return &WebsiteEventController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *WebsiteEventController) Active(ctx *fiber.Ctx) error {
	if err := ctx.Params("id"); err == "" {
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Active(ctx.Context(), ctx.Params("id"))
	if err != nil {
		c.Log.Warnf("Failed to get website events : %+v", err)
		return err
	}

	return ctx.JSON(response)
}

func (c *WebsiteEventController) Stats(ctx *fiber.Ctx) error {
	if err := ctx.Params("id"); err == "" {
		return fiber.ErrBadRequest
	}

	query := new(model.WEStatsReq)
	if err := ctx.QueryParser(query); err != nil {
		c.Log.Warnf("Failed to parse query params : %+v", err)
		return fiber.ErrBadRequest
	}

	query.WebsiteID = ctx.Params("id")

	response, err := c.UseCase.Stats(ctx.Context(), query)
	if err != nil {
		c.Log.Warnf("Failed to get website events : %+v", err)
		return err
	}

	return ctx.JSON(response)
}

func (c *WebsiteEventController) PageViews(ctx *fiber.Ctx) error {
	if err := ctx.Params("id"); err == "" {
		return fiber.ErrBadRequest
	}

	query := new(model.WEPageViewsReq)
	if err := ctx.QueryParser(query); err != nil {
		c.Log.Warnf("Failed to parse query params : %+v", err)
		return fiber.ErrBadRequest
	}

	query.WebsiteID = ctx.Params("id")

	response, err := c.UseCase.PageViews(ctx.Context(), query)
	if err != nil {
		c.Log.Warnf("Failed to get website page views : %+v", err)
		return err
	}

	return ctx.JSON(response)
}

func (c *WebsiteEventController) Metrics(ctx *fiber.Ctx) error {
	if err := ctx.Params("id"); err == "" {
		return fiber.ErrBadRequest
	}

	query := new(model.WEMetricsReq)
	if err := ctx.QueryParser(query); err != nil {
		c.Log.Warnf("Failed to parse query params : %+v", err)
		return fiber.ErrBadRequest
	}

	query.WebsiteID = ctx.Params("id")

	response, err := c.UseCase.Metrics(ctx.Context(), query)
	if err != nil {
		c.Log.Warnf("Failed to get website page views : %+v", err)
		return err
	}

	return ctx.JSON(response)
}
