package util

import (
	"net/url"
	"reflect"
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
