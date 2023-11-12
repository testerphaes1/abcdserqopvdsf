package tasks

import (
	"encoding/json"
	"github.com/hibiken/asynq"
	"log"
	"test-manager/repos"
	"test-manager/tasks/task_models"
)

func AggregateEndpointStats(group string, tasks []*asynq.Task) *asynq.Task {
	log.Printf("Aggregating %d tasks from group %q", len(tasks), group)
	var b = make([]repos.WriteEndpointStatsOptions, len(tasks))
	for i, t := range tasks {
		json.Unmarshal(t.Payload(), &b[i])
	}
	t, _ := json.Marshal(b)
	return asynq.NewTask(task_models.TypeEndpointStore, t)
}
