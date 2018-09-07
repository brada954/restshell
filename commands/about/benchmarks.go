package about

import (
	"fmt"
	"io"
)

type BenchmarkTopic struct {
	Key         string
	Title       string
	Description string
	About       string
}

var localBenchmarkTopic = &BenchmarkTopic{
	Key:         "BENCHMARK",
	Title:       "Benchmarks",
	Description: "General benchmarking and the REST API support",
	About: `Benchmarking support includes functions to collect iteration data
and display the data in a consistent manner. The BMRESV and BMGET commands
are good examples to build on.

The benchmarking command set supports configuring iterations, throttling,
concurrency (in some cases) as well as output formating.

Benchmarking Control
Benchmarks provide control over the number of iterations, throttling, and
REST client initialization.

The --reconnect option causes a connection to be re-established on each
iteration so the connection time is included in each iteration.

Concurrency
Concurency is limited in support to bmget (and non-standard ccresv).

Concurrency testing also has an impact from reconnection issues. Each worker
thread that is run concurrently will generally create new connection if
all existing ones are in use.

To reduce this issue, the --warming option will process a concurrent job
before starting the benchmark. The number of jobs processed will equal the
the number of workers or the concurrency value. This however has not
been a panacea to the problem as observations have shown less consistency.
A long connect time can also stall the first jobs of the benchmark
as there may only be a 350ms delay waiting for those warming
connections to complete before continuing. (More work is necessary)

There are also settings for max idle connections and max idle connections
per host that affect consistency; MaxIdleConnsPerHost is boosted to 1000.

Output Formating

There are two CSV format options. The --csv option provides raw data 
in CSV format and --csv-fmt option may format things like times to
be easily read by humans.

REST API controls
Benchmarking also depends on the capabilities of the Rest client options.
Key options include controls on reconnect or validation of SSL certificates.

Consult the commands help for more information.

It is recommended to review the verbose output to determine if the tests
are performing as expected and the options enabled are working (until fully
implemented)
`,
}

func NewBenchmarkTopic() *BenchmarkTopic {
	return localBenchmarkTopic
}

func (a *BenchmarkTopic) GetKey() string {
	return a.Key
}

func (a *BenchmarkTopic) GetTitle() string {
	return a.Title
}

func (a *BenchmarkTopic) GetDescription() string {
	return a.Description
}

func (a *BenchmarkTopic) WriteAbout(o io.Writer) error {
	fmt.Fprintf(o, a.About)
	return nil
}
