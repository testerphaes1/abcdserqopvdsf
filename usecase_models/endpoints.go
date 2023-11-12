package usecase_models

type EndpointRules struct {
	EndpointName    string            `json:"endpoint_name"`
	Url             string            `json:"url"`
	Method          string            `json:"method"`
	Body            string            `json:"body"`
	Header          map[string]string `json:"header"`
	AcceptanceModel AcceptanceModel   `json:"acceptance_model"` // check keys with their type and status
}

type AcceptanceModel struct {
	Statuses       []string        `json:"statuses"`
	ResponseBodies []KeyValueModel `json:"response_bodies"`
}

type KeyValueModel struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type EndpointResponses struct {
	HeaderResponses map[string]map[string][]string `json:"header_responses"`
	BodyResponses   map[string]string              `json:"body_responses"`
	TimeResponses   map[string]float64             `json:"time_responses"`
	StatusResponses map[string]int                 `json:"status_responses"`
}
