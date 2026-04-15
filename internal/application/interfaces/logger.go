package interfaces

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
	Command(msg string)
	Success(msg string)
}
