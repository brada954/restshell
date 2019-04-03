package shell

import (
	"errors"
	"testing"
	"time"
)

func TestSiegeMarkIterations3(t *testing.T) {
	tm := time.Now()
	defer mockNowCleanup(mockNow([]time.Time{
		tm.Add(time.Millisecond * 10),
		tm.Add(time.Millisecond * 16), //I1 - B0
		tm.Add(time.Millisecond * 500),
		tm.Add(time.Millisecond * 650), //I2 - B0
		tm.Add(time.Millisecond * 1000),
		tm.Add(time.Millisecond * 1114), // I3 - B1
		tm.Add(time.Millisecond * 1250), // I3 - B1
		tm.Add(time.Millisecond * 1300), // I4 - B1
		tm.Add(time.Millisecond * 1400), // I4 - B1
		tm.Add(time.Millisecond * 4000), // I5 - B3
		tm.Add(time.Millisecond * 4300), // I5 - B3
		tm.Add(time.Millisecond * 4100), // I6 - B3
		tm.Add(time.Millisecond * 4500), // I6 - B4
		tm.Add(time.Millisecond * 5510),
		tm.Add(time.Millisecond * 11200),
		tm.Add(time.Millisecond * 11250),
	}))

	// var expectedFmt string
	// var expectedVal float64

	sm := NewSiegemark(time.Second*time.Duration(10), 10) // Consumes one time now value
	for i := 0; i < 6; i++ {
		sm.StartIteration(i)
		sm.EndIteration(i)
	}

	opts := GetStdOptions()
	disabled := false
	enabled := true
	opts.csvOutputOption = &disabled
	opts.fullOutputOption = &enabled
	sm.Dump("Test", opts, true)
}

func TestSiegemarkDump(t *testing.T) {
	tm := time.Now()
	defer mockNowCleanup(mockNow([]time.Time{
		tm.Add(time.Second * 2),
		tm.Add(time.Second * 4),
		tm.Add(time.Second * 6),
		tm.Add(time.Second * 8),
		tm.Add(time.Second * 9),
		tm.Add(time.Second * 10),
		tm.Add(time.Second * 11),
		tm.Add(time.Second * 12),
		tm.Add(time.Second * 13),
	}))

	// var expectedFmt string
	// var expectedVal float64

	sm := NewSiegemark(time.Second*time.Duration(12), 10)
	for i := 0; i < 6; i++ {
		sm.StartIteration(i)
		sm.EndIteration(i)
		if i == 2 {
			sm.SetIterationStatus(i, errors.New("Fake failure!"))
		} else {
			sm.SetIterationStatus(i, nil)
		}
	}

	// This test just dumps output to get visual representation
	opts := GetStdOptions()
	csv := false
	opts.csvOutputOption = &csv
	sm.Dump("Test", opts, true)

	csv = true
	sm.Dump("CSV Test", opts, true)
}
