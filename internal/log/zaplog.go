package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
)

type log interface {
	debug(string, ...any)
	info(string, ...any)
	error(string, ...any)
	sync() error
}

type customizeLog struct {
	log  *zap.Logger
	std  zapcore.WriteSyncer
	file string
}

var logger log

func init() {
	var cl = new(customizeLog)
	cl.newLog()
	logger = cl
}

func (c *customizeLog) createDir() {
	//path, _ := os.Executable()
	//dir := filepath.Dir(path)
	logDir := "/tmp/log"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, 0755) // 创建目录及其父目录，权限为 0755
		if err != nil {
			panic(fmt.Sprintf("Catalog creation failed: %v\n", err))
			return
		}
	}
	c.file = filepath.Join(logDir, "detection.log")
	return
}

func (c *customizeLog) fileRecorder() {
	c.std = zapcore.AddSync(&lumberjack.Logger{
		Filename:   c.file,
		MaxSize:    50, // megabytes
		MaxBackups: 1,
		MaxAge:     7, // days
	})
	return
}

func (c *customizeLog) newLog() {
	c.createDir()
	c.fileRecorder()
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, os.Stdout, zap.DebugLevel),
		//zapcore.NewCore(consoleEncoder, c.std, zap.DebugLevel),
	)
	c.log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))
}

func (c *customizeLog) debug(msg string, args ...any) {
	c.log.Debug(fmt.Sprintf(msg, args...))
}
func (c *customizeLog) info(msg string, args ...any) {
	c.log.Info(fmt.Sprintf(msg, args...))
}
func (c *customizeLog) error(msg string, args ...any) {
	c.log.Error(fmt.Sprintf(msg, args...))
}
func (c *customizeLog) sync() error {
	return c.log.Sync()
}

func Sync() {
	_ = logger.sync()
}

func Debug(format string, args ...any) {
	logger.debug(format, args...)
}

func Info(format string, args ...any) {
	logger.info(format, args...)
}

func Error(format string, args ...any) {
	logger.error(format, args...)
}
