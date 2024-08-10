package repository

import (
	"umami-go/internal/entity"
	"umami-go/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WebsiteRepository struct {
	Repository[entity.Website]
	Log *logrus.Logger
}

func NewWebsiteRepository(log *logrus.Logger) *WebsiteRepository {
	return &WebsiteRepository{
		Log: log,
	}
}

func (r *WebsiteRepository) FindByWebsiteID(db *gorm.DB, website *entity.Website, websiteID string) error {
	query := db.Where("website_id = ?", websiteID)

	if err := query.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("user_id", "username")
	}).Error; err != nil {
		return err
	}

	err := query.First(website).Error

	return err
}

func (r *WebsiteRepository) FindByUserWebsites(db *gorm.DB, websites *[]entity.Website, request *model.WebsitesRequest) error {
	q := db.Model(&entity.Website{}).Where("user_id = ?", request.UserID)
	q = q.Where("name LIKE ?", "%"+request.Query+"%")
	q = q.Order(request.OrderBy)
	offset := (request.Page - 1) * request.PageSize
	q = q.Limit(request.PageSize).Offset(offset)

	err := q.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("username", "user_id")
	}).Find(websites).Error
	if err != nil {
		return err
	}

	return nil
}
