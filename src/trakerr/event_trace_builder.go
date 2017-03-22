package trakerr

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

//EventTraceBuilder is a static-eqsue struct which has methods assosiated with it to parse stacktraces.
type EventTraceBuilder struct {
}

//GetEventTraces takes an error, along with a depth and a skip and returns an array of parsed inner stacktraces.
//The depth is how many lines the stacktrace is long, to parse; while the skip bumps up the stacktrace.
//Without a skip, since the error is being passed to EventTraceBuilder the stacktrace of the current
//thread will start here, and go through main to where the error occured. In this case; skip removes
//the trakerr API internals from the trace, since they're not relevent, and depth is the final top of the trace.
func (tb *EventTraceBuilder) GetEventTraces(err interface{}, depth int, skip int) []InnerStackTrace {
	if err == nil {
		return nil
	}

	var traces = []InnerStackTrace{}

	return tb.AddStackTrace(traces, err, depth, skip+1)
}

//AddStackTrace adds a filled inner stacktrace to a list and returns it.
func (tb *EventTraceBuilder) AddStackTrace(traces []InnerStackTrace, err interface{}, depth int, skip int) []InnerStackTrace {
	var innerTrace = InnerStackTrace{}

	innerTrace.TraceLines = tb.GetTraceLines(err, depth, skip+1)
	innerTrace.Message = fmt.Sprint(err)
	innerTrace.Type_ = fmt.Sprintf("%T", err)

	traces = append(traces, innerTrace)
	return traces
}

//GetTraceLines parses each line of the stacktrace and returns an array of lines to populate InnerStackTrace
func (tb *EventTraceBuilder) GetTraceLines(err interface{}, depth int, skip int) []StackTraceLine {
	var traceLines = []StackTraceLine{}
	var goPath = tb.FileErrorHandler(filepath.Abs(os.Getenv("GOPATH")))
	var goRuntime = tb.FileErrorHandler(filepath.Abs(runtime.GOROOT()))
	for i := 0; i < depth; i++ {
		pc, file, line, ok := runtime.Caller(skip + 1 + i)
		if !ok {
			break
		}
		var localFilePath = tb.FileErrorHandler(filepath.Abs(file))

		var function = runtime.FuncForPC(pc)
		stLine := StackTraceLine{}

		var finalstring string
		if strings.Contains(strings.ToLower(localFilePath), strings.ToLower(goPath)) { //If it's goPath stacktrace
			finalstring = localFilePath[len(goPath):]
		} else if strings.Contains(strings.ToLower(localFilePath), strings.ToLower(goRuntime)) { //Otherwise its called from the runtime.
			finalstring = localFilePath[len(goRuntime):]
		} else {
			finalstring = localFilePath
		}

		stLine.File = strings.TrimLeft(finalstring, "\\/ ")
		stLine.Line = int32(line)
		stLine.Function = function.Name()
		traceLines = append(traceLines, stLine)
	}

	return traceLines
}

//FileErrorHandler is a small error handler for calls to find the paths of the files for path output parsing.
func (tb *EventTraceBuilder) FileErrorHandler(str string, er error) string {
	if er != nil {
		panic(er)
	}

	return str
}
