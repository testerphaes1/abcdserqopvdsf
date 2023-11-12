package task_models

type Type string

const (
	TypeEndpoint               = "end:point"
	TypeNotification           = "notification"
	TypeEndpointStore          = "endpoint_store"
	TypeAggregateEndpointStore = "aggregate_endpoint_store"
)

//TypeNetCats                = "net:cats"
//TypePageSpeeds             = "page:speeds"
//TypePings                  = "ping:s"
//TypeTraceRoutes            = "trace:routes"

const (
	QueueEndpoint      = "endpoint"
	QueueNotification  = "notification"
	QueueEndpointStore = "endpoint_store"
)

//QueueNetCats                = "net_cats"
//QueuePageSpeeds             = "page_speeds"
//QueuePings                  = "pings"
//QueueTraceRoutes            = "trace_routes"

const (
	GroupAggregateEndpointStore = "group_aggregate_endpoint_store"
)
