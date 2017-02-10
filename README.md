# Go API client for trakerr_client

Get your application events and errors to Trakerr via the *Trakerr API*.

## Overview
- API version: 1.0.0
- Package version: 1.0.0

## Installation
From the root directory issue the following
```
    go get github.com/trakerr-io/trakerr-go/src/trakerr
```

And then in your imports add this

```
    "github.com/trakerr-io/trakerr-go/src/trakerr"

```
## Getting Started

There are a few options (illustrated below with comment Option-#) to send events to Trakerr. The easiest of
which is to send only errors to Trakerr (Option-1).

### Creating a new client


```$golang
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
```$golang
	err := errors.New("Something bad happened here")

	// Option-1: send error
	client.SendError(err)
```

### Option-2: Send an error to trakerr programmatically with custom properties
```$golang
	err := errors.New("Something bad happened here")

	// Option-2: send error with custom properties
	appEventWithErr := client.CreateAppEventFromError(err)

	// set any custom data on appEvent
	appEventWithErr.CustomProperties.StringData.CustomData1 = "foo"
	appEventWithErr.CustomProperties.StringData.CustomData2 = "bar"

	client.SendEvent(appEventWithErr)
```

### Option-3: Send an event including non-exceptions to Trakerr
```$golang
	err := errors.New("Something bad happened here")

	// Option-3: send event manually
	appEvent := client.NewAppEvent("Info", "SomeType", "SomeMessage")

	// set any custom data on appEvent
	appEvent.CustomProperties.StringData.CustomData1 = "foo"
	appEvent.CustomProperties.StringData.CustomData2 = "bar"

	client.SendEvent(appEvent)
```

## Documentation For Models

 - [AppEvent](https://github.com/trakerr-io/trakerr-go/blob/master/src/trakerr/docs/AppEvent.md)



