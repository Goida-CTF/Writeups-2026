package usecases

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode"

	piston "github.com/milindmadhukar/go-piston"
	"go.uber.org/zap"

	"znanie-drevnih/internal/models"
	"znanie-drevnih/internal/models/taskmodels"
)

var errTaskIsEmpty = errors.New("task is empty")

func (u *UseCases) RunTask(ctx context.Context,
	task *taskmodels.Task,
	bxxCode string,
) (*models.TaskRunResult, error) {
	if task == nil {
		return nil, errTaskIsEmpty
	}

	containsNonRuChars := containsNonRuChars(bxxCode)

	result := &models.TaskRunResult{
		Passed:             false,
		CouldRun:           true,
		ContainsNonRuChars: containsNonRuChars,
	}

	for _, test := range task.Tests {
		if len(test.AnyOf) > 0 {
			passed := false
			couldRun := false
			for _, candidate := range test.AnyOf {
				ok, runOk, err := u.runSingleTest(ctx, bxxCode, candidate.Stdin, candidate.Stdout)
				if err != nil {
					return nil, err
				}
				if !runOk {
					result.CouldRun = false
					return result, nil
				}
				couldRun = true
				if ok {
					passed = true
					break
				}
			}
			if !couldRun {
				result.CouldRun = false
				return result, nil
			}
			if !passed {
				return result, nil
			}
			continue
		}

		ok, runOk, err := u.runSingleTest(ctx, bxxCode, test.Stdin, test.Stdout)
		if err != nil {
			return nil, fmt.Errorf("u.runSingleTest: %w", err)
		}
		if !runOk {
			result.CouldRun = false
			return result, nil
		}
		if !ok {
			return result, nil
		}
	}

	result.Passed = true
	return result, nil
}

func (u *UseCases) runSingleTest(ctx context.Context,
	bxxCode, stdin, expected string,
) (bool, bool, error) {
	bxxCode = strings.Replace(bxxCode, "СлавянскийC++.h", "Ве_крест_крест.h", 1) // Ugly, but quick fix

	res, err := u.client.RunBXXCode(ctx, bxxCode, stdin)
	if err != nil {
		return false, false, fmt.Errorf("u.client.RunBXXCode: %w", err)
	}

	if isCompileError(res) {
		u.logger.Info("compile failed",
			zap.String("compileMessage", res.Compile.Message),
			zap.String("compileStderr", res.Compile.Stderr),
			zap.Int("compileCode", res.Compile.Code),
		)
		return false, false, nil
	}
	if res.Run.Code != 0 || strings.EqualFold(res.Run.Status, "error") {
		u.logger.Info("run failed",
			zap.String("runMessage", res.Run.Message),
			zap.String("runStderr", res.Run.Stderr),
			zap.Int("runCode", res.Run.Code),
		)
		return false, false, nil
	}

	output := res.Run.Output
	if output == "" {
		output = res.Run.Stdout
	}

	return normalizeOutput(output) == normalizeOutput(expected), true, nil
}

func isCompileError(res *piston.PistonExecution) bool {
	if res == nil {
		return true
	}
	if res.Compile.Code != 0 {
		return true
	}
	if strings.EqualFold(res.Compile.Status, "error") {
		return true
	}
	if res.Compile.Message != "" || res.Compile.Stderr != "" {
		return true
	}
	return false
}

func normalizeOutput(output string) string {
	output = strings.ToLower(output)
	output = strings.ReplaceAll(output, "\r\n", "\n")
	lines := strings.Split(output, "\n")
	for i := range lines {
		lines[i] = strings.TrimRightFunc(lines[i], unicode.IsSpace)
	}
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return strings.Join(lines, "\n")
}

func containsNonRuChars(code string) bool {
	for line := range strings.SplitSeq(code, "\n") {
		scanLine, allowRanges := includeAllowRanges(line)
		for idx, r := range scanLine {
			if !unicode.IsLetter(r) {
				continue
			}
			if unicode.In(r, unicode.Cyrillic) {
				continue
			}
			if unicode.In(r, unicode.Latin) && inRanges(idx, allowRanges) {
				continue
			}
			return true
		}
	}
	return false
}

type byteRange struct {
	start int
	end   int
}

func includeAllowRanges(line string) (string, []byteRange) {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "#include") {
		return line, nil
	}
	ranges := []byteRange{
		{start: 0, end: len("#include")},
	}

	if start := strings.Index(trimmed, "<"); start != -1 {
		if end := strings.Index(trimmed[start+1:], ">"); end != -1 {
			end += start + 1
			ranges = append(ranges, byteRange{start: start + 1, end: end})
		}
	} else if start := strings.Index(trimmed, "\""); start != -1 {
		if end := strings.Index(trimmed[start+1:], "\""); end != -1 {
			end += start + 1
			ranges = append(ranges, byteRange{start: start + 1, end: end})
		}
	}

	return trimmed, ranges
}

func inRanges(idx int, ranges []byteRange) bool {
	for _, r := range ranges {
		if idx >= r.start && idx < r.end {
			return true
		}
	}
	return false
}
