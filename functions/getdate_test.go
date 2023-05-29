package functions

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

var debugEnabled = false
var OutputDateFormat = "2006-01-02T15:04:05.000"

func TestSanityCheckLocalDateTimeOffsetAssumingEastCoast(t *testing.T) {
	offset := GetOffsetForLocal()
	if offset != -1*4*60*60 {
		t.Errorf("Expected 4 hours offset, but found %d", offset)
	}
}

func TestGetDateInLocalFormatMatchesNow(t *testing.T) {
	// Arrange
	dateNow := time.Now() // This is local time with offset from UTC
	dateFormat := "2006-01-02 15:04:05 MST"
	expectedDateString := dateNow.Format(dateFormat)

	// Act
	str, cache := GetDateSubstitute(nil, "getdate", "local", dateFormat)

	// Assert
	if date, ok := cache.(time.Time); ok {
		if !IsDateRangeWithinTollerance(date, dateNow, 1) {
			t.Errorf(
				"Returned date(%s) is not within tollerance of Now(%s)",
				date.Format(OutputDateFormat),
				dateNow.Format(OutputDateFormat),
			)
		}

		if str != expectedDateString {
			t.Errorf("Expected date formated as %s but got %s", expectedDateString, str)
		}
	} else {
		t.Errorf("Expected time.Time value returned but got: %v", reflect.TypeOf(cache))
	}
}

func TestGetDateDisplaysLocalFormatWithTimeZone(t *testing.T) {
	//Arrange
	dt := time.Date(2021, 3, 9, 17, 29, 32, 100000000, time.Local)
	expectedDateString := "2021-03-09T17:29:32.100"
	expectedTimeZone, _ := dt.Zone()
	expectedDateString = expectedDateString + " " + expectedTimeZone

	// Act
	str, date := GetDateSubstitute(dt, "getdate", "local", "2006-01-02T15:04:05.000 MST")

	// Assert
	if str != expectedDateString {
		t.Errorf("Expected date %s but got date %s", expectedDateString, str)
	}

	if dt != date {
		t.Errorf("Expected date value to remain unchanged")
	}
}

func TestGetDateLocalUnixMatchesUtcUnix(t *testing.T) {
	//Arrange
	dt := time.Date(2021, 3, 9, 17, 29, 32, 100000000, time.Local)
	dtUtc := dt.UTC()
	expectedTimeStamp := dtUtc.Unix()

	// Act
	str, date := GetDateSubstitute(dt, "getdate", "unix", "2006-01-02T15:04:05.000 MST")

	// Assert
	if fmt.Sprint(expectedTimeStamp) != str {
		t.Errorf("Expected timestamp %d but got date %s", expectedTimeStamp, str)
	}

	if dt != date {
		t.Errorf("Expected date value to remain unchanged")
	}
}

func TestGetDateUtcIsOffsetFromLocalNow(t *testing.T) {
	// Arrange
	dateFormat := "2006-01-02 15:04:05" // Format excludes offset/zone information
	dateNow := time.Now()               // This is local time with offset from UTC

	// Note: What if transitioned between standard and daylight savinges?
	offset := GetOffsetForLocal()
	dateAdjusted := dateNow.Add(-time.Second * time.Duration(offset))
	expectedDateString := dateAdjusted.Format(dateFormat)

	// Act
	str, cache := GetDateSubstitute(nil, "getdate", "utc", dateFormat)

	// Assert
	if date, ok := cache.(time.Time); ok {
		if !IsDateRangeWithinTollerance(date, dateNow, 1) {
			t.Errorf(
				"Returned date(%s) is not within tollerance of Now(%s)",
				date.Format(OutputDateFormat),
				dateNow.Format(OutputDateFormat),
			)
		}

		if str != expectedDateString {
			t.Errorf("Expected date formatted as %s but got %s", expectedDateString, str)
		}

		/// Need to consider zone and offset scenarios closer
		/// It appears all dates in this test are local (EST/EDT)

		// expectedZone, _ := dateAdjusted.Zone()
		// actualZone, _ := date.Zone()
		// if actualZone != "UTC" {
		// 	t.Errorf("Expected date to belong to UTC zone but got: %s", actualZone)
		// }
		// if expectedZone == actualZone {
		// 	t.Errorf("Expected zone unexpectantly matched actual zone: %s", actualZone)
		// }

	} else {
		t.Errorf("Expected time.Time value returned but got: %v", reflect.TypeOf(cache))
	}
}

