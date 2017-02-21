//Package trakerr gives the client access to client side constructors for initializing
//and using trakerr.
package trakerr

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

//TrakerrClient ...
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
	eventsAPI               EventsAPI
	eventTraceBuilder       EventTraceBuilder
}

//NewTrakerrClientWithDefaults creates a new TrakerrClient and return it with the data.
// Most parameters are optional i.e. empty (pass "" to use defaults) with the exception of apiKey which is required.
func NewTrakerrClientWithDefaults(
	apiKey string,
	contextAppVersion string,
	contextEnvName string) *TrakerrClient {
	return NewTrakerrClient(apiKey, "", contextAppVersion, contextEnvName, "", "", "", "", "", "")
}

// NewTrakerrClient creates a new TrakerrClient and return it with the data.
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

	if contextEnvName == "" {
		contextEnvName = "development"
	}
	if contextAppVersion == "" {
		contextAppVersion = "1.0"
	}
	if contextEnvHostname == "" {
		contextEnvHostname, _ = os.Hostname()
	}

	if contextAppOS == "" {
		contextAppOS = runtime.GOOS
		switch contextAppOS {
		case "windows": // Keep an eye out for OS's using carrage returns
			//cmd1 := exec.Command("systeminfo")
			cmd := exec.Command("cmd", "/C", "systeminfo") //exec.Command("findstr", "/C:\"OS Name\"")
			var out bytes.Buffer
			cmd.Stdout = &out

			err1 := cmd.Run()
			if err1 != nil {
				contextAppOS = runtime.GOOS
				contextAppOSVersion = "N/A (arch:" + runtime.GOARCH + ")"
			} else {
				var output = out.String()
				contextAppOS = getTextFromLine(output, "OS Name:", "\n")
				version := getTextFromLine(output, "OS Version:", "\n")
				versionstringsplit := strings.Split(version, " ")

				if len(versionstringsplit) >= 1 {
					contextAppOSVersion = versionstringsplit[0]
				} else {
					contextAppOSVersion = version
				}
			}

		case "darwin": // ...
			cmd := exec.Command("bash", "-c", "system_profiler SPSoftwareDataType") //exec.Command("findstr", "/C:\"OS Name\"")
			var out bytes.Buffer
			cmd.Stdout = &out

			err1 := cmd.Run()
			if err1 != nil {
				contextAppOS = runtime.GOOS
				contextAppOSVersion = "N/A (arch:" + runtime.GOARCH + ")"
			} else {
				var output = out.String()

				contextAppOS = getTextFromLine(output, "System Version:", "(")
				contextAppOSVersion = getTextFromLine(output, "Kernel Version:", "\n")
			}

		default:
			cmd := exec.Command("bash", "-c", "uname -s") //Uname -r and -s
			var out bytes.Buffer
			cmd.Stdout = &out

			err1 := cmd.Run()
			var output string
			if err1 != nil {
				contextAppOS = runtime.GOOS
			} else {
				output = out.String()

				contextAppOS = strings.Trim(output, " \r\n")
			}

			cmd = exec.Command("bash", "-c", "uname -r") //Uname -r and -s
			err1 = cmd.Run()
			if err1 != nil {
				contextAppOSVersion = "N/A (arch:" + runtime.GOARCH + ")"
			} else {
				output = out.String()

				contextAppOSVersion = strings.Trim(output, " \r\n")
			}

		}

	}
	var eventsAPI EventsAPI

	if url != "" {
		eventsAPI = *NewEventsAPIWithBasePath(url)
	} else {
		eventsAPI = *NewEventsAPI()
	}

	return &TrakerrClient{
		apiKey:                  apiKey,
		url:                     url,
		contextAppVersion:       contextAppVersion,
		contextEnvName:          contextEnvName,
		contextEnvVersion:       contextEnvVersion,
		contextEnvHostname:      contextEnvHostname,
		contextAppOS:            contextAppOS,
		contextAppOSVersion:     contextAppOSVersion,
		contextDataCenter:       contextDataCenter,
		contextDataCenterRegion: contextDataCenterRegion,
		eventsAPI:               eventsAPI,
		eventTraceBuilder:       EventTraceBuilder{}}
}

func getTextFromLine(text string, prefix string, suffix string) string {
	var startindex = strings.Index(text, prefix)
	if startindex == -1 {
		return ""
	}
	var newstringfromprefix = text[startindex+len(prefix):]
	var endindex = strings.Index(newstringfromprefix, suffix)
	if endindex == -1 {
		return ""
	}

	return strings.Trim(newstringfromprefix[0:endindex], " \n\r")
}

//NewAppEvent ...
func (trakerrClient *TrakerrClient) NewAppEvent(classification string, eventType string, eventMessage string) *AppEvent {
	if classification == "" {
		classification = "Error"
	}
	if eventType == "" {
		eventType = "unknown"
	}
	if eventMessage == "" {
		eventMessage = "unknown"
	}
	return trakerrClient.FillDefaults(&AppEvent{Classification: classification, EventType: eventType, EventMessage: eventMessage})
}

//NewErrorEvent ...
func (trakerrClient *TrakerrClient) NewEmptyEvent() *AppEvent {
	return trakerrClient.NewAppEvent(" ", "", "")
}

//SendEvent ...
func (trakerrClient *TrakerrClient) SendEvent(appEvent *AppEvent) (*APIResponse, error) {
	return trakerrClient.eventsAPI.EventsPost(*trakerrClient.FillDefaults(appEvent))
}

