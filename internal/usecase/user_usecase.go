package usecase

import (
	"context"
	"time"
	"umami-go/internal/entity"
	"umami-go/internal/model"
	"umami-go/internal/model/converter"
	"umami-go/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB             *gorm.DB
	Redis          *redis.Client
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
}

func NewUserUseCase(db *gorm.DB, redis *redis.Client, logger *logrus.Logger, validate *validator.Validate, userRepository *repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Redis:          redis,
		Log:            logger,
		Validate:       validate,
		UserRepository: userRepository,
	}
}

func (c *UserUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.LoginUserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", request)
		return nil, err
	}

	user := new(entity.User)
	if err := c.UserRepository.FindByUsername(tx, user, request.Username); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Log.Warnf("Failed compare password : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := c.Redis.Set(ctx, user.ID, user.ID, 0).Err(); err != nil {
		c.Log.Warnf("Failed set redis : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	claims := jwt.MapClaims{
		"userId":  user.ID,
		"isAdmin": user.Role == "admin",
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 1 day
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		c.Log.Warnf("Failed signed token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserLoginToResponse(user, t), nil
}

func (c *UserUseCase) Verify(ctx context.Context, id string) (*model.LoginUser, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := new(entity.User)
	if err := c.UserRepository.FindByUserID(tx, user, id); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserVerifyToResponse(user), nil
}
