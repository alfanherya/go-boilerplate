package usecase

import (
	"context"
	"umami-go/internal/entity"
	"umami-go/internal/model"
	"umami-go/internal/model/converter"
	"umami-go/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WebsiteEventUseCase struct {
	DB                *gorm.DB
	Redis             *redis.Client
	Log               *logrus.Logger
	Validate          *validator.Validate
	Repository        *repository.WebsiteEventRepository
	WebsiteRepository *repository.WebsiteRepository
}

func NewWebsiteEventUseCase(db *gorm.DB, redis *redis.Client, logger *logrus.Logger, validate *validator.Validate, websiteEventRepository *repository.WebsiteEventRepository, websiteRepository *repository.WebsiteRepository) *WebsiteEventUseCase {
	return &WebsiteEventUseCase{
		DB:                db,
		Redis:             redis,
		Log:               logger,
		Validate:          validate,
		Repository:        websiteEventRepository,
		WebsiteRepository: websiteRepository,
	}
}

func (u *WebsiteEventUseCase) Active(ctx context.Context, id string) (*model.WECount, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.WebsiteRepository.FindByWebsiteID(tx, new(entity.Website), id); err != nil {
		u.Log.Warnf("Failed to get website : %+v", err)
		return nil, fiber.ErrNotFound
	}

	count, err := u.Repository.Active(tx, new(entity.WebsiteEvent), id)
	if err != nil {
		u.Log.Warnf("Failed get website event : %+v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return &model.WECount{X: count}, nil
}

func (u *WebsiteEventUseCase) Stats(ctx context.Context, req *model.WEStatsReq) (*model.WEStatsRes, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(req); err != nil {
		u.Log.Warnf("Invalid request body : %+v", req)
		return nil, err
	}

	if err := u.WebsiteRepository.FindByWebsiteID(tx, new(entity.Website), req.WebsiteID); err != nil {
		u.Log.Warnf("Failed to get website : %+v", err)
		return nil, fiber.ErrNotFound
	}

	statsNow := new(model.WEStats)
	if err := u.Repository.Stats(tx, statsNow, req); err != nil {
		u.Log.Warnf("Failed to get website event stats : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	req.StartAt = req.StartAt - 86400000 // 86400000 is 1 day in milliseconds
	req.EndAt = req.EndAt - 86400000     // 86400000 is 1 day in milliseconds

	statsPrev := new(model.WEStats)
	if err := u.Repository.Stats(tx, statsPrev, req); err != nil {
		u.Log.Warnf("Failed to get website event stats : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WEStatsToResponse(statsNow, statsPrev), nil
}

func (u *WebsiteEventUseCase) PageViews(ctx context.Context, req *model.WEPageViewsReq) (*model.WEPageViewsRes, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(req); err != nil {
		u.Log.Warnf("Invalid request body : %+v", req)
		return nil, err
	}

	if err := u.WebsiteRepository.FindByWebsiteID(tx, new(entity.Website), req.WebsiteID); err != nil {
		u.Log.Warnf("Failed to get website : %+v", err)
		return nil, fiber.ErrNotFound
	}

	pageStats := new([]model.XY)
	sessionStats := new([]model.XY)
	if err := u.Repository.PageViews(tx, pageStats, sessionStats, req); err != nil {
		u.Log.Warnf("Failed to get website event page views : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WEPageViewsToResponse(pageStats, sessionStats), nil
}

func (u *WebsiteEventUseCase) Metrics(ctx context.Context, req *model.WEMetricsReq) (*[]model.XY, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(req); err != nil {
		u.Log.Warnf("Invalid request body : %+v", req)
		return nil, err
	}

	if err := u.WebsiteRepository.FindByWebsiteID(tx, new(entity.Website), req.WebsiteID); err != nil {
		u.Log.Warnf("Failed to get website : %+v", err)
		return nil, fiber.ErrNotFound
	}

	result := new([]model.XY)

	if err := u.Repository.Metrics(tx, result, req); err != nil {
		u.Log.Warnf("Failed to get website event metrics : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if req.Type == "language" {
		return converter.WEMetricsLangToResponse(result), nil
	}

	return result, nil
}
