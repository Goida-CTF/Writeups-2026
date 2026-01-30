package usecases

import "znanie-drevnih/internal/game/session"

func (u *UseCases) StartSession() (*session.Session, error) {
	return u.game.StartSession()
}
