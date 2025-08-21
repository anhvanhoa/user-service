package pkglog

import (
	serviceLogger "cms-server/domain/service/logger"
	pkgres "cms-server/infrastructure/service/response"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger interface {
	Log(c *fiber.Ctx, err error) error
	serviceLogger.Logger
}

type logger struct {
	Logger *zap.Logger
}

func NewConfig() *lumberjack.Logger {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", os.ModePerm)
	}
	date := time.Now().Format("2006-01-02")
	return &lumberjack.Logger{
		Filename:   fmt.Sprintf("logs/%s.log", date), // Tạo file theo ngày
		MaxSize:    10,                               // MB, giới hạn file log
		MaxBackups: 7,                                // Giữ lại 7 file log cũ
		MaxAge:     30,                               // Giữ log trong 30 ngày
		Compress:   true,                             // Nén log cũ
	}
}

// InitLogger thiết lập Logger với Lumberjack và Zap
func InitLogger(config *lumberjack.Logger, logLevel zapcore.Level, logFile bool) Logger {
	encoderConfig := zapcore.EncoderConfig{
		LevelKey:         "level",
		MessageKey:       "message",
		CallerKey:        "caller",
		TimeKey:          "time",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeTime:       zapcore.ISO8601TimeEncoder, // Format thời gian
		EncodeCaller:     zapcore.ShortCallerEncoder, // Hiển thị file.go:line
		EncodeLevel:      zapcore.CapitalLevelEncoder,
		ConsoleSeparator: " | ",
	}

	var coreLogs []zapcore.Core

	// Nếu log ra file thì thêm vào coreLogs
	if logFile {
		configFile := encoderConfig
		fileEncoder := zapcore.NewJSONEncoder(configFile)
		fileWriter := zapcore.AddSync(config)
		coreLogs = append(coreLogs, zapcore.NewCore(fileEncoder, fileWriter, logLevel))
	}

	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Hiển thị màu
	encoderConfig.StacktraceKey = "stack"
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleWriter := zapcore.Lock(os.Stdout)
	coreLogs = append(coreLogs, zapcore.NewCore(consoleEncoder, consoleWriter, logLevel))

	// Kết hợp nhiều writer
	core := zapcore.NewTee(coreLogs...)

	return &logger{
		Logger: zap.New(
			core,
			zap.AddCaller(),
			zap.AddCallerSkip(2),
			zap.AddStacktrace(zapcore.ErrorLevel),
		),
	}
}

func convertToZapFields(fields ...any) []zap.Field {
	var zapFields []zap.Field
	for i := 0; i < len(fields)-1; i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			continue // bỏ qua nếu key không phải string
		}
		value := fields[i+1]
		zapFields = append(zapFields, zap.Any(key, value))
	}
	return zapFields
}

// Các hàm tiện ích
func (l *logger) Info(msg string, fields ...any) {
	zapFields := convertToZapFields(fields...)
	if len(zapFields) > 0 {
		zapFields = append(zapFields, zap.String("info", msg))
	}
	l.Logger.Info(msg, zapFields...)
}

func (l *logger) Debug(msg string, fields ...any) {
	zapFields := convertToZapFields(fields...)
	if len(zapFields) > 0 {
		zapFields = append(zapFields, zap.String("debug", msg))
	}
	l.Logger.Debug(msg, zapFields...)
}

func (l *logger) Warn(msg string, fields ...any) {
	zapFields := convertToZapFields(fields...)
	if len(zapFields) > 0 {
		zapFields = append(zapFields, zap.String("warning", msg))
	}
	l.Logger.Warn(msg, zapFields...)
}

func (l *logger) Error(msg string, fields ...any) {
	zapFields := convertToZapFields(fields...)
	if len(zapFields) > 0 {
		zapFields = append(zapFields, zap.String("error", msg))
	}
	l.Logger.Error(msg, zapFields...)
}

func (l *logger) Fatal(msg string, fields ...any) {
	zapFields := convertToZapFields(fields...)
	if len(zapFields) > 0 {
		zapFields = append(zapFields, zap.String("fatal", msg))
	}
	l.Logger.Fatal(msg, zapFields...)
}

func (l *logger) Log(c *fiber.Ctx, err error) error {
	l.Error(
		err.Error(),
		zap.String("path", c.Path()),
		zap.String("method", c.Method()),
		zap.String("ip", c.IP()),
		zap.String("user-agent", c.Get("User-Agent")),
	)
	var er error
	switch e := err.(type) {
	case *pkgres.ErrorApp:
		er = c.Status(e.GetCode()).JSON(e)
	case *fiber.Error:
		er = c.Status(e.Code).JSON(fiber.Map{
			"Message": e.Message,
		})
	default:
		er = c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Message": e.Error(),
		})
	}
	return er
}
