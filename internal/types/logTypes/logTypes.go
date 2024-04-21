package logTypes

type LogType string

const (
	DEBUG LogType = "debug"
	INFO  LogType = "info"
	WARN  LogType = "warn"
	ERROR LogType = "error"
	PANIC LogType = "panic"
	FATAL LogType = "fatal"
)
