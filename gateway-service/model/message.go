package model

type MessageRequest struct {
	Message string                 `json:"message"`
	UserRef string                 `json:"userRef"`        //end-user reference
	Feature string                 `json:"feature"`        //Feature
	Data    map[string]interface{} `json:"data,omitempty"` //attachment
}

type MessageResponseData struct {
	MessageRef string `json:"messageRef,omitempty"` //message reference
}
