# AppEvent

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApiKey** | **string** | API key generated for the application | [default to null]
**Classification** | **string** | one of &#39;debug&#39;,&#39;info&#39;,&#39;warning&#39;,&#39;error&#39; or a custom string | [default to null]
**EventType** | **string** | type or event or error (eg. NullPointerException) | [default to null]
**EventMessage** | **string** | message containing details of the event or error | [default to null]
**EventTime** | **int64** | (optional) event time in ms since epoch | [optional] [default to null]
**EventStacktrace** | [**Stacktrace**](Stacktrace.md) |  | [optional] [default to null]
**EventUser** | **string** | (optional) event user identifying a user | [optional] [default to null]
**EventSession** | **string** | (optional) session identification | [optional] [default to null]
**ContextAppVersion** | **string** | (optional) application version information | [optional] [default to null]
**ContextEnvName** | **string** | (optional) one of &#39;development&#39;,&#39;staging&#39;,&#39;production&#39; or a custom string | [optional] [default to null]
**ContextEnvVersion** | **string** | (optional) version of environment | [optional] [default to null]
**ContextEnvHostname** | **string** | (optional) hostname or ID of environment | [optional] [default to null]
**ContextAppBrowser** | **string** | (optional) browser name if running in a browser (eg. Chrome) | [optional] [default to null]
**ContextAppBrowserVersion** | **string** | (optional) browser version if running in a browser | [optional] [default to null]
**ContextAppOS** | **string** | (optional) OS the application is running on | [optional] [default to null]
**ContextAppOSVersion** | **string** | (optional) OS version the application is running on | [optional] [default to null]
**ContextDataCenter** | **string** | (optional) Data center the application is running on or connected to | [optional] [default to null]
**ContextDataCenterRegion** | **string** | (optional) Data center region | [optional] [default to null]
**CustomProperties** | [**CustomData**](CustomData.md) |  | [optional] [default to null]
**CustomSegments** | [**CustomData**](CustomData.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


