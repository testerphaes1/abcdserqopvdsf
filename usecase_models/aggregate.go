package usecase_models

type AggregateAllRuleSubWorks struct {
	Endpoints   []*Endpoints   `json:"endpoints"`
	TraceRoutes []*TraceRoutes `json:"trace_routes"`
	NetCats     []*NetCats     `json:"net_cats"`
	Pings       []*Pings       `json:"pings"`
	PageSpeed   []*PageSpeeds  `json:"page_speeds"`
}
