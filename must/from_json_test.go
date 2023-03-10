package must_test

import (
	_ "embed"
	"reflect"
	"testing"

	"github.com/jahkeup/testthings/must"
)

type typeOne struct {
	One string
}

type typeTwo struct {
	Two [2]string
}

func TestMustFromJSON(t *testing.T) {
	testcases := map[string]struct {
		json     string
		expected any
	}{
		"type one": {
			json: `{"One": "1"}`,
			expected: typeOne{
				One: "1",
			},
		},
		"type two": {
			json: `{"Two": ["2", "3"]}`,
			expected: typeTwo{
				Two: [2]string{"2", "3"},
			},
		},
		"arstar": {
			json: `{"Three": ["4"]}`,
			expected: map[string]interface{}{
				"Three": []interface{}{"4"},
			},
		},
	}

	for name, tc := range testcases {
		jsonText := []byte(tc.json)

		t.Run(name, func(t *testing.T) {
			switch expected := tc.expected.(type) {
			case typeOne:
				actual := must.FromJSON[typeOne](t, jsonText)
				if actual != expected {
					t.Errorf("%#v != %#v", actual, expected)
				}

			case typeTwo:
				actual := must.FromJSON[typeTwo](t, jsonText)
				if actual != expected {
					t.Errorf("%#v != %#v", actual, expected)
				}
			default:
				actual := must.FromJSON[any](t, jsonText)
				if !reflect.DeepEqual(actual, expected) {
					t.Logf("actual:\t%#v", actual)
					t.Logf("expected:\t%#v", expected)
					t.Error("not equal")
				}
			}
		})
	}
}
