package nasync

const (
	LEVEL_UNKNOWN = iota
	LEVEL_DEBUG
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
)

type LogFunc func(level int, logStr string)
