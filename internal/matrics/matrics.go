package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	TasksPending = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "tasks_pending_total",
			Help: "Tasks waiting to be processed",
		},
	)

	TasksRunning = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "tasks_running_total",
			Help: "Tasks currently running",
		},
	)

	TasksCompleted = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tasks_completed_total",
			Help: "Successfully completed tasks",
		},
	)

	TasksFailed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tasks_failed_total",
			Help: "Failed tasks",
		},
	)

	TaskDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "task_execution_duration_seconds",
			Help:    "Task execution duration",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"task_type"},
	)
)

func Register() {
	prometheus.MustRegister(
		TasksPending,
		TasksRunning,
		TasksCompleted,
		TasksFailed,
		TaskDuration,
	)
}
