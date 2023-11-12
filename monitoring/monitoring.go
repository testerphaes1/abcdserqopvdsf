package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	EndpointStoreDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "endpoint_store_duration",
			Help:    "Endpoint store duration to complete.",
			Buckets: []float64{0.3, 0.5, 1, 10},
		},
	)

	EndpointTaskDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "endpoint_tasks_duration",
			Help:    "Endpoint tasks duration to complete.",
			Buckets: []float64{60, 300},
		},
	)

	EndpointTaskDatacenterDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "endpoint_tasks_datacenter_duration",
			Help:    "Endpoint tasks datacenter duration to complete.",
			Buckets: []float64{0.3, 0.5, 1, 10},
		},
		[]string{"datacenter"},
	)

	NotificationTaskCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "notification_tasks_counter",
			Help: "Notification tasks counter to complete.",
		},
		[]string{"state"},
	)

	ActiveSchedulerTasksGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_scheduler_tasks_total",
			Help: "The total number of scheduler tasks",
		},
	)

	ErrorEndpointTaskFetchingDatacenterGauge = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "error_endpoint_task_fetching_datacenter",
			Help: "there is a problem on fetching data center on endpoint task",
		},
	)

	DisabledEndpointRuleReasonFunctionality = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "disabled_endpoint_rule_reason_functionality",
			Help: "endpoint rule is disabled because of functionality problem",
		},
	)

	ProcessedTasksGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "processed_tasks_total",
			Help: "The total number of processed tasks",
		},
		[]string{"queue_name"},
	)

	FailedTasksGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "failed_tasks_total",
			Help: "The total number of times processing failed",
		},
		[]string{"queue_name"},
	)

	ActiveTasksGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "active_tasks_total",
			Help: "The number of tasks currently being processed",
		},
		[]string{"queue_name"},
	)

	PendingTasksGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "pending_tasks_total",
			Help: "The number of tasks ready to be processed",
		},
		[]string{"queue_name"},
	)

	ScheduledTasksGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "scheduled_tasks_total",
			Help: "The number of tasks scheduled for future",
		},
		[]string{"queue_name"},
	)

	ArchivedTasksGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "archived_tasks_total",
			Help: "The number of tasks that reached max-retry",
		},
		[]string{"queue_name"},
	)

	RetriedTasksGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "retried_tasks_total",
			Help: "The number of tasks that retried",
		},
		[]string{"queue_name"},
	)

	QueueMemoryUsageGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "queue_total_memory_bytes",
			Help: "Total number of bytes that the queue and its tasks require to be stored in redis",
		},
		[]string{"queue_name"},
	)

	LatencyTasksGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "queue_latency",
			Help: "Latency of the queue, measured by the oldest pending task in the queue.",
		},
		[]string{"queue_name"},
	)
	SizeTasksGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "queue_total_tasks",
			Help: "Size is the total number of tasks in the queue.",
		},
		[]string{"queue_name"},
	)
	CompletedTasksGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "total_completed_tasks",
			Help: "Number of stored completed tasks.",
		},
		[]string{"queue_name"},
	)
	ProcessedTotalTasksGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "total_processed_total_tasks",
			Help: "Total number of tasks processed (cumulative).",
		},
		[]string{"queue_name"},
	)
	FailedTotalTasksGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "total_failed_total_tasks",
			Help: "Total number of tasks failed to be processed within the given date (counter resets daily).",
		},
		[]string{"queue_name"},
	)
)

// Init initialize monitoring
func Init() {
}
