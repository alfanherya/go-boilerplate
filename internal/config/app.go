package config

import (
	"umami-go/internal/delivery/http"
	"umami-go/internal/delivery/http/middleware"
	"umami-go/internal/delivery/http/route"
	"umami-go/internal/repository"
	"umami-go/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
	Redis    *redis.Client
}

func Bootstrap(config *BootstrapConfig) {
	// setup repositories
	healthRepository := repository.NewHealthRepository(config.Log)
	userRepository := repository.NewUserRepository(config.Log)
	websiteRepository := repository.NewWebsiteRepository(config.Log)
	websiteEventRepository := repository.NewWebsiteEventRepository(config.Log)

	// setup use cases
	healthUseCase := usecase.NewHealthUseCase(config.DB, config.Redis, config.Log, config.Validate, healthRepository)
	userUseCase := usecase.NewUserUseCase(config.DB, config.Redis, config.Log, config.Validate, userRepository)
	websiteUseCase := usecase.NewWebsiteUseCase(config.DB, config.Redis, config.Log, config.Validate, websiteRepository)
	websiteEventUseCase := usecase.NewWebsiteEventUseCase(config.DB, config.Redis, config.Log, config.Validate, websiteEventRepository, websiteRepository)

	// setup controller
	healthController := http.NewHealthController(healthUseCase, config.Log)
	userController := http.NewUserController(userUseCase, config.Log)
	websiteController := http.NewWebsiteController(websiteUseCase, config.Log)
	websiteEventController := http.NewWebsiteEventController(websiteEventUseCase, config.Log)

	// setup middleware
	authMiddleware := middleware.NewAuth(config.Config)

	routeConfig := route.RouteConfig{
		App:                    config.App,
		HealthController:       healthController,
		UserController:         userController,
		WebsiteController:      websiteController,
		WebsiteEventController: websiteEventController,
		AuthMiddleware:         authMiddleware,
	}

	routeConfig.Setup()
}
