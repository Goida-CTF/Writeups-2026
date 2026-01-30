package ws

import (
	"encoding/json"
	"fmt"

	"znanie-drevnih/internal/models/wsmodels"
)

type clientMessage struct {
	Type     wsmodels.ClientMessageType
	Submit   *wsmodels.ClientMessageSubmit
	Continue *wsmodels.ClientMessageContinue
	Start    *wsmodels.ClientMessageStart
}

func parseClientMessage(payload []byte) (*clientMessage, error) {
	var envelope struct {
		Type wsmodels.ClientMessageType `json:"type"`
	}
	if err := json.Unmarshal(payload, &envelope); err != nil {
		return nil, err
	}

	switch envelope.Type {
	case wsmodels.ClientMessageTypeSubmit:
		var msg wsmodels.ClientMessageSubmit
		if err := json.Unmarshal(payload, &msg); err != nil {
			return nil, fmt.Errorf("submit: %w", err)
		}
		return &clientMessage{
			Type:   envelope.Type,
			Submit: &msg,
		}, nil
	case wsmodels.ClientMessageTypeContinue:
		var msg wsmodels.ClientMessageContinue
		if err := json.Unmarshal(payload, &msg); err != nil {
			return nil, fmt.Errorf("continue: %w", err)
		}
		return &clientMessage{
			Type:     envelope.Type,
			Continue: &msg,
		}, nil
	case wsmodels.ClientMessageTypeStart:
		var msg wsmodels.ClientMessageStart
		if err := json.Unmarshal(payload, &msg); err != nil {
			return nil, fmt.Errorf("start: %w", err)
		}
		return &clientMessage{
			Type:  envelope.Type,
			Start: &msg,
		}, nil
	default:
		return &clientMessage{Type: envelope.Type}, nil
	}
}
