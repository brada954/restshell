package shell

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type SiegeBucket struct {
	Label         string
	Period        int
	End           time.Time
	StartedJobs   int
	CompletedJobs int
	Errors        int
	mux           sync.Mutex // Protects the counters in slice
}

type Siegemark struct {
	StartTime           time.Time
	Duration            time.Duration
	BucketDuration      time.Duration
	LateStarts          int
	Note                string
	message             string
	summarized          bool
	successful          int
	started             int
	failed              int
	startedPerSecond    float64
	completedPerSecond  float64
	successfulPerSecond float64
	failedPerSecond     float64
	mux                 sync.Mutex // Protects counters like LateStarts
	Buckets             []SiegeBucket
	custom              interface{}
}

func NewSiegemark(duration time.Duration, buckets int) Siegemark {
	bd := duration / time.Duration(buckets)
	result := Siegemark{
		summarized:     false,
		StartTime:      mockableTimeNow(),
		BucketDuration: bd,
		Buckets:        make([]SiegeBucket, buckets),
	}

	endTime := result.StartTime
	for i := 0; i < len(result.Buckets); i++ {
		endTime = endTime.Add(bd)
		result.Buckets[i].Period = i
		result.Buckets[i].Label = strconv.Itoa(i)
		result.Buckets[i].End = endTime
	}
	return result
}

func (sm *Siegemark) Start() {
	sm.StartTime = mockableTimeNow()

	for i := 0; i < len(sm.Buckets); i++ {
		sm.Buckets[i].End = sm.StartTime.Add(sm.BucketDuration * time.Duration(i+1))
	}
}

func (sm *Siegemark) End() {
	sm.Duration = mockableTimeSince(sm.StartTime)
}

func (sm *Siegemark) StartIteration(i int) {
	var bucket *SiegeBucket
	now := mockableTimeNow()
	for i = 0; i < len(sm.Buckets); i++ {
		if sm.Buckets[i].End.Before(now) {
			continue
		}
		bucket = &sm.Buckets[i]
		break
	}

	if bucket == nil {
		// We do not count jobs starting late
		sm.mux.Lock()
		sm.LateStarts++
		sm.mux.Unlock()
		fmt.Println("Started Late: ", sm.StartTime, "-", now)
	} else {
		bucket.mux.Lock()
		bucket.StartedJobs++
		bucket.mux.Unlock()
		fmt.Println("Bucket: ", now, "-", bucket.End)
	}
}

func (sm *Siegemark) EndIteration(i int) {
	sm.EndIterationWithError(i, nil)
}

func (sm *Siegemark) EndIterationWithError(i int, err error) {
	now := mockableTimeNow()
	var bucket = &sm.Buckets[len(sm.Buckets)-1] // Default to last bucket
	for i = 0; i < len(sm.Buckets); i++ {
		if sm.Buckets[i].End.Before(now) {
			continue
		}
		bucket = &sm.Buckets[i]
		break
	}

	// Use last bucket even if job finished late as we are counting those.
	bucket.mux.Lock()
	bucket.CompletedJobs++
	if err != nil {
		bucket.Errors++
	}
	bucket.mux.Unlock()
}

func (sm *Siegemark) SetIterationStatus(i int, err error) {
	sm.UpdateIterationError(i, err)
}

func (sm *Siegemark) UpdateIterationError(i int, err error) {
	if err == nil {
		return
	}

	now := mockableTimeNow()
	var bucket = &sm.Buckets[len(sm.Buckets)-1] // Default to last bucket
	for i = 0; i < len(sm.Buckets); i++ {
		if sm.Buckets[i].End.Before(now) {
			continue
		}
		bucket = &sm.Buckets[i]
		break
	}

	// Use last bucket even if job finished late as we are counting those.
	bucket.mux.Lock()
	bucket.Errors++
	bucket.mux.Unlock()
}

func (sm *Siegemark) summarize() {

	// Totals calculations
	for _, v := range sm.Buckets {
		sm.started = sm.started + v.StartedJobs
		sm.successful = sm.successful + v.CompletedJobs
		sm.failed = sm.failed + v.Errors
	}

	// Avoid potential for divide by zero
	if sm.Duration == 0 {
		sm.Duration = mockableTimeSince(sm.StartTime)
	}
	if sm.Duration < time.Nanosecond {
		sm.Duration = time.Nanosecond
	}

	sm.startedPerSecond = float64(sm.started) / sm.Duration.Seconds()
	sm.completedPerSecond = float64(sm.successful+sm.failed) / sm.Duration.Seconds()
	sm.successfulPerSecond = float64(sm.successful) / sm.Duration.Seconds()
	sm.failedPerSecond = float64(sm.failed) / sm.Duration.Seconds()
	sm.summarized = true
}

