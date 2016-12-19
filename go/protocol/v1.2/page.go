package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
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
	Id             string    `json:"id"`             // Frame unique identifier.
	ParentId       string    `json:"parentId"`       // Parent frame identifier.
	LoaderId       *LoaderId `json:"loaderId"`       // Identifier of the loader associated with this frame.
	Name           string    `json:"name"`           // Frame's name as specified in the tag.
	Url            string    `json:"url"`            // Frame document's URL.
	SecurityOrigin string    `json:"securityOrigin"` // Frame document's security origin.
	MimeType       string    `json:"mimeType"`       // Frame document's mimeType as determined by the browser.
}

// Information about the Resource on the page.
type FrameResource struct {
	Url      string       `json:"url"`      // Resource URL.
	Type     ResourceType `json:"type"`     // Type of this resource.
	MimeType string       `json:"mimeType"` // Resource mimeType as determined by the browser.
	Failed   bool         `json:"failed"`   // True if the resource failed to load.
	Canceled bool         `json:"canceled"` // True if the resource was canceled during loading.
}

// Information about the Frame hierarchy along with their cached resources.
type FrameResourceTree struct {
	Frame       *Frame               `json:"frame"`       // Frame information for this tree item.
	ChildFrames []*FrameResourceTree `json:"childFrames"` // Child frames.
	Resources   []*FrameResource     `json:"resources"`   // Information about frame resources.
}

// Unique script identifier.
type ScriptIdentifier string

// Navigation history entry.
type NavigationEntry struct {
	Id    int    `json:"id"`    // Unique id of the navigation history entry.
	Url   string `json:"url"`   // URL of the navigation history entry.
	Title string `json:"title"` // Title of the navigation history entry.
}

// Screencast frame metadata.
type ScreencastFrameMetadata struct {
	OffsetTop       int `json:"offsetTop"`       // Top offset in DIP.
	PageScaleFactor int `json:"pageScaleFactor"` // Page scale factor.
	DeviceWidth     int `json:"deviceWidth"`     // Device screen width in DIP.
	DeviceHeight    int `json:"deviceHeight"`    // Device screen height in DIP.
	ScrollOffsetX   int `json:"scrollOffsetX"`   // Position of horizontal scroll in CSS pixels.
	ScrollOffsetY   int `json:"scrollOffsetY"`   // Position of vertical scroll in CSS pixels.
	Timestamp       int `json:"timestamp"`       // Frame swap timestamp.
}

// Javascript dialog type.
type DialogType string

const DialogTypeAlert DialogType = "alert"
const DialogTypeConfirm DialogType = "confirm"
const DialogTypePrompt DialogType = "prompt"
const DialogTypeBeforeunload DialogType = "beforeunload"

// Error while paring app manifest.
type AppManifestError struct {
	Message  string `json:"message"`  // Error message.
	Critical int    `json:"critical"` // If criticial, this is a non-recoverable parse error.
	Line     int    `json:"line"`     // Error line.
	Column   int    `json:"column"`   // Error column.
}

// Proceed: allow the navigation; Cancel: cancel the navigation; CancelAndIgnore: cancels the navigation and makes the requester of the navigation acts like the request was never made.
type NavigationResponse string

const NavigationResponseProceed NavigationResponse = "Proceed"
const NavigationResponseCancel NavigationResponse = "Cancel"
const NavigationResponseCancelAndIgnore NavigationResponse = "CancelAndIgnore"

type PageEnableCB func(err error)

// Enables page domain notifications.
type PageEnableCommand struct {
	cb PageEnableCB
}

func NewPageEnableCommand(cb PageEnableCB) *PageEnableCommand {
	return &PageEnableCommand{
		cb: cb,
	}
}

func (cmd *PageEnableCommand) Name() string {
	return "Page.enable"
}

func (cmd *PageEnableCommand) Params() interface{} {
	return nil
}

