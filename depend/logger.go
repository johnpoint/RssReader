package depend

import (
	"RssReader/config"
	"RssReader/pkg/bootstrap"
	"RssReader/pkg/log"
	"context"
	"go.uber.org/zap/zapcore"
)

type Logger struct{}

var _ bootstrap.Component = (*Logger)(nil)

func (d *Logger) Init(ctx context.Context) error {
	var options = make([]log.Option, 0)
	if config.Config.Environment == "local" {
		options = append(options,
			log.WithConsoleEncoding(),
			log.WrapLevelEncoder(zapcore.CapitalColorLevelEncoder),
			log.WrapTimeEncoder(zapcore.ISO8601TimeEncoder),
		)
	} else {
		options = append(options,
			log.WithJSONEncoding(),
			log.WrapLevelEncoder(zapcore.CapitalLevelEncoder),
		)
	}

	log.OverrideLoggerWithOption(map[string]interface{}{
		"service-name": config.Config.ServiceName,
	}, options...)
	return nil
}
