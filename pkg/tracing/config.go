package tracing

type Config struct {
	ServiceName   string
	AgentHost     string
	AgentPort     string
	EnableTracing bool
	OtlpEndpoint  string
}
