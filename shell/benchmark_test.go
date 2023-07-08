package shell

import (
	"errors"
	"testing"
	"time"
)

func TestBenchMarkIterations3(t *testing.T) {
	defer mockSinceCleanup(mockSince([]time.Duration{2000000, 3600000, 4000000}))

	var expectedFmt string
	var expectedVal float64

	bm := NewBenchmark(3)
	for k := range bm.Iterations {
		jc := bm.StartIteration(k)
		jc.EndIteration(nil)
		bm.FinalizeIteration(jc)
	}

	avg := bm.WallAverageInMs()
	expectedVal = float64(9600000/3) / 1000000.0
	if avg != expectedVal {
		t.Errorf("Error wall average; expected %f got %f", expectedVal, avg)
	}

	avgfmt := bm.WallAverageFmt()
	expectedFmt = "3.200ms"
	if avgfmt != expectedFmt {
		t.Errorf("Error wall average format; expected %s got %s", expectedFmt, avgfmt)
	}

	avg = bm.HlAverageInMs()
	expectedVal = float64(3600000) / 1000000.0
	if avg != expectedVal {
		t.Errorf("Error wall average; expected %f got %f", expectedVal, avg)
	}

	avgfmt = bm.HlAverageFmt()
	expectedFmt = "3.600ms"
	if avgfmt != expectedFmt {
		t.Errorf("Error wall average format; expected %s got %s", expectedFmt, avgfmt)
	}
}

func TestBenchMarkDump(t *testing.T) {
	testdata := []time.Duration{2000000, 3600000, 4000000, 3800000}
	defer mockSinceCleanup(mockSince(testdata))

	bm := NewBenchmark(len(testdata))
	for k := range bm.Iterations {
		jc := bm.StartIteration(k)
		jc.EndIteration(nil)
		if k == 2 {
			jc.UpdateError(errors.New("Fake failure!"))
		}
		bm.FinalizeIteration(jc)
	}

	// This test just dumps output to get visual representation
	opts := GetStdOptions()
	csv := false
	opts.csvOutputOption = &csv
	bm.Dump("Test", opts, true)

	csv = true
	bm.Dump("CSV Test", opts, true)
}

func TestBenchMarkWallAvgFormats(t *testing.T) {
	var val, expected string
	var bm = Benchmark{}

	bm.summarized = true
	bm.avgWallTimeMs = 3422.25
	expected = "3.422S"
	val = bm.WallAverageFmt()
	if val != expected {
		t.Errorf("Unexpected format for %s; got %s", expected, val)
	}

	bm.avgWallTimeMs = .0000232
	expected = "0.000023ms"
	val = bm.WallAverageFmt()
	if val != expected {
		t.Errorf("Unexpected format for %s; got %s", expected, val)
	}
}
