package healthz

import (
	"errors"
	"reflect"
)

var (
	ErrorHealthCheckNotAFunc = errors.New("the given health check is not a func")
	ErrorHealthCheckNoError  = errors.New("the given health check does not return *healthz.Error")

	errorInterface = reflect.TypeOf((*error)(nil)).Elem()
)
