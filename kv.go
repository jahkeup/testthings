package testthings

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// formatBasicKeyPair is used when no format is provided.
const formatBasicKeyPair = "%[1]v=%[2]q"

// KV is a convenience type to hold, format, and log contextual data. Helpers
// are available both at the package level and on the type itself.
type KV map[string]any

// Log prints the KV to the logger, one key-pair on each line.
func (k KV) Log(testingT Logger) {
	LogKV(testingT, k)
}

// Log prints the KV to the logger, one key-pair on each line formatted
// accordint to the given format string.
func (k KV) Logf(testingT Logger, format string) {
	logfKV(testingT, format, k)
}

// Format formats the entire KV into a string. See FormatKV for details on
// format strings are handled.
func (k KV) Format(format string) string {
	return FormatKV(format, k)
}

// Strings returns a list of strings where each key-pair has been formatted. The
// results are lexicographically sorted by their key.
func (k KV) Strings(format string) []string {
	return formatIntercepted(format, k, defaultInterceptor)
}

// LogKV logs one line per key-pair to provide contextual output within test
// cases.
func LogKV(testingT Logger, kv KV) {
	logfKV(testingT, formatBasicKeyPair, kv)
}

func logfKV(testingT Logger, format string, kv KV) {
	strs := formatIntercepted(format, kv, defaultInterceptor)

	for _, s := range strs {
		testingT.Log(s)
	}
}

// FormatKV produces a string with each key-pair formatted with the provided
// string. Trailing whitespace (and `,`) are treated as formatted string
// separators and is implicitly used to join the formatted strings together.
func FormatKV(kvFormat string, kv KV) string {
	var sep string
	if kvFormat == "" {
		kvFormat = formatBasicKeyPair
		sep = " "
	}

	trimmedFormat := strings.TrimRightFunc(kvFormat, func(r rune) bool {
		if unicode.IsSpace(r) {
			return true
		}

		switch r {
		case ',':
			return true
		}

		return false
	})
	if trimmedFormat != kvFormat {
		sep = kvFormat[len(trimmedFormat):]
	}

	strs := formatIntercepted(trimmedFormat, kv, defaultInterceptor)

	return strings.Join(strs, sep)
}

func formatIntercepted(kvFormat string, kv KV, intercept interceptorFactory) []string {
	keys := []string{}
	formatted := map[string]string{}
	for k, v := range kv {
		keys = append(keys, k)
		if intercept != nil {
			ik, iv := intercept(k), intercept(v)
			formatted[k] = fmt.Sprintf(kvFormat, ik, iv)
		} else {
			formatted[k] = fmt.Sprintf(kvFormat, k, v)
		}
	}

	sort.Strings(keys)
	strs := []string{}
	for _, k := range keys {
		strs = append(strs, formatted[k])
	}

	return strs
}

type interceptorFactory = func(v any) interceptorInstance

var defaultInterceptor = newInterceptor(nil, nil)

// newInterceptor creates an interception handler factory that is used to modify
// the string printed in formatting strings.
func newInterceptor(stringHandler func(any) string, gostringHandler func(any) string) interceptorFactory {
	return func(v any) interceptorInstance {
		return formatInterception{
			Value:           v,

			StringHandler:   stringHandler,
			GoStringHandler: gostringHandler,
		}
	}
}

type interceptorInstance interface {
	fmt.Stringer
	fmt.GoStringer

	StringV(any) string
	GoStringV(any) string
}

type formatInterception struct {
	Value any

	StringHandler   func(any) string
	GoStringHandler func(any) string
}

var _ interceptorInstance = (*formatInterception)(nil)

func (fi formatInterception) String() string {
	return fi.StringV(fi.Value)
}

func (fi formatInterception) GoString() string {
	return fi.GoStringV(fi.Value)
}

func (fi formatInterception) GoStringV(v any) string {
	if fi.GoStringHandler != nil {
		if s := fi.GoStringHandler(v); s != "" {
			return s
		}
	}

	return fmt.Sprintf("%#v", fi.Value)
}

func (fi formatInterception) StringV(v any) string {
	if fi.StringHandler != nil {
		if s := fi.StringHandler(v); s != "" {
			return s
		}
	}

	return fmt.Sprintf("%v", fi.Value)
}
