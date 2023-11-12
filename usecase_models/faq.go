package usecase_models

import (
	"github.com/volatiletech/null/v8"
	"time"
)

type Faq struct {
	ID        int         `json:"id"`
	Question  null.String `json:"question"`
	Answer    null.String `json:"answer"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	DeletedAt null.Time   `json:"deleted_at"`
}

type CreateFaqResponse struct {
	FaqId int `json:"faq_id"`
}
