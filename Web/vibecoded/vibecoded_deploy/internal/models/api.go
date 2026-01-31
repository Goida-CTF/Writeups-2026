package models

type ResponseStatus string

const (
	ResponseStatusOK    ResponseStatus = "ok"
	ResponseStatusError ResponseStatus = "error"
)

type JsonResponse struct {
	Status ResponseStatus `json:"status"`
	Error  string         `json:"error,omitempty"`
	Data   any            `json:"data,omitempty"`
}

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegistrationForm struct {
	LoginForm
}

type PasswordComplexityRequest struct {
	Password string `json:"password"`
}

type PasswordComplexityResult struct {
	Ok    bool   `json:"ok"`
	Level int    `json:"level"`
	Text  string `json:"text"`
}

type Note struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type NoteUpdate struct {
	ID int `json:"id"`
	Note
}
