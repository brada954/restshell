package shell

import (
	"fmt"
	"math"
	"time"
)

type Benchmark struct {
	Iterations        []BenchmarkIteration
	StartTime         time.Time
	Duration          time.Duration
	Note              string
	message           string
	summarized        bool
	errors            int
	totalTimeMs       float64
	avgWallTimeMs     float64
	hlAvgWallTimeMs   float64
	highTimeMs        float64
	lowTimeMs         float64
	highTimeIndex     int
	lowTimeIndex      int
	standardDeviation float64
	custom            interface{}
}

type BenchmarkIteration struct {
	Iteration int
	WallTime  int64 // Nano-seconds
	Err       error
	Messages  []string
	Custom    interface{}
	start     time.Time
	end       time.Time
}

func NewBenchmark(iterations int) Benchmark {
	result := Benchmark{
		summarized: false,
		Iterations: make([]BenchmarkIteration, iterations, iterations),
		StartTime:  time.Now(),
	}
	for i := range result.Iterations {
		result.Iterations[i].Iteration = i
		result.Iterations[i].Messages = make([]string, 0)
	}
	return result
}

func (bm *Benchmark) summarize() {
	var successful int = 0
	var totalWall int64 = 0
	var highTime int64 = 0
	var lowTime int64 = math.MaxInt64
	var highTimeIndex int = 0
	var lowTimeIndex int = 0

	// Averge calculations
	for k, v := range bm.Iterations {
		if v.Err == nil {
			successful = successful + 1
			totalWall = totalWall + int64(v.WallTime)
			if v.WallTime > highTime {
				highTime = int64(v.WallTime)
				highTimeIndex = k
			}
			if v.WallTime < lowTime {
				lowTime = v.WallTime
				lowTimeIndex = k
			}
		}
	}

	bm.errors = len(bm.Iterations) - successful
	bm.totalTimeMs = float64(totalWall) / 1000000.0
	bm.avgWallTimeMs = (float64(totalWall) / float64(successful)) / 1000000.0
	bm.hlAvgWallTimeMs = bm.avgWallTimeMs
	if successful > 2 {
		bm.hlAvgWallTimeMs = (float64(totalWall-highTime-lowTime) / float64(successful-2)) / 1000000.0
	}
	bm.highTimeMs = float64(highTime) / 1000000.0
	bm.highTimeIndex = highTimeIndex
	bm.lowTimeMs = float64(lowTime) / 1000000.0
	bm.lowTimeIndex = lowTimeIndex
	bm.summarized = true

	// Standard deviation calculation
	var n float64 = 0.0
	var sumOfSq float64 = 0.0
	for _, v := range bm.Iterations {
		if v.Err == nil {
			n = n + 1
			sumOfSq = sumOfSq + math.Pow((v.WallTimeInMs()-bm.avgWallTimeMs), 2.0)
		}
	}
	if n > 0 {
		bm.standardDeviation = math.Sqrt(sumOfSq / n)
		if bm.hlAvgWallTimeMs < (bm.standardDeviation * 10) {
			bm.message = fmt.Sprintf("StdDev (%.3f) exceeds 10%%", bm.standardDeviation)
		}
	} else {
		bm.standardDeviation = 0
	}

	if bm.Duration == 0 {
		bm.Duration = time.Since(bm.StartTime)
	}
}

func (bm *Benchmark) Start() {
	bm.StartTime = time.Now()
}

func (bm *Benchmark) End() {
	bm.Duration = time.Since(bm.StartTime)
}

func (bm *Benchmark) WallTimeInMs() float64 {
	return float64(bm.Duration) / float64(time.Millisecond)
}

func (bm *Benchmark) StartIteration(i int) {
	bm.Iterations[i].start = time.Now()
}

func (bm *Benchmark) EndIteration(i int) {
	bm.EndIterationWithError(i, nil)
}

