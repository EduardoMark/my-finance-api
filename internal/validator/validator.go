package validator

import (
	"context"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Validator interface {
	Valid(context.Context) Evaluator
}

type Evaluator map[string]any

func (e *Evaluator) AddFieldError(key, message string) {
	if *e == nil {
		*e = make(map[string]any)
	}

	if _, exists := (*e)[key]; !exists {
		(*e)[key] = message
	}
}

func (e *Evaluator) CheckField(ok bool, key, message string) {
	if !ok {
		e.AddFieldError(key, message)
	}
}

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MinChars(value string, len int) bool {
	return utf8.RuneCountInString(value) >= len
}

func MaxChars(value string, len int) bool {
	return utf8.RuneCountInString(value) <= len
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func TransactionType(value string) bool {
	value = strings.ToLower(value)
	return value == "income" || value == "expense"
}

func CheckBalance(value float64) bool {
	return value >= 0
}
