package witai

const (
	DefaultDevice = "default"
)

type Verbosity uint

const (
	_ Verbosity = iota
	Error
	Warn
	Info
	Debug
)
