package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

// Screen orientation.
type ScreenOrientation struct {
	Type  string `json:"type"`  // Orientation type.
	Angle int    `json:"angle"` // Orientation angle.
}

// advance: If the scheduler runs out of immediate work, the virtual time base may fast forward to allow the next delayed task (if any) to run; pause: The virtual time base may not advance; pauseIfNetworkFetchesPending: The virtual time base may not advance if there are any pending resource fetches.
// @experimental
type VirtualTimePolicy string

const VirtualTimePolicyAdvance VirtualTimePolicy = "advance"
const VirtualTimePolicyPause VirtualTimePolicy = "pause"
const VirtualTimePolicyPauseIfNetworkFetchesPending VirtualTimePolicy = "pauseIfNetworkFetchesPending"

type EmulationSetDeviceMetricsOverrideParams struct {
	Width             int                `json:"width"`                       // Overriding width value in pixels (minimum 0, maximum 10000000). 0 disables the override.
	Height            int                `json:"height"`                      // Overriding height value in pixels (minimum 0, maximum 10000000). 0 disables the override.
	DeviceScaleFactor float64            `json:"deviceScaleFactor"`           // Overriding device scale factor value. 0 disables the override.
	Mobile            bool               `json:"mobile"`                      // Whether to emulate mobile device. This includes viewport meta tag, overlay scrollbars, text autosizing and more.
	FitWindow         bool               `json:"fitWindow"`                   // Whether a view that exceeds the available browser window area should be scaled down to fit.
	Scale             float64            `json:"scale,omitempty"`             // Scale to apply to resulting view image. Ignored in |fitWindow| mode.
	OffsetX           float64            `json:"offsetX,omitempty"`           // Not used.
	OffsetY           float64            `json:"offsetY,omitempty"`           // Not used.
	ScreenWidth       int                `json:"screenWidth,omitempty"`       // Overriding screen width value in pixels (minimum 0, maximum 10000000). Only used for |mobile==true|.
	ScreenHeight      int                `json:"screenHeight,omitempty"`      // Overriding screen height value in pixels (minimum 0, maximum 10000000). Only used for |mobile==true|.
	PositionX         int                `json:"positionX,omitempty"`         // Overriding view X position on screen in pixels (minimum 0, maximum 10000000). Only used for |mobile==true|.
	PositionY         int                `json:"positionY,omitempty"`         // Overriding view Y position on screen in pixels (minimum 0, maximum 10000000). Only used for |mobile==true|.
	ScreenOrientation *ScreenOrientation `json:"screenOrientation,omitempty"` // Screen orientation override.
}

// Overrides the values of device screen dimensions (window.screen.width, window.screen.height, window.innerWidth, window.innerHeight, and "device-width"/"device-height"-related CSS media query results).

type EmulationSetDeviceMetricsOverrideCommand struct {
	params *EmulationSetDeviceMetricsOverrideParams
	wg     sync.WaitGroup
	err    error
}

func NewEmulationSetDeviceMetricsOverrideCommand(params *EmulationSetDeviceMetricsOverrideParams) *EmulationSetDeviceMetricsOverrideCommand {
	return &EmulationSetDeviceMetricsOverrideCommand{
		params: params,
	}
}

func (cmd *EmulationSetDeviceMetricsOverrideCommand) Name() string {
	return "Emulation.setDeviceMetricsOverride"
}

func (cmd *EmulationSetDeviceMetricsOverrideCommand) Params() interface{} {
	return cmd.params
}

