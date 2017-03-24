package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

// Unique loader identifier.
type LoaderId string

// Unique request identifier.
type RequestId string

// Number of seconds since epoch.
type NetworkTimestamp float64

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
	RequestTime       float64 `json:"requestTime"`       // Timing's requestTime is a baseline in seconds, while the other numbers are ticks in milliseconds relatively to this requestTime.
	ProxyStart        float64 `json:"proxyStart"`        // Started resolving proxy.
	ProxyEnd          float64 `json:"proxyEnd"`          // Finished resolving proxy.
	DnsStart          float64 `json:"dnsStart"`          // Started DNS address resolve.
	DnsEnd            float64 `json:"dnsEnd"`            // Finished DNS address resolve.
	ConnectStart      float64 `json:"connectStart"`      // Started connecting to the remote host.
	ConnectEnd        float64 `json:"connectEnd"`        // Connected to the remote host.
	SslStart          float64 `json:"sslStart"`          // Started SSL handshake.
	SslEnd            float64 `json:"sslEnd"`            // Finished SSL handshake.
	WorkerStart       float64 `json:"workerStart"`       // Started running ServiceWorker.
	WorkerReady       float64 `json:"workerReady"`       // Finished Starting ServiceWorker.
	SendStart         float64 `json:"sendStart"`         // Started sending request.
	SendEnd           float64 `json:"sendEnd"`           // Finished sending request.
	PushStart         float64 `json:"pushStart"`         // Time the server started pushing request.
	PushEnd           float64 `json:"pushEnd"`           // Time the server finished pushing request.
	ReceiveHeadersEnd float64 `json:"receiveHeadersEnd"` // Finished receiving response headers.
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
	Url              string           `json:"url"`                        // Request URL.
	Method           string           `json:"method"`                     // HTTP request method.
	Headers          *Headers         `json:"headers"`                    // HTTP request headers.
	PostData         string           `json:"postData,omitempty"`         // HTTP POST request data.
	MixedContentType string           `json:"mixedContentType,omitempty"` // The mixed content status of the request, as defined in http://www.w3.org/TR/mixed-content/
	InitialPriority  ResourcePriority `json:"initialPriority"`            // Priority of the resource request at the time request is sent.
	ReferrerPolicy   string           `json:"referrerPolicy"`             // The referrer policy of the request, as defined in https://www.w3.org/TR/referrer-policy/
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
	KeyExchange                    string                        `json:"keyExchange"`                    // Key Exchange used by the connection, or the empty string if not applicable.
	KeyExchangeGroup               string                        `json:"keyExchangeGroup,omitempty"`     // (EC)DH group used by the connection, if applicable.
	Cipher                         string                        `json:"cipher"`                         // Cipher name.
	Mac                            string                        `json:"mac,omitempty"`                  // TLS MAC. Note that AEAD ciphers do not have separate MACs.
	CertificateId                  *CertificateId                `json:"certificateId"`                  // Certificate ID value.
	SubjectName                    string                        `json:"subjectName"`                    // Certificate subject name.
	SanList                        []string                      `json:"sanList"`                        // Subject Alternative Name (SAN) DNS names and IP addresses.
	Issuer                         string                        `json:"issuer"`                         // Name of the issuing CA.
	ValidFrom                      NetworkTimestamp              `json:"validFrom"`                      // Certificate valid from date.
	ValidTo                        NetworkTimestamp              `json:"validTo"`                        // Certificate valid to (expiration) date
	SignedCertificateTimestampList []*SignedCertificateTimestamp `json:"signedCertificateTimestampList"` // List of signed certificate timestamps (SCTs).
}

// The reason why request was blocked.
// @experimental
type BlockedReason string

const BlockedReasonCsp BlockedReason = "csp"
const BlockedReasonMixedContent BlockedReason = "mixed-content"
const BlockedReasonOrigin BlockedReason = "origin"
const BlockedReasonInspector BlockedReason = "inspector"
const BlockedReasonSubresourceFilter BlockedReason = "subresource-filter"
const BlockedReasonOther BlockedReason = "other"

// HTTP response data.
type Response struct {
	Url                string           `json:"url"`                          // Response URL. This URL can be different from CachedResource.url in case of redirect.
	Status             float64          `json:"status"`                       // HTTP response status code.
	StatusText         string           `json:"statusText"`                   // HTTP response status text.
	Headers            *Headers         `json:"headers"`                      // HTTP response headers.
	HeadersText        string           `json:"headersText,omitempty"`        // HTTP response headers text.
	MimeType           string           `json:"mimeType"`                     // Resource mimeType as determined by the browser.
	RequestHeaders     *Headers         `json:"requestHeaders,omitempty"`     // Refined HTTP request headers that were actually transmitted over the network.
	RequestHeadersText string           `json:"requestHeadersText,omitempty"` // HTTP request headers text.
	ConnectionReused   bool             `json:"connectionReused"`             // Specifies whether physical connection was actually reused for this request.
	ConnectionId       float64          `json:"connectionId"`                 // Physical connection id that was actually used for this request.
	RemoteIPAddress    string           `json:"remoteIPAddress,omitempty"`    // Remote IP address.
	RemotePort         int              `json:"remotePort,omitempty"`         // Remote port.
	FromDiskCache      bool             `json:"fromDiskCache,omitempty"`      // Specifies that the request was served from the disk cache.
	FromServiceWorker  bool             `json:"fromServiceWorker,omitempty"`  // Specifies that the request was served from the ServiceWorker.
	EncodedDataLength  float64          `json:"encodedDataLength"`            // Total number of bytes received for this request so far.
	Timing             *ResourceTiming  `json:"timing,omitempty"`             // Timing information for the given request.
	Protocol           string           `json:"protocol,omitempty"`           // Protocol used to fetch this request.
	SecurityState      *SecurityState   `json:"securityState"`                // Security state of the request resource.
	SecurityDetails    *SecurityDetails `json:"securityDetails,omitempty"`    // Security details for the request.
}

// WebSocket request data.
// @experimental
type WebSocketRequest struct {
	Headers *Headers `json:"headers"` // HTTP request headers.
}

// WebSocket response data.
// @experimental
type WebSocketResponse struct {
	Status             float64  `json:"status"`                       // HTTP response status code.
	StatusText         string   `json:"statusText"`                   // HTTP response status text.
	Headers            *Headers `json:"headers"`                      // HTTP response headers.
	HeadersText        string   `json:"headersText,omitempty"`        // HTTP response headers text.
	RequestHeaders     *Headers `json:"requestHeaders,omitempty"`     // HTTP request headers.
	RequestHeadersText string   `json:"requestHeadersText,omitempty"` // HTTP request headers text.
}

// WebSocket frame data.
// @experimental
type WebSocketFrame struct {
	Opcode      float64 `json:"opcode"`      // WebSocket frame opcode.
	Mask        bool    `json:"mask"`        // WebSocke frame mask.
	PayloadData string  `json:"payloadData"` // WebSocke frame payload data.
}

// Information about the cached resource.
type CachedResource struct {
	Url      string        `json:"url"`                // Resource URL. This is the url of the original network request.
	Type     *ResourceType `json:"type"`               // Type of this resource.
	Response *Response     `json:"response,omitempty"` // Cached response data.
	BodySize float64       `json:"bodySize"`           // Cached response body size.
}

