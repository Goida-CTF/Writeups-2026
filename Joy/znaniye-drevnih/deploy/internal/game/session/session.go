package session

import (
	"crypto/rand"
	"math/big"

	"znanie-drevnih/internal/models/taskmodels"
)

const stepPhraseChancePercent = 30

type Stage int

const (
	StageDialog Stage = iota
	StageTasks
	StageDone
)

type Session struct {
	speech        taskmodels.SpeechMap
	tasks         []taskmodels.Task
	partsRequired uint64

	stage     Stage
	dialogKey string
	dialogIdx int
	taskIdx   int
	taskOrder []int
}

func New(data *taskmodels.GameData, partsRequired uint64) *Session {
	s := &Session{
		speech:        data.Speech,
		tasks:         data.Tasks,
		partsRequired: partsRequired,
	}
	s.Reset()
	return s
}

func (s *Session) Reset() {
	s.stage = StageDialog
	s.dialogKey = "start"
	s.dialogIdx = 0
	s.taskIdx = 0
	s.taskOrder = buildTaskOrder(s.tasks, s.partsRequired)
}

func (s *Session) Stage() Stage {
	return s.stage
}

func (s *Session) SetStage(stage Stage) {
	s.stage = stage
}

func (s *Session) CurrentDialog() (*taskmodels.Dialog, bool) {
	dialogs, ok := s.speech.Dialogs[s.dialogKey]
	if !ok || s.dialogIdx >= len(dialogs) {
		return nil, false
	}
	return &dialogs[s.dialogIdx], true
}

func (s *Session) IncDialog() {
	s.dialogIdx++
}

func (s *Session) CurrentTask() (*taskmodels.Task, bool) {
	if s.taskIdx >= len(s.taskOrder) {
		return nil, false
	}
	idx := s.taskOrder[s.taskIdx]
	if idx < 0 || idx >= len(s.tasks) {
		return nil, false
	}
	return &s.tasks[idx], true
}

func (s *Session) IncTask() {
	s.taskIdx++
}

func (s *Session) RandomFailPhrase() (*taskmodels.Phrase, bool) {
	return randomPhrase(s.speech.Phrases.Fail)
}

func (s *Session) RandomStepPhrase() (*taskmodels.Phrase, bool) {
	return randomPhrase(s.speech.Phrases.Steps)
}

func (s *Session) RandomWinPhrase() (*taskmodels.Phrase, bool) {
	return randomPhrase(s.speech.Phrases.Win)
}

func (s *Session) ShouldSendStepPhrase() bool {
	value, ok := randInt(100)
	if !ok {
		return false
	}
	return value < stepPhraseChancePercent
}

func randomPhrase(phrases []taskmodels.Phrase) (*taskmodels.Phrase, bool) {
	if len(phrases) == 0 {
		return nil, false
	}
	index, ok := randInt(len(phrases))
	if !ok {
		return nil, false
	}
	return &phrases[index], true
}

func buildTaskOrder(tasks []taskmodels.Task, partsRequired uint64) []int {
	if len(tasks) == 0 {
		return nil
	}

	required := int(partsRequired)
	if required <= 0 {
		required = len(tasks)
	}

	idxZero := -1
	var others []int
	for i, task := range tasks {
		if task.ID == 0 && idxZero == -1 {
			idxZero = i
			continue
		}
		others = append(others, i)
	}
	shuffleInts(others)

	order := make([]int, 0, required)
	if idxZero != -1 {
		order = append(order, idxZero)
	}

	remaining := required - len(order)
	if remaining <= 0 {
		return order[:required]
	}

	if partsRequired > uint64(len(tasks)) {
		for i := 0; i < remaining; i++ {
			if len(others) == 0 {
				order = append(order, idxZero)
				continue
			}
			if i > 0 && i%len(others) == 0 {
				shuffleInts(others)
			}
			order = append(order, others[i%len(others)])
		}
		return order
	}

	if remaining > len(others) {
		remaining = len(others)
	}
	order = append(order, others[:remaining]...)
	return order
}

func shuffleInts(values []int) {
	for i := len(values) - 1; i > 0; i-- {
		j, ok := randInt(i + 1)
		if !ok {
			return
		}
		values[i], values[j] = values[j], values[i]
	}
}

func randInt(max int) (int, bool) {
	if max <= 0 {
		return 0, false
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, false
	}
	return int(n.Int64()), true
}
