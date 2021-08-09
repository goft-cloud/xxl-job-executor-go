package xxl

import (
	"fmt"
	"log"
)

// LogFunc 应用日志
type LogFunc func(req LogReq, res *LogRes) []byte

// Logger 系统日志
type Logger interface {
	Info(format string, a ...interface{})
	Error(format string, a ...interface{})
}

func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{}
}
func NewDefaultLoggerByArgs(logPath, runPath string) *DefaultLogger {
	return &DefaultLogger{
		Path:    logPath,
		RunPath: runPath,
	}
}

type DefaultLogger struct {
	Path    string // 日志路径
	RunPath string // 运行脚本路径
}

func (l *DefaultLogger) Init() {
	if len(l.Path) == 0 {
		l.Path = DefaultLogPath
	}

	if len(l.RunPath) == 0 {
		l.RunPath = DefaultRunPath
	}
}

func (l *DefaultLogger) Info(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(format, a...))
}

func (l *DefaultLogger) Error(format string, a ...interface{}) {
	log.Println(fmt.Sprintf(format, a...))
}
