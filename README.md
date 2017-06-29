# Trakerr - Go API client

Get your application events and errors to Trakerr via the *Trakerr API*.

You can send both errors and non-errors (plain log statements, for example) to Trakerr with this API.

## Overview

The options presented in the integration guides are primarily oriented around sending errors or warnings without additional parameters.
**Option-4 in the detailed integration guide** describes how you could send a non-error (or any log statement) along with additional parameters.

The SDK takes performance impact seriously and all communication between the SDK <=> Trakerr avoids blocking the calling function. The SDK also applies asynchronous patterns where applicable.

A Trakerr *Event* is what is created by this SDK and is then sent to Trakerr for Trakerr to capture. A Trakerr *Event* can consist of various parameters as described here in [Model.AppEvent](https://github.com/trakerr-io/trakerr-csharp/blob/master/generated/docs/AppEvent.md).
Some of these parameters are populated by default and others are optional and can be supplied by you.

Since some of these parameters are common across all event's, the API has the option of setting these on the
TrakerrClient instance (described towards the bottom) and offers a factory API for creating AppEvent's.

### Key AppEvent Properties

#### Log Level, Event Type and Classification
* **Log Level** This enum specifies the logging level to be used for this event ('debug','info','warning','error' or 'fatal')
* **Event Type** This defines the type of event or logger name. This is automatically set for errors.
* **Classification** This is a user settable property that controls how the events are grouped. Defaults to 'Issue'. Set this to a different value to group this event in a different group.

#### Event User, Event Session and Correlation ID
* **Event User** This is the user that is associated with this event. This can be any user data or could be encrypted if privacy is required.
* **Event Session** This is any session specific information associated with this event.
* **Cross App Correlation ID** This is an additional ID that can be used for cross-application correlation of the same event.

#### Operation Time
* **Operation Time** This property in milliseconds measures the operation time for this specific event.

#### Custom properties and segments
In addition to the above, you can use custom properties and segments to send custom event, performance data. These
can then be visualized in Trakerr's dashboards.

### Requirements
go version 1.7.5+


## Installation
From the root directory issue the following
```bash
    go get github.com/trakerr-io/trakerr-go/src/trakerr
```

And then in your imports add this

```golang
    import "github.com/trakerr-io/trakerr-go/src/trakerr"

```

## Detailed Integration Guide
There are a few options (illustrated below with comment Option-#) to send events to Trakerr. If you would like to generate some sample events quickly and see the API in action, you can:

```bash
go get github.com/trakerr-io/trakerr-go/src/test
```
and then cd into the the folder. Once there, you can simply do:

```bash
go run sample_app.go <api-key>
```
To generate an error and get started on the site.

### Creating a new client
Once you have gotten the library with the steps above you can import the trakerr library. You may also wish to import the error standard library. Then you can initialize TrakerrClient from any method:

```golang
package main

import (
	"github.com/trakerr-io/trakerr-go/src/trakerr"
	"errors"
)

func main() {
	client := trakerr.NewTrakerrClient(
		"<api-key>",
		"<your app version>",
		"<your development stage>")
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
	appEvent *trakerr.AppEvent //store as many AppEvents as you need to cover.
}
```

And initializing it simply with the above methods.

```golang
ts := TestSession{client, appEvent}
```

We can then simply pass the code section relevent AppEvent from the struct when we want to keep a precautionary defer through TrakerrClient's recovery methods.

```golang
defer ts.client.RecoverWithAppEvent("Error", ts.appEvent)
```

Recover catches the panic and recover, while sending the error to Trakerr. If you wish to handle the error your own way,

```golang
defer ts.client.NotifyWithAppEvent("Error", ts.appEvent)
```

will catch the error, send it to Trakerr and then repanic in the same method.


### Option-2: Send an error to trakerr programmatically
You can manually send an error without using the panic subroutines. Create a new error manually as a result of an action and then pass it to TrakerrClient's `SendError()` function.

```golang
	err := errors.New("Something bad happened here")

	// Option-1: send error
	client.SendError("Error", "", err)
```

### Option-3: Send an error to trakerr programmatically with custom properties
You can follow the above steps, but instead pass the error into `CreateAppEventFromError()`. This will allow you to send custom properties to the app event before you send it.

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
Sending a non-error uses a similar process as above, but skips the step with the error. Be sure to fill in the type and message when sending a non-error!

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

The TrakerrClient struct however has a lot of exposed properties. The benefit to setting these after you create the TrakerrClient is that AppEvent will default it's values against the TrakerClient that created it. This way if there is a value that all your AppEvents uses, and the constructor default value currently doesn't suit you; it may be easier to change it in TrakerrClient as it will become the default value for all AppEvents created after. A lot of these are populated by default value by the constructor, but you can populate them with whatever string data you want. The following table provides an in depth look at each of those.

Name | Type | Description | Notes
------------ | ------------- | -------------  | -------------
**apiKey** | **string** | API key generated for the application | 
**contextAppVersion** | **string** | Application version information. | Default value: `1.0`
**contextDevelopmentStage** | **string** | One of development, staging, production; or a custom string. | Default Value: `development`
**contextEnvLanguage** | **string** | Constant string representing the language the application is in. | Default value: "Golang"
**contextEnvName** | **string** | OS and Arch name the compiler is targeting for the application. | Default Value: runtime.GOOS + " " + runtime.GOARCH
**contextEnvVersion** | **string** | Version of the go runtime the program is compiled on. | Default Value: runtime.Version()
**contextEnvHostname** | **string** | Hostname or ID of environment. | Default value: os.hostname()
**contextAppOS** | **string** | OS the application is running on. | Default value: OS name (ie. Windows, MacOS).
**contextAppOSVersion** | **string** | OS Version the application is running on. | Default value: OS Version.
**contextAppOSBrowser** | **string** | An optional string browser name the application is running on. | Defaults to `empty string` (`""`)
**contextAppOSBrowserVersion** | **string** | An optional string browser version the application is running on. | Defaults to `empty string` (`""`)
**contextDataCenter** | **string** | Data center the application is running on or connected to. | Defaults to `empty string` (`""`)
**contextDataCenterRegion** | **string** | Data center region. | Defaults to `empty string` (`""`)

## Documentation For Models

 - [AppEvent](https://github.com/trakerr-io/trakerr-go/blob/master/src/trakerr/docs/AppEvent.md)



