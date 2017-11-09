package zap

import (
	"github.com/gliderlabs/com/objects"
	"go.uber.org/zap"
)

// TODO: wrap for With
// TODO: component for config

// Register the zap logger component with a registry
func Register(registry *objects.Registry) error {
	logger, _ := zap.NewDevelopment()
	return registry.Register(&objects.Object{Value: logger.Sugar()})
}
