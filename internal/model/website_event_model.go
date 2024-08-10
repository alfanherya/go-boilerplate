package model

// Common fields shared by multiple request structs
type CommonRequestFields struct {
	WebsiteID string
	StartAt   int64  `query:"startAt" validate:"required"`
	EndAt     int64  `query:"endAt" validate:"required"`
	URL       string `query:"url"`
	Referrer  string `query:"referrer"`
	Title     string `query:"title"`
	OS        string `query:"os"`
	Browser   string `query:"browser"`
	Device    string `query:"device"`
	Country   string `query:"country"`
	Region    string `query:"region"`
	City      string `query:"city"`
}

type WEStatsReq struct {
	CommonRequestFields
	Query    string `query:"query"`
	Event    string `query:"event"`
	Language string `query:"language"`
}

type WEPageViewsReq struct {
	CommonRequestFields
	Unit     string `query:"unit" validate:"required"`
	Timezone string `query:"timezone" validate:"required"`
}

type WEMetricsReq struct {
	CommonRequestFields
	Type     string `query:"type" validate:"required"`
	Query    string `query:"query"`
	Event    string `query:"event"`
	Language string `query:"language"`
	Limit    int64  `query:"limit"`
	Offset   int64  `query:"offset"`
	Search   string `query:"search"`
}

type WECount struct {
	X int64 `json:"x"`
}

type WEStatsRes struct {
	PageViews ChangeValue `json:"pageViews"`
	Visitors  ChangeValue `json:"visitors"`
	Visits    ChangeValue `json:"visits"`
	Bounces   ChangeValue `json:"bounces"`
	TotalTime ChangeValue `json:"totalTime"`
}

// WEStats struct (for database mapping)
type WEStats struct {
	PageViews int64 `gorm:"column:pageviews"`
	Visitors  int64 `gorm:"column:visitors"`
	Visits    int64 `gorm:"column:visits"`
	Bounces   int64 `gorm:"column:bounces"`
	TotalTime int64 `gorm:"column:totaltime"`
}

type ChangeValue struct {
	Value  int64 `json:"value"`
	Change int64 `json:"change"`
}

// XY struct (for graph data)
type XY struct {
	X string `json:"x"`
	Y int64  `json:"y"`
}

type WEPageViewsRes struct {
	PageViews []XY `json:"pageViews"`
	Sessions  []XY `json:"sessions"`
}
