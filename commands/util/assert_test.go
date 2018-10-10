package util

import (
	"testing"
)

func TestIsEqual(t *testing.T) {
	testIsEqual(t, 5, "5", true)
	testIsEqual(t, 23, "2343", false)
	testIsEqual(t, "teststr", "teststr", true)
	testIsEqual(t, "teststr", "different", false)
	testIsEqual(t, 234.645, "234.645", true)
	testIsEqual(t, 234.64, "234.645", false)
}

func TestIsGreater(t *testing.T) {
	testIsGt(t, 5, "5", false)
	testIsGt(t, 23, "2343", false)
	testIsGt(t, 23, "21", true)
	testIsGt(t, "teststr", "teststr", false)
	testIsGt(t, "teststr", "different", true)
	testIsGt(t, "teststr", "zdifferent", false)
	testIsGt(t, 234.645, "234.645", false)
	testIsGt(t, 234.64, "234.645", false)
	testIsGt(t, 234.646, "234.645", true)
}

func TestIsGreaterOrEqual(t *testing.T) {
	testIsGte(t, 5, "5", true)
	testIsGte(t, 23, "2343", false)
	testIsGte(t, 23, "21", true)
	testIsGte(t, "teststr", "teststr", true)
	testIsGte(t, "teststr", "different", true)
	testIsGte(t, "teststr", "zdifferent", false)
	testIsGte(t, 234.645, "234.645", true)
	testIsGte(t, 234.64, "234.645", false)
	testIsGte(t, 234.646, "234.645", true)
}

func TestIsLessor(t *testing.T) {
	testIsLt(t, 5, "5", false)
	testIsLt(t, 23, "2343", true)
	testIsLt(t, 23, "21", false)
	testIsLt(t, "teststr", "teststr", false)
	testIsLt(t, "teststr", "different", false)
	testIsLt(t, "teststr", "zdifferent", true)
	testIsLt(t, 234.645, "234.645", false)
	testIsLt(t, 234.64, "234.645", true)
	testIsLt(t, 234.646, "234.645", false)
}

func TestIsLessOrEqual(t *testing.T) {
	testIsLte(t, 5, "5", true)
	testIsLte(t, 23, "2343", true)
	testIsLte(t, 23, "21", false)
	testIsLte(t, "teststr", "teststr", true)
	testIsLte(t, "teststr", "different", false)
	testIsLte(t, "teststr", "zdifferent", true)
	testIsLte(t, 234.645, "234.645", true)
	testIsLte(t, 234.64, "234.645", true)
	testIsLte(t, 234.646, "234.645", false)
}

func TestIsDateAssert(t *testing.T) {
	testDateParsed(t, "2016-05-12")
	testDateParsed(t, "2016-05-12T15:03:23.123Z")
	testDateParsed(t, "2016-05-12T15:03:23Z")
	testDateParsed(t, "2017-09-19T11:23:12.9674575-04:00")
}

func TestIsNotDateAssert(t *testing.T) {
	testDateNotParsed(t, "2016-05-12T15:03:23")
	testDateNotParsed(t, "3:15PM")
	testDateNotParsed(t, "2016-05-12 3:15PM")
	testDateNotParsed(t, "2016-05-12T15:03:23.123Z -500")
}

func testDateParsed(t *testing.T, d string) {
	if err := isDate(d); err != nil {
		t.Errorf(err.Error())
	}
}

func testDateNotParsed(t *testing.T, d string) {
	if err := isDate(d); err == nil {
		t.Error("Date parsed unexpectedly: " + d)
	}
}

func testIsEqual(t *testing.T, i interface{}, value string, success bool) {
	err := isEqual(i, value)
	if err != nil && success {
		t.Errorf("IsEqual failed unexpectedly: %s", err.Error())
	} else if err == nil && !success {
		t.Errorf("IsEqual unexpected succeeded: %v==%s", i, value)
	}
}

func testIsGt(t *testing.T, i interface{}, value string, success bool) {
	err := isGt(i, value)
	if err != nil && success {
		t.Errorf("IsGt failed unexpectedly: %s", err.Error())
	} else if err == nil && !success {
		t.Errorf("IsGt unexpected succeeded: %v!<%s", i, value)
	}
}

func testIsGte(t *testing.T, i interface{}, value string, success bool) {
	err := isGte(i, value)
	if err != nil && success {
		t.Errorf("IsGte failed unexpectedly: %s", err.Error())
	} else if err == nil && !success {
		t.Errorf("IsGte unexpected succeeded: %v!<=%s", i, value)
	}
}

func testIsLt(t *testing.T, i interface{}, value string, success bool) {
	err := isLt(i, value)
	if err != nil && success {
		t.Errorf("IsLt failed unexpectedly: %s", err.Error())
	} else if err == nil && !success {
		t.Errorf("IsLt unexpected succeeded: %v!>%s", i, value)
	}
}

func testIsLte(t *testing.T, i interface{}, value string, success bool) {
	err := isLte(i, value)
	if err != nil && success {
		t.Errorf("IsLte failed unexpectedly: %s", err.Error())
	} else if err == nil && !success {
		t.Errorf("IsLte unexpected succeeded: %v!>=%s", i, value)
	}
}
