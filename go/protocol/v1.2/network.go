package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
)

// Unique loader identifier.
type LoaderId string

// Unique request identifier.
type RequestId string

// Number of seconds since epoch.
type NetworkTimestamp int

// Request / response headers as keys / values of JSON object.
type Headers struct {
}

// Loading priority of a resource request.
type ConnectionType string

const ConnectionTypeNone ConnectionType = "none"
const ConnectionTypeCellular2g ConnectionType = "cellular2g"
const ConnectionTypeCellular3g ConnectionType = "cellular3g"
const ConnectionTypeCellular4g ConnectionType = "cellular4g"
const ConnectionTypeBluetooth ConnectionType = "bluetooth"
const ConnectionTypeEthernet ConnectionType = "ethernet"
const ConnectionTypeWifi ConnectionType = "wifi"
const ConnectionTypeWimax ConnectionType = "wimax"
const ConnectionTypeOther ConnectionType = "other"

// Represents the cookie's 'SameSite' status: https://tools.ietf.org/html/draft-west-first-party-cookies
type CookieSameSite string

const CookieSameSiteStrict CookieSameSite = "Strict"
const CookieSameSiteLax CookieSameSite = "Lax"

// Timing information for the request.
type ResourceTiming struct {
	RequestTime       int `json:"requestTime"`       // Timing's requestTime is a baseline in seconds, while the other numbers are ticks in milliseconds relatively to this requestTime.
	ProxyStart        int `json:"proxyStart"`        // Started resolving proxy.
	ProxyEnd          int `json:"proxyEnd"`          // Finished resolving proxy.
	DnsStart          int `json:"dnsStart"`          // Started DNS address resolve.
	DnsEnd            int `json:"dnsEnd"`            // Finished DNS address resolve.
	ConnectStart      int `json:"connectStart"`      // Started connecting to the remote host.
	ConnectEnd        int `json:"connectEnd"`        // Connected to the remote host.
	SslStart          int `json:"sslStart"`          // Started SSL handshake.
	SslEnd            int `json:"sslEnd"`            // Finished SSL handshake.
	WorkerStart       int `json:"workerStart"`       // Started running ServiceWorker.
	WorkerReady       int `json:"workerReady"`       // Finished Starting ServiceWorker.
	SendStart         int `json:"sendStart"`         // Started sending request.
	SendEnd           int `json:"sendEnd"`           // Finished sending request.
	PushStart         int `json:"pushStart"`         // Time the server started pushing request.
	PushEnd           int `json:"pushEnd"`           // Time the server finished pushing request.
	ReceiveHeadersEnd int `json:"receiveHeadersEnd"` // Finished receiving response headers.
}

// Loading priority of a resource request.
type ResourcePriority string

const ResourcePriorityVeryLow ResourcePriority = "VeryLow"
const ResourcePriorityLow ResourcePriority = "Low"
const ResourcePriorityMedium ResourcePriority = "Medium"
const ResourcePriorityHigh ResourcePriority = "High"
const ResourcePriorityVeryHigh ResourcePriority = "VeryHigh"

// HTTP request data.
type Request struct {
	Url              string           `json:"url"`              // Request URL.
	Method           string           `json:"method"`           // HTTP request method.
	Headers          *Headers         `json:"headers"`          // HTTP request headers.
	PostData         string           `json:"postData"`         // HTTP POST request data.
	MixedContentType string           `json:"mixedContentType"` // The mixed content status of the request, as defined in http://www.w3.org/TR/mixed-content/
	InitialPriority  ResourcePriority `json:"initialPriority"`  // Priority of the resource request at the time request is sent.
}

// Details of a signed certificate timestamp (SCT).
type SignedCertificateTimestamp struct {
	Status             string           `json:"status"`             // Validation status.
	Origin             string           `json:"origin"`             // Origin.
	LogDescription     string           `json:"logDescription"`     // Log name / description.
	LogId              string           `json:"logId"`              // Log ID.
	Timestamp          NetworkTimestamp `json:"timestamp"`          // Issuance date.
	HashAlgorithm      string           `json:"hashAlgorithm"`      // Hash algorithm.
	SignatureAlgorithm string           `json:"signatureAlgorithm"` // Signature algorithm.
	SignatureData      string           `json:"signatureData"`      // Signature data.
}

// Security details about a request.
type SecurityDetails struct {
	Protocol                       string                        `json:"protocol"`                       // Protocol name (e.g. "TLS 1.2" or "QUIC").
	KeyExchange                    string                        `json:"keyExchange"`                    // Key Exchange used by the connection.
	KeyExchangeGroup               string                        `json:"keyExchangeGroup"`               // (EC)DH group used by the connection, if applicable.
	Cipher                         string                        `json:"cipher"`                         // Cipher name.
	Mac                            string                        `json:"mac"`                            // TLS MAC. Note that AEAD ciphers do not have separate MACs.
	CertificateId                  *CertificateId                `json:"certificateId"`                  // Certificate ID value.
	SubjectName                    string                        `json:"subjectName"`                    // Certificate subject name.
	SanList                        []string                      `json:"sanList"`                        // Subject Alternative Name (SAN) DNS names and IP addresses.
	Issuer                         string                        `json:"issuer"`                         // Name of the issuing CA.
	ValidFrom                      NetworkTimestamp              `json:"validFrom"`                      // Certificate valid from date.
	ValidTo                        NetworkTimestamp              `json:"validTo"`                        // Certificate valid to (expiration) date
	SignedCertificateTimestampList []*SignedCertificateTimestamp `json:"signedCertificateTimestampList"` // List of signed certificate timestamps (SCTs).
}

// The reason why request was blocked.
type BlockedReason string

const BlockedReasonCsp BlockedReason = "csp"
const BlockedReasonMixedContent BlockedReason = "mixed-content"
const BlockedReasonOrigin BlockedReason = "origin"
const BlockedReasonInspector BlockedReason = "inspector"
const BlockedReasonSubresourceFilter BlockedReason = "subresource-filter"
const BlockedReasonOther BlockedReason = "other"

