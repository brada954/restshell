package shell

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

// SiegeBucket -- structure for histogram of siege benchmark
type SiegeBucket struct {
	// Initialized data
	Label  string
	Period int
	End    time.Time

	// Updated data protected by mutex
	StartedJobs      int
	SuccessfulJobs   int
	SuccessDurations time.Duration
	FailedJobs       int
	FailedDurations  time.Duration
	TotalDuration    time.Duration

	mux sync.Mutex // Protects the updated data in bucket
}

// Siegemark -- JobMonitor for siege benchmarking
type Siegemark struct {
	StartTime         time.Time
	Duration          time.Duration
	BucketDuration    time.Duration
	LateStarts        int
	Note              string
	message           string
	summarized        bool
	totalJobs         int // Total jobs initiated
	totalDuration     time.Duration
	successfulJobs    int // Total jobs completed without API call
	successDuration   time.Duration
	failedJobs        int // Total jobs failed the API call
	failedDuration    time.Duration
	requestsPerSecond float64
	avgRequest        float64
	failuresPerSecond float64
	mux               sync.Mutex // Protects counters like LateStarts
	Buckets           []SiegeBucket
	custom            interface{}
}

// SiegemarkIteration -- track a single rest call
type SiegemarkIteration struct {
	Iteration int
	StartTime time.Time
	EndTime   time.Time
	Error     error
}

// NewSiegemark -- Create siege benchmarking job monitor
func NewSiegemark(duration time.Duration, buckets int) *Siegemark {

	// Calculate bucket durations and initialize benchmark parameters
	bd := duration / time.Duration(buckets)
	result := Siegemark{
		summarized:     false,
		BucketDuration: bd,
		Buckets:        make([]SiegeBucket, buckets),
	}

	// Initialize period data in buckets
	for i := 0; i < len(result.Buckets); i++ {
		result.Buckets[i].Period = i
		result.Buckets[i].Label = strconv.Itoa(i)
	}
	return &result
}

// Start -- Start a benchmark
func (sm *Siegemark) Start() {
	sm.StartTime = mockableTimeNow()
	for i := 0; i < len(sm.Buckets); i++ {
		// Because of calculation below, use time in next bucket as start >= End
		sm.Buckets[i].End = sm.StartTime.Add(sm.BucketDuration * time.Duration(i+1))
	}
}

// StartIteration -- Start an iteration
func (sm *Siegemark) StartIteration(i int) JobContext {
	now := mockableTimeNow()

	return &SiegemarkIteration{
		Iteration: i,
		StartTime: now,
	}
}

// End -- Record the End a benchmark
func (sm *Siegemark) End() {
	sm.Duration = mockableTimeSince(sm.StartTime)
}

// FinalizeIteration -- Fold iteration data into aggregated bucket data
func (sm *Siegemark) FinalizeIteration(jc JobContext) {
	if si, ok := jc.(*SiegemarkIteration); ok {
		bucket := sm.getBucket(si.EndTime)
		if bucket == nil {
			// We do not count jobs starting late
			sm.mux.Lock()
			sm.LateStarts++
			sm.mux.Unlock()
		} else {
			dur := si.EndTime.Sub(si.StartTime)
			bucket.mux.Lock()
			bucket.StartedJobs++
			bucket.TotalDuration += dur
			if si.Error != nil {
				bucket.FailedJobs++
				bucket.FailedDurations += dur
			} else {
				bucket.SuccessfulJobs++
				bucket.SuccessDurations += dur
			}
			bucket.mux.Unlock()
		}
	}
}

// EndIteration -- Collect completion data on iteration
func (si *SiegemarkIteration) EndIteration(err error) {
	now := mockableTimeNow()
	si.EndTime = now
	si.Error = err
}

// UpdateError -- Update the error on an iteration
func (si *SiegemarkIteration) UpdateError(err error) {
	si.Error = err
}

func (sm *Siegemark) getBucket(now time.Time) *SiegeBucket {
	for i := 0; i < len(sm.Buckets); i++ {
		var bucket = &sm.Buckets[i]
		if now.Before(bucket.End) {
			return bucket
		}
	}
	return &sm.Buckets[len(sm.Buckets)-1]
}