func TestSetDateCreatesDate(t *testing.T) {
	// Arrange
	inputDate := "2021-02-25T17:21:05"
	expectedDate := time.Date(2021, 02, 25, 17, 21, 05, 0, time.Local)

	// Act
	str, cache := SetDateSubstitute(nil, "setdate", "local", inputDate)

	// Assert
	if date, ok := cache.(time.Time); ok {
		if !IsDateRangeWithinTollerance(date, expectedDate, 0) {
			t.Errorf(
				"Expected date (%s) but received (%s)",
				expectedDate.Format(OutputDateFormat),
				date.Format(OutputDateFormat))
		}
	} else {
		t.Errorf("Expected time.Time value returned but got: %v", reflect.TypeOf(cache))
	}

	if str != "" {
		t.Errorf("Expected empty string but received: %s", str)
	}
}

func TestSetDateDoesNotModifyCachedDate(t *testing.T) {
	// Arrange
	inputDate := "2021-02-25T17:21:05"
	altDate := "2018-11-02T04:45:52"
	expectedDate := time.Date(2021, 02, 25, 17, 21, 05, 0, time.Local)

	// Act
	var str string
	var cache interface{}
	_, cache = SetDateSubstitute(nil, "setdate", "local", inputDate)
	str, cache = SetDateSubstitute(cache, "setdate", "local", altDate)

	// Assert
	if date, ok := cache.(time.Time); ok {
		if !IsDateRangeWithinTollerance(date, expectedDate, 0) {
			t.Errorf(
				"Expected date (%s) but received (%s)",
				expectedDate.Format(OutputDateFormat),
				date.Format(OutputDateFormat))
		}
	} else {
		t.Errorf("Expected time.Time value returned but got: %v", reflect.TypeOf(cache))
	}

	if str != "" {
		t.Errorf("Expected empty string but received: %s", str)
	}
}

func TestSetDateReturnsMinDateForInvalidFormattedDate(t *testing.T) {
	// Arrange
	inputDate := "04:23:48"
	expectedDateMin := time.Time{}

	// Act
	str, cache := SetDateSubstitute(nil, "setdate", "local", inputDate)

	// Assert
	if date, ok := cache.(time.Time); ok {
		if !IsDateRangeWithinTollerance(date, expectedDateMin, 0) {
			t.Errorf(
				"Expected date (%s) but received (%s)",
				expectedDateMin.Format(OutputDateFormat),
				date.Format(OutputDateFormat))
		}
	} else {
		t.Errorf("Expected time.Time value returned but got: %v", reflect.TypeOf(cache))
	}

	if str != "" {
		t.Errorf("Expected empty string but received: %s", str)
	}
}

func TestSetDateOffset30MinutesBeforeLocalNow(t *testing.T) {
	// Arrange
	dateNow := time.Now() // This is local time with offset from UTC
	expectedDate := dateNow.
		Add(time.Minute * -30).
		Add(time.Second * 4)

	// Act
	str, cache := SetDateOffsetSubstitute(nil, "setdateoffset", "local", "n=-30;s=4")

	// Assert
	if date, ok := cache.(time.Time); ok {
		if !IsDateRangeWithinTollerance(date, expectedDate, 1) {
			t.Errorf(
				"Returned date(%s) is not within tollerance of expected date(%s)",
				date.Format(OutputDateFormat),
				expectedDate.Format(OutputDateFormat))
		}

		if str != "" {
			t.Errorf("Expected empty date string but got %s", str)
		}
	} else {
		t.Errorf("Expected time.Time value returned but got: %v", reflect.TypeOf(cache))
	}
}

/// IsDateRangeWithinTollerance - validate date2 is equal to date1 or no later than tollerance in seconds
func IsDateRangeWithinTollerance(date1, date2 time.Time, tollerance int64) bool {
	difference := date2.Unix() - date1.Unix()
	if debugEnabled {
		fmt.Printf("DEBUG: %d-%d=%d\n", date2.Unix(), date1.Unix(), difference)
	}
	return difference >= 0 && difference <= tollerance
}

func GetOffsetForLocal() int {
	dateNow := time.Now()
	_, offset := dateNow.Zone()
	return offset
}
