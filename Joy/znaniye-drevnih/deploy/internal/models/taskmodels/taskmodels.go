package taskmodels

type SpeechMap struct {
	Dialogs map[string][]Dialog `yaml:"dialogs" json:"dialogs"`
	Phrases Phrases             `yaml:"phrases" json:"phrases"`
}

type Dialog struct {
	ID    int    `yaml:"id,omitempty" json:"id,omitempty"`
	Role  string `yaml:"role" json:"role"`
	Audio string `yaml:"audio" json:"audio"`
	Text  string `yaml:"text" json:"text"`
}

type Phrases struct {
	Fail  []Phrase `yaml:"fail" json:"fail"`
	Steps []Phrase `yaml:"steps" json:"steps"`
	Win   []Phrase `yaml:"win" json:"win"`
}

type Phrase struct {
	Role  string `yaml:"role" json:"role"`
	Audio string `yaml:"audio" json:"audio"`
	Text  string `yaml:"text" json:"text"`
}

type TaskFile struct {
	Tasks []Task `yaml:"tasks"`
}

type Task struct {
	ID    int      `yaml:"id" json:"id"`
	Text  string   `yaml:"text" json:"text"`
	Task  TaskMeta `yaml:"task" json:"task"`
	Tests []Test   `yaml:"tests" json:"-"`
}

type TaskMeta struct {
	Type string `yaml:"type" json:"type"`
}

type Test struct {
	Stdin  string     `yaml:"stdin,omitempty" json:"-"`
	Stdout string     `yaml:"stdout,omitempty" json:"-"`
	AnyOf  []TestCase `yaml:"anyOf,omitempty" json:"-"`
}

type TestCase struct {
	Stdin  string `yaml:"stdin,omitempty" json:"-"`
	Stdout string `yaml:"stdout,omitempty" json:"-"`
}

type GameData struct {
	Speech SpeechMap
	Tasks  []Task
}
