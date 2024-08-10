package entity

import "time"

type Session struct {
	ID            string         `gorm:"primaryKey;type:uuid;column:session_id"`
	WebsiteID     string         `gorm:"type:uuid;column:website_id"`
	Hostname      *string        `gorm:"type:varchar(100)"`
	Browser       *string        `gorm:"type:varchar(20)"`
	OS            *string        `gorm:"type:varchar(20)"`
	Device        *string        `gorm:"type:varchar(20)"`
	Screen        *string        `gorm:"type:varchar(11)"`
	Language      *string        `gorm:"type:varchar(35)"`
	Country       *string        `gorm:"type:char(2)"`
	Subdivision1  *string        `gorm:"type:varchar(20)"`
	Subdivision2  *string        `gorm:"type:varchar(50)"`
	City          *string        `gorm:"type:varchar(50)"`
	CreatedAt     *time.Time     `gorm:"type:timestamptz;column:created_at;default:now()"`
	WebsiteEvents []WebsiteEvent `gorm:"foreignKey:SessionID"`
	SessionData   []SessionData  `gorm:"foreignKey:SessionID"`
}

func (u *Session) TableName() string {
	return "session"
}
