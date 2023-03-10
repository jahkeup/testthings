package must

import (
	"encoding/json"
	"fmt"

	"github.com/jahkeup/testthings"
)

// FromJSON will parse the given json text into an output object. If the object
// does not parse, the test will be failed. Alternatively, this function may be
// used in initialization contexts by passing nil for the testingT object in
// which case errors will cause the program to panic at runtime.
func FromJSON[T any](testingT testthings.Terminator, text []byte) T {
	if th, ok := testingT.(interface {
		Helper()
	}); ok {
		th.Helper()
	}

	return Must[T](testingT, func() T {
		out := new(T)
		if err := json.Unmarshal(text, out); err != nil {
			msg := fmt.Sprintf("cannot unmarshal json into %T: %v", out, err)
			if testingT == nil {
				panic(msg)
			}
			testingT.Fatal(msg)
		}

		return *out
	})
}
