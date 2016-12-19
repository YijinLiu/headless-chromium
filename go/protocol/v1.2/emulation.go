package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
)

// Screen orientation.
type ScreenOrientation struct {
	Type  string `json:"type"`  // Orientation type.
	Angle int    `json:"angle"` // Orientation angle.
}

// advance: If the scheduler runs out of immediate work, the virtual time base may fast forward to allow the next delayed task (if any) to run; pause: The virtual time base may not advance; pauseIfNetworkFetchesPending: The virtual time base may not advance if there are any pending resource fetches.
type VirtualTimePolicy string

const VirtualTimePolicyAdvance VirtualTimePolicy = "advance"
const VirtualTimePolicyPause VirtualTimePolicy = "pause"
const VirtualTimePolicyPauseIfNetworkFetchesPending VirtualTimePolicy = "pauseIfNetworkFetchesPending"

type EmulationSetDeviceMetricsOverrideParams struct {
	Width             int                `json:"width"`             // Overriding width value in pixels (minimum 0, maximum 10000000). 0 disables the override.
	Height            int                `json:"height"`            // Overriding height value in pixels (minimum 0, maximum 10000000). 0 disables the override.
	DeviceScaleFactor int                `json:"deviceScaleFactor"` // Overriding device scale factor value. 0 disables the override.
	Mobile            bool               `json:"mobile"`            // Whether to emulate mobile device. This includes viewport meta tag, overlay scrollbars, text autosizing and more.
	FitWindow         bool               `json:"fitWindow"`         // Whether a view that exceeds the available browser window area should be scaled down to fit.
	Scale             int                `json:"scale"`             // Scale to apply to resulting view image. Ignored in |fitWindow| mode.
	OffsetX           int                `json:"offsetX"`           // Not used.
	OffsetY           int                `json:"offsetY"`           // Not used.
	ScreenWidth       int                `json:"screenWidth"`       // Overriding screen width value in pixels (minimum 0, maximum 10000000). Only used for |mobile==true|.
	ScreenHeight      int                `json:"screenHeight"`      // Overriding screen height value in pixels (minimum 0, maximum 10000000). Only used for |mobile==true|.
	PositionX         int                `json:"positionX"`         // Overriding view X position on screen in pixels (minimum 0, maximum 10000000). Only used for |mobile==true|.
	PositionY         int                `json:"positionY"`         // Overriding view Y position on screen in pixels (minimum 0, maximum 10000000). Only used for |mobile==true|.
	ScreenOrientation *ScreenOrientation `json:"screenOrientation"` // Screen orientation override.
}

type EmulationSetDeviceMetricsOverrideCB func(err error)

// Overrides the values of device screen dimensions (window.screen.width, window.screen.height, window.innerWidth, window.innerHeight, and "device-width"/"device-height"-related CSS media query results).
type EmulationSetDeviceMetricsOverrideCommand struct {
	params *EmulationSetDeviceMetricsOverrideParams
	cb     EmulationSetDeviceMetricsOverrideCB
}

func NewEmulationSetDeviceMetricsOverrideCommand(params *EmulationSetDeviceMetricsOverrideParams, cb EmulationSetDeviceMetricsOverrideCB) *EmulationSetDeviceMetricsOverrideCommand {
	return &EmulationSetDeviceMetricsOverrideCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *EmulationSetDeviceMetricsOverrideCommand) Name() string {
	return "Emulation.setDeviceMetricsOverride"
}

func (cmd *EmulationSetDeviceMetricsOverrideCommand) Params() interface{} {
	return cmd.params
}

