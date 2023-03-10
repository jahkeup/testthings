package testerr

type TestingError string

func (a TestingError) Error() string {
	if a != "" {
		return string(a)
	}
	return "testing error value"
}

var Ignore = TestingError("IGNORE THIS ERROR!")

var Expected = TestingError("this error is expected!")

var TODO = TestingError("TODO: an error")

var HACK = TestingError("HACK: this is a hack")

var Any TestingError

var NilPointer = TestingError("unexpected nil pointer")
