package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"umami-go/internal/model"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type HealthRepository struct {
	DB  *gorm.DB
	Log *logrus.Logger
}

func NewHealthRepository(log *logrus.Logger) *HealthRepository {
	return &HealthRepository{
		Log: log,
	}
}

func (r *HealthRepository) CheckDB(tx *gorm.DB, health *model.CheckDBResponse) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Ping the database
	db, err := tx.DB()
	if err != nil {
		health.Status = "down"
		health.Message = fmt.Sprintf("db down: %v", err)
		r.Log.Fatalf(fmt.Sprintf("db down: %v", err))
		return err
	}

	err = db.PingContext(ctx)
	if err != nil {
		health.Status = "down"
		health.Message = fmt.Sprintf("db down: %v", err)
		r.Log.Fatalf(fmt.Sprintf("db down: %v", err))
		return err
	}

	// Database is up, add more statistics
	health.Status = "up"
	health.Message = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := db.Stats()
	health.OpenConnections = strconv.Itoa(dbStats.OpenConnections)
	health.InUse = strconv.Itoa(dbStats.InUse)
	health.Idle = strconv.Itoa(dbStats.Idle)
	health.WaitCount = strconv.FormatInt(dbStats.WaitCount, 10)
	health.WaitDuration = dbStats.WaitDuration.String()
	health.MaxIdleClosed = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	health.MaxLifetimeClosed = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		health.Message = "The database is experiencing heavy load."
	}
	if dbStats.WaitCount > 1000 {
		health.Message = "The database has a high number of wait events, indicating potential bottlenecks."
	}
	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		health.Message = "Many idle connections are being closed, consider revising the connection pool settings."
	}
	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		health.Message = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return nil
}

func (r *HealthRepository) CheckRedis(ctx context.Context, redisClient *redis.Client, health *model.CheckRedisResponse) error {
	status := redisClient.Ping(ctx)
	if status.Err() != nil {
		health.Status = "down"
		health.Message = fmt.Sprintf("redis down: %v", status.Err())
		r.Log.Fatalf(fmt.Sprintf("redis down: %v", status.Err()))
		return status.Err()
	}

	// Redis is up, add more statistics
	health.Status = "up"
	health.Message = "It's healthy"

	// Get Redis stats
	info, err := redisClient.Info(ctx).Result()
	if err != nil {
		health.Status = "up"
		health.Message = fmt.Sprintf("failed to get redis info: %v", err)
		r.Log.Warnf(fmt.Sprintf("failed to get redis info: %v", err))
		return err
	}

	// Parse the info string (this is a simplified example, in reality you may want to parse specific fields)
	parsedInfo := parseRedisInfo(info)
	health.Info = parsedInfo

	return nil
}

func parseRedisInfo(info string) map[string]interface{} {
	parsedInfo := make(map[string]interface{})
	lines := strings.Split(info, "\r\n")
	currentSection := ""

	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			currentSection = strings.ToLower(strings.TrimSpace(strings.TrimPrefix(line, "#")))
			parsedInfo[currentSection] = make(map[string]string)
		} else if line != "" {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				sectionMap := parsedInfo[currentSection].(map[string]string)
				sectionMap[parts[0]] = parts[1]
			}
		}
	}

	return parsedInfo
}
