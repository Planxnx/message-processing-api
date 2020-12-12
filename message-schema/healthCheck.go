package messageschema

type HealthCheckMessageFormat struct {
	Feature     string        `json:"feature"`
	Description string        `json:"description"`
	ExecuteMode []ExecuteMode `json:"executeMode"`
	ServiceName string        `json:"serviceName"`
}
