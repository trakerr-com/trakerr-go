package trakerr

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

type EventTraceBuilder struct {
}

//GetEventTraces ...
func (tb *EventTraceBuilder) GetEventTraces(err interface{}, depth int) []InnerStackTrace {
	if err == nil {
		return nil
	}

	var traces = []InnerStackTrace{}

	return tb.AddStackTrace(traces, err, depth)
}

//AddStackTrace ...
func (tb *EventTraceBuilder) AddStackTrace(traces []InnerStackTrace, err interface{}, depth int) []InnerStackTrace {
	var innerTrace = InnerStackTrace{}

	innerTrace.TraceLines = tb.GetTraceLines(err, depth)
	innerTrace.Message = fmt.Sprint(err)
	innerTrace.Type_ = fmt.Sprintf("%T", err)

	traces = append(traces, innerTrace)
	return traces
}

//GetTraceLines ..
func (tb *EventTraceBuilder) GetTraceLines(err interface{}, depth int) []StackTraceLine {
	var traceLines = []StackTraceLine{}
	var goPath = os.Getenv("GOPATH")

	for i := 0; i < depth; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		var function = runtime.FuncForPC(pc)
		stLine := StackTraceLine{}
		stLine.File = strings.TrimPrefix(file, goPath)
		stLine.Line = int32(line)
		stLine.Function = function.Name()
		traceLines = append(traceLines, stLine)
	}

	return traceLines
}
