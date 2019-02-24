package helpers

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type causer interface {
	Cause() error
}

func PrintErr(err error) {
	if err == nil {
		fmt.Println(nil)
	}

	errStr := err.Error()
	parts := strings.Split(errStr, ": ")
	for _, part := range parts {
		fmt.Println(part)
	}

	stackTrace := errors.StackTrace{}
	for err != nil {
		if tracer, ok := err.(stackTracer); ok {
			stackTrace = tracer.StackTrace()
		}
		if causer, ok := err.(causer); ok {
			err = causer.Cause()
		} else {
			err = nil
		}
	}

	for _, frame := range stackTrace {
		fmt.Printf("%+v\n", frame)
	}
}
