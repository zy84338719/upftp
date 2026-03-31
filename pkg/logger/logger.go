package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var std *zap.Logger

// LoggerConfig 日志配置结构体
type LoggerConfig struct {
	Level      string `yaml:"level"`       // 日志级别: debug, info, warn, error
	Format     string `yaml:"format"`      // 日志格式: text, json
	Output     string `yaml:"output"`      // 输出方式: console, file, both
	File       string `yaml:"file"`        // 日志文件路径
	MaxSize    int    `yaml:"max_size"`    // 最大文件大小(MB)
	MaxBackups int    `yaml:"max_backups"` // 最大备份数量
	MaxAge     int    `yaml:"max_age"`     // 最大保留天数
	Compress   bool   `yaml:"compress"`    // 是否压缩
}

// Init 初始化日志系统
func Init(level, format string) {
	config := &LoggerConfig{
		Level:      level,
		Format:     format,
		Output:     "console",
		File:       "./logs/upftp.log",
		MaxSize:    100,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   true,
	}

	initWithConfig(config)
}

// InitWithConfig 使用配置初始化日志系统
func InitWithConfig(config *LoggerConfig) {
	initWithConfig(config)
}

// initWithConfig 内部初始化函数
func initWithConfig(config *LoggerConfig) {
	var zapLevel zapcore.Level
	switch strings.ToLower(config.Level) {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn", "warning":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	var encoderConfig zapcore.EncoderConfig
	if config.Format == "json" {
		encoderConfig = zap.NewProductionEncoderConfig()
		encoderConfig.TimeKey = "timestamp"
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// 创建输出
	var cores []zapcore.Core

	// 控制台输出
	if config.Output == "console" || config.Output == "both" {
		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			zapLevel,
		)
		cores = append(cores, consoleCore)
	}

	// 文件输出
	if config.Output == "file" || config.Output == "both" {
		// 确保日志目录存在
		logDir := filepath.Dir(config.File)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			fmt.Printf("Failed to create log directory: %v\n", err)
		} else {
			// 使用 lumberjack 进行日志轮转
			lumberjackWriter := &lumberjack.Logger{
				Filename:   config.File,
				MaxSize:    config.MaxSize,
				MaxBackups: config.MaxBackups,
				MaxAge:     config.MaxAge,
				Compress:   config.Compress,
			}

			// 对于文件输出，使用不带颜色的编码器
			fileEncoderConfig := encoderConfig
			if config.Format != "json" {
				fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
			}

			fileCore := zapcore.NewCore(
				zapcore.NewConsoleEncoder(fileEncoderConfig),
				zapcore.AddSync(lumberjackWriter),
				zapLevel,
			)
			cores = append(cores, fileCore)
		}
	}

	// 如果没有配置输出，默认使用控制台输出
	if len(cores) == 0 {
		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			zapLevel,
		)
		cores = append(cores, consoleCore)
	}

	// 创建日志器
	core := zapcore.NewTee(cores...)
	std = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))

	Info("Logger initialized with level: %s, format: %s, output: %s", config.Level, config.Format, config.Output)
}

// Debug 输出调试级别日志
func Debug(format string, args ...interface{}) {
	if std != nil {
		std.Sugar().Debugf(format, args...)
	}
}

// Info 输出信息级别日志
func Info(format string, args ...interface{}) {
	if std != nil {
		std.Sugar().Infof(format, args...)
	}
}

// Warn 输出警告级别日志
func Warn(format string, args ...interface{}) {
	if std != nil {
		std.Sugar().Warnf(format, args...)
	}
}

// Error 输出错误级别日志
func Error(format string, args ...interface{}) {
	if std != nil {
		std.Sugar().Errorf(format, args...)
	}
}

// Fatal 输出致命级别日志并退出
func Fatal(format string, args ...interface{}) {
	if std != nil {
		std.Sugar().Fatalf(format, args...)
	}
	os.Exit(1)
}

// WithFields 添加字段并返回新的日志器
func WithFields(fields map[string]interface{}) *zap.SugaredLogger {
	if std == nil {
		Init("info", "text")
	}

	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}

	return std.With(zapFields...).Sugar()
}

// WithField 添加单个字段并返回新的日志器
func WithField(key string, value interface{}) *zap.SugaredLogger {
	if std == nil {
		Init("info", "text")
	}

	return std.With(zap.Any(key, value)).Sugar()
}
