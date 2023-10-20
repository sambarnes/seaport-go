package logging

import "go.uber.org/zap"

var Log *zap.Logger

func init() {
	Log, _ = zap.NewProduction()
}
