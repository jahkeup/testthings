package must_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/jahkeup/testthings/must"
	"github.com/jahkeup/testthings/testerr"
)

func TestMust(t *testing.T) {
	t.Run("func() T", func(t *testing.T) {
		f := func() int { return 42 }
		actual := must.Must[int](t, f)
		if actual != 42 {
			t.Fatal("not passing right")
		}
	})

	t.Run("func(*T)", func(t *testing.T) {
		f := func(v *int) { *v = 42 }
		actual := must.Must[int](t, f)
		if actual != 42 {
			t.Fatal("not passing right")
		}
	})

	t.Run("panic", func(t *testing.T) {
		t.Run("func() (T, error)", func(t *testing.T) {
			f := func() (int, error) {
				return 42, testerr.Expected
			}

			wg := &sync.WaitGroup{}
			defer wg.Wait()

			wg.Add(1)
			go func() {
				defer func() {
					if r := recover(); r == nil {
						t.Error("did not panic")
					}
					wg.Done()
				}()
				var actual int

				actual = must.Must[int](nil, f)
				if actual == 42 {
					t.Error("shouldn't have even gotten here")
				}
			}()
		})

		t.Run("func() T", func(t *testing.T) {
			f := func() int {
				panic("bai")
				return 42 // shouldn't be returned
			}

			wg := &sync.WaitGroup{}
			defer wg.Wait()

			wg.Add(1)
			go func() {
				defer func() {
					if r := recover(); r == nil {
						t.Error("did not panic")
					}
					wg.Done()
				}()
				var actual int

				actual = must.Must[int](nil, f)
				if actual == 42 {
					t.Error("shouldn't have even gotten here")
				}
			}()
		})
	})
}

type Foo struct {
	Name           string `json:"name,omitempty"`
	FavoriteNumber int    `json:"favorite_number,omitempty"`
}

func NewFoo(name string, num int) (*Foo, error) {
	// pretend condition to justify including an 'error' in the signature.
	if num%3 == 0 && num%5 == 0 {
		return nil, errors.New("better luck next time")
	}
	return &Foo{Name: name, FavoriteNumber: num}, nil
}

func TestNewFooing(t *testing.T) {
	const (
		errorNum   = 630 // 42 * 3 * 5
		noErrorNum = 631 // + 1 ;)
	)
	testcases := map[string]struct {
		foo            Foo
		expectedNumber int
	}{
		"happy path": {
			foo: must.Must[Foo](t, func() (*Foo, error) {
				return NewFoo("Foo The Great", noErrorNum)
			}),
			expectedNumber: 631,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			if tc.foo.FavoriteNumber != tc.expectedNumber {
				t.Fatalf("should have been %d but was %d", tc.expectedNumber, tc.foo.FavoriteNumber)
			}
		})
	}

	t.Run("foo panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic")
			}
		}()

		// TODO: keep trying to remove the type parameter from this.. Go will
		// get there, eventually.
		var foo Foo = must.Must[Foo](nil, func() (*Foo, error) {
			return NewFoo("im going to make Must panic", errorNum)
		})

		t.Fatal("should not get to here")
		foo.FavoriteNumber -= 1
	})
}
