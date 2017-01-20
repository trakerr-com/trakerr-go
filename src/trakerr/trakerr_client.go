package trakerr

import (
	"trakerr_client"
	"os"
	"runtime"
	"time"
	"fmt"
)

type TrakerrClient struct {
	apiKey                  string
	url                     string
	contextAppVersion       string
	contextEnvName          string
	contextEnvVersion       string
	contextEnvHostname      string
	contextAppOS            string
	contextAppOSVersion     string
	contextDataCenter       string
	contextDataCenterRegion string
	eventsApi               trakerr_client.EventsApi
	eventTraceBuilder       EventTraceBuilder
}

// Create a new TrakerrClient and return it with the data.
//
// Most parameters are optional i.e. empty (pass "" to use defaults) with the exception of apiKey which is required.
func NewTrakerrClientWithDefaults(
	apiKey string,
	contextAppVersion string,
	contextEnvName string) *TrakerrClient {
	return NewTrakerrClient(apiKey, "", contextAppVersion, contextEnvName, "", "", "", "", "", "")
}

// Create a new TrakerrClient and return it with the data.
//
// Most parameters are optional i.e. empty (pass "" to use defaults) with the exception of apiKey which is required.
// url is the location of the serverr service, if "" is passed it defaults to https://trakerr.io/api/v1
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
	contextDataCenterRegion string) *TrakerrClient {

	if contextEnvName == "" { contextEnvName = "development" }
	if contextAppVersion == "" { contextAppVersion = "1.0" }
	if contextEnvHostname == "" { contextEnvHostname, _ = os.Hostname() }

	if contextAppOS == "" {
		contextAppOS = runtime.GOOS
		contextAppOSVersion = "N/A (arch:" + runtime.GOARCH + ")"
	}
	var eventsApi trakerr_client.EventsApi

	if url != "" {
		eventsApi = *trakerr_client.NewEventsApiWithBasePath(url);
	} else {
		eventsApi = *trakerr_client.NewEventsApi()
	}

	return &TrakerrClient{
		apiKey: apiKey,
		url: url,
		contextAppVersion: contextAppVersion,
		contextEnvName: contextEnvName,
		contextEnvVersion: contextEnvVersion,
		contextEnvHostname: contextEnvHostname,
		contextAppOS: contextAppOS,
		contextAppOSVersion: contextAppOSVersion,
		contextDataCenter: contextDataCenter,
		contextDataCenterRegion: contextDataCenterRegion,
		eventsApi: eventsApi,
		eventTraceBuilder: EventTraceBuilder{} }
}

func (trakerrClient *TrakerrClient) NewAppEvent(classification string, eventType string, eventMessage string) *trakerr_client.AppEvent {
	if classification == "" { classification = "Error" }
	if eventType == "" { eventType = "unknown" }
	if eventMessage == "" { eventMessage = "unknown "}
	return trakerrClient.FillDefaults(&trakerr_client.AppEvent{Classification: classification, EventType:eventType, EventMessage: eventMessage })
}

func (trakerrClient *TrakerrClient) SendEvent(appEvent *trakerr_client.AppEvent) (*trakerr_client.APIResponse, error) {
	return trakerrClient.eventsApi.EventsPost(*trakerrClient.FillDefaults(appEvent))
}

func (trakerrClient *TrakerrClient) SendError(err interface{}) (*trakerr_client.APIResponse, error) {
	appEvent := trakerrClient.CreateAppEventFromError(err)

	return trakerrClient.eventsApi.EventsPost(*appEvent)
}

func (trakerrClient *TrakerrClient) CreateAppEventFromError(err interface{}) *trakerr_client.AppEvent {
	stacktrace := trakerrClient.eventTraceBuilder.GetEventTraces(err, 4)
	event := trakerr_client.AppEvent{}
	event.EventType = fmt.Sprintf("%T", err)
	event.EventMessage = fmt.Sprint(err)
	event.Classification = "Error"

	result := trakerrClient.FillDefaults(&event)
	event.EventStacktrace = stacktrace
	return result
}

func (trakerrClient *TrakerrClient) FillDefaults(appEvent *trakerr_client.AppEvent) *trakerr_client.AppEvent {
	if appEvent.ApiKey == "" {
		appEvent.ApiKey = trakerrClient.apiKey
	}

	if (appEvent.ContextAppVersion == "") {
		appEvent.ContextAppVersion = trakerrClient.contextAppVersion
	}

	if (appEvent.ContextEnvName == "") {
		appEvent.ContextEnvName = trakerrClient.contextEnvName
	}
	if (appEvent.ContextEnvVersion == "") {
		appEvent.ContextEnvVersion = trakerrClient.contextEnvVersion
	}
	if (appEvent.ContextEnvHostname == "") {
		appEvent.ContextEnvHostname = trakerrClient.contextEnvHostname
	}

	if (appEvent.ContextAppOS == "") {
		appEvent.ContextAppOS = trakerrClient.contextAppOS
		appEvent.ContextAppOSVersion = trakerrClient.contextAppOSVersion
	}

	if (appEvent.ContextDataCenter == "") {
		appEvent.ContextDataCenter = trakerrClient.contextDataCenter
	}
	if (appEvent.ContextDataCenterRegion == "") {
		appEvent.ContextDataCenterRegion = trakerrClient.contextDataCenterRegion
	}

	if (appEvent.EventTime <= 0) {
		appEvent.EventTime = makeTimestamp()
	}
	return appEvent
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}



