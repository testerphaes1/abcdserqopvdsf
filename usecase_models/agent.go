package usecase_models

type AgentCurlRequest struct {
	Url    string              `json:"url"`
	Method string              `json:"method"`
	Header map[string][]string `json:"header"`
	Body   string              `json:"body"`
}

type AgentCurlResponse struct {
	Message    string `json:"message"`
	Status     int    `json:"status"`
	Statistics struct {
		ResponseTime float64             `json:"response_time"`
		StatusCode   int                 `json:"status_code"`
		Header       map[string][]string `json:"header"`
		Body         string              `json:"body"`
	} `json:"statistics"`
}

type AgentNetCatRequest struct {
	Address string `json:"address"`
	Port    string `json:"port"`
	Type    string `json:"type"`
	TimeOut int    `json:"time_out"`
}

type AgentNetCatResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type AgentPageSpeedRequest struct {
	Url string `json:"url"`
}

type AgentPageSpeedResponse struct {
	Status int `json:"status"`
}

type AgentPingRequest struct {
	Address string `json:"address"`
	Count   int    `json:"count"`
	TimeOut int    `json:"time_out"`
}

type AgentPingResponse struct {
	Status     int `json:"status"`
	Statistics struct {
		PacketsReceive int `json:"packets_receive"`
		PacketsSent    int `json:"packets_sent"`
		PacketLoss     int `json:"packet_loss"`
		IpAddress      struct {
			IP   string `json:"IP"`
			Zone string `json:"Zone"`
		} `json:"ip_address"`
		Address string `json:"address"`
		AvgRtt  int    `json:"avg_rtt"`
	} `json:"statistics"`
	Message string `json:"message"`
}

type AgentTraceRouteRequest struct {
	Address string `json:"address"`
	Retry   int    `json:"retry"`
	Hop     int    `json:"hop"`
}

type AgentTraceRouteResponse struct {
	Status  int      `json:"status"`
	Hop     []string `json:"hop"`
	Message string   `json:"message"`
}
