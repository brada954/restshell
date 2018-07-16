package shell

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

var json1 = `{
	"testint" : 45,
	"teststr" : "my home town",
	"car" : {
	"color" : "red",
	"make" : "dodge",
	"strarray" : [ "s1", "s2", "s2" ]
	},
	"dataarray" : [ {"d1" : 1}, {"d2" : 2} ]
}`

var jsonArray = `[
	{ "name" : "object1", "value" : "first"},
	{ "name" : "object2", "value" : "second"}
]`

// Test the alternate support to test a string/non-JSON result
// from a REST service. The node is "/" for non-json results.
func TestNonJsonPushResult(t *testing.T) {
	expectedString := "This is a test string"
	err := errors.New(expectedString)

	PushError(err)

	if r, err := PeekResult(0); err != nil {
		t.Errorf("Unexpected failure to peek result: %s", err.Error())
		return
	} else {
		node, err := GetJsonNode("/", r.Map)
		if err != nil {
			t.Errorf("failure to get root value from map: %s", err.Error())
		} else if strResult, ok := node.(string); !ok {
			t.Errorf("Failed to get string result")
		} else if strResult != expectedString {
			t.Errorf("Failed to get string match: %s!=%s", expectedString, strResult)
		}
	}
}

func TestPushPeek(t *testing.T) {
	err := PushResponse(makeRestResponse(json1, "application/json", 200), nil)
	if err != nil {
		t.Errorf("Error pushing json: %s", err.Error())
		return
	}
	r, err := PeekResult(0)
	if err != nil {
		t.Errorf("Error peeking at the history")
	} else {
		if r.Text != string(json1) {
			t.Errorf("Text in peek result does not match expected for json1")
		}
		if r.Map == nil {
			t.Errorf("Map does not appear to be initialized")
		}
		if m, ok := r.Map.(map[string]interface{}); !ok {
			t.Errorf("Map has the wrong type: %v", reflect.TypeOf(m))
		}
	}
}

func TestRootAccessToMap(t *testing.T) {
	err := PushResponse(makeRestResponse(json1, "application/json", 200), nil)
	if err != nil {
		t.Errorf("Error pushing json: %s", err.Error())
		return
	}

	r, err := PeekResult(0)
	if err != nil {
		t.Errorf("PeakResult faild with error: " + err.Error())
	}

	n, err := GetJsonNode("/", r.Map)
	if err != nil {
		t.Errorf("Unexpected root err: %s", err.Error())
		return
	}

	// 4 nodes in map
	m, ok := n.(map[string]interface{})
	if !ok {
		t.Errorf("Unexpected error converting node to map; was type: %v", reflect.TypeOf(n))
	} else if len(m) != 4 {
		t.Errorf("Expected 4 nodes in map, got %d", len(m))
	}
}

func TestFirstLevelObjectFloat64Success(t *testing.T) {
	err := PushResponse(makeRestResponse(json1, "application/json", 200), nil)
	if err != nil {
		t.Errorf("Error pushing json: %s", err.Error())
		return
	}

	r, err := PeekResult(0)
	if err != nil {
		t.Errorf("PeakResult faild with error: " + err.Error())
	}

	assertNodeFloat64(t, "testint", r.Map, float64(45))
}

func TestSecondLevelObjectStringSuccess(t *testing.T) {
	err := PushResponse(makeRestResponse(json1, "application/json", 200), nil)
	if err != nil {
		t.Errorf("Error pushing json: %s", err.Error())
		return
	}

	r, err := PeekResult(0)
	if err != nil {
		t.Errorf("PeakResult faild with error: " + err.Error())
	}

	assertNodeString(t, "car.make", r.Map, "dodge")
}

func TestStringArraySuccess(t *testing.T) {
	err := PushResponse(makeRestResponse(json1, "application/json", 200), nil)
	if err != nil {
		t.Errorf("Error pushing json: %s", err.Error())
		return
	}

	r, err := PeekResult(0)
	if err != nil {
		t.Errorf("PeakResult faild with error: " + err.Error())
	}

	assertNodeString(t, "car.strarray[1]", r.Map, "s2")
}

func TestMissingObjectStringNodes(t *testing.T) {
	err := PushResponse(makeRestResponse(json1, "application/json", 200), nil)
	if err != nil {
		t.Errorf("Error pushing json: %s", err.Error())
		return
	}

	r, err := PeekResult(0)
	if err != nil {
		t.Errorf("PeakResult faild with error: " + err.Error())
	}

	assertNodeNotFound(t, "home", r.Map)
	assertNodeNotFound(t, "car.model", r.Map)
	assertNodeNotFound(t, "carx.make", r.Map)
}

func TestFirstItemOfObjectArray(t *testing.T) {
	err := PushResponse(makeRestResponse(json1, "application/json", 200), nil)
	if err != nil {
		t.Errorf("Error pushing json: %s", err.Error())
		return
	}

	r, err := PeekResult(0)
	if err != nil {
		t.Errorf("PeakResult faild with error: " + err.Error())
	}

	key := "dataarray[0].d1"
	x, err := GetJsonNode(key, r.Map)
	if err != nil {
		t.Errorf("Unexpected error for node (%s): %s", key, err.Error())
	}

	v, ok := x.(float64)
	if !ok {
		t.Errorf("Did not get an float64 type: got %v", reflect.TypeOf(x))
		return
	}

	if v != float64(1) {
		t.Errorf("Expected value 1; got %f", v)
	}
}

