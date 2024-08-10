package entity

import "time"

type EventData struct {
	ID             string       `gorm:"primaryKey;type:uuid;column:event_data_id"`
	WebsiteID      string       `gorm:"type:uuid;column:website_id"`
	WebsiteEventID string       `gorm:"type:uuid;column:website_event_id"`
	DataKey        string       `gorm:"type:varchar(500);column:data_key"`
	StringValue    *string      `gorm:"type:varchar(500);column:string_value"`
	NumberValue    *float64     `gorm:"type:decimal(19,4);column:number_value"`
	DateValue      *time.Time   `gorm:"type:timestamptz;column:date_value"`
	DataType       int          `gorm:"type:integer;column:data_type"`
	CreatedAt      *time.Time   `gorm:"type:timestamptz;column:created_at;default:now()"`
	Website        Website      `gorm:"foreignKey:WebsiteID"`
	WebsiteEvent   WebsiteEvent `gorm:"foreignKey:WebsiteEventID"`
}

func (u *EventData) TableName() string {
	return "event_data"
}
