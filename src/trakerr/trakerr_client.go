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

//makeTimestamp is a package level function that formats the time now into a Trakerr readable format.
func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

//getTextFromLine is a package level function that takes a text in a string and
//returns the string in between the end of the prefix string and the start of the suffix string.
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

//TrakerrClient is the class that sends events to Trakerr.
//In a normal use case, Trakerr populates the field below with default values, but can be set after constuction to a custom value.
//The discription for the fields are below, and the struct associated methods below that.
type TrakerrClient struct {
	apiKey                     string
	contextAppVersion          string
	contextDeploymentStage     string
	contextEnvLanguage         string
	contextEnvName             string
	contextEnvVersion          string
	contextEnvHostname         string
	contextAppOS               string
	contextAppOSVersion        string
	contextAppOSBrowser        string
	contextAppOSBrowserVersion string
	contextDataCenter          string
	contextDataCenterRegion    string
	eventsAPI                  EventsApi
	eventTraceBuilder          EventTraceBuilder
}

//apiKey is your API key string.
//contextAppVersion is the version of the application.
//contextDeploymentStage is the deployment stage of the application.
//contextEnvLanguage is the constant string representing the language the application is in.
//contextEnvName is the OS and Arch name the compiler is targeting for the application.
//contextEnvVersion is the version of golang the application is run on.
//contextEnvHostname is hostname of the pc running the application.
//contextAppOS is the OS the application is running on.
//contextAppOSVersion is the version of the OS the application is running on.
//contextAppBrowser is an optional string browser name the application is running on.
//contextAppBrowserVersion is an optional string browser version the application is running on.
//contextDatacenter is the optional datacenter the code may be running on.
//contextDatacenterRegion is the optional datacenter region the code may be running on.

// NewTrakerrClient creates a new TrakerrClient and return it with the data.
// Most parameters are optional i.e. empty (pass "" to use defaults) with the exception of apiKey which is required.
func NewTrakerrClient(
	apiKey string,
	contextAppVersion string,
	contextDeploymentStage string) *TrakerrClient {

	if contextDeploymentStage == "" {
		contextDeploymentStage = "development"
	}
	if contextAppVersion == "" {
		contextAppVersion = "1.0"
	}

	contextEnvLanguage := "GoLang"
	//Go is a compiled language; the interpreter doesn't matter.
	//Instead it seems to compile profiles for OS (for the OS depended libs) and Arch.
	//I've provided the OS and arch target of the running program.
	contextEnvName := runtime.GOOS + "/" + runtime.GOARCH
	contextEnvVersion := runtime.Version()
	contextEnvHostname, _ := os.Hostname()

	contextAppOS := runtime.GOOS
	var contextAppOSVersion string
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

	var eventsAPI EventsApi
	eventsAPI = *NewEventsApi()

	return &TrakerrClient{
		apiKey:                  apiKey,
		contextAppVersion:       contextAppVersion,
		contextDeploymentStage:  contextDeploymentStage,
		contextEnvLanguage:      contextEnvLanguage,
		contextEnvName:          contextEnvName,
		contextEnvVersion:       contextEnvVersion,
		contextEnvHostname:      contextEnvHostname,
		contextAppOS:            contextAppOS,
		contextAppOSVersion:     contextAppOSVersion,
		contextDataCenter:       "",
		contextDataCenterRegion: "",
		eventsAPI:               eventsAPI,
		eventTraceBuilder:       EventTraceBuilder{}}
}

//NewAppEvent returns an AppEvent pointer with the classification eventType and eventMessage filled.
func (trakerrClient *TrakerrClient) NewAppEvent(loglevel string, classification string, eventType string, eventMessage string) *AppEvent {
	loglevel = strings.ToLower(loglevel)

	visitedURL := map[string]bool{
		"info":    true,
		"debug":   true,
		"warn":    true,
		"warning": true,
		"fatal":   true,
	}
	if !visitedURL[loglevel] {
		loglevel = "error"
	}
	if classification == "" {
		classification = "issue"
	}
	if eventType == "" {
		eventType = "unknown"
	}
	if eventMessage == "" {
		eventMessage = "unknown"
	}
	return trakerrClient.FillDefaults(&AppEvent{Classification: classification, EventType: eventType, EventMessage: eventMessage})
}

//NewEmptyEvent returns a Appevent pointer which is empty. If the AppEvent is passed into a defer later, classification, eventType, and eventMessage
//will be filled by the error parsing.
func (trakerrClient *TrakerrClient) NewEmptyEvent() *AppEvent {
	return trakerrClient.NewAppEvent("", "", "", "")
}

//SendEvent sends the event to trakerr.
func (trakerrClient *TrakerrClient) SendEvent(appEvent *AppEvent) (*APIResponse, error) {
	return trakerrClient.eventsAPI.EventsPost(*trakerrClient.FillDefaults(appEvent))
}

