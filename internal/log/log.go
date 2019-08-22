package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"synthesis/pkg/config"
)

var (
	log            = logrus.New()
	ENV_PRODUCTION = config.GetBool("general.production")
)

func init() {
	if ENV_PRODUCTION {
		file, err := os.OpenFile("log.json", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err == nil {
			log.Out = file
		} else {
			log.Info("Failed to log to file, using default stderr")
		}
		log.Formatter = &logrus.JSONFormatter{}
		log.Level = logrus.InfoLevel
	} else {
		//log.Out = os.Stdout
		//log.Level = logrus.TraceLevel
		////log.ReportCaller = true
		//log.Formatter = &logrus.TextFormatter{
		//FullTimestamp: true,
		//}
		file, err := os.OpenFile("log.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err == nil {
			log.Out = file
		} else {
			log.Info("Failed to log to file, using default stderr")
		}
		//log.Formatter = &logrus.JSONFormatter{
		//PrettyPrint: true,
		//}
		log.Formatter = &logrus.TextFormatter{
			DisableColors: true,
			//DisableLevelTruncation: false, // only work in color mode
		}
		log.Level = logrus.TraceLevel
	}
}

type M map[string]interface{}

func (m M) raw() M {
	for k, v := range m {
		m[k] = fmt.Sprintf("%v", v)
	}
	return m
}

func Vs(fields M) *logrus.Entry {
	for k, v := range fields {
		fields[k] = fmt.Sprintf("%#v", v)
	}
	return log.WithFields(logrus.Fields(fields))
}

func V(key string, value interface{}) *logrus.Entry {
	return log.WithField(key, fmt.Sprintf("%#v", value))
}

func T(args ...interface{}) {
	log.Trace(S(args...)...)
}

func Tv(fields M) {
	log.WithFields(logrus.Fields(fields.raw())).Trace()
}

func Tf(format string, args ...interface{}) {
	log.Tracef(format, S(args...)...)
}

func D(args ...interface{}) {
	log.Debug(S(args...)...)
}

func Dv(fields M) {
	log.WithFields(logrus.Fields(fields.raw())).Debug()
}

func Df(format string, args ...interface{}) {
	log.Debugf(format, S(args...)...)
}

func I(args ...interface{}) {
	log.Info(S(args...)...)
}

func Iv(fields M) {
	log.WithFields(logrus.Fields(fields.raw())).Info()
}

func If(format string, args ...interface{}) {
	log.Infof(format, S(args...)...)
}

func W(args ...interface{}) {
	log.Warn(S(args...)...)
}

func Wv(fields M) {
	log.WithFields(logrus.Fields(fields.raw())).Warn()
}

func Wf(format string, args ...interface{}) {
	log.Warnf(format, S(args...)...)
}

func E(args ...interface{}) {
	log.Error(S(args...)...)
}

func Ev(fields M) {
	log.WithFields(logrus.Fields(fields.raw())).Error()
}

func Ef(format string, args ...interface{}) {
	log.Errorf(format, S(args...)...)
}

func F(args ...interface{}) {
	log.Fatal(S(args...)...)
}

func Fv(fields M) {
	log.WithFields(logrus.Fields(fields.raw())).Fatal()
}

func Ff(format string, args ...interface{}) {
	log.Fatalf(format, S(args...)...)
}

func P(args ...interface{}) {
	log.Panic(S(args...)...)
}

func Pv(fields M) {
	log.WithFields(logrus.Fields(fields.raw())).Panic()
}

func Pf(format string, args ...interface{}) {
	log.Panicf(format, S(args...)...)
}

func S(args ...interface{}) []interface{} {
	r := []interface{}{}
	for _, v := range args {
		r = append(r, fmt.Sprintf("%v", v))
	}
	return r
}