func (bm *Benchmark) EndIterationWithError(i int, err error) {
	bm.Iterations[i].WallTime = int64(mockableTimeSince(bm.Iterations[i].start))
	bm.Iterations[i].end = bm.Iterations[i].start.Add(time.Duration(bm.Iterations[i].WallTime))
	bm.Iterations[i].Err = err
}

func (bm *Benchmark) SetIterationStatus(i int, err error) {
	bm.UpdateIterationError(i, err)
}

func (bm *Benchmark) UpdateIterationError(i int, err error) {
	bm.Iterations[i].Err = err
}

func (bm *Benchmark) AddIterationMessage(i int, msg string) {
	bm.Iterations[i].Messages = append(bm.Iterations[i].Messages, msg)
}

func (bm *Benchmark) Dump(label string, opts StandardOptions, showIterations bool) {
	if !bm.summarized {
		bm.summarize()
	}

	var highIndexLen = numOfDigits(bm.highTimeIndex)
	var lowIndexLen = numOfDigits(bm.lowTimeIndex)

	if !opts.IsHeaderDisabled() {
		var note string
		if len(bm.Note) > 0 {
			note = fmt.Sprintf(" -- %s", bm.Note)
		}

		if opts.IsCsvOutputEnabled() {
			var headingFmt = "%[1]s,%[2]s,%[3]s,%[4]s,%[5]s,%[6]s,%[7]s,%[8]s,%[9]s,%[10]s,%[11]s,%[12]s\n"
			fmt.Fprintf(OutputWriter(),
				headingFmt,
				"Label", "Count", "Err", "Avg", "High", "HI", "Low", "LI", "Avg-(HL)", "Tot", "Message", note)
		} else {
			var headingFmt = "%-14[1]s  %5[2]s  %5[3]s  %8[4]s  %8[5]s%[13]*[6]s  %8[7]s%[14]*[8]s  %8[9]s  %8[10]s %8[11]s %[12]s\n"
			fmt.Fprintf(OutputWriter(),
				headingFmt,
				"Label", "Count", "Err", "Avg", "High", "", "Low", "", "Avg-(HL)", "Tot", "Message",
				note,
				highIndexLen+2,
				lowIndexLen+2)
		}
	}

	if opts.IsFormattedCsvEnabled() {
		var displayFmt = "%s,%d,%d,%s,%s,%d,%s,%d,%s,%s,%s\n"
		fmt.Fprintf(OutputWriter(),
			displayFmt,
			label,
			len(bm.Iterations),
			bm.errors,
			bm.WallAverageFmt(),
			bm.HighTimeFmt(),
			bm.highTimeIndex,
			bm.LowTimeFmt(),
			bm.lowTimeIndex,
			bm.HlAverageFmt(),
			bm.WallTimeFmt(),
			bm.message,
		)
	} else if opts.IsCsvOutputEnabled() {
		var displayFmt = "%s,%d,%d,%f,%f,%d,%f,%d,%f,%f,%s\n"
		fmt.Fprintf(OutputWriter(),
			displayFmt,
			label,
			len(bm.Iterations),
			bm.errors,
			bm.WallAverageInMs(),
			bm.HighTimeInMs(),
			bm.highTimeIndex,
			bm.LowTimeInMs(),
			bm.lowTimeIndex,
			bm.HlAverageInMs(),
			bm.WallTimeInMs(),
			bm.message,
		)
	} else {
		var displayFmt = "%-14s  %5d  %5d %9s %9s(%*d) %9s(%*d) %9s %9s  %s\n"
		fmt.Fprintf(OutputWriter(),
			displayFmt,
			label,
			len(bm.Iterations),
			bm.errors,
			bm.WallAverageFmt(),
			bm.HighTimeFmt(),
			highIndexLen,
			bm.highTimeIndex,
			bm.LowTimeFmt(),
			lowIndexLen,
			bm.lowTimeIndex,
			bm.HlAverageFmt(),
			bm.WallTimeFmt(),
			bm.message,
		)
	}

	if showIterations {
		bm.DumpIterations(opts)
	}
}

