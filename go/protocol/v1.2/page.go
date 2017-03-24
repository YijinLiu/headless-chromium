package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

// Resource type as it was perceived by the rendering engine.
type ResourceType string

const ResourceTypeDocument ResourceType = "Document"
const ResourceTypeStylesheet ResourceType = "Stylesheet"
const ResourceTypeImage ResourceType = "Image"
const ResourceTypeMedia ResourceType = "Media"
const ResourceTypeFont ResourceType = "Font"
const ResourceTypeScript ResourceType = "Script"
const ResourceTypeTextTrack ResourceType = "TextTrack"
const ResourceTypeXHR ResourceType = "XHR"
const ResourceTypeFetch ResourceType = "Fetch"
const ResourceTypeEventSource ResourceType = "EventSource"
const ResourceTypeWebSocket ResourceType = "WebSocket"
const ResourceTypeManifest ResourceType = "Manifest"
const ResourceTypeOther ResourceType = "Other"

// Unique frame identifier.
type FrameId string

// Information about the Frame on the page.
type Frame struct {
	Id             string    `json:"id"`                 // Frame unique identifier.
	ParentId       string    `json:"parentId,omitempty"` // Parent frame identifier.
	LoaderId       *LoaderId `json:"loaderId"`           // Identifier of the loader associated with this frame.
	Name           string    `json:"name,omitempty"`     // Frame's name as specified in the tag.
	Url            string    `json:"url"`                // Frame document's URL.
	SecurityOrigin string    `json:"securityOrigin"`     // Frame document's security origin.
	MimeType       string    `json:"mimeType"`           // Frame document's mimeType as determined by the browser.
}

// Information about the Resource on the page.
// @experimental
type FrameResource struct {
	Url          string            `json:"url"`                    // Resource URL.
	Type         ResourceType      `json:"type"`                   // Type of this resource.
	MimeType     string            `json:"mimeType"`               // Resource mimeType as determined by the browser.
	LastModified *NetworkTimestamp `json:"lastModified,omitempty"` // last-modified timestamp as reported by server.
	ContentSize  float64           `json:"contentSize,omitempty"`  // Resource content size.
	Failed       bool              `json:"failed,omitempty"`       // True if the resource failed to load.
	Canceled     bool              `json:"canceled,omitempty"`     // True if the resource was canceled during loading.
}

// Information about the Frame hierarchy along with their cached resources.
// @experimental
type FrameResourceTree struct {
	Frame       *Frame               `json:"frame"`                 // Frame information for this tree item.
	ChildFrames []*FrameResourceTree `json:"childFrames,omitempty"` // Child frames.
	Resources   []*FrameResource     `json:"resources"`             // Information about frame resources.
}

// Unique script identifier.
// @experimental
type ScriptIdentifier string

// Navigation history entry.
// @experimental
type NavigationEntry struct {
	Id    int    `json:"id"`    // Unique id of the navigation history entry.
	Url   string `json:"url"`   // URL of the navigation history entry.
	Title string `json:"title"` // Title of the navigation history entry.
}

// Screencast frame metadata.
// @experimental
type ScreencastFrameMetadata struct {
	OffsetTop       float64 `json:"offsetTop"`           // Top offset in DIP.
	PageScaleFactor float64 `json:"pageScaleFactor"`     // Page scale factor.
	DeviceWidth     float64 `json:"deviceWidth"`         // Device screen width in DIP.
	DeviceHeight    float64 `json:"deviceHeight"`        // Device screen height in DIP.
	ScrollOffsetX   float64 `json:"scrollOffsetX"`       // Position of horizontal scroll in CSS pixels.
	ScrollOffsetY   float64 `json:"scrollOffsetY"`       // Position of vertical scroll in CSS pixels.
	Timestamp       float64 `json:"timestamp,omitempty"` // Frame swap timestamp.
}

// Javascript dialog type.
// @experimental
type DialogType string

const DialogTypeAlert DialogType = "alert"
const DialogTypeConfirm DialogType = "confirm"
const DialogTypePrompt DialogType = "prompt"
const DialogTypeBeforeunload DialogType = "beforeunload"

// Error while paring app manifest.
// @experimental
type AppManifestError struct {
	Message  string `json:"message"`  // Error message.
	Critical int    `json:"critical"` // If criticial, this is a non-recoverable parse error.
	Line     int    `json:"line"`     // Error line.
	Column   int    `json:"column"`   // Error column.
}

// Proceed: allow the navigation; Cancel: cancel the navigation; CancelAndIgnore: cancels the navigation and makes the requester of the navigation acts like the request was never made.
// @experimental
type NavigationResponse string

const NavigationResponseProceed NavigationResponse = "Proceed"
const NavigationResponseCancel NavigationResponse = "Cancel"
const NavigationResponseCancelAndIgnore NavigationResponse = "CancelAndIgnore"

// Layout viewport position and dimensions.
// @experimental
type LayoutViewport struct {
	PageX        int `json:"pageX"`        // Horizontal offset relative to the document (CSS pixels).
	PageY        int `json:"pageY"`        // Vertical offset relative to the document (CSS pixels).
	ClientWidth  int `json:"clientWidth"`  // Width (CSS pixels), excludes scrollbar if present.
	ClientHeight int `json:"clientHeight"` // Height (CSS pixels), excludes scrollbar if present.
}

// Visual viewport position, dimensions, and scale.
// @experimental
type VisualViewport struct {
	OffsetX      float64 `json:"offsetX"`      // Horizontal offset relative to the layout viewport (CSS pixels).
	OffsetY      float64 `json:"offsetY"`      // Vertical offset relative to the layout viewport (CSS pixels).
	PageX        float64 `json:"pageX"`        // Horizontal offset relative to the document (CSS pixels).
	PageY        float64 `json:"pageY"`        // Vertical offset relative to the document (CSS pixels).
	ClientWidth  float64 `json:"clientWidth"`  // Width (CSS pixels), excludes scrollbar if present.
	ClientHeight float64 `json:"clientHeight"` // Height (CSS pixels), excludes scrollbar if present.
	Scale        float64 `json:"scale"`        // Scale relative to the ideal viewport (size at width=device-width).
}

// Enables page domain notifications.

type PageEnableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewPageEnableCommand() *PageEnableCommand {
	return &PageEnableCommand{}
}

func (cmd *PageEnableCommand) Name() string {
	return "Page.enable"
}

func (cmd *PageEnableCommand) Params() interface{} {
	return nil
}