func (cmd *PageEnableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type PageDisableCB func(err error)

// Disables page domain notifications.
type PageDisableCommand struct {
	cb PageDisableCB
}

func NewPageDisableCommand(cb PageDisableCB) *PageDisableCommand {
	return &PageDisableCommand{
		cb: cb,
	}
}

func (cmd *PageDisableCommand) Name() string {
	return "Page.disable"
}

func (cmd *PageDisableCommand) Params() interface{} {
	return nil
}

func (cmd *PageDisableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type AddScriptToEvaluateOnLoadParams struct {
	ScriptSource string `json:"scriptSource"`
}

type AddScriptToEvaluateOnLoadResult struct {
	Identifier ScriptIdentifier `json:"identifier"` // Identifier of the added script.
}

type AddScriptToEvaluateOnLoadCB func(result *AddScriptToEvaluateOnLoadResult, err error)

type AddScriptToEvaluateOnLoadCommand struct {
	params *AddScriptToEvaluateOnLoadParams
	cb     AddScriptToEvaluateOnLoadCB
}

func NewAddScriptToEvaluateOnLoadCommand(params *AddScriptToEvaluateOnLoadParams, cb AddScriptToEvaluateOnLoadCB) *AddScriptToEvaluateOnLoadCommand {
	return &AddScriptToEvaluateOnLoadCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AddScriptToEvaluateOnLoadCommand) Name() string {
	return "Page.addScriptToEvaluateOnLoad"
}

func (cmd *AddScriptToEvaluateOnLoadCommand) Params() interface{} {
	return cmd.params
}

func (cmd *AddScriptToEvaluateOnLoadCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj AddScriptToEvaluateOnLoadResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type RemoveScriptToEvaluateOnLoadParams struct {
	Identifier ScriptIdentifier `json:"identifier"`
}

type RemoveScriptToEvaluateOnLoadCB func(err error)

type RemoveScriptToEvaluateOnLoadCommand struct {
	params *RemoveScriptToEvaluateOnLoadParams
	cb     RemoveScriptToEvaluateOnLoadCB
}

func NewRemoveScriptToEvaluateOnLoadCommand(params *RemoveScriptToEvaluateOnLoadParams, cb RemoveScriptToEvaluateOnLoadCB) *RemoveScriptToEvaluateOnLoadCommand {
	return &RemoveScriptToEvaluateOnLoadCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *RemoveScriptToEvaluateOnLoadCommand) Name() string {
	return "Page.removeScriptToEvaluateOnLoad"
}

func (cmd *RemoveScriptToEvaluateOnLoadCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveScriptToEvaluateOnLoadCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetAutoAttachToCreatedPagesParams struct {
	AutoAttach bool `json:"autoAttach"` // If true, browser will open a new inspector window for every page created from this one.
}

type SetAutoAttachToCreatedPagesCB func(err error)

// Controls whether browser will open a new inspector window for connected pages.
type SetAutoAttachToCreatedPagesCommand struct {
	params *SetAutoAttachToCreatedPagesParams
	cb     SetAutoAttachToCreatedPagesCB
}

func NewSetAutoAttachToCreatedPagesCommand(params *SetAutoAttachToCreatedPagesParams, cb SetAutoAttachToCreatedPagesCB) *SetAutoAttachToCreatedPagesCommand {
	return &SetAutoAttachToCreatedPagesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetAutoAttachToCreatedPagesCommand) Name() string {
	return "Page.setAutoAttachToCreatedPages"
}

func (cmd *SetAutoAttachToCreatedPagesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetAutoAttachToCreatedPagesCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type ReloadParams struct {
	IgnoreCache            bool   `json:"ignoreCache"`            // If true, browser cache is ignored (as if the user pressed Shift+refresh).
	ScriptToEvaluateOnLoad string `json:"scriptToEvaluateOnLoad"` // If set, the script will be injected into all frames of the inspected page after reload.
}

type ReloadCB func(err error)

// Reloads given page optionally ignoring the cache.
type ReloadCommand struct {
	params *ReloadParams
	cb     ReloadCB
}

func NewReloadCommand(params *ReloadParams, cb ReloadCB) *ReloadCommand {
	return &ReloadCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *ReloadCommand) Name() string {
	return "Page.reload"
}

func (cmd *ReloadCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ReloadCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type NavigateParams struct {
	Url string `json:"url"` // URL to navigate the page to.
}

type NavigateResult struct {
	FrameId FrameId `json:"frameId"` // Frame id that will be navigated.
}

type NavigateCB func(result *NavigateResult, err error)

// Navigates current page to the given URL.
type NavigateCommand struct {
	params *NavigateParams
	cb     NavigateCB
}

func NewNavigateCommand(params *NavigateParams, cb NavigateCB) *NavigateCommand {
	return &NavigateCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *NavigateCommand) Name() string {
	return "Page.navigate"
}

func (cmd *NavigateCommand) Params() interface{} {
	return cmd.params
}

func (cmd *NavigateCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj NavigateResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type GetNavigationHistoryResult struct {
	CurrentIndex int                `json:"currentIndex"` // Index of the current navigation history entry.
	Entries      []*NavigationEntry `json:"entries"`      // Array of navigation history entries.
}

type GetNavigationHistoryCB func(result *GetNavigationHistoryResult, err error)

// Returns navigation history for the current page.
type GetNavigationHistoryCommand struct {
	cb GetNavigationHistoryCB
}

func NewGetNavigationHistoryCommand(cb GetNavigationHistoryCB) *GetNavigationHistoryCommand {
	return &GetNavigationHistoryCommand{
		cb: cb,
	}
}

func (cmd *GetNavigationHistoryCommand) Name() string {
	return "Page.getNavigationHistory"
}

func (cmd *GetNavigationHistoryCommand) Params() interface{} {
	return nil
}

func (cmd *GetNavigationHistoryCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetNavigationHistoryResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type NavigateToHistoryEntryParams struct {
	EntryId int `json:"entryId"` // Unique id of the entry to navigate to.
}

type NavigateToHistoryEntryCB func(err error)

// Navigates current page to the given history entry.
type NavigateToHistoryEntryCommand struct {
	params *NavigateToHistoryEntryParams
	cb     NavigateToHistoryEntryCB
}

func NewNavigateToHistoryEntryCommand(params *NavigateToHistoryEntryParams, cb NavigateToHistoryEntryCB) *NavigateToHistoryEntryCommand {
	return &NavigateToHistoryEntryCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *NavigateToHistoryEntryCommand) Name() string {
	return "Page.navigateToHistoryEntry"
}

func (cmd *NavigateToHistoryEntryCommand) Params() interface{} {
	return cmd.params
}

func (cmd *NavigateToHistoryEntryCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type PageGetCookiesResult struct {
	Cookies []*Cookie `json:"cookies"` // Array of cookie objects.
}

type PageGetCookiesCB func(result *PageGetCookiesResult, err error)

// Returns all browser cookies. Depending on the backend support, will return detailed cookie information in the cookies field.
type PageGetCookiesCommand struct {
	cb PageGetCookiesCB
}

func NewPageGetCookiesCommand(cb PageGetCookiesCB) *PageGetCookiesCommand {
	return &PageGetCookiesCommand{
		cb: cb,
	}
}

func (cmd *PageGetCookiesCommand) Name() string {
	return "Page.getCookies"
}

func (cmd *PageGetCookiesCommand) Params() interface{} {
	return nil
}

func (cmd *PageGetCookiesCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj PageGetCookiesResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type PageDeleteCookieParams struct {
	CookieName string `json:"cookieName"` // Name of the cookie to remove.
	Url        string `json:"url"`        // URL to match cooke domain and path.
}

type PageDeleteCookieCB func(err error)

// Deletes browser cookie with given name, domain and path.
type PageDeleteCookieCommand struct {
	params *PageDeleteCookieParams
	cb     PageDeleteCookieCB
}

func NewPageDeleteCookieCommand(params *PageDeleteCookieParams, cb PageDeleteCookieCB) *PageDeleteCookieCommand {
	return &PageDeleteCookieCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *PageDeleteCookieCommand) Name() string {
	return "Page.deleteCookie"
}

func (cmd *PageDeleteCookieCommand) Params() interface{} {
	return cmd.params
}

func (cmd *PageDeleteCookieCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type GetResourceTreeResult struct {
	FrameTree *FrameResourceTree `json:"frameTree"` // Present frame / resource tree structure.
}

type GetResourceTreeCB func(result *GetResourceTreeResult, err error)

// Returns present frame / resource tree structure.
type GetResourceTreeCommand struct {
	cb GetResourceTreeCB
}

func NewGetResourceTreeCommand(cb GetResourceTreeCB) *GetResourceTreeCommand {
	return &GetResourceTreeCommand{
		cb: cb,
	}
}

func (cmd *GetResourceTreeCommand) Name() string {
	return "Page.getResourceTree"
}

func (cmd *GetResourceTreeCommand) Params() interface{} {
	return nil
}

func (cmd *GetResourceTreeCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetResourceTreeResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
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

type GetResourceContentCB func(result *GetResourceContentResult, err error)

// Returns content of the given resource.
type GetResourceContentCommand struct {
	params *GetResourceContentParams
	cb     GetResourceContentCB
}

func NewGetResourceContentCommand(params *GetResourceContentParams, cb GetResourceContentCB) *GetResourceContentCommand {
	return &GetResourceContentCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetResourceContentCommand) Name() string {
	return "Page.getResourceContent"
}

func (cmd *GetResourceContentCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetResourceContentCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetResourceContentResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type SearchInResourceParams struct {
	FrameId       FrameId `json:"frameId"`       // Frame id for resource to search in.
	Url           string  `json:"url"`           // URL of the resource to search in.
	Query         string  `json:"query"`         // String to search for.
	CaseSensitive bool    `json:"caseSensitive"` // If true, search is case sensitive.
	IsRegex       bool    `json:"isRegex"`       // If true, treats string parameter as regex.
}

type SearchInResourceResult struct {
	Result []*SearchMatch `json:"result"` // List of search matches.
}

type SearchInResourceCB func(result *SearchInResourceResult, err error)

// Searches for given string in resource content.
type SearchInResourceCommand struct {
	params *SearchInResourceParams
	cb     SearchInResourceCB
}

func NewSearchInResourceCommand(params *SearchInResourceParams, cb SearchInResourceCB) *SearchInResourceCommand {
	return &SearchInResourceCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SearchInResourceCommand) Name() string {
	return "Page.searchInResource"
}

func (cmd *SearchInResourceCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SearchInResourceCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj SearchInResourceResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type SetDocumentContentParams struct {
	FrameId FrameId `json:"frameId"` // Frame id to set HTML for.
	Html    string  `json:"html"`    // HTML content to set.
}

type SetDocumentContentCB func(err error)

// Sets given markup as the document's HTML.
type SetDocumentContentCommand struct {
	params *SetDocumentContentParams
	cb     SetDocumentContentCB
}

func NewSetDocumentContentCommand(params *SetDocumentContentParams, cb SetDocumentContentCB) *SetDocumentContentCommand {
	return &SetDocumentContentCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetDocumentContentCommand) Name() string {
	return "Page.setDocumentContent"
}

func (cmd *SetDocumentContentCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetDocumentContentCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type PageSetDeviceMetricsOverrideParams struct {
	Width             int                `json:"width"`             // Overriding width value in pixels (minimum 0, maximum 10000000). 0 disables the override.
	Height            int                `json:"height"`            // Overriding height value in pixels (minimum 0, maximum 10000000). 0 disables the override.
	DeviceScaleFactor int                `json:"deviceScaleFactor"` // Overriding device scale factor value. 0 disables the override.
	Mobile            bool               `json:"mobile"`            // Whether to emulate mobile device. This includes viewport meta tag, overlay scrollbars, text autosizing and more.
	FitWindow         bool               `json:"fitWindow"`         // Whether a view that exceeds the available browser window area should be scaled down to fit.
	Scale             int                `json:"scale"`             // Scale to apply to resulting view image. Ignored in |fitWindow| mode.
	OffsetX           int                `json:"offsetX"`           // X offset to shift resulting view image by. Ignored in |fitWindow| mode.
	OffsetY           int                `json:"offsetY"`           // Y offset to shift resulting view image by. Ignored in |fitWindow| mode.
	ScreenWidth       int                `json:"screenWidth"`       // Overriding screen width value in pixels (minimum 0, maximum 10000000). Only used for |mobile==true|.
	ScreenHeight      int                `json:"screenHeight"`      // Overriding screen height value in pixels (minimum 0, maximum 10000000). Only used for |mobile==true|.
	PositionX         int                `json:"positionX"`         // Overriding view X position on screen in pixels (minimum 0, maximum 10000000). Only used for |mobile==true|.
	PositionY         int                `json:"positionY"`         // Overriding view Y position on screen in pixels (minimum 0, maximum 10000000). Only used for |mobile==true|.
	ScreenOrientation *ScreenOrientation `json:"screenOrientation"` // Screen orientation override.
}

type PageSetDeviceMetricsOverrideCB func(err error)

// Overrides the values of device screen dimensions (window.screen.width, window.screen.height, window.innerWidth, window.innerHeight, and "device-width"/"device-height"-related CSS media query results).
type PageSetDeviceMetricsOverrideCommand struct {
	params *PageSetDeviceMetricsOverrideParams
	cb     PageSetDeviceMetricsOverrideCB
}

func NewPageSetDeviceMetricsOverrideCommand(params *PageSetDeviceMetricsOverrideParams, cb PageSetDeviceMetricsOverrideCB) *PageSetDeviceMetricsOverrideCommand {
	return &PageSetDeviceMetricsOverrideCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *PageSetDeviceMetricsOverrideCommand) Name() string {
	return "Page.setDeviceMetricsOverride"
}

func (cmd *PageSetDeviceMetricsOverrideCommand) Params() interface{} {
	return cmd.params
}

func (cmd *PageSetDeviceMetricsOverrideCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type PageClearDeviceMetricsOverrideCB func(err error)

// Clears the overriden device metrics.
type PageClearDeviceMetricsOverrideCommand struct {
	cb PageClearDeviceMetricsOverrideCB
}

func NewPageClearDeviceMetricsOverrideCommand(cb PageClearDeviceMetricsOverrideCB) *PageClearDeviceMetricsOverrideCommand {
	return &PageClearDeviceMetricsOverrideCommand{
		cb: cb,
	}
}

func (cmd *PageClearDeviceMetricsOverrideCommand) Name() string {
	return "Page.clearDeviceMetricsOverride"
}

func (cmd *PageClearDeviceMetricsOverrideCommand) Params() interface{} {
	return nil
}

func (cmd *PageClearDeviceMetricsOverrideCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type PageSetGeolocationOverrideParams struct {
	Latitude  int `json:"latitude"`  // Mock latitude
	Longitude int `json:"longitude"` // Mock longitude
	Accuracy  int `json:"accuracy"`  // Mock accuracy
}

type PageSetGeolocationOverrideCB func(err error)

// Overrides the Geolocation Position or Error. Omitting any of the parameters emulates position unavailable.
type PageSetGeolocationOverrideCommand struct {
	params *PageSetGeolocationOverrideParams
	cb     PageSetGeolocationOverrideCB
}

func NewPageSetGeolocationOverrideCommand(params *PageSetGeolocationOverrideParams, cb PageSetGeolocationOverrideCB) *PageSetGeolocationOverrideCommand {
	return &PageSetGeolocationOverrideCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *PageSetGeolocationOverrideCommand) Name() string {
	return "Page.setGeolocationOverride"
}

func (cmd *PageSetGeolocationOverrideCommand) Params() interface{} {
	return cmd.params
}

func (cmd *PageSetGeolocationOverrideCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type PageClearGeolocationOverrideCB func(err error)

// Clears the overriden Geolocation Position and Error.
type PageClearGeolocationOverrideCommand struct {
	cb PageClearGeolocationOverrideCB
}

func NewPageClearGeolocationOverrideCommand(cb PageClearGeolocationOverrideCB) *PageClearGeolocationOverrideCommand {
	return &PageClearGeolocationOverrideCommand{
		cb: cb,
	}
}

func (cmd *PageClearGeolocationOverrideCommand) Name() string {
	return "Page.clearGeolocationOverride"
}

func (cmd *PageClearGeolocationOverrideCommand) Params() interface{} {
	return nil
}

func (cmd *PageClearGeolocationOverrideCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type PageSetDeviceOrientationOverrideParams struct {
	Alpha int `json:"alpha"` // Mock alpha
	Beta  int `json:"beta"`  // Mock beta
	Gamma int `json:"gamma"` // Mock gamma
}

type PageSetDeviceOrientationOverrideCB func(err error)

// Overrides the Device Orientation.
type PageSetDeviceOrientationOverrideCommand struct {
	params *PageSetDeviceOrientationOverrideParams
	cb     PageSetDeviceOrientationOverrideCB
}

func NewPageSetDeviceOrientationOverrideCommand(params *PageSetDeviceOrientationOverrideParams, cb PageSetDeviceOrientationOverrideCB) *PageSetDeviceOrientationOverrideCommand {
	return &PageSetDeviceOrientationOverrideCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *PageSetDeviceOrientationOverrideCommand) Name() string {
	return "Page.setDeviceOrientationOverride"
}

func (cmd *PageSetDeviceOrientationOverrideCommand) Params() interface{} {
	return cmd.params
}

func (cmd *PageSetDeviceOrientationOverrideCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type PageClearDeviceOrientationOverrideCB func(err error)

// Clears the overridden Device Orientation.
type PageClearDeviceOrientationOverrideCommand struct {
	cb PageClearDeviceOrientationOverrideCB
}

func NewPageClearDeviceOrientationOverrideCommand(cb PageClearDeviceOrientationOverrideCB) *PageClearDeviceOrientationOverrideCommand {
	return &PageClearDeviceOrientationOverrideCommand{
		cb: cb,
	}
}

func (cmd *PageClearDeviceOrientationOverrideCommand) Name() string {
	return "Page.clearDeviceOrientationOverride"
}

func (cmd *PageClearDeviceOrientationOverrideCommand) Params() interface{} {
	return nil
}

func (cmd *PageClearDeviceOrientationOverrideCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type PageSetTouchEmulationEnabledParams struct {
	Enabled       bool   `json:"enabled"`       // Whether the touch event emulation should be enabled.
	Configuration string `json:"configuration"` // Touch/gesture events configuration. Default: current platform.
}

type PageSetTouchEmulationEnabledCB func(err error)

// Toggles mouse event-based touch event emulation.
type PageSetTouchEmulationEnabledCommand struct {
	params *PageSetTouchEmulationEnabledParams
	cb     PageSetTouchEmulationEnabledCB
}

func NewPageSetTouchEmulationEnabledCommand(params *PageSetTouchEmulationEnabledParams, cb PageSetTouchEmulationEnabledCB) *PageSetTouchEmulationEnabledCommand {
	return &PageSetTouchEmulationEnabledCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *PageSetTouchEmulationEnabledCommand) Name() string {
	return "Page.setTouchEmulationEnabled"
}

func (cmd *PageSetTouchEmulationEnabledCommand) Params() interface{} {
	return cmd.params
}

func (cmd *PageSetTouchEmulationEnabledCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type CaptureScreenshotResult struct {
	Data string `json:"data"` // Base64-encoded image data (PNG).
}

type CaptureScreenshotCB func(result *CaptureScreenshotResult, err error)

// Capture page screenshot.
type CaptureScreenshotCommand struct {
	cb CaptureScreenshotCB
}

func NewCaptureScreenshotCommand(cb CaptureScreenshotCB) *CaptureScreenshotCommand {
	return &CaptureScreenshotCommand{
		cb: cb,
	}
}

func (cmd *CaptureScreenshotCommand) Name() string {
	return "Page.captureScreenshot"
}

func (cmd *CaptureScreenshotCommand) Params() interface{} {
	return nil
}

func (cmd *CaptureScreenshotCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj CaptureScreenshotResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type StartScreencastParams struct {
	Format        string `json:"format"`        // Image compression format.
	Quality       int    `json:"quality"`       // Compression quality from range [0..100].
	MaxWidth      int    `json:"maxWidth"`      // Maximum screenshot width.
	MaxHeight     int    `json:"maxHeight"`     // Maximum screenshot height.
	EveryNthFrame int    `json:"everyNthFrame"` // Send every n-th frame.
}

type StartScreencastCB func(err error)

// Starts sending each frame using the screencastFrame event.
type StartScreencastCommand struct {
	params *StartScreencastParams
	cb     StartScreencastCB
}

func NewStartScreencastCommand(params *StartScreencastParams, cb StartScreencastCB) *StartScreencastCommand {
	return &StartScreencastCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *StartScreencastCommand) Name() string {
	return "Page.startScreencast"
}

func (cmd *StartScreencastCommand) Params() interface{} {
	return cmd.params
}

func (cmd *StartScreencastCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type StopScreencastCB func(err error)

// Stops sending each frame in the screencastFrame.
type StopScreencastCommand struct {
	cb StopScreencastCB
}

func NewStopScreencastCommand(cb StopScreencastCB) *StopScreencastCommand {
	return &StopScreencastCommand{
		cb: cb,
	}
}

func (cmd *StopScreencastCommand) Name() string {
	return "Page.stopScreencast"
}

func (cmd *StopScreencastCommand) Params() interface{} {
	return nil
}

func (cmd *StopScreencastCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type ScreencastFrameAckParams struct {
	SessionId int `json:"sessionId"` // Frame number.
}

type ScreencastFrameAckCB func(err error)

// Acknowledges that a screencast frame has been received by the frontend.
type ScreencastFrameAckCommand struct {
	params *ScreencastFrameAckParams
	cb     ScreencastFrameAckCB
}

func NewScreencastFrameAckCommand(params *ScreencastFrameAckParams, cb ScreencastFrameAckCB) *ScreencastFrameAckCommand {
	return &ScreencastFrameAckCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *ScreencastFrameAckCommand) Name() string {
	return "Page.screencastFrameAck"
}

func (cmd *ScreencastFrameAckCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ScreencastFrameAckCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type HandleJavaScriptDialogParams struct {
	Accept     bool   `json:"accept"`     // Whether to accept or dismiss the dialog.
	PromptText string `json:"promptText"` // The text to enter into the dialog prompt before accepting. Used only if this is a prompt dialog.
}

type HandleJavaScriptDialogCB func(err error)

// Accepts or dismisses a JavaScript initiated dialog (alert, confirm, prompt, or onbeforeunload).
type HandleJavaScriptDialogCommand struct {
	params *HandleJavaScriptDialogParams
	cb     HandleJavaScriptDialogCB
}

func NewHandleJavaScriptDialogCommand(params *HandleJavaScriptDialogParams, cb HandleJavaScriptDialogCB) *HandleJavaScriptDialogCommand {
	return &HandleJavaScriptDialogCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *HandleJavaScriptDialogCommand) Name() string {
	return "Page.handleJavaScriptDialog"
}

func (cmd *HandleJavaScriptDialogCommand) Params() interface{} {
	return cmd.params
}

func (cmd *HandleJavaScriptDialogCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetColorPickerEnabledParams struct {
	Enabled bool `json:"enabled"` // Shows / hides color picker
}

type SetColorPickerEnabledCB func(err error)

// Shows / hides color picker
type SetColorPickerEnabledCommand struct {
	params *SetColorPickerEnabledParams
	cb     SetColorPickerEnabledCB
}

func NewSetColorPickerEnabledCommand(params *SetColorPickerEnabledParams, cb SetColorPickerEnabledCB) *SetColorPickerEnabledCommand {
	return &SetColorPickerEnabledCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetColorPickerEnabledCommand) Name() string {
	return "Page.setColorPickerEnabled"
}

func (cmd *SetColorPickerEnabledCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetColorPickerEnabledCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type ConfigureOverlayParams struct {
	Suspended bool   `json:"suspended"` // Whether overlay should be suspended and not consume any resources.
	Message   string `json:"message"`   // Overlay message to display.
}

type ConfigureOverlayCB func(err error)

// Configures overlay.
type ConfigureOverlayCommand struct {
	params *ConfigureOverlayParams
	cb     ConfigureOverlayCB
}

func NewConfigureOverlayCommand(params *ConfigureOverlayParams, cb ConfigureOverlayCB) *ConfigureOverlayCommand {
	return &ConfigureOverlayCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *ConfigureOverlayCommand) Name() string {
	return "Page.configureOverlay"
}

func (cmd *ConfigureOverlayCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ConfigureOverlayCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type GetAppManifestResult struct {
	Url    string              `json:"url"` // Manifest location.
	Errors []*AppManifestError `json:"errors"`
	Data   string              `json:"data"` // Manifest content.
}

type GetAppManifestCB func(result *GetAppManifestResult, err error)

type GetAppManifestCommand struct {
	cb GetAppManifestCB
}

func NewGetAppManifestCommand(cb GetAppManifestCB) *GetAppManifestCommand {
	return &GetAppManifestCommand{
		cb: cb,
	}
}

func (cmd *GetAppManifestCommand) Name() string {
	return "Page.getAppManifest"
}

func (cmd *GetAppManifestCommand) Params() interface{} {
	return nil
}

func (cmd *GetAppManifestCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetAppManifestResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type RequestAppBannerCB func(err error)

type RequestAppBannerCommand struct {
	cb RequestAppBannerCB
}

func NewRequestAppBannerCommand(cb RequestAppBannerCB) *RequestAppBannerCommand {
	return &RequestAppBannerCommand{
		cb: cb,
	}
}

func (cmd *RequestAppBannerCommand) Name() string {
	return "Page.requestAppBanner"
}

func (cmd *RequestAppBannerCommand) Params() interface{} {
	return nil
}

func (cmd *RequestAppBannerCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetBlockedEventsWarningThresholdParams struct {
	Threshold int `json:"threshold"` // If set to a positive number, specifies threshold in seconds for input event latency that will cause a console warning about blocked event to be issued. If zero or less, the warning is disabled.
}

type SetBlockedEventsWarningThresholdCB func(err error)

type SetBlockedEventsWarningThresholdCommand struct {
	params *SetBlockedEventsWarningThresholdParams
	cb     SetBlockedEventsWarningThresholdCB
}

func NewSetBlockedEventsWarningThresholdCommand(params *SetBlockedEventsWarningThresholdParams, cb SetBlockedEventsWarningThresholdCB) *SetBlockedEventsWarningThresholdCommand {
	return &SetBlockedEventsWarningThresholdCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetBlockedEventsWarningThresholdCommand) Name() string {
	return "Page.setBlockedEventsWarningThreshold"
}

func (cmd *SetBlockedEventsWarningThresholdCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetBlockedEventsWarningThresholdCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetControlNavigationsParams struct {
	Enabled bool `json:"enabled"`
}

type SetControlNavigationsCB func(err error)

// Toggles navigation throttling which allows programatic control over navigation and redirect response.
type SetControlNavigationsCommand struct {
	params *SetControlNavigationsParams
	cb     SetControlNavigationsCB
}

func NewSetControlNavigationsCommand(params *SetControlNavigationsParams, cb SetControlNavigationsCB) *SetControlNavigationsCommand {
	return &SetControlNavigationsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetControlNavigationsCommand) Name() string {
	return "Page.setControlNavigations"
}

func (cmd *SetControlNavigationsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetControlNavigationsCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type ProcessNavigationParams struct {
	Response     NavigationResponse `json:"response"`
	NavigationId int                `json:"navigationId"`
}

type ProcessNavigationCB func(err error)

// Should be sent in response to a navigationRequested or a redirectRequested event, telling the browser how to handle the navigation.
type ProcessNavigationCommand struct {
	params *ProcessNavigationParams
	cb     ProcessNavigationCB
}

func NewProcessNavigationCommand(params *ProcessNavigationParams, cb ProcessNavigationCB) *ProcessNavigationCommand {
	return &ProcessNavigationCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *ProcessNavigationCommand) Name() string {
	return "Page.processNavigation"
}

func (cmd *ProcessNavigationCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ProcessNavigationCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type DomContentEventFiredEvent struct {
	Timestamp int `json:"timestamp"`
}

type DomContentEventFiredEventSink struct {
	events chan *DomContentEventFiredEvent
}

func NewDomContentEventFiredEventSink(bufSize int) *DomContentEventFiredEventSink {
	return &DomContentEventFiredEventSink{
		events: make(chan *DomContentEventFiredEvent, bufSize),
	}
}

func (s *DomContentEventFiredEventSink) Name() string {
	return "Page.domContentEventFired"
}

func (s *DomContentEventFiredEventSink) OnEvent(params []byte) {
	evt := &DomContentEventFiredEvent{}
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

type LoadEventFiredEvent struct {
	Timestamp int `json:"timestamp"`
}

type LoadEventFiredEventSink struct {
	events chan *LoadEventFiredEvent
}

func NewLoadEventFiredEventSink(bufSize int) *LoadEventFiredEventSink {
	return &LoadEventFiredEventSink{
		events: make(chan *LoadEventFiredEvent, bufSize),
	}
}

func (s *LoadEventFiredEventSink) Name() string {
	return "Page.loadEventFired"
}

func (s *LoadEventFiredEventSink) OnEvent(params []byte) {
	evt := &LoadEventFiredEvent{}
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

type FrameAttachedEvent struct {
	FrameId       FrameId `json:"frameId"`       // Id of the frame that has been attached.
	ParentFrameId FrameId `json:"parentFrameId"` // Parent frame identifier.
}

// Fired when frame has been attached to its parent.
type FrameAttachedEventSink struct {
	events chan *FrameAttachedEvent
}

func NewFrameAttachedEventSink(bufSize int) *FrameAttachedEventSink {
	return &FrameAttachedEventSink{
		events: make(chan *FrameAttachedEvent, bufSize),
	}
}

func (s *FrameAttachedEventSink) Name() string {
	return "Page.frameAttached"
}

func (s *FrameAttachedEventSink) OnEvent(params []byte) {
	evt := &FrameAttachedEvent{}
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

type FrameNavigatedEvent struct {
	Frame *Frame `json:"frame"` // Frame object.
}

// Fired once navigation of the frame has completed. Frame is now associated with the new loader.
type FrameNavigatedEventSink struct {
	events chan *FrameNavigatedEvent
}

func NewFrameNavigatedEventSink(bufSize int) *FrameNavigatedEventSink {
	return &FrameNavigatedEventSink{
		events: make(chan *FrameNavigatedEvent, bufSize),
	}
}

func (s *FrameNavigatedEventSink) Name() string {
	return "Page.frameNavigated"
}

func (s *FrameNavigatedEventSink) OnEvent(params []byte) {
	evt := &FrameNavigatedEvent{}
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

type FrameDetachedEvent struct {
	FrameId FrameId `json:"frameId"` // Id of the frame that has been detached.
}

// Fired when frame has been detached from its parent.
type FrameDetachedEventSink struct {
	events chan *FrameDetachedEvent
}

func NewFrameDetachedEventSink(bufSize int) *FrameDetachedEventSink {
	return &FrameDetachedEventSink{
		events: make(chan *FrameDetachedEvent, bufSize),
	}
}

func (s *FrameDetachedEventSink) Name() string {
	return "Page.frameDetached"
}

func (s *FrameDetachedEventSink) OnEvent(params []byte) {
	evt := &FrameDetachedEvent{}
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

type FrameStartedLoadingEvent struct {
	FrameId FrameId `json:"frameId"` // Id of the frame that has started loading.
}

// Fired when frame has started loading.
type FrameStartedLoadingEventSink struct {
	events chan *FrameStartedLoadingEvent
}

func NewFrameStartedLoadingEventSink(bufSize int) *FrameStartedLoadingEventSink {
	return &FrameStartedLoadingEventSink{
		events: make(chan *FrameStartedLoadingEvent, bufSize),
	}
}

func (s *FrameStartedLoadingEventSink) Name() string {
	return "Page.frameStartedLoading"
}

func (s *FrameStartedLoadingEventSink) OnEvent(params []byte) {
	evt := &FrameStartedLoadingEvent{}
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

type FrameStoppedLoadingEvent struct {
	FrameId FrameId `json:"frameId"` // Id of the frame that has stopped loading.
}

// Fired when frame has stopped loading.
type FrameStoppedLoadingEventSink struct {
	events chan *FrameStoppedLoadingEvent
}

func NewFrameStoppedLoadingEventSink(bufSize int) *FrameStoppedLoadingEventSink {
	return &FrameStoppedLoadingEventSink{
		events: make(chan *FrameStoppedLoadingEvent, bufSize),
	}
}

func (s *FrameStoppedLoadingEventSink) Name() string {
	return "Page.frameStoppedLoading"
}

func (s *FrameStoppedLoadingEventSink) OnEvent(params []byte) {
	evt := &FrameStoppedLoadingEvent{}
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

type FrameScheduledNavigationEvent struct {
	FrameId FrameId `json:"frameId"` // Id of the frame that has scheduled a navigation.
	Delay   int     `json:"delay"`   // Delay (in seconds) until the navigation is scheduled to begin. The navigation is not guaranteed to start.
}

// Fired when frame schedules a potential navigation.
type FrameScheduledNavigationEventSink struct {
	events chan *FrameScheduledNavigationEvent
}

func NewFrameScheduledNavigationEventSink(bufSize int) *FrameScheduledNavigationEventSink {
	return &FrameScheduledNavigationEventSink{
		events: make(chan *FrameScheduledNavigationEvent, bufSize),
	}
}

func (s *FrameScheduledNavigationEventSink) Name() string {
	return "Page.frameScheduledNavigation"
}

func (s *FrameScheduledNavigationEventSink) OnEvent(params []byte) {
	evt := &FrameScheduledNavigationEvent{}
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

type FrameClearedScheduledNavigationEvent struct {
	FrameId FrameId `json:"frameId"` // Id of the frame that has cleared its scheduled navigation.
}

// Fired when frame no longer has a scheduled navigation.
type FrameClearedScheduledNavigationEventSink struct {
	events chan *FrameClearedScheduledNavigationEvent
}

func NewFrameClearedScheduledNavigationEventSink(bufSize int) *FrameClearedScheduledNavigationEventSink {
	return &FrameClearedScheduledNavigationEventSink{
		events: make(chan *FrameClearedScheduledNavigationEvent, bufSize),
	}
}

func (s *FrameClearedScheduledNavigationEventSink) Name() string {
	return "Page.frameClearedScheduledNavigation"
}

func (s *FrameClearedScheduledNavigationEventSink) OnEvent(params []byte) {
	evt := &FrameClearedScheduledNavigationEvent{}
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

type FrameResizedEvent struct {
}

type FrameResizedEventSink struct {
	events chan *FrameResizedEvent
}

func NewFrameResizedEventSink(bufSize int) *FrameResizedEventSink {
	return &FrameResizedEventSink{
		events: make(chan *FrameResizedEvent, bufSize),
	}
}

func (s *FrameResizedEventSink) Name() string {
	return "Page.frameResized"
}

func (s *FrameResizedEventSink) OnEvent(params []byte) {
	evt := &FrameResizedEvent{}
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

type JavascriptDialogOpeningEvent struct {
	Message string     `json:"message"` // Message that will be displayed by the dialog.
	Type    DialogType `json:"type"`    // Dialog type.
}

// Fired when a JavaScript initiated dialog (alert, confirm, prompt, or onbeforeunload) is about to open.
type JavascriptDialogOpeningEventSink struct {
	events chan *JavascriptDialogOpeningEvent
}

func NewJavascriptDialogOpeningEventSink(bufSize int) *JavascriptDialogOpeningEventSink {
	return &JavascriptDialogOpeningEventSink{
		events: make(chan *JavascriptDialogOpeningEvent, bufSize),
	}
}

func (s *JavascriptDialogOpeningEventSink) Name() string {
	return "Page.javascriptDialogOpening"
}

func (s *JavascriptDialogOpeningEventSink) OnEvent(params []byte) {
	evt := &JavascriptDialogOpeningEvent{}
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

type JavascriptDialogClosedEvent struct {
	Result bool `json:"result"` // Whether dialog was confirmed.
}

// Fired when a JavaScript initiated dialog (alert, confirm, prompt, or onbeforeunload) has been closed.
type JavascriptDialogClosedEventSink struct {
	events chan *JavascriptDialogClosedEvent
}

func NewJavascriptDialogClosedEventSink(bufSize int) *JavascriptDialogClosedEventSink {
	return &JavascriptDialogClosedEventSink{
		events: make(chan *JavascriptDialogClosedEvent, bufSize),
	}
}

func (s *JavascriptDialogClosedEventSink) Name() string {
	return "Page.javascriptDialogClosed"
}

func (s *JavascriptDialogClosedEventSink) OnEvent(params []byte) {
	evt := &JavascriptDialogClosedEvent{}
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

type ScreencastFrameEvent struct {
	Data      string                   `json:"data"`      // Base64-encoded compressed image.
	Metadata  *ScreencastFrameMetadata `json:"metadata"`  // Screencast frame metadata.
	SessionId int                      `json:"sessionId"` // Frame number.
}

// Compressed image data requested by the startScreencast.
type ScreencastFrameEventSink struct {
	events chan *ScreencastFrameEvent
}

func NewScreencastFrameEventSink(bufSize int) *ScreencastFrameEventSink {
	return &ScreencastFrameEventSink{
		events: make(chan *ScreencastFrameEvent, bufSize),
	}
}

func (s *ScreencastFrameEventSink) Name() string {
	return "Page.screencastFrame"
}

func (s *ScreencastFrameEventSink) OnEvent(params []byte) {
	evt := &ScreencastFrameEvent{}
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

type ScreencastVisibilityChangedEvent struct {
	Visible bool `json:"visible"` // True if the page is visible.
}

// Fired when the page with currently enabled screencast was shown or hidden .
type ScreencastVisibilityChangedEventSink struct {
	events chan *ScreencastVisibilityChangedEvent
}

func NewScreencastVisibilityChangedEventSink(bufSize int) *ScreencastVisibilityChangedEventSink {
	return &ScreencastVisibilityChangedEventSink{
		events: make(chan *ScreencastVisibilityChangedEvent, bufSize),
	}
}

func (s *ScreencastVisibilityChangedEventSink) Name() string {
	return "Page.screencastVisibilityChanged"
}

func (s *ScreencastVisibilityChangedEventSink) OnEvent(params []byte) {
	evt := &ScreencastVisibilityChangedEvent{}
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

type ColorPickedEvent struct {
	Color *RGBA `json:"color"` // RGBA of the picked color.
}

// Fired when a color has been picked.
type ColorPickedEventSink struct {
	events chan *ColorPickedEvent
}

func NewColorPickedEventSink(bufSize int) *ColorPickedEventSink {
	return &ColorPickedEventSink{
		events: make(chan *ColorPickedEvent, bufSize),
	}
}

func (s *ColorPickedEventSink) Name() string {
	return "Page.colorPicked"
}

func (s *ColorPickedEventSink) OnEvent(params []byte) {
	evt := &ColorPickedEvent{}
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

type InterstitialShownEvent struct {
}

// Fired when interstitial page was shown
type InterstitialShownEventSink struct {
	events chan *InterstitialShownEvent
}

func NewInterstitialShownEventSink(bufSize int) *InterstitialShownEventSink {
	return &InterstitialShownEventSink{
		events: make(chan *InterstitialShownEvent, bufSize),
	}
}

func (s *InterstitialShownEventSink) Name() string {
	return "Page.interstitialShown"
}

func (s *InterstitialShownEventSink) OnEvent(params []byte) {
	evt := &InterstitialShownEvent{}
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

type InterstitialHiddenEvent struct {
}

// Fired when interstitial page was hidden
type InterstitialHiddenEventSink struct {
	events chan *InterstitialHiddenEvent
}

func NewInterstitialHiddenEventSink(bufSize int) *InterstitialHiddenEventSink {
	return &InterstitialHiddenEventSink{
		events: make(chan *InterstitialHiddenEvent, bufSize),
	}
}

func (s *InterstitialHiddenEventSink) Name() string {
	return "Page.interstitialHidden"
}

func (s *InterstitialHiddenEventSink) OnEvent(params []byte) {
	evt := &InterstitialHiddenEvent{}
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

type NavigationRequestedEvent struct {
	IsInMainFrame bool   `json:"isInMainFrame"` // Whether the navigation is taking place in the main frame or in a subframe.
	IsRedirect    bool   `json:"isRedirect"`    // Whether the navigation has encountered a server redirect or not.
	NavigationId  int    `json:"navigationId"`
	Url           string `json:"url"` // URL of requested navigation.
}

// Fired when a navigation is started if navigation throttles are enabled.  The navigation will be deferred until processNavigation is called.
type NavigationRequestedEventSink struct {
	events chan *NavigationRequestedEvent
}

func NewNavigationRequestedEventSink(bufSize int) *NavigationRequestedEventSink {
	return &NavigationRequestedEventSink{
		events: make(chan *NavigationRequestedEvent, bufSize),
	}
}

func (s *NavigationRequestedEventSink) Name() string {
	return "Page.navigationRequested"
}

func (s *NavigationRequestedEventSink) OnEvent(params []byte) {
	evt := &NavigationRequestedEvent{}
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
