package log

import (
	"github.com/sirupsen/logrus"
	"os"
	"posam/config"
)

var (
	log            = logrus.New()
	ENV_PRODUCTION = config.GetBool("general.production")
)

func init() {
	if ENV_PRODUCTION {
		file, err := os.OpenFile("log.json", os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			log.Out = file
		} else {
			log.Info("Failed to log to file, using default stderr")
		}
		log.Formatter = &logrus.JSONFormatter{}
		log.Level = logrus.InfoLevel
	} else {
		log.Out = os.Stdout
		log.Level = logrus.TraceLevel
		log.ReportCaller = true
		log.Formatter = &logrus.TextFormatter{
			FullTimestamp: true,
		}
	}
}

type Fields logrus.Fields

func WithFields(fields Fields) *logrus.Entry {
	return log.WithFields(logrus.Fields(fields))
}

func Data(key string, value interface{}) *logrus.Entry {
	return log.WithField(key, fmt.Sprintf("%#v", value))
}

func Trace(args ...interface{}) {
	log.Trace(args)
}

func Tracef(format string, args ...interface{}) {
	log.Tracef(format, args)
}

func Debug(args ...interface{}) {
	log.Debug(args)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args)
}

func Info(args ...interface{}) {
	log.Info(args)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args)
}

func Warn(args ...interface{}) {
	log.Warn(args)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args)
}

func Error(args ...interface{}) {
	log.Error(args)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args)
}

func Fatal(args ...interface{}) {
	log.Fatal(args)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args)
}

func Panic(args ...interface{}) {
	log.Panic(args)
}

func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args)
}
