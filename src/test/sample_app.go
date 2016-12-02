package main

import (
	"trakerr"
	"errors"
)

func main() {
	client := trakerr.NewTrakerrClientWithDefaults(
		"ceba200baf79b1b5e9dc73d4054d6c9618388477122",
		"http://192.168.0.117:3000/api/v1",
		"1.0",
		"development")
	err := errors.New("Something bad happened here")

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

	client.SendEvent(appEvent)
}