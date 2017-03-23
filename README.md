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
	client := trakerr.NewTrakerrClient(
		"<replace with your API key>",
		"1.0",
		"development")
    ...
}
```

### Option-1: Catch errors automatically with a defer
Once you've created a client, you can set up an exception prepared for an area which may cause panics:

```golang
appEvent := client.NewAppEvent("Error", "", "", "")
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
	client.SendError("Error", "", err)
```

### Option-3: Send an error to trakerr programmatically with custom properties
```golang
	err := errors.New("Something bad happened here")

	// Option-2: send error with custom properties
	appEventWithErr := client.CreateAppEventFromError("Error", "", err)

	// set any custom data on appEvent
	appEventWithErr.CustomProperties.StringData.CustomData1 = "foo"
	appEventWithErr.CustomProperties.StringData.CustomData2 = "bar"

	client.SendEvent(appEventWithErr)
```

### Option-4: Send an event including non-exceptions to Trakerr
```golang
	// Option-3: send event manually
	appEvent := client.NewAppEvent("Info", "custom classification", "SomeType", "SomeMessage")

	// set any custom data on appEvent
	appEvent.CustomProperties.StringData.CustomData1 = "foo"
	appEvent.CustomProperties.StringData.CustomData2 = "bar"

	client.SendEvent(appEvent)
```

## Initializing Trakerr
Due to the nature of golang, Trakerr is initalized to default values with the constructor.
```golang
func NewTrakerrClientWithDefaults(
	apiKey string,
	contextAppVersion string,
	contextEnvName string) *TrakerrClient
```
The contextEnvName name is intended to be used as a string identifier as to what your codebase is for; release, development, prototype. You can use it for whatever you denote as useful. The contextAppVersion is useful for a codebase version identifier, or perhaps some other useful metric for the error.

The TrakerrClient struct however has a lot of exposed properties. The benifit to setting these after you create the TrakerrClient is that AppEvent will default it's values against the TrakerClient that created it. This way if there is a value that all your AppEvents uses, and the constructor default value currently doesn't suit you; it may be easier to change it in TrakerrClient as it will become the default value for all AppEvents created after. A lot of these are populated by default value by the constructor, but you can populate them with whatever string data you want. The following table provides an in depth look at each of those.

Name | Type | Description | Notes
------------ | ------------- | -------------  | -------------
**apiKey** | **string** | API key generated for the application | 
**contextAppVersion** | **string** | Application version information. | Default value: "1.0" 
**contextDevelopmentStage** | **string** | One of development, staging, production; or a custom string. | Default Value: "develoment"
**contextEnvLanguage** | **string** | OS and Arch name the compiler is targeting for the application. | Default value: "Golang"
**contextEnvName** | **string** | Constant string representing the language the application is in. | Default Value: runtime.GOOS + " " + runtime.GOARCH
**contextEnvHostname** | **string** | Hostname or ID of environment. | Default value: os.hostname()
**contextAppOS** | **string** | OS the application is running on. | Default value: OS name (ie. Windows, MacOS).
**contextAppOSVersion** | **string** | OS Version the application is running on. | Default value: OS Version.
**contextAppOSBrowser** | **string** | An optional string browser name the application is running on. | Defaults to empty (`""`)
**contextAppOSBrowserVersion** | **string** | An optional string browser version the application is running on. | Defaults to empty (`""`)
**contextDataCenter** | **string** | Data center the application is running on or connected to. | Defaults to empty (`""`)
**contextDataCenterRegion** | **string** | Data center region. | Defaults to empty (`""`)

## Documentation For Models

 - [AppEvent](https://github.com/trakerr-io/trakerr-go/blob/master/src/trakerr/docs/AppEvent.md)



