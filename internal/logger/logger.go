package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var (
	levelNames = map[Level]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		WARN:  "WARN",
		ERROR: "ERROR",
	}
	levelColors = map[Level]string{
		DEBUG: "\033[36m",
		INFO:  "\033[32m",
		WARN:  "\033[33m",
		ERROR: "\033[31m",
	}
)

type Logger struct {
	level  Level
	format string
	out    io.Writer
	prefix string
}

var std *Logger

func Init(level, format string) {
	l := parseLevel(level)
	std = &Logger{
		level:  l,
		format: format,
		out:    os.Stdout,
		prefix: "upftp",
	}
}

func parseLevel(s string) Level {
	switch strings.ToLower(s) {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn", "warning":
		return WARN
	case "error":
		return ERROR
	default:
		return INFO
	}
}

func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, args...)

	if l.format == "json" {
		log.Printf(`{"time":"%s","level":"%s","msg":"%s"}`, now, levelNames[level], msg)
	} else {
		color := levelColors[level]
		reset := "\033[0m"
		log.Printf("%s[%s]%s [%s] %s", color, levelNames[level], reset, now, msg)
	}
}

func Debug(format string, args ...interface{}) {
	if std != nil {
		std.log(DEBUG, format, args...)
	}
}

func Info(format string, args ...interface{}) {
	if std != nil {
		std.log(INFO, format, args...)
	}
}

func Warn(format string, args ...interface{}) {
	if std != nil {
		std.log(WARN, format, args...)
	}
}

func Error(format string, args ...interface{}) {
	if std != nil {
		std.log(ERROR, format, args...)
	}
}

func Fatal(format string, args ...interface{}) {
	if std != nil {
		std.log(ERROR, format, args...)
	}
	os.Exit(1)
}
