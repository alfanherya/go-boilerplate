package entity

import "time"

type TeamUser struct {
	ID        string     `gorm:"primaryKey;type:uuid;column:team_user_id"`
	TeamID    string     `gorm:"type:uuid;column:team_id"`
	UserID    string     `gorm:"type:uuid;column:user_id"`
	Role      string     `gorm:"type:varchar(50)"`
	CreatedAt *time.Time `gorm:"type:timestamptz;column:created_at;default:now()"`
	UpdatedAt *time.Time `gorm:"type:timestamptz;column:updated_at"`
	Team      Team       `gorm:"foreignKey:TeamID"`
	User      User       `gorm:"foreignKey:UserID"`
}

func (u *TeamUser) TableName() string {
	return "team_user"
}