func (bm *Benchmark) DumpIterations(opts StandardOptions) {
	var headingFmt = "-,%s,%s,%s,%s,%s,%s\n"
	var displayFmt = "-,%d,%f,%d,%d,%s,%s\n"
	if opts.IsFormattedCsvEnabled() {
		displayFmt = "-,%d,%s,%d,%d,%s,%s\n"
	} else if opts.IsCsvOutputEnabled() {
		// CSV which is default
	} else {
		headingFmt = "  %9s  %8s  %10s  %10s  %3s  %s\n"
		displayFmt = "  %9d %9s  %10d  %10d  %3s  %s\n"
	}
	if len(headingFmt) > 0 && !opts.IsHeaderDisabled() {
		fmt.Fprintf(OutputWriter(), headingFmt, "Iteration", "WallTime", "Start", "End", "Err", "Message")
	}

	for _, v := range bm.Iterations {
		v.DumpLine(opts, displayFmt, bm.StartTime)
	}
}

func (bi *BenchmarkIteration) DumpLine(opts StandardOptions, displayFmt string, start time.Time) {
	errMark := ""
	errMsg := ""
	if bi.Err != nil {
		errMark = " X "
		errMsg = bi.Err.Error()
	}

	if opts.IsCsvOutputEnabled() && !opts.IsFormattedCsvEnabled() {
		fmt.Fprintf(OutputWriter(),
			displayFmt,
			bi.Iteration,
			bi.WallTimeInMs(),
			bi.start.Sub(start)/time.Microsecond,
			bi.end.Sub(start)/time.Microsecond,
			errMark,
			errMsg,
		)
	} else {
		fmt.Fprintf(OutputWriter(),
			displayFmt,
			bi.Iteration,
			bi.WallTimeFmt(),
			bi.start.Sub(start)/time.Microsecond,
			bi.end.Sub(start)/time.Microsecond,
			errMark,
			errMsg,
		)
	}
}

func (bm *Benchmark) WallAverageFmt() string {
	return FormatMsTime(bm.WallAverageInMs())
}

func (bm *Benchmark) HighTimeFmt() string {
	return FormatMsTime(bm.HighTimeInMs())
}

func (bm *Benchmark) LowTimeFmt() string {
	return FormatMsTime(bm.LowTimeInMs())
}

func (bm *Benchmark) WallTimeFmt() string {
	return FormatMsTime(bm.WallTimeInMs())
}

func (bm *Benchmark) HlAverageFmt() string {
	return FormatMsTime(bm.HlAverageInMs())
}

func (bm *Benchmark) WallAverageInMs() float64 {
	if !bm.summarized {
		bm.summarize()
	}
	return bm.avgWallTimeMs
}

func (bm *Benchmark) HighTimeInMs() float64 {
	if !bm.summarized {
		bm.summarize()
	}
	return bm.highTimeMs
}

func (bm *Benchmark) LowTimeInMs() float64 {
	if !bm.summarized {
		bm.summarize()
	}
	return bm.lowTimeMs
}

func (bm *Benchmark) HlAverageInMs() float64 {
	if !bm.summarized {
		bm.summarize()
	}
	return bm.hlAvgWallTimeMs
}

func (bi *BenchmarkIteration) WallTimeFmt() string {
	return FormatMsTime(bi.WallTimeInMs())
}

func (bi *BenchmarkIteration) WallTimeInMs() float64 {
	return float64(bi.WallTime) / 1000000.0
}

// Requires positive x
func numOfDigits(x int) int {
	if x < 0 {
		panic("Invalid number to get digits")
	}

	var start int64 = 9
	var i int = 0
	for ; start < (math.MaxInt32/10)-9; i++ {
		if int64(x) <= start {
			return i + 1
		}
		start = (start * 10) + 9
	}
	return i
}
