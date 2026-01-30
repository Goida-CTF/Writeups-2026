package ws

import (
	"strings"

	"znanie-drevnih/internal/models"
	"znanie-drevnih/internal/models/taskmodels"
	"znanie-drevnih/internal/models/wsmodels"
)

type messageKind int

const (
	msgUnknown messageKind = iota
	msgNext
	msgRestart
	msgCode
)

func classifyMessage(msg *clientMessage) messageKind {
	if msg == nil {
		return msgUnknown
	}
	switch strings.TrimSpace(strings.ToLower(string(msg.Type))) {
	case string(wsmodels.ClientMessageTypeStart):
		return msgRestart
	case string(wsmodels.ClientMessageTypeContinue):
		return msgNext
	case string(wsmodels.ClientMessageTypeSubmit):
		return msgCode
	default:
		return msgUnknown
	}
}

func newDialogMessage(dialog *taskmodels.Dialog) *wsmodels.DialogServerMessage {
	msg := &wsmodels.DialogServerMessage{
		Role:  toRole(dialog.Role),
		Audio: dialog.Audio,
		Text:  dialog.Text,
	}
	msg.Type = wsmodels.ServerMessageTypeDialog
	return msg
}

func newPhraseMessage(phrase *taskmodels.Phrase, flag string) *wsmodels.PhraseServerMessage {
	msg := &wsmodels.PhraseServerMessage{
		Role:  toRole(phrase.Role),
		Audio: phrase.Audio,
		Text:  phrase.Text,
		Flag:  flag,
	}
	msg.Type = wsmodels.ServerMessageTypePhrase
	return msg
}

func newTaskMessage(task *taskmodels.Task) *wsmodels.TaskServerMessage {
	msg := &wsmodels.TaskServerMessage{
		Text: task.Text,
	}
	msg.Type = wsmodels.ServerMessageTypeTask
	return msg
}

func newResultMessage(result *models.TaskRunResult) *wsmodels.ResultServerMessage {
	msg := &wsmodels.ResultServerMessage{
		CouldRun:           result.CouldRun,
		ContainsNonRuChars: result.ContainsNonRuChars,
		Success:            result.Passed,
	}
	msg.Type = wsmodels.ServerMessageTypeResult
	return msg
}

func newGameStatusMessage(status wsmodels.GameStatus) *wsmodels.GameStatusServerMessage {
	msg := &wsmodels.GameStatusServerMessage{
		Status: status,
	}
	msg.Type = wsmodels.ServerMessageTypeGame
	return msg
}

func toRole(role string) wsmodels.Role {
	switch strings.ToLower(role) {
	case string(wsmodels.RoleSlav):
		return wsmodels.RoleSlav
	case string(wsmodels.RoleZuck):
		return wsmodels.RoleZuck
	default:
		return wsmodels.Role(role)
	}
}
