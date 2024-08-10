package model

type HealthResponse struct {
	Database CheckDBResponse    `json:"database"`
	Redis    CheckRedisResponse `json:"redis"`
}

type CheckDBResponse struct {
	Idle              string `json:"idle"`
	InUse             string `json:"in_use"`
	MaxIdleClosed     string `json:"max_idle_closed"`
	MaxLifetimeClosed string `json:"max_lifetime_closed"`
	Message           string `json:"message"`
	OpenConnections   string `json:"open_connections"`
	Status            string `json:"status"`
	WaitCount         string `json:"wait_count"`
	WaitDuration      string `json:"wait_duration"`
}

type CheckRedisResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Info    map[string]interface{} `json:"info"`
}
