package testthings

// Cleanuper describe types that can cleanup after themselves and that allow
// adding functions to be called when they're cleaning up. In practice, this is
// really just a scoped-down testing.TB.
type Cleanuper interface {
	Cleanup(func())
}

// Logger describes types that can emit log messages. In practice, this is
// really just a scoped-down testing.TB.
type Logger interface {
	Log(args ...any)
}
