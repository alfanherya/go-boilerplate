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

type WebsiteUseCase struct {
	DB                *gorm.DB
	Redis             *redis.Client
	Log               *logrus.Logger
	Validate          *validator.Validate
	WebsiteRepository *repository.WebsiteRepository
}

func NewWebsiteUseCase(db *gorm.DB, redis *redis.Client, logger *logrus.Logger, validate *validator.Validate, websiteRepository *repository.WebsiteRepository) *WebsiteUseCase {
	return &WebsiteUseCase{
		DB:                db,
		Redis:             redis,
		Log:               logger,
		Validate:          validate,
		WebsiteRepository: websiteRepository,
	}
}

func (u *WebsiteUseCase) Create(ctx context.Context, request *model.WebsiteCreateRequest) (*model.WebsiteResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body : %+v", request)
		return nil, err
	}

	website := &entity.Website{
		ID:     u.WebsiteRepository.GenerateUUID(),
		Name:   request.Name,
		Domain: &request.Domain,
	}

	if request.TeamID == "" { // need to validate teamID
		website.UserID = &request.UserID
	}

	if err := u.WebsiteRepository.Create(tx, website); err != nil {
		u.Log.Warnf("Failed create website : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("error creating website")
		return nil, fiber.ErrInternalServerError
	}

	return converter.WebsiteToResponse(website), nil
}

func (u *WebsiteUseCase) Websites(ctx context.Context, query *model.WebsitesRequest) (*model.WebsitesResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	websites := new([]entity.Website) // <- *[]entity.Website

	if err := u.WebsiteRepository.FindByUserWebsites(tx, websites, query); err != nil {
		u.Log.Warnf("Failed find websites : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WebsitesToResponse(websites, query), nil
}

func (u *WebsiteUseCase) Website(ctx context.Context, websiteID string) (*model.WebsiteResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	website := new(entity.Website)
	if err := u.WebsiteRepository.FindByWebsiteID(tx, website, websiteID); err != nil {
		u.Log.Warnf("Failed find website : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WebsiteToResponse(website), nil
}

func (u *WebsiteUseCase) Delete(ctx context.Context, websiteID string) (*model.WebsiteResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	website := new(entity.Website)

	if err := u.WebsiteRepository.FindByWebsiteID(tx, website, websiteID); err != nil {
		u.Log.Warnf("Failed find website : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if err := u.WebsiteRepository.Delete(tx, website); err != nil {
		u.Log.Warnf("Failed delete website : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("error deleting website")
		return nil, fiber.ErrInternalServerError
	}

	return converter.WebsiteToResponse(website), nil
}

func (u *WebsiteUseCase) Update(ctx context.Context, request *model.WebsiteUpdateRequest) (*model.WebsiteResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body : %+v", request)
		return nil, err
	}

	website := new(entity.Website)

	if err := u.WebsiteRepository.FindByWebsiteID(tx, website, request.ID); err != nil {
		u.Log.Warnf("Failed find website : %+v", err)
		return nil, fiber.ErrNotFound
	}

	website.Name = request.Name
	website.Domain = &request.Domain
	website.ShareID = &request.ShareID

	if err := u.WebsiteRepository.Update(tx, website); err != nil {
		u.Log.Warnf("Failed update website : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("error updating website")
		return nil, fiber.ErrInternalServerError
	}

	return converter.WebsiteToResponse(website), nil
}
