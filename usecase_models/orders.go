package usecase_models

import (
	"github.com/volatiletech/null/v8"
	"time"
)

const (
	OrderStatusCreated  = 0
	OrderStatusPending  = 1
	OrderStatusVerified = 2
	OrderStatusDone     = 2
	OrderStatusCanceled = 3
)

type Order struct {
	ID        int       `json:"id"`
	AccountId int       `json:"account_id"`
	ProjectId int       `json:"project_id"`
	PackageId int       `json:"package_id"`
	GatewayId int       `json:"gateway_id"`
	Status    string    `json:"status"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at"`
}

type VerifyOrder struct {
	GatewayOrderId string `json:"gateway_order_id"`
}

type CreateOrderResponse struct {
	OrderId     int    `json:"order_id"`
	GatewayLink string `json:"gateway_link"`
}
