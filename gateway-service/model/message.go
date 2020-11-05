package model

type MessageRequest struct {
	Message  string                 `json:"message"`
	UserRef  string                 `json:"userRef"`        //end-user reference
	Features map[string]bool        `json:"features"`       //Feature this message will uses next
	Data     map[string]interface{} `json:"data,omitempty"` //attachment
}

type MessageResponseData struct {
	MessageRef string `json:"messageRef,omitempty"` //message reference
}
