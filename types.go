package testthings

type Cleanuper interface {
	Cleanup(func())
}

type Logger interface {
	Log(args ...any)
}