func (sm *Siegemark) AddIterationMessage(i int, msg string) {
	//sm.Bucket[i].Messages = append(bm.Iterations[i].Messages, msg)
}

func (sm *Siegemark) Dump(label string, opts StandardOptions, showIterations bool) {
	if !sm.summarized {
		sm.summarize()
	}

	if !opts.IsHeaderDisabled() {
		var note string
		if len(sm.Note) > 0 {
			note = fmt.Sprintf(" -- %s", sm.Note)
		}

		if opts.IsCsvOutputEnabled() {
			var headingFmt = "%[1]s,%[2]s,%[3]s,%[4]s,%[5]s,%[6]s,%[7]s,%[8]s,%[9]s,%[10]s\n"
			fmt.Fprintf(OutputWriter(),
				headingFmt,
				"Label", "Total", "Success", "Error", "Avg Total", "Avg Success", "Avg Error", "Late", "Message", note)
		} else {
			var headingFmt = "%-14[1]s  %6[2]s  %8[3]s  %12[4]s  %12[5]s %12[6]s  %8[7]s  %8[8]s %[9]s %[10]s\n"
			fmt.Fprintf(OutputWriter(),
				headingFmt,
				"Label", "Count", "Success", "Error", "Avg Total", "Avg Success", "Avg Error", "Late", "Message", note)
		}
	}

	if opts.IsFormattedCsvEnabled() {
		var displayFmt = "%s,%d,%d,%d,%f,%f,%f,%d,%s\n"
		fmt.Fprintf(OutputWriter(),
			displayFmt,
			label,
			sm.started,
			sm.successful,
			sm.failed,
			sm.completedPerSecond,
			sm.successfulPerSecond,
			sm.failedPerSecond,
			sm.LateStarts,
			sm.message,
		)
	} else if opts.IsCsvOutputEnabled() {
		var displayFmt = "%s,%d,%d,%d,%f,%f,%f,%d,%s\n"
		fmt.Fprintf(OutputWriter(),
			displayFmt,
			label,
			sm.started,
			sm.successful,
			sm.failed,
			sm.completedPerSecond,
			sm.successfulPerSecond,
			sm.failedPerSecond,
			sm.LateStarts,
			sm.message,
		)
	} else {
		var displayFmt = "%-14s  %6d  %8d  %12d  %12f %12f  %8f  %8d %s\n"
		fmt.Fprintf(OutputWriter(),
			displayFmt,
			label,
			sm.started,
			sm.successful,
			sm.failed,
			sm.completedPerSecond,
			sm.successfulPerSecond,
			sm.failedPerSecond,
			sm.LateStarts,
			sm.message,
		)
	}

	if showIterations {
		sm.DumpIterations(opts)
	}
}

func (sm *Siegemark) DumpIterations(opts StandardOptions) {
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

	for _, v := range sm.Buckets {
		v.DumpLine(opts, displayFmt /*, bm.StartTime*/)
	}
}

func (sb *SiegeBucket) DumpLine(opts StandardOptions, format string) {
	fmt.Println("Bucket", sb.Label, sb.StartedJobs, sb.CompletedJobs, sb.Errors)
}

// func (sm *Siegemark) WallAverageFmt() string {
// 	return FormatMsTime(bm.WallAverageInMs())
// }

// func (sm *Siegemark) HighTimeFmt() string {
// 	return FormatMsTime(bm.HighTimeInMs())
// }

// func (sm *Siegemark) LowTimeFmt() string {
// 	return FormatMsTime(bm.LowTimeInMs())
// }

// func (sm *Siegemark) WallTimeFmt() string {
// 	return FormatMsTime(bm.WallTimeInMs())
// }

// func (sm *Siegemark) HlAverageFmt() string {
// 	return FormatMsTime(bm.HlAverageInMs())
// }

// func (sm *Siegemark) WallAverageInMs() float64 {
// 	if !bm.summarized {
// 		bm.summarize()
// 	}
// 	return bm.avgWallTimeMs
// }

// func (sm *Siegemark) HighTimeInMs() float64 {
// 	if !bm.summarized {
// 		bm.summarize()
// 	}
// 	return bm.highTimeMs
// }

// func (sm *Siegemark) LowTimeInMs() float64 {
// 	if !bm.summarized {
// 		bm.summarize()
// 	}
// 	return bm.lowTimeMs
// }

// func (sm *Siegemark) HlAverageInMs() float64 {
// 	if !bm.summarized {
// 		bm.summarize()
// 	}
// 	return bm.hlAvgWallTimeMs
// }
