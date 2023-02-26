package testthings_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jahkeup/testthings"
)

type logged struct {
	Msgs []any
}

func (l *logged) Log(args ...any) {
	l.Msgs = append(l.Msgs, args...)
}

func TestLogKV(t *testing.T) {
	kv := testthings.KV{
		"foo":  "bar",
		"baz":  "qux",
		"test": "value",
	}
	logger := &logged{}
	testthings.LogKV(logger, kv)
	assert.Len(t, logger.Msgs, len(kv))
}

func TestLogKV_style(t *testing.T) {
	testthings.KV{
		"foo": "bar",
	}.Log(t)

	testthings.LogKV(t, testthings.KV{
		"foo": "bar",
	})
}

func TestFormatKV(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Empty(t, testthings.FormatKV("foo", testthings.KV{}))
	})

	t.Run("default", func(t *testing.T) {
		actual := testthings.FormatKV("", testthings.KV{
			"foo": "bar",
			"baz": "qux",
		})
		assert.NotEmpty(t, actual)
		assert.Equal(t, `baz="qux" foo="bar"`, actual)
	})

	t.Run("default", func(t *testing.T) {
		actual := testthings.FormatKV("", testthings.KV{
			"foo": "bar",
			"baz": "qux",
		})
		assert.NotEmpty(t, actual)
		assert.Equal(t, `baz="qux" foo="bar"`, actual)
	})

	t.Run("inferred sep", func(t *testing.T) {
		t.Run("comma", func(t *testing.T) {
			actual := testthings.FormatKV("%v=%q,", testthings.KV{
				"foo": "bar",
				"baz": "qux",
			})
			assert.NotEmpty(t, actual)
			assert.Equal(t, `baz="qux",foo="bar"`, actual)
		})
		t.Run("comma+space", func(t *testing.T) {
			actual := testthings.FormatKV("%v=%q, ", testthings.KV{
				"foo":  "bar",
				"baz":  "qux",
				"neat": "o",
			})
			assert.NotEmpty(t, actual)
			assert.Equal(t, `baz="qux", foo="bar", neat="o"`, actual)
		})

		t.Run("tab", func(t *testing.T) {
			actual := testthings.FormatKV("%v=%q\t", testthings.KV{
				"foo":  "bar",
				"baz":  "qux",
				"neat": "o",
			})
			assert.NotEmpty(t, actual)
			assert.Equal(t, `baz="qux"	foo="bar"	neat="o"`, actual)
		})
	})

	t.Run("style", func(t *testing.T) {
		actual := testthings.KV{
			"foo": "bar",
			"baz": 1,
		}.Format("%v%v")
		require.NotEmpty(t, actual)
		assert.Equal(t, "baz1foobar", actual)
	})

}
