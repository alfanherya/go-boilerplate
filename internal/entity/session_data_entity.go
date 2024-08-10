package entity

import "time"

type SessionData struct {
	ID          string     `gorm:"primaryKey;type:uuid;column:session_data_id"`
	WebsiteID   string     `gorm:"type:uuid;column:website_id"`
	SessionID   string     `gorm:"type:uuid;column:session_id"`
	DataKey     string     `gorm:"type:varchar(500);column:data_key"`
	StringValue *string    `gorm:"type:varchar(500);column:string_value"`
	NumberValue *float64   `gorm:"type:decimal(19,4);column:number_value"`
	DateValue   *time.Time `gorm:"type:timestamptz;column:date_value"`
	DataType    int        `gorm:"type:integer;column:data_type"`
	CreatedAt   *time.Time `gorm:"type:timestamptz;column:created_at;default:now()"`
	Website     Website    `gorm:"foreignKey:WebsiteID"`
	Session     Session    `gorm:"foreignKey:SessionID"`
}

func (u *SessionData) TableName() string {
	return "session_data"
}
