package witai

const (
	// DefaultDevice points to the default recording device on your machine.
	DefaultDevice = "default"
)

// Verbosity controls the verbosity of wit_* commands' logging levels.
type Verbosity uint

const (
	// Ignore the first verbosity level.
	_ Verbosity = iota

	Error // 1
	Warn  // 2
	Info  // 3
	Debug // 4
)
