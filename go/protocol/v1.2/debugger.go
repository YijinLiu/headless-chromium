package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
)

// Breakpoint identifier.
type BreakpointId string

// Call frame identifier.
type CallFrameId string

// Location in the source code.
type Location struct {
	ScriptId     *ScriptId `json:"scriptId"`     // Script identifier as reported in the Debugger.scriptParsed.
	LineNumber   int       `json:"lineNumber"`   // Line number in the script (0-based).
	ColumnNumber int       `json:"columnNumber"` // Column number in the script (0-based).
}

// Location in the source code.
type ScriptPosition struct {
	LineNumber   int `json:"lineNumber"`
	ColumnNumber int `json:"columnNumber"`
}

// JavaScript call frame. Array of call frames form the call stack.
type DebuggerCallFrame struct {
	CallFrameId      CallFrameId   `json:"callFrameId"`      // Call frame identifier. This identifier is only valid while the virtual machine is paused.
	FunctionName     string        `json:"functionName"`     // Name of the JavaScript function called on this call frame.
	FunctionLocation *Location     `json:"functionLocation"` // Location in the source code.
	Location         *Location     `json:"location"`         // Location in the source code.
	ScopeChain       []*Scope      `json:"scopeChain"`       // Scope chain for this call frame.
	This             *RemoteObject `json:"this"`             // this object for this call frame.
	ReturnValue      *RemoteObject `json:"returnValue"`      // The value being returned, if the function is at return point.
}

// Scope description.
type Scope struct {
	Type          string        `json:"type"`   // Scope type.
	Object        *RemoteObject `json:"object"` // Object representing the scope. For global and with scopes it represents the actual object; for the rest of the scopes, it is artificial transient object enumerating scope variables as its properties.
	Name          string        `json:"name"`
	StartLocation *Location     `json:"startLocation"` // Location in the source code where scope starts
	EndLocation   *Location     `json:"endLocation"`   // Location in the source code where scope ends
}

// Search match for resource.
type SearchMatch struct {
	LineNumber  int    `json:"lineNumber"`  // Line number in resource content.
	LineContent string `json:"lineContent"` // Line with match content.
}

type DebuggerEnableCB func(err error)

// Enables debugger for the given page. Clients should not assume that the debugging has been enabled until the result for this command is received.
type DebuggerEnableCommand struct {
	cb DebuggerEnableCB
}

func NewDebuggerEnableCommand(cb DebuggerEnableCB) *DebuggerEnableCommand {
	return &DebuggerEnableCommand{
		cb: cb,
	}
}

func (cmd *DebuggerEnableCommand) Name() string {
	return "Debugger.enable"
}

func (cmd *DebuggerEnableCommand) Params() interface{} {
	return nil
}

