package ws

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"

	"znanie-drevnih/internal/game/session"
	"znanie-drevnih/internal/models/taskmodels"
	"znanie-drevnih/internal/models/wsmodels"
)

func (s *Service) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.l.Error("upgrade to WS", zap.Error(err))
		return
	}
	defer func() { _ = conn.Close() }()

	conn.SetReadLimit(maxMessageSize)
	_ = conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		_ = conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	incoming := make(chan *clientMessage)
	outgoing := make(chan any, 8)
	errs := make(chan error, 1)
	done := make(chan struct{})

	go readLoop(conn, incoming, errs, done, s.l)
	go writeLoop(conn, outgoing, errs, done)

	var sess *session.Session
	ctx := context.Background()
	partTimeout := s.uc.TaskPartTimeout()
	taskFlag := s.uc.TaskFlag()

	connTimer := time.NewTimer(time.Hour)
	defer connTimer.Stop()

	var taskTimer *time.Timer
	var taskTimeout <-chan time.Time
	resetTaskTimer := func() {
		if taskTimer != nil {
			if !taskTimer.Stop() {
				select {
				case <-taskTimer.C:
				default:
				}
			}
		}
		taskTimer = time.NewTimer(partTimeout)
		taskTimeout = taskTimer.C
	}
	stopTaskTimer := func() {
		if taskTimer == nil {
			return
		}
		if !taskTimer.Stop() {
			select {
			case <-taskTimer.C:
			default:
			}
		}
		taskTimeout = nil
	}
	startTaskTimer := func(task *taskmodels.Task) {
		if task == nil || task.ID == 0 {
			stopTaskTimer()
			return
		}
		resetTaskTimer()
	}
	endSessionWithError := func(err error) {
		if err != nil {
			s.l.Error("ws session error", zap.Error(err))
		}
		stopTaskTimer()
		sess = nil
		outgoing <- newGameStatusMessage(wsmodels.GameStatusError)
	}

	sendInitial := func() {
		if sess == nil {
			return
		}
		if dialog, ok := sess.CurrentDialog(); ok {
			outgoing <- newDialogMessage(dialog)
			return
		}
		sess.SetStage(session.StageTasks)
		if task, ok := sess.CurrentTask(); ok {
			outgoing <- newTaskMessage(task)
			startTaskTimer(task)
			return
		}
		sess.SetStage(session.StageDone)
	}

	for {
		select {
		case msg := <-incoming:
			switch classifyMessage(msg) {
			case msgRestart:
				stopTaskTimer()
				newSess, err := s.uc.StartSession()
				if err != nil {
					endSessionWithError(err)
					break
				}
				sess = newSess
				sendInitial()
			case msgNext:
				if sess == nil || sess.Stage() != session.StageDialog {
					break
				}
				sess.IncDialog()
				if dialog, ok := sess.CurrentDialog(); ok {
					outgoing <- newDialogMessage(dialog)
					break
				}
				sess.SetStage(session.StageTasks)
				if task, ok := sess.CurrentTask(); ok {
					outgoing <- newTaskMessage(task)
					startTaskTimer(task)
				} else {
					sess.SetStage(session.StageDone)
				}
			case msgCode:
				if sess == nil || sess.Stage() != session.StageTasks {
					break
				}
				task, ok := sess.CurrentTask()
				if !ok {
					sess.SetStage(session.StageDone)
					break
				}
				code := ""
				if msg != nil && msg.Submit != nil {
					code = msg.Submit.Code
				}
				if code == "" {
					startTaskTimer(task)
					break
				}
				stopTaskTimer()
				result, err := s.uc.RunTask(ctx, task, code)
				if err != nil {
					endSessionWithError(err)
					break
				}
				if result != nil {
					outgoing <- newResultMessage(result)
				}
				if result == nil || !result.Passed {
					if phrase, ok := sess.RandomFailPhrase(); ok {
						outgoing <- newPhraseMessage(phrase, "")
					}
					outgoing <- newGameStatusMessage(wsmodels.GameStatusWrong)
					startTaskTimer(task)
					break
				}

				if sess.ShouldSendStepPhrase() {
					if phrase, ok := sess.RandomStepPhrase(); ok {
						outgoing <- newPhraseMessage(phrase, "")
					}
				}

				sess.IncTask()
				if task, ok := sess.CurrentTask(); ok {
					outgoing <- newTaskMessage(task)
					startTaskTimer(task)
					break
				}

				sess.SetStage(session.StageDone)
				if phrase, ok := sess.RandomWinPhrase(); ok {
					outgoing <- newPhraseMessage(phrase, taskFlag)
				}
				outgoing <- newGameStatusMessage(wsmodels.GameStatusWin)
			default:
			}
		case <-taskTimeout:
			outgoing <- newGameStatusMessage(wsmodels.GameStatusTimeout)
			stopTaskTimer()
			sess = nil
		case <-connTimer.C:
			close(done)
			return
		case err := <-errs:
			if err != nil {
				s.l.Info("ws connection closed", zap.Error(err))
			}
			close(done)
			return
		}
	}
}
