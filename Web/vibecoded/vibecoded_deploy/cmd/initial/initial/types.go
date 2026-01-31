package initial

type InitialNote struct {
	Title   string
	Content string
}

type InitialUser struct {
	Username string
	Password string
	IsAdmin  bool
	Notes    []InitialNote
}