//SendError ...
func (trakerrClient *TrakerrClient) SendError(classification string, err interface{}) {
	trakerrClient.SendErrorWithSkip(err, classification, 4)
}

//SendErrorWithSkip ...
func (trakerrClient *TrakerrClient) SendErrorWithSkip(err interface{}, classification string, skip int) (*APIResponse, error) {
	appEvent := trakerrClient.CreateAppEventFromErrorWithSkip(err, classification, skip+1)

	return trakerrClient.eventsAPI.EventsPost(*appEvent)
}

//CreateAppEventFromError ...
func (trakerrClient *TrakerrClient) CreateAppEventFromError(classification string, err interface{}) *AppEvent {
	return trakerrClient.CreateAppEventFromErrorWithSkip(err, classification, 4)

}

//CreateAppEventFromErrorWithSkip ...
func (trakerrClient *TrakerrClient) CreateAppEventFromErrorWithSkip(err interface{}, classification string, skip int) *AppEvent {
	stacktrace := trakerrClient.eventTraceBuilder.GetEventTraces(err, 50, skip+1)
	event := AppEvent{}
	event.EventType = fmt.Sprintf("%T", err)
	event.EventMessage = fmt.Sprint(err)
	event.Classification = classification

	result := trakerrClient.FillDefaults(&event)
	event.EventStacktrace = stacktrace
	return result
}

//AddStackTraceToAppEvent ...
func (trakerrClient *TrakerrClient) AddStackTraceToAppEvent(appEvent *AppEvent, classification string, err interface{}, skip int) {
	stacktrace := trakerrClient.eventTraceBuilder.GetEventTraces(err, 50, skip+1)
	var event = appEvent
	if event.EventType == "" || event.EventMessage == "unknown" {
		event.EventType = fmt.Sprintf("%T", err)
	}
	if event.EventMessage == "" || event.EventMessage == "unknown" {
		event.EventMessage = fmt.Sprint(err)
	}
	event.Classification = classification

	event.EventStacktrace = stacktrace
}

//Recover recovers from a panic and sends the error to Trakerr.
//Use in a Defer statement.
func (trakerrClient *TrakerrClient) Recover(classification string) {
	if err := recover(); err != nil {
		response, apierr := trakerrClient.SendErrorWithSkip(err, classification, 4)
		if response.StatusCode > 399 {
			fmt.Println(response.Status)
		}
		if apierr != nil {
			panic(apierr)
		}
	}
}

//RecoverWithAppEvent ...
func (trakerrClient *TrakerrClient) RecoverWithAppEvent(classification string, appEvent *AppEvent) {
	if err := recover(); err != nil {
		trakerrClient.AddStackTraceToAppEvent(appEvent, classification, err, 4)
		response, apierr := trakerrClient.SendEvent(appEvent)
		if response.StatusCode > 399 {
			fmt.Println(response.Status)
		}
		if apierr != nil {
			panic(apierr)
		}

	}
}

//Notify recovers from an error and then repanics after sending the error to Trakerr,
//so that the panic can be picked up by the program error handler.
//Use in a Defer statement.
func (trakerrClient *TrakerrClient) Notify(classification string) {
	if err := recover(); err != nil {
		response, apierr := trakerrClient.SendErrorWithSkip(err, classification, 4)
		if response.StatusCode > 399 {
			fmt.Println(response.Status)
		}
		if apierr != nil {
			panic(apierr)
		}
		panic(err)
	}
}

//NotifyWithAppEvent recovers from an error and then repanics after sending the error to Trakerr,
//so that the panic can be picked up by the program error handler.
//Use in a Defer statement.
func (trakerrClient *TrakerrClient) NotifyWithAppEvent(classification string, appEvent *AppEvent) {
	if err := recover(); err != nil {
		trakerrClient.AddStackTraceToAppEvent(appEvent, classification, err, 4)
		response, apierr := trakerrClient.SendEvent(appEvent)
		if response.StatusCode > 399 {
			fmt.Println(response.Status)
		}
		if apierr != nil {
			panic(apierr)
		}
		panic(err)
	}
}

//FillDefaults ...
func (trakerrClient *TrakerrClient) FillDefaults(appEvent *AppEvent) *AppEvent {
	if appEvent.ApiKey == "" {
		appEvent.ApiKey = trakerrClient.apiKey
	}

	if appEvent.ContextAppVersion == "" {
		appEvent.ContextAppVersion = trakerrClient.contextAppVersion
	}

	if appEvent.ContextEnvName == "" {
		appEvent.ContextEnvName = trakerrClient.contextEnvName
	}
	if appEvent.ContextEnvVersion == "" {
		appEvent.ContextEnvVersion = trakerrClient.contextEnvVersion
	}
	if appEvent.ContextEnvHostname == "" {
		appEvent.ContextEnvHostname = trakerrClient.contextEnvHostname
	}

	if appEvent.ContextAppOS == "" {
		appEvent.ContextAppOS = trakerrClient.contextAppOS
		appEvent.ContextAppOSVersion = trakerrClient.contextAppOSVersion
	}

	if appEvent.ContextDataCenter == "" {
		appEvent.ContextDataCenter = trakerrClient.contextDataCenter
	}
	if appEvent.ContextDataCenterRegion == "" {
		appEvent.ContextDataCenterRegion = trakerrClient.contextDataCenterRegion
	}

	if appEvent.EventTime <= 0 {
		appEvent.EventTime = makeTimestamp()
	}
	return appEvent
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