// Information about the request initiator.
type Initiator struct {
	Type       string      `json:"type"`                 // Type of this initiator.
	Stack      *StackTrace `json:"stack,omitempty"`      // Initiator JavaScript stack trace, set for Script only.
	Url        string      `json:"url,omitempty"`        // Initiator URL, set for Parser type only.
	LineNumber float64     `json:"lineNumber,omitempty"` // Initiator line number, set for Parser type only (0-based).
}

// Cookie object
// @experimental
type Cookie struct {
	Name     string         `json:"name"`               // Cookie name.
	Value    string         `json:"value"`              // Cookie value.
	Domain   string         `json:"domain"`             // Cookie domain.
	Path     string         `json:"path"`               // Cookie path.
	Expires  float64        `json:"expires"`            // Cookie expiration date as the number of seconds since the UNIX epoch.
	Size     int            `json:"size"`               // Cookie size.
	HttpOnly bool           `json:"httpOnly"`           // True if cookie is http-only.
	Secure   bool           `json:"secure"`             // True if cookie is secure.
	Session  bool           `json:"session"`            // True in case of session cookie.
	SameSite CookieSameSite `json:"sameSite,omitempty"` // Cookie SameSite type.
}

type NetworkEnableParams struct {
	MaxTotalBufferSize    int `json:"maxTotalBufferSize,omitempty"`    // Buffer size in bytes to use when preserving network payloads (XHRs, etc).
	MaxResourceBufferSize int `json:"maxResourceBufferSize,omitempty"` // Per-resource buffer size in bytes to use when preserving network payloads (XHRs, etc).
}

// Enables network tracking, network events will now be delivered to the client.

type NetworkEnableCommand struct {
	params *NetworkEnableParams
	wg     sync.WaitGroup
	err    error
}

func NewNetworkEnableCommand(params *NetworkEnableParams) *NetworkEnableCommand {
	return &NetworkEnableCommand{
		params: params,
	}
}

func (cmd *NetworkEnableCommand) Name() string {
	return "Network.enable"
}

func (cmd *NetworkEnableCommand) Params() interface{} {
	return cmd.params
}

