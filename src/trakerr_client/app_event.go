/* 
 * Trakerr API
 *
 * Get your application events and errors to Trakerr via the *Trakerr API*.
 *
 * OpenAPI spec version: 1.0.0
 * 
 * Generated by: https://github.com/swagger-api/swagger-codegen.git
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package trakerr_client

type AppEvent struct {

	// API key generated for the application
	ApiKey string `json:"apiKey,omitempty"`

	// one of 'debug','info','warning','error' or a custom string
	Classification string `json:"classification,omitempty"`

	// type or event or error (eg. NullPointerException)
	EventType string `json:"eventType,omitempty"`

	// message containing details of the event or error
	EventMessage string `json:"eventMessage,omitempty"`

	// (optional) event time in ms since epoch
	EventTime int64 `json:"eventTime,omitempty"`

	EventStacktrace []InnerStackTrace `json:"eventStacktrace,omitempty"`

	// (optional) event user identifying a user
	EventUser string `json:"eventUser,omitempty"`

	// (optional) session identification
	EventSession string `json:"eventSession,omitempty"`

	// (optional) application version information
	ContextAppVersion string `json:"contextAppVersion,omitempty"`

	// (optional) one of 'development','staging','production' or a custom string
	ContextEnvName string `json:"contextEnvName,omitempty"`

	// (optional) version of environment
	ContextEnvVersion string `json:"contextEnvVersion,omitempty"`

	// (optional) hostname or ID of environment
	ContextEnvHostname string `json:"contextEnvHostname,omitempty"`

	// (optional) browser name if running in a browser (eg. Chrome)
	ContextAppBrowser string `json:"contextAppBrowser,omitempty"`

	// (optional) browser version if running in a browser
	ContextAppBrowserVersion string `json:"contextAppBrowserVersion,omitempty"`

	// (optional) OS the application is running on
	ContextAppOS string `json:"contextAppOS,omitempty"`

	// (optional) OS version the application is running on
	ContextAppOSVersion string `json:"contextAppOSVersion,omitempty"`

	// (optional) Data center the application is running on or connected to
	ContextDataCenter string `json:"contextDataCenter,omitempty"`

	// (optional) Data center region
	ContextDataCenterRegion string `json:"contextDataCenterRegion,omitempty"`

	CustomProperties CustomData `json:"customProperties,omitempty"`

	CustomSegments CustomData `json:"customSegments,omitempty"`
}
