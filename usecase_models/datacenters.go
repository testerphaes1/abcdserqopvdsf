package usecase_models

import (
	"github.com/volatiletech/null/v8"
	"time"
)

type Datacenter struct {
	ID             int          `json:"id"`
	Baseurl        string       `json:"baseurl"`
	Title          string       `json:"title"`
	ConnectionRate null.Int     `json:"connection_rate"`
	Lat            null.Float64 `json:"lat"`
	LNG            null.Float64 `json:"lng"`
	LocationName   null.String  `json:"location_name"`
	CountryName    null.String  `json:"country_name"`
	UpdatedAt      time.Time    `json:"updated_at"`
	CreatedAt      time.Time    `json:"created_at"`
	DeletedAt      null.Time    `json:"deleted_at"`
}

type CreateDatacenterResponse struct {
	DatacenterId int `json:"datacenter_id"`
}
