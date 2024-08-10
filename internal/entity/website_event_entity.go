package entity

import "time"

type WebsiteEvent struct {
	ID             string      `gorm:"primaryKey;type:uuid;column:event_id"`
	WebsiteID      string      `gorm:"type:uuid;column:website_id"`
	SessionID      string      `gorm:"type:uuid;column:session_id"`
	VisitID        string      `gorm:"type:uuid;column:visit_id"`
	CreatedAt      *time.Time  `gorm:"type:timestamptz;column:created_at;default:now()"`
	URLPath        string      `gorm:"type:varchar(500);column:url_path"`
	URLQuery       *string     `gorm:"type:varchar(500);column:url_query"`
	ReferrerPath   *string     `gorm:"type:varchar(500);column:referrer_path"`
	ReferrerQuery  *string     `gorm:"type:varchar(500);column:referrer_query"`
	ReferrerDomain *string     `gorm:"type:varchar(500);column:referrer_domain"`
	PageTitle      *string     `gorm:"type:varchar(500);column:page_title"`
	EventType      int         `gorm:"type:integer;column:event_type;default:1"`
	EventName      *string     `gorm:"type:varchar(50);column:event_name"`
	EventData      []EventData `gorm:"foreignKey:WebsiteEventID"`
	Session        Session     `gorm:"foreignKey:SessionID"`
}

func (u *WebsiteEvent) TableName() string {
	return "website_event"
}
