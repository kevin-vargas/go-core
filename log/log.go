package log

type Logger interface {
	Info(...any)
	Infof(string, ...any)
	Error(error)
	Debug(...any)
	Debugf(string, ...any)
	SetLevel(Level)
}

type Level int64

const (
	Undefined Level = iota
	Debug
	Info
	Error
)

type Option func(Logger)

func WithLevel(lvl Level) Option {
	return func(l Logger) {
		l.SetLevel(lvl)
	}
}
