package config

import (
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewSentry(config *viper.Viper, log *logrus.Logger) {
	dsn := config.GetString("sentry.dsn")

	sentry.Init(sentry.ClientOptions{
		Dsn: dsn,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			if hint.Context != nil {
				if c, ok := hint.Context.Value(sentry.RequestContextKey).(*fiber.Ctx); ok {
					log.WithFields(logrus.Fields{"hostname": utils.CopyString(c.Hostname())}).Error("error")
				}
			}
			return event
		},
		Debug:              true,
		AttachStacktrace:   true,
		EnableTracing:      true,
		ProfilesSampleRate: 1.0,
		TracesSampleRate:   1.0,
	})
}
