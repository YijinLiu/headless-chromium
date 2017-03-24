package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

// Breakpoint identifier.
type BreakpointId string

// Call frame identifier.
type CallFrameId string

// Location in the source code.
type Location struct {
	ScriptId     *ScriptId `json:"scriptId"`               // Script identifier as reported in the Debugger.scriptParsed.
	LineNumber   int       `json:"lineNumber"`             // Line number in the script (0-based).
	ColumnNumber int       `json:"columnNumber,omitempty"` // Column number in the script (0-based).
}

// Location in the source code.
// @experimental
type ScriptPosition struct {
	LineNumber   int `json:"lineNumber"`
	ColumnNumber int `json:"columnNumber"`
}

// JavaScript call frame. Array of call frames form the call stack.
type DebuggerCallFrame struct {
	CallFrameId      CallFrameId   `json:"callFrameId"`                // Call frame identifier. This identifier is only valid while the virtual machine is paused.
	FunctionName     string        `json:"functionName"`               // Name of the JavaScript function called on this call frame.
	FunctionLocation *Location     `json:"functionLocation,omitempty"` // Location in the source code.
	Location         *Location     `json:"location"`                   // Location in the source code.
	ScopeChain       []*Scope      `json:"scopeChain"`                 // Scope chain for this call frame.
	This             *RemoteObject `json:"this"`                       // this object for this call frame.
	ReturnValue      *RemoteObject `json:"returnValue,omitempty"`      // The value being returned, if the function is at return point.
}

// Scope description.
type Scope struct {
	Type          string        `json:"type"`   // Scope type.
	Object        *RemoteObject `json:"object"` // Object representing the scope. For global and with scopes it represents the actual object; for the rest of the scopes, it is artificial transient object enumerating scope variables as its properties.
	Name          string        `json:"name,omitempty"`
	StartLocation *Location     `json:"startLocation,omitempty"` // Location in the source code where scope starts
	EndLocation   *Location     `json:"endLocation,omitempty"`   // Location in the source code where scope ends
}

// Search match for resource.
// @experimental
type SearchMatch struct {
	LineNumber  float64 `json:"lineNumber"`  // Line number in resource content.
	LineContent string  `json:"lineContent"` // Line with match content.
}

// Enables debugger for the given page. Clients should not assume that the debugging has been enabled until the result for this command is received.

type DebuggerEnableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewDebuggerEnableCommand() *DebuggerEnableCommand {
	return &DebuggerEnableCommand{}
}

func (cmd *DebuggerEnableCommand) Name() string {
	return "Debugger.enable"
}

func (cmd *DebuggerEnableCommand) Params() interface{} {
	return nil
}