// HTTP response data.
type Response struct {
	Url                string           `json:"url"`                // Response URL. This URL can be different from CachedResource.url in case of redirect.
	Status             int              `json:"status"`             // HTTP response status code.
	StatusText         string           `json:"statusText"`         // HTTP response status text.
	Headers            *Headers         `json:"headers"`            // HTTP response headers.
	HeadersText        string           `json:"headersText"`        // HTTP response headers text.
	MimeType           string           `json:"mimeType"`           // Resource mimeType as determined by the browser.
	RequestHeaders     *Headers         `json:"requestHeaders"`     // Refined HTTP request headers that were actually transmitted over the network.
	RequestHeadersText string           `json:"requestHeadersText"` // HTTP request headers text.
	ConnectionReused   bool             `json:"connectionReused"`   // Specifies whether physical connection was actually reused for this request.
	ConnectionId       int              `json:"connectionId"`       // Physical connection id that was actually used for this request.
	RemoteIPAddress    string           `json:"remoteIPAddress"`    // Remote IP address.
	RemotePort         int              `json:"remotePort"`         // Remote port.
	FromDiskCache      bool             `json:"fromDiskCache"`      // Specifies that the request was served from the disk cache.
	FromServiceWorker  bool             `json:"fromServiceWorker"`  // Specifies that the request was served from the ServiceWorker.
	EncodedDataLength  int              `json:"encodedDataLength"`  // Total number of bytes received for this request so far.
	Timing             *ResourceTiming  `json:"timing"`             // Timing information for the given request.
	Protocol           string           `json:"protocol"`           // Protocol used to fetch this request.
	SecurityState      *SecurityState   `json:"securityState"`      // Security state of the request resource.
	SecurityDetails    *SecurityDetails `json:"securityDetails"`    // Security details for the request.
}

// WebSocket request data.
type WebSocketRequest struct {
	Headers *Headers `json:"headers"` // HTTP request headers.
}

// WebSocket response data.
type WebSocketResponse struct {
	Status             int      `json:"status"`             // HTTP response status code.
	StatusText         string   `json:"statusText"`         // HTTP response status text.
	Headers            *Headers `json:"headers"`            // HTTP response headers.
	HeadersText        string   `json:"headersText"`        // HTTP response headers text.
	RequestHeaders     *Headers `json:"requestHeaders"`     // HTTP request headers.
	RequestHeadersText string   `json:"requestHeadersText"` // HTTP request headers text.
}

// WebSocket frame data.
type WebSocketFrame struct {
	Opcode      int    `json:"opcode"`      // WebSocket frame opcode.
	Mask        bool   `json:"mask"`        // WebSocke frame mask.
	PayloadData string `json:"payloadData"` // WebSocke frame payload data.
}

// Information about the cached resource.
type CachedResource struct {
	Url      string        `json:"url"`      // Resource URL. This is the url of the original network request.
	Type     *ResourceType `json:"type"`     // Type of this resource.
	Response *Response     `json:"response"` // Cached response data.
	BodySize int           `json:"bodySize"` // Cached response body size.
}

// Information about the request initiator.
type Initiator struct {
	Type       string      `json:"type"`       // Type of this initiator.
	Stack      *StackTrace `json:"stack"`      // Initiator JavaScript stack trace, set for Script only.
	Url        string      `json:"url"`        // Initiator URL, set for Parser type only.
	LineNumber int         `json:"lineNumber"` // Initiator line number, set for Parser type only (0-based).
}

// Cookie object
type Cookie struct {
	Name     string         `json:"name"`     // Cookie name.
	Value    string         `json:"value"`    // Cookie value.
	Domain   string         `json:"domain"`   // Cookie domain.
	Path     string         `json:"path"`     // Cookie path.
	Expires  int            `json:"expires"`  // Cookie expiration date as the number of seconds since the UNIX epoch.
	Size     int            `json:"size"`     // Cookie size.
	HttpOnly bool           `json:"httpOnly"` // True if cookie is http-only.
	Secure   bool           `json:"secure"`   // True if cookie is secure.
	Session  bool           `json:"session"`  // True in case of session cookie.
	SameSite CookieSameSite `json:"sameSite"` // Cookie SameSite type.
}

type NetworkEnableParams struct {
	MaxTotalBufferSize    int `json:"maxTotalBufferSize"`    // Buffer size in bytes to use when preserving network payloads (XHRs, etc).
	MaxResourceBufferSize int `json:"maxResourceBufferSize"` // Per-resource buffer size in bytes to use when preserving network payloads (XHRs, etc).
}

type NetworkEnableCB func(err error)

// Enables network tracking, network events will now be delivered to the client.
type NetworkEnableCommand struct {
	params *NetworkEnableParams
	cb     NetworkEnableCB
}

func NewNetworkEnableCommand(params *NetworkEnableParams, cb NetworkEnableCB) *NetworkEnableCommand {
	return &NetworkEnableCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *NetworkEnableCommand) Name() string {
	return "Network.enable"
}

func (cmd *NetworkEnableCommand) Params() interface{} {
	return cmd.params
}

