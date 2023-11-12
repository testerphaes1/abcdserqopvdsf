package gateway

import "context"

type Gateway interface {
	CreateOrder(ctx context.Context, request CreateOrderRequest) (CreateOrderResponse, error)
	VerifyOrder(ctx context.Context, request VerifyOrderRequest) (VerifyOrderResponse, error)
}

type CreateOrderRequest struct {
	Amount         int    `json:"amount"`
	CallBackUrl    string `json:"call_back_url"`
	ServerUniqueId string `json:"server_unique_id"`
}

type CreateOrderResponse struct {
	GatewayOrderId string `json:"gateway_order_id"`
	OrderLink      string `json:"order_link"`
}

type VerifyOrderRequest struct {
	ServerOrderId  string `json:"server_order_id"`
	GatewayOrderId string `json:"gateway_order_id"`
}

type VerifyOrderResponse struct {
	Status         int    `json:"status"`
	TrackId        string `json:"track_id"`
	ServerOrderId  string `json:"server_order_id"`
	GatewayOrderId string `json:"gateway_order_id"`
	Amount         string `json:"amount"`
	GatewayTime    string `json:"gateway_time"`
}
