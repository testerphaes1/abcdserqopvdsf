package task_models

type NotificationsPayload struct {
	Type                string `json:"type"`
	State               string `json:"state"`
	ProjectId           int    `json:"project_id"`
	Username            string `json:"username"`
	Address             string `json:"address"`
	PipelineName        string `json:"pipeline_name"`
	Datacenters         string `json:"datacenters"`
	RootCause           string `json:"root_cause"`
	Time                string `json:"time"`
	ResolvedDatacenters string `json:"resolved_datacenters"`
	FailedDatacenters   string `json:"failed_datacenters"`
}
