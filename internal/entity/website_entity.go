package entity

import (
	"time"

	"gorm.io/gorm"
)

type Website struct {
	ID          string         `gorm:"primaryKey;type:uuid;column:website_id"`
	Name        string         `gorm:"type:varchar(100)"`
	Domain      *string        `gorm:"type:varchar(500)"`
	ShareID     *string        `gorm:"unique;type:varchar(50);column:share_id"`
	ResetAt     *time.Time     `gorm:"type:timestamptz;column:reset_at"`
	UserID      *string        `gorm:"type:uuid;column:user_id"`
	TeamID      *string        `gorm:"type:uuid;column:team_id"`
	CreatedBy   *string        `gorm:"type:uuid;column:created_by"`
	CreatedAt   *time.Time     `gorm:"type:timestamptz;column:created_at;default:now()"`
	UpdatedAt   *time.Time     `gorm:"type:timestamptz;column:updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"type:timestamptz;column:deleted_at"`
	User        User           `gorm:"foreignKey:UserID"`
	CreateUser  User           `gorm:"foreignKey:CreatedBy"`
	Team        Team           `gorm:"foreignKey:TeamID"`
	EventData   []EventData    `gorm:"foreignKey:WebsiteID"`
	Reports     []Report       `gorm:"foreignKey:WebsiteID"`
	SessionData []SessionData  `gorm:"foreignKey:WebsiteID"`
}

func (u *Website) TableName() string {
	return "website"
}
