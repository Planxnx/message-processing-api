package model

import "time"

type HealthCheckResponseData struct {
	Feature     string     `json:"feature"`
	Description string     `json:"description"`
	ExecuteMode []string   `json:"executeMode"`
	ServiceName string     `json:"serviceName"`
	Status      bool       `json:"status"`
	LastOnline  *time.Time `json:"lastOnline,omitempty"`
}
