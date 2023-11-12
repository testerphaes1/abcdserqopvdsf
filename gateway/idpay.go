package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type IdpayData struct {
	CallBackUrl string `json:"call_back_url"`
}

type idpay struct {
	Baseurl string `json:"baseurl"`
	ApiKey  string `json:"api_key"`
}

func NewIdpayGateway(baseurl string, apiKey string) Gateway {
	return &idpay{
		Baseurl: baseurl,
		ApiKey:  apiKey,
	}
}

func (i idpay) CreateOrder(ctx context.Context, request CreateOrderRequest) (CreateOrderResponse, error) {
	url := i.Baseurl + "/v1.1/payment"
	reqB, _ := json.Marshal(idpayCreateOrderRequest{
		OrderId:     request.ServerUniqueId,
		CallBackUrl: request.CallBackUrl,
		Amount:      request.Amount,
	})
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqB))
	if err != nil {
		return CreateOrderResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", i.ApiKey)
	req.Header.Set("X-SANDBOX", "1")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return CreateOrderResponse{}, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var respM idpayCreateOrderResponse
	err = json.Unmarshal(respBody, &respM)
	if err != nil {
		return CreateOrderResponse{}, err
	}
	return CreateOrderResponse{
		GatewayOrderId: respM.Id,
		OrderLink:      respM.Link,
	}, nil
}

func (i idpay) VerifyOrder(ctx context.Context, request VerifyOrderRequest) (VerifyOrderResponse, error) {
	url := i.Baseurl + "/v1.1/payment/verify"
	reqB, _ := json.Marshal(idpayVerifyOrderRequest{
		Id:      request.GatewayOrderId,
		OrderId: request.ServerOrderId,
	})
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqB))
	if err != nil {
		return VerifyOrderResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", i.ApiKey)
	req.Header.Set("X-SANDBOX", "1")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return VerifyOrderResponse{}, err
	}
	if resp.StatusCode != 200 {
		return VerifyOrderResponse{}, errors.New("problem in calling idpay gateway")
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var respM idpayVerifyOrderResponse
	err = json.Unmarshal(respBody, &respM)
	if err != nil {
		return VerifyOrderResponse{}, err
	}
	return VerifyOrderResponse{
		Status:         respM.Status,
		TrackId:        respM.TrackId,
		ServerOrderId:  respM.OrderId,
		GatewayOrderId: respM.Id,
		Amount:         respM.Amount,
		GatewayTime:    respM.Date,
	}, nil
}

type idpayCreateOrderRequest struct {
	OrderId     string `json:"order_id"`
	CallBackUrl string `json:"callback"`
	Amount      int    `json:"amount"`
}

type idpayCreateOrderResponse struct {
	Id   string `json:"id"`
	Link string `json:"link"`
}

type idpayVerifyOrderRequest struct {
	Id      string `json:"id"`
	OrderId string `json:"order_id"`
}
type idpayVerifyOrderResponse struct {
	Status  int    `json:"status"`
	TrackId string `json:"track_id"`
	Id      string `json:"id"`
	OrderId string `json:"order_id"`
	Amount  string `json:"amount"`
	Date    string `json:"date"`
	Payment struct {
		TrackId      string `json:"track_id"`
		Amount       string `json:"amount"`
		CardNo       string `json:"card_no"`
		HashedCardNo string `json:"hashed_card_no"`
		Date         string `json:"date"`
	} `json:"payment"`
	Verify struct {
		Date interface{} `json:"date"`
	} `json:"verify"`
}
