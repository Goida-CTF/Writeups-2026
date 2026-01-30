package models

type TaskRunResult struct {
	Passed             bool
	CouldRun           bool
	ContainsNonRuChars bool
}
