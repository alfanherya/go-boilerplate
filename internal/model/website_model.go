package model

import "time"

type WebsitesRequest struct {
	Query    string `query:"query"`
	Page     int    `query:"page"`
	PageSize int    `query:"pageSize"`
	OrderBy  string `query:"orderBy"`
	UserID   string
}

type WebsitesResponse struct {
	Data     []WebsiteResponse `json:"data"`
	Count    int               `json:"count"`
	Page     int               `json:"page"`
	PageSize int               `json:"pageSize"`
	OrderBy  string            `json:"orderBy"`
}

type WebsiteResponse struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Domain    string     `json:"domain"`
	ShareID   string     `json:"shareId"`
	ResetAt   *time.Time `json:"resetAt"`
	UserID    string     `json:"userId"`
	TeamID    string     `json:"teamId"`
	CreatedBy string     `json:"createdBy"`
	CreateAt  *time.Time `json:"createdAt"`
	UpdateAt  *time.Time `json:"updatedAt"`
	DeleteAt  *time.Time `json:"deletedAt"`
}

type WebsiteCreateRequest struct {
	Name    string `json:"name" validate:"required,max=100"`
	Domain  string `json:"domain" validate:"required,max=500"`
	ShareID string `json:"shareId"`
	TeamID  string `json:"teamId"`
	UserID  string `json:"userId"`
}

type WebsiteUpdateRequest struct {
	ID      string `json:"id" validate:"required"`
	Name    string `json:"name" validate:"required,max=100"`
	Domain  string `json:"domain" validate:"required,max=500"`
	ShareID string `json:"shareId"`
}