func (cmd *EmulationSetDeviceMetricsOverrideCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func EmulationSetDeviceMetricsOverride(params *EmulationSetDeviceMetricsOverrideParams, conn *hc.Conn) (err error) {
	cmd := NewEmulationSetDeviceMetricsOverrideCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type EmulationSetDeviceMetricsOverrideCB func(err error)

// Overrides the values of device screen dimensions (window.screen.width, window.screen.height, window.innerWidth, window.innerHeight, and "device-width"/"device-height"-related CSS media query results).

type AsyncEmulationSetDeviceMetricsOverrideCommand struct {
	params *EmulationSetDeviceMetricsOverrideParams
	cb     EmulationSetDeviceMetricsOverrideCB
}

func NewAsyncEmulationSetDeviceMetricsOverrideCommand(params *EmulationSetDeviceMetricsOverrideParams, cb EmulationSetDeviceMetricsOverrideCB) *AsyncEmulationSetDeviceMetricsOverrideCommand {
	return &AsyncEmulationSetDeviceMetricsOverrideCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncEmulationSetDeviceMetricsOverrideCommand) Name() string {
	return "Emulation.setDeviceMetricsOverride"
}

func (cmd *AsyncEmulationSetDeviceMetricsOverrideCommand) Params() interface{} {
	return cmd.params
}

func (cmd *EmulationSetDeviceMetricsOverrideCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncEmulationSetDeviceMetricsOverrideCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Clears the overriden device metrics.

type EmulationClearDeviceMetricsOverrideCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewEmulationClearDeviceMetricsOverrideCommand() *EmulationClearDeviceMetricsOverrideCommand {
	return &EmulationClearDeviceMetricsOverrideCommand{}
}

func (cmd *EmulationClearDeviceMetricsOverrideCommand) Name() string {
	return "Emulation.clearDeviceMetricsOverride"
}

func (cmd *EmulationClearDeviceMetricsOverrideCommand) Params() interface{} {
	return nil
}

func (cmd *EmulationClearDeviceMetricsOverrideCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func EmulationClearDeviceMetricsOverride(conn *hc.Conn) (err error) {
	cmd := NewEmulationClearDeviceMetricsOverrideCommand()
	cmd.Run(conn)
	return cmd.err
}

type EmulationClearDeviceMetricsOverrideCB func(err error)

// Clears the overriden device metrics.

type AsyncEmulationClearDeviceMetricsOverrideCommand struct {
	cb EmulationClearDeviceMetricsOverrideCB
}

func NewAsyncEmulationClearDeviceMetricsOverrideCommand(cb EmulationClearDeviceMetricsOverrideCB) *AsyncEmulationClearDeviceMetricsOverrideCommand {
	return &AsyncEmulationClearDeviceMetricsOverrideCommand{
		cb: cb,
	}
}

func (cmd *AsyncEmulationClearDeviceMetricsOverrideCommand) Name() string {
	return "Emulation.clearDeviceMetricsOverride"
}

func (cmd *AsyncEmulationClearDeviceMetricsOverrideCommand) Params() interface{} {
	return nil
}

func (cmd *EmulationClearDeviceMetricsOverrideCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncEmulationClearDeviceMetricsOverrideCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type ForceViewportParams struct {
	X     float64 `json:"x"`     // X coordinate of top-left corner of the area (CSS pixels).
	Y     float64 `json:"y"`     // Y coordinate of top-left corner of the area (CSS pixels).
	Scale float64 `json:"scale"` // Scale to apply to the area (relative to a page scale of 1.0).
}

// Overrides the visible area of the page. The change is hidden from the page, i.e. the observable scroll position and page scale does not change. In effect, the command moves the specified area of the page into the top-left corner of the frame.
// @experimental
type ForceViewportCommand struct {
	params *ForceViewportParams
	wg     sync.WaitGroup
	err    error
}

func NewForceViewportCommand(params *ForceViewportParams) *ForceViewportCommand {
	return &ForceViewportCommand{
		params: params,
	}
}

func (cmd *ForceViewportCommand) Name() string {
	return "Emulation.forceViewport"
}

func (cmd *ForceViewportCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ForceViewportCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ForceViewport(params *ForceViewportParams, conn *hc.Conn) (err error) {
	cmd := NewForceViewportCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type ForceViewportCB func(err error)

// Overrides the visible area of the page. The change is hidden from the page, i.e. the observable scroll position and page scale does not change. In effect, the command moves the specified area of the page into the top-left corner of the frame.
// @experimental
type AsyncForceViewportCommand struct {
	params *ForceViewportParams
	cb     ForceViewportCB
}

func NewAsyncForceViewportCommand(params *ForceViewportParams, cb ForceViewportCB) *AsyncForceViewportCommand {
	return &AsyncForceViewportCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncForceViewportCommand) Name() string {
	return "Emulation.forceViewport"
}

func (cmd *AsyncForceViewportCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ForceViewportCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncForceViewportCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Resets the visible area of the page to the original viewport, undoing any effects of the forceViewport command.
// @experimental
type ResetViewportCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewResetViewportCommand() *ResetViewportCommand {
	return &ResetViewportCommand{}
}

func (cmd *ResetViewportCommand) Name() string {
	return "Emulation.resetViewport"
}

func (cmd *ResetViewportCommand) Params() interface{} {
	return nil
}

func (cmd *ResetViewportCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ResetViewport(conn *hc.Conn) (err error) {
	cmd := NewResetViewportCommand()
	cmd.Run(conn)
	return cmd.err
}

type ResetViewportCB func(err error)

// Resets the visible area of the page to the original viewport, undoing any effects of the forceViewport command.
// @experimental
type AsyncResetViewportCommand struct {
	cb ResetViewportCB
}

func NewAsyncResetViewportCommand(cb ResetViewportCB) *AsyncResetViewportCommand {
	return &AsyncResetViewportCommand{
		cb: cb,
	}
}

func (cmd *AsyncResetViewportCommand) Name() string {
	return "Emulation.resetViewport"
}

func (cmd *AsyncResetViewportCommand) Params() interface{} {
	return nil
}

func (cmd *ResetViewportCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncResetViewportCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Requests that page scale factor is reset to initial values.
// @experimental
type ResetPageScaleFactorCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewResetPageScaleFactorCommand() *ResetPageScaleFactorCommand {
	return &ResetPageScaleFactorCommand{}
}

func (cmd *ResetPageScaleFactorCommand) Name() string {
	return "Emulation.resetPageScaleFactor"
}

func (cmd *ResetPageScaleFactorCommand) Params() interface{} {
	return nil
}

func (cmd *ResetPageScaleFactorCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ResetPageScaleFactor(conn *hc.Conn) (err error) {
	cmd := NewResetPageScaleFactorCommand()
	cmd.Run(conn)
	return cmd.err
}

type ResetPageScaleFactorCB func(err error)

// Requests that page scale factor is reset to initial values.
// @experimental
type AsyncResetPageScaleFactorCommand struct {
	cb ResetPageScaleFactorCB
}

func NewAsyncResetPageScaleFactorCommand(cb ResetPageScaleFactorCB) *AsyncResetPageScaleFactorCommand {
	return &AsyncResetPageScaleFactorCommand{
		cb: cb,
	}
}

func (cmd *AsyncResetPageScaleFactorCommand) Name() string {
	return "Emulation.resetPageScaleFactor"
}

func (cmd *AsyncResetPageScaleFactorCommand) Params() interface{} {
	return nil
}

func (cmd *ResetPageScaleFactorCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncResetPageScaleFactorCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetPageScaleFactorParams struct {
	PageScaleFactor float64 `json:"pageScaleFactor"` // Page scale factor.
}

// Sets a specified page scale factor.
// @experimental
type SetPageScaleFactorCommand struct {
	params *SetPageScaleFactorParams
	wg     sync.WaitGroup
	err    error
}

func NewSetPageScaleFactorCommand(params *SetPageScaleFactorParams) *SetPageScaleFactorCommand {
	return &SetPageScaleFactorCommand{
		params: params,
	}
}

func (cmd *SetPageScaleFactorCommand) Name() string {
	return "Emulation.setPageScaleFactor"
}

func (cmd *SetPageScaleFactorCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetPageScaleFactorCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetPageScaleFactor(params *SetPageScaleFactorParams, conn *hc.Conn) (err error) {
	cmd := NewSetPageScaleFactorCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetPageScaleFactorCB func(err error)

// Sets a specified page scale factor.
// @experimental
type AsyncSetPageScaleFactorCommand struct {
	params *SetPageScaleFactorParams
	cb     SetPageScaleFactorCB
}

func NewAsyncSetPageScaleFactorCommand(params *SetPageScaleFactorParams, cb SetPageScaleFactorCB) *AsyncSetPageScaleFactorCommand {
	return &AsyncSetPageScaleFactorCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetPageScaleFactorCommand) Name() string {
	return "Emulation.setPageScaleFactor"
}

func (cmd *AsyncSetPageScaleFactorCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetPageScaleFactorCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetPageScaleFactorCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetVisibleSizeParams struct {
	Width  int `json:"width"`  // Frame width (DIP).
	Height int `json:"height"` // Frame height (DIP).
}

// Resizes the frame/viewport of the page. Note that this does not affect the frame's container (e.g. browser window). Can be used to produce screenshots of the specified size. Not supported on Android.
// @experimental
type SetVisibleSizeCommand struct {
	params *SetVisibleSizeParams
	wg     sync.WaitGroup
	err    error
}

func NewSetVisibleSizeCommand(params *SetVisibleSizeParams) *SetVisibleSizeCommand {
	return &SetVisibleSizeCommand{
		params: params,
	}
}

func (cmd *SetVisibleSizeCommand) Name() string {
	return "Emulation.setVisibleSize"
}

func (cmd *SetVisibleSizeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetVisibleSizeCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetVisibleSize(params *SetVisibleSizeParams, conn *hc.Conn) (err error) {
	cmd := NewSetVisibleSizeCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetVisibleSizeCB func(err error)

// Resizes the frame/viewport of the page. Note that this does not affect the frame's container (e.g. browser window). Can be used to produce screenshots of the specified size. Not supported on Android.
// @experimental
type AsyncSetVisibleSizeCommand struct {
	params *SetVisibleSizeParams
	cb     SetVisibleSizeCB
}

func NewAsyncSetVisibleSizeCommand(params *SetVisibleSizeParams, cb SetVisibleSizeCB) *AsyncSetVisibleSizeCommand {
	return &AsyncSetVisibleSizeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetVisibleSizeCommand) Name() string {
	return "Emulation.setVisibleSize"
}

func (cmd *AsyncSetVisibleSizeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetVisibleSizeCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetVisibleSizeCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetScriptExecutionDisabledParams struct {
	Value bool `json:"value"` // Whether script execution should be disabled in the page.
}

// Switches script execution in the page.
// @experimental
type SetScriptExecutionDisabledCommand struct {
	params *SetScriptExecutionDisabledParams
	wg     sync.WaitGroup
	err    error
}

func NewSetScriptExecutionDisabledCommand(params *SetScriptExecutionDisabledParams) *SetScriptExecutionDisabledCommand {
	return &SetScriptExecutionDisabledCommand{
		params: params,
	}
}

func (cmd *SetScriptExecutionDisabledCommand) Name() string {
	return "Emulation.setScriptExecutionDisabled"
}

func (cmd *SetScriptExecutionDisabledCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetScriptExecutionDisabledCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetScriptExecutionDisabled(params *SetScriptExecutionDisabledParams, conn *hc.Conn) (err error) {
	cmd := NewSetScriptExecutionDisabledCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetScriptExecutionDisabledCB func(err error)

// Switches script execution in the page.
// @experimental
type AsyncSetScriptExecutionDisabledCommand struct {
	params *SetScriptExecutionDisabledParams
	cb     SetScriptExecutionDisabledCB
}

func NewAsyncSetScriptExecutionDisabledCommand(params *SetScriptExecutionDisabledParams, cb SetScriptExecutionDisabledCB) *AsyncSetScriptExecutionDisabledCommand {
	return &AsyncSetScriptExecutionDisabledCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetScriptExecutionDisabledCommand) Name() string {
	return "Emulation.setScriptExecutionDisabled"
}

func (cmd *AsyncSetScriptExecutionDisabledCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetScriptExecutionDisabledCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetScriptExecutionDisabledCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type EmulationSetGeolocationOverrideParams struct {
	Latitude  float64 `json:"latitude,omitempty"`  // Mock latitude
	Longitude float64 `json:"longitude,omitempty"` // Mock longitude
	Accuracy  float64 `json:"accuracy,omitempty"`  // Mock accuracy
}

// Overrides the Geolocation Position or Error. Omitting any of the parameters emulates position unavailable.
// @experimental
type EmulationSetGeolocationOverrideCommand struct {
	params *EmulationSetGeolocationOverrideParams
	wg     sync.WaitGroup
	err    error
}

func NewEmulationSetGeolocationOverrideCommand(params *EmulationSetGeolocationOverrideParams) *EmulationSetGeolocationOverrideCommand {
	return &EmulationSetGeolocationOverrideCommand{
		params: params,
	}
}

func (cmd *EmulationSetGeolocationOverrideCommand) Name() string {
	return "Emulation.setGeolocationOverride"
}

func (cmd *EmulationSetGeolocationOverrideCommand) Params() interface{} {
	return cmd.params
}

func (cmd *EmulationSetGeolocationOverrideCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func EmulationSetGeolocationOverride(params *EmulationSetGeolocationOverrideParams, conn *hc.Conn) (err error) {
	cmd := NewEmulationSetGeolocationOverrideCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type EmulationSetGeolocationOverrideCB func(err error)

// Overrides the Geolocation Position or Error. Omitting any of the parameters emulates position unavailable.
// @experimental
type AsyncEmulationSetGeolocationOverrideCommand struct {
	params *EmulationSetGeolocationOverrideParams
	cb     EmulationSetGeolocationOverrideCB
}

func NewAsyncEmulationSetGeolocationOverrideCommand(params *EmulationSetGeolocationOverrideParams, cb EmulationSetGeolocationOverrideCB) *AsyncEmulationSetGeolocationOverrideCommand {
	return &AsyncEmulationSetGeolocationOverrideCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncEmulationSetGeolocationOverrideCommand) Name() string {
	return "Emulation.setGeolocationOverride"
}

func (cmd *AsyncEmulationSetGeolocationOverrideCommand) Params() interface{} {
	return cmd.params
}

func (cmd *EmulationSetGeolocationOverrideCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncEmulationSetGeolocationOverrideCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Clears the overriden Geolocation Position and Error.
// @experimental
type EmulationClearGeolocationOverrideCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewEmulationClearGeolocationOverrideCommand() *EmulationClearGeolocationOverrideCommand {
	return &EmulationClearGeolocationOverrideCommand{}
}

func (cmd *EmulationClearGeolocationOverrideCommand) Name() string {
	return "Emulation.clearGeolocationOverride"
}

func (cmd *EmulationClearGeolocationOverrideCommand) Params() interface{} {
	return nil
}

func (cmd *EmulationClearGeolocationOverrideCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func EmulationClearGeolocationOverride(conn *hc.Conn) (err error) {
	cmd := NewEmulationClearGeolocationOverrideCommand()
	cmd.Run(conn)
	return cmd.err
}

type EmulationClearGeolocationOverrideCB func(err error)

// Clears the overriden Geolocation Position and Error.
// @experimental
type AsyncEmulationClearGeolocationOverrideCommand struct {
	cb EmulationClearGeolocationOverrideCB
}

func NewAsyncEmulationClearGeolocationOverrideCommand(cb EmulationClearGeolocationOverrideCB) *AsyncEmulationClearGeolocationOverrideCommand {
	return &AsyncEmulationClearGeolocationOverrideCommand{
		cb: cb,
	}
}

func (cmd *AsyncEmulationClearGeolocationOverrideCommand) Name() string {
	return "Emulation.clearGeolocationOverride"
}

func (cmd *AsyncEmulationClearGeolocationOverrideCommand) Params() interface{} {
	return nil
}

func (cmd *EmulationClearGeolocationOverrideCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncEmulationClearGeolocationOverrideCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type EmulationSetTouchEmulationEnabledParams struct {
	Enabled       bool   `json:"enabled"`                 // Whether the touch event emulation should be enabled.
	Configuration string `json:"configuration,omitempty"` // Touch/gesture events configuration. Default: current platform.
}

// Toggles mouse event-based touch event emulation.

type EmulationSetTouchEmulationEnabledCommand struct {
	params *EmulationSetTouchEmulationEnabledParams
	wg     sync.WaitGroup
	err    error
}

func NewEmulationSetTouchEmulationEnabledCommand(params *EmulationSetTouchEmulationEnabledParams) *EmulationSetTouchEmulationEnabledCommand {
	return &EmulationSetTouchEmulationEnabledCommand{
		params: params,
	}
}

func (cmd *EmulationSetTouchEmulationEnabledCommand) Name() string {
	return "Emulation.setTouchEmulationEnabled"
}

func (cmd *EmulationSetTouchEmulationEnabledCommand) Params() interface{} {
	return cmd.params
}

func (cmd *EmulationSetTouchEmulationEnabledCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func EmulationSetTouchEmulationEnabled(params *EmulationSetTouchEmulationEnabledParams, conn *hc.Conn) (err error) {
	cmd := NewEmulationSetTouchEmulationEnabledCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type EmulationSetTouchEmulationEnabledCB func(err error)

// Toggles mouse event-based touch event emulation.

type AsyncEmulationSetTouchEmulationEnabledCommand struct {
	params *EmulationSetTouchEmulationEnabledParams
	cb     EmulationSetTouchEmulationEnabledCB
}

func NewAsyncEmulationSetTouchEmulationEnabledCommand(params *EmulationSetTouchEmulationEnabledParams, cb EmulationSetTouchEmulationEnabledCB) *AsyncEmulationSetTouchEmulationEnabledCommand {
	return &AsyncEmulationSetTouchEmulationEnabledCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncEmulationSetTouchEmulationEnabledCommand) Name() string {
	return "Emulation.setTouchEmulationEnabled"
}

func (cmd *AsyncEmulationSetTouchEmulationEnabledCommand) Params() interface{} {
	return cmd.params
}

func (cmd *EmulationSetTouchEmulationEnabledCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncEmulationSetTouchEmulationEnabledCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetEmulatedMediaParams struct {
	Media string `json:"media"` // Media type to emulate. Empty string disables the override.
}

// Emulates the given media for CSS media queries.

type SetEmulatedMediaCommand struct {
	params *SetEmulatedMediaParams
	wg     sync.WaitGroup
	err    error
}

func NewSetEmulatedMediaCommand(params *SetEmulatedMediaParams) *SetEmulatedMediaCommand {
	return &SetEmulatedMediaCommand{
		params: params,
	}
}

func (cmd *SetEmulatedMediaCommand) Name() string {
	return "Emulation.setEmulatedMedia"
}

func (cmd *SetEmulatedMediaCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetEmulatedMediaCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetEmulatedMedia(params *SetEmulatedMediaParams, conn *hc.Conn) (err error) {
	cmd := NewSetEmulatedMediaCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetEmulatedMediaCB func(err error)

// Emulates the given media for CSS media queries.

type AsyncSetEmulatedMediaCommand struct {
	params *SetEmulatedMediaParams
	cb     SetEmulatedMediaCB
}

func NewAsyncSetEmulatedMediaCommand(params *SetEmulatedMediaParams, cb SetEmulatedMediaCB) *AsyncSetEmulatedMediaCommand {
	return &AsyncSetEmulatedMediaCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetEmulatedMediaCommand) Name() string {
	return "Emulation.setEmulatedMedia"
}

func (cmd *AsyncSetEmulatedMediaCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetEmulatedMediaCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetEmulatedMediaCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetCPUThrottlingRateParams struct {
	Rate float64 `json:"rate"` // Throttling rate as a slowdown factor (1 is no throttle, 2 is 2x slowdown, etc).
}

// Enables CPU throttling to emulate slow CPUs.
// @experimental
type SetCPUThrottlingRateCommand struct {
	params *SetCPUThrottlingRateParams
	wg     sync.WaitGroup
	err    error
}

func NewSetCPUThrottlingRateCommand(params *SetCPUThrottlingRateParams) *SetCPUThrottlingRateCommand {
	return &SetCPUThrottlingRateCommand{
		params: params,
	}
}

func (cmd *SetCPUThrottlingRateCommand) Name() string {
	return "Emulation.setCPUThrottlingRate"
}

func (cmd *SetCPUThrottlingRateCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetCPUThrottlingRateCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetCPUThrottlingRate(params *SetCPUThrottlingRateParams, conn *hc.Conn) (err error) {
	cmd := NewSetCPUThrottlingRateCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetCPUThrottlingRateCB func(err error)

// Enables CPU throttling to emulate slow CPUs.
// @experimental
type AsyncSetCPUThrottlingRateCommand struct {
	params *SetCPUThrottlingRateParams
	cb     SetCPUThrottlingRateCB
}

func NewAsyncSetCPUThrottlingRateCommand(params *SetCPUThrottlingRateParams, cb SetCPUThrottlingRateCB) *AsyncSetCPUThrottlingRateCommand {
	return &AsyncSetCPUThrottlingRateCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetCPUThrottlingRateCommand) Name() string {
	return "Emulation.setCPUThrottlingRate"
}

func (cmd *AsyncSetCPUThrottlingRateCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetCPUThrottlingRateCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetCPUThrottlingRateCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type CanEmulateResult struct {
	Result bool `json:"result"` // True if emulation is supported.
}

// Tells whether emulation is supported.
// @experimental
type CanEmulateCommand struct {
	result CanEmulateResult
	wg     sync.WaitGroup
	err    error
}

func NewCanEmulateCommand() *CanEmulateCommand {
	return &CanEmulateCommand{}
}

func (cmd *CanEmulateCommand) Name() string {
	return "Emulation.canEmulate"
}

func (cmd *CanEmulateCommand) Params() interface{} {
	return nil
}

func (cmd *CanEmulateCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func CanEmulate(conn *hc.Conn) (result *CanEmulateResult, err error) {
	cmd := NewCanEmulateCommand()
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type CanEmulateCB func(result *CanEmulateResult, err error)

// Tells whether emulation is supported.
// @experimental
type AsyncCanEmulateCommand struct {
	cb CanEmulateCB
}

func NewAsyncCanEmulateCommand(cb CanEmulateCB) *AsyncCanEmulateCommand {
	return &AsyncCanEmulateCommand{
		cb: cb,
	}
}

func (cmd *AsyncCanEmulateCommand) Name() string {
	return "Emulation.canEmulate"
}

func (cmd *AsyncCanEmulateCommand) Params() interface{} {
	return nil
}

func (cmd *CanEmulateCommand) Result() *CanEmulateResult {
	return &cmd.result
}

func (cmd *CanEmulateCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncCanEmulateCommand) Done(data []byte, err error) {
	var result CanEmulateResult
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

type SetVirtualTimePolicyParams struct {
	Policy VirtualTimePolicy `json:"policy"`
	Budget int               `json:"budget,omitempty"` // If set, after this many virtual milliseconds have elapsed virtual time will be paused and a virtualTimeBudgetExpired event is sent.
}

// Turns on virtual time for all frames (replacing real-time with a synthetic time source) and sets the current virtual time policy.  Note this supersedes any previous time budget.
// @experimental
type SetVirtualTimePolicyCommand struct {
	params *SetVirtualTimePolicyParams
	wg     sync.WaitGroup
	err    error
}

func NewSetVirtualTimePolicyCommand(params *SetVirtualTimePolicyParams) *SetVirtualTimePolicyCommand {
	return &SetVirtualTimePolicyCommand{
		params: params,
	}
}

func (cmd *SetVirtualTimePolicyCommand) Name() string {
	return "Emulation.setVirtualTimePolicy"
}

func (cmd *SetVirtualTimePolicyCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetVirtualTimePolicyCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetVirtualTimePolicy(params *SetVirtualTimePolicyParams, conn *hc.Conn) (err error) {
	cmd := NewSetVirtualTimePolicyCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetVirtualTimePolicyCB func(err error)

// Turns on virtual time for all frames (replacing real-time with a synthetic time source) and sets the current virtual time policy.  Note this supersedes any previous time budget.
// @experimental
type AsyncSetVirtualTimePolicyCommand struct {
	params *SetVirtualTimePolicyParams
	cb     SetVirtualTimePolicyCB
}

func NewAsyncSetVirtualTimePolicyCommand(params *SetVirtualTimePolicyParams, cb SetVirtualTimePolicyCB) *AsyncSetVirtualTimePolicyCommand {
	return &AsyncSetVirtualTimePolicyCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetVirtualTimePolicyCommand) Name() string {
	return "Emulation.setVirtualTimePolicy"
}

func (cmd *AsyncSetVirtualTimePolicyCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetVirtualTimePolicyCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetVirtualTimePolicyCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Notification sent after the virual time budget for the current VirtualTimePolicy has run out.
// @experimental
type VirtualTimeBudgetExpiredEvent struct {
}

func OnVirtualTimeBudgetExpired(conn *hc.Conn, cb func(evt *VirtualTimeBudgetExpiredEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &VirtualTimeBudgetExpiredEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Emulation.virtualTimeBudgetExpired", sink)
}
