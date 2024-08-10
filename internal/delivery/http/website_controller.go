package http

import (
	"umami-go/internal/delivery/http/middleware"
	"umami-go/internal/model"
	"umami-go/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type WebsiteController struct {
	Log     *logrus.Logger
	UseCase *usecase.WebsiteUseCase
}

func NewWebsiteController(useCase *usecase.WebsiteUseCase, logger *logrus.Logger) *WebsiteController {
	return &WebsiteController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *WebsiteController) Create(ctx *fiber.Ctx) error {
	request := new(model.WebsiteCreateRequest)

	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	if request.TeamID == "" {
		userID := middleware.GetUserID(ctx)
		request.UserID = userID
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create website : %+v", err)
		return err
	}

	return ctx.JSON(response)
}

func (c *WebsiteController) Websites(ctx *fiber.Ctx) error {
	query := new(model.WebsitesRequest)
	if err := ctx.QueryParser(query); err != nil {
		c.Log.Warnf("Failed to parse query params : %+v", err)
		return fiber.ErrBadRequest
	}

	query.UserID = middleware.GetUserID(ctx)

	if query.Query == "" {
		query.Query = "%" + query.Query + "%"
	}

	if query.OrderBy == "" {
		query.OrderBy = "created_at desc" // butuh di tuning
	}

	if query.Page < 1 {
		query.Page = 1
	}

	if query.PageSize < 1 {
		query.PageSize = 10
	}

	response, err := c.UseCase.Websites(ctx.UserContext(), query)
	if err != nil {
		c.Log.Warnf("Failed to get websites : %+v", err)
		return err
	}

	return ctx.JSON(response)
}

func (c *WebsiteController) Website(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Website(ctx.UserContext(), id)
	if err != nil {
		c.Log.Warnf("Failed to get website : %+v", err)
		return err
	}

	return ctx.JSON(response)
}

func (c *WebsiteController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Delete(ctx.UserContext(), id)
	if err != nil {
		c.Log.Warnf("Failed to delete website : %+v", err)
		return err
	}

	return ctx.JSON(response)
}

func (c *WebsiteController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.ErrBadRequest
	}

	request := new(model.WebsiteUpdateRequest)
	request.ID = id

	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to delete website : %+v", err)
		return err
	}

	return ctx.JSON(response)
}
