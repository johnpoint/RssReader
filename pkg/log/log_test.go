package log

import (
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestLogger_Error(t *testing.T) {
	Error("test", String("info", "test"))
}

func TestLogger_Info(t *testing.T) {
	Info("test", String("info", "test"))
}

func TestOverrideLoggerWithOption(t *testing.T) {
	OverrideLoggerWithOption(map[string]interface{}{
		"service-name": "test_service",
	}, WithJSONEncoding(), WrapLevelEncoder(zapcore.CapitalLevelEncoder))
	Info("test", String("info", "test"))
	Error("test", String("info", "test"))
}
