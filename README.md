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

There are a few options (illustrated below with comment Option-#) to send events to Trakerr. The easiest of
which is to send only errors to Trakerr (Option-1).

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

### Option-1: Send an error to trakerr programmatically
```golang
	err := errors.New("Something bad happened here")

	// Option-1: send error
	client.SendError(err)
```

### Option-2: Send an error to trakerr programmatically with custom properties
```golang
	err := errors.New("Something bad happened here")

	// Option-2: send error with custom properties
	appEventWithErr := client.CreateAppEventFromError(err)

	// set any custom data on appEvent
	appEventWithErr.CustomProperties.StringData.CustomData1 = "foo"
	appEventWithErr.CustomProperties.StringData.CustomData2 = "bar"

	client.SendEvent(appEventWithErr)
```

### Option-3: Send an event including non-exceptions to Trakerr
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
client := trakerr.NewTrakerrClientWithDefaults(
		"API key",
		"ContextAppVersion",
		"ContextEnvName")
```
The ContextEnvName name is intended to be used as a string identifier as to what your codebase is for; release, development, prototype. You can use it for whatever you denote as useful. The ContextApp Version is useful for a codebase version identifier, or perhaps some other useful metric for the error.

Looking at the second call we're exposed to a lot of what the first call defaults.

```golang
    client := trakerr.NewTrakerrClient(
        "API Key Here",
        "URL",
        "ContextAppVersion",
        "ContextEnvName",
        "ContextEnvHostName",
        "ContextAppos",
        "ContextApposVersion",
        "ContextDataCenter",
        "ContextDataCenterRegion")
```
A lot of these are populated by default value by the first call, but you can populate them with whatever string data you want. Here is an indepth look at each of those.

Name | Type | Description | Notes
------------ | ------------- | -------------  | -------------
**ApiKey** | **string** | API key generated for the application | 
**URL** | **string** |(optional) The URL to send to. You will mostly want to leave this empty string to send to trakerr. | [optional if passed `""`] Default value: "1.0"
**ContextAppVersion** | **string** | (optional) application version information. | [optional if passed `""`] Default value: "1.0" 
**ContextEnvName** | **string** | (optional) one of development, staging, production; or a custom string. | [optional if passed `""`] Default Value: "develoment"
**ContextEnvHostname** | **string** | (optional) hostname or ID of environment. | [optional if passed `""`] Default value: os.hostname()
**ContextAppOS** | **string** | (optional) OS the application is running on. | [optional if passed `""`] Default value: OS name (ie. Windows, MacOS) (Currently being reworked).
**ContextAppOSVersion** | **string** | (optional) OS Version the application is running on. | [optional if passed `""`] Default value: System architecture string (Currently being reworked).
**ContextDataCenter** | **string** | (optional) Data center the application is running on or connected to. | [optional if passed `""`]
**ContextDataCenterRegion** | **string** | (optional) Data center region. | [optional if passed `""`]


## Documentation For Models

 - [AppEvent](https://github.com/trakerr-io/trakerr-go/blob/master/src/trakerr/docs/AppEvent.md)



