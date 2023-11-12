package usecase_models

type Scheduling struct {
	PipelineId     int    `json:"pipeline_id"`
	PipelineName   string `json:"pipeline_name"`
	IsUp           bool   `json:"is_up"`
	ProjectId      int    `json:"project_id"`
	Duration       int    `json:"duration"`
	EndAt          string `json:"end_at"`
	IsActive       bool   `json:"is_active"`
	DataCentersIds []int  `json:"data_centers"` // datacenter id
	IsHeartBeat    bool   `json:"is_heart_beat"`
}
