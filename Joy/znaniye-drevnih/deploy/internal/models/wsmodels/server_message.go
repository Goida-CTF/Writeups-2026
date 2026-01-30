package wsmodels

type ServerMessageType string

const (
	ServerMessageTypeDialog ServerMessageType = "dialog"
	ServerMessageTypePhrase ServerMessageType = "phrase"
	ServerMessageTypeTask   ServerMessageType = "task"
	ServerMessageTypeResult ServerMessageType = "result"
	ServerMessageTypeGame   ServerMessageType = "game"
)

type serverMessageBase struct {
	Type ServerMessageType `json:"type"`
}

type Role string

const (
	RoleSlav Role = "slav"
	RoleZuck Role = "zuck"
)

type DialogServerMessage struct {
	serverMessageBase
	Role  Role   `json:"role"`
	Audio string `json:"audio"` // Base64-encoded mp3 data
	Text  string `json:"text"`
}

type PhraseServerMessage struct {
	serverMessageBase
	Role  Role   `json:"role"`
	Audio string `json:"audio"` // Base64-encoded mp3 data
	Text  string `json:"text"`
	Flag  string `json:"flag,omitempty"` // Only for win step
}

type TaskServerMessage struct {
	serverMessageBase
	Text string `json:"text"`
}

type ResultServerMessage struct {
	serverMessageBase
	CouldRun           bool `json:"couldRun"`
	ContainsNonRuChars bool `json:"containsNonRuChars"`
	Success            bool `json:"success"` // Success determines if task part was completed
}

type GameStatus string

const (
	GameStatusTimeout GameStatus = "timeout"
	GameStatusWrong   GameStatus = "wrong"
	GameStatusWin     GameStatus = "win"
	GameStatusError   GameStatus = "error"
)

type GameStatusServerMessage struct {
	serverMessageBase
	Status GameStatus `json:"status"`
}
