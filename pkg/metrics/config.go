package metrics

import "time"

type Config struct {
	ServiceName    string
	EnableMetrics  bool
	OtlpEndpoint   string
	ExportInterval time.Duration
	ExportTimeout  time.Duration
}
