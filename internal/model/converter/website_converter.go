package converter

import (
	"time"
	"umami-go/internal/entity"
	"umami-go/internal/model"
)

func WebsitesToResponse(websites *[]entity.Website, query *model.WebsitesRequest) *model.WebsitesResponse {
	var websiteResponses []model.WebsiteResponse

	for _, website := range *websites {
		websiteResponses = append(websiteResponses, *WebsiteToResponse(&website))
	}

	return &model.WebsitesResponse{
		Data:     websiteResponses,
		Count:    len(*websites),
		OrderBy:  query.OrderBy,
		Page:     query.Page,
		PageSize: query.PageSize,
	}
}

func WebsiteToResponse(website *entity.Website) *model.WebsiteResponse {
	var domain, shareID string
	var resetAt, createdAt, updatedAt, deletedAt *time.Time

	if website.Domain != nil {
		domain = *website.Domain
	}

	if website.ShareID != nil {
		shareID = *website.ShareID
	}

	if website.ResetAt != nil {
		resetAt = website.ResetAt
	}

	if website.CreatedAt != nil {
		createdAt = website.CreatedAt
	}

	if website.UpdatedAt != nil {
		updatedAt = website.UpdatedAt
	}

	if website.DeletedAt.Valid {
		deletedAt = &website.DeletedAt.Time
	}

	var userID, teamID, createdBy string
	if website.UserID != nil {
		userID = *website.UserID
	}

	if website.TeamID != nil {
		teamID = *website.TeamID
	}

	if website.CreatedBy != nil {
		createdBy = *website.CreatedBy
	}

	return &model.WebsiteResponse{
		ID:        website.ID,
		Name:      website.Name,
		Domain:    domain,
		ShareID:   shareID,
		ResetAt:   resetAt,
		UserID:    userID,
		TeamID:    teamID,
		CreatedBy: createdBy,
		CreateAt:  createdAt,
		UpdateAt:  updatedAt,
		DeleteAt:  deletedAt,
	}
}
