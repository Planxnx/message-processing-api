package provider

import "time"

type ProviderData struct {
	ID        string    `bson:"id" json:"id"`
	Name      string    `bson:"name" json:"name"`
	Webhook   string    `bson:"webhook" json:"webhook"`
	Secret    string    `bson:"secret" json:"secret"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type CreateProviderResult struct {
	ID     string `bson:"id" json:"id"`
	Secret string `bson:"secret" json:"secret"`
}
