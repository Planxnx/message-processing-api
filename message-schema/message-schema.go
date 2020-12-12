package messageschema

import "time"

type DefaultMessageFormat struct {
	Message       string                 `json:"message"`
	Ref1          string                 `json:"ref1"`  //client reference
	Ref2          string                 `json:"ref2"`  //message reference
	Ref3          string                 `json:"ref3"`  //end-user reference
	Owner         string                 `json:"owner"` //service reference
	PublishedBy   string                 `json:"publishedBy"`
	PublishedAt   time.Time              `json:"publishedAt"` //messageQueue published time
	Features      map[string]bool        `json:"features"`    //Feature this message will uses next
	Data          map[string]interface{} `json:"data"`        //attachment
	Type          string                 `json:"type"`        //message type eg. reply message, notification
	ExcuteMode    ExecuteMode            `json:"excuteMode"`
	CallbackTopic string                 `json:"callbackTopic"` //specific callback topic for
	Error         string                 `json:"error"`
}

//ExecuteMode is execution mode async or sync
type ExecuteMode string

const (
	SynchronousMode  ExecuteMode = "SYNC"
	AsynchronousMode ExecuteMode = "ASYNC"
)

const (
	CommonMessage   string = "commonMessage"
	ReplyMessage    string = "replyMessage"
	Notification    string = "notification"
	ChitchatMessage string = "chitchat"
)
