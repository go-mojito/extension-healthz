package healthz

import (
	"time"

	"github.com/go-mojito/mojito"
)

func handleHealthz(ctx mojito.Context) error {
	isHealthy := IsHealthy()
	if !isHealthy {
		ctx.Response().WriteHeader(500)
	}
	return ctx.PrettyJSON(struct {
		Healthy   bool               `json:"healthy"`
		Liveness  map[string]*Health `json:"liveness"`
		Readiness map[string]*Health `json:"readiness"`
		Timestamp time.Time          `json:"timestamp"`
	}{
		Healthy:   isHealthy,
		Liveness:  livenessChecks.ToMap(),
		Readiness: readinessChecks.ToMap(),
		Timestamp: time.Now(),
	})
}

func handleLiveness(ctx mojito.Context) error {
	isHealthy := IsLive()
	if !isHealthy {
		ctx.Response().WriteHeader(500)
	}
	return ctx.PrettyJSON(struct {
		Healthy   bool               `json:"healthy"`
		Checks    map[string]*Health `json:"checks"`
		Timestamp time.Time          `json:"timestamp"`
	}{
		Healthy:   isHealthy,
		Checks:    livenessChecks.ToMap(),
		Timestamp: time.Now(),
	})
}

func handleReadiness(ctx mojito.Context) error {
	isHealthy := IsReady()
	if !isHealthy {
		ctx.Response().WriteHeader(500)
	}
	return ctx.PrettyJSON(struct {
		Healthy   bool               `json:"healthy"`
		Checks    map[string]*Health `json:"checks"`
		Timestamp time.Time          `json:"timestamp"`
	}{
		Healthy:   isHealthy,
		Checks:    readinessChecks.ToMap(),
		Timestamp: time.Now(),
	})
}

func handleLivenessCheck(ctx mojito.Context) error {
	check := livenessChecks.GetOrDefault(ctx.Request().Param("name"), nil)
	if check == nil {
		ctx.Response().WriteHeader(404)
		return nil
	}

	if !check.Healthy {
		ctx.Response().WriteHeader(500)
	}
	return ctx.PrettyJSON(check)
}

func handleReadinessCheck(ctx mojito.Context) error {
	check := readinessChecks.GetOrDefault(ctx.Request().Param("name"), nil)
	if check == nil {
		ctx.Response().WriteHeader(404)
		return nil
	}

	if !check.Healthy {
		ctx.Response().WriteHeader(500)
	}
	return ctx.PrettyJSON(check)
}
