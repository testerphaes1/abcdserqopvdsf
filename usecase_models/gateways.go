package usecase_models

import (
	"github.com/volatiletech/null/v8"
	"time"
)

type Gateway struct {
	ID             int         `json:"id"`
	Baseurl        string      `json:"baseurl"`
	Title          string      `json:"title"`
	ConnectionRate int         `json:"connection_rate"`
	IsActive       bool        `json:"is_active"`
	IsDefault      bool        `json:"is_default"`
	Data           interface{} `json:"data"`
	UpdatedAt      time.Time   `json:"updated_at"`
	CreatedAt      time.Time   `json:"created_at"`
	DeletedAt      null.Time   `json:"deleted_at"`
}

type CreateGatewayResponse struct {
	GatewayId int `json:"gateway_id"`
}