func TestInvalidKeyForItemInObjectArray(t *testing.T) {
	err := PushResponse(makeRestResponse(json1, "application/json", 200), nil)
	if err != nil {
		t.Errorf("Error pushing json: %s", err.Error())
		return
	}

	r, err := PeekResult(0)
	if err != nil {
		t.Errorf("PeakResult faild with error: " + err.Error())
	}

	key := "dataarray.d1"
	x, err := GetJsonNode(key, r.Map)
	if err == nil {
		t.Errorf("Unexpected success for node (%s) type returned: %v", key, reflect.TypeOf(x))
	}

	if !strings.Contains(err.Error(), "nexpected type") {
		t.Errorf("Unexpected error message for invalid type: %s", err.Error())
	}
}

func TestStringArrayOutOfBounds(t *testing.T) {
	err := PushResponse(makeRestResponse(json1, "application/json", 200), nil)
	if err != nil {
		t.Errorf("Error pushing json: %s", err.Error())
		return
	}

	r, err := PeekResult(0)
	if err != nil {
		t.Errorf("PeakResult faild with error: " + err.Error())
	}

	key := "car.strarray[3]"
	x, err := GetJsonNode(key, r.Map)
	if err == nil {
		t.Errorf("Unexpected success for node (%s) type returned: %v", key, reflect.TypeOf(x))
	}

	if !strings.Contains(err.Error(), "out of bounds") {
		t.Errorf("Unexpected error message for invalid type: %s", err.Error())
	}
}

func TestDataArrayOutOfBounds(t *testing.T) {
	err := PushResponse(makeRestResponse(json1, "application/json", 200), nil)
	if err != nil {
		t.Errorf("Error pushing json: %s", err.Error())
		return
	}

	r, err := PeekResult(0)
	if err != nil {
		t.Errorf("PeakResult faild with error: " + err.Error())
	}

	key := "dataarray[2].d1"
	x, err := GetJsonNode(key, r.Map)
	if err == nil {
		t.Errorf("Unexpected success for node (%s) type returned: %v", key, reflect.TypeOf(x))
	}

	if !strings.Contains(err.Error(), "out of bounds") {
		t.Errorf("Unexpected error message for invalid type: %s", err.Error())
	}
}

func TestJSONArray(t *testing.T) {
	err := PushResponse(makeRestResponse(jsonArray, "application/json", 200), nil)
	if err != nil {
		t.Errorf("Error pushing json: %s", err.Error())
		return
	}

	_, err = PeekResult(0)
	if err != nil {
		t.Errorf("PeakResult faild with error: " + err.Error())
	}
	//
	// AssertNodeNotFound(t, "home", r.Map)
	// AssertNodeNotFound(t, "car.model", r.Map)
	// AssertNodeNotFound(t, "carx.make", r.Map)
}

func makeRestResponse(data string, contentType string, status int) *RestResponse {
	result := RestResponse{}
	result.Text = data
	result.httpResp = &http.Response{
		Status:     strconv.Itoa(status) + "Some Error",
		StatusCode: status,
		Header:     http.Header(make(map[string][]string, 0)),
	}
	result.httpResp.Header.Add("Content-Type", contentType)
	return &result
}

func assertNodeNotFound(t *testing.T, key string, m interface{}) {
	x, err := GetJsonNode(key, m)
	if err == nil {
		t.Errorf("Unexpected success for node (%s)", key)
	}
	if x != nil {
		t.Errorf("Unexpected value found for node (%s): %v(%v)", key, reflect.TypeOf(x), x)
	}
	if err != ErrNotFound {
		t.Errorf("Unexpected error value for not found node (%s): %s", key, err.Error())
	}
}

func assertNodeFloat64(t *testing.T, key string, m interface{}, value float64) {
	x, err := GetJsonNode(key, m)
	if err != nil {
		t.Errorf("Unexpected error for node (%s): %s", key, err.Error())
	}
	switch v := x.(type) {
	case float64:
		if v != value {
			t.Errorf("Unexpected value for (%s) expected %f; got %f", key, value, v)
		}
	default:
		t.Errorf("Invalid type for (%s) expected a float64; got %v", key, reflect.TypeOf(x))
	}
}

func assertNodeString(t *testing.T, key string, m interface{}, value string) {
	x, err := GetJsonNode(key, m)
	if err != nil {
		t.Errorf("Unexpected error for node (%s): %s", key, err.Error())
	}
	switch v := x.(type) {
	case string:
		if v != value {
			t.Errorf("Unexpected value for (%s) expected %s; got %s", key, value, v)
		}
	default:
		t.Errorf("Invalid type for (%s) expected a string; got %v", key, reflect.TypeOf(x))
	}
}
