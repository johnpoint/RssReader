package log

type Option interface {
	apply(*Logger)
}

type optionFunc func(*Logger)

func (f optionFunc) apply(log *Logger) {
	f(log)
}

// WithJSONEncoding 可以指定 json 编码方式
func WithJSONEncoding() Option {
	return optionFunc(func(log *Logger) {
		log.encoding = "json"
	})
}

// WithConsoleEncoding 可以指定 console 编码方式
func WithConsoleEncoding() Option {
	return optionFunc(func(log *Logger) {
		log.encoding = "console"
	})
}

// WrapLevelEncoder 自定义 level 编码方式
func WrapLevelEncoder(levelEncoder LevelEncoder) Option {
	return optionFunc(func(logger *Logger) {
		logger.levelEncoder = levelEncoder
	})
}

// WrapTimeEncoder 自定义 time 编码方式
func WrapTimeEncoder(timeEncoder TimeEncoder) Option {
	return optionFunc(func(logger *Logger) {
		logger.timeEncoder = timeEncoder
	})
}

func IsDev() Option {
	return optionFunc(func(logger *Logger) {
		logger.isDev = true
	})
}

// WarpOutputPaths 自定义输出
func WarpOutputPaths(paths []string) Option {
	return optionFunc(func(logger *Logger) {
		logger.outputPaths = paths
	})
}

// WarpErrorOutputPaths 自定义error信息输出
func WarpErrorOutputPaths(paths []string) Option {
	return optionFunc(func(logger *Logger) {
		logger.errorOutputPaths = paths
	})
}
