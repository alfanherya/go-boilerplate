package repository

import (
	"time"
	"umami-go/internal/entity"
	"umami-go/internal/libs"
	"umami-go/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WebsiteEventRepository struct {
	Repository[entity.WebsiteEvent]
	Log *logrus.Logger
}

func NewWebsiteEventRepository(log *logrus.Logger) *WebsiteEventRepository {
	return &WebsiteEventRepository{
		Log: log,
	}
}

func (r *WebsiteEventRepository) Active(db *gorm.DB, websiteEvent *entity.WebsiteEvent, websiteID string) (int64, error) {
	var count int64
	startAt := time.Now().Add(-1 * time.Minute) // 1 minutes

	if err := db.Model(websiteEvent).Distinct("session_id").Where("website_id = ?", websiteID).Where("created_at >= ?", startAt).Count(&count).Error; err != nil {
		r.Log.Warnf("Failed to count website event : %+v", err)
		return 0, err
	}

	return count, nil
}

func (r *WebsiteEventRepository) Stats(db *gorm.DB, stats *model.WEStats, req *model.WEStatsReq) error {
	subQuery := db.Table("website_event").
		Select("website_event.session_id, website_event.visit_id, COUNT(*) AS count, MIN(website_event.created_at) AS min_time, MAX(website_event.created_at) AS max_time").
		Joins("INNER JOIN session ON website_event.session_id = session.session_id").
		Where("website_event.website_id = ?", req.WebsiteID).
		Where("EXTRACT(EPOCH FROM website_event.created_at) * 1000 BETWEEN ? AND ?", req.StartAt, req.EndAt).
		Where("website_event.event_type = 1").
		Group("website_event.session_id, website_event.visit_id")

	for key, value := range libs.FilterColumns {
		if fieldValue := libs.GetFilterColumns(req, key); fieldValue != "" {
			subQuery = subQuery.Where(value+" = ?", fieldValue)
		}
	}

	return db.Raw(`
		SELECT SUM(t.count) AS pageviews,
			   COUNT(DISTINCT t.session_id) AS visitors,
			   COUNT(DISTINCT t.visit_id) AS visits,
			   SUM(CASE WHEN t.count = 1 THEN 1 ELSE 0 END) AS bounces,
			   SUM(FLOOR(EXTRACT(EPOCH FROM (t.min_time - t.max_time)))) AS totaltime
		FROM (?) as t
	`, subQuery).Scan(&stats).Error
}

func (r *WebsiteEventRepository) PageViews(db *gorm.DB, pageStats, sessionStats *[]model.XY, req *model.WEPageViewsReq) error {
	page := db.Table("website_event").
		Select(`TO_CHAR(DATE_TRUNC('DAY', WEBSITE_EVENT.CREATED_AT at TIME zone 'ASIA/JAKARTA'), 'YYYY-MM-DD') AS x,
                COUNT(*) AS y`).
		Joins("INNER JOIN session ON WEBSITE_EVENT.SESSION_ID = SESSION.SESSION_ID").
		Where("WEBSITE_EVENT.WEBSITE_ID = ?", req.WebsiteID).
		Where("EXTRACT(EPOCH FROM WEBSITE_EVENT.CREATED_AT) * 1000 BETWEEN ? AND ?", req.StartAt, req.EndAt).
		Where("EVENT_TYPE = 1").
		Group("x").
		Scan(&pageStats).Error

	if page != nil {
		return page
	}

	session := db.Table("website_event").
		Select(`TO_CHAR(DATE_TRUNC('DAY', WEBSITE_EVENT.CREATED_AT at TIME zone 'ASIA/JAKARTA'), 'YYYY-MM-DD') AS x,
                COUNT(distinct website_event.session_id) AS y`).
		Joins("INNER JOIN session ON WEBSITE_EVENT.SESSION_ID = SESSION.SESSION_ID").
		Where("WEBSITE_EVENT.WEBSITE_ID = ?", req.WebsiteID).
		Where("EXTRACT(EPOCH FROM WEBSITE_EVENT.CREATED_AT) * 1000 BETWEEN ? AND ?", req.StartAt, req.EndAt).
		Where("EVENT_TYPE = 1").
		Group("x").
		Scan(&sessionStats).Error

	if session != nil {
		return session
	}

	return nil
}

func (r *WebsiteEventRepository) Metrics(db *gorm.DB, pageStats *[]model.XY, req *model.WEMetricsReq) error {
	q := db.Table("website_event").
		Select(`language x, count(distinct website_event.session_id) y`).
		Joins("inner join session on website_event.session_id = session.session_id").
		Where("website_event.website_id = ?", req.WebsiteID).
		Where("website_event.event_type = ?", 1).
		Where("extract(epoch from website_event.created_at) * 1000 between ? and ?", req.StartAt, req.EndAt).
		Group("x").
		Order("y desc")

	if req.Limit > 0 {
		q = q.Limit(int(req.Limit))
	} else {
		q = q.Limit(500)
	}

	if req.Offset > 0 {
		q = q.Offset(int(req.Offset))
	}

	return q.Scan(&pageStats).Error
}
