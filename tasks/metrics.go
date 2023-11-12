package tasks

import (
	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
	"test-manager/monitoring"
	"test-manager/utils"
	"time"
)

// ReportMetrics report prometheus metrics
func ReportMetrics(refreshInterrval time.Duration, redisConf *redis.Options) {
	inspector := utils.AsynqInspector(redisConf)

	go func() {
		var queues []string
		var err error
		for {
			time.Sleep(refreshInterrval)

			// Fetch only active queues list, we get error in GetQueueInfo if queue was not active
			// So we can't just use constant list of queues
			queues, err = inspector.Queues()
			if err != nil {
				log.Errorf("cant get list of queues: %s", err)
			}

			for _, queue := range queues {
				stats, err := inspector.GetQueueInfo(queue)
				if err != nil {
					log.Errorf("cant get current stat for queue: %s : %s", queue, err)
					continue
				}

				monitoring.LatencyTasksGauge.WithLabelValues(queue).Set(stats.Latency.Seconds())
				monitoring.SizeTasksGauge.WithLabelValues(queue).Set(float64(stats.Size))
				monitoring.CompletedTasksGauge.WithLabelValues(queue).Set(float64(stats.Completed))
				monitoring.ProcessedTotalTasksGauge.WithLabelValues(queue).Set(float64(stats.ProcessedTotal))
				monitoring.FailedTotalTasksGauge.WithLabelValues(queue).Set(float64(stats.FailedTotal))
				monitoring.ScheduledTasksGauge.WithLabelValues(queue).Set(float64(stats.Scheduled))
				monitoring.PendingTasksGauge.WithLabelValues(queue).Set(float64(stats.Pending))
				monitoring.ActiveTasksGauge.WithLabelValues(queue).Set(float64(stats.Active))
				monitoring.FailedTasksGauge.WithLabelValues(queue).Set(float64(stats.Failed))
				monitoring.RetriedTasksGauge.WithLabelValues(queue).Set(float64(stats.Retry))
				monitoring.ArchivedTasksGauge.WithLabelValues(queue).Set(float64(stats.Archived))
				monitoring.ProcessedTasksGauge.WithLabelValues(queue).Set(float64(stats.Processed))
				monitoring.QueueMemoryUsageGauge.WithLabelValues(queue).Set(float64(stats.MemoryUsage))
			}
		}
	}()
}
