package logger

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log = logrus.New()

func InitLogger(levelStr string) {
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FieldMap:        logrus.FieldMap{logrus.FieldKeyTime: "timestamp"},
	})

	writer := &lumberjack.Logger{
		Filename:   "app.log",
		MaxSize:    10, // MB
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
		LocalTime:  true,
	}

	Log.SetOutput(writer)
	Log.SetReportCaller(false)

	switch levelStr {
	case "debug":
		Log.SetLevel(logrus.DebugLevel)
	case "warn":
		Log.SetLevel(logrus.WarnLevel)
	case "error":
		Log.SetLevel(logrus.ErrorLevel)
	default:
		Log.SetLevel(logrus.InfoLevel)
	}

	Log.WithFields(logrus.Fields{"level": levelStr}).Info("Logger initialized")
}

func AuditLog(action string, userID uint, ip string, err error) {
	fields := logrus.Fields{
		"action":  action,
		"user_id": userID,
		"ip":      ip,
		"type":    "audit",
	}
	if err != nil {
		Log.WithFields(fields).WithError(err).Warn("Audit warning")
	} else {
		Log.WithFields(fields).Info("Audit success")
	}
}
