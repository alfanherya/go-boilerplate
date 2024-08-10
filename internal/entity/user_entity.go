package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          string         `gorm:"primaryKey;type:uuid;column:user_id"`
	Username    string         `gorm:"unique;type:varchar(255)"`
	Password    string         `gorm:"type:varchar(60)"`
	Role        string         `gorm:"type:varchar(50)"`
	LogoURL     *string        `gorm:"type:varchar(2183);column:logo_url"`
	DisplayName *string        `gorm:"type:varchar(255);column:display_name"`
	CreatedAt   *time.Time     `gorm:"type:timestamptz;column:created_at;default:now()"`
	UpdatedAt   *time.Time     `gorm:"type:timestamptz;column:updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"type:timestamptz;column:deleted_at"`
	Websites    []Website      `gorm:"foreignKey:UserID"`
	Teams       []TeamUser     `gorm:"foreignKey:UserID"`
	Reports     []Report       `gorm:"foreignKey:UserID"`
}

func (u *User) TableName() string {
	return "user"
}
