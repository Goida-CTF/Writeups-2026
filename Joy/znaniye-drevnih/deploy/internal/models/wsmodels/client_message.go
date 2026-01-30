package wsmodels

type ClientMessageType string

const (
	ClientMessageTypeSubmit   ClientMessageType = "submit"
	ClientMessageTypeContinue ClientMessageType = "continue"
	ClientMessageTypeStart    ClientMessageType = "start"
)

type clientMessageBase struct {
	Type ClientMessageType `json:"type"`
}

type ClientMessageSubmit struct {
	clientMessageBase
	Code string `json:"code"`
}

type ClientMessageContinue struct {
	clientMessageBase
}

type ClientMessageStart struct {
	clientMessageBase
}
