package game

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"znanie-drevnih/internal/models/taskmodels"
)

const (
	speechMapFile = "speech_map.yaml"
	tasksFile     = "tasks.yaml"
)

func loadGameData(taskDataPath string) (*taskmodels.GameData, error) {
	speechPath := filepath.Join(taskDataPath, speechMapFile)
	tasksPath := filepath.Join(taskDataPath, tasksFile)

	speechBytes, err := os.ReadFile(speechPath)
	if err != nil {
		return nil, fmt.Errorf("read speech map: %w", err)
	}
	var speech taskmodels.SpeechMap
	if err := yaml.Unmarshal(speechBytes, &speech); err != nil {
		return nil, fmt.Errorf("unmarshal speech map: %w", err)
	}
	if err := embedAudioBase64(taskDataPath, &speech); err != nil {
		return nil, fmt.Errorf("embed audio: %w", err)
	}

	tasksBytes, err := os.ReadFile(tasksPath)
	if err != nil {
		return nil, fmt.Errorf("read tasks: %w", err)
	}
	var taskFile taskmodels.TaskFile
	if err := yaml.Unmarshal(tasksBytes, &taskFile); err != nil {
		return nil, fmt.Errorf("unmarshal tasks: %w", err)
	}

	return &taskmodels.GameData{
		Speech: speech,
		Tasks:  taskFile.Tasks,
	}, nil
}

func embedAudioBase64(taskDataPath string, speech *taskmodels.SpeechMap) error {
	for key, dialogs := range speech.Dialogs {
		for i := range dialogs {
			encoded, err := loadAudioBase64(taskDataPath, dialogs[i].Audio)
			if err != nil {
				return fmt.Errorf("dialog audio %s: %w", dialogs[i].Audio, err)
			}
			dialogs[i].Audio = encoded
		}
		speech.Dialogs[key] = dialogs
	}

	for i := range speech.Phrases.Fail {
		encoded, err := loadAudioBase64(taskDataPath, speech.Phrases.Fail[i].Audio)
		if err != nil {
			return fmt.Errorf("fail phrase audio %s: %w", speech.Phrases.Fail[i].Audio, err)
		}
		speech.Phrases.Fail[i].Audio = encoded
	}
	for i := range speech.Phrases.Steps {
		encoded, err := loadAudioBase64(taskDataPath, speech.Phrases.Steps[i].Audio)
		if err != nil {
			return fmt.Errorf("step phrase audio %s: %w", speech.Phrases.Steps[i].Audio, err)
		}
		speech.Phrases.Steps[i].Audio = encoded
	}
	for i := range speech.Phrases.Win {
		encoded, err := loadAudioBase64(taskDataPath, speech.Phrases.Win[i].Audio)
		if err != nil {
			return fmt.Errorf("win phrase audio %s: %w", speech.Phrases.Win[i].Audio, err)
		}
		speech.Phrases.Win[i].Audio = encoded
	}

	return nil
}

func loadAudioBase64(taskDataPath, audioPath string) (string, error) {
	if audioPath == "" {
		return "", nil
	}
	fullPath := filepath.Join(taskDataPath, "audio", filepath.FromSlash(audioPath))
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}