func (cmd *DebuggerEnableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type DebuggerDisableCB func(err error)

// Disables debugger for given page.
type DebuggerDisableCommand struct {
	cb DebuggerDisableCB
}

func NewDebuggerDisableCommand(cb DebuggerDisableCB) *DebuggerDisableCommand {
	return &DebuggerDisableCommand{
		cb: cb,
	}
}

func (cmd *DebuggerDisableCommand) Name() string {
	return "Debugger.disable"
}

func (cmd *DebuggerDisableCommand) Params() interface{} {
	return nil
}

func (cmd *DebuggerDisableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetBreakpointsActiveParams struct {
	Active bool `json:"active"` // New value for breakpoints active state.
}

type SetBreakpointsActiveCB func(err error)

// Activates / deactivates all breakpoints on the page.
type SetBreakpointsActiveCommand struct {
	params *SetBreakpointsActiveParams
	cb     SetBreakpointsActiveCB
}

func NewSetBreakpointsActiveCommand(params *SetBreakpointsActiveParams, cb SetBreakpointsActiveCB) *SetBreakpointsActiveCommand {
	return &SetBreakpointsActiveCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetBreakpointsActiveCommand) Name() string {
	return "Debugger.setBreakpointsActive"
}

func (cmd *SetBreakpointsActiveCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetBreakpointsActiveCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetSkipAllPausesParams struct {
	Skip bool `json:"skip"` // New value for skip pauses state.
}

type SetSkipAllPausesCB func(err error)

// Makes page not interrupt on any pauses (breakpoint, exception, dom exception etc).
type SetSkipAllPausesCommand struct {
	params *SetSkipAllPausesParams
	cb     SetSkipAllPausesCB
}

func NewSetSkipAllPausesCommand(params *SetSkipAllPausesParams, cb SetSkipAllPausesCB) *SetSkipAllPausesCommand {
	return &SetSkipAllPausesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetSkipAllPausesCommand) Name() string {
	return "Debugger.setSkipAllPauses"
}

func (cmd *SetSkipAllPausesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetSkipAllPausesCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetBreakpointByUrlParams struct {
	LineNumber   int    `json:"lineNumber"`   // Line number to set breakpoint at.
	Url          string `json:"url"`          // URL of the resources to set breakpoint on.
	UrlRegex     string `json:"urlRegex"`     // Regex pattern for the URLs of the resources to set breakpoints on. Either url or urlRegex must be specified.
	ColumnNumber int    `json:"columnNumber"` // Offset in the line to set breakpoint at.
	Condition    string `json:"condition"`    // Expression to use as a breakpoint condition. When specified, debugger will only stop on the breakpoint if this expression evaluates to true.
}

type SetBreakpointByUrlResult struct {
	BreakpointId BreakpointId `json:"breakpointId"` // Id of the created breakpoint for further reference.
	Locations    []*Location  `json:"locations"`    // List of the locations this breakpoint resolved into upon addition.
}

type SetBreakpointByUrlCB func(result *SetBreakpointByUrlResult, err error)

// Sets JavaScript breakpoint at given location specified either by URL or URL regex. Once this command is issued, all existing parsed scripts will have breakpoints resolved and returned in locations property. Further matching script parsing will result in subsequent breakpointResolved events issued. This logical breakpoint will survive page reloads.
type SetBreakpointByUrlCommand struct {
	params *SetBreakpointByUrlParams
	cb     SetBreakpointByUrlCB
}

func NewSetBreakpointByUrlCommand(params *SetBreakpointByUrlParams, cb SetBreakpointByUrlCB) *SetBreakpointByUrlCommand {
	return &SetBreakpointByUrlCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetBreakpointByUrlCommand) Name() string {
	return "Debugger.setBreakpointByUrl"
}

func (cmd *SetBreakpointByUrlCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetBreakpointByUrlCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj SetBreakpointByUrlResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type SetBreakpointParams struct {
	Location  *Location `json:"location"`  // Location to set breakpoint in.
	Condition string    `json:"condition"` // Expression to use as a breakpoint condition. When specified, debugger will only stop on the breakpoint if this expression evaluates to true.
}

type SetBreakpointResult struct {
	BreakpointId   BreakpointId `json:"breakpointId"`   // Id of the created breakpoint for further reference.
	ActualLocation *Location    `json:"actualLocation"` // Location this breakpoint resolved into.
}

type SetBreakpointCB func(result *SetBreakpointResult, err error)

// Sets JavaScript breakpoint at a given location.
type SetBreakpointCommand struct {
	params *SetBreakpointParams
	cb     SetBreakpointCB
}

func NewSetBreakpointCommand(params *SetBreakpointParams, cb SetBreakpointCB) *SetBreakpointCommand {
	return &SetBreakpointCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetBreakpointCommand) Name() string {
	return "Debugger.setBreakpoint"
}

func (cmd *SetBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetBreakpointCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj SetBreakpointResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type RemoveBreakpointParams struct {
	BreakpointId BreakpointId `json:"breakpointId"`
}

type RemoveBreakpointCB func(err error)

// Removes JavaScript breakpoint.
type RemoveBreakpointCommand struct {
	params *RemoveBreakpointParams
	cb     RemoveBreakpointCB
}

func NewRemoveBreakpointCommand(params *RemoveBreakpointParams, cb RemoveBreakpointCB) *RemoveBreakpointCommand {
	return &RemoveBreakpointCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *RemoveBreakpointCommand) Name() string {
	return "Debugger.removeBreakpoint"
}

func (cmd *RemoveBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveBreakpointCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type ContinueToLocationParams struct {
	Location *Location `json:"location"` // Location to continue to.
}

type ContinueToLocationCB func(err error)

// Continues execution until specific location is reached.
type ContinueToLocationCommand struct {
	params *ContinueToLocationParams
	cb     ContinueToLocationCB
}

func NewContinueToLocationCommand(params *ContinueToLocationParams, cb ContinueToLocationCB) *ContinueToLocationCommand {
	return &ContinueToLocationCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *ContinueToLocationCommand) Name() string {
	return "Debugger.continueToLocation"
}

func (cmd *ContinueToLocationCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ContinueToLocationCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type StepOverCB func(err error)

// Steps over the statement.
type StepOverCommand struct {
	cb StepOverCB
}

func NewStepOverCommand(cb StepOverCB) *StepOverCommand {
	return &StepOverCommand{
		cb: cb,
	}
}

func (cmd *StepOverCommand) Name() string {
	return "Debugger.stepOver"
}

func (cmd *StepOverCommand) Params() interface{} {
	return nil
}

func (cmd *StepOverCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type StepIntoCB func(err error)

// Steps into the function call.
type StepIntoCommand struct {
	cb StepIntoCB
}

func NewStepIntoCommand(cb StepIntoCB) *StepIntoCommand {
	return &StepIntoCommand{
		cb: cb,
	}
}

func (cmd *StepIntoCommand) Name() string {
	return "Debugger.stepInto"
}

func (cmd *StepIntoCommand) Params() interface{} {
	return nil
}

func (cmd *StepIntoCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type StepOutCB func(err error)

// Steps out of the function call.
type StepOutCommand struct {
	cb StepOutCB
}

func NewStepOutCommand(cb StepOutCB) *StepOutCommand {
	return &StepOutCommand{
		cb: cb,
	}
}

func (cmd *StepOutCommand) Name() string {
	return "Debugger.stepOut"
}

func (cmd *StepOutCommand) Params() interface{} {
	return nil
}

func (cmd *StepOutCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type PauseCB func(err error)

// Stops on the next JavaScript statement.
type PauseCommand struct {
	cb PauseCB
}

func NewPauseCommand(cb PauseCB) *PauseCommand {
	return &PauseCommand{
		cb: cb,
	}
}

func (cmd *PauseCommand) Name() string {
	return "Debugger.pause"
}

func (cmd *PauseCommand) Params() interface{} {
	return nil
}

func (cmd *PauseCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type ResumeCB func(err error)

// Resumes JavaScript execution.
type ResumeCommand struct {
	cb ResumeCB
}

func NewResumeCommand(cb ResumeCB) *ResumeCommand {
	return &ResumeCommand{
		cb: cb,
	}
}

func (cmd *ResumeCommand) Name() string {
	return "Debugger.resume"
}

func (cmd *ResumeCommand) Params() interface{} {
	return nil
}

func (cmd *ResumeCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SearchInContentParams struct {
	ScriptId      *ScriptId `json:"scriptId"`      // Id of the script to search in.
	Query         string    `json:"query"`         // String to search for.
	CaseSensitive bool      `json:"caseSensitive"` // If true, search is case sensitive.
	IsRegex       bool      `json:"isRegex"`       // If true, treats string parameter as regex.
}

type SearchInContentResult struct {
	Result []*SearchMatch `json:"result"` // List of search matches.
}

type SearchInContentCB func(result *SearchInContentResult, err error)

// Searches for given string in script content.
type SearchInContentCommand struct {
	params *SearchInContentParams
	cb     SearchInContentCB
}

func NewSearchInContentCommand(params *SearchInContentParams, cb SearchInContentCB) *SearchInContentCommand {
	return &SearchInContentCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SearchInContentCommand) Name() string {
	return "Debugger.searchInContent"
}

func (cmd *SearchInContentCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SearchInContentCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj SearchInContentResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type SetScriptSourceParams struct {
	ScriptId     *ScriptId `json:"scriptId"`     // Id of the script to edit.
	ScriptSource string    `json:"scriptSource"` // New content of the script.
	DryRun       bool      `json:"dryRun"`       //  If true the change will not actually be applied. Dry run may be used to get result description without actually modifying the code.
}

type SetScriptSourceResult struct {
	CallFrames       []*DebuggerCallFrame `json:"callFrames"`       // New stack trace in case editing has happened while VM was stopped.
	StackChanged     bool                 `json:"stackChanged"`     // Whether current call stack  was modified after applying the changes.
	AsyncStackTrace  *StackTrace          `json:"asyncStackTrace"`  // Async stack trace, if any.
	ExceptionDetails *ExceptionDetails    `json:"exceptionDetails"` // Exception details if any.
}

type SetScriptSourceCB func(result *SetScriptSourceResult, err error)

// Edits JavaScript source live.
type SetScriptSourceCommand struct {
	params *SetScriptSourceParams
	cb     SetScriptSourceCB
}

func NewSetScriptSourceCommand(params *SetScriptSourceParams, cb SetScriptSourceCB) *SetScriptSourceCommand {
	return &SetScriptSourceCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetScriptSourceCommand) Name() string {
	return "Debugger.setScriptSource"
}

func (cmd *SetScriptSourceCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetScriptSourceCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj SetScriptSourceResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type RestartFrameParams struct {
	CallFrameId CallFrameId `json:"callFrameId"` // Call frame identifier to evaluate on.
}

type RestartFrameResult struct {
	CallFrames      []*DebuggerCallFrame `json:"callFrames"`      // New stack trace.
	AsyncStackTrace *StackTrace          `json:"asyncStackTrace"` // Async stack trace, if any.
}

type RestartFrameCB func(result *RestartFrameResult, err error)

// Restarts particular call frame from the beginning.
type RestartFrameCommand struct {
	params *RestartFrameParams
	cb     RestartFrameCB
}

func NewRestartFrameCommand(params *RestartFrameParams, cb RestartFrameCB) *RestartFrameCommand {
	return &RestartFrameCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *RestartFrameCommand) Name() string {
	return "Debugger.restartFrame"
}

func (cmd *RestartFrameCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RestartFrameCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj RestartFrameResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type GetScriptSourceParams struct {
	ScriptId *ScriptId `json:"scriptId"` // Id of the script to get source for.
}

type GetScriptSourceResult struct {
	ScriptSource string `json:"scriptSource"` // Script source.
}

type GetScriptSourceCB func(result *GetScriptSourceResult, err error)

// Returns source for the script with given id.
type GetScriptSourceCommand struct {
	params *GetScriptSourceParams
	cb     GetScriptSourceCB
}

func NewGetScriptSourceCommand(params *GetScriptSourceParams, cb GetScriptSourceCB) *GetScriptSourceCommand {
	return &GetScriptSourceCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetScriptSourceCommand) Name() string {
	return "Debugger.getScriptSource"
}

func (cmd *GetScriptSourceCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetScriptSourceCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetScriptSourceResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type SetPauseOnExceptionsParams struct {
	State string `json:"state"` // Pause on exceptions mode.
}

type SetPauseOnExceptionsCB func(err error)

// Defines pause on exceptions state. Can be set to stop on all exceptions, uncaught exceptions or no exceptions. Initial pause on exceptions state is none.
type SetPauseOnExceptionsCommand struct {
	params *SetPauseOnExceptionsParams
	cb     SetPauseOnExceptionsCB
}

func NewSetPauseOnExceptionsCommand(params *SetPauseOnExceptionsParams, cb SetPauseOnExceptionsCB) *SetPauseOnExceptionsCommand {
	return &SetPauseOnExceptionsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetPauseOnExceptionsCommand) Name() string {
	return "Debugger.setPauseOnExceptions"
}

func (cmd *SetPauseOnExceptionsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetPauseOnExceptionsCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type EvaluateOnCallFrameParams struct {
	CallFrameId           CallFrameId `json:"callFrameId"`           // Call frame identifier to evaluate on.
	Expression            string      `json:"expression"`            // Expression to evaluate.
	ObjectGroup           string      `json:"objectGroup"`           // String object group name to put result into (allows rapid releasing resulting object handles using releaseObjectGroup).
	IncludeCommandLineAPI bool        `json:"includeCommandLineAPI"` // Specifies whether command line API should be available to the evaluated expression, defaults to false.
	Silent                bool        `json:"silent"`                // In silent mode exceptions thrown during evaluation are not reported and do not pause execution. Overrides setPauseOnException state.
	ReturnByValue         bool        `json:"returnByValue"`         // Whether the result is expected to be a JSON object that should be sent by value.
	GeneratePreview       bool        `json:"generatePreview"`       // Whether preview should be generated for the result.
}

type EvaluateOnCallFrameResult struct {
	Result           *RemoteObject     `json:"result"`           // Object wrapper for the evaluation result.
	ExceptionDetails *ExceptionDetails `json:"exceptionDetails"` // Exception details.
}

type EvaluateOnCallFrameCB func(result *EvaluateOnCallFrameResult, err error)

// Evaluates expression on a given call frame.
type EvaluateOnCallFrameCommand struct {
	params *EvaluateOnCallFrameParams
	cb     EvaluateOnCallFrameCB
}

func NewEvaluateOnCallFrameCommand(params *EvaluateOnCallFrameParams, cb EvaluateOnCallFrameCB) *EvaluateOnCallFrameCommand {
	return &EvaluateOnCallFrameCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *EvaluateOnCallFrameCommand) Name() string {
	return "Debugger.evaluateOnCallFrame"
}

func (cmd *EvaluateOnCallFrameCommand) Params() interface{} {
	return cmd.params
}

func (cmd *EvaluateOnCallFrameCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj EvaluateOnCallFrameResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type SetVariableValueParams struct {
	ScopeNumber  int           `json:"scopeNumber"`  // 0-based number of scope as was listed in scope chain. Only 'local', 'closure' and 'catch' scope types are allowed. Other scopes could be manipulated manually.
	VariableName string        `json:"variableName"` // Variable name.
	NewValue     *CallArgument `json:"newValue"`     // New variable value.
	CallFrameId  CallFrameId   `json:"callFrameId"`  // Id of callframe that holds variable.
}

type SetVariableValueCB func(err error)

// Changes value of variable in a callframe. Object-based scopes are not supported and must be mutated manually.
type SetVariableValueCommand struct {
	params *SetVariableValueParams
	cb     SetVariableValueCB
}

func NewSetVariableValueCommand(params *SetVariableValueParams, cb SetVariableValueCB) *SetVariableValueCommand {
	return &SetVariableValueCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetVariableValueCommand) Name() string {
	return "Debugger.setVariableValue"
}

func (cmd *SetVariableValueCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetVariableValueCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetAsyncCallStackDepthParams struct {
	MaxDepth int `json:"maxDepth"` // Maximum depth of async call stacks. Setting to 0 will effectively disable collecting async call stacks (default).
}

type SetAsyncCallStackDepthCB func(err error)

// Enables or disables async call stacks tracking.
type SetAsyncCallStackDepthCommand struct {
	params *SetAsyncCallStackDepthParams
	cb     SetAsyncCallStackDepthCB
}

func NewSetAsyncCallStackDepthCommand(params *SetAsyncCallStackDepthParams, cb SetAsyncCallStackDepthCB) *SetAsyncCallStackDepthCommand {
	return &SetAsyncCallStackDepthCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetAsyncCallStackDepthCommand) Name() string {
	return "Debugger.setAsyncCallStackDepth"
}

func (cmd *SetAsyncCallStackDepthCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetAsyncCallStackDepthCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetBlackboxPatternsParams struct {
	Patterns []string `json:"patterns"` // Array of regexps that will be used to check script url for blackbox state.
}

type SetBlackboxPatternsCB func(err error)

// Replace previous blackbox patterns with passed ones. Forces backend to skip stepping/pausing in scripts with url matching one of the patterns. VM will try to leave blackboxed script by performing 'step in' several times, finally resorting to 'step out' if unsuccessful.
type SetBlackboxPatternsCommand struct {
	params *SetBlackboxPatternsParams
	cb     SetBlackboxPatternsCB
}

func NewSetBlackboxPatternsCommand(params *SetBlackboxPatternsParams, cb SetBlackboxPatternsCB) *SetBlackboxPatternsCommand {
	return &SetBlackboxPatternsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetBlackboxPatternsCommand) Name() string {
	return "Debugger.setBlackboxPatterns"
}

func (cmd *SetBlackboxPatternsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetBlackboxPatternsCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetBlackboxedRangesParams struct {
	ScriptId  *ScriptId         `json:"scriptId"` // Id of the script.
	Positions []*ScriptPosition `json:"positions"`
}

type SetBlackboxedRangesCB func(err error)

// Makes backend skip steps in the script in blackboxed ranges. VM will try leave blacklisted scripts by performing 'step in' several times, finally resorting to 'step out' if unsuccessful. Positions array contains positions where blackbox state is changed. First interval isn't blackboxed. Array should be sorted.
type SetBlackboxedRangesCommand struct {
	params *SetBlackboxedRangesParams
	cb     SetBlackboxedRangesCB
}

func NewSetBlackboxedRangesCommand(params *SetBlackboxedRangesParams, cb SetBlackboxedRangesCB) *SetBlackboxedRangesCommand {
	return &SetBlackboxedRangesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetBlackboxedRangesCommand) Name() string {
	return "Debugger.setBlackboxedRanges"
}

func (cmd *SetBlackboxedRangesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetBlackboxedRangesCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

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

// Fired when virtual machine parses script. This event is also fired for all known and uncollected scripts upon enabling debugger.
type ScriptParsedEventSink struct {
	events chan *ScriptParsedEvent
}

func NewScriptParsedEventSink(bufSize int) *ScriptParsedEventSink {
	return &ScriptParsedEventSink{
		events: make(chan *ScriptParsedEvent, bufSize),
	}
}

func (s *ScriptParsedEventSink) Name() string {
	return "Debugger.scriptParsed"
}

func (s *ScriptParsedEventSink) OnEvent(params []byte) {
	evt := &ScriptParsedEvent{}
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

// Fired when virtual machine fails to parse the script.
type ScriptFailedToParseEventSink struct {
	events chan *ScriptFailedToParseEvent
}

func NewScriptFailedToParseEventSink(bufSize int) *ScriptFailedToParseEventSink {
	return &ScriptFailedToParseEventSink{
		events: make(chan *ScriptFailedToParseEvent, bufSize),
	}
}

func (s *ScriptFailedToParseEventSink) Name() string {
	return "Debugger.scriptFailedToParse"
}

func (s *ScriptFailedToParseEventSink) OnEvent(params []byte) {
	evt := &ScriptFailedToParseEvent{}
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

type BreakpointResolvedEvent struct {
	BreakpointId BreakpointId `json:"breakpointId"` // Breakpoint unique identifier.
	Location     *Location    `json:"location"`     // Actual breakpoint location.
}

// Fired when breakpoint is resolved to an actual script and location.
type BreakpointResolvedEventSink struct {
	events chan *BreakpointResolvedEvent
}

func NewBreakpointResolvedEventSink(bufSize int) *BreakpointResolvedEventSink {
	return &BreakpointResolvedEventSink{
		events: make(chan *BreakpointResolvedEvent, bufSize),
	}
}

func (s *BreakpointResolvedEventSink) Name() string {
	return "Debugger.breakpointResolved"
}

func (s *BreakpointResolvedEventSink) OnEvent(params []byte) {
	evt := &BreakpointResolvedEvent{}
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

type PausedEvent struct {
	CallFrames      []*DebuggerCallFrame `json:"callFrames"`      // Call stack the virtual machine stopped on.
	Reason          string               `json:"reason"`          // Pause reason.
	Data            map[string]string    `json:"data"`            // Object containing break-specific auxiliary properties.
	HitBreakpoints  []string             `json:"hitBreakpoints"`  // Hit breakpoints IDs
	AsyncStackTrace *StackTrace          `json:"asyncStackTrace"` // Async stack trace, if any.
}

// Fired when the virtual machine stopped on breakpoint or exception or any other stop criteria.
type PausedEventSink struct {
	events chan *PausedEvent
}

func NewPausedEventSink(bufSize int) *PausedEventSink {
	return &PausedEventSink{
		events: make(chan *PausedEvent, bufSize),
	}
}

func (s *PausedEventSink) Name() string {
	return "Debugger.paused"
}

func (s *PausedEventSink) OnEvent(params []byte) {
	evt := &PausedEvent{}
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

type ResumedEvent struct {
}

// Fired when the virtual machine resumed execution.
type ResumedEventSink struct {
	events chan *ResumedEvent
}

func NewResumedEventSink(bufSize int) *ResumedEventSink {
	return &ResumedEventSink{
		events: make(chan *ResumedEvent, bufSize),
	}
}

func (s *ResumedEventSink) Name() string {
	return "Debugger.resumed"
}

func (s *ResumedEventSink) OnEvent(params []byte) {
	evt := &ResumedEvent{}
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
