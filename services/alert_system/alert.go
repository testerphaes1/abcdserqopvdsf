package alert_system

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type AlertHandler interface {
	SendAlert(ctx context.Context, request AlertRequest) error
	AlertLogs(ctx context.Context, request AlertLogsRequest, page, perPage int) (AlertLogsResponse, error)
}

type alertHandler struct {
	AlertUrl string `json:"alert_url"`
}

type AlertLogsRequest struct {
	Id        *string `json:"id"`
	UserId    *string `json:"user_id"`
	AlertType *string `json:"alert_type"`
	Target    *string `json:"target"`
	Status    *string `json:"status"`
}

type AlertLogsResponse struct {
	Limit      int    `json:"limit"`
	Page       int    `json:"page"`
	Sort       string `json:"sort"`
	TotalRows  int    `json:"total_rows"`
	TotalPages int    `json:"total_pages"`
	Rows       []struct {
		Id           string    `json:"id"`
		UserId       string    `json:"user_id"`
		AlertType    string    `json:"alert_type"`
		Target       string    `json:"target"`
		Status       string    `json:"status"`
		Message      string    `json:"message"`
		InsertedTime time.Time `json:"inserted_time"`
	} `json:"rows"`
}

func NewAlertHandler(alertUrl string) AlertHandler {
	return &alertHandler{AlertUrl: alertUrl}
}

type AlertRequest struct {
	AlertType        string            `json:"alert_type"`
	UserId           string            `json:"user_id"`
	Targets          []string          `json:"targets"`
	Subject          string            `json:"subject"`
	Message          string            `json:"message"`
	IsTemplate       bool              `json:"is_template"`
	Template         string            `json:"template"`
	TemplateKeyPairs map[string]string `json:"template_key_pairs"`
	AdditionalData   map[string]string `json:"additional_data"`
}

func (a *alertHandler) SendAlert(ctx context.Context, request AlertRequest) error {
	reqB, _ := json.Marshal(request)
	req, err := http.NewRequestWithContext(ctx, "POST", a.AlertUrl+"/admin/alert", bytes.NewBuffer(reqB))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return errors.New("failed to send alert")
	}
	return nil
}

func (a *alertHandler) AlertLogs(ctx context.Context, request AlertLogsRequest, page, perPage int) (AlertLogsResponse, error) {
	reqB, _ := json.Marshal(request)
	req, err := http.NewRequestWithContext(ctx, "GET", a.AlertUrl+"/admin/alert/log"+fmt.Sprintf("?page=%d&&limit=%d", page, perPage), bytes.NewBuffer(reqB))
	if err != nil {
		return AlertLogsResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return AlertLogsResponse{}, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	var respM AlertLogsResponse
	err = json.Unmarshal(respBody, &respM)
	if err != nil {
		return AlertLogsResponse{}, err
	}
	return respM, nil
}
