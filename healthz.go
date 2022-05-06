package healthz

import (
	"reflect"
	"time"

	"github.com/go-mojito/mojito"
	"github.com/go-mojito/mojito/log"
	"github.com/go-mojito/mojito/pkg/router"
	"github.com/infinytum/injector"
	"github.com/infinytum/structures"
)

var (
	livenessChecks  = structures.NewMap[string, *Health]()
	readinessChecks = structures.NewMap[string, *Health]()
)

// IsHealthy will return true if both liveness and readiness checks succeed
func IsHealthy() bool {
	return IsLive() && IsReady()
}

// IsLive will return true if all liveness checks succeed
func IsLive() bool {
	for _, err := range livenessChecks.ToMap() {
		if err != nil && !err.Healthy {
			return false
		}
	}
	return true
}

// IsReady will return true if all readiness checks succeed
func IsReady() bool {
	for _, err := range readinessChecks.ToMap() {
		if err != nil && !err.Healthy {
			return false
		}
	}
	return true
}

// Liveness registers a new readiness check
func Liveness(name string, check interface{}, cycle time.Duration) error {
	checkType := reflect.TypeOf(check)
	if checkType.Kind() != reflect.Func {
		return ErrorHealthCheckNotAFunc
	}

	if checkType.NumOut() != 1 || checkType.Out(0) != reflect.TypeOf(&Health{}) {
		return ErrorHealthCheckNoError
	}
	go checkExecutor(name, check, cycle, livenessChecks)
	return livenessChecks.Set(name, nil)
}

// Readiness registers a new readiness check
func Readiness(name string, check interface{}, cycle time.Duration) error {
	checkType := reflect.TypeOf(check)
	if checkType.Kind() != reflect.Func {
		return ErrorHealthCheckNotAFunc
	}

	if checkType.NumOut() != 1 || checkType.Out(0) != reflect.TypeOf(&Health{}) {
		return ErrorHealthCheckNoError
	}
	go checkExecutor(name, check, cycle, readinessChecks)
	return readinessChecks.Set(name, nil)
}

// OnDefault registers the healthz endpoint on the default router
func OnDefault() {
	On(mojito.DefaultRouter())
}

// On registers the healthz endpoint on the given router
func On(r mojito.Router) {
	r.WithGroup("/healthz", func(group router.Group) {
		group.GET("/", handleHealthz)
		group.GET("/live", handleLiveness)
		group.GET("/live/:name", handleLivenessCheck)
		group.GET("/ready", handleReadiness)
		group.GET("/ready/:name", handleReadinessCheck)
	})
}

func checkExecutor(name string, check interface{}, cycle time.Duration, table structures.Map[string, *Health]) {
	for {
		healthErr, err := injector.CallT[*Health](check)
		if err != nil {
			log.Error(err)
			return
		}
		table.Set(name, healthErr)
		<-time.After(cycle)
	}
}