func (cmd *EmulationSetDeviceMetricsOverrideCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type EmulationClearDeviceMetricsOverrideCB func(err error)

// Clears the overriden device metrics.
type EmulationClearDeviceMetricsOverrideCommand struct {
	cb EmulationClearDeviceMetricsOverrideCB
}

func NewEmulationClearDeviceMetricsOverrideCommand(cb EmulationClearDeviceMetricsOverrideCB) *EmulationClearDeviceMetricsOverrideCommand {
	return &EmulationClearDeviceMetricsOverrideCommand{
		cb: cb,
	}
}

func (cmd *EmulationClearDeviceMetricsOverrideCommand) Name() string {
	return "Emulation.clearDeviceMetricsOverride"
}

func (cmd *EmulationClearDeviceMetricsOverrideCommand) Params() interface{} {
	return nil
}

func (cmd *EmulationClearDeviceMetricsOverrideCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type ForceViewportParams struct {
	X     int `json:"x"`     // X coordinate of top-left corner of the area (CSS pixels).
	Y     int `json:"y"`     // Y coordinate of top-left corner of the area (CSS pixels).
	Scale int `json:"scale"` // Scale to apply to the area (relative to a page scale of 1.0).
}

type ForceViewportCB func(err error)

// Overrides the visible area of the page. The change is hidden from the page, i.e. the observable scroll position and page scale does not change. In effect, the command moves the specified area of the page into the top-left corner of the frame.
type ForceViewportCommand struct {
	params *ForceViewportParams
	cb     ForceViewportCB
}

func NewForceViewportCommand(params *ForceViewportParams, cb ForceViewportCB) *ForceViewportCommand {
	return &ForceViewportCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *ForceViewportCommand) Name() string {
	return "Emulation.forceViewport"
}

func (cmd *ForceViewportCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ForceViewportCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type ResetViewportCB func(err error)

// Resets the visible area of the page to the original viewport, undoing any effects of the forceViewport command.
type ResetViewportCommand struct {
	cb ResetViewportCB
}

func NewResetViewportCommand(cb ResetViewportCB) *ResetViewportCommand {
	return &ResetViewportCommand{
		cb: cb,
	}
}

func (cmd *ResetViewportCommand) Name() string {
	return "Emulation.resetViewport"
}

func (cmd *ResetViewportCommand) Params() interface{} {
	return nil
}

func (cmd *ResetViewportCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type ResetPageScaleFactorCB func(err error)

// Requests that page scale factor is reset to initial values.
type ResetPageScaleFactorCommand struct {
	cb ResetPageScaleFactorCB
}

func NewResetPageScaleFactorCommand(cb ResetPageScaleFactorCB) *ResetPageScaleFactorCommand {
	return &ResetPageScaleFactorCommand{
		cb: cb,
	}
}

func (cmd *ResetPageScaleFactorCommand) Name() string {
	return "Emulation.resetPageScaleFactor"
}

func (cmd *ResetPageScaleFactorCommand) Params() interface{} {
	return nil
}

func (cmd *ResetPageScaleFactorCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetPageScaleFactorParams struct {
	PageScaleFactor int `json:"pageScaleFactor"` // Page scale factor.
}

type SetPageScaleFactorCB func(err error)

// Sets a specified page scale factor.
type SetPageScaleFactorCommand struct {
	params *SetPageScaleFactorParams
	cb     SetPageScaleFactorCB
}

func NewSetPageScaleFactorCommand(params *SetPageScaleFactorParams, cb SetPageScaleFactorCB) *SetPageScaleFactorCommand {
	return &SetPageScaleFactorCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetPageScaleFactorCommand) Name() string {
	return "Emulation.setPageScaleFactor"
}

func (cmd *SetPageScaleFactorCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetPageScaleFactorCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetVisibleSizeParams struct {
	Width  int `json:"width"`  // Frame width (DIP).
	Height int `json:"height"` // Frame height (DIP).
}

type SetVisibleSizeCB func(err error)

// Resizes the frame/viewport of the page. Note that this does not affect the frame's container (e.g. browser window). Can be used to produce screenshots of the specified size. Not supported on Android.
type SetVisibleSizeCommand struct {
	params *SetVisibleSizeParams
	cb     SetVisibleSizeCB
}

func NewSetVisibleSizeCommand(params *SetVisibleSizeParams, cb SetVisibleSizeCB) *SetVisibleSizeCommand {
	return &SetVisibleSizeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetVisibleSizeCommand) Name() string {
	return "Emulation.setVisibleSize"
}

func (cmd *SetVisibleSizeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetVisibleSizeCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetScriptExecutionDisabledParams struct {
	Value bool `json:"value"` // Whether script execution should be disabled in the page.
}

type SetScriptExecutionDisabledCB func(err error)

// Switches script execution in the page.
type SetScriptExecutionDisabledCommand struct {
	params *SetScriptExecutionDisabledParams
	cb     SetScriptExecutionDisabledCB
}

func NewSetScriptExecutionDisabledCommand(params *SetScriptExecutionDisabledParams, cb SetScriptExecutionDisabledCB) *SetScriptExecutionDisabledCommand {
	return &SetScriptExecutionDisabledCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetScriptExecutionDisabledCommand) Name() string {
	return "Emulation.setScriptExecutionDisabled"
}

func (cmd *SetScriptExecutionDisabledCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetScriptExecutionDisabledCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type EmulationSetGeolocationOverrideParams struct {
	Latitude  int `json:"latitude"`  // Mock latitude
	Longitude int `json:"longitude"` // Mock longitude
	Accuracy  int `json:"accuracy"`  // Mock accuracy
}

type EmulationSetGeolocationOverrideCB func(err error)

// Overrides the Geolocation Position or Error. Omitting any of the parameters emulates position unavailable.
type EmulationSetGeolocationOverrideCommand struct {
	params *EmulationSetGeolocationOverrideParams
	cb     EmulationSetGeolocationOverrideCB
}

func NewEmulationSetGeolocationOverrideCommand(params *EmulationSetGeolocationOverrideParams, cb EmulationSetGeolocationOverrideCB) *EmulationSetGeolocationOverrideCommand {
	return &EmulationSetGeolocationOverrideCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *EmulationSetGeolocationOverrideCommand) Name() string {
	return "Emulation.setGeolocationOverride"
}

func (cmd *EmulationSetGeolocationOverrideCommand) Params() interface{} {
	return cmd.params
}

func (cmd *EmulationSetGeolocationOverrideCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type EmulationClearGeolocationOverrideCB func(err error)

// Clears the overriden Geolocation Position and Error.
type EmulationClearGeolocationOverrideCommand struct {
	cb EmulationClearGeolocationOverrideCB
}

func NewEmulationClearGeolocationOverrideCommand(cb EmulationClearGeolocationOverrideCB) *EmulationClearGeolocationOverrideCommand {
	return &EmulationClearGeolocationOverrideCommand{
		cb: cb,
	}
}

func (cmd *EmulationClearGeolocationOverrideCommand) Name() string {
	return "Emulation.clearGeolocationOverride"
}

func (cmd *EmulationClearGeolocationOverrideCommand) Params() interface{} {
	return nil
}

func (cmd *EmulationClearGeolocationOverrideCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type EmulationSetTouchEmulationEnabledParams struct {
	Enabled       bool   `json:"enabled"`       // Whether the touch event emulation should be enabled.
	Configuration string `json:"configuration"` // Touch/gesture events configuration. Default: current platform.
}

type EmulationSetTouchEmulationEnabledCB func(err error)

// Toggles mouse event-based touch event emulation.
type EmulationSetTouchEmulationEnabledCommand struct {
	params *EmulationSetTouchEmulationEnabledParams
	cb     EmulationSetTouchEmulationEnabledCB
}

func NewEmulationSetTouchEmulationEnabledCommand(params *EmulationSetTouchEmulationEnabledParams, cb EmulationSetTouchEmulationEnabledCB) *EmulationSetTouchEmulationEnabledCommand {
	return &EmulationSetTouchEmulationEnabledCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *EmulationSetTouchEmulationEnabledCommand) Name() string {
	return "Emulation.setTouchEmulationEnabled"
}

func (cmd *EmulationSetTouchEmulationEnabledCommand) Params() interface{} {
	return cmd.params
}

func (cmd *EmulationSetTouchEmulationEnabledCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetEmulatedMediaParams struct {
	Media string `json:"media"` // Media type to emulate. Empty string disables the override.
}

type SetEmulatedMediaCB func(err error)

// Emulates the given media for CSS media queries.
type SetEmulatedMediaCommand struct {
	params *SetEmulatedMediaParams
	cb     SetEmulatedMediaCB
}

func NewSetEmulatedMediaCommand(params *SetEmulatedMediaParams, cb SetEmulatedMediaCB) *SetEmulatedMediaCommand {
	return &SetEmulatedMediaCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetEmulatedMediaCommand) Name() string {
	return "Emulation.setEmulatedMedia"
}

func (cmd *SetEmulatedMediaCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetEmulatedMediaCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetCPUThrottlingRateParams struct {
	Rate int `json:"rate"` // Throttling rate as a slowdown factor (1 is no throttle, 2 is 2x slowdown, etc).
}

type SetCPUThrottlingRateCB func(err error)

// Enables CPU throttling to emulate slow CPUs.
type SetCPUThrottlingRateCommand struct {
	params *SetCPUThrottlingRateParams
	cb     SetCPUThrottlingRateCB
}

func NewSetCPUThrottlingRateCommand(params *SetCPUThrottlingRateParams, cb SetCPUThrottlingRateCB) *SetCPUThrottlingRateCommand {
	return &SetCPUThrottlingRateCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetCPUThrottlingRateCommand) Name() string {
	return "Emulation.setCPUThrottlingRate"
}

func (cmd *SetCPUThrottlingRateCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetCPUThrottlingRateCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type CanEmulateResult struct {
	Result bool `json:"result"` // True if emulation is supported.
}

type CanEmulateCB func(result *CanEmulateResult, err error)

// Tells whether emulation is supported.
type CanEmulateCommand struct {
	cb CanEmulateCB
}

func NewCanEmulateCommand(cb CanEmulateCB) *CanEmulateCommand {
	return &CanEmulateCommand{
		cb: cb,
	}
}

func (cmd *CanEmulateCommand) Name() string {
	return "Emulation.canEmulate"
}

func (cmd *CanEmulateCommand) Params() interface{} {
	return nil
}

func (cmd *CanEmulateCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj CanEmulateResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type SetVirtualTimePolicyParams struct {
	Policy VirtualTimePolicy `json:"policy"`
	Budget int               `json:"budget"` // If set, after this many virtual milliseconds have elapsed virtual time will be paused and a virtualTimeBudgetExpired event is sent.
}

type SetVirtualTimePolicyCB func(err error)

// Turns on virtual time for all frames (replacing real-time with a synthetic time source) and sets the current virtual time policy.  Note this supersedes any previous time budget.
type SetVirtualTimePolicyCommand struct {
	params *SetVirtualTimePolicyParams
	cb     SetVirtualTimePolicyCB
}

func NewSetVirtualTimePolicyCommand(params *SetVirtualTimePolicyParams, cb SetVirtualTimePolicyCB) *SetVirtualTimePolicyCommand {
	return &SetVirtualTimePolicyCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetVirtualTimePolicyCommand) Name() string {
	return "Emulation.setVirtualTimePolicy"
}

func (cmd *SetVirtualTimePolicyCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetVirtualTimePolicyCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type VirtualTimeBudgetExpiredEvent struct {
}

// Notification sent after the virual time budget for the current VirtualTimePolicy has run out.
type VirtualTimeBudgetExpiredEventSink struct {
	events chan *VirtualTimeBudgetExpiredEvent
}

func NewVirtualTimeBudgetExpiredEventSink(bufSize int) *VirtualTimeBudgetExpiredEventSink {
	return &VirtualTimeBudgetExpiredEventSink{
		events: make(chan *VirtualTimeBudgetExpiredEvent, bufSize),
	}
}

func (s *VirtualTimeBudgetExpiredEventSink) Name() string {
	return "Emulation.virtualTimeBudgetExpired"
}

func (s *VirtualTimeBudgetExpiredEventSink) OnEvent(params []byte) {
	evt := &VirtualTimeBudgetExpiredEvent{}
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
