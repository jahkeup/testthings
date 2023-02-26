package testthings_test

import (
	"context"
	"testing"

	"github.com/jahkeup/testthings"
)

func TestContexts(t *testing.T) {
	topctx := testthings.C(t)
	t.Log("top context will remain alive")

	var subctx context.Context
	t.Run("sub", func(t *testing.T) {
		subctx = testthings.C(t)
		select {
		case _ = <-topctx.Done():
			t.Fatal("topctx should still be alive")
		case _ = <-subctx.Done():
			t.Fatal("subctx should be alive")
		default:
		}
	})

	t.Log("checking context")
	select {
	case _, ok := <-subctx.Done():
		if ok {
			t.Fatal("subctx should NOT be alive")
		}
	default:
	}

	select {
	case _ = <-topctx.Done():
		t.Fatal("topctx should still be alive")
	default:
	}

}
