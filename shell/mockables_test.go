package shell

import (
	"time"
)

// External references that can be mocked:
//  var MockableTimeSince = time.Since
//  var MockableTimeNow = time.Now

///////////////////////////////////////////////
// Mock time.Since for various iterations
var useSinceData []time.Duration
var currDataIdx = 0

func mockSince(data []time.Duration) func(time.Time) time.Duration {
	mockableTimeSince = mockSinceImpl
	useSinceData = data
	currDataIdx = 0
	return time.Since
}

func mockSinceImpl(time.Time) time.Duration {
	result := useSinceData[currDataIdx]
	if currDataIdx < (len(useSinceData) - 1) {
		currDataIdx = currDataIdx + 1
	}
	return result
}

func mockSinceCleanup(fn func(time.Time) time.Duration) {
	mockableTimeSince = fn
}

var mockNowIndex int
var mockNowTimes []time.Time

func mockNow(times []time.Time) func() time.Time {
	mockNowIndex = 0
	mockNowTimes = times
	mockableTimeNow = mockNowImpl
	return time.Now
}

func mockNowImpl() time.Time {
	result := mockNowTimes[mockNowIndex]
	mockNowIndex++
	if mockNowIndex == len(mockNowTimes) {
		mockNowIndex = 0
	}
	return result
}

func mockNowCleanup(fn func() time.Time) {
	mockableTimeNow = fn
}
