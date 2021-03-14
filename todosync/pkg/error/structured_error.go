// Package error contains common error handling primitives.
// Most notably it contains a type which allows for more structured error handling.
package error

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"time"
)

// shamelessly copied/inspired from https://www.orsolabs.com/post/go-errors-and-logs/
type StructuredError struct {
	debugID string
	time    time.Time
	msg     string
	err     error
	pairs   []interface{}
}

func (e StructuredError) Error() string {
	return fmt.Sprintf("%s: %v", e.msg, e.err)
}

func (e StructuredError) Unwrap() error {
	return e.err
}

func (e StructuredError) DebugID() string {
	return e.debugID
}

func (e StructuredError) Time() time.Time {
	return e.time
}

func (e StructuredError) Root() *StructuredError {
	root := &e

	for {
		_, ok := root.err.(*StructuredError)
		if !ok {
			break
		}

		root, _ = root.err.(*StructuredError)
	}

	return root
}

func Err(msg string, err error, pairs ...interface{}) error {
	var id string
	var t time.Time

	_, ok := err.(*StructuredError)
	if !ok {
		// we only require the id and time for the innermost error
		id = GenDebugID()
		t = time.Now()

		pairs = append(pairs, "debugID", id, "time", t.String())
	}

	return &StructuredError{
		debugID: id,
		time:    t,
		msg:     msg,
		err:     err,
		pairs:   pairs,
	}
}

func GenDebugID() string {
	return uuid.NewV4().String()[:8]
}

func DeepKeyVals(err error, pairs ...interface{}) []interface{} {
	keyvals := []interface{}{"error", err.Error()}
	keyvals = append(keyvals, pairs...)

	structErr, ok := err.(*StructuredError)
	for {
		if !ok {
			return keyvals
		}

		keyvals = append(keyvals, structErr.pairs...)
		structErr, ok = structErr.err.(*StructuredError)
	}
}
