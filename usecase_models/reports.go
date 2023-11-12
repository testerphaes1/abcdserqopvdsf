package usecase_models

import "time"

type EndpointReportResult struct {
	Time  time.Time   `json:"_time"`
	Field string      `json:"_field"`
	Value interface{} `json:"_value"`
}
