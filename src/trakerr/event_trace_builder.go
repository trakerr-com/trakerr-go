package trakerr

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type EventTraceBuilder struct {
}

func (tb *EventTraceBuilder) GetEventTraces(err interface{}, depth int, skip int) []InnerStackTrace {
	if err == nil {
		return nil
	}

	var traces = []InnerStackTrace{}

	return tb.AddStackTrace(traces, err, depth, skip+1)
}

func (tb *EventTraceBuilder) AddStackTrace(traces []InnerStackTrace, err interface{}, depth int, skip int) []InnerStackTrace {
	var innerTrace = InnerStackTrace{}

	innerTrace.TraceLines = tb.GetTraceLines(err, depth, skip+1)
	innerTrace.Message = fmt.Sprint(err)
	innerTrace.Type_ = fmt.Sprintf("%T", err)

	traces = append(traces, innerTrace)
	return traces
}

func (tb *EventTraceBuilder) GetTraceLines(err interface{}, depth int, skip int) []StackTraceLine {
	var traceLines = []StackTraceLine{}
	var goPath = os.Getenv("GOPATH")
	var goFilePath = tb.FileErrorHandler(filepath.Abs(goPath))
	for i := 0; i < depth; i++ {
		pc, file, line, ok := runtime.Caller(skip + 1 + i)
		if !ok {
			break
		}
		var localFilePath = tb.FileErrorHandler(filepath.Abs(file))

		var function = runtime.FuncForPC(pc)
		stLine := StackTraceLine{}
		stLine.File = strings.TrimPrefix(strings.ToLower(localFilePath), strings.ToLower(goFilePath))
		stLine.Line = int32(line)
		stLine.Function = function.Name()
		traceLines = append(traceLines, stLine)
	}

	return traceLines
}

//FileErrorHandler ...
func (tb *EventTraceBuilder) FileErrorHandler(str string, er error) string {
	if er != nil {
		panic(er)
	}

	return str
}
