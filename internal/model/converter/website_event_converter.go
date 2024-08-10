package converter

import (
	"fmt"
	"strings"
	"umami-go/internal/model"
)

func WEStatsToResponse(statsNow, statsPrev *model.WEStats) *model.WEStatsRes {
	return &model.WEStatsRes{
		PageViews: model.ChangeValue{
			Value:  statsNow.PageViews,
			Change: statsNow.PageViews - statsPrev.PageViews,
		},
		Visitors: model.ChangeValue{
			Value:  statsNow.Visitors,
			Change: statsNow.Visitors - statsPrev.Visitors,
		},
		Visits: model.ChangeValue{
			Value:  statsNow.Visits,
			Change: statsNow.Visits - statsPrev.Visits,
		},
		Bounces: model.ChangeValue{
			Value:  statsNow.Bounces,
			Change: statsNow.Bounces - statsPrev.Bounces,
		},
		TotalTime: model.ChangeValue{
			Value:  statsNow.TotalTime,
			Change: statsNow.TotalTime - statsPrev.TotalTime,
		},
	}
}

func WEPageViewsToResponse(pageStats, sessionStats *[]model.XY) *model.WEPageViewsRes {
	return &model.WEPageViewsRes{
		PageViews: *pageStats,
		Sessions:  *sessionStats,
	}
}

func WEMetricsLangToResponse(data *[]model.XY) *[]model.XY {
	if data == nil {
		return nil
	}

	combined := make(map[string]model.XY)
	for _, d := range *data {
		key := strings.ToLower(strings.Split(fmt.Sprint(d.X), "-")[0])
		if val, ok := combined[key]; !ok {
			combined[key] = model.XY{X: key, Y: d.Y}
		} else {
			val.Y += d.Y
			combined[key] = val
		}
	}

	combinedValues := make([]model.XY, 0, len(combined))
	for _, v := range combined {
		combinedValues = append(combinedValues, v)
	}

	return &combinedValues
}
