package usecase_models

type RulesRequest struct {
	Endpoints   Endpoints   `json:"endpoints"`
	TraceRoutes TraceRoutes `json:"trace_routes"`
	NetCats     NetCats     `json:"net_cats"`
	Pings       Pings       `json:"pings"`
	PageSpeed   PageSpeeds  `json:"page_speeds"`
}

type Endpoints struct {
	Endpoints  []EndpointRules `json:"endpoints"`
	Scheduling Scheduling      `json:"scheduling"`
}

type TraceRoutes struct {
	TraceRouts []TraceRouteRules `json:"trace_routes"`
	Scheduling Scheduling        `json:"scheduling"`
}

type NetCats struct {
	NetCats    []NetCatsRules `json:"net_cats"`
	Scheduling Scheduling     `json:"scheduling"`
}

type Pings struct {
	Pings      []PingsRules `json:"pings"`
	Scheduling Scheduling   `json:"scheduling"`
}

type PageSpeeds struct {
	PageSpeed  []PageSpeedRules `json:"page_speeds"`
	Scheduling Scheduling       `json:"scheduling"`
}
