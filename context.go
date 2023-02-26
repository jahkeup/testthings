package testthings

import "context"

func C(testingT Cleanuper) context.Context {
	ctx, _ := NewContext(testingT)
	return ctx
}

func NewContext(testingT Cleanuper) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	testingT.Cleanup(cancel)
	return ctx, cancel
}
