package amplitude

type Logger interface {
	Debug(message string)
	Info(message string)
	Warn(message string)
	Error(message string)
}