func (cmd *NetworkEnableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type NetworkDisableCB func(err error)

// Disables network tracking, prevents network events from being sent to the client.
type NetworkDisableCommand struct {
	cb NetworkDisableCB
}

func NewNetworkDisableCommand(cb NetworkDisableCB) *NetworkDisableCommand {
	return &NetworkDisableCommand{
		cb: cb,
	}
}

func (cmd *NetworkDisableCommand) Name() string {
	return "Network.disable"
}

func (cmd *NetworkDisableCommand) Params() interface{} {
	return nil
}

func (cmd *NetworkDisableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetUserAgentOverrideParams struct {
	UserAgent string `json:"userAgent"` // User agent to use.
}

type SetUserAgentOverrideCB func(err error)

// Allows overriding user agent with the given string.
type SetUserAgentOverrideCommand struct {
	params *SetUserAgentOverrideParams
	cb     SetUserAgentOverrideCB
}

func NewSetUserAgentOverrideCommand(params *SetUserAgentOverrideParams, cb SetUserAgentOverrideCB) *SetUserAgentOverrideCommand {
	return &SetUserAgentOverrideCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetUserAgentOverrideCommand) Name() string {
	return "Network.setUserAgentOverride"
}

func (cmd *SetUserAgentOverrideCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetUserAgentOverrideCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetExtraHTTPHeadersParams struct {
	Headers *Headers `json:"headers"` // Map with extra HTTP headers.
}

type SetExtraHTTPHeadersCB func(err error)

// Specifies whether to always send extra HTTP headers with the requests from this page.
type SetExtraHTTPHeadersCommand struct {
	params *SetExtraHTTPHeadersParams
	cb     SetExtraHTTPHeadersCB
}

func NewSetExtraHTTPHeadersCommand(params *SetExtraHTTPHeadersParams, cb SetExtraHTTPHeadersCB) *SetExtraHTTPHeadersCommand {
	return &SetExtraHTTPHeadersCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetExtraHTTPHeadersCommand) Name() string {
	return "Network.setExtraHTTPHeaders"
}

func (cmd *SetExtraHTTPHeadersCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetExtraHTTPHeadersCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type GetResponseBodyParams struct {
	RequestId RequestId `json:"requestId"` // Identifier of the network request to get content for.
}

type GetResponseBodyResult struct {
	Body          string `json:"body"`          // Response body.
	Base64Encoded bool   `json:"base64Encoded"` // True, if content was sent as base64.
}

type GetResponseBodyCB func(result *GetResponseBodyResult, err error)

// Returns content served for the given request.
type GetResponseBodyCommand struct {
	params *GetResponseBodyParams
	cb     GetResponseBodyCB
}

func NewGetResponseBodyCommand(params *GetResponseBodyParams, cb GetResponseBodyCB) *GetResponseBodyCommand {
	return &GetResponseBodyCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetResponseBodyCommand) Name() string {
	return "Network.getResponseBody"
}

func (cmd *GetResponseBodyCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetResponseBodyCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetResponseBodyResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type AddBlockedURLParams struct {
	Url string `json:"url"` // URL to block.
}

type AddBlockedURLCB func(err error)

// Blocks specific URL from loading.
type AddBlockedURLCommand struct {
	params *AddBlockedURLParams
	cb     AddBlockedURLCB
}

func NewAddBlockedURLCommand(params *AddBlockedURLParams, cb AddBlockedURLCB) *AddBlockedURLCommand {
	return &AddBlockedURLCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AddBlockedURLCommand) Name() string {
	return "Network.addBlockedURL"
}

func (cmd *AddBlockedURLCommand) Params() interface{} {
	return cmd.params
}

func (cmd *AddBlockedURLCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type RemoveBlockedURLParams struct {
	Url string `json:"url"` // URL to stop blocking.
}

type RemoveBlockedURLCB func(err error)

// Cancels blocking of a specific URL from loading.
type RemoveBlockedURLCommand struct {
	params *RemoveBlockedURLParams
	cb     RemoveBlockedURLCB
}

func NewRemoveBlockedURLCommand(params *RemoveBlockedURLParams, cb RemoveBlockedURLCB) *RemoveBlockedURLCommand {
	return &RemoveBlockedURLCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *RemoveBlockedURLCommand) Name() string {
	return "Network.removeBlockedURL"
}

func (cmd *RemoveBlockedURLCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveBlockedURLCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type ReplayXHRParams struct {
	RequestId RequestId `json:"requestId"` // Identifier of XHR to replay.
}

type ReplayXHRCB func(err error)

// This method sends a new XMLHttpRequest which is identical to the original one. The following parameters should be identical: method, url, async, request body, extra headers, withCredentials attribute, user, password.
type ReplayXHRCommand struct {
	params *ReplayXHRParams
	cb     ReplayXHRCB
}

func NewReplayXHRCommand(params *ReplayXHRParams, cb ReplayXHRCB) *ReplayXHRCommand {
	return &ReplayXHRCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *ReplayXHRCommand) Name() string {
	return "Network.replayXHR"
}

func (cmd *ReplayXHRCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ReplayXHRCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetMonitoringXHREnabledParams struct {
	Enabled bool `json:"enabled"` // Monitoring enabled state.
}

type SetMonitoringXHREnabledCB func(err error)

// Toggles monitoring of XMLHttpRequest. If true, console will receive messages upon each XHR issued.
type SetMonitoringXHREnabledCommand struct {
	params *SetMonitoringXHREnabledParams
	cb     SetMonitoringXHREnabledCB
}

func NewSetMonitoringXHREnabledCommand(params *SetMonitoringXHREnabledParams, cb SetMonitoringXHREnabledCB) *SetMonitoringXHREnabledCommand {
	return &SetMonitoringXHREnabledCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetMonitoringXHREnabledCommand) Name() string {
	return "Network.setMonitoringXHREnabled"
}

func (cmd *SetMonitoringXHREnabledCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetMonitoringXHREnabledCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type CanClearBrowserCacheResult struct {
	Result bool `json:"result"` // True if browser cache can be cleared.
}

type CanClearBrowserCacheCB func(result *CanClearBrowserCacheResult, err error)

// Tells whether clearing browser cache is supported.
type CanClearBrowserCacheCommand struct {
	cb CanClearBrowserCacheCB
}

func NewCanClearBrowserCacheCommand(cb CanClearBrowserCacheCB) *CanClearBrowserCacheCommand {
	return &CanClearBrowserCacheCommand{
		cb: cb,
	}
}

func (cmd *CanClearBrowserCacheCommand) Name() string {
	return "Network.canClearBrowserCache"
}

func (cmd *CanClearBrowserCacheCommand) Params() interface{} {
	return nil
}

func (cmd *CanClearBrowserCacheCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj CanClearBrowserCacheResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type ClearBrowserCacheCB func(err error)

// Clears browser cache.
type ClearBrowserCacheCommand struct {
	cb ClearBrowserCacheCB
}

func NewClearBrowserCacheCommand(cb ClearBrowserCacheCB) *ClearBrowserCacheCommand {
	return &ClearBrowserCacheCommand{
		cb: cb,
	}
}

func (cmd *ClearBrowserCacheCommand) Name() string {
	return "Network.clearBrowserCache"
}

func (cmd *ClearBrowserCacheCommand) Params() interface{} {
	return nil
}

func (cmd *ClearBrowserCacheCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type CanClearBrowserCookiesResult struct {
	Result bool `json:"result"` // True if browser cookies can be cleared.
}

type CanClearBrowserCookiesCB func(result *CanClearBrowserCookiesResult, err error)

// Tells whether clearing browser cookies is supported.
type CanClearBrowserCookiesCommand struct {
	cb CanClearBrowserCookiesCB
}

func NewCanClearBrowserCookiesCommand(cb CanClearBrowserCookiesCB) *CanClearBrowserCookiesCommand {
	return &CanClearBrowserCookiesCommand{
		cb: cb,
	}
}

func (cmd *CanClearBrowserCookiesCommand) Name() string {
	return "Network.canClearBrowserCookies"
}

func (cmd *CanClearBrowserCookiesCommand) Params() interface{} {
	return nil
}

func (cmd *CanClearBrowserCookiesCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj CanClearBrowserCookiesResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type ClearBrowserCookiesCB func(err error)

// Clears browser cookies.
type ClearBrowserCookiesCommand struct {
	cb ClearBrowserCookiesCB
}

func NewClearBrowserCookiesCommand(cb ClearBrowserCookiesCB) *ClearBrowserCookiesCommand {
	return &ClearBrowserCookiesCommand{
		cb: cb,
	}
}

func (cmd *ClearBrowserCookiesCommand) Name() string {
	return "Network.clearBrowserCookies"
}

func (cmd *ClearBrowserCookiesCommand) Params() interface{} {
	return nil
}

func (cmd *ClearBrowserCookiesCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type NetworkGetCookiesResult struct {
	Cookies []*Cookie `json:"cookies"` // Array of cookie objects.
}

type NetworkGetCookiesCB func(result *NetworkGetCookiesResult, err error)

// Returns all browser cookies. Depending on the backend support, will return detailed cookie information in the cookies field.
type NetworkGetCookiesCommand struct {
	cb NetworkGetCookiesCB
}

func NewNetworkGetCookiesCommand(cb NetworkGetCookiesCB) *NetworkGetCookiesCommand {
	return &NetworkGetCookiesCommand{
		cb: cb,
	}
}

func (cmd *NetworkGetCookiesCommand) Name() string {
	return "Network.getCookies"
}

func (cmd *NetworkGetCookiesCommand) Params() interface{} {
	return nil
}

func (cmd *NetworkGetCookiesCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj NetworkGetCookiesResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type NetworkDeleteCookieParams struct {
	CookieName string `json:"cookieName"` // Name of the cookie to remove.
	Url        string `json:"url"`        // URL to match cooke domain and path.
}

type NetworkDeleteCookieCB func(err error)

// Deletes browser cookie with given name, domain and path.
type NetworkDeleteCookieCommand struct {
	params *NetworkDeleteCookieParams
	cb     NetworkDeleteCookieCB
}

func NewNetworkDeleteCookieCommand(params *NetworkDeleteCookieParams, cb NetworkDeleteCookieCB) *NetworkDeleteCookieCommand {
	return &NetworkDeleteCookieCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *NetworkDeleteCookieCommand) Name() string {
	return "Network.deleteCookie"
}

func (cmd *NetworkDeleteCookieCommand) Params() interface{} {
	return cmd.params
}

func (cmd *NetworkDeleteCookieCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetCookieParams struct {
	Url            string           `json:"url"`            // The request-URI to associate with the setting of the cookie. This value can affect the default domain and path values of the created cookie.
	Name           string           `json:"name"`           // The name of the cookie.
	Value          string           `json:"value"`          // The value of the cookie.
	Domain         string           `json:"domain"`         // If omitted, the cookie becomes a host-only cookie.
	Path           string           `json:"path"`           // Defaults to the path portion of the url parameter.
	Secure         bool             `json:"secure"`         // Defaults ot false.
	HttpOnly       bool             `json:"httpOnly"`       // Defaults to false.
	SameSite       CookieSameSite   `json:"sameSite"`       // Defaults to browser default behavior.
	ExpirationDate NetworkTimestamp `json:"expirationDate"` // If omitted, the cookie becomes a session cookie.
}

type SetCookieResult struct {
	Success bool `json:"success"` // True if successfully set cookie.
}

type SetCookieCB func(result *SetCookieResult, err error)

// Sets a cookie with the given cookie data; may overwrite equivalent cookies if they exist.
type SetCookieCommand struct {
	params *SetCookieParams
	cb     SetCookieCB
}

func NewSetCookieCommand(params *SetCookieParams, cb SetCookieCB) *SetCookieCommand {
	return &SetCookieCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetCookieCommand) Name() string {
	return "Network.setCookie"
}

func (cmd *SetCookieCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetCookieCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj SetCookieResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type CanEmulateNetworkConditionsResult struct {
	Result bool `json:"result"` // True if emulation of network conditions is supported.
}

type CanEmulateNetworkConditionsCB func(result *CanEmulateNetworkConditionsResult, err error)

// Tells whether emulation of network conditions is supported.
type CanEmulateNetworkConditionsCommand struct {
	cb CanEmulateNetworkConditionsCB
}

func NewCanEmulateNetworkConditionsCommand(cb CanEmulateNetworkConditionsCB) *CanEmulateNetworkConditionsCommand {
	return &CanEmulateNetworkConditionsCommand{
		cb: cb,
	}
}

func (cmd *CanEmulateNetworkConditionsCommand) Name() string {
	return "Network.canEmulateNetworkConditions"
}

func (cmd *CanEmulateNetworkConditionsCommand) Params() interface{} {
	return nil
}

func (cmd *CanEmulateNetworkConditionsCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj CanEmulateNetworkConditionsResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type EmulateNetworkConditionsParams struct {
	Offline            bool           `json:"offline"`            // True to emulate internet disconnection.
	Latency            int            `json:"latency"`            // Additional latency (ms).
	DownloadThroughput int            `json:"downloadThroughput"` // Maximal aggregated download throughput.
	UploadThroughput   int            `json:"uploadThroughput"`   // Maximal aggregated upload throughput.
	ConnectionType     ConnectionType `json:"connectionType"`     // Connection type if known.
}

type EmulateNetworkConditionsCB func(err error)

// Activates emulation of network conditions.
type EmulateNetworkConditionsCommand struct {
	params *EmulateNetworkConditionsParams
	cb     EmulateNetworkConditionsCB
}

func NewEmulateNetworkConditionsCommand(params *EmulateNetworkConditionsParams, cb EmulateNetworkConditionsCB) *EmulateNetworkConditionsCommand {
	return &EmulateNetworkConditionsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *EmulateNetworkConditionsCommand) Name() string {
	return "Network.emulateNetworkConditions"
}

func (cmd *EmulateNetworkConditionsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *EmulateNetworkConditionsCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetCacheDisabledParams struct {
	CacheDisabled bool `json:"cacheDisabled"` // Cache disabled state.
}

type SetCacheDisabledCB func(err error)

// Toggles ignoring cache for each request. If true, cache will not be used.
type SetCacheDisabledCommand struct {
	params *SetCacheDisabledParams
	cb     SetCacheDisabledCB
}

func NewSetCacheDisabledCommand(params *SetCacheDisabledParams, cb SetCacheDisabledCB) *SetCacheDisabledCommand {
	return &SetCacheDisabledCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetCacheDisabledCommand) Name() string {
	return "Network.setCacheDisabled"
}

func (cmd *SetCacheDisabledCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetCacheDisabledCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetBypassServiceWorkerParams struct {
	Bypass bool `json:"bypass"` // Bypass service worker and load from network.
}

type SetBypassServiceWorkerCB func(err error)

// Toggles ignoring of service worker for each request.
type SetBypassServiceWorkerCommand struct {
	params *SetBypassServiceWorkerParams
	cb     SetBypassServiceWorkerCB
}

func NewSetBypassServiceWorkerCommand(params *SetBypassServiceWorkerParams, cb SetBypassServiceWorkerCB) *SetBypassServiceWorkerCommand {
	return &SetBypassServiceWorkerCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetBypassServiceWorkerCommand) Name() string {
	return "Network.setBypassServiceWorker"
}

func (cmd *SetBypassServiceWorkerCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetBypassServiceWorkerCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetDataSizeLimitsForTestParams struct {
	MaxTotalSize    int `json:"maxTotalSize"`    // Maximum total buffer size.
	MaxResourceSize int `json:"maxResourceSize"` // Maximum per-resource size.
}

type SetDataSizeLimitsForTestCB func(err error)

// For testing.
type SetDataSizeLimitsForTestCommand struct {
	params *SetDataSizeLimitsForTestParams
	cb     SetDataSizeLimitsForTestCB
}

func NewSetDataSizeLimitsForTestCommand(params *SetDataSizeLimitsForTestParams, cb SetDataSizeLimitsForTestCB) *SetDataSizeLimitsForTestCommand {
	return &SetDataSizeLimitsForTestCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetDataSizeLimitsForTestCommand) Name() string {
	return "Network.setDataSizeLimitsForTest"
}

func (cmd *SetDataSizeLimitsForTestCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetDataSizeLimitsForTestCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type GetCertificateParams struct {
	Origin string `json:"origin"` // Origin to get certificate for.
}

type GetCertificateResult struct {
	TableNames []string `json:"tableNames"`
}

type GetCertificateCB func(result *GetCertificateResult, err error)

// Returns the DER-encoded certificate.
type GetCertificateCommand struct {
	params *GetCertificateParams
	cb     GetCertificateCB
}

func NewGetCertificateCommand(params *GetCertificateParams, cb GetCertificateCB) *GetCertificateCommand {
	return &GetCertificateCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetCertificateCommand) Name() string {
	return "Network.getCertificate"
}

func (cmd *GetCertificateCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetCertificateCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetCertificateResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type ResourceChangedPriorityEvent struct {
	RequestId   RequestId        `json:"requestId"`   // Request identifier.
	NewPriority ResourcePriority `json:"newPriority"` // New priority
	Timestamp   NetworkTimestamp `json:"timestamp"`   // Timestamp.
}

// Fired when resource loading priority is changed
type ResourceChangedPriorityEventSink struct {
	events chan *ResourceChangedPriorityEvent
}

func NewResourceChangedPriorityEventSink(bufSize int) *ResourceChangedPriorityEventSink {
	return &ResourceChangedPriorityEventSink{
		events: make(chan *ResourceChangedPriorityEvent, bufSize),
	}
}

func (s *ResourceChangedPriorityEventSink) Name() string {
	return "Network.resourceChangedPriority"
}

func (s *ResourceChangedPriorityEventSink) OnEvent(params []byte) {
	evt := &ResourceChangedPriorityEvent{}
	if err := json.Unmarshal(params, evt); err != nil {
		logging.Vlog(-1, err)
	} else {
		select {
		case s.events <- evt:
			// Do nothing.
		default:
			logging.Vlogf(0, "Dropped one event(%v).", evt)
		}
	}
}

type RequestWillBeSentEvent struct {
	RequestId        RequestId        `json:"requestId"`        // Request identifier.
	FrameId          *FrameId         `json:"frameId"`          // Frame identifier.
	LoaderId         LoaderId         `json:"loaderId"`         // Loader identifier.
	DocumentURL      string           `json:"documentURL"`      // URL of the document this request is loaded for.
	Request          *Request         `json:"request"`          // Request data.
	Timestamp        NetworkTimestamp `json:"timestamp"`        // Timestamp.
	WallTime         NetworkTimestamp `json:"wallTime"`         // UTC Timestamp.
	Initiator        *Initiator       `json:"initiator"`        // Request initiator.
	RedirectResponse *Response        `json:"redirectResponse"` // Redirect response data.
	Type             *ResourceType    `json:"type"`             // Type of this resource.
}

// Fired when page is about to send HTTP request.
type RequestWillBeSentEventSink struct {
	events chan *RequestWillBeSentEvent
}

func NewRequestWillBeSentEventSink(bufSize int) *RequestWillBeSentEventSink {
	return &RequestWillBeSentEventSink{
		events: make(chan *RequestWillBeSentEvent, bufSize),
	}
}

func (s *RequestWillBeSentEventSink) Name() string {
	return "Network.requestWillBeSent"
}

func (s *RequestWillBeSentEventSink) OnEvent(params []byte) {
	evt := &RequestWillBeSentEvent{}
	if err := json.Unmarshal(params, evt); err != nil {
		logging.Vlog(-1, err)
	} else {
		select {
		case s.events <- evt:
			// Do nothing.
		default:
			logging.Vlogf(0, "Dropped one event(%v).", evt)
		}
	}
}

type RequestServedFromCacheEvent struct {
	RequestId RequestId `json:"requestId"` // Request identifier.
}

// Fired if request ended up loading from cache.
type RequestServedFromCacheEventSink struct {
	events chan *RequestServedFromCacheEvent
}

func NewRequestServedFromCacheEventSink(bufSize int) *RequestServedFromCacheEventSink {
	return &RequestServedFromCacheEventSink{
		events: make(chan *RequestServedFromCacheEvent, bufSize),
	}
}

func (s *RequestServedFromCacheEventSink) Name() string {
	return "Network.requestServedFromCache"
}

func (s *RequestServedFromCacheEventSink) OnEvent(params []byte) {
	evt := &RequestServedFromCacheEvent{}
	if err := json.Unmarshal(params, evt); err != nil {
		logging.Vlog(-1, err)
	} else {
		select {
		case s.events <- evt:
			// Do nothing.
		default:
			logging.Vlogf(0, "Dropped one event(%v).", evt)
		}
	}
}

type ResponseReceivedEvent struct {
	RequestId RequestId        `json:"requestId"` // Request identifier.
	FrameId   *FrameId         `json:"frameId"`   // Frame identifier.
	LoaderId  LoaderId         `json:"loaderId"`  // Loader identifier.
	Timestamp NetworkTimestamp `json:"timestamp"` // Timestamp.
	Type      *ResourceType    `json:"type"`      // Resource type.
	Response  *Response        `json:"response"`  // Response data.
}

// Fired when HTTP response is available.
type ResponseReceivedEventSink struct {
	events chan *ResponseReceivedEvent
}

func NewResponseReceivedEventSink(bufSize int) *ResponseReceivedEventSink {
	return &ResponseReceivedEventSink{
		events: make(chan *ResponseReceivedEvent, bufSize),
	}
}

func (s *ResponseReceivedEventSink) Name() string {
	return "Network.responseReceived"
}

func (s *ResponseReceivedEventSink) OnEvent(params []byte) {
	evt := &ResponseReceivedEvent{}
	if err := json.Unmarshal(params, evt); err != nil {
		logging.Vlog(-1, err)
	} else {
		select {
		case s.events <- evt:
			// Do nothing.
		default:
			logging.Vlogf(0, "Dropped one event(%v).", evt)
		}
	}
}

type DataReceivedEvent struct {
	RequestId         RequestId        `json:"requestId"`         // Request identifier.
	Timestamp         NetworkTimestamp `json:"timestamp"`         // Timestamp.
	DataLength        int              `json:"dataLength"`        // Data chunk length.
	EncodedDataLength int              `json:"encodedDataLength"` // Actual bytes received (might be less than dataLength for compressed encodings).
}

// Fired when data chunk was received over the network.
type DataReceivedEventSink struct {
	events chan *DataReceivedEvent
}

func NewDataReceivedEventSink(bufSize int) *DataReceivedEventSink {
	return &DataReceivedEventSink{
		events: make(chan *DataReceivedEvent, bufSize),
	}
}

func (s *DataReceivedEventSink) Name() string {
	return "Network.dataReceived"
}

func (s *DataReceivedEventSink) OnEvent(params []byte) {
	evt := &DataReceivedEvent{}
	if err := json.Unmarshal(params, evt); err != nil {
		logging.Vlog(-1, err)
	} else {
		select {
		case s.events <- evt:
			// Do nothing.
		default:
			logging.Vlogf(0, "Dropped one event(%v).", evt)
		}
	}
}

type LoadingFinishedEvent struct {
	RequestId         RequestId        `json:"requestId"`         // Request identifier.
	Timestamp         NetworkTimestamp `json:"timestamp"`         // Timestamp.
	EncodedDataLength int              `json:"encodedDataLength"` // Total number of bytes received for this request.
}

// Fired when HTTP request has finished loading.
type LoadingFinishedEventSink struct {
	events chan *LoadingFinishedEvent
}

func NewLoadingFinishedEventSink(bufSize int) *LoadingFinishedEventSink {
	return &LoadingFinishedEventSink{
		events: make(chan *LoadingFinishedEvent, bufSize),
	}
}

func (s *LoadingFinishedEventSink) Name() string {
	return "Network.loadingFinished"
}

func (s *LoadingFinishedEventSink) OnEvent(params []byte) {
	evt := &LoadingFinishedEvent{}
	if err := json.Unmarshal(params, evt); err != nil {
		logging.Vlog(-1, err)
	} else {
		select {
		case s.events <- evt:
			// Do nothing.
		default:
			logging.Vlogf(0, "Dropped one event(%v).", evt)
		}
	}
}

type LoadingFailedEvent struct {
	RequestId     RequestId        `json:"requestId"`     // Request identifier.
	Timestamp     NetworkTimestamp `json:"timestamp"`     // Timestamp.
	Type          *ResourceType    `json:"type"`          // Resource type.
	ErrorText     string           `json:"errorText"`     // User friendly error message.
	Canceled      bool             `json:"canceled"`      // True if loading was canceled.
	BlockedReason BlockedReason    `json:"blockedReason"` // The reason why loading was blocked, if any.
}

// Fired when HTTP request has failed to load.
type LoadingFailedEventSink struct {
	events chan *LoadingFailedEvent
}

func NewLoadingFailedEventSink(bufSize int) *LoadingFailedEventSink {
	return &LoadingFailedEventSink{
		events: make(chan *LoadingFailedEvent, bufSize),
	}
}

func (s *LoadingFailedEventSink) Name() string {
	return "Network.loadingFailed"
}

func (s *LoadingFailedEventSink) OnEvent(params []byte) {
	evt := &LoadingFailedEvent{}
	if err := json.Unmarshal(params, evt); err != nil {
		logging.Vlog(-1, err)
	} else {
		select {
		case s.events <- evt:
			// Do nothing.
		default:
			logging.Vlogf(0, "Dropped one event(%v).", evt)
		}
	}
}

type WebSocketWillSendHandshakeRequestEvent struct {
	RequestId RequestId         `json:"requestId"` // Request identifier.
	Timestamp NetworkTimestamp  `json:"timestamp"` // Timestamp.
	WallTime  NetworkTimestamp  `json:"wallTime"`  // UTC Timestamp.
	Request   *WebSocketRequest `json:"request"`   // WebSocket request data.
}

// Fired when WebSocket is about to initiate handshake.
type WebSocketWillSendHandshakeRequestEventSink struct {
	events chan *WebSocketWillSendHandshakeRequestEvent
}

func NewWebSocketWillSendHandshakeRequestEventSink(bufSize int) *WebSocketWillSendHandshakeRequestEventSink {
	return &WebSocketWillSendHandshakeRequestEventSink{
		events: make(chan *WebSocketWillSendHandshakeRequestEvent, bufSize),
	}
}

func (s *WebSocketWillSendHandshakeRequestEventSink) Name() string {
	return "Network.webSocketWillSendHandshakeRequest"
}

func (s *WebSocketWillSendHandshakeRequestEventSink) OnEvent(params []byte) {
	evt := &WebSocketWillSendHandshakeRequestEvent{}
	if err := json.Unmarshal(params, evt); err != nil {
		logging.Vlog(-1, err)
	} else {
		select {
		case s.events <- evt:
			// Do nothing.
		default:
			logging.Vlogf(0, "Dropped one event(%v).", evt)
		}
	}
}

type WebSocketHandshakeResponseReceivedEvent struct {
	RequestId RequestId          `json:"requestId"` // Request identifier.
	Timestamp NetworkTimestamp   `json:"timestamp"` // Timestamp.
	Response  *WebSocketResponse `json:"response"`  // WebSocket response data.
}

// Fired when WebSocket handshake response becomes available.
type WebSocketHandshakeResponseReceivedEventSink struct {
	events chan *WebSocketHandshakeResponseReceivedEvent
}

func NewWebSocketHandshakeResponseReceivedEventSink(bufSize int) *WebSocketHandshakeResponseReceivedEventSink {
	return &WebSocketHandshakeResponseReceivedEventSink{
		events: make(chan *WebSocketHandshakeResponseReceivedEvent, bufSize),
	}
}

func (s *WebSocketHandshakeResponseReceivedEventSink) Name() string {
	return "Network.webSocketHandshakeResponseReceived"
}

func (s *WebSocketHandshakeResponseReceivedEventSink) OnEvent(params []byte) {
	evt := &WebSocketHandshakeResponseReceivedEvent{}
	if err := json.Unmarshal(params, evt); err != nil {
		logging.Vlog(-1, err)
	} else {
		select {
		case s.events <- evt:
			// Do nothing.
		default:
			logging.Vlogf(0, "Dropped one event(%v).", evt)
		}
	}
}

type WebSocketCreatedEvent struct {
	RequestId RequestId  `json:"requestId"` // Request identifier.
	Url       string     `json:"url"`       // WebSocket request URL.
	Initiator *Initiator `json:"initiator"` // Request initiator.
}

// Fired upon WebSocket creation.
type WebSocketCreatedEventSink struct {
	events chan *WebSocketCreatedEvent
}

func NewWebSocketCreatedEventSink(bufSize int) *WebSocketCreatedEventSink {
	return &WebSocketCreatedEventSink{
		events: make(chan *WebSocketCreatedEvent, bufSize),
	}
}

func (s *WebSocketCreatedEventSink) Name() string {
	return "Network.webSocketCreated"
}

func (s *WebSocketCreatedEventSink) OnEvent(params []byte) {
	evt := &WebSocketCreatedEvent{}
	if err := json.Unmarshal(params, evt); err != nil {
		logging.Vlog(-1, err)
	} else {
		select {
		case s.events <- evt:
			// Do nothing.
		default:
			logging.Vlogf(0, "Dropped one event(%v).", evt)
		}
	}
}

type WebSocketClosedEvent struct {
	RequestId RequestId        `json:"requestId"` // Request identifier.
	Timestamp NetworkTimestamp `json:"timestamp"` // Timestamp.
}

// Fired when WebSocket is closed.
type WebSocketClosedEventSink struct {
	events chan *WebSocketClosedEvent
}

func NewWebSocketClosedEventSink(bufSize int) *WebSocketClosedEventSink {
	return &WebSocketClosedEventSink{
		events: make(chan *WebSocketClosedEvent, bufSize),
	}
}

func (s *WebSocketClosedEventSink) Name() string {
	return "Network.webSocketClosed"
}

func (s *WebSocketClosedEventSink) OnEvent(params []byte) {
	evt := &WebSocketClosedEvent{}
	if err := json.Unmarshal(params, evt); err != nil {
		logging.Vlog(-1, err)
	} else {
		select {
		case s.events <- evt:
			// Do nothing.
		default:
			logging.Vlogf(0, "Dropped one event(%v).", evt)
		}
	}
}

type WebSocketFrameReceivedEvent struct {
	RequestId RequestId        `json:"requestId"` // Request identifier.
	Timestamp NetworkTimestamp `json:"timestamp"` // Timestamp.
	Response  *WebSocketFrame  `json:"response"`  // WebSocket response data.
}

// Fired when WebSocket frame is received.
type WebSocketFrameReceivedEventSink struct {
	events chan *WebSocketFrameReceivedEvent
}

func NewWebSocketFrameReceivedEventSink(bufSize int) *WebSocketFrameReceivedEventSink {
	return &WebSocketFrameReceivedEventSink{
		events: make(chan *WebSocketFrameReceivedEvent, bufSize),
	}
}

func (s *WebSocketFrameReceivedEventSink) Name() string {
	return "Network.webSocketFrameReceived"
}

func (s *WebSocketFrameReceivedEventSink) OnEvent(params []byte) {
	evt := &WebSocketFrameReceivedEvent{}
	if err := json.Unmarshal(params, evt); err != nil {
		logging.Vlog(-1, err)
	} else {
		select {
		case s.events <- evt:
			// Do nothing.
		default:
			logging.Vlogf(0, "Dropped one event(%v).", evt)
		}
	}
}

type WebSocketFrameErrorEvent struct {
	RequestId    RequestId        `json:"requestId"`    // Request identifier.
	Timestamp    NetworkTimestamp `json:"timestamp"`    // Timestamp.
	ErrorMessage string           `json:"errorMessage"` // WebSocket frame error message.
}

// Fired when WebSocket frame error occurs.
type WebSocketFrameErrorEventSink struct {
	events chan *WebSocketFrameErrorEvent
}

func NewWebSocketFrameErrorEventSink(bufSize int) *WebSocketFrameErrorEventSink {
	return &WebSocketFrameErrorEventSink{
		events: make(chan *WebSocketFrameErrorEvent, bufSize),
	}
}

func (s *WebSocketFrameErrorEventSink) Name() string {
	return "Network.webSocketFrameError"
}

func (s *WebSocketFrameErrorEventSink) OnEvent(params []byte) {
	evt := &WebSocketFrameErrorEvent{}
	if err := json.Unmarshal(params, evt); err != nil {
		logging.Vlog(-1, err)
	} else {
		select {
		case s.events <- evt:
			// Do nothing.
		default:
			logging.Vlogf(0, "Dropped one event(%v).", evt)
		}
	}
}

type WebSocketFrameSentEvent struct {
	RequestId RequestId        `json:"requestId"` // Request identifier.
	Timestamp NetworkTimestamp `json:"timestamp"` // Timestamp.
	Response  *WebSocketFrame  `json:"response"`  // WebSocket response data.
}

// Fired when WebSocket frame is sent.
type WebSocketFrameSentEventSink struct {
	events chan *WebSocketFrameSentEvent
}

func NewWebSocketFrameSentEventSink(bufSize int) *WebSocketFrameSentEventSink {
	return &WebSocketFrameSentEventSink{
		events: make(chan *WebSocketFrameSentEvent, bufSize),
	}
}

func (s *WebSocketFrameSentEventSink) Name() string {
	return "Network.webSocketFrameSent"
}

func (s *WebSocketFrameSentEventSink) OnEvent(params []byte) {
	evt := &WebSocketFrameSentEvent{}
	if err := json.Unmarshal(params, evt); err != nil {
		logging.Vlog(-1, err)
	} else {
		select {
		case s.events <- evt:
			// Do nothing.
		default:
			logging.Vlogf(0, "Dropped one event(%v).", evt)
		}
	}
}

type EventSourceMessageReceivedEvent struct {
	RequestId RequestId        `json:"requestId"` // Request identifier.
	Timestamp NetworkTimestamp `json:"timestamp"` // Timestamp.
	EventName string           `json:"eventName"` // Message type.
	EventId   string           `json:"eventId"`   // Message identifier.
	Data      string           `json:"data"`      // Message content.
}

// Fired when EventSource message is received.
type EventSourceMessageReceivedEventSink struct {
	events chan *EventSourceMessageReceivedEvent
}

func NewEventSourceMessageReceivedEventSink(bufSize int) *EventSourceMessageReceivedEventSink {
	return &EventSourceMessageReceivedEventSink{
		events: make(chan *EventSourceMessageReceivedEvent, bufSize),
	}
}

func (s *EventSourceMessageReceivedEventSink) Name() string {
	return "Network.eventSourceMessageReceived"
}

func (s *EventSourceMessageReceivedEventSink) OnEvent(params []byte) {
	evt := &EventSourceMessageReceivedEvent{}
	if err := json.Unmarshal(params, evt); err != nil {
		logging.Vlog(-1, err)
	} else {
		select {
		case s.events <- evt:
			// Do nothing.
		default:
			logging.Vlogf(0, "Dropped one event(%v).", evt)
		}
	}
}
