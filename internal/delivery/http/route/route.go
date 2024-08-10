package route

import (
	"umami-go/internal/delivery/http"

	"github.com/gofiber/contrib/fibersentry"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

type RouteConfig struct {
	App                    *fiber.App
	HealthController       *http.HealthController
	UserController         *http.UserController
	WebsiteController      *http.WebsiteController
	WebsiteEventController *http.WebsiteEventController
	AuthMiddleware         fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.App.Get("/metrics", monitor.New())
	c.App.Use(fibersentry.New(fibersentry.Config{Repanic: true, WaitForDelivery: true}))

	c.SetupHealthRoute()

	c.SetupAuthRoute()
	c.SetupUserRoute()
	c.SetupWebsiteRoute()
}

func (c *RouteConfig) SetupAuthRoute() {
	group := c.App.Group("/auth")
	group.Post("/login", c.UserController.Login)
}

func (c *RouteConfig) SetupUserRoute() {
	c.App.Use(c.AuthMiddleware)
	group := c.App.Group("/user")
	group.Get("/verify", c.UserController.Verify)
}

func (c *RouteConfig) SetupWebsiteRoute() {
	c.App.Use(c.AuthMiddleware)
	group := c.App.Group("/websites")
	group.Get("/", c.WebsiteController.Websites)
	group.Post("/", c.WebsiteController.Create)
	group.Get("/:id", c.WebsiteController.Website)
	group.Post("/:id", c.WebsiteController.Update)
	group.Delete("/:id", c.WebsiteController.Delete)

	group.Get("/:id/active", c.WebsiteEventController.Active)
	group.Get("/:id/stats", c.WebsiteEventController.Stats)
	group.Get("/:id/page-views", c.WebsiteEventController.PageViews)
	group.Get("/:id/metrics", c.WebsiteEventController.Metrics)
}

func (c *RouteConfig) SetupHealthRoute() {
	group := c.App.Group("/health")
	group.Get("/", c.HealthController.All)
	group.Get("/db", c.HealthController.CheckDB)
	group.Get("/redis", c.HealthController.CheckRedis)
}
