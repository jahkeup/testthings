package must

import (
	"fmt"

	"github.com/jahkeup/testthings"
	"github.com/jahkeup/testthings/testerr"
)

type Mustable[T any] interface {
	~func(*T) | ~func() T | ~func() *T | ~func() (T, error) | ~func() (*T, error)
}

func Must[T any, F Mustable[T]](testingT testthings.Terminator, mustable F) T {
	if th, ok := testingT.(interface {
		Helper()
	}); ok {
		th.Helper()
	}

	// ret is the final value to be returned. Declared here to allow the panic
	// recovery to clear the value.
	var ret T

	// When not hooking into testingT, set up to recover any panics that might
	// propagate from the inner function.
	if testingT != nil {
		recoverAsFatal := func() {
			if r := recover(); r != nil {
				testingT.Fatal(fmt.Sprintf("must! but:\n%v", r))
				// when recovering, ensure the type's zero value is returned
				var zero T
				ret = zero
			}
		}
		defer recoverAsFatal()
	}

	doFail := func(v any) {
		msg := fmt.Sprintf("must! but: %v", v)
		if testingT != nil {
			testingT.Fatal(msg)
		} else {
			panic(msg)
		}
	}

	// Unfortunately Go's generics don't eliminate the need to reflect on the
	// signatures, so we do that and invoke the right form. Nonetheless, its
	// nice to see call sites cleaned up by this :)
	switch fn := any(mustable).(type) {
	case func(*T): // injected argument to mutate and return
		// inject is used as dedicated working-memory space, to control the
		// indirection of the object given to callers.
		var inject *T
		inject = new(T)
		fn(inject)
		ret = *inject
		return ret

	case func() T: // constructor-like without errors (can panic)
		ret = fn()
		return ret
	case func() *T: // constructor-like without errors (can panic)
		retptr := fn()
		if retptr == nil {
			doFail(testerr.NilPointer)
		}
		return ret

	case func() (T, error): // constructor-like with error
		var err error
		ret, err = fn()
		if err != nil {
			doFail(err)
		}
		return ret
	case func() (*T, error): // constructor-like with error
		retptr, err := fn()
		if err != nil {
			doFail(err)
		}
		ret = *retptr
		return ret

	default:
		panic("uhhhhhh, nice. make an issue via GitHub please :)")
	}
}
