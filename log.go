package xxl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

// LogFunc 应用日志
type LogFunc func(req LogReq, res *LogRes) []byte

// Logger 系统日志
type Logger interface {
	Info(format string, a ...interface{})
	InfoJob(jobId int64, format string, a ...interface{})
	Error(format string, a ...interface{})
	ErrorJob(jobId int64, format string, a ...interface{})
	ReadLog(req *LogReq) *LogRes
}

func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{}
}
func NewDefaultLoggerByArgs(logPath string) *DefaultLogger {
	return &DefaultLogger{
		Path: logPath,
	}
}

type DefaultLogger struct {
	Path string // 日志路径
}

func (l *DefaultLogger) Init() {
	if len(l.Path) == 0 {
		l.Path = DefaultLogPath
	}
}

func (l *DefaultLogger) ReadLog(req *LogReq) *LogRes {
	logPath := fmt.Sprintf("%s/%d.log", l.Path, req.LogID)
	context, toLineNum, end := l.readFromLine(req.FromLineNum, DefaultLogPageSize, logPath)

	return &LogRes{Code: 200, Msg: "", Content: LogResContent{
		FromLineNum: req.FromLineNum,
		ToLineNum:   toLineNum,
		LogContent:  context,
		IsEnd:       end,
	}}
}

func (l *DefaultLogger) readFromLine(start, size int, fileName string) (string, int, bool) {
	file, _ := os.Open(fileName)
	fileScanner := bufio.NewScanner(file)
	lineCount := 1
	message := ""
	for fileScanner.Scan() {
		lineCount++
		if lineCount < start {
			continue
		}

		oneLine := fileScanner.Text()
		message += oneLine + "\n"
		if oneLine == TaskLogEnd {
			return message, lineCount, true
		}

		if start+size < lineCount {
			break
		}
	}

	defer file.Close()
	return message, lineCount, false
}
func (l *DefaultLogger) Info(format string, a ...interface{}) {
	//fmt.Println(fmt.Sprintf(format, a...))
}

func (l *DefaultLogger) Error(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(format, a...))
}

func (l *DefaultLogger) InfoJob(jobId int64, format string, a ...interface{}) {
	l.writeLog(jobId, "info", fmt.Sprintf(format, a...))
}

func (l *DefaultLogger) ErrorJob(jobId int64, format string, a ...interface{}) {
	l.writeLog(jobId, "info", fmt.Sprintf(format, a...))
}

func (l *DefaultLogger) writeLog(jobId int64, level, message string) {
	if !IsDir(l.Path) {
		CreateDir(l.Path)
	}

	logMsg := message
	if message != TaskLogStart && message != TaskLogEnd {
		logMsg = fmt.Sprintf("%s %s %s", time.Now().Format(TimestampFormatter), level, message)
	}

	logPath := fmt.Sprintf("%s/%d.log", l.Path, jobId)
	writeString := []byte(logMsg + "\n")
	err := l.appendLog(logPath, writeString)
	if err != nil {
		fmt.Println("日志写入失败 err=" + err.Error())
	}
}

func (l *DefaultLogger) appendLog(filename string, data []byte) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}
