package gateway

import "context"

type zarinpal struct {
	Baseurl  string
	ApiToken string
}

func NewZarinpalGateway(baseurl string, apiToken string) Gateway {
	return &zarinpal{
		Baseurl:  baseurl,
		ApiToken: apiToken,
	}
}

func (i zarinpal) CreateOrder(ctx context.Context, request CreateOrderRequest) (CreateOrderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (i zarinpal) VerifyOrder(ctx context.Context, request VerifyOrderRequest) (VerifyOrderResponse, error) {
	//TODO implement me
	panic("implement me")
}