func (cmd *NetworkEnableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func NetworkEnable(params *NetworkEnableParams, conn *hc.Conn) (err error) {
	cmd := NewNetworkEnableCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type NetworkEnableCB func(err error)

// Enables network tracking, network events will now be delivered to the client.

type AsyncNetworkEnableCommand struct {
	params *NetworkEnableParams
	cb     NetworkEnableCB
}

func NewAsyncNetworkEnableCommand(params *NetworkEnableParams, cb NetworkEnableCB) *AsyncNetworkEnableCommand {
	return &AsyncNetworkEnableCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncNetworkEnableCommand) Name() string {
	return "Network.enable"
}

func (cmd *AsyncNetworkEnableCommand) Params() interface{} {
	return cmd.params
}

func (cmd *NetworkEnableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncNetworkEnableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Disables network tracking, prevents network events from being sent to the client.

type NetworkDisableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewNetworkDisableCommand() *NetworkDisableCommand {
	return &NetworkDisableCommand{}
}

func (cmd *NetworkDisableCommand) Name() string {
	return "Network.disable"
}

func (cmd *NetworkDisableCommand) Params() interface{} {
	return nil
}

func (cmd *NetworkDisableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func NetworkDisable(conn *hc.Conn) (err error) {
	cmd := NewNetworkDisableCommand()
	cmd.Run(conn)
	return cmd.err
}

type NetworkDisableCB func(err error)

// Disables network tracking, prevents network events from being sent to the client.

type AsyncNetworkDisableCommand struct {
	cb NetworkDisableCB
}

func NewAsyncNetworkDisableCommand(cb NetworkDisableCB) *AsyncNetworkDisableCommand {
	return &AsyncNetworkDisableCommand{
		cb: cb,
	}
}

func (cmd *AsyncNetworkDisableCommand) Name() string {
	return "Network.disable"
}

func (cmd *AsyncNetworkDisableCommand) Params() interface{} {
	return nil
}

func (cmd *NetworkDisableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncNetworkDisableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetUserAgentOverrideParams struct {
	UserAgent string `json:"userAgent"` // User agent to use.
}

// Allows overriding user agent with the given string.

type SetUserAgentOverrideCommand struct {
	params *SetUserAgentOverrideParams
	wg     sync.WaitGroup
	err    error
}

func NewSetUserAgentOverrideCommand(params *SetUserAgentOverrideParams) *SetUserAgentOverrideCommand {
	return &SetUserAgentOverrideCommand{
		params: params,
	}
}

func (cmd *SetUserAgentOverrideCommand) Name() string {
	return "Network.setUserAgentOverride"
}

func (cmd *SetUserAgentOverrideCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetUserAgentOverrideCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetUserAgentOverride(params *SetUserAgentOverrideParams, conn *hc.Conn) (err error) {
	cmd := NewSetUserAgentOverrideCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetUserAgentOverrideCB func(err error)

// Allows overriding user agent with the given string.

type AsyncSetUserAgentOverrideCommand struct {
	params *SetUserAgentOverrideParams
	cb     SetUserAgentOverrideCB
}

func NewAsyncSetUserAgentOverrideCommand(params *SetUserAgentOverrideParams, cb SetUserAgentOverrideCB) *AsyncSetUserAgentOverrideCommand {
	return &AsyncSetUserAgentOverrideCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetUserAgentOverrideCommand) Name() string {
	return "Network.setUserAgentOverride"
}

func (cmd *AsyncSetUserAgentOverrideCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetUserAgentOverrideCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetUserAgentOverrideCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetExtraHTTPHeadersParams struct {
	Headers *Headers `json:"headers"` // Map with extra HTTP headers.
}

// Specifies whether to always send extra HTTP headers with the requests from this page.

type SetExtraHTTPHeadersCommand struct {
	params *SetExtraHTTPHeadersParams
	wg     sync.WaitGroup
	err    error
}

func NewSetExtraHTTPHeadersCommand(params *SetExtraHTTPHeadersParams) *SetExtraHTTPHeadersCommand {
	return &SetExtraHTTPHeadersCommand{
		params: params,
	}
}

func (cmd *SetExtraHTTPHeadersCommand) Name() string {
	return "Network.setExtraHTTPHeaders"
}

func (cmd *SetExtraHTTPHeadersCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetExtraHTTPHeadersCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetExtraHTTPHeaders(params *SetExtraHTTPHeadersParams, conn *hc.Conn) (err error) {
	cmd := NewSetExtraHTTPHeadersCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetExtraHTTPHeadersCB func(err error)

// Specifies whether to always send extra HTTP headers with the requests from this page.

type AsyncSetExtraHTTPHeadersCommand struct {
	params *SetExtraHTTPHeadersParams
	cb     SetExtraHTTPHeadersCB
}

func NewAsyncSetExtraHTTPHeadersCommand(params *SetExtraHTTPHeadersParams, cb SetExtraHTTPHeadersCB) *AsyncSetExtraHTTPHeadersCommand {
	return &AsyncSetExtraHTTPHeadersCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetExtraHTTPHeadersCommand) Name() string {
	return "Network.setExtraHTTPHeaders"
}

func (cmd *AsyncSetExtraHTTPHeadersCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetExtraHTTPHeadersCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetExtraHTTPHeadersCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type GetResponseBodyParams struct {
	RequestId RequestId `json:"requestId"` // Identifier of the network request to get content for.
}

type GetResponseBodyResult struct {
	Body          string `json:"body"`          // Response body.
	Base64Encoded bool   `json:"base64Encoded"` // True, if content was sent as base64.
}

// Returns content served for the given request.

type GetResponseBodyCommand struct {
	params *GetResponseBodyParams
	result GetResponseBodyResult
	wg     sync.WaitGroup
	err    error
}

func NewGetResponseBodyCommand(params *GetResponseBodyParams) *GetResponseBodyCommand {
	return &GetResponseBodyCommand{
		params: params,
	}
}

func (cmd *GetResponseBodyCommand) Name() string {
	return "Network.getResponseBody"
}

func (cmd *GetResponseBodyCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetResponseBodyCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetResponseBody(params *GetResponseBodyParams, conn *hc.Conn) (result *GetResponseBodyResult, err error) {
	cmd := NewGetResponseBodyCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetResponseBodyCB func(result *GetResponseBodyResult, err error)

// Returns content served for the given request.

type AsyncGetResponseBodyCommand struct {
	params *GetResponseBodyParams
	cb     GetResponseBodyCB
}

func NewAsyncGetResponseBodyCommand(params *GetResponseBodyParams, cb GetResponseBodyCB) *AsyncGetResponseBodyCommand {
	return &AsyncGetResponseBodyCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetResponseBodyCommand) Name() string {
	return "Network.getResponseBody"
}

func (cmd *AsyncGetResponseBodyCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetResponseBodyCommand) Result() *GetResponseBodyResult {
	return &cmd.result
}

func (cmd *GetResponseBodyCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetResponseBodyCommand) Done(data []byte, err error) {
	var result GetResponseBodyResult
	if err == nil {
		err = json.Unmarshal(data, &result)
	}
	if cmd.cb == nil {
		logging.Vlog(-1, err)
	} else if err != nil {
		cmd.cb(nil, err)
	} else {
		cmd.cb(&result, nil)
	}
}

type AddBlockedURLParams struct {
	Url string `json:"url"` // URL to block.
}

// Blocks specific URL from loading.
// @experimental
type AddBlockedURLCommand struct {
	params *AddBlockedURLParams
	wg     sync.WaitGroup
	err    error
}

func NewAddBlockedURLCommand(params *AddBlockedURLParams) *AddBlockedURLCommand {
	return &AddBlockedURLCommand{
		params: params,
	}
}

func (cmd *AddBlockedURLCommand) Name() string {
	return "Network.addBlockedURL"
}

func (cmd *AddBlockedURLCommand) Params() interface{} {
	return cmd.params
}

func (cmd *AddBlockedURLCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func AddBlockedURL(params *AddBlockedURLParams, conn *hc.Conn) (err error) {
	cmd := NewAddBlockedURLCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type AddBlockedURLCB func(err error)

// Blocks specific URL from loading.
// @experimental
type AsyncAddBlockedURLCommand struct {
	params *AddBlockedURLParams
	cb     AddBlockedURLCB
}

func NewAsyncAddBlockedURLCommand(params *AddBlockedURLParams, cb AddBlockedURLCB) *AsyncAddBlockedURLCommand {
	return &AsyncAddBlockedURLCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncAddBlockedURLCommand) Name() string {
	return "Network.addBlockedURL"
}

func (cmd *AsyncAddBlockedURLCommand) Params() interface{} {
	return cmd.params
}

func (cmd *AddBlockedURLCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncAddBlockedURLCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type RemoveBlockedURLParams struct {
	Url string `json:"url"` // URL to stop blocking.
}

// Cancels blocking of a specific URL from loading.
// @experimental
type RemoveBlockedURLCommand struct {
	params *RemoveBlockedURLParams
	wg     sync.WaitGroup
	err    error
}

func NewRemoveBlockedURLCommand(params *RemoveBlockedURLParams) *RemoveBlockedURLCommand {
	return &RemoveBlockedURLCommand{
		params: params,
	}
}

func (cmd *RemoveBlockedURLCommand) Name() string {
	return "Network.removeBlockedURL"
}

func (cmd *RemoveBlockedURLCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveBlockedURLCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func RemoveBlockedURL(params *RemoveBlockedURLParams, conn *hc.Conn) (err error) {
	cmd := NewRemoveBlockedURLCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type RemoveBlockedURLCB func(err error)

// Cancels blocking of a specific URL from loading.
// @experimental
type AsyncRemoveBlockedURLCommand struct {
	params *RemoveBlockedURLParams
	cb     RemoveBlockedURLCB
}

func NewAsyncRemoveBlockedURLCommand(params *RemoveBlockedURLParams, cb RemoveBlockedURLCB) *AsyncRemoveBlockedURLCommand {
	return &AsyncRemoveBlockedURLCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncRemoveBlockedURLCommand) Name() string {
	return "Network.removeBlockedURL"
}

func (cmd *AsyncRemoveBlockedURLCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveBlockedURLCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncRemoveBlockedURLCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type ReplayXHRParams struct {
	RequestId RequestId `json:"requestId"` // Identifier of XHR to replay.
}

// This method sends a new XMLHttpRequest which is identical to the original one. The following parameters should be identical: method, url, async, request body, extra headers, withCredentials attribute, user, password.
// @experimental
type ReplayXHRCommand struct {
	params *ReplayXHRParams
	wg     sync.WaitGroup
	err    error
}

func NewReplayXHRCommand(params *ReplayXHRParams) *ReplayXHRCommand {
	return &ReplayXHRCommand{
		params: params,
	}
}

func (cmd *ReplayXHRCommand) Name() string {
	return "Network.replayXHR"
}

func (cmd *ReplayXHRCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ReplayXHRCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ReplayXHR(params *ReplayXHRParams, conn *hc.Conn) (err error) {
	cmd := NewReplayXHRCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type ReplayXHRCB func(err error)

// This method sends a new XMLHttpRequest which is identical to the original one. The following parameters should be identical: method, url, async, request body, extra headers, withCredentials attribute, user, password.
// @experimental
type AsyncReplayXHRCommand struct {
	params *ReplayXHRParams
	cb     ReplayXHRCB
}

func NewAsyncReplayXHRCommand(params *ReplayXHRParams, cb ReplayXHRCB) *AsyncReplayXHRCommand {
	return &AsyncReplayXHRCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncReplayXHRCommand) Name() string {
	return "Network.replayXHR"
}

func (cmd *AsyncReplayXHRCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ReplayXHRCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncReplayXHRCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetMonitoringXHREnabledParams struct {
	Enabled bool `json:"enabled"` // Monitoring enabled state.
}

// Toggles monitoring of XMLHttpRequest. If true, console will receive messages upon each XHR issued.
// @experimental
type SetMonitoringXHREnabledCommand struct {
	params *SetMonitoringXHREnabledParams
	wg     sync.WaitGroup
	err    error
}

func NewSetMonitoringXHREnabledCommand(params *SetMonitoringXHREnabledParams) *SetMonitoringXHREnabledCommand {
	return &SetMonitoringXHREnabledCommand{
		params: params,
	}
}

func (cmd *SetMonitoringXHREnabledCommand) Name() string {
	return "Network.setMonitoringXHREnabled"
}

func (cmd *SetMonitoringXHREnabledCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetMonitoringXHREnabledCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetMonitoringXHREnabled(params *SetMonitoringXHREnabledParams, conn *hc.Conn) (err error) {
	cmd := NewSetMonitoringXHREnabledCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetMonitoringXHREnabledCB func(err error)

// Toggles monitoring of XMLHttpRequest. If true, console will receive messages upon each XHR issued.
// @experimental
type AsyncSetMonitoringXHREnabledCommand struct {
	params *SetMonitoringXHREnabledParams
	cb     SetMonitoringXHREnabledCB
}

func NewAsyncSetMonitoringXHREnabledCommand(params *SetMonitoringXHREnabledParams, cb SetMonitoringXHREnabledCB) *AsyncSetMonitoringXHREnabledCommand {
	return &AsyncSetMonitoringXHREnabledCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetMonitoringXHREnabledCommand) Name() string {
	return "Network.setMonitoringXHREnabled"
}

func (cmd *AsyncSetMonitoringXHREnabledCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetMonitoringXHREnabledCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetMonitoringXHREnabledCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type CanClearBrowserCacheResult struct {
	Result bool `json:"result"` // True if browser cache can be cleared.
}

// Tells whether clearing browser cache is supported.

type CanClearBrowserCacheCommand struct {
	result CanClearBrowserCacheResult
	wg     sync.WaitGroup
	err    error
}

func NewCanClearBrowserCacheCommand() *CanClearBrowserCacheCommand {
	return &CanClearBrowserCacheCommand{}
}

func (cmd *CanClearBrowserCacheCommand) Name() string {
	return "Network.canClearBrowserCache"
}

func (cmd *CanClearBrowserCacheCommand) Params() interface{} {
	return nil
}

func (cmd *CanClearBrowserCacheCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func CanClearBrowserCache(conn *hc.Conn) (result *CanClearBrowserCacheResult, err error) {
	cmd := NewCanClearBrowserCacheCommand()
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type CanClearBrowserCacheCB func(result *CanClearBrowserCacheResult, err error)

// Tells whether clearing browser cache is supported.

type AsyncCanClearBrowserCacheCommand struct {
	cb CanClearBrowserCacheCB
}

func NewAsyncCanClearBrowserCacheCommand(cb CanClearBrowserCacheCB) *AsyncCanClearBrowserCacheCommand {
	return &AsyncCanClearBrowserCacheCommand{
		cb: cb,
	}
}

func (cmd *AsyncCanClearBrowserCacheCommand) Name() string {
	return "Network.canClearBrowserCache"
}

func (cmd *AsyncCanClearBrowserCacheCommand) Params() interface{} {
	return nil
}

func (cmd *CanClearBrowserCacheCommand) Result() *CanClearBrowserCacheResult {
	return &cmd.result
}

func (cmd *CanClearBrowserCacheCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncCanClearBrowserCacheCommand) Done(data []byte, err error) {
	var result CanClearBrowserCacheResult
	if err == nil {
		err = json.Unmarshal(data, &result)
	}
	if cmd.cb == nil {
		logging.Vlog(-1, err)
	} else if err != nil {
		cmd.cb(nil, err)
	} else {
		cmd.cb(&result, nil)
	}
}

// Clears browser cache.

type ClearBrowserCacheCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewClearBrowserCacheCommand() *ClearBrowserCacheCommand {
	return &ClearBrowserCacheCommand{}
}

func (cmd *ClearBrowserCacheCommand) Name() string {
	return "Network.clearBrowserCache"
}

func (cmd *ClearBrowserCacheCommand) Params() interface{} {
	return nil
}

func (cmd *ClearBrowserCacheCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ClearBrowserCache(conn *hc.Conn) (err error) {
	cmd := NewClearBrowserCacheCommand()
	cmd.Run(conn)
	return cmd.err
}

type ClearBrowserCacheCB func(err error)

// Clears browser cache.

type AsyncClearBrowserCacheCommand struct {
	cb ClearBrowserCacheCB
}

func NewAsyncClearBrowserCacheCommand(cb ClearBrowserCacheCB) *AsyncClearBrowserCacheCommand {
	return &AsyncClearBrowserCacheCommand{
		cb: cb,
	}
}

func (cmd *AsyncClearBrowserCacheCommand) Name() string {
	return "Network.clearBrowserCache"
}

func (cmd *AsyncClearBrowserCacheCommand) Params() interface{} {
	return nil
}

func (cmd *ClearBrowserCacheCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncClearBrowserCacheCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type CanClearBrowserCookiesResult struct {
	Result bool `json:"result"` // True if browser cookies can be cleared.
}

// Tells whether clearing browser cookies is supported.

type CanClearBrowserCookiesCommand struct {
	result CanClearBrowserCookiesResult
	wg     sync.WaitGroup
	err    error
}

func NewCanClearBrowserCookiesCommand() *CanClearBrowserCookiesCommand {
	return &CanClearBrowserCookiesCommand{}
}

func (cmd *CanClearBrowserCookiesCommand) Name() string {
	return "Network.canClearBrowserCookies"
}

func (cmd *CanClearBrowserCookiesCommand) Params() interface{} {
	return nil
}

func (cmd *CanClearBrowserCookiesCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func CanClearBrowserCookies(conn *hc.Conn) (result *CanClearBrowserCookiesResult, err error) {
	cmd := NewCanClearBrowserCookiesCommand()
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type CanClearBrowserCookiesCB func(result *CanClearBrowserCookiesResult, err error)

// Tells whether clearing browser cookies is supported.

type AsyncCanClearBrowserCookiesCommand struct {
	cb CanClearBrowserCookiesCB
}

func NewAsyncCanClearBrowserCookiesCommand(cb CanClearBrowserCookiesCB) *AsyncCanClearBrowserCookiesCommand {
	return &AsyncCanClearBrowserCookiesCommand{
		cb: cb,
	}
}

func (cmd *AsyncCanClearBrowserCookiesCommand) Name() string {
	return "Network.canClearBrowserCookies"
}

func (cmd *AsyncCanClearBrowserCookiesCommand) Params() interface{} {
	return nil
}

func (cmd *CanClearBrowserCookiesCommand) Result() *CanClearBrowserCookiesResult {
	return &cmd.result
}

func (cmd *CanClearBrowserCookiesCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncCanClearBrowserCookiesCommand) Done(data []byte, err error) {
	var result CanClearBrowserCookiesResult
	if err == nil {
		err = json.Unmarshal(data, &result)
	}
	if cmd.cb == nil {
		logging.Vlog(-1, err)
	} else if err != nil {
		cmd.cb(nil, err)
	} else {
		cmd.cb(&result, nil)
	}
}

// Clears browser cookies.

type ClearBrowserCookiesCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewClearBrowserCookiesCommand() *ClearBrowserCookiesCommand {
	return &ClearBrowserCookiesCommand{}
}

func (cmd *ClearBrowserCookiesCommand) Name() string {
	return "Network.clearBrowserCookies"
}

func (cmd *ClearBrowserCookiesCommand) Params() interface{} {
	return nil
}

func (cmd *ClearBrowserCookiesCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ClearBrowserCookies(conn *hc.Conn) (err error) {
	cmd := NewClearBrowserCookiesCommand()
	cmd.Run(conn)
	return cmd.err
}

type ClearBrowserCookiesCB func(err error)

// Clears browser cookies.

type AsyncClearBrowserCookiesCommand struct {
	cb ClearBrowserCookiesCB
}

func NewAsyncClearBrowserCookiesCommand(cb ClearBrowserCookiesCB) *AsyncClearBrowserCookiesCommand {
	return &AsyncClearBrowserCookiesCommand{
		cb: cb,
	}
}

func (cmd *AsyncClearBrowserCookiesCommand) Name() string {
	return "Network.clearBrowserCookies"
}

func (cmd *AsyncClearBrowserCookiesCommand) Params() interface{} {
	return nil
}

func (cmd *ClearBrowserCookiesCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncClearBrowserCookiesCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type NetworkGetCookiesResult struct {
	Cookies []*Cookie `json:"cookies"` // Array of cookie objects.
}

// Returns all browser cookies for the current URL. Depending on the backend support, will return detailed cookie information in the cookies field.
// @experimental
type NetworkGetCookiesCommand struct {
	result NetworkGetCookiesResult
	wg     sync.WaitGroup
	err    error
}

func NewNetworkGetCookiesCommand() *NetworkGetCookiesCommand {
	return &NetworkGetCookiesCommand{}
}

func (cmd *NetworkGetCookiesCommand) Name() string {
	return "Network.getCookies"
}

func (cmd *NetworkGetCookiesCommand) Params() interface{} {
	return nil
}

func (cmd *NetworkGetCookiesCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func NetworkGetCookies(conn *hc.Conn) (result *NetworkGetCookiesResult, err error) {
	cmd := NewNetworkGetCookiesCommand()
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type NetworkGetCookiesCB func(result *NetworkGetCookiesResult, err error)

// Returns all browser cookies for the current URL. Depending on the backend support, will return detailed cookie information in the cookies field.
// @experimental
type AsyncNetworkGetCookiesCommand struct {
	cb NetworkGetCookiesCB
}

func NewAsyncNetworkGetCookiesCommand(cb NetworkGetCookiesCB) *AsyncNetworkGetCookiesCommand {
	return &AsyncNetworkGetCookiesCommand{
		cb: cb,
	}
}

func (cmd *AsyncNetworkGetCookiesCommand) Name() string {
	return "Network.getCookies"
}

func (cmd *AsyncNetworkGetCookiesCommand) Params() interface{} {
	return nil
}

func (cmd *NetworkGetCookiesCommand) Result() *NetworkGetCookiesResult {
	return &cmd.result
}

func (cmd *NetworkGetCookiesCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncNetworkGetCookiesCommand) Done(data []byte, err error) {
	var result NetworkGetCookiesResult
	if err == nil {
		err = json.Unmarshal(data, &result)
	}
	if cmd.cb == nil {
		logging.Vlog(-1, err)
	} else if err != nil {
		cmd.cb(nil, err)
	} else {
		cmd.cb(&result, nil)
	}
}

type GetAllCookiesResult struct {
	Cookies []*Cookie `json:"cookies"` // Array of cookie objects.
}

// Returns all browser cookies. Depending on the backend support, will return detailed cookie information in the cookies field.
// @experimental
type GetAllCookiesCommand struct {
	result GetAllCookiesResult
	wg     sync.WaitGroup
	err    error
}

func NewGetAllCookiesCommand() *GetAllCookiesCommand {
	return &GetAllCookiesCommand{}
}

func (cmd *GetAllCookiesCommand) Name() string {
	return "Network.getAllCookies"
}

func (cmd *GetAllCookiesCommand) Params() interface{} {
	return nil
}

func (cmd *GetAllCookiesCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetAllCookies(conn *hc.Conn) (result *GetAllCookiesResult, err error) {
	cmd := NewGetAllCookiesCommand()
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetAllCookiesCB func(result *GetAllCookiesResult, err error)

// Returns all browser cookies. Depending on the backend support, will return detailed cookie information in the cookies field.
// @experimental
type AsyncGetAllCookiesCommand struct {
	cb GetAllCookiesCB
}

func NewAsyncGetAllCookiesCommand(cb GetAllCookiesCB) *AsyncGetAllCookiesCommand {
	return &AsyncGetAllCookiesCommand{
		cb: cb,
	}
}

func (cmd *AsyncGetAllCookiesCommand) Name() string {
	return "Network.getAllCookies"
}

func (cmd *AsyncGetAllCookiesCommand) Params() interface{} {
	return nil
}

func (cmd *GetAllCookiesCommand) Result() *GetAllCookiesResult {
	return &cmd.result
}

func (cmd *GetAllCookiesCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetAllCookiesCommand) Done(data []byte, err error) {
	var result GetAllCookiesResult
	if err == nil {
		err = json.Unmarshal(data, &result)
	}
	if cmd.cb == nil {
		logging.Vlog(-1, err)
	} else if err != nil {
		cmd.cb(nil, err)
	} else {
		cmd.cb(&result, nil)
	}
}

type NetworkDeleteCookieParams struct {
	CookieName string `json:"cookieName"` // Name of the cookie to remove.
	Url        string `json:"url"`        // URL to match cooke domain and path.
}

// Deletes browser cookie with given name, domain and path.
// @experimental
type NetworkDeleteCookieCommand struct {
	params *NetworkDeleteCookieParams
	wg     sync.WaitGroup
	err    error
}

func NewNetworkDeleteCookieCommand(params *NetworkDeleteCookieParams) *NetworkDeleteCookieCommand {
	return &NetworkDeleteCookieCommand{
		params: params,
	}
}

func (cmd *NetworkDeleteCookieCommand) Name() string {
	return "Network.deleteCookie"
}

func (cmd *NetworkDeleteCookieCommand) Params() interface{} {
	return cmd.params
}

func (cmd *NetworkDeleteCookieCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func NetworkDeleteCookie(params *NetworkDeleteCookieParams, conn *hc.Conn) (err error) {
	cmd := NewNetworkDeleteCookieCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type NetworkDeleteCookieCB func(err error)

// Deletes browser cookie with given name, domain and path.
// @experimental
type AsyncNetworkDeleteCookieCommand struct {
	params *NetworkDeleteCookieParams
	cb     NetworkDeleteCookieCB
}

func NewAsyncNetworkDeleteCookieCommand(params *NetworkDeleteCookieParams, cb NetworkDeleteCookieCB) *AsyncNetworkDeleteCookieCommand {
	return &AsyncNetworkDeleteCookieCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncNetworkDeleteCookieCommand) Name() string {
	return "Network.deleteCookie"
}

func (cmd *AsyncNetworkDeleteCookieCommand) Params() interface{} {
	return cmd.params
}

func (cmd *NetworkDeleteCookieCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncNetworkDeleteCookieCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetCookieParams struct {
	Url            string           `json:"url"`                      // The request-URI to associate with the setting of the cookie. This value can affect the default domain and path values of the created cookie.
	Name           string           `json:"name"`                     // The name of the cookie.
	Value          string           `json:"value"`                    // The value of the cookie.
	Domain         string           `json:"domain,omitempty"`         // If omitted, the cookie becomes a host-only cookie.
	Path           string           `json:"path,omitempty"`           // Defaults to the path portion of the url parameter.
	Secure         bool             `json:"secure,omitempty"`         // Defaults ot false.
	HttpOnly       bool             `json:"httpOnly,omitempty"`       // Defaults to false.
	SameSite       CookieSameSite   `json:"sameSite,omitempty"`       // Defaults to browser default behavior.
	ExpirationDate NetworkTimestamp `json:"expirationDate,omitempty"` // If omitted, the cookie becomes a session cookie.
}

type SetCookieResult struct {
	Success bool `json:"success"` // True if successfully set cookie.
}

// Sets a cookie with the given cookie data; may overwrite equivalent cookies if they exist.
// @experimental
type SetCookieCommand struct {
	params *SetCookieParams
	result SetCookieResult
	wg     sync.WaitGroup
	err    error
}

func NewSetCookieCommand(params *SetCookieParams) *SetCookieCommand {
	return &SetCookieCommand{
		params: params,
	}
}

func (cmd *SetCookieCommand) Name() string {
	return "Network.setCookie"
}

func (cmd *SetCookieCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetCookieCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetCookie(params *SetCookieParams, conn *hc.Conn) (result *SetCookieResult, err error) {
	cmd := NewSetCookieCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type SetCookieCB func(result *SetCookieResult, err error)

// Sets a cookie with the given cookie data; may overwrite equivalent cookies if they exist.
// @experimental
type AsyncSetCookieCommand struct {
	params *SetCookieParams
	cb     SetCookieCB
}

func NewAsyncSetCookieCommand(params *SetCookieParams, cb SetCookieCB) *AsyncSetCookieCommand {
	return &AsyncSetCookieCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetCookieCommand) Name() string {
	return "Network.setCookie"
}

func (cmd *AsyncSetCookieCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetCookieCommand) Result() *SetCookieResult {
	return &cmd.result
}

func (cmd *SetCookieCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetCookieCommand) Done(data []byte, err error) {
	var result SetCookieResult
	if err == nil {
		err = json.Unmarshal(data, &result)
	}
	if cmd.cb == nil {
		logging.Vlog(-1, err)
	} else if err != nil {
		cmd.cb(nil, err)
	} else {
		cmd.cb(&result, nil)
	}
}

type CanEmulateNetworkConditionsResult struct {
	Result bool `json:"result"` // True if emulation of network conditions is supported.
}

// Tells whether emulation of network conditions is supported.
// @experimental
type CanEmulateNetworkConditionsCommand struct {
	result CanEmulateNetworkConditionsResult
	wg     sync.WaitGroup
	err    error
}

func NewCanEmulateNetworkConditionsCommand() *CanEmulateNetworkConditionsCommand {
	return &CanEmulateNetworkConditionsCommand{}
}

func (cmd *CanEmulateNetworkConditionsCommand) Name() string {
	return "Network.canEmulateNetworkConditions"
}

func (cmd *CanEmulateNetworkConditionsCommand) Params() interface{} {
	return nil
}

func (cmd *CanEmulateNetworkConditionsCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func CanEmulateNetworkConditions(conn *hc.Conn) (result *CanEmulateNetworkConditionsResult, err error) {
	cmd := NewCanEmulateNetworkConditionsCommand()
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type CanEmulateNetworkConditionsCB func(result *CanEmulateNetworkConditionsResult, err error)

// Tells whether emulation of network conditions is supported.
// @experimental
type AsyncCanEmulateNetworkConditionsCommand struct {
	cb CanEmulateNetworkConditionsCB
}

func NewAsyncCanEmulateNetworkConditionsCommand(cb CanEmulateNetworkConditionsCB) *AsyncCanEmulateNetworkConditionsCommand {
	return &AsyncCanEmulateNetworkConditionsCommand{
		cb: cb,
	}
}

func (cmd *AsyncCanEmulateNetworkConditionsCommand) Name() string {
	return "Network.canEmulateNetworkConditions"
}

func (cmd *AsyncCanEmulateNetworkConditionsCommand) Params() interface{} {
	return nil
}

func (cmd *CanEmulateNetworkConditionsCommand) Result() *CanEmulateNetworkConditionsResult {
	return &cmd.result
}

func (cmd *CanEmulateNetworkConditionsCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncCanEmulateNetworkConditionsCommand) Done(data []byte, err error) {
	var result CanEmulateNetworkConditionsResult
	if err == nil {
		err = json.Unmarshal(data, &result)
	}
	if cmd.cb == nil {
		logging.Vlog(-1, err)
	} else if err != nil {
		cmd.cb(nil, err)
	} else {
		cmd.cb(&result, nil)
	}
}

type EmulateNetworkConditionsParams struct {
	Offline            bool           `json:"offline"`                  // True to emulate internet disconnection.
	Latency            float64        `json:"latency"`                  // Additional latency (ms).
	DownloadThroughput float64        `json:"downloadThroughput"`       // Maximal aggregated download throughput.
	UploadThroughput   float64        `json:"uploadThroughput"`         // Maximal aggregated upload throughput.
	ConnectionType     ConnectionType `json:"connectionType,omitempty"` // Connection type if known.
}

// Activates emulation of network conditions.

type EmulateNetworkConditionsCommand struct {
	params *EmulateNetworkConditionsParams
	wg     sync.WaitGroup
	err    error
}

func NewEmulateNetworkConditionsCommand(params *EmulateNetworkConditionsParams) *EmulateNetworkConditionsCommand {
	return &EmulateNetworkConditionsCommand{
		params: params,
	}
}

func (cmd *EmulateNetworkConditionsCommand) Name() string {
	return "Network.emulateNetworkConditions"
}

func (cmd *EmulateNetworkConditionsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *EmulateNetworkConditionsCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func EmulateNetworkConditions(params *EmulateNetworkConditionsParams, conn *hc.Conn) (err error) {
	cmd := NewEmulateNetworkConditionsCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type EmulateNetworkConditionsCB func(err error)

// Activates emulation of network conditions.

type AsyncEmulateNetworkConditionsCommand struct {
	params *EmulateNetworkConditionsParams
	cb     EmulateNetworkConditionsCB
}

func NewAsyncEmulateNetworkConditionsCommand(params *EmulateNetworkConditionsParams, cb EmulateNetworkConditionsCB) *AsyncEmulateNetworkConditionsCommand {
	return &AsyncEmulateNetworkConditionsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncEmulateNetworkConditionsCommand) Name() string {
	return "Network.emulateNetworkConditions"
}

func (cmd *AsyncEmulateNetworkConditionsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *EmulateNetworkConditionsCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncEmulateNetworkConditionsCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetCacheDisabledParams struct {
	CacheDisabled bool `json:"cacheDisabled"` // Cache disabled state.
}

// Toggles ignoring cache for each request. If true, cache will not be used.

type SetCacheDisabledCommand struct {
	params *SetCacheDisabledParams
	wg     sync.WaitGroup
	err    error
}

func NewSetCacheDisabledCommand(params *SetCacheDisabledParams) *SetCacheDisabledCommand {
	return &SetCacheDisabledCommand{
		params: params,
	}
}

func (cmd *SetCacheDisabledCommand) Name() string {
	return "Network.setCacheDisabled"
}

func (cmd *SetCacheDisabledCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetCacheDisabledCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetCacheDisabled(params *SetCacheDisabledParams, conn *hc.Conn) (err error) {
	cmd := NewSetCacheDisabledCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetCacheDisabledCB func(err error)

// Toggles ignoring cache for each request. If true, cache will not be used.

type AsyncSetCacheDisabledCommand struct {
	params *SetCacheDisabledParams
	cb     SetCacheDisabledCB
}

func NewAsyncSetCacheDisabledCommand(params *SetCacheDisabledParams, cb SetCacheDisabledCB) *AsyncSetCacheDisabledCommand {
	return &AsyncSetCacheDisabledCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetCacheDisabledCommand) Name() string {
	return "Network.setCacheDisabled"
}

func (cmd *AsyncSetCacheDisabledCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetCacheDisabledCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetCacheDisabledCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetBypassServiceWorkerParams struct {
	Bypass bool `json:"bypass"` // Bypass service worker and load from network.
}

// Toggles ignoring of service worker for each request.
// @experimental
type SetBypassServiceWorkerCommand struct {
	params *SetBypassServiceWorkerParams
	wg     sync.WaitGroup
	err    error
}

func NewSetBypassServiceWorkerCommand(params *SetBypassServiceWorkerParams) *SetBypassServiceWorkerCommand {
	return &SetBypassServiceWorkerCommand{
		params: params,
	}
}

func (cmd *SetBypassServiceWorkerCommand) Name() string {
	return "Network.setBypassServiceWorker"
}

func (cmd *SetBypassServiceWorkerCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetBypassServiceWorkerCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetBypassServiceWorker(params *SetBypassServiceWorkerParams, conn *hc.Conn) (err error) {
	cmd := NewSetBypassServiceWorkerCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetBypassServiceWorkerCB func(err error)

// Toggles ignoring of service worker for each request.
// @experimental
type AsyncSetBypassServiceWorkerCommand struct {
	params *SetBypassServiceWorkerParams
	cb     SetBypassServiceWorkerCB
}

func NewAsyncSetBypassServiceWorkerCommand(params *SetBypassServiceWorkerParams, cb SetBypassServiceWorkerCB) *AsyncSetBypassServiceWorkerCommand {
	return &AsyncSetBypassServiceWorkerCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetBypassServiceWorkerCommand) Name() string {
	return "Network.setBypassServiceWorker"
}

func (cmd *AsyncSetBypassServiceWorkerCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetBypassServiceWorkerCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetBypassServiceWorkerCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetDataSizeLimitsForTestParams struct {
	MaxTotalSize    int `json:"maxTotalSize"`    // Maximum total buffer size.
	MaxResourceSize int `json:"maxResourceSize"` // Maximum per-resource size.
}

// For testing.
// @experimental
type SetDataSizeLimitsForTestCommand struct {
	params *SetDataSizeLimitsForTestParams
	wg     sync.WaitGroup
	err    error
}

func NewSetDataSizeLimitsForTestCommand(params *SetDataSizeLimitsForTestParams) *SetDataSizeLimitsForTestCommand {
	return &SetDataSizeLimitsForTestCommand{
		params: params,
	}
}

func (cmd *SetDataSizeLimitsForTestCommand) Name() string {
	return "Network.setDataSizeLimitsForTest"
}

func (cmd *SetDataSizeLimitsForTestCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetDataSizeLimitsForTestCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetDataSizeLimitsForTest(params *SetDataSizeLimitsForTestParams, conn *hc.Conn) (err error) {
	cmd := NewSetDataSizeLimitsForTestCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetDataSizeLimitsForTestCB func(err error)

// For testing.
// @experimental
type AsyncSetDataSizeLimitsForTestCommand struct {
	params *SetDataSizeLimitsForTestParams
	cb     SetDataSizeLimitsForTestCB
}

func NewAsyncSetDataSizeLimitsForTestCommand(params *SetDataSizeLimitsForTestParams, cb SetDataSizeLimitsForTestCB) *AsyncSetDataSizeLimitsForTestCommand {
	return &AsyncSetDataSizeLimitsForTestCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetDataSizeLimitsForTestCommand) Name() string {
	return "Network.setDataSizeLimitsForTest"
}

func (cmd *AsyncSetDataSizeLimitsForTestCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetDataSizeLimitsForTestCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetDataSizeLimitsForTestCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type GetCertificateParams struct {
	Origin string `json:"origin"` // Origin to get certificate for.
}

type GetCertificateResult struct {
	TableNames []string `json:"tableNames"`
}

// Returns the DER-encoded certificate.
// @experimental
type GetCertificateCommand struct {
	params *GetCertificateParams
	result GetCertificateResult
	wg     sync.WaitGroup
	err    error
}

func NewGetCertificateCommand(params *GetCertificateParams) *GetCertificateCommand {
	return &GetCertificateCommand{
		params: params,
	}
}

func (cmd *GetCertificateCommand) Name() string {
	return "Network.getCertificate"
}

func (cmd *GetCertificateCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetCertificateCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetCertificate(params *GetCertificateParams, conn *hc.Conn) (result *GetCertificateResult, err error) {
	cmd := NewGetCertificateCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetCertificateCB func(result *GetCertificateResult, err error)

// Returns the DER-encoded certificate.
// @experimental
type AsyncGetCertificateCommand struct {
	params *GetCertificateParams
	cb     GetCertificateCB
}

func NewAsyncGetCertificateCommand(params *GetCertificateParams, cb GetCertificateCB) *AsyncGetCertificateCommand {
	return &AsyncGetCertificateCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetCertificateCommand) Name() string {
	return "Network.getCertificate"
}

func (cmd *AsyncGetCertificateCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetCertificateCommand) Result() *GetCertificateResult {
	return &cmd.result
}

func (cmd *GetCertificateCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetCertificateCommand) Done(data []byte, err error) {
	var result GetCertificateResult
	if err == nil {
		err = json.Unmarshal(data, &result)
	}
	if cmd.cb == nil {
		logging.Vlog(-1, err)
	} else if err != nil {
		cmd.cb(nil, err)
	} else {
		cmd.cb(&result, nil)
	}
}

// Fired when resource loading priority is changed
// @experimental
type ResourceChangedPriorityEvent struct {
	RequestId   RequestId        `json:"requestId"`   // Request identifier.
	NewPriority ResourcePriority `json:"newPriority"` // New priority
	Timestamp   NetworkTimestamp `json:"timestamp"`   // Timestamp.
}

func OnResourceChangedPriority(conn *hc.Conn, cb func(evt *ResourceChangedPriorityEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &ResourceChangedPriorityEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Network.resourceChangedPriority", sink)
}

// Fired when page is about to send HTTP request.

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

func OnRequestWillBeSent(conn *hc.Conn, cb func(evt *RequestWillBeSentEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &RequestWillBeSentEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Network.requestWillBeSent", sink)
}

// Fired if request ended up loading from cache.

type RequestServedFromCacheEvent struct {
	RequestId RequestId `json:"requestId"` // Request identifier.
}

func OnRequestServedFromCache(conn *hc.Conn, cb func(evt *RequestServedFromCacheEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &RequestServedFromCacheEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Network.requestServedFromCache", sink)
}

// Fired when HTTP response is available.

type ResponseReceivedEvent struct {
	RequestId RequestId        `json:"requestId"` // Request identifier.
	FrameId   *FrameId         `json:"frameId"`   // Frame identifier.
	LoaderId  LoaderId         `json:"loaderId"`  // Loader identifier.
	Timestamp NetworkTimestamp `json:"timestamp"` // Timestamp.
	Type      *ResourceType    `json:"type"`      // Resource type.
	Response  *Response        `json:"response"`  // Response data.
}

func OnResponseReceived(conn *hc.Conn, cb func(evt *ResponseReceivedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &ResponseReceivedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Network.responseReceived", sink)
}

// Fired when data chunk was received over the network.

type DataReceivedEvent struct {
	RequestId         RequestId        `json:"requestId"`         // Request identifier.
	Timestamp         NetworkTimestamp `json:"timestamp"`         // Timestamp.
	DataLength        int              `json:"dataLength"`        // Data chunk length.
	EncodedDataLength int              `json:"encodedDataLength"` // Actual bytes received (might be less than dataLength for compressed encodings).
}

func OnDataReceived(conn *hc.Conn, cb func(evt *DataReceivedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &DataReceivedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Network.dataReceived", sink)
}

// Fired when HTTP request has finished loading.

type LoadingFinishedEvent struct {
	RequestId         RequestId        `json:"requestId"`         // Request identifier.
	Timestamp         NetworkTimestamp `json:"timestamp"`         // Timestamp.
	EncodedDataLength float64          `json:"encodedDataLength"` // Total number of bytes received for this request.
}

func OnLoadingFinished(conn *hc.Conn, cb func(evt *LoadingFinishedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &LoadingFinishedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Network.loadingFinished", sink)
}

// Fired when HTTP request has failed to load.

type LoadingFailedEvent struct {
	RequestId     RequestId        `json:"requestId"`     // Request identifier.
	Timestamp     NetworkTimestamp `json:"timestamp"`     // Timestamp.
	Type          *ResourceType    `json:"type"`          // Resource type.
	ErrorText     string           `json:"errorText"`     // User friendly error message.
	Canceled      bool             `json:"canceled"`      // True if loading was canceled.
	BlockedReason BlockedReason    `json:"blockedReason"` // The reason why loading was blocked, if any.
}

func OnLoadingFailed(conn *hc.Conn, cb func(evt *LoadingFailedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &LoadingFailedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Network.loadingFailed", sink)
}

// Fired when WebSocket is about to initiate handshake.
// @experimental
type WebSocketWillSendHandshakeRequestEvent struct {
	RequestId RequestId         `json:"requestId"` // Request identifier.
	Timestamp NetworkTimestamp  `json:"timestamp"` // Timestamp.
	WallTime  NetworkTimestamp  `json:"wallTime"`  // UTC Timestamp.
	Request   *WebSocketRequest `json:"request"`   // WebSocket request data.
}

func OnWebSocketWillSendHandshakeRequest(conn *hc.Conn, cb func(evt *WebSocketWillSendHandshakeRequestEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &WebSocketWillSendHandshakeRequestEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Network.webSocketWillSendHandshakeRequest", sink)
}

// Fired when WebSocket handshake response becomes available.
// @experimental
type WebSocketHandshakeResponseReceivedEvent struct {
	RequestId RequestId          `json:"requestId"` // Request identifier.
	Timestamp NetworkTimestamp   `json:"timestamp"` // Timestamp.
	Response  *WebSocketResponse `json:"response"`  // WebSocket response data.
}

func OnWebSocketHandshakeResponseReceived(conn *hc.Conn, cb func(evt *WebSocketHandshakeResponseReceivedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &WebSocketHandshakeResponseReceivedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Network.webSocketHandshakeResponseReceived", sink)
}

// Fired upon WebSocket creation.
// @experimental
type WebSocketCreatedEvent struct {
	RequestId RequestId  `json:"requestId"` // Request identifier.
	Url       string     `json:"url"`       // WebSocket request URL.
	Initiator *Initiator `json:"initiator"` // Request initiator.
}

func OnWebSocketCreated(conn *hc.Conn, cb func(evt *WebSocketCreatedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &WebSocketCreatedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Network.webSocketCreated", sink)
}

// Fired when WebSocket is closed.
// @experimental
type WebSocketClosedEvent struct {
	RequestId RequestId        `json:"requestId"` // Request identifier.
	Timestamp NetworkTimestamp `json:"timestamp"` // Timestamp.
}

func OnWebSocketClosed(conn *hc.Conn, cb func(evt *WebSocketClosedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &WebSocketClosedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Network.webSocketClosed", sink)
}

// Fired when WebSocket frame is received.
// @experimental
type WebSocketFrameReceivedEvent struct {
	RequestId RequestId        `json:"requestId"` // Request identifier.
	Timestamp NetworkTimestamp `json:"timestamp"` // Timestamp.
	Response  *WebSocketFrame  `json:"response"`  // WebSocket response data.
}

func OnWebSocketFrameReceived(conn *hc.Conn, cb func(evt *WebSocketFrameReceivedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &WebSocketFrameReceivedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Network.webSocketFrameReceived", sink)
}

// Fired when WebSocket frame error occurs.
// @experimental
type WebSocketFrameErrorEvent struct {
	RequestId    RequestId        `json:"requestId"`    // Request identifier.
	Timestamp    NetworkTimestamp `json:"timestamp"`    // Timestamp.
	ErrorMessage string           `json:"errorMessage"` // WebSocket frame error message.
}

func OnWebSocketFrameError(conn *hc.Conn, cb func(evt *WebSocketFrameErrorEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &WebSocketFrameErrorEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Network.webSocketFrameError", sink)
}

// Fired when WebSocket frame is sent.
// @experimental
type WebSocketFrameSentEvent struct {
	RequestId RequestId        `json:"requestId"` // Request identifier.
	Timestamp NetworkTimestamp `json:"timestamp"` // Timestamp.
	Response  *WebSocketFrame  `json:"response"`  // WebSocket response data.
}

func OnWebSocketFrameSent(conn *hc.Conn, cb func(evt *WebSocketFrameSentEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &WebSocketFrameSentEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Network.webSocketFrameSent", sink)
}

// Fired when EventSource message is received.
// @experimental
type EventSourceMessageReceivedEvent struct {
	RequestId RequestId        `json:"requestId"` // Request identifier.
	Timestamp NetworkTimestamp `json:"timestamp"` // Timestamp.
	EventName string           `json:"eventName"` // Message type.
	EventId   string           `json:"eventId"`   // Message identifier.
	Data      string           `json:"data"`      // Message content.
}

func OnEventSourceMessageReceived(conn *hc.Conn, cb func(evt *EventSourceMessageReceivedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &EventSourceMessageReceivedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Network.eventSourceMessageReceived", sink)
}
