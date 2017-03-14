package main

import (
	"os"

	"github.com/trakerr-io/trakerr-go/src/trakerr"
)

func main() {
	var client *trakerr.TrakerrClient
	if len(os.Args) > 1 {
		client = trakerr.NewTrakerrClientWithDefaults(
			os.Args[1],
			"1.0",
			"development")
	} else {
		client = trakerr.NewTrakerrClientWithDefaults(
			"<Api Key here>",
			"1.0",
			"development")
	}

	//Option-1: Global error handling

	appEvent := client.NewEmptyEvent()
	// set any custom data on appEvent
	appEvent.CustomProperties.StringData.CustomData1 = "foo"
	appEvent.CustomProperties.StringData.CustomData2 = "bar"
	appEvent.EventUser = "john@trakerr.io"
	appEvent.EventSession = "12"

	ts := TestSession{client, appEvent}
	buf := []int{1, 2, 3}
	te := TestError{}
	te.BufferOverflowError(buf, 4, ts)

	// Option-2: send error
	/*err := errors.New("Something bad happened here")
	client.SendError("Error", err)

	// Option-3: send error with custom properties
	appEventWithErr := client.CreateAppEventFromError("Error", err)

	// set any custom data on appEvent
	appEventWithErr.CustomProperties.StringData.CustomData1 = "foo"
	appEventWithErr.CustomProperties.StringData.CustomData2 = "bar"

	client.SendEvent(appEventWithErr)

	// Option-4: send event manually
	appEventCustom := client.NewAppEvent("Info", "SomeType", "SomeMessage")

	// set any custom data on appEvent
	appEventCustom.CustomProperties.StringData.CustomData1 = "foo"
	appEventCustom.CustomProperties.StringData.CustomData2 = "bar"

	client.SendEvent(appEventCustom)*/
}

type TestSession struct {
	client   *trakerr.TrakerrClient
	appEvent *trakerr.AppEvent
}

type TestError struct {
}

//BufferOverflowError ...
func (testError *TestError) BufferOverflowError(buf []int, i int, session TestSession) (x int) {
	//defer client.Recover()
	defer session.client.RecoverWithAppEvent("Error", session.appEvent)

	x = buf[i]
	return x
}
