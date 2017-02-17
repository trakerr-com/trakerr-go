package main

import "github.com/trakerr-io/trakerr-go/src/trakerr"

func main() {
	client := trakerr.NewTrakerrClientWithDefaults(
		"API Key Here",
		"1.0",
		"development")
	/*err := errors.New("Something bad happened here")

	// Option-1: send error
	client.SendError(err)

	// Option-2: send error with custom properties
	appEventWithErr := client.CreateAppEventFromError(err)

	// set any custom data on appEvent
	appEventWithErr.CustomProperties.StringData.CustomData1 = "foo"
	appEventWithErr.CustomProperties.StringData.CustomData2 = "bar"

	client.SendEvent(appEventWithErr)

	// Option-3: send event manually
	appEvent := client.NewAppEvent("Info", "SomeType", "SomeMessage")

	// set any custom data on appEvent
	appEvent.CustomProperties.StringData.CustomData1 = "foo"
	appEvent.CustomProperties.StringData.CustomData2 = "bar"

	client.SendEvent(appEvent)*/

	//Option 4: Global error handling
	buf := []int{1, 2, 3}
	te := TestError{}
	te.BufferOverflowError(buf, 4, client)
}

func GlobalHandler(client *trakerr.TrakerrClient) {
	if err := recover(); err != nil {

		appEvent := client.CreateAppEventFromError(err)
		// set any custom data on appEvent
		appEvent.CustomProperties.StringData.CustomData1 = "foo"
		appEvent.CustomProperties.StringData.CustomData2 = "bar"
		appEvent.EventUser = "John Doe"
		appEvent.EventSession = "12"

		client.SendEvent(appEvent)
		//panic(err) //Resets the panic on the local machine.
	}
}

type TestError struct {
}

//BufferOverflowError ...
func (testError *TestError) BufferOverflowError(buf []int, i int, client *trakerr.TrakerrClient) (x int) {
	defer GlobalHandler(client)
	//defer client.Recover()

	x = buf[i]
	return x
}
