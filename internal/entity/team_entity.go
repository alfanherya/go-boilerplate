package entity

import (
	"time"

	"gorm.io/gorm"
)

type Team struct {
	ID         string         `gorm:"primaryKey;type:uuid;column:team_id"`
	Name       string         `gorm:"type:varchar(50)"`
	AccessCode *string        `gorm:"unique;type:varchar(50);column:access_code"`
	LogoURL    *string        `gorm:"type:varchar(2183);column:logo_url"`
	CreatedAt  *time.Time     `gorm:"type:timestamptz;column:created_at;default:now()"`
	UpdatedAt  *time.Time     `gorm:"type:timestamptz;column:updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"type:timestamptz;column:deleted_at"`
	Websites   []Website      `gorm:"foreignKey:TeamID"`
	TeamUsers  []TeamUser     `gorm:"foreignKey:TeamID"`
}

func (u *Team) TableName() string {
	return "team"
}
