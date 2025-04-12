package metrics

type Config struct {
	ServiceName   string
	AgentHost     string
	AgentPort     string
	EnableMetrics bool
	OtlpEndpoint  string
}
