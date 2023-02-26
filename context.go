package testthings

import "context"

// C provides a context cleaned up with tests. A convenience function.
func C(testingT Cleanuper) context.Context {
	// callers wanting short and sweet don't get the cancel function :P
	ctx, _ := NewContext(testingT)
	return ctx
}

// NewContext creates a context that's cancelled when the testing.TB scope ends.
func NewContext(testingT Cleanuper) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	testingT.Cleanup(cancel)
	return ctx, cancel
}
