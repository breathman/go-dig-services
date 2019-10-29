package util

import (
	"errors"
	"fmt"
	"strings"
)

// ErrorSlice slice of errors
type ErrorSlice struct {
	error
	Errors []error
}

// NewErrorSlice - constructor for ErrorSlice
func NewErrorSlice(slice []string) ErrorSlice {
	errs := make([]error, len(slice))
	for i, s := range slice {
		errs[i] = errors.New(s)
	}
	return ErrorSlice{Errors: errs}
}

func (es *ErrorSlice) Error() string {
	return strings.Join(es.Messages(), "\n")
}

// Append mutates ErrorSlice state
func (es *ErrorSlice) Append(err error) {
	if err == nil {
		return
	}
	es.Errors = append(es.Errors, err)
}

// Messages maps errors to error messages as strings
func (es ErrorSlice) Messages() []string {
	msgs := make([]string, len(es.Errors))
	for i, e := range es.Errors {
		msgs[i] = e.Error()
	}
	return msgs
}

// Join joins error messages with given separator string
func (es ErrorSlice) Join(sep string) error {
	errMsg := strings.Join(es.Messages(), sep)
	return fmt.Errorf(errMsg)
}

// Joinnl - join with default separator '\n'
func (es ErrorSlice) Joinnl() error {
	return es.Join("\n")
}

// Empty - check slice is empty
func (es ErrorSlice) Empty() bool {
	return len(es.Errors) == 0
}

// NilIfEmpty - return nil if no errors
func (es *ErrorSlice) NilIfEmpty() error {
	if len(es.Errors) > 0 {
		return es
	}
	return nil
}
