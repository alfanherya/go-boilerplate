package libs

import "umami-go/internal/model"

var FilterColumns = map[string]string{
	"url":      "url_path",
	"referrer": "referrer_domain",
	"title":    "page_title",
	"query":    "url_query",
	"os":       "os",
	"browser":  "browser",
	"device":   "device",
	"country":  "country",
	"region":   "subdivision1",
	"city":     "city",
	"language": "language",
	"event":    "event_name",
}

func GetFilterColumns(query *model.WEStatsReq, field string) string {
	switch field {
	case "url":
		return query.URL
	case "referrer":
		return query.Referrer
	case "title":
		return query.Title
	case "query":
		return query.Query
	case "os":
		return query.OS
	case "browser":
		return query.Browser
	case "device":
		return query.Device
	case "country":
		return query.Country
	case "region":
		return query.Region
	case "city":
		return query.City
	case "language":
		return query.Language
	case "event":
		return query.Event
	default:
		return ""
	}
}