//SendError outward facing method that creates an event and takes a classification and an error.
func (trakerrClient *TrakerrClient) SendError(loglevel string, classification string, err interface{}) {
	trakerrClient.SendErrorWithSkip(err, loglevel, classification, 4)
}

//SendErrorWithSkip internal method that handles creating an app event and gets the stacktrace before sending.
func (trakerrClient *TrakerrClient) SendErrorWithSkip(err interface{}, loglevel string, classification string, skip int) (*APIResponse, error) {
	appEvent := trakerrClient.CreateAppEventFromErrorWithSkip(err, loglevel, classification, skip+1)

	return trakerrClient.eventsAPI.EventsPost(*appEvent)
}

//CreateAppEventFromError internal method that provides some default values for CreateAppEventFromErrorWithSkip.
func (trakerrClient *TrakerrClient) CreateAppEventFromError(loglevel string, classification string, err interface{}) *AppEvent {
	return trakerrClient.CreateAppEventFromErrorWithSkip(err, loglevel, classification, 4)

}

//CreateAppEventFromErrorWithSkip internal method which calls eventTraceBuilder to parse the stacktrace and creates an app event with it.
//Pass "" to use the default value classification
func (trakerrClient *TrakerrClient) CreateAppEventFromErrorWithSkip(err interface{}, loglevel string, classification string, skip int) *AppEvent {
	stacktrace := trakerrClient.eventTraceBuilder.GetEventTraces(err, 50, skip+1)
	event := trakerrClient.NewAppEvent(loglevel, classification, fmt.Sprintf("%T", err), fmt.Sprint(err))

	result := trakerrClient.FillDefaults(event)
	result.EventStacktrace = stacktrace
	return result
}

//AddStackTraceToAppEvent internal method to add a stack trace to an already exisiting AppEvent.
//Useful for creating your app event first to populate custom data.
//appEvent's EventType and EventMessaage are filled with the details from the error.
func (trakerrClient *TrakerrClient) AddStackTraceToAppEvent(appEvent *AppEvent, err interface{}, skip int) {
	stacktrace := trakerrClient.eventTraceBuilder.GetEventTraces(err, 50, skip+1)
	var event = appEvent
	if event.EventType == "" || event.EventMessage == "unknown" {
		event.EventType = fmt.Sprintf("%T", err)
	}
	if event.EventMessage == "" || event.EventMessage == "unknown" {
		event.EventMessage = fmt.Sprint(err)
	}

	event.EventStacktrace = stacktrace
}

//Recover recovers from a panic and sends the error to Trakerr. Creates the AppEvent
//Use in a Defer statement. The loglevel is the the string classifiction of the error (ie: "Error", "Info", ect).
func (trakerrClient *TrakerrClient) Recover(loglevel string, classification string) {
	if err := recover(); err != nil {
		response, apierr := trakerrClient.SendErrorWithSkip(err, loglevel, classification, 4)
		if response.StatusCode > 399 {
			fmt.Println(response.Status)
		}
		if apierr != nil {
			panic(apierr)
		}
	}
}

//RecoverWithAppEvent recovers from a panic and sends the error to Trakerr from a defer statement.
//This function takes in an AppEvent so could popultate the AppEvent with custom data and then attach the err from the defer.
func (trakerrClient *TrakerrClient) RecoverWithAppEvent(appEvent *AppEvent) {
	if err := recover(); err != nil {
		trakerrClient.AddStackTraceToAppEvent(appEvent, err, 4)
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
func (trakerrClient *TrakerrClient) Notify(loglevel string, classification string) {
	if err := recover(); err != nil {
		response, apierr := trakerrClient.SendErrorWithSkip(err, loglevel, classification, 4)
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
//so that the panic can be picked up by the program error handler. Use in a Defer statement.
//This function takes in an AppEvent so could popultate the AppEvent with custom data and then attach the err from the defer.
func (trakerrClient *TrakerrClient) NotifyWithAppEvent(appEvent *AppEvent) {
	if err := recover(); err != nil {
		trakerrClient.AddStackTraceToAppEvent(appEvent, err, 4)
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

//FillDefaults Populates the appevent with the TrakerrClient defaults.
func (trakerrClient *TrakerrClient) FillDefaults(appEvent *AppEvent) *AppEvent {
	if appEvent.ApiKey == "" {
		appEvent.ApiKey = trakerrClient.apiKey
	}
	if appEvent.ContextAppVersion == "" {
		appEvent.ContextAppVersion = trakerrClient.contextAppVersion
	}
	if appEvent.DeploymentStage == "" {
		appEvent.DeploymentStage = trakerrClient.contextDeploymentStage
	}

	if appEvent.ContextEnvLanguage == "" {
		appEvent.ContextEnvLanguage = trakerrClient.contextEnvLanguage
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