func (sm *Siegemark) summarize() {

	if sm.summarized {
		return
	}

	// Totals calculations
	for idx := range sm.Buckets {
		v := &sm.Buckets[idx]
		sm.totalJobs = sm.totalJobs + v.StartedJobs
		sm.totalDuration = sm.totalDuration + v.TotalDuration

		sm.successfulJobs = sm.successfulJobs + v.SuccessfulJobs
		sm.successDuration = sm.successDuration + v.SuccessDurations

		sm.failedJobs = sm.failedJobs + v.FailedJobs
		sm.failedDuration = sm.failedDuration + v.FailedDurations
	}

	// Avoid potential for divide by zero
	if sm.Duration < time.Nanosecond {
		sm.Duration = time.Nanosecond
	}

	// requestsPerSecond are total requests within the wall time of benchmark
	sm.requestsPerSecond = float64(sm.totalJobs) / sm.Duration.Seconds()

	// Note: avgRequest uses the duration total of each request
	sm.avgRequest = sm.totalDuration.Seconds() / float64(sm.totalJobs)

	// Note failures per second is based on wall time of benchmark
	sm.failuresPerSecond = float64(sm.failedJobs) / sm.Duration.Seconds()
	sm.summarized = true
}

func (sm *Siegemark) AddIterationMessage(i int, msg string) {
	//sm.Bucket[i].Messages = append(bm.Iterations[i].Messages, msg)
}

// Dump -- Dump to output stream the information formated as options requested
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
			var headingFmt = "%-16[1]s  %8[2]s  %8[3]s  %8[4]s  %12[5]s %12[6]s  %8[7]s  %12[8]s %[9]s %[10]s\n"
			fmt.Fprintf(OutputWriter(),
				headingFmt,
				"Label", "Count", "Success", "Error", "Avg Req", "Req/Sec", "Err/Sec", "Duration(S)", "Message", note)
		}
	}

	if opts.IsFormattedCsvEnabled() {
		var displayFmt = "%s,%d,%d,%d,%f,%f,%f,%d,%s\n"
		fmt.Fprintf(OutputWriter(),
			displayFmt,
			label,
			sm.totalJobs,
			sm.successfulJobs,
			sm.failedJobs,
			sm.avgRequest,
			sm.requestsPerSecond,
			sm.failuresPerSecond,
			sm.LateStarts,
			sm.message)
	} else if opts.IsCsvOutputEnabled() {
		var displayFmt = "%s,%d,%d,%d,%f,%f,%f,%d,%s\n"
		fmt.Fprintf(OutputWriter(),
			displayFmt,
			label,
			sm.totalJobs,
			sm.successfulJobs,
			sm.failedJobs,
			sm.avgRequest,
			sm.requestsPerSecond,
			sm.failuresPerSecond,
			sm.LateStarts,
			sm.message,
		)
	} else {
		var displayFmt = "%-16.16s  %8d  %8d  %8d  %12f %12f  %8f  %12f %s\n"
		fmt.Fprintf(OutputWriter(),
			displayFmt,
			label,
			sm.totalJobs,
			sm.successfulJobs,
			sm.failedJobs,
			sm.avgRequest,
			sm.requestsPerSecond,
			sm.failuresPerSecond,
			sm.Duration.Seconds(),
			sm.message,
		)
	}

	if showIterations {
		sm.DumpIterations(opts)
	}
}

// DumpIterations -- Dumps the iterations or with Siegemark the buckets
func (sm *Siegemark) DumpIterations(opts StandardOptions) {
	var headingFmt = "%s,%s,%s,%s,%s,%s\n"
	var displayFmt = "%s,%d,%d,%d,%f,%s\n"
	if opts.IsFormattedCsvEnabled() {
		displayFmt = "%s,%d,%d,%d,%f,%s\n"
	} else if opts.IsCsvOutputEnabled() {
		// CSV which is default
	} else {
		headingFmt = "  %9s  %10s  %10s  %10s  %10s %s\n"
		displayFmt = "  %9s  %10d  %10d  %10d  %10f %s\n"
	}
	if len(headingFmt) > 0 && !opts.IsHeaderDisabled() {
		fmt.Fprintf(OutputWriter(), headingFmt, "Bucket", "Start", "End", "Err", "AvgReq", "Message")
	}

	for idx := range sm.Buckets {
		v := &sm.Buckets[idx]
		v.dumpLine(opts, displayFmt /*, bm.StartTime*/)
	}
}

func (sb *SiegeBucket) dumpLine(opts StandardOptions, format string) {
	avgReq := 0.0
	if sb.StartedJobs > 0 {
		avgReq = sb.TotalDuration.Seconds() / float64(sb.StartedJobs)
	}
	fmt.Fprintf(OutputWriter(), format, sb.Label, sb.StartedJobs, sb.SuccessfulJobs, sb.FailedJobs, avgReq, "")
}
