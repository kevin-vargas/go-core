package logrus

import (
	"github.com/kevin-vargas/go-core/log"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

type logger struct {
	client *logrus.Logger
}

func (l *logger) Info(args ...any) {
	l.client.Infoln(args...)
}
func (l *logger) Infof(format string, args ...any) {
	l.client.Infof(format, args...)
}
func (l *logger) Error(err error) {
	l.client.Errorln(err.Error())

}
func (l *logger) Debug(args ...any) {
	l.client.Debug(args...)
}
func (l *logger) Debugf(format string, args ...any) {
	l.client.Debugf(format, args...)
}

var mapLevel = map[log.Level]logrus.Level{
	log.Debug: logrus.InfoLevel,
	log.Error: logrus.ErrorLevel,
	log.Info:  logrus.InfoLevel,
}

func (l *logger) SetLevel(lvl log.Level) {
	l.client.SetLevel(mapLevel[lvl])
}

func New(ops ...log.Option) log.Logger {
	client := logrus.New()
	client.Formatter = &easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "\t[%lvl%]\t: %msg% \n",
	}
	var l log.Logger = &logger{
		client: client,
	}
	for _, o := range ops {
		o(l)
	}
	return l
}
