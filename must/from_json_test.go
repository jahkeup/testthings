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

//go:embed testdata/some.json
var some_json []byte

var someJSONData = must.FromJSON[struct{ Hi []string }](nil, some_json)

var derivedData = must.Must[[]string](nil, func() []string {
	// assumption assertions, etc
	if len(someJSONData.Hi) == 0 {
		panic("no data?!")
	}
	return someJSONData.Hi
})

func TestFromJSON_simple(t *testing.T) {
	t.Logf("someJSONData: %#v", someJSONData)
	if len(someJSONData.Hi) < 1 {
		t.Error("should have loaded the json")
	}
	if len(someJSONData.Hi) != len(derivedData) {
		t.Fatal("should have derived the same data")
	}
	for i := range someJSONData.Hi {
		a, b := someJSONData.Hi[i], derivedData[i]
		if a != b {
			t.Errorf("expected element %d to be %q (a), but was %q (b)", i, a, b)
		}
	}
}