func (cmd *DebuggerEnableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func DebuggerEnable(conn *hc.Conn) (err error) {
	cmd := NewDebuggerEnableCommand()
	cmd.Run(conn)
	return cmd.err
}

type DebuggerEnableCB func(err error)

// Enables debugger for the given page. Clients should not assume that the debugging has been enabled until the result for this command is received.

type AsyncDebuggerEnableCommand struct {
	cb DebuggerEnableCB
}

func NewAsyncDebuggerEnableCommand(cb DebuggerEnableCB) *AsyncDebuggerEnableCommand {
	return &AsyncDebuggerEnableCommand{
		cb: cb,
	}
}

func (cmd *AsyncDebuggerEnableCommand) Name() string {
	return "Debugger.enable"
}

func (cmd *AsyncDebuggerEnableCommand) Params() interface{} {
	return nil
}

func (cmd *DebuggerEnableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncDebuggerEnableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Disables debugger for given page.

type DebuggerDisableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewDebuggerDisableCommand() *DebuggerDisableCommand {
	return &DebuggerDisableCommand{}
}

func (cmd *DebuggerDisableCommand) Name() string {
	return "Debugger.disable"
}

func (cmd *DebuggerDisableCommand) Params() interface{} {
	return nil
}

func (cmd *DebuggerDisableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func DebuggerDisable(conn *hc.Conn) (err error) {
	cmd := NewDebuggerDisableCommand()
	cmd.Run(conn)
	return cmd.err
}

type DebuggerDisableCB func(err error)

// Disables debugger for given page.

type AsyncDebuggerDisableCommand struct {
	cb DebuggerDisableCB
}

func NewAsyncDebuggerDisableCommand(cb DebuggerDisableCB) *AsyncDebuggerDisableCommand {
	return &AsyncDebuggerDisableCommand{
		cb: cb,
	}
}

func (cmd *AsyncDebuggerDisableCommand) Name() string {
	return "Debugger.disable"
}

func (cmd *AsyncDebuggerDisableCommand) Params() interface{} {
	return nil
}

func (cmd *DebuggerDisableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncDebuggerDisableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetBreakpointsActiveParams struct {
	Active bool `json:"active"` // New value for breakpoints active state.
}

// Activates / deactivates all breakpoints on the page.

type SetBreakpointsActiveCommand struct {
	params *SetBreakpointsActiveParams
	wg     sync.WaitGroup
	err    error
}

func NewSetBreakpointsActiveCommand(params *SetBreakpointsActiveParams) *SetBreakpointsActiveCommand {
	return &SetBreakpointsActiveCommand{
		params: params,
	}
}

func (cmd *SetBreakpointsActiveCommand) Name() string {
	return "Debugger.setBreakpointsActive"
}

func (cmd *SetBreakpointsActiveCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetBreakpointsActiveCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetBreakpointsActive(params *SetBreakpointsActiveParams, conn *hc.Conn) (err error) {
	cmd := NewSetBreakpointsActiveCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetBreakpointsActiveCB func(err error)

// Activates / deactivates all breakpoints on the page.

type AsyncSetBreakpointsActiveCommand struct {
	params *SetBreakpointsActiveParams
	cb     SetBreakpointsActiveCB
}

func NewAsyncSetBreakpointsActiveCommand(params *SetBreakpointsActiveParams, cb SetBreakpointsActiveCB) *AsyncSetBreakpointsActiveCommand {
	return &AsyncSetBreakpointsActiveCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetBreakpointsActiveCommand) Name() string {
	return "Debugger.setBreakpointsActive"
}

func (cmd *AsyncSetBreakpointsActiveCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetBreakpointsActiveCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetBreakpointsActiveCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetSkipAllPausesParams struct {
	Skip bool `json:"skip"` // New value for skip pauses state.
}

// Makes page not interrupt on any pauses (breakpoint, exception, dom exception etc).

type SetSkipAllPausesCommand struct {
	params *SetSkipAllPausesParams
	wg     sync.WaitGroup
	err    error
}

func NewSetSkipAllPausesCommand(params *SetSkipAllPausesParams) *SetSkipAllPausesCommand {
	return &SetSkipAllPausesCommand{
		params: params,
	}
}

func (cmd *SetSkipAllPausesCommand) Name() string {
	return "Debugger.setSkipAllPauses"
}

func (cmd *SetSkipAllPausesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetSkipAllPausesCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetSkipAllPauses(params *SetSkipAllPausesParams, conn *hc.Conn) (err error) {
	cmd := NewSetSkipAllPausesCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetSkipAllPausesCB func(err error)

// Makes page not interrupt on any pauses (breakpoint, exception, dom exception etc).

type AsyncSetSkipAllPausesCommand struct {
	params *SetSkipAllPausesParams
	cb     SetSkipAllPausesCB
}

func NewAsyncSetSkipAllPausesCommand(params *SetSkipAllPausesParams, cb SetSkipAllPausesCB) *AsyncSetSkipAllPausesCommand {
	return &AsyncSetSkipAllPausesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetSkipAllPausesCommand) Name() string {
	return "Debugger.setSkipAllPauses"
}

func (cmd *AsyncSetSkipAllPausesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetSkipAllPausesCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetSkipAllPausesCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetBreakpointByUrlParams struct {
	LineNumber   int    `json:"lineNumber"`             // Line number to set breakpoint at.
	Url          string `json:"url,omitempty"`          // URL of the resources to set breakpoint on.
	UrlRegex     string `json:"urlRegex,omitempty"`     // Regex pattern for the URLs of the resources to set breakpoints on. Either url or urlRegex must be specified.
	ColumnNumber int    `json:"columnNumber,omitempty"` // Offset in the line to set breakpoint at.
	Condition    string `json:"condition,omitempty"`    // Expression to use as a breakpoint condition. When specified, debugger will only stop on the breakpoint if this expression evaluates to true.
}

type SetBreakpointByUrlResult struct {
	BreakpointId BreakpointId `json:"breakpointId"` // Id of the created breakpoint for further reference.
	Locations    []*Location  `json:"locations"`    // List of the locations this breakpoint resolved into upon addition.
}

// Sets JavaScript breakpoint at given location specified either by URL or URL regex. Once this command is issued, all existing parsed scripts will have breakpoints resolved and returned in locations property. Further matching script parsing will result in subsequent breakpointResolved events issued. This logical breakpoint will survive page reloads.

type SetBreakpointByUrlCommand struct {
	params *SetBreakpointByUrlParams
	result SetBreakpointByUrlResult
	wg     sync.WaitGroup
	err    error
}

func NewSetBreakpointByUrlCommand(params *SetBreakpointByUrlParams) *SetBreakpointByUrlCommand {
	return &SetBreakpointByUrlCommand{
		params: params,
	}
}

func (cmd *SetBreakpointByUrlCommand) Name() string {
	return "Debugger.setBreakpointByUrl"
}

func (cmd *SetBreakpointByUrlCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetBreakpointByUrlCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetBreakpointByUrl(params *SetBreakpointByUrlParams, conn *hc.Conn) (result *SetBreakpointByUrlResult, err error) {
	cmd := NewSetBreakpointByUrlCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type SetBreakpointByUrlCB func(result *SetBreakpointByUrlResult, err error)

// Sets JavaScript breakpoint at given location specified either by URL or URL regex. Once this command is issued, all existing parsed scripts will have breakpoints resolved and returned in locations property. Further matching script parsing will result in subsequent breakpointResolved events issued. This logical breakpoint will survive page reloads.

type AsyncSetBreakpointByUrlCommand struct {
	params *SetBreakpointByUrlParams
	cb     SetBreakpointByUrlCB
}

func NewAsyncSetBreakpointByUrlCommand(params *SetBreakpointByUrlParams, cb SetBreakpointByUrlCB) *AsyncSetBreakpointByUrlCommand {
	return &AsyncSetBreakpointByUrlCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetBreakpointByUrlCommand) Name() string {
	return "Debugger.setBreakpointByUrl"
}

func (cmd *AsyncSetBreakpointByUrlCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetBreakpointByUrlCommand) Result() *SetBreakpointByUrlResult {
	return &cmd.result
}

func (cmd *SetBreakpointByUrlCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetBreakpointByUrlCommand) Done(data []byte, err error) {
	var result SetBreakpointByUrlResult
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

type SetBreakpointParams struct {
	Location  *Location `json:"location"`            // Location to set breakpoint in.
	Condition string    `json:"condition,omitempty"` // Expression to use as a breakpoint condition. When specified, debugger will only stop on the breakpoint if this expression evaluates to true.
}

type SetBreakpointResult struct {
	BreakpointId   BreakpointId `json:"breakpointId"`   // Id of the created breakpoint for further reference.
	ActualLocation *Location    `json:"actualLocation"` // Location this breakpoint resolved into.
}

// Sets JavaScript breakpoint at a given location.

type SetBreakpointCommand struct {
	params *SetBreakpointParams
	result SetBreakpointResult
	wg     sync.WaitGroup
	err    error
}

func NewSetBreakpointCommand(params *SetBreakpointParams) *SetBreakpointCommand {
	return &SetBreakpointCommand{
		params: params,
	}
}

func (cmd *SetBreakpointCommand) Name() string {
	return "Debugger.setBreakpoint"
}

func (cmd *SetBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetBreakpointCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetBreakpoint(params *SetBreakpointParams, conn *hc.Conn) (result *SetBreakpointResult, err error) {
	cmd := NewSetBreakpointCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type SetBreakpointCB func(result *SetBreakpointResult, err error)

// Sets JavaScript breakpoint at a given location.

type AsyncSetBreakpointCommand struct {
	params *SetBreakpointParams
	cb     SetBreakpointCB
}

func NewAsyncSetBreakpointCommand(params *SetBreakpointParams, cb SetBreakpointCB) *AsyncSetBreakpointCommand {
	return &AsyncSetBreakpointCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetBreakpointCommand) Name() string {
	return "Debugger.setBreakpoint"
}

func (cmd *AsyncSetBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetBreakpointCommand) Result() *SetBreakpointResult {
	return &cmd.result
}

func (cmd *SetBreakpointCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetBreakpointCommand) Done(data []byte, err error) {
	var result SetBreakpointResult
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

type RemoveBreakpointParams struct {
	BreakpointId BreakpointId `json:"breakpointId"`
}

// Removes JavaScript breakpoint.

type RemoveBreakpointCommand struct {
	params *RemoveBreakpointParams
	wg     sync.WaitGroup
	err    error
}

func NewRemoveBreakpointCommand(params *RemoveBreakpointParams) *RemoveBreakpointCommand {
	return &RemoveBreakpointCommand{
		params: params,
	}
}

func (cmd *RemoveBreakpointCommand) Name() string {
	return "Debugger.removeBreakpoint"
}

func (cmd *RemoveBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveBreakpointCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func RemoveBreakpoint(params *RemoveBreakpointParams, conn *hc.Conn) (err error) {
	cmd := NewRemoveBreakpointCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type RemoveBreakpointCB func(err error)

// Removes JavaScript breakpoint.

type AsyncRemoveBreakpointCommand struct {
	params *RemoveBreakpointParams
	cb     RemoveBreakpointCB
}

func NewAsyncRemoveBreakpointCommand(params *RemoveBreakpointParams, cb RemoveBreakpointCB) *AsyncRemoveBreakpointCommand {
	return &AsyncRemoveBreakpointCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncRemoveBreakpointCommand) Name() string {
	return "Debugger.removeBreakpoint"
}

func (cmd *AsyncRemoveBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveBreakpointCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncRemoveBreakpointCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type GetPossibleBreakpointsParams struct {
	Start *Location `json:"start"`         // Start of range to search possible breakpoint locations in.
	End   *Location `json:"end,omitempty"` // End of range to search possible breakpoint locations in (excluding). When not specifed, end of scripts is used as end of range.
}

type GetPossibleBreakpointsResult struct {
	Locations []*Location `json:"locations"` // List of the possible breakpoint locations.
}

// Returns possible locations for breakpoint. scriptId in start and end range locations should be the same.
// @experimental
type GetPossibleBreakpointsCommand struct {
	params *GetPossibleBreakpointsParams
	result GetPossibleBreakpointsResult
	wg     sync.WaitGroup
	err    error
}

func NewGetPossibleBreakpointsCommand(params *GetPossibleBreakpointsParams) *GetPossibleBreakpointsCommand {
	return &GetPossibleBreakpointsCommand{
		params: params,
	}
}

func (cmd *GetPossibleBreakpointsCommand) Name() string {
	return "Debugger.getPossibleBreakpoints"
}

func (cmd *GetPossibleBreakpointsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetPossibleBreakpointsCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetPossibleBreakpoints(params *GetPossibleBreakpointsParams, conn *hc.Conn) (result *GetPossibleBreakpointsResult, err error) {
	cmd := NewGetPossibleBreakpointsCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetPossibleBreakpointsCB func(result *GetPossibleBreakpointsResult, err error)

// Returns possible locations for breakpoint. scriptId in start and end range locations should be the same.
// @experimental
type AsyncGetPossibleBreakpointsCommand struct {
	params *GetPossibleBreakpointsParams
	cb     GetPossibleBreakpointsCB
}

func NewAsyncGetPossibleBreakpointsCommand(params *GetPossibleBreakpointsParams, cb GetPossibleBreakpointsCB) *AsyncGetPossibleBreakpointsCommand {
	return &AsyncGetPossibleBreakpointsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetPossibleBreakpointsCommand) Name() string {
	return "Debugger.getPossibleBreakpoints"
}

func (cmd *AsyncGetPossibleBreakpointsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetPossibleBreakpointsCommand) Result() *GetPossibleBreakpointsResult {
	return &cmd.result
}

func (cmd *GetPossibleBreakpointsCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetPossibleBreakpointsCommand) Done(data []byte, err error) {
	var result GetPossibleBreakpointsResult
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

type ContinueToLocationParams struct {
	Location *Location `json:"location"` // Location to continue to.
}

// Continues execution until specific location is reached.

type ContinueToLocationCommand struct {
	params *ContinueToLocationParams
	wg     sync.WaitGroup
	err    error
}

func NewContinueToLocationCommand(params *ContinueToLocationParams) *ContinueToLocationCommand {
	return &ContinueToLocationCommand{
		params: params,
	}
}

func (cmd *ContinueToLocationCommand) Name() string {
	return "Debugger.continueToLocation"
}

func (cmd *ContinueToLocationCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ContinueToLocationCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ContinueToLocation(params *ContinueToLocationParams, conn *hc.Conn) (err error) {
	cmd := NewContinueToLocationCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type ContinueToLocationCB func(err error)

// Continues execution until specific location is reached.

type AsyncContinueToLocationCommand struct {
	params *ContinueToLocationParams
	cb     ContinueToLocationCB
}

func NewAsyncContinueToLocationCommand(params *ContinueToLocationParams, cb ContinueToLocationCB) *AsyncContinueToLocationCommand {
	return &AsyncContinueToLocationCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncContinueToLocationCommand) Name() string {
	return "Debugger.continueToLocation"
}

func (cmd *AsyncContinueToLocationCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ContinueToLocationCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncContinueToLocationCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Steps over the statement.

type StepOverCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewStepOverCommand() *StepOverCommand {
	return &StepOverCommand{}
}

func (cmd *StepOverCommand) Name() string {
	return "Debugger.stepOver"
}

func (cmd *StepOverCommand) Params() interface{} {
	return nil
}

func (cmd *StepOverCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func StepOver(conn *hc.Conn) (err error) {
	cmd := NewStepOverCommand()
	cmd.Run(conn)
	return cmd.err
}

type StepOverCB func(err error)

// Steps over the statement.

type AsyncStepOverCommand struct {
	cb StepOverCB
}

func NewAsyncStepOverCommand(cb StepOverCB) *AsyncStepOverCommand {
	return &AsyncStepOverCommand{
		cb: cb,
	}
}

func (cmd *AsyncStepOverCommand) Name() string {
	return "Debugger.stepOver"
}

func (cmd *AsyncStepOverCommand) Params() interface{} {
	return nil
}

func (cmd *StepOverCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncStepOverCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Steps into the function call.

type StepIntoCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewStepIntoCommand() *StepIntoCommand {
	return &StepIntoCommand{}
}

func (cmd *StepIntoCommand) Name() string {
	return "Debugger.stepInto"
}

func (cmd *StepIntoCommand) Params() interface{} {
	return nil
}

func (cmd *StepIntoCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func StepInto(conn *hc.Conn) (err error) {
	cmd := NewStepIntoCommand()
	cmd.Run(conn)
	return cmd.err
}

type StepIntoCB func(err error)

// Steps into the function call.

type AsyncStepIntoCommand struct {
	cb StepIntoCB
}

func NewAsyncStepIntoCommand(cb StepIntoCB) *AsyncStepIntoCommand {
	return &AsyncStepIntoCommand{
		cb: cb,
	}
}

func (cmd *AsyncStepIntoCommand) Name() string {
	return "Debugger.stepInto"
}

func (cmd *AsyncStepIntoCommand) Params() interface{} {
	return nil
}

func (cmd *StepIntoCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncStepIntoCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Steps out of the function call.

type StepOutCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewStepOutCommand() *StepOutCommand {
	return &StepOutCommand{}
}

func (cmd *StepOutCommand) Name() string {
	return "Debugger.stepOut"
}

func (cmd *StepOutCommand) Params() interface{} {
	return nil
}

func (cmd *StepOutCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func StepOut(conn *hc.Conn) (err error) {
	cmd := NewStepOutCommand()
	cmd.Run(conn)
	return cmd.err
}

type StepOutCB func(err error)

// Steps out of the function call.

type AsyncStepOutCommand struct {
	cb StepOutCB
}

func NewAsyncStepOutCommand(cb StepOutCB) *AsyncStepOutCommand {
	return &AsyncStepOutCommand{
		cb: cb,
	}
}

func (cmd *AsyncStepOutCommand) Name() string {
	return "Debugger.stepOut"
}

func (cmd *AsyncStepOutCommand) Params() interface{} {
	return nil
}

func (cmd *StepOutCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncStepOutCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Stops on the next JavaScript statement.

type PauseCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewPauseCommand() *PauseCommand {
	return &PauseCommand{}
}

func (cmd *PauseCommand) Name() string {
	return "Debugger.pause"
}

func (cmd *PauseCommand) Params() interface{} {
	return nil
}

func (cmd *PauseCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func Pause(conn *hc.Conn) (err error) {
	cmd := NewPauseCommand()
	cmd.Run(conn)
	return cmd.err
}

type PauseCB func(err error)

// Stops on the next JavaScript statement.

type AsyncPauseCommand struct {
	cb PauseCB
}

func NewAsyncPauseCommand(cb PauseCB) *AsyncPauseCommand {
	return &AsyncPauseCommand{
		cb: cb,
	}
}

func (cmd *AsyncPauseCommand) Name() string {
	return "Debugger.pause"
}

func (cmd *AsyncPauseCommand) Params() interface{} {
	return nil
}

func (cmd *PauseCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncPauseCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Resumes JavaScript execution.

type ResumeCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewResumeCommand() *ResumeCommand {
	return &ResumeCommand{}
}

func (cmd *ResumeCommand) Name() string {
	return "Debugger.resume"
}

func (cmd *ResumeCommand) Params() interface{} {
	return nil
}

func (cmd *ResumeCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func Resume(conn *hc.Conn) (err error) {
	cmd := NewResumeCommand()
	cmd.Run(conn)
	return cmd.err
}

type ResumeCB func(err error)

// Resumes JavaScript execution.

type AsyncResumeCommand struct {
	cb ResumeCB
}

func NewAsyncResumeCommand(cb ResumeCB) *AsyncResumeCommand {
	return &AsyncResumeCommand{
		cb: cb,
	}
}

func (cmd *AsyncResumeCommand) Name() string {
	return "Debugger.resume"
}

func (cmd *AsyncResumeCommand) Params() interface{} {
	return nil
}

func (cmd *ResumeCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncResumeCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SearchInContentParams struct {
	ScriptId      *ScriptId `json:"scriptId"`                // Id of the script to search in.
	Query         string    `json:"query"`                   // String to search for.
	CaseSensitive bool      `json:"caseSensitive,omitempty"` // If true, search is case sensitive.
	IsRegex       bool      `json:"isRegex,omitempty"`       // If true, treats string parameter as regex.
}

type SearchInContentResult struct {
	Result []*SearchMatch `json:"result"` // List of search matches.
}

// Searches for given string in script content.
// @experimental
type SearchInContentCommand struct {
	params *SearchInContentParams
	result SearchInContentResult
	wg     sync.WaitGroup
	err    error
}

func NewSearchInContentCommand(params *SearchInContentParams) *SearchInContentCommand {
	return &SearchInContentCommand{
		params: params,
	}
}

func (cmd *SearchInContentCommand) Name() string {
	return "Debugger.searchInContent"
}

func (cmd *SearchInContentCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SearchInContentCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SearchInContent(params *SearchInContentParams, conn *hc.Conn) (result *SearchInContentResult, err error) {
	cmd := NewSearchInContentCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type SearchInContentCB func(result *SearchInContentResult, err error)

// Searches for given string in script content.
// @experimental
type AsyncSearchInContentCommand struct {
	params *SearchInContentParams
	cb     SearchInContentCB
}

func NewAsyncSearchInContentCommand(params *SearchInContentParams, cb SearchInContentCB) *AsyncSearchInContentCommand {
	return &AsyncSearchInContentCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSearchInContentCommand) Name() string {
	return "Debugger.searchInContent"
}

func (cmd *AsyncSearchInContentCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SearchInContentCommand) Result() *SearchInContentResult {
	return &cmd.result
}

func (cmd *SearchInContentCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSearchInContentCommand) Done(data []byte, err error) {
	var result SearchInContentResult
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

type SetScriptSourceParams struct {
	ScriptId     *ScriptId `json:"scriptId"`         // Id of the script to edit.
	ScriptSource string    `json:"scriptSource"`     // New content of the script.
	DryRun       bool      `json:"dryRun,omitempty"` //  If true the change will not actually be applied. Dry run may be used to get result description without actually modifying the code.
}

type SetScriptSourceResult struct {
	CallFrames       []*DebuggerCallFrame `json:"callFrames"`       // New stack trace in case editing has happened while VM was stopped.
	StackChanged     bool                 `json:"stackChanged"`     // Whether current call stack  was modified after applying the changes.
	AsyncStackTrace  *StackTrace          `json:"asyncStackTrace"`  // Async stack trace, if any.
	ExceptionDetails *ExceptionDetails    `json:"exceptionDetails"` // Exception details if any.
}

// Edits JavaScript source live.

type SetScriptSourceCommand struct {
	params *SetScriptSourceParams
	result SetScriptSourceResult
	wg     sync.WaitGroup
	err    error
}

func NewSetScriptSourceCommand(params *SetScriptSourceParams) *SetScriptSourceCommand {
	return &SetScriptSourceCommand{
		params: params,
	}
}

func (cmd *SetScriptSourceCommand) Name() string {
	return "Debugger.setScriptSource"
}

func (cmd *SetScriptSourceCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetScriptSourceCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetScriptSource(params *SetScriptSourceParams, conn *hc.Conn) (result *SetScriptSourceResult, err error) {
	cmd := NewSetScriptSourceCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type SetScriptSourceCB func(result *SetScriptSourceResult, err error)

// Edits JavaScript source live.

type AsyncSetScriptSourceCommand struct {
	params *SetScriptSourceParams
	cb     SetScriptSourceCB
}

func NewAsyncSetScriptSourceCommand(params *SetScriptSourceParams, cb SetScriptSourceCB) *AsyncSetScriptSourceCommand {
	return &AsyncSetScriptSourceCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetScriptSourceCommand) Name() string {
	return "Debugger.setScriptSource"
}

func (cmd *AsyncSetScriptSourceCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetScriptSourceCommand) Result() *SetScriptSourceResult {
	return &cmd.result
}

func (cmd *SetScriptSourceCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetScriptSourceCommand) Done(data []byte, err error) {
	var result SetScriptSourceResult
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

type RestartFrameParams struct {
	CallFrameId CallFrameId `json:"callFrameId"` // Call frame identifier to evaluate on.
}

type RestartFrameResult struct {
	CallFrames      []*DebuggerCallFrame `json:"callFrames"`      // New stack trace.
	AsyncStackTrace *StackTrace          `json:"asyncStackTrace"` // Async stack trace, if any.
}

// Restarts particular call frame from the beginning.

type RestartFrameCommand struct {
	params *RestartFrameParams
	result RestartFrameResult
	wg     sync.WaitGroup
	err    error
}

func NewRestartFrameCommand(params *RestartFrameParams) *RestartFrameCommand {
	return &RestartFrameCommand{
		params: params,
	}
}

func (cmd *RestartFrameCommand) Name() string {
	return "Debugger.restartFrame"
}

func (cmd *RestartFrameCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RestartFrameCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func RestartFrame(params *RestartFrameParams, conn *hc.Conn) (result *RestartFrameResult, err error) {
	cmd := NewRestartFrameCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type RestartFrameCB func(result *RestartFrameResult, err error)

// Restarts particular call frame from the beginning.

type AsyncRestartFrameCommand struct {
	params *RestartFrameParams
	cb     RestartFrameCB
}

func NewAsyncRestartFrameCommand(params *RestartFrameParams, cb RestartFrameCB) *AsyncRestartFrameCommand {
	return &AsyncRestartFrameCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncRestartFrameCommand) Name() string {
	return "Debugger.restartFrame"
}

func (cmd *AsyncRestartFrameCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RestartFrameCommand) Result() *RestartFrameResult {
	return &cmd.result
}

func (cmd *RestartFrameCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncRestartFrameCommand) Done(data []byte, err error) {
	var result RestartFrameResult
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

type GetScriptSourceParams struct {
	ScriptId *ScriptId `json:"scriptId"` // Id of the script to get source for.
}

type GetScriptSourceResult struct {
	ScriptSource string `json:"scriptSource"` // Script source.
}

// Returns source for the script with given id.

type GetScriptSourceCommand struct {
	params *GetScriptSourceParams
	result GetScriptSourceResult
	wg     sync.WaitGroup
	err    error
}

func NewGetScriptSourceCommand(params *GetScriptSourceParams) *GetScriptSourceCommand {
	return &GetScriptSourceCommand{
		params: params,
	}
}

func (cmd *GetScriptSourceCommand) Name() string {
	return "Debugger.getScriptSource"
}

func (cmd *GetScriptSourceCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetScriptSourceCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetScriptSource(params *GetScriptSourceParams, conn *hc.Conn) (result *GetScriptSourceResult, err error) {
	cmd := NewGetScriptSourceCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetScriptSourceCB func(result *GetScriptSourceResult, err error)

// Returns source for the script with given id.

type AsyncGetScriptSourceCommand struct {
	params *GetScriptSourceParams
	cb     GetScriptSourceCB
}

func NewAsyncGetScriptSourceCommand(params *GetScriptSourceParams, cb GetScriptSourceCB) *AsyncGetScriptSourceCommand {
	return &AsyncGetScriptSourceCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetScriptSourceCommand) Name() string {
	return "Debugger.getScriptSource"
}

func (cmd *AsyncGetScriptSourceCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetScriptSourceCommand) Result() *GetScriptSourceResult {
	return &cmd.result
}

func (cmd *GetScriptSourceCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetScriptSourceCommand) Done(data []byte, err error) {
	var result GetScriptSourceResult
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

type SetPauseOnExceptionsParams struct {
	State string `json:"state"` // Pause on exceptions mode.
}

// Defines pause on exceptions state. Can be set to stop on all exceptions, uncaught exceptions or no exceptions. Initial pause on exceptions state is none.

type SetPauseOnExceptionsCommand struct {
	params *SetPauseOnExceptionsParams
	wg     sync.WaitGroup
	err    error
}

func NewSetPauseOnExceptionsCommand(params *SetPauseOnExceptionsParams) *SetPauseOnExceptionsCommand {
	return &SetPauseOnExceptionsCommand{
		params: params,
	}
}

func (cmd *SetPauseOnExceptionsCommand) Name() string {
	return "Debugger.setPauseOnExceptions"
}

func (cmd *SetPauseOnExceptionsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetPauseOnExceptionsCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetPauseOnExceptions(params *SetPauseOnExceptionsParams, conn *hc.Conn) (err error) {
	cmd := NewSetPauseOnExceptionsCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetPauseOnExceptionsCB func(err error)

// Defines pause on exceptions state. Can be set to stop on all exceptions, uncaught exceptions or no exceptions. Initial pause on exceptions state is none.

type AsyncSetPauseOnExceptionsCommand struct {
	params *SetPauseOnExceptionsParams
	cb     SetPauseOnExceptionsCB
}

func NewAsyncSetPauseOnExceptionsCommand(params *SetPauseOnExceptionsParams, cb SetPauseOnExceptionsCB) *AsyncSetPauseOnExceptionsCommand {
	return &AsyncSetPauseOnExceptionsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetPauseOnExceptionsCommand) Name() string {
	return "Debugger.setPauseOnExceptions"
}

func (cmd *AsyncSetPauseOnExceptionsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetPauseOnExceptionsCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetPauseOnExceptionsCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type EvaluateOnCallFrameParams struct {
	CallFrameId           CallFrameId `json:"callFrameId"`                     // Call frame identifier to evaluate on.
	Expression            string      `json:"expression"`                      // Expression to evaluate.
	ObjectGroup           string      `json:"objectGroup,omitempty"`           // String object group name to put result into (allows rapid releasing resulting object handles using releaseObjectGroup).
	IncludeCommandLineAPI bool        `json:"includeCommandLineAPI,omitempty"` // Specifies whether command line API should be available to the evaluated expression, defaults to false.
	Silent                bool        `json:"silent,omitempty"`                // In silent mode exceptions thrown during evaluation are not reported and do not pause execution. Overrides setPauseOnException state.
	ReturnByValue         bool        `json:"returnByValue,omitempty"`         // Whether the result is expected to be a JSON object that should be sent by value.
	GeneratePreview       bool        `json:"generatePreview,omitempty"`       // Whether preview should be generated for the result.
}

type EvaluateOnCallFrameResult struct {
	Result           *RemoteObject     `json:"result"`           // Object wrapper for the evaluation result.
	ExceptionDetails *ExceptionDetails `json:"exceptionDetails"` // Exception details.
}

// Evaluates expression on a given call frame.

type EvaluateOnCallFrameCommand struct {
	params *EvaluateOnCallFrameParams
	result EvaluateOnCallFrameResult
	wg     sync.WaitGroup
	err    error
}

func NewEvaluateOnCallFrameCommand(params *EvaluateOnCallFrameParams) *EvaluateOnCallFrameCommand {
	return &EvaluateOnCallFrameCommand{
		params: params,
	}
}

func (cmd *EvaluateOnCallFrameCommand) Name() string {
	return "Debugger.evaluateOnCallFrame"
}

func (cmd *EvaluateOnCallFrameCommand) Params() interface{} {
	return cmd.params
}

func (cmd *EvaluateOnCallFrameCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func EvaluateOnCallFrame(params *EvaluateOnCallFrameParams, conn *hc.Conn) (result *EvaluateOnCallFrameResult, err error) {
	cmd := NewEvaluateOnCallFrameCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type EvaluateOnCallFrameCB func(result *EvaluateOnCallFrameResult, err error)

// Evaluates expression on a given call frame.

type AsyncEvaluateOnCallFrameCommand struct {
	params *EvaluateOnCallFrameParams
	cb     EvaluateOnCallFrameCB
}

func NewAsyncEvaluateOnCallFrameCommand(params *EvaluateOnCallFrameParams, cb EvaluateOnCallFrameCB) *AsyncEvaluateOnCallFrameCommand {
	return &AsyncEvaluateOnCallFrameCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncEvaluateOnCallFrameCommand) Name() string {
	return "Debugger.evaluateOnCallFrame"
}

func (cmd *AsyncEvaluateOnCallFrameCommand) Params() interface{} {
	return cmd.params
}

func (cmd *EvaluateOnCallFrameCommand) Result() *EvaluateOnCallFrameResult {
	return &cmd.result
}

func (cmd *EvaluateOnCallFrameCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncEvaluateOnCallFrameCommand) Done(data []byte, err error) {
	var result EvaluateOnCallFrameResult
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

type SetVariableValueParams struct {
	ScopeNumber  int           `json:"scopeNumber"`  // 0-based number of scope as was listed in scope chain. Only 'local', 'closure' and 'catch' scope types are allowed. Other scopes could be manipulated manually.
	VariableName string        `json:"variableName"` // Variable name.
	NewValue     *CallArgument `json:"newValue"`     // New variable value.
	CallFrameId  CallFrameId   `json:"callFrameId"`  // Id of callframe that holds variable.
}

// Changes value of variable in a callframe. Object-based scopes are not supported and must be mutated manually.

type SetVariableValueCommand struct {
	params *SetVariableValueParams
	wg     sync.WaitGroup
	err    error
}

func NewSetVariableValueCommand(params *SetVariableValueParams) *SetVariableValueCommand {
	return &SetVariableValueCommand{
		params: params,
	}
}

func (cmd *SetVariableValueCommand) Name() string {
	return "Debugger.setVariableValue"
}

func (cmd *SetVariableValueCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetVariableValueCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetVariableValue(params *SetVariableValueParams, conn *hc.Conn) (err error) {
	cmd := NewSetVariableValueCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetVariableValueCB func(err error)

// Changes value of variable in a callframe. Object-based scopes are not supported and must be mutated manually.

type AsyncSetVariableValueCommand struct {
	params *SetVariableValueParams
	cb     SetVariableValueCB
}

func NewAsyncSetVariableValueCommand(params *SetVariableValueParams, cb SetVariableValueCB) *AsyncSetVariableValueCommand {
	return &AsyncSetVariableValueCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetVariableValueCommand) Name() string {
	return "Debugger.setVariableValue"
}

func (cmd *AsyncSetVariableValueCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetVariableValueCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetVariableValueCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetAsyncCallStackDepthParams struct {
	MaxDepth int `json:"maxDepth"` // Maximum depth of async call stacks. Setting to 0 will effectively disable collecting async call stacks (default).
}

// Enables or disables async call stacks tracking.

type SetAsyncCallStackDepthCommand struct {
	params *SetAsyncCallStackDepthParams
	wg     sync.WaitGroup
	err    error
}

func NewSetAsyncCallStackDepthCommand(params *SetAsyncCallStackDepthParams) *SetAsyncCallStackDepthCommand {
	return &SetAsyncCallStackDepthCommand{
		params: params,
	}
}

func (cmd *SetAsyncCallStackDepthCommand) Name() string {
	return "Debugger.setAsyncCallStackDepth"
}

func (cmd *SetAsyncCallStackDepthCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetAsyncCallStackDepthCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetAsyncCallStackDepth(params *SetAsyncCallStackDepthParams, conn *hc.Conn) (err error) {
	cmd := NewSetAsyncCallStackDepthCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetAsyncCallStackDepthCB func(err error)

// Enables or disables async call stacks tracking.

type AsyncSetAsyncCallStackDepthCommand struct {
	params *SetAsyncCallStackDepthParams
	cb     SetAsyncCallStackDepthCB
}

func NewAsyncSetAsyncCallStackDepthCommand(params *SetAsyncCallStackDepthParams, cb SetAsyncCallStackDepthCB) *AsyncSetAsyncCallStackDepthCommand {
	return &AsyncSetAsyncCallStackDepthCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetAsyncCallStackDepthCommand) Name() string {
	return "Debugger.setAsyncCallStackDepth"
}

func (cmd *AsyncSetAsyncCallStackDepthCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetAsyncCallStackDepthCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetAsyncCallStackDepthCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetBlackboxPatternsParams struct {
	Patterns []string `json:"patterns"` // Array of regexps that will be used to check script url for blackbox state.
}

// Replace previous blackbox patterns with passed ones. Forces backend to skip stepping/pausing in scripts with url matching one of the patterns. VM will try to leave blackboxed script by performing 'step in' several times, finally resorting to 'step out' if unsuccessful.
// @experimental
type SetBlackboxPatternsCommand struct {
	params *SetBlackboxPatternsParams
	wg     sync.WaitGroup
	err    error
}

func NewSetBlackboxPatternsCommand(params *SetBlackboxPatternsParams) *SetBlackboxPatternsCommand {
	return &SetBlackboxPatternsCommand{
		params: params,
	}
}

func (cmd *SetBlackboxPatternsCommand) Name() string {
	return "Debugger.setBlackboxPatterns"
}

func (cmd *SetBlackboxPatternsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetBlackboxPatternsCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetBlackboxPatterns(params *SetBlackboxPatternsParams, conn *hc.Conn) (err error) {
	cmd := NewSetBlackboxPatternsCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetBlackboxPatternsCB func(err error)

// Replace previous blackbox patterns with passed ones. Forces backend to skip stepping/pausing in scripts with url matching one of the patterns. VM will try to leave blackboxed script by performing 'step in' several times, finally resorting to 'step out' if unsuccessful.
// @experimental
type AsyncSetBlackboxPatternsCommand struct {
	params *SetBlackboxPatternsParams
	cb     SetBlackboxPatternsCB
}

func NewAsyncSetBlackboxPatternsCommand(params *SetBlackboxPatternsParams, cb SetBlackboxPatternsCB) *AsyncSetBlackboxPatternsCommand {
	return &AsyncSetBlackboxPatternsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetBlackboxPatternsCommand) Name() string {
	return "Debugger.setBlackboxPatterns"
}

func (cmd *AsyncSetBlackboxPatternsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetBlackboxPatternsCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetBlackboxPatternsCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetBlackboxedRangesParams struct {
	ScriptId  *ScriptId         `json:"scriptId"` // Id of the script.
	Positions []*ScriptPosition `json:"positions"`
}

// Makes backend skip steps in the script in blackboxed ranges. VM will try leave blacklisted scripts by performing 'step in' several times, finally resorting to 'step out' if unsuccessful. Positions array contains positions where blackbox state is changed. First interval isn't blackboxed. Array should be sorted.
// @experimental
type SetBlackboxedRangesCommand struct {
	params *SetBlackboxedRangesParams
	wg     sync.WaitGroup
	err    error
}

func NewSetBlackboxedRangesCommand(params *SetBlackboxedRangesParams) *SetBlackboxedRangesCommand {
	return &SetBlackboxedRangesCommand{
		params: params,
	}
}

func (cmd *SetBlackboxedRangesCommand) Name() string {
	return "Debugger.setBlackboxedRanges"
}

func (cmd *SetBlackboxedRangesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetBlackboxedRangesCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetBlackboxedRanges(params *SetBlackboxedRangesParams, conn *hc.Conn) (err error) {
	cmd := NewSetBlackboxedRangesCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetBlackboxedRangesCB func(err error)

// Makes backend skip steps in the script in blackboxed ranges. VM will try leave blacklisted scripts by performing 'step in' several times, finally resorting to 'step out' if unsuccessful. Positions array contains positions where blackbox state is changed. First interval isn't blackboxed. Array should be sorted.
// @experimental
type AsyncSetBlackboxedRangesCommand struct {
	params *SetBlackboxedRangesParams
	cb     SetBlackboxedRangesCB
}

func NewAsyncSetBlackboxedRangesCommand(params *SetBlackboxedRangesParams, cb SetBlackboxedRangesCB) *AsyncSetBlackboxedRangesCommand {
	return &AsyncSetBlackboxedRangesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetBlackboxedRangesCommand) Name() string {
	return "Debugger.setBlackboxedRanges"
}

func (cmd *AsyncSetBlackboxedRangesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetBlackboxedRangesCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetBlackboxedRangesCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Fired when virtual machine parses script. This event is also fired for all known and uncollected scripts upon enabling debugger.

type ScriptParsedEvent struct {
	ScriptId                *ScriptId           `json:"scriptId"`                // Identifier of the script parsed.
	Url                     string              `json:"url"`                     // URL or name of the script parsed (if any).
	StartLine               int                 `json:"startLine"`               // Line offset of the script within the resource with given URL (for script tags).
	StartColumn             int                 `json:"startColumn"`             // Column offset of the script within the resource with given URL.
	EndLine                 int                 `json:"endLine"`                 // Last line of the script.
	EndColumn               int                 `json:"endColumn"`               // Length of the last line of the script.
	ExecutionContextId      *ExecutionContextId `json:"executionContextId"`      // Specifies script creation context.
	Hash                    string              `json:"hash"`                    // Content hash of the script.
	ExecutionContextAuxData map[string]string   `json:"executionContextAuxData"` // Embedder-specific auxiliary data.
	IsLiveEdit              bool                `json:"isLiveEdit"`              // True, if this script is generated as a result of the live edit operation.
	SourceMapURL            string              `json:"sourceMapURL"`            // URL of source map associated with script (if any).
	HasSourceURL            bool                `json:"hasSourceURL"`            // True, if this script has sourceURL.
}

func OnScriptParsed(conn *hc.Conn, cb func(evt *ScriptParsedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &ScriptParsedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Debugger.scriptParsed", sink)
}

// Fired when virtual machine fails to parse the script.

type ScriptFailedToParseEvent struct {
	ScriptId                *ScriptId           `json:"scriptId"`                // Identifier of the script parsed.
	Url                     string              `json:"url"`                     // URL or name of the script parsed (if any).
	StartLine               int                 `json:"startLine"`               // Line offset of the script within the resource with given URL (for script tags).
	StartColumn             int                 `json:"startColumn"`             // Column offset of the script within the resource with given URL.
	EndLine                 int                 `json:"endLine"`                 // Last line of the script.
	EndColumn               int                 `json:"endColumn"`               // Length of the last line of the script.
	ExecutionContextId      *ExecutionContextId `json:"executionContextId"`      // Specifies script creation context.
	Hash                    string              `json:"hash"`                    // Content hash of the script.
	ExecutionContextAuxData map[string]string   `json:"executionContextAuxData"` // Embedder-specific auxiliary data.
	SourceMapURL            string              `json:"sourceMapURL"`            // URL of source map associated with script (if any).
	HasSourceURL            bool                `json:"hasSourceURL"`            // True, if this script has sourceURL.
}

func OnScriptFailedToParse(conn *hc.Conn, cb func(evt *ScriptFailedToParseEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &ScriptFailedToParseEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Debugger.scriptFailedToParse", sink)
}

// Fired when breakpoint is resolved to an actual script and location.

type BreakpointResolvedEvent struct {
	BreakpointId BreakpointId `json:"breakpointId"` // Breakpoint unique identifier.
	Location     *Location    `json:"location"`     // Actual breakpoint location.
}

func OnBreakpointResolved(conn *hc.Conn, cb func(evt *BreakpointResolvedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &BreakpointResolvedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Debugger.breakpointResolved", sink)
}

// Fired when the virtual machine stopped on breakpoint or exception or any other stop criteria.

type PausedEvent struct {
	CallFrames      []*DebuggerCallFrame `json:"callFrames"`      // Call stack the virtual machine stopped on.
	Reason          string               `json:"reason"`          // Pause reason.
	Data            map[string]string    `json:"data"`            // Object containing break-specific auxiliary properties.
	HitBreakpoints  []string             `json:"hitBreakpoints"`  // Hit breakpoints IDs
	AsyncStackTrace *StackTrace          `json:"asyncStackTrace"` // Async stack trace, if any.
}

func OnPaused(conn *hc.Conn, cb func(evt *PausedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &PausedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Debugger.paused", sink)
}

// Fired when the virtual machine resumed execution.

type ResumedEvent struct {
}

func OnResumed(conn *hc.Conn, cb func(evt *ResumedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &ResumedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Debugger.resumed", sink)
}
