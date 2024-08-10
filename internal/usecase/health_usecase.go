package usecase

import (
	"context"
	"sync"
	"time"
	"umami-go/internal/model"
	"umami-go/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type HealthUseCase struct {
	DB               *gorm.DB
	Redis            *redis.Client
	Log              *logrus.Logger
	Validate         *validator.Validate
	HealthRepository *repository.HealthRepository
}

func NewHealthUseCase(db *gorm.DB, redis *redis.Client, logger *logrus.Logger, validate *validator.Validate, healthRepository *repository.HealthRepository) *HealthUseCase {
	return &HealthUseCase{
		DB:               db,
		Redis:            redis,
		Log:              logger,
		Validate:         validate,
		HealthRepository: healthRepository,
	}
}

func (c *HealthUseCase) All(ctx context.Context) (*model.HealthResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	health := new(model.HealthResponse)
	var wg sync.WaitGroup
	var dbErr, redisErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := c.HealthRepository.CheckDB(tx, &health.Database); err != nil {
			c.Log.Warnf("Failed check db : %+v", err)
			dbErr = fiber.ErrInternalServerError
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := c.HealthRepository.CheckRedis(ctx, c.Redis, &health.Redis); err != nil {
			c.Log.Warnf("Failed check redis : %+v", err)
			redisErr = fiber.ErrInternalServerError
		}
	}()

	wg.Wait()

	if dbErr != nil {
		return nil, dbErr
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if redisErr != nil {
		return nil, redisErr
	}

	return health, nil
}

func (c *HealthUseCase) CheckDB(ctx context.Context) (*model.CheckDBResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	health := new(model.CheckDBResponse)

	if err := c.HealthRepository.CheckDB(tx, health); err != nil {
		c.Log.Warnf("Failed check db : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return health, nil
}

func (c *HealthUseCase) CheckRedis(ctx context.Context) (*model.CheckRedisResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	health := new(model.CheckRedisResponse)

	if err := c.HealthRepository.CheckRedis(ctx, c.Redis, health); err != nil {
		c.Log.Warnf("Failed check redis : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return health, nil
}
