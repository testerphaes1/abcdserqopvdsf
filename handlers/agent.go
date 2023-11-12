package handlers

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"test-manager/usecase_models"
)

type AgentHandler interface {
	SendCurl(ctx context.Context, dataCenterUrl string, request usecase_models.AgentCurlRequest) (response string, responseHeader map[string][]string, status int, responseTime float64, err error)
	SendNetCat(ctx context.Context, dataCenterUrl string, request usecase_models.AgentNetCatRequest) (response usecase_models.AgentNetCatResponse, err error)
	SendPageSpeed(ctx context.Context, dataCenterUrl string, request usecase_models.AgentPageSpeedRequest) (response usecase_models.AgentPageSpeedResponse, err error)
	SendPing(ctx context.Context, dataCenterUrl string, request usecase_models.AgentPingRequest) (response usecase_models.AgentPingResponse, err error)
	SendTraceRoute(ctx context.Context, dataCenterUrl string, request usecase_models.AgentTraceRouteRequest) (response usecase_models.AgentTraceRouteResponse, err error)
}

type agentHandler struct {
}

func NewAgentHandler() AgentHandler {
	return &agentHandler{}
}

func (a *agentHandler) SendCurl(ctx context.Context, dataCenterUrl string, request usecase_models.AgentCurlRequest) (response string, responseHeader map[string][]string, status int, responseTime float64, err error) {
	reqB, _ := json.Marshal(request)
	req, err := http.NewRequestWithContext(ctx, "POST", dataCenterUrl+"/v1/curl", bytes.NewBuffer(reqB))
	if err != nil {
		return response, responseHeader, status, responseTime, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Ignore invalid certificates
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return response, responseHeader, status, responseTime, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	var respM usecase_models.AgentCurlResponse
	err = json.Unmarshal(respBody, &respM)
	if err != nil {
		return response, responseHeader, status, responseTime, err
	}

	if respM.Statistics.StatusCode == 0 {
		return respM.Message, respM.Statistics.Header, respM.Statistics.StatusCode, respM.Statistics.ResponseTime, nil
	}
	return respM.Statistics.Body, respM.Statistics.Header, respM.Statistics.StatusCode, respM.Statistics.ResponseTime, nil
}

func (a *agentHandler) SendNetCat(ctx context.Context, dataCenterUrl string, request usecase_models.AgentNetCatRequest) (response usecase_models.AgentNetCatResponse, err error) {
	reqB, _ := json.Marshal(request)
	req, err := http.NewRequestWithContext(ctx, "POST", dataCenterUrl+"/v1/netcat", bytes.NewBuffer(reqB))
	if err != nil {
		return response, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	var respM usecase_models.AgentNetCatResponse
	err = json.Unmarshal(respBody, &respM)
	if err != nil {
		return response, err
	}
	return respM, nil
}

func (a *agentHandler) SendPageSpeed(ctx context.Context, dataCenterUrl string, request usecase_models.AgentPageSpeedRequest) (response usecase_models.AgentPageSpeedResponse, err error) {
	reqB, _ := json.Marshal(request)
	req, err := http.NewRequestWithContext(ctx, "POST", dataCenterUrl+"/v1/pagespeed", bytes.NewBuffer(reqB))
	if err != nil {
		return response, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	var respM usecase_models.AgentPageSpeedResponse
	err = json.Unmarshal(respBody, &respM)
	if err != nil {
		return response, err
	}
	return respM, nil
}

func (a *agentHandler) SendPing(ctx context.Context, dataCenterUrl string, request usecase_models.AgentPingRequest) (response usecase_models.AgentPingResponse, err error) {
	reqB, _ := json.Marshal(request)
	req, err := http.NewRequestWithContext(ctx, "POST", dataCenterUrl+"/v1/ping", bytes.NewBuffer(reqB))
	if err != nil {
		return response, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	var respM usecase_models.AgentPingResponse
	err = json.Unmarshal(respBody, &respM)
	if err != nil {
		return response, err
	}
	return respM, nil
}

func (a *agentHandler) SendTraceRoute(ctx context.Context, dataCenterUrl string, request usecase_models.AgentTraceRouteRequest) (response usecase_models.AgentTraceRouteResponse, err error) {
	reqB, _ := json.Marshal(request)
	req, err := http.NewRequestWithContext(ctx, "POST", dataCenterUrl+"/v1/traceroute", bytes.NewBuffer(reqB))
	if err != nil {
		return response, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	var respM usecase_models.AgentTraceRouteResponse
	err = json.Unmarshal(respBody, &respM)
	if err != nil {
		return response, err
	}
	return respM, nil
}
