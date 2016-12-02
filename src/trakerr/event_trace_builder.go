package trakerr

import (
	"trakerr_client"
	"runtime"
	"fmt"
)

type EventTraceBuilder struct {

}


func (tb *EventTraceBuilder) GetEventTraces(err interface{}, depth int) []trakerr_client.InnerStackTrace {
	if(err == nil) { return nil }

	var traces = []trakerr_client.InnerStackTrace{}

	return tb.AddStackTrace(traces, err, depth)
}

func (tb *EventTraceBuilder) AddStackTrace(traces []trakerr_client.InnerStackTrace, err interface{}, depth int) []trakerr_client.InnerStackTrace {
	var innerTrace = trakerr_client.InnerStackTrace{}

	innerTrace.TraceLines = tb.GetTraceLines(err, depth);
	innerTrace.Message = fmt.Sprint(err)
	innerTrace.Type_ = fmt.Sprintf("%T", err)

	traces = append(traces, innerTrace)
	return traces
}

func (tb *EventTraceBuilder) GetTraceLines(err interface{}, depth int) []trakerr_client.StackTraceLine {
	var traceLines = []trakerr_client.StackTraceLine{};

	for i:= 0;i< depth;i++ {
		pc, file, line, ok := runtime.Caller(i)
		if(!ok) { break; }

		var function = runtime.FuncForPC(pc)
		stLine := trakerr_client.StackTraceLine{}
		stLine.File  = file
		stLine.Line = int32(line)
		stLine.Function = function.Name()
		traceLines = append(traceLines, stLine)
	}

	return traceLines
}
