package entity

import "time"

type Report struct {
	ID          string     `gorm:"primaryKey;type:uuid;column:report_id"`
	UserID      string     `gorm:"type:uuid;column:user_id"`
	WebsiteID   string     `gorm:"type:uuid;column:website_id"`
	Type        string     `gorm:"type:varchar(200)"`
	Name        string     `gorm:"type:varchar(200)"`
	Description string     `gorm:"type:varchar(500)"`
	Parameters  string     `gorm:"type:varchar(6000)"`
	CreatedAt   *time.Time `gorm:"type:timestamptz;column:created_at;default:now()"`
	UpdatedAt   *time.Time `gorm:"type:timestamptz;column:updated_at"`
	User        User       `gorm:"foreignKey:UserID"`
	Website     Website    `gorm:"foreignKey:WebsiteID"`
}

func (u *Report) TableName() string {
	return "report"
}
