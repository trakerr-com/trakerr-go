# Go API client for trakerr_client

Get your application events and errors to Trakerr via the *Trakerr API*.

## Overview
- API version: 1.0.0
- Package version: 1.0.0

## Installation
From the root directory issue the following
```bash
    go get github.com/trakerr-io/trakerr-go/src/trakerr
```

And then in your imports add this

```golang
    import "github.com/trakerr-io/trakerr-go/src/trakerr"

```
## Getting Started

There are a few options (illustrated below with comment Option-#) to send events to Trakerr. If you would like to generate some sample events quickly and see the API in action, you can:

```bash
go get github.com/trakerr-io/trakerr-go/src/test
```
and then cd into the the folder. Once there, you can simple do:

```bash
go run sample_app.go <API key here>
```
To generate an error and get started on the site.

### Creating a new client


```golang
package main

import (
	"github.com/trakerr-io/trakerr-go/src/trakerr"
	"errors"
)

func main() {
	client := trakerr.NewTrakerrClientWithDefaults(
		"<replace with your API key>",
		"1.0",
		"development")
    ...
}
```

### Option-1: Catch errors automatically with a defer
Once you've created a client, you can set up an exception prepared for an area which may cause panics:

```golang
appEvent := client.NewEmptyEvent()
// set any custom data on appEvent
appEvent.CustomProperties.StringData.CustomData1 = "foo"
appEvent.CustomProperties.StringData.CustomData2 = "bar"
appEvent.EventUser = "john@trakerr.io"
appEvent.EventSession = "12"
```

We suggest storing these in a struct for global error handling:
```golang
type TestSession struct {
	client   *trakerr.TrakerrClient
	appEvent *trakerr.AppEvent
}
```

And initializing it simply with the above methods.

```golang
ts := TestSession{client, appEvent}
```

We can then simply access the struct when we want to keep a precautionary defer. And then calling one of the methods that recover from the panic in Traker Client.

```golang
defer session.client.RecoverWithAppEvent("Error", session.appEvent)
```

Recover catches the panic and recover, while sending the error to Trakerr. If you wish to handle the error your own way,

```golang
defer session.client.NotifyWithAppEvent("Error", session.appEvent)
```

will catch the error, send it to Trakerr and then repanic in the same method.


### Option-2: Send an error to trakerr programmatically
```golang
	err := errors.New("Something bad happened here")

	// Option-1: send error
	client.SendError("Error", err)
```

### Option-3: Send an error to trakerr programmatically with custom properties
```golang
	err := errors.New("Something bad happened here")

	// Option-2: send error with custom properties
	appEventWithErr := client.CreateAppEventFromError("Error", err)

	// set any custom data on appEvent
	appEventWithErr.CustomProperties.StringData.CustomData1 = "foo"
	appEventWithErr.CustomProperties.StringData.CustomData2 = "bar"

	client.SendEvent(appEventWithErr)
```

### Option-4: Send an event including non-exceptions to Trakerr
```golang
	err := errors.New("Something bad happened here")

	// Option-3: send event manually
	appEvent := client.NewAppEvent("Info", "SomeType", "SomeMessage")

	// set any custom data on appEvent
	appEvent.CustomProperties.StringData.CustomData1 = "foo"
	appEvent.CustomProperties.StringData.CustomData2 = "bar"

	client.SendEvent(appEvent)
```

## Initializing Trakerr
Due to the nature of golang, Trakerr can be initialized in one of two ways. The first way with defaults, is relatively self explanatory.
```golang
func NewTrakerrClientWithDefaults(
	apiKey string,
	contextAppVersion string,
	contextEnvName string) *TrakerrClient
```
The ContextEnvName name is intended to be used as a string identifier as to what your codebase is for; release, development, prototype. You can use it for whatever you denote as useful. The ContextApp Version is useful for a codebase version identifier, or perhaps some other useful metric for the error.

Looking at the second call we're exposed to a lot of what the first call defaults.

```golang
func NewTrakerrClient(
	apiKey string,
	url string,
	contextAppVersion string,
	contextEnvName string,
	contextEnvVersion string,
	contextEnvHostname string,
	contextAppOS string,
	contextAppOSVersion string,
	contextDataCenter string,
	contextDataCenterRegion string) *TrakerrClient
```
A lot of these are populated by default value by the first call, but you can populate them with whatever string data you want. Here is an indepth look at each of those.

Name | Type | Description | Notes
------------ | ------------- | -------------  | -------------
**apiKey** | **string** | API key generated for the application | 
**url** | **string** |(optional) The URL to send to. You will mostly want to leave this empty string to send to trakerr. | [optional if passed `""`]
**contextAppVersion** | **string** | (optional) application version information. | [optional if passed `""`] Default value: "1.0" 
**contextEnvName** | **string** | (optional) one of development, staging, production; or a custom string. | [optional if passed `""`] Default Value: "develoment"
**contextEnvHostname** | **string** | (optional) hostname or ID of environment. | [optional if passed `""`] Default value: os.hostname()
**contextAppOS** | **string** | (optional) OS the application is running on. | [optional if passed `""`] Default value: OS name (ie. Windows, MacOS) (Currently being reworked).
**contextAppOSVersion** | **string** | (optional) OS Version the application is running on. | [optional if passed `""`] Default value: System architecture string (Currently being reworked).
**contextDataCenter** | **string** | (optional) Data center the application is running on or connected to. | [optional if passed `""`]
**contextDataCenterRegion** | **string** | (optional) Data center region. | [optional if passed `""`]


## Documentation For Models

 - [AppEvent](https://github.com/trakerr-io/trakerr-go/blob/master/src/trakerr/docs/AppEvent.md)