func (cmd *PageEnableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func PageEnable(conn *hc.Conn) (err error) {
	cmd := NewPageEnableCommand()
	cmd.Run(conn)
	return cmd.err
}

type PageEnableCB func(err error)

// Enables page domain notifications.

type AsyncPageEnableCommand struct {
	cb PageEnableCB
}

func NewAsyncPageEnableCommand(cb PageEnableCB) *AsyncPageEnableCommand {
	return &AsyncPageEnableCommand{
		cb: cb,
	}
}

func (cmd *AsyncPageEnableCommand) Name() string {
	return "Page.enable"
}

func (cmd *AsyncPageEnableCommand) Params() interface{} {
	return nil
}

func (cmd *PageEnableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncPageEnableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Disables page domain notifications.

type PageDisableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewPageDisableCommand() *PageDisableCommand {
	return &PageDisableCommand{}
}

func (cmd *PageDisableCommand) Name() string {
	return "Page.disable"
}

func (cmd *PageDisableCommand) Params() interface{} {
	return nil
}

func (cmd *PageDisableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func PageDisable(conn *hc.Conn) (err error) {
	cmd := NewPageDisableCommand()
	cmd.Run(conn)
	return cmd.err
}

type PageDisableCB func(err error)

// Disables page domain notifications.

type AsyncPageDisableCommand struct {
	cb PageDisableCB
}

func NewAsyncPageDisableCommand(cb PageDisableCB) *AsyncPageDisableCommand {
	return &AsyncPageDisableCommand{
		cb: cb,
	}
}

func (cmd *AsyncPageDisableCommand) Name() string {
	return "Page.disable"
}

func (cmd *AsyncPageDisableCommand) Params() interface{} {
	return nil
}

func (cmd *PageDisableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncPageDisableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type AddScriptToEvaluateOnLoadParams struct {
	ScriptSource string `json:"scriptSource"`
}

type AddScriptToEvaluateOnLoadResult struct {
	Identifier ScriptIdentifier `json:"identifier"` // Identifier of the added script.
}

// @experimental
type AddScriptToEvaluateOnLoadCommand struct {
	params *AddScriptToEvaluateOnLoadParams
	result AddScriptToEvaluateOnLoadResult
	wg     sync.WaitGroup
	err    error
}

func NewAddScriptToEvaluateOnLoadCommand(params *AddScriptToEvaluateOnLoadParams) *AddScriptToEvaluateOnLoadCommand {
	return &AddScriptToEvaluateOnLoadCommand{
		params: params,
	}
}

func (cmd *AddScriptToEvaluateOnLoadCommand) Name() string {
	return "Page.addScriptToEvaluateOnLoad"
}

func (cmd *AddScriptToEvaluateOnLoadCommand) Params() interface{} {
	return cmd.params
}

func (cmd *AddScriptToEvaluateOnLoadCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func AddScriptToEvaluateOnLoad(params *AddScriptToEvaluateOnLoadParams, conn *hc.Conn) (result *AddScriptToEvaluateOnLoadResult, err error) {
	cmd := NewAddScriptToEvaluateOnLoadCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type AddScriptToEvaluateOnLoadCB func(result *AddScriptToEvaluateOnLoadResult, err error)

// @experimental
type AsyncAddScriptToEvaluateOnLoadCommand struct {
	params *AddScriptToEvaluateOnLoadParams
	cb     AddScriptToEvaluateOnLoadCB
}

func NewAsyncAddScriptToEvaluateOnLoadCommand(params *AddScriptToEvaluateOnLoadParams, cb AddScriptToEvaluateOnLoadCB) *AsyncAddScriptToEvaluateOnLoadCommand {
	return &AsyncAddScriptToEvaluateOnLoadCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncAddScriptToEvaluateOnLoadCommand) Name() string {
	return "Page.addScriptToEvaluateOnLoad"
}

func (cmd *AsyncAddScriptToEvaluateOnLoadCommand) Params() interface{} {
	return cmd.params
}

func (cmd *AddScriptToEvaluateOnLoadCommand) Result() *AddScriptToEvaluateOnLoadResult {
	return &cmd.result
}

func (cmd *AddScriptToEvaluateOnLoadCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncAddScriptToEvaluateOnLoadCommand) Done(data []byte, err error) {
	var result AddScriptToEvaluateOnLoadResult
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

type RemoveScriptToEvaluateOnLoadParams struct {
	Identifier ScriptIdentifier `json:"identifier"`
}

// @experimental
type RemoveScriptToEvaluateOnLoadCommand struct {
	params *RemoveScriptToEvaluateOnLoadParams
	wg     sync.WaitGroup
	err    error
}

func NewRemoveScriptToEvaluateOnLoadCommand(params *RemoveScriptToEvaluateOnLoadParams) *RemoveScriptToEvaluateOnLoadCommand {
	return &RemoveScriptToEvaluateOnLoadCommand{
		params: params,
	}
}

func (cmd *RemoveScriptToEvaluateOnLoadCommand) Name() string {
	return "Page.removeScriptToEvaluateOnLoad"
}

func (cmd *RemoveScriptToEvaluateOnLoadCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveScriptToEvaluateOnLoadCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func RemoveScriptToEvaluateOnLoad(params *RemoveScriptToEvaluateOnLoadParams, conn *hc.Conn) (err error) {
	cmd := NewRemoveScriptToEvaluateOnLoadCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type RemoveScriptToEvaluateOnLoadCB func(err error)

// @experimental
type AsyncRemoveScriptToEvaluateOnLoadCommand struct {
	params *RemoveScriptToEvaluateOnLoadParams
	cb     RemoveScriptToEvaluateOnLoadCB
}

func NewAsyncRemoveScriptToEvaluateOnLoadCommand(params *RemoveScriptToEvaluateOnLoadParams, cb RemoveScriptToEvaluateOnLoadCB) *AsyncRemoveScriptToEvaluateOnLoadCommand {
	return &AsyncRemoveScriptToEvaluateOnLoadCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncRemoveScriptToEvaluateOnLoadCommand) Name() string {
	return "Page.removeScriptToEvaluateOnLoad"
}

func (cmd *AsyncRemoveScriptToEvaluateOnLoadCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveScriptToEvaluateOnLoadCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncRemoveScriptToEvaluateOnLoadCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetAutoAttachToCreatedPagesParams struct {
	AutoAttach bool `json:"autoAttach"` // If true, browser will open a new inspector window for every page created from this one.
}

// Controls whether browser will open a new inspector window for connected pages.
// @experimental
type SetAutoAttachToCreatedPagesCommand struct {
	params *SetAutoAttachToCreatedPagesParams
	wg     sync.WaitGroup
	err    error
}

func NewSetAutoAttachToCreatedPagesCommand(params *SetAutoAttachToCreatedPagesParams) *SetAutoAttachToCreatedPagesCommand {
	return &SetAutoAttachToCreatedPagesCommand{
		params: params,
	}
}

func (cmd *SetAutoAttachToCreatedPagesCommand) Name() string {
	return "Page.setAutoAttachToCreatedPages"
}

func (cmd *SetAutoAttachToCreatedPagesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetAutoAttachToCreatedPagesCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetAutoAttachToCreatedPages(params *SetAutoAttachToCreatedPagesParams, conn *hc.Conn) (err error) {
	cmd := NewSetAutoAttachToCreatedPagesCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetAutoAttachToCreatedPagesCB func(err error)

// Controls whether browser will open a new inspector window for connected pages.
// @experimental
type AsyncSetAutoAttachToCreatedPagesCommand struct {
	params *SetAutoAttachToCreatedPagesParams
	cb     SetAutoAttachToCreatedPagesCB
}

func NewAsyncSetAutoAttachToCreatedPagesCommand(params *SetAutoAttachToCreatedPagesParams, cb SetAutoAttachToCreatedPagesCB) *AsyncSetAutoAttachToCreatedPagesCommand {
	return &AsyncSetAutoAttachToCreatedPagesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetAutoAttachToCreatedPagesCommand) Name() string {
	return "Page.setAutoAttachToCreatedPages"
}

func (cmd *AsyncSetAutoAttachToCreatedPagesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetAutoAttachToCreatedPagesCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetAutoAttachToCreatedPagesCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type ReloadParams struct {
	IgnoreCache            bool   `json:"ignoreCache,omitempty"`            // If true, browser cache is ignored (as if the user pressed Shift+refresh).
	ScriptToEvaluateOnLoad string `json:"scriptToEvaluateOnLoad,omitempty"` // If set, the script will be injected into all frames of the inspected page after reload.
}

// Reloads given page optionally ignoring the cache.

type ReloadCommand struct {
	params *ReloadParams
	wg     sync.WaitGroup
	err    error
}

func NewReloadCommand(params *ReloadParams) *ReloadCommand {
	return &ReloadCommand{
		params: params,
	}
}

func (cmd *ReloadCommand) Name() string {
	return "Page.reload"
}

func (cmd *ReloadCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ReloadCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func Reload(params *ReloadParams, conn *hc.Conn) (err error) {
	cmd := NewReloadCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type ReloadCB func(err error)

// Reloads given page optionally ignoring the cache.

type AsyncReloadCommand struct {
	params *ReloadParams
	cb     ReloadCB
}

func NewAsyncReloadCommand(params *ReloadParams, cb ReloadCB) *AsyncReloadCommand {
	return &AsyncReloadCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncReloadCommand) Name() string {
	return "Page.reload"
}

func (cmd *AsyncReloadCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ReloadCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncReloadCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type NavigateParams struct {
	Url string `json:"url"` // URL to navigate the page to.
}

type NavigateResult struct {
	FrameId FrameId `json:"frameId"` // Frame id that will be navigated.
}

// Navigates current page to the given URL.

type NavigateCommand struct {
	params *NavigateParams
	result NavigateResult
	wg     sync.WaitGroup
	err    error
}

func NewNavigateCommand(params *NavigateParams) *NavigateCommand {
	return &NavigateCommand{
		params: params,
	}
}

func (cmd *NavigateCommand) Name() string {
	return "Page.navigate"
}

func (cmd *NavigateCommand) Params() interface{} {
	return cmd.params
}

func (cmd *NavigateCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func Navigate(params *NavigateParams, conn *hc.Conn) (result *NavigateResult, err error) {
	cmd := NewNavigateCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type NavigateCB func(result *NavigateResult, err error)

// Navigates current page to the given URL.

type AsyncNavigateCommand struct {
	params *NavigateParams
	cb     NavigateCB
}

func NewAsyncNavigateCommand(params *NavigateParams, cb NavigateCB) *AsyncNavigateCommand {
	return &AsyncNavigateCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncNavigateCommand) Name() string {
	return "Page.navigate"
}

func (cmd *AsyncNavigateCommand) Params() interface{} {
	return cmd.params
}

func (cmd *NavigateCommand) Result() *NavigateResult {
	return &cmd.result
}

func (cmd *NavigateCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncNavigateCommand) Done(data []byte, err error) {
	var result NavigateResult
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

// Force the page stop all navigations and pending resource fetches.
// @experimental
type StopLoadingCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewStopLoadingCommand() *StopLoadingCommand {
	return &StopLoadingCommand{}
}

func (cmd *StopLoadingCommand) Name() string {
	return "Page.stopLoading"
}

func (cmd *StopLoadingCommand) Params() interface{} {
	return nil
}

func (cmd *StopLoadingCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func StopLoading(conn *hc.Conn) (err error) {
	cmd := NewStopLoadingCommand()
	cmd.Run(conn)
	return cmd.err
}

type StopLoadingCB func(err error)

// Force the page stop all navigations and pending resource fetches.
// @experimental
type AsyncStopLoadingCommand struct {
	cb StopLoadingCB
}

func NewAsyncStopLoadingCommand(cb StopLoadingCB) *AsyncStopLoadingCommand {
	return &AsyncStopLoadingCommand{
		cb: cb,
	}
}

func (cmd *AsyncStopLoadingCommand) Name() string {
	return "Page.stopLoading"
}

func (cmd *AsyncStopLoadingCommand) Params() interface{} {
	return nil
}

func (cmd *StopLoadingCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncStopLoadingCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type GetNavigationHistoryResult struct {
	CurrentIndex int                `json:"currentIndex"` // Index of the current navigation history entry.
	Entries      []*NavigationEntry `json:"entries"`      // Array of navigation history entries.
}

// Returns navigation history for the current page.
// @experimental
type GetNavigationHistoryCommand struct {
	result GetNavigationHistoryResult
	wg     sync.WaitGroup
	err    error
}

func NewGetNavigationHistoryCommand() *GetNavigationHistoryCommand {
	return &GetNavigationHistoryCommand{}
}

func (cmd *GetNavigationHistoryCommand) Name() string {
	return "Page.getNavigationHistory"
}

func (cmd *GetNavigationHistoryCommand) Params() interface{} {
	return nil
}

func (cmd *GetNavigationHistoryCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetNavigationHistory(conn *hc.Conn) (result *GetNavigationHistoryResult, err error) {
	cmd := NewGetNavigationHistoryCommand()
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetNavigationHistoryCB func(result *GetNavigationHistoryResult, err error)

// Returns navigation history for the current page.
// @experimental
type AsyncGetNavigationHistoryCommand struct {
	cb GetNavigationHistoryCB
}

func NewAsyncGetNavigationHistoryCommand(cb GetNavigationHistoryCB) *AsyncGetNavigationHistoryCommand {
	return &AsyncGetNavigationHistoryCommand{
		cb: cb,
	}
}

func (cmd *AsyncGetNavigationHistoryCommand) Name() string {
	return "Page.getNavigationHistory"
}

func (cmd *AsyncGetNavigationHistoryCommand) Params() interface{} {
	return nil
}

func (cmd *GetNavigationHistoryCommand) Result() *GetNavigationHistoryResult {
	return &cmd.result
}

func (cmd *GetNavigationHistoryCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetNavigationHistoryCommand) Done(data []byte, err error) {
	var result GetNavigationHistoryResult
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

type NavigateToHistoryEntryParams struct {
	EntryId int `json:"entryId"` // Unique id of the entry to navigate to.
}

// Navigates current page to the given history entry.
// @experimental
type NavigateToHistoryEntryCommand struct {
	params *NavigateToHistoryEntryParams
	wg     sync.WaitGroup
	err    error
}

func NewNavigateToHistoryEntryCommand(params *NavigateToHistoryEntryParams) *NavigateToHistoryEntryCommand {
	return &NavigateToHistoryEntryCommand{
		params: params,
	}
}

func (cmd *NavigateToHistoryEntryCommand) Name() string {
	return "Page.navigateToHistoryEntry"
}

func (cmd *NavigateToHistoryEntryCommand) Params() interface{} {
	return cmd.params
}

func (cmd *NavigateToHistoryEntryCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func NavigateToHistoryEntry(params *NavigateToHistoryEntryParams, conn *hc.Conn) (err error) {
	cmd := NewNavigateToHistoryEntryCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type NavigateToHistoryEntryCB func(err error)

// Navigates current page to the given history entry.
// @experimental
type AsyncNavigateToHistoryEntryCommand struct {
	params *NavigateToHistoryEntryParams
	cb     NavigateToHistoryEntryCB
}

func NewAsyncNavigateToHistoryEntryCommand(params *NavigateToHistoryEntryParams, cb NavigateToHistoryEntryCB) *AsyncNavigateToHistoryEntryCommand {
	return &AsyncNavigateToHistoryEntryCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncNavigateToHistoryEntryCommand) Name() string {
	return "Page.navigateToHistoryEntry"
}

func (cmd *AsyncNavigateToHistoryEntryCommand) Params() interface{} {
	return cmd.params
}

func (cmd *NavigateToHistoryEntryCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncNavigateToHistoryEntryCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type PageGetCookiesResult struct {
	Cookies []*Cookie `json:"cookies"` // Array of cookie objects.
}

// Returns all browser cookies. Depending on the backend support, will return detailed cookie information in the cookies field.
// @experimental
type PageGetCookiesCommand struct {
	result PageGetCookiesResult
	wg     sync.WaitGroup
	err    error
}

func NewPageGetCookiesCommand() *PageGetCookiesCommand {
	return &PageGetCookiesCommand{}
}

func (cmd *PageGetCookiesCommand) Name() string {
	return "Page.getCookies"
}

func (cmd *PageGetCookiesCommand) Params() interface{} {
	return nil
}

func (cmd *PageGetCookiesCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func PageGetCookies(conn *hc.Conn) (result *PageGetCookiesResult, err error) {
	cmd := NewPageGetCookiesCommand()
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type PageGetCookiesCB func(result *PageGetCookiesResult, err error)

// Returns all browser cookies. Depending on the backend support, will return detailed cookie information in the cookies field.
// @experimental
type AsyncPageGetCookiesCommand struct {
	cb PageGetCookiesCB
}

func NewAsyncPageGetCookiesCommand(cb PageGetCookiesCB) *AsyncPageGetCookiesCommand {
	return &AsyncPageGetCookiesCommand{
		cb: cb,
	}
}

func (cmd *AsyncPageGetCookiesCommand) Name() string {
	return "Page.getCookies"
}

func (cmd *AsyncPageGetCookiesCommand) Params() interface{} {
	return nil
}

func (cmd *PageGetCookiesCommand) Result() *PageGetCookiesResult {
	return &cmd.result
}

func (cmd *PageGetCookiesCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncPageGetCookiesCommand) Done(data []byte, err error) {
	var result PageGetCookiesResult
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

type PageDeleteCookieParams struct {
	CookieName string `json:"cookieName"` // Name of the cookie to remove.
	Url        string `json:"url"`        // URL to match cooke domain and path.
}

// Deletes browser cookie with given name, domain and path.
// @experimental
type PageDeleteCookieCommand struct {
	params *PageDeleteCookieParams
	wg     sync.WaitGroup
	err    error
}

func NewPageDeleteCookieCommand(params *PageDeleteCookieParams) *PageDeleteCookieCommand {
	return &PageDeleteCookieCommand{
		params: params,
	}
}

func (cmd *PageDeleteCookieCommand) Name() string {
	return "Page.deleteCookie"
}

func (cmd *PageDeleteCookieCommand) Params() interface{} {
	return cmd.params
}

func (cmd *PageDeleteCookieCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func PageDeleteCookie(params *PageDeleteCookieParams, conn *hc.Conn) (err error) {
	cmd := NewPageDeleteCookieCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type PageDeleteCookieCB func(err error)

// Deletes browser cookie with given name, domain and path.
// @experimental
type AsyncPageDeleteCookieCommand struct {
	params *PageDeleteCookieParams
	cb     PageDeleteCookieCB
}

func NewAsyncPageDeleteCookieCommand(params *PageDeleteCookieParams, cb PageDeleteCookieCB) *AsyncPageDeleteCookieCommand {
	return &AsyncPageDeleteCookieCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncPageDeleteCookieCommand) Name() string {
	return "Page.deleteCookie"
}

func (cmd *AsyncPageDeleteCookieCommand) Params() interface{} {
	return cmd.params
}

func (cmd *PageDeleteCookieCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncPageDeleteCookieCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type GetResourceTreeResult struct {
	FrameTree *FrameResourceTree `json:"frameTree"` // Present frame / resource tree structure.
}

// Returns present frame / resource tree structure.
// @experimental
type GetResourceTreeCommand struct {
	result GetResourceTreeResult
	wg     sync.WaitGroup
	err    error
}

func NewGetResourceTreeCommand() *GetResourceTreeCommand {
	return &GetResourceTreeCommand{}
}

func (cmd *GetResourceTreeCommand) Name() string {
	return "Page.getResourceTree"
}

func (cmd *GetResourceTreeCommand) Params() interface{} {
	return nil
}

func (cmd *GetResourceTreeCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetResourceTree(conn *hc.Conn) (result *GetResourceTreeResult, err error) {
	cmd := NewGetResourceTreeCommand()
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetResourceTreeCB func(result *GetResourceTreeResult, err error)

// Returns present frame / resource tree structure.
// @experimental
type AsyncGetResourceTreeCommand struct {
	cb GetResourceTreeCB
}

func NewAsyncGetResourceTreeCommand(cb GetResourceTreeCB) *AsyncGetResourceTreeCommand {
	return &AsyncGetResourceTreeCommand{
		cb: cb,
	}
}

func (cmd *AsyncGetResourceTreeCommand) Name() string {
	return "Page.getResourceTree"
}

func (cmd *AsyncGetResourceTreeCommand) Params() interface{} {
	return nil
}

func (cmd *GetResourceTreeCommand) Result() *GetResourceTreeResult {
	return &cmd.result
}

func (cmd *GetResourceTreeCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetResourceTreeCommand) Done(data []byte, err error) {
	var result GetResourceTreeResult
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

type GetResourceContentParams struct {
	FrameId FrameId `json:"frameId"` // Frame id to get resource for.
	Url     string  `json:"url"`     // URL of the resource to get content for.
}

type GetResourceContentResult struct {
	Content       string `json:"content"`       // Resource content.
	Base64Encoded bool   `json:"base64Encoded"` // True, if content was served as base64.
}

// Returns content of the given resource.
// @experimental
type GetResourceContentCommand struct {
	params *GetResourceContentParams
	result GetResourceContentResult
	wg     sync.WaitGroup
	err    error
}

func NewGetResourceContentCommand(params *GetResourceContentParams) *GetResourceContentCommand {
	return &GetResourceContentCommand{
		params: params,
	}
}

func (cmd *GetResourceContentCommand) Name() string {
	return "Page.getResourceContent"
}

func (cmd *GetResourceContentCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetResourceContentCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetResourceContent(params *GetResourceContentParams, conn *hc.Conn) (result *GetResourceContentResult, err error) {
	cmd := NewGetResourceContentCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetResourceContentCB func(result *GetResourceContentResult, err error)

// Returns content of the given resource.
// @experimental
type AsyncGetResourceContentCommand struct {
	params *GetResourceContentParams
	cb     GetResourceContentCB
}

func NewAsyncGetResourceContentCommand(params *GetResourceContentParams, cb GetResourceContentCB) *AsyncGetResourceContentCommand {
	return &AsyncGetResourceContentCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetResourceContentCommand) Name() string {
	return "Page.getResourceContent"
}

func (cmd *AsyncGetResourceContentCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetResourceContentCommand) Result() *GetResourceContentResult {
	return &cmd.result
}

func (cmd *GetResourceContentCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetResourceContentCommand) Done(data []byte, err error) {
	var result GetResourceContentResult
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

type SearchInResourceParams struct {
	FrameId       FrameId `json:"frameId"`                 // Frame id for resource to search in.
	Url           string  `json:"url"`                     // URL of the resource to search in.
	Query         string  `json:"query"`                   // String to search for.
	CaseSensitive bool    `json:"caseSensitive,omitempty"` // If true, search is case sensitive.
	IsRegex       bool    `json:"isRegex,omitempty"`       // If true, treats string parameter as regex.
}

type SearchInResourceResult struct {
	Result []*SearchMatch `json:"result"` // List of search matches.
}

// Searches for given string in resource content.
// @experimental
type SearchInResourceCommand struct {
	params *SearchInResourceParams
	result SearchInResourceResult
	wg     sync.WaitGroup
	err    error
}

func NewSearchInResourceCommand(params *SearchInResourceParams) *SearchInResourceCommand {
	return &SearchInResourceCommand{
		params: params,
	}
}

func (cmd *SearchInResourceCommand) Name() string {
	return "Page.searchInResource"
}

func (cmd *SearchInResourceCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SearchInResourceCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SearchInResource(params *SearchInResourceParams, conn *hc.Conn) (result *SearchInResourceResult, err error) {
	cmd := NewSearchInResourceCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type SearchInResourceCB func(result *SearchInResourceResult, err error)

// Searches for given string in resource content.
// @experimental
type AsyncSearchInResourceCommand struct {
	params *SearchInResourceParams
	cb     SearchInResourceCB
}

func NewAsyncSearchInResourceCommand(params *SearchInResourceParams, cb SearchInResourceCB) *AsyncSearchInResourceCommand {
	return &AsyncSearchInResourceCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSearchInResourceCommand) Name() string {
	return "Page.searchInResource"
}

func (cmd *AsyncSearchInResourceCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SearchInResourceCommand) Result() *SearchInResourceResult {
	return &cmd.result
}

func (cmd *SearchInResourceCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSearchInResourceCommand) Done(data []byte, err error) {
	var result SearchInResourceResult
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

type SetDocumentContentParams struct {
	FrameId FrameId `json:"frameId"` // Frame id to set HTML for.
	Html    string  `json:"html"`    // HTML content to set.
}

// Sets given markup as the document's HTML.
// @experimental
type SetDocumentContentCommand struct {
	params *SetDocumentContentParams
	wg     sync.WaitGroup
	err    error
}

func NewSetDocumentContentCommand(params *SetDocumentContentParams) *SetDocumentContentCommand {
	return &SetDocumentContentCommand{
		params: params,
	}
}

func (cmd *SetDocumentContentCommand) Name() string {
	return "Page.setDocumentContent"
}

func (cmd *SetDocumentContentCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetDocumentContentCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetDocumentContent(params *SetDocumentContentParams, conn *hc.Conn) (err error) {
	cmd := NewSetDocumentContentCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetDocumentContentCB func(err error)

// Sets given markup as the document's HTML.
// @experimental
type AsyncSetDocumentContentCommand struct {
	params *SetDocumentContentParams
	cb     SetDocumentContentCB
}

func NewAsyncSetDocumentContentCommand(params *SetDocumentContentParams, cb SetDocumentContentCB) *AsyncSetDocumentContentCommand {
	return &AsyncSetDocumentContentCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetDocumentContentCommand) Name() string {
	return "Page.setDocumentContent"
}

func (cmd *AsyncSetDocumentContentCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetDocumentContentCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetDocumentContentCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type PageSetDeviceMetricsOverrideParams struct {
	Width             int                `json:"width"`                       // Overriding width value in pixels (minimum 0, maximum 10000000). 0 disables the override.
	Height            int                `json:"height"`                      // Overriding height value in pixels (minimum 0, maximum 10000000). 0 disables the override.
	DeviceScaleFactor float64            `json:"deviceScaleFactor"`           // Overriding device scale factor value. 0 disables the override.
	Mobile            bool               `json:"mobile"`                      // Whether to emulate mobile device. This includes viewport meta tag, overlay scrollbars, text autosizing and more.
	FitWindow         bool               `json:"fitWindow"`                   // Whether a view that exceeds the available browser window area should be scaled down to fit.
	Scale             float64            `json:"scale,omitempty"`             // Scale to apply to resulting view image. Ignored in |fitWindow| mode.
	OffsetX           float64            `json:"offsetX,omitempty"`           // X offset to shift resulting view image by. Ignored in |fitWindow| mode.
	OffsetY           float64            `json:"offsetY,omitempty"`           // Y offset to shift resulting view image by. Ignored in |fitWindow| mode.
	ScreenWidth       int                `json:"screenWidth,omitempty"`       // Overriding screen width value in pixels (minimum 0, maximum 10000000). Only used for |mobile==true|.
	ScreenHeight      int                `json:"screenHeight,omitempty"`      // Overriding screen height value in pixels (minimum 0, maximum 10000000). Only used for |mobile==true|.
	PositionX         int                `json:"positionX,omitempty"`         // Overriding view X position on screen in pixels (minimum 0, maximum 10000000). Only used for |mobile==true|.
	PositionY         int                `json:"positionY,omitempty"`         // Overriding view Y position on screen in pixels (minimum 0, maximum 10000000). Only used for |mobile==true|.
	ScreenOrientation *ScreenOrientation `json:"screenOrientation,omitempty"` // Screen orientation override.
}

// Overrides the values of device screen dimensions (window.screen.width, window.screen.height, window.innerWidth, window.innerHeight, and "device-width"/"device-height"-related CSS media query results).
// @experimental
type PageSetDeviceMetricsOverrideCommand struct {
	params *PageSetDeviceMetricsOverrideParams
	wg     sync.WaitGroup
	err    error
}

func NewPageSetDeviceMetricsOverrideCommand(params *PageSetDeviceMetricsOverrideParams) *PageSetDeviceMetricsOverrideCommand {
	return &PageSetDeviceMetricsOverrideCommand{
		params: params,
	}
}

func (cmd *PageSetDeviceMetricsOverrideCommand) Name() string {
	return "Page.setDeviceMetricsOverride"
}

func (cmd *PageSetDeviceMetricsOverrideCommand) Params() interface{} {
	return cmd.params
}

func (cmd *PageSetDeviceMetricsOverrideCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func PageSetDeviceMetricsOverride(params *PageSetDeviceMetricsOverrideParams, conn *hc.Conn) (err error) {
	cmd := NewPageSetDeviceMetricsOverrideCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type PageSetDeviceMetricsOverrideCB func(err error)

// Overrides the values of device screen dimensions (window.screen.width, window.screen.height, window.innerWidth, window.innerHeight, and "device-width"/"device-height"-related CSS media query results).
// @experimental
type AsyncPageSetDeviceMetricsOverrideCommand struct {
	params *PageSetDeviceMetricsOverrideParams
	cb     PageSetDeviceMetricsOverrideCB
}

func NewAsyncPageSetDeviceMetricsOverrideCommand(params *PageSetDeviceMetricsOverrideParams, cb PageSetDeviceMetricsOverrideCB) *AsyncPageSetDeviceMetricsOverrideCommand {
	return &AsyncPageSetDeviceMetricsOverrideCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncPageSetDeviceMetricsOverrideCommand) Name() string {
	return "Page.setDeviceMetricsOverride"
}

func (cmd *AsyncPageSetDeviceMetricsOverrideCommand) Params() interface{} {
	return cmd.params
}

func (cmd *PageSetDeviceMetricsOverrideCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncPageSetDeviceMetricsOverrideCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Clears the overriden device metrics.
// @experimental
type PageClearDeviceMetricsOverrideCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewPageClearDeviceMetricsOverrideCommand() *PageClearDeviceMetricsOverrideCommand {
	return &PageClearDeviceMetricsOverrideCommand{}
}

func (cmd *PageClearDeviceMetricsOverrideCommand) Name() string {
	return "Page.clearDeviceMetricsOverride"
}

func (cmd *PageClearDeviceMetricsOverrideCommand) Params() interface{} {
	return nil
}

func (cmd *PageClearDeviceMetricsOverrideCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func PageClearDeviceMetricsOverride(conn *hc.Conn) (err error) {
	cmd := NewPageClearDeviceMetricsOverrideCommand()
	cmd.Run(conn)
	return cmd.err
}

type PageClearDeviceMetricsOverrideCB func(err error)

// Clears the overriden device metrics.
// @experimental
type AsyncPageClearDeviceMetricsOverrideCommand struct {
	cb PageClearDeviceMetricsOverrideCB
}

func NewAsyncPageClearDeviceMetricsOverrideCommand(cb PageClearDeviceMetricsOverrideCB) *AsyncPageClearDeviceMetricsOverrideCommand {
	return &AsyncPageClearDeviceMetricsOverrideCommand{
		cb: cb,
	}
}

func (cmd *AsyncPageClearDeviceMetricsOverrideCommand) Name() string {
	return "Page.clearDeviceMetricsOverride"
}

func (cmd *AsyncPageClearDeviceMetricsOverrideCommand) Params() interface{} {
	return nil
}

func (cmd *PageClearDeviceMetricsOverrideCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncPageClearDeviceMetricsOverrideCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type PageSetGeolocationOverrideParams struct {
	Latitude  float64 `json:"latitude,omitempty"`  // Mock latitude
	Longitude float64 `json:"longitude,omitempty"` // Mock longitude
	Accuracy  float64 `json:"accuracy,omitempty"`  // Mock accuracy
}

// Overrides the Geolocation Position or Error. Omitting any of the parameters emulates position unavailable.

type PageSetGeolocationOverrideCommand struct {
	params *PageSetGeolocationOverrideParams
	wg     sync.WaitGroup
	err    error
}

func NewPageSetGeolocationOverrideCommand(params *PageSetGeolocationOverrideParams) *PageSetGeolocationOverrideCommand {
	return &PageSetGeolocationOverrideCommand{
		params: params,
	}
}

func (cmd *PageSetGeolocationOverrideCommand) Name() string {
	return "Page.setGeolocationOverride"
}

func (cmd *PageSetGeolocationOverrideCommand) Params() interface{} {
	return cmd.params
}

func (cmd *PageSetGeolocationOverrideCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func PageSetGeolocationOverride(params *PageSetGeolocationOverrideParams, conn *hc.Conn) (err error) {
	cmd := NewPageSetGeolocationOverrideCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type PageSetGeolocationOverrideCB func(err error)

// Overrides the Geolocation Position or Error. Omitting any of the parameters emulates position unavailable.

type AsyncPageSetGeolocationOverrideCommand struct {
	params *PageSetGeolocationOverrideParams
	cb     PageSetGeolocationOverrideCB
}

func NewAsyncPageSetGeolocationOverrideCommand(params *PageSetGeolocationOverrideParams, cb PageSetGeolocationOverrideCB) *AsyncPageSetGeolocationOverrideCommand {
	return &AsyncPageSetGeolocationOverrideCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncPageSetGeolocationOverrideCommand) Name() string {
	return "Page.setGeolocationOverride"
}

func (cmd *AsyncPageSetGeolocationOverrideCommand) Params() interface{} {
	return cmd.params
}

func (cmd *PageSetGeolocationOverrideCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncPageSetGeolocationOverrideCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Clears the overriden Geolocation Position and Error.

type PageClearGeolocationOverrideCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewPageClearGeolocationOverrideCommand() *PageClearGeolocationOverrideCommand {
	return &PageClearGeolocationOverrideCommand{}
}

func (cmd *PageClearGeolocationOverrideCommand) Name() string {
	return "Page.clearGeolocationOverride"
}

func (cmd *PageClearGeolocationOverrideCommand) Params() interface{} {
	return nil
}

func (cmd *PageClearGeolocationOverrideCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func PageClearGeolocationOverride(conn *hc.Conn) (err error) {
	cmd := NewPageClearGeolocationOverrideCommand()
	cmd.Run(conn)
	return cmd.err
}

type PageClearGeolocationOverrideCB func(err error)

// Clears the overriden Geolocation Position and Error.

type AsyncPageClearGeolocationOverrideCommand struct {
	cb PageClearGeolocationOverrideCB
}

func NewAsyncPageClearGeolocationOverrideCommand(cb PageClearGeolocationOverrideCB) *AsyncPageClearGeolocationOverrideCommand {
	return &AsyncPageClearGeolocationOverrideCommand{
		cb: cb,
	}
}

func (cmd *AsyncPageClearGeolocationOverrideCommand) Name() string {
	return "Page.clearGeolocationOverride"
}

func (cmd *AsyncPageClearGeolocationOverrideCommand) Params() interface{} {
	return nil
}

func (cmd *PageClearGeolocationOverrideCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncPageClearGeolocationOverrideCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type PageSetDeviceOrientationOverrideParams struct {
	Alpha float64 `json:"alpha"` // Mock alpha
	Beta  float64 `json:"beta"`  // Mock beta
	Gamma float64 `json:"gamma"` // Mock gamma
}

// Overrides the Device Orientation.
// @experimental
type PageSetDeviceOrientationOverrideCommand struct {
	params *PageSetDeviceOrientationOverrideParams
	wg     sync.WaitGroup
	err    error
}

func NewPageSetDeviceOrientationOverrideCommand(params *PageSetDeviceOrientationOverrideParams) *PageSetDeviceOrientationOverrideCommand {
	return &PageSetDeviceOrientationOverrideCommand{
		params: params,
	}
}

func (cmd *PageSetDeviceOrientationOverrideCommand) Name() string {
	return "Page.setDeviceOrientationOverride"
}

func (cmd *PageSetDeviceOrientationOverrideCommand) Params() interface{} {
	return cmd.params
}

func (cmd *PageSetDeviceOrientationOverrideCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func PageSetDeviceOrientationOverride(params *PageSetDeviceOrientationOverrideParams, conn *hc.Conn) (err error) {
	cmd := NewPageSetDeviceOrientationOverrideCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type PageSetDeviceOrientationOverrideCB func(err error)

// Overrides the Device Orientation.
// @experimental
type AsyncPageSetDeviceOrientationOverrideCommand struct {
	params *PageSetDeviceOrientationOverrideParams
	cb     PageSetDeviceOrientationOverrideCB
}

func NewAsyncPageSetDeviceOrientationOverrideCommand(params *PageSetDeviceOrientationOverrideParams, cb PageSetDeviceOrientationOverrideCB) *AsyncPageSetDeviceOrientationOverrideCommand {
	return &AsyncPageSetDeviceOrientationOverrideCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncPageSetDeviceOrientationOverrideCommand) Name() string {
	return "Page.setDeviceOrientationOverride"
}

func (cmd *AsyncPageSetDeviceOrientationOverrideCommand) Params() interface{} {
	return cmd.params
}

func (cmd *PageSetDeviceOrientationOverrideCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncPageSetDeviceOrientationOverrideCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Clears the overridden Device Orientation.
// @experimental
type PageClearDeviceOrientationOverrideCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewPageClearDeviceOrientationOverrideCommand() *PageClearDeviceOrientationOverrideCommand {
	return &PageClearDeviceOrientationOverrideCommand{}
}

func (cmd *PageClearDeviceOrientationOverrideCommand) Name() string {
	return "Page.clearDeviceOrientationOverride"
}

func (cmd *PageClearDeviceOrientationOverrideCommand) Params() interface{} {
	return nil
}

func (cmd *PageClearDeviceOrientationOverrideCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func PageClearDeviceOrientationOverride(conn *hc.Conn) (err error) {
	cmd := NewPageClearDeviceOrientationOverrideCommand()
	cmd.Run(conn)
	return cmd.err
}

type PageClearDeviceOrientationOverrideCB func(err error)

// Clears the overridden Device Orientation.
// @experimental
type AsyncPageClearDeviceOrientationOverrideCommand struct {
	cb PageClearDeviceOrientationOverrideCB
}

func NewAsyncPageClearDeviceOrientationOverrideCommand(cb PageClearDeviceOrientationOverrideCB) *AsyncPageClearDeviceOrientationOverrideCommand {
	return &AsyncPageClearDeviceOrientationOverrideCommand{
		cb: cb,
	}
}

func (cmd *AsyncPageClearDeviceOrientationOverrideCommand) Name() string {
	return "Page.clearDeviceOrientationOverride"
}

func (cmd *AsyncPageClearDeviceOrientationOverrideCommand) Params() interface{} {
	return nil
}

func (cmd *PageClearDeviceOrientationOverrideCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncPageClearDeviceOrientationOverrideCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type PageSetTouchEmulationEnabledParams struct {
	Enabled       bool   `json:"enabled"`                 // Whether the touch event emulation should be enabled.
	Configuration string `json:"configuration,omitempty"` // Touch/gesture events configuration. Default: current platform.
}

// Toggles mouse event-based touch event emulation.
// @experimental
type PageSetTouchEmulationEnabledCommand struct {
	params *PageSetTouchEmulationEnabledParams
	wg     sync.WaitGroup
	err    error
}

func NewPageSetTouchEmulationEnabledCommand(params *PageSetTouchEmulationEnabledParams) *PageSetTouchEmulationEnabledCommand {
	return &PageSetTouchEmulationEnabledCommand{
		params: params,
	}
}

func (cmd *PageSetTouchEmulationEnabledCommand) Name() string {
	return "Page.setTouchEmulationEnabled"
}

func (cmd *PageSetTouchEmulationEnabledCommand) Params() interface{} {
	return cmd.params
}

func (cmd *PageSetTouchEmulationEnabledCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func PageSetTouchEmulationEnabled(params *PageSetTouchEmulationEnabledParams, conn *hc.Conn) (err error) {
	cmd := NewPageSetTouchEmulationEnabledCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type PageSetTouchEmulationEnabledCB func(err error)

// Toggles mouse event-based touch event emulation.
// @experimental
type AsyncPageSetTouchEmulationEnabledCommand struct {
	params *PageSetTouchEmulationEnabledParams
	cb     PageSetTouchEmulationEnabledCB
}

func NewAsyncPageSetTouchEmulationEnabledCommand(params *PageSetTouchEmulationEnabledParams, cb PageSetTouchEmulationEnabledCB) *AsyncPageSetTouchEmulationEnabledCommand {
	return &AsyncPageSetTouchEmulationEnabledCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncPageSetTouchEmulationEnabledCommand) Name() string {
	return "Page.setTouchEmulationEnabled"
}

func (cmd *AsyncPageSetTouchEmulationEnabledCommand) Params() interface{} {
	return cmd.params
}

func (cmd *PageSetTouchEmulationEnabledCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncPageSetTouchEmulationEnabledCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type CaptureScreenshotResult struct {
	Data string `json:"data"` // Base64-encoded image data (PNG).
}

// Capture page screenshot.
// @experimental
type CaptureScreenshotCommand struct {
	result CaptureScreenshotResult
	wg     sync.WaitGroup
	err    error
}

func NewCaptureScreenshotCommand() *CaptureScreenshotCommand {
	return &CaptureScreenshotCommand{}
}

func (cmd *CaptureScreenshotCommand) Name() string {
	return "Page.captureScreenshot"
}

func (cmd *CaptureScreenshotCommand) Params() interface{} {
	return nil
}

func (cmd *CaptureScreenshotCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func CaptureScreenshot(conn *hc.Conn) (result *CaptureScreenshotResult, err error) {
	cmd := NewCaptureScreenshotCommand()
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type CaptureScreenshotCB func(result *CaptureScreenshotResult, err error)

// Capture page screenshot.
// @experimental
type AsyncCaptureScreenshotCommand struct {
	cb CaptureScreenshotCB
}

func NewAsyncCaptureScreenshotCommand(cb CaptureScreenshotCB) *AsyncCaptureScreenshotCommand {
	return &AsyncCaptureScreenshotCommand{
		cb: cb,
	}
}

func (cmd *AsyncCaptureScreenshotCommand) Name() string {
	return "Page.captureScreenshot"
}

func (cmd *AsyncCaptureScreenshotCommand) Params() interface{} {
	return nil
}

func (cmd *CaptureScreenshotCommand) Result() *CaptureScreenshotResult {
	return &cmd.result
}

func (cmd *CaptureScreenshotCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncCaptureScreenshotCommand) Done(data []byte, err error) {
	var result CaptureScreenshotResult
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

type StartScreencastParams struct {
	Format        string `json:"format,omitempty"`        // Image compression format.
	Quality       int    `json:"quality,omitempty"`       // Compression quality from range [0..100].
	MaxWidth      int    `json:"maxWidth,omitempty"`      // Maximum screenshot width.
	MaxHeight     int    `json:"maxHeight,omitempty"`     // Maximum screenshot height.
	EveryNthFrame int    `json:"everyNthFrame,omitempty"` // Send every n-th frame.
}

// Starts sending each frame using the screencastFrame event.
// @experimental
type StartScreencastCommand struct {
	params *StartScreencastParams
	wg     sync.WaitGroup
	err    error
}

func NewStartScreencastCommand(params *StartScreencastParams) *StartScreencastCommand {
	return &StartScreencastCommand{
		params: params,
	}
}

func (cmd *StartScreencastCommand) Name() string {
	return "Page.startScreencast"
}

func (cmd *StartScreencastCommand) Params() interface{} {
	return cmd.params
}

func (cmd *StartScreencastCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func StartScreencast(params *StartScreencastParams, conn *hc.Conn) (err error) {
	cmd := NewStartScreencastCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type StartScreencastCB func(err error)

// Starts sending each frame using the screencastFrame event.
// @experimental
type AsyncStartScreencastCommand struct {
	params *StartScreencastParams
	cb     StartScreencastCB
}

func NewAsyncStartScreencastCommand(params *StartScreencastParams, cb StartScreencastCB) *AsyncStartScreencastCommand {
	return &AsyncStartScreencastCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncStartScreencastCommand) Name() string {
	return "Page.startScreencast"
}

func (cmd *AsyncStartScreencastCommand) Params() interface{} {
	return cmd.params
}

func (cmd *StartScreencastCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncStartScreencastCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Stops sending each frame in the screencastFrame.
// @experimental
type StopScreencastCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewStopScreencastCommand() *StopScreencastCommand {
	return &StopScreencastCommand{}
}

func (cmd *StopScreencastCommand) Name() string {
	return "Page.stopScreencast"
}

func (cmd *StopScreencastCommand) Params() interface{} {
	return nil
}

func (cmd *StopScreencastCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func StopScreencast(conn *hc.Conn) (err error) {
	cmd := NewStopScreencastCommand()
	cmd.Run(conn)
	return cmd.err
}

type StopScreencastCB func(err error)

// Stops sending each frame in the screencastFrame.
// @experimental
type AsyncStopScreencastCommand struct {
	cb StopScreencastCB
}

func NewAsyncStopScreencastCommand(cb StopScreencastCB) *AsyncStopScreencastCommand {
	return &AsyncStopScreencastCommand{
		cb: cb,
	}
}

func (cmd *AsyncStopScreencastCommand) Name() string {
	return "Page.stopScreencast"
}

func (cmd *AsyncStopScreencastCommand) Params() interface{} {
	return nil
}

func (cmd *StopScreencastCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncStopScreencastCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type ScreencastFrameAckParams struct {
	SessionId int `json:"sessionId"` // Frame number.
}

// Acknowledges that a screencast frame has been received by the frontend.
// @experimental
type ScreencastFrameAckCommand struct {
	params *ScreencastFrameAckParams
	wg     sync.WaitGroup
	err    error
}

func NewScreencastFrameAckCommand(params *ScreencastFrameAckParams) *ScreencastFrameAckCommand {
	return &ScreencastFrameAckCommand{
		params: params,
	}
}

func (cmd *ScreencastFrameAckCommand) Name() string {
	return "Page.screencastFrameAck"
}

func (cmd *ScreencastFrameAckCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ScreencastFrameAckCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ScreencastFrameAck(params *ScreencastFrameAckParams, conn *hc.Conn) (err error) {
	cmd := NewScreencastFrameAckCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type ScreencastFrameAckCB func(err error)

// Acknowledges that a screencast frame has been received by the frontend.
// @experimental
type AsyncScreencastFrameAckCommand struct {
	params *ScreencastFrameAckParams
	cb     ScreencastFrameAckCB
}

func NewAsyncScreencastFrameAckCommand(params *ScreencastFrameAckParams, cb ScreencastFrameAckCB) *AsyncScreencastFrameAckCommand {
	return &AsyncScreencastFrameAckCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncScreencastFrameAckCommand) Name() string {
	return "Page.screencastFrameAck"
}

func (cmd *AsyncScreencastFrameAckCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ScreencastFrameAckCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncScreencastFrameAckCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type HandleJavaScriptDialogParams struct {
	Accept     bool   `json:"accept"`               // Whether to accept or dismiss the dialog.
	PromptText string `json:"promptText,omitempty"` // The text to enter into the dialog prompt before accepting. Used only if this is a prompt dialog.
}

// Accepts or dismisses a JavaScript initiated dialog (alert, confirm, prompt, or onbeforeunload).

type HandleJavaScriptDialogCommand struct {
	params *HandleJavaScriptDialogParams
	wg     sync.WaitGroup
	err    error
}

func NewHandleJavaScriptDialogCommand(params *HandleJavaScriptDialogParams) *HandleJavaScriptDialogCommand {
	return &HandleJavaScriptDialogCommand{
		params: params,
	}
}

func (cmd *HandleJavaScriptDialogCommand) Name() string {
	return "Page.handleJavaScriptDialog"
}

func (cmd *HandleJavaScriptDialogCommand) Params() interface{} {
	return cmd.params
}

func (cmd *HandleJavaScriptDialogCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func HandleJavaScriptDialog(params *HandleJavaScriptDialogParams, conn *hc.Conn) (err error) {
	cmd := NewHandleJavaScriptDialogCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type HandleJavaScriptDialogCB func(err error)

// Accepts or dismisses a JavaScript initiated dialog (alert, confirm, prompt, or onbeforeunload).

type AsyncHandleJavaScriptDialogCommand struct {
	params *HandleJavaScriptDialogParams
	cb     HandleJavaScriptDialogCB
}

func NewAsyncHandleJavaScriptDialogCommand(params *HandleJavaScriptDialogParams, cb HandleJavaScriptDialogCB) *AsyncHandleJavaScriptDialogCommand {
	return &AsyncHandleJavaScriptDialogCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncHandleJavaScriptDialogCommand) Name() string {
	return "Page.handleJavaScriptDialog"
}

func (cmd *AsyncHandleJavaScriptDialogCommand) Params() interface{} {
	return cmd.params
}

func (cmd *HandleJavaScriptDialogCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncHandleJavaScriptDialogCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetColorPickerEnabledParams struct {
	Enabled bool `json:"enabled"` // Shows / hides color picker
}

// Shows / hides color picker
// @experimental
type SetColorPickerEnabledCommand struct {
	params *SetColorPickerEnabledParams
	wg     sync.WaitGroup
	err    error
}

func NewSetColorPickerEnabledCommand(params *SetColorPickerEnabledParams) *SetColorPickerEnabledCommand {
	return &SetColorPickerEnabledCommand{
		params: params,
	}
}

func (cmd *SetColorPickerEnabledCommand) Name() string {
	return "Page.setColorPickerEnabled"
}

func (cmd *SetColorPickerEnabledCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetColorPickerEnabledCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetColorPickerEnabled(params *SetColorPickerEnabledParams, conn *hc.Conn) (err error) {
	cmd := NewSetColorPickerEnabledCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetColorPickerEnabledCB func(err error)

// Shows / hides color picker
// @experimental
type AsyncSetColorPickerEnabledCommand struct {
	params *SetColorPickerEnabledParams
	cb     SetColorPickerEnabledCB
}

func NewAsyncSetColorPickerEnabledCommand(params *SetColorPickerEnabledParams, cb SetColorPickerEnabledCB) *AsyncSetColorPickerEnabledCommand {
	return &AsyncSetColorPickerEnabledCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetColorPickerEnabledCommand) Name() string {
	return "Page.setColorPickerEnabled"
}

func (cmd *AsyncSetColorPickerEnabledCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetColorPickerEnabledCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetColorPickerEnabledCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type ConfigureOverlayParams struct {
	Suspended bool   `json:"suspended,omitempty"` // Whether overlay should be suspended and not consume any resources.
	Message   string `json:"message,omitempty"`   // Overlay message to display.
}

// Configures overlay.
// @experimental
type ConfigureOverlayCommand struct {
	params *ConfigureOverlayParams
	wg     sync.WaitGroup
	err    error
}

func NewConfigureOverlayCommand(params *ConfigureOverlayParams) *ConfigureOverlayCommand {
	return &ConfigureOverlayCommand{
		params: params,
	}
}

func (cmd *ConfigureOverlayCommand) Name() string {
	return "Page.configureOverlay"
}

func (cmd *ConfigureOverlayCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ConfigureOverlayCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ConfigureOverlay(params *ConfigureOverlayParams, conn *hc.Conn) (err error) {
	cmd := NewConfigureOverlayCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type ConfigureOverlayCB func(err error)

// Configures overlay.
// @experimental
type AsyncConfigureOverlayCommand struct {
	params *ConfigureOverlayParams
	cb     ConfigureOverlayCB
}

func NewAsyncConfigureOverlayCommand(params *ConfigureOverlayParams, cb ConfigureOverlayCB) *AsyncConfigureOverlayCommand {
	return &AsyncConfigureOverlayCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncConfigureOverlayCommand) Name() string {
	return "Page.configureOverlay"
}

func (cmd *AsyncConfigureOverlayCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ConfigureOverlayCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncConfigureOverlayCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type GetAppManifestResult struct {
	Url    string              `json:"url"` // Manifest location.
	Errors []*AppManifestError `json:"errors"`
	Data   string              `json:"data"` // Manifest content.
}

// @experimental
type GetAppManifestCommand struct {
	result GetAppManifestResult
	wg     sync.WaitGroup
	err    error
}

func NewGetAppManifestCommand() *GetAppManifestCommand {
	return &GetAppManifestCommand{}
}

func (cmd *GetAppManifestCommand) Name() string {
	return "Page.getAppManifest"
}

func (cmd *GetAppManifestCommand) Params() interface{} {
	return nil
}

func (cmd *GetAppManifestCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetAppManifest(conn *hc.Conn) (result *GetAppManifestResult, err error) {
	cmd := NewGetAppManifestCommand()
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetAppManifestCB func(result *GetAppManifestResult, err error)

// @experimental
type AsyncGetAppManifestCommand struct {
	cb GetAppManifestCB
}

func NewAsyncGetAppManifestCommand(cb GetAppManifestCB) *AsyncGetAppManifestCommand {
	return &AsyncGetAppManifestCommand{
		cb: cb,
	}
}

func (cmd *AsyncGetAppManifestCommand) Name() string {
	return "Page.getAppManifest"
}

func (cmd *AsyncGetAppManifestCommand) Params() interface{} {
	return nil
}

func (cmd *GetAppManifestCommand) Result() *GetAppManifestResult {
	return &cmd.result
}

func (cmd *GetAppManifestCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetAppManifestCommand) Done(data []byte, err error) {
	var result GetAppManifestResult
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

// @experimental
type RequestAppBannerCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewRequestAppBannerCommand() *RequestAppBannerCommand {
	return &RequestAppBannerCommand{}
}

func (cmd *RequestAppBannerCommand) Name() string {
	return "Page.requestAppBanner"
}

func (cmd *RequestAppBannerCommand) Params() interface{} {
	return nil
}

func (cmd *RequestAppBannerCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func RequestAppBanner(conn *hc.Conn) (err error) {
	cmd := NewRequestAppBannerCommand()
	cmd.Run(conn)
	return cmd.err
}

type RequestAppBannerCB func(err error)

// @experimental
type AsyncRequestAppBannerCommand struct {
	cb RequestAppBannerCB
}

func NewAsyncRequestAppBannerCommand(cb RequestAppBannerCB) *AsyncRequestAppBannerCommand {
	return &AsyncRequestAppBannerCommand{
		cb: cb,
	}
}

func (cmd *AsyncRequestAppBannerCommand) Name() string {
	return "Page.requestAppBanner"
}

func (cmd *AsyncRequestAppBannerCommand) Params() interface{} {
	return nil
}

func (cmd *RequestAppBannerCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncRequestAppBannerCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetControlNavigationsParams struct {
	Enabled bool `json:"enabled"`
}

// Toggles navigation throttling which allows programatic control over navigation and redirect response.
// @experimental
type SetControlNavigationsCommand struct {
	params *SetControlNavigationsParams
	wg     sync.WaitGroup
	err    error
}

func NewSetControlNavigationsCommand(params *SetControlNavigationsParams) *SetControlNavigationsCommand {
	return &SetControlNavigationsCommand{
		params: params,
	}
}

func (cmd *SetControlNavigationsCommand) Name() string {
	return "Page.setControlNavigations"
}

func (cmd *SetControlNavigationsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetControlNavigationsCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetControlNavigations(params *SetControlNavigationsParams, conn *hc.Conn) (err error) {
	cmd := NewSetControlNavigationsCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetControlNavigationsCB func(err error)

// Toggles navigation throttling which allows programatic control over navigation and redirect response.
// @experimental
type AsyncSetControlNavigationsCommand struct {
	params *SetControlNavigationsParams
	cb     SetControlNavigationsCB
}

func NewAsyncSetControlNavigationsCommand(params *SetControlNavigationsParams, cb SetControlNavigationsCB) *AsyncSetControlNavigationsCommand {
	return &AsyncSetControlNavigationsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetControlNavigationsCommand) Name() string {
	return "Page.setControlNavigations"
}

func (cmd *AsyncSetControlNavigationsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetControlNavigationsCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetControlNavigationsCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type ProcessNavigationParams struct {
	Response     NavigationResponse `json:"response"`
	NavigationId int                `json:"navigationId"`
}

// Should be sent in response to a navigationRequested or a redirectRequested event, telling the browser how to handle the navigation.
// @experimental
type ProcessNavigationCommand struct {
	params *ProcessNavigationParams
	wg     sync.WaitGroup
	err    error
}

func NewProcessNavigationCommand(params *ProcessNavigationParams) *ProcessNavigationCommand {
	return &ProcessNavigationCommand{
		params: params,
	}
}

func (cmd *ProcessNavigationCommand) Name() string {
	return "Page.processNavigation"
}

func (cmd *ProcessNavigationCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ProcessNavigationCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ProcessNavigation(params *ProcessNavigationParams, conn *hc.Conn) (err error) {
	cmd := NewProcessNavigationCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type ProcessNavigationCB func(err error)

// Should be sent in response to a navigationRequested or a redirectRequested event, telling the browser how to handle the navigation.
// @experimental
type AsyncProcessNavigationCommand struct {
	params *ProcessNavigationParams
	cb     ProcessNavigationCB
}

func NewAsyncProcessNavigationCommand(params *ProcessNavigationParams, cb ProcessNavigationCB) *AsyncProcessNavigationCommand {
	return &AsyncProcessNavigationCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncProcessNavigationCommand) Name() string {
	return "Page.processNavigation"
}

func (cmd *AsyncProcessNavigationCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ProcessNavigationCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncProcessNavigationCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type GetLayoutMetricsResult struct {
	LayoutViewport *LayoutViewport `json:"layoutViewport"` // Metrics relating to the layout viewport.
	VisualViewport *VisualViewport `json:"visualViewport"` // Metrics relating to the visual viewport.
}

// Returns metrics relating to the layouting of the page, such as viewport bounds/scale.
// @experimental
type GetLayoutMetricsCommand struct {
	result GetLayoutMetricsResult
	wg     sync.WaitGroup
	err    error
}

func NewGetLayoutMetricsCommand() *GetLayoutMetricsCommand {
	return &GetLayoutMetricsCommand{}
}

func (cmd *GetLayoutMetricsCommand) Name() string {
	return "Page.getLayoutMetrics"
}

func (cmd *GetLayoutMetricsCommand) Params() interface{} {
	return nil
}

func (cmd *GetLayoutMetricsCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetLayoutMetrics(conn *hc.Conn) (result *GetLayoutMetricsResult, err error) {
	cmd := NewGetLayoutMetricsCommand()
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetLayoutMetricsCB func(result *GetLayoutMetricsResult, err error)

// Returns metrics relating to the layouting of the page, such as viewport bounds/scale.
// @experimental
type AsyncGetLayoutMetricsCommand struct {
	cb GetLayoutMetricsCB
}

func NewAsyncGetLayoutMetricsCommand(cb GetLayoutMetricsCB) *AsyncGetLayoutMetricsCommand {
	return &AsyncGetLayoutMetricsCommand{
		cb: cb,
	}
}

func (cmd *AsyncGetLayoutMetricsCommand) Name() string {
	return "Page.getLayoutMetrics"
}

func (cmd *AsyncGetLayoutMetricsCommand) Params() interface{} {
	return nil
}

func (cmd *GetLayoutMetricsCommand) Result() *GetLayoutMetricsResult {
	return &cmd.result
}

func (cmd *GetLayoutMetricsCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetLayoutMetricsCommand) Done(data []byte, err error) {
	var result GetLayoutMetricsResult
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

type DomContentEventFiredEvent struct {
	Timestamp float64 `json:"timestamp"`
}

func OnDomContentEventFired(conn *hc.Conn, cb func(evt *DomContentEventFiredEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &DomContentEventFiredEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Page.domContentEventFired", sink)
}

type LoadEventFiredEvent struct {
	Timestamp float64 `json:"timestamp"`
}

func OnLoadEventFired(conn *hc.Conn, cb func(evt *LoadEventFiredEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &LoadEventFiredEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Page.loadEventFired", sink)
}

// Fired when frame has been attached to its parent.

type FrameAttachedEvent struct {
	FrameId       FrameId `json:"frameId"`       // Id of the frame that has been attached.
	ParentFrameId FrameId `json:"parentFrameId"` // Parent frame identifier.
}

func OnFrameAttached(conn *hc.Conn, cb func(evt *FrameAttachedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &FrameAttachedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Page.frameAttached", sink)
}

// Fired once navigation of the frame has completed. Frame is now associated with the new loader.

type FrameNavigatedEvent struct {
	Frame *Frame `json:"frame"` // Frame object.
}

func OnFrameNavigated(conn *hc.Conn, cb func(evt *FrameNavigatedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &FrameNavigatedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Page.frameNavigated", sink)
}

// Fired when frame has been detached from its parent.

type FrameDetachedEvent struct {
	FrameId FrameId `json:"frameId"` // Id of the frame that has been detached.
}

func OnFrameDetached(conn *hc.Conn, cb func(evt *FrameDetachedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &FrameDetachedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Page.frameDetached", sink)
}

// Fired when frame has started loading.
// @experimental
type FrameStartedLoadingEvent struct {
	FrameId FrameId `json:"frameId"` // Id of the frame that has started loading.
}

func OnFrameStartedLoading(conn *hc.Conn, cb func(evt *FrameStartedLoadingEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &FrameStartedLoadingEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Page.frameStartedLoading", sink)
}

// Fired when frame has stopped loading.
// @experimental
type FrameStoppedLoadingEvent struct {
	FrameId FrameId `json:"frameId"` // Id of the frame that has stopped loading.
}

func OnFrameStoppedLoading(conn *hc.Conn, cb func(evt *FrameStoppedLoadingEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &FrameStoppedLoadingEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Page.frameStoppedLoading", sink)
}

// Fired when frame schedules a potential navigation.
// @experimental
type FrameScheduledNavigationEvent struct {
	FrameId FrameId `json:"frameId"` // Id of the frame that has scheduled a navigation.
	Delay   float64 `json:"delay"`   // Delay (in seconds) until the navigation is scheduled to begin. The navigation is not guaranteed to start.
}

func OnFrameScheduledNavigation(conn *hc.Conn, cb func(evt *FrameScheduledNavigationEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &FrameScheduledNavigationEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Page.frameScheduledNavigation", sink)
}

// Fired when frame no longer has a scheduled navigation.
// @experimental
type FrameClearedScheduledNavigationEvent struct {
	FrameId FrameId `json:"frameId"` // Id of the frame that has cleared its scheduled navigation.
}

func OnFrameClearedScheduledNavigation(conn *hc.Conn, cb func(evt *FrameClearedScheduledNavigationEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &FrameClearedScheduledNavigationEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Page.frameClearedScheduledNavigation", sink)
}

// @experimental
type FrameResizedEvent struct {
}

func OnFrameResized(conn *hc.Conn, cb func(evt *FrameResizedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &FrameResizedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Page.frameResized", sink)
}

// Fired when a JavaScript initiated dialog (alert, confirm, prompt, or onbeforeunload) is about to open.

type JavascriptDialogOpeningEvent struct {
	Message string     `json:"message"` // Message that will be displayed by the dialog.
	Type    DialogType `json:"type"`    // Dialog type.
}

func OnJavascriptDialogOpening(conn *hc.Conn, cb func(evt *JavascriptDialogOpeningEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &JavascriptDialogOpeningEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Page.javascriptDialogOpening", sink)
}

// Fired when a JavaScript initiated dialog (alert, confirm, prompt, or onbeforeunload) has been closed.

type JavascriptDialogClosedEvent struct {
	Result bool `json:"result"` // Whether dialog was confirmed.
}

func OnJavascriptDialogClosed(conn *hc.Conn, cb func(evt *JavascriptDialogClosedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &JavascriptDialogClosedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Page.javascriptDialogClosed", sink)
}

// Compressed image data requested by the startScreencast.
// @experimental
type ScreencastFrameEvent struct {
	Data      string                   `json:"data"`      // Base64-encoded compressed image.
	Metadata  *ScreencastFrameMetadata `json:"metadata"`  // Screencast frame metadata.
	SessionId int                      `json:"sessionId"` // Frame number.
}

func OnScreencastFrame(conn *hc.Conn, cb func(evt *ScreencastFrameEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &ScreencastFrameEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Page.screencastFrame", sink)
}

// Fired when the page with currently enabled screencast was shown or hidden .
// @experimental
type ScreencastVisibilityChangedEvent struct {
	Visible bool `json:"visible"` // True if the page is visible.
}

func OnScreencastVisibilityChanged(conn *hc.Conn, cb func(evt *ScreencastVisibilityChangedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &ScreencastVisibilityChangedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Page.screencastVisibilityChanged", sink)
}

// Fired when a color has been picked.
// @experimental
type ColorPickedEvent struct {
	Color *RGBA `json:"color"` // RGBA of the picked color.
}

func OnColorPicked(conn *hc.Conn, cb func(evt *ColorPickedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &ColorPickedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Page.colorPicked", sink)
}

// Fired when interstitial page was shown

type InterstitialShownEvent struct {
}

func OnInterstitialShown(conn *hc.Conn, cb func(evt *InterstitialShownEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &InterstitialShownEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Page.interstitialShown", sink)
}

// Fired when interstitial page was hidden

type InterstitialHiddenEvent struct {
}

func OnInterstitialHidden(conn *hc.Conn, cb func(evt *InterstitialHiddenEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &InterstitialHiddenEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Page.interstitialHidden", sink)
}

// Fired when a navigation is started if navigation throttles are enabled.  The navigation will be deferred until processNavigation is called.

type NavigationRequestedEvent struct {
	IsInMainFrame bool   `json:"isInMainFrame"` // Whether the navigation is taking place in the main frame or in a subframe.
	IsRedirect    bool   `json:"isRedirect"`    // Whether the navigation has encountered a server redirect or not.
	NavigationId  int    `json:"navigationId"`
	Url           string `json:"url"` // URL of requested navigation.
}

func OnNavigationRequested(conn *hc.Conn, cb func(evt *NavigationRequestedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &NavigationRequestedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Page.navigationRequested", sink)
}
