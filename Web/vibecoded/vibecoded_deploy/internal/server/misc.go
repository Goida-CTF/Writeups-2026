package server

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"reflect"
	"strings"
)

type ctxKey int

const CtxKeyJWTToken ctxKey = iota

func DecodeJSON(w http.ResponseWriter, r *http.Request, v any) bool {
	contentType := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		HandleClientError(w, r,
			fmt.Errorf("mime.ParseMediaType: %w", err),
			"Error parsing Content-Type", http.StatusBadRequest)
		return false
	}
	if mediatype != "application/json" {
		HandleClientError(w, r,
			ErrContentTypeNotJSON,
			ErrContentTypeNotJSON.Error(), http.StatusUnsupportedMediaType)
		return false
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		HandleClientError(w, r,
			fmt.Errorf("io.ReadAll: %w", err),
			"Error reading request body", http.StatusBadRequest)
		return false
	}
	defer func() {
		_ = r.Body.Close()
	}()

	if len(strings.TrimSpace(string(body))) == 0 {
		HandleClientError(w, r,
			ErrEmptyRequestBody,
			ErrEmptyRequestBody.Error(), http.StatusBadRequest)
		return false
	}

	if err := json.Unmarshal(body, v); err != nil {
		HandleClientError(w, r,
			fmt.Errorf("json.NewDecoder: %w", err),
			"Invalid JSON", http.StatusBadRequest)
		return false
	}

	return true
}

func renderJSONWithCode(w http.ResponseWriter, code int, v any) error {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("renderJSONWithCode: json.Marshal: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)

	if _, err = w.Write(jsonBytes); err != nil {
		return fmt.Errorf("renderJSONWithCode: w.Write: %w", err)
	}

	return nil
}

func RenderJSON(w http.ResponseWriter, r *http.Request, v any) {
	if err := renderJSONWithCode(w, http.StatusOK, &v); err != nil {
		HandleInternalServerError(w, r, err)
	}
}

func CheckAllFieldsExist(w http.ResponseWriter, r *http.Request, v any) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		HandleInternalServerError(w, r,
			fmt.Errorf("CheckAllFieldsExist: Invalid model provided"))
		return false
	}

	rv = rv.Elem()
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		fieldType := rv.Type().Field(i)
		jsonTag := fieldType.Tag.Get("json")

		if jsonTag == "-" || jsonTag == "" {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			if field.String() != "" {
				continue
			}

			err := fmt.Errorf("%w: %v", ErrEmptyField, jsonTag)
			HandleClientError(w, r,
				err,
				err.Error(), http.StatusBadRequest)
			return false
		}
	}
	return true
}
