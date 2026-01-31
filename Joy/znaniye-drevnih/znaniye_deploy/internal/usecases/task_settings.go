package usecases

import "time"

func (u *UseCases) TaskPartTimeout() time.Duration {
	return u.game.TaskPartTimeout()
}

func (u *UseCases) TaskFlag() string {
	return u.game.TaskFlag()
}
