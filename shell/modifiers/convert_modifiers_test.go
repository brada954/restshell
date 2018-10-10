package modifiers

import (
	"testing"
)

func TestLengthModifier(t *testing.T) {
	testLengthMod(t, "testthis", 8)
	testLengthMod(t, "234", 3)
	testLengthMod(t, 234, 3)
	testLengthMod(t, "", 0)
	testLengthMod(t, nil, 0)
	testLengthMod(t, "001232", 6)
	testLengthMod(t, 01232, 3) // this is an octal number (Decimal 666)
	testLengthMod(t, -43223, 6)
}

func testLengthMod(t *testing.T, value interface{}, expected int) {
	valueModifierFunc := NullModifier
	valueModifierFunc = MakeLengthModifier(valueModifierFunc)

	result, err := valueModifierFunc(value)
	if value == nil {
		if err == nil {
			t.Errorf("Expected an error for a nil value")
		}
		if result != nil {
			t.Errorf("Expected nil value for a nil value")
		}
		return
	}

	if err != nil {
		t.Errorf("Length modifier error %d!=%v: %s", expected, result, err.Error())
	}

	// Consider having a 0 value for nil??? TBD
	if newval, ok := result.(int); !ok {
		t.Errorf("Length Modifier did not return a length for %v", value)
	} else if newval != expected {
		t.Errorf("Length Modifier did not return expected length: %d!=%d", expected, newval)
	}
}

func TestConvertToIntModifier(t *testing.T) {
	testConvertToIntMod(t, "324", 324)
	testConvertToIntMod(t, "-123", -123)
	testConvertToIntMod(t, 54.87, 54)
	testConvertToIntMod(t, 54.27, 54)
	testConvertToIntMod(t, -24.27, -24)
	testConvertToIntMod(t, -24.97, -24)
	testConvertToIntMod(t, "0", 0)

	// These should probably be non-errors
	testConvertToIntModError(t, "677.143", "strconv.Atoi: parsing \"677.143\": invalid syntax")
	testConvertToIntModError(t, "-37.343", "strconv.Atoi: parsing \"-37.343\": invalid syntax")

	testConvertToIntModError(t, "", "strconv.Atoi: parsing \"\": invalid syntax")
}

func testConvertToIntMod(t *testing.T, value interface{}, expected int) {
	valueModifierFunc := NullModifier
	valueModifierFunc = MakeToIntModifier(valueModifierFunc)

	result, err := valueModifierFunc(value)
	if err != nil {
		t.Errorf("ConvertToInt error processing (%v): %s", value, err.Error())
	}

	// Consider having a 0 value for nil??? TBD
	if newval, ok := result.(int); !ok {
		if bigint, ok2 := result.(int64); !ok2 {
			t.Errorf("ConvertToInt did not return an integer for %v", value)
		} else if bigint != int64(expected) {
			t.Errorf("ConvertToInt did not return expected int64: %d!=%d", expected, newval)
		}
	} else if newval != expected {
		t.Errorf("ConvertToInt did not return expected integer: %d!=%d", expected, newval)
	}
}

func testConvertToIntModError(t *testing.T, value interface{}, expected string) {
	valueModifierFunc := NullModifier
	valueModifierFunc = MakeToIntModifier(valueModifierFunc)

	_, err := valueModifierFunc(value)
	if err == nil {
		t.Errorf("ConvertToInt unexpected success converting (%v)", value)
	}

	// Consider having a 0 value for nil??? TBD
	if err.Error() != expected {
		t.Errorf("ConvertToInt returned unexpected error: %s!=%d", expected, err.Error())
	}
}

func TestConvertToFloatModifier(t *testing.T) {
	testConvertToFloatMod(t, "324", 324.0)
	testConvertToFloatMod(t, "-123", -123.0)
	testConvertToFloatMod(t, 54.87, 54.87)
	testConvertToFloatMod(t, 54.27, 54.27)
	testConvertToFloatMod(t, -24.27, -24.27)
	testConvertToFloatMod(t, -24.97, -24.97)
	testConvertToFloatMod(t, "0", 0.0)

	// These should probably be non-errors
	testConvertToFloatModError(t, "", "strconv.Atoi: parsing \"\": invalid syntax")
}

func testConvertToFloatMod(t *testing.T, value interface{}, expected float64) {
	valueModifierFunc := NullModifier
	valueModifierFunc = MakeToFloatModifier(valueModifierFunc)

	result, err := valueModifierFunc(value)
	if err != nil {
		t.Errorf("ConvertToFloat error processing (%v): %s", value, err.Error())
	}

	// Consider having a 0 value for nil??? TBD
	if newval, ok := result.(float64); !ok {
		t.Errorf("ConvertToFloat did not return an float64 for %v", value)
	} else if newval != expected {
		t.Errorf("ConvertToFloat did not return expected integer: %f!=%f", expected, newval)
	}
}

func testConvertToFloatModError(t *testing.T, value interface{}, expected string) {
	valueModifierFunc := NullModifier
	valueModifierFunc = MakeToIntModifier(valueModifierFunc)

	_, err := valueModifierFunc(value)
	if err == nil {
		t.Errorf("ConvertToFloat unexpected success converting (%v)", value)
	}

	// Consider having a 0 value for nil??? TBD
	if err.Error() != expected {
		t.Errorf("ConvertToFloat returned unexpected error message: %s!=%s", expected, err.Error())
	}
}
