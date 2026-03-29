package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"

	"github.com/gin-gonic/gin"
)

func HasAtLeastOneField(v interface{}) bool {
	val := reflect.ValueOf(v)

	// Se for ponteiro, pega o valor
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Itera sobre os campos da struct
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		// Verifica se o campo é um ponteiro e não é nil
		if field.Kind() == reflect.Ptr && !field.IsNil() {
			return true
		}

		// Para campos não-ponteiro, verifica se não é zero value
		if field.Kind() != reflect.Ptr && !field.IsZero() {
			return true
		}
	}

	return false
}

func GetApiErr(err error) *HTTPError {
	var apiErr *HTTPError

	if errors.As(err, &apiErr) {
		return apiErr
	} else {
		return NewHTTPError(http.StatusInternalServerError, "an unexpected error occurred")
	}
}

func GetJSONKeys(c *gin.Context) ([]string, error) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var jsonData map[string]interface{}
	if err := json.Unmarshal(body, &jsonData); err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(jsonData))
	for key := range jsonData {
		keys = append(keys, key)
	}

	return keys, nil
}

func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func ContainsSome(slice []string, items []string) bool {
	for _, item := range items {
		if Contains(slice, item) {
			return true
		}
	}
	return false
}

func IsValidURL(str string) bool {
	u, err := url.ParseRequestURI(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func NewValue[T any](v T) OptionalNullable[T] {
	return OptionalNullable[T]{Set: true, Value: &v}
}

func NewNull[T any]() OptionalNullable[T] {
	return OptionalNullable[T]{Set: true, Value: nil}
}

func Unset[T any]() OptionalNullable[T] {
	return OptionalNullable[T]{Set: false}
}
