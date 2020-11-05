package model

type ProviderResgisterRequest struct {
	ProviderID   string `json:"providerID"`
	ProviderName string `json:"providerName"`
	Webhook      string `json:"webhook"`
}

type ProviderResgisterResponseData struct {
	ProviderID string `json:"providerID,omitempty"`
	Secret     string `json:"secret,omitempty"`
}
