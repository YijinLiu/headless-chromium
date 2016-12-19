package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
)

// Unique script identifier.
type ScriptId string

// Unique object identifier.
type RemoteObjectId string

// Primitive value which cannot be JSON-stringified.
type UnserializableValue string

const UnserializableValueInfinity UnserializableValue = "Infinity"
const UnserializableValueNaN UnserializableValue = "NaN"
const UnserializableValueNInfinity UnserializableValue = "-Infinity"
const UnserializableValueN0 UnserializableValue = "-0"

// Mirror object referencing original JavaScript object.
type RemoteObject struct {
	Type                string              `json:"type"`                // Object type.
	Subtype             string              `json:"subtype"`             // Object subtype hint. Specified for object type values only.
	ClassName           string              `json:"className"`           // Object class (constructor) name. Specified for object type values only.
	Value               string              `json:"value"`               // Remote object value in case of primitive values or JSON values (if it was requested).
	UnserializableValue UnserializableValue `json:"unserializableValue"` // Primitive value which can not be JSON-stringified does not have value, but gets this property.
	Description         string              `json:"description"`         // String representation of the object.
	ObjectId            RemoteObjectId      `json:"objectId"`            // Unique object identifier (for non-primitive values).
	Preview             *ObjectPreview      `json:"preview"`             // Preview containing abbreviated property values. Specified for object type values only.
	CustomPreview       *CustomPreview      `json:"customPreview"`
}

type CustomPreview struct {
	Header                     string         `json:"header"`
	HasBody                    bool           `json:"hasBody"`
	FormatterObjectId          RemoteObjectId `json:"formatterObjectId"`
	BindRemoteObjectFunctionId RemoteObjectId `json:"bindRemoteObjectFunctionId"`
	ConfigObjectId             RemoteObjectId `json:"configObjectId"`
}

// Object containing abbreviated remote object value.
type ObjectPreview struct {
	Type        string             `json:"type"`        // Object type.
	Subtype     string             `json:"subtype"`     // Object subtype hint. Specified for object type values only.
	Description string             `json:"description"` // String representation of the object.
	Overflow    bool               `json:"overflow"`    // True iff some of the properties or entries of the original object did not fit.
	Properties  []*PropertyPreview `json:"properties"`  // List of the properties.
	Entries     []*EntryPreview    `json:"entries"`     // List of the entries. Specified for map and set subtype values only.
}

type PropertyPreview struct {
	Name         string         `json:"name"`         // Property name.
	Type         string         `json:"type"`         // Object type. Accessor means that the property itself is an accessor property.
	Value        string         `json:"value"`        // User-friendly property value string.
	ValuePreview *ObjectPreview `json:"valuePreview"` // Nested value preview.
	Subtype      string         `json:"subtype"`      // Object subtype hint. Specified for object type values only.
}

type EntryPreview struct {
	Key   *ObjectPreview `json:"key"`   // Preview of the key. Specified for map-like collection entries.
	Value *ObjectPreview `json:"value"` // Preview of the value.
}

// Object property descriptor.
type PropertyDescriptor struct {
	Name         string        `json:"name"`         // Property name or symbol description.
	Value        *RemoteObject `json:"value"`        // The value associated with the property.
	Writable     bool          `json:"writable"`     // True if the value associated with the property may be changed (data descriptors only).
	Get          *RemoteObject `json:"get"`          // A function which serves as a getter for the property, or undefined if there is no getter (accessor descriptors only).
	Set          *RemoteObject `json:"set"`          // A function which serves as a setter for the property, or undefined if there is no setter (accessor descriptors only).
	Configurable bool          `json:"configurable"` // True if the type of this property descriptor may be changed and if the property may be deleted from the corresponding object.
	Enumerable   bool          `json:"enumerable"`   // True if this property shows up during enumeration of the properties on the corresponding object.
	WasThrown    bool          `json:"wasThrown"`    // True if the result was thrown during the evaluation.
	IsOwn        bool          `json:"isOwn"`        // True if the property is owned for the object.
	Symbol       *RemoteObject `json:"symbol"`       // Property symbol object, if the property is of the symbol type.
}

// Object internal property descriptor. This property isn't normally visible in JavaScript code.
type InternalPropertyDescriptor struct {
	Name  string        `json:"name"`  // Conventional property name.
	Value *RemoteObject `json:"value"` // The value associated with the property.
}

// Represents function call argument. Either remote object id objectId, primitive value, unserializable primitive value or neither of (for undefined) them should be specified.
type CallArgument struct {
	Value               string              `json:"value"`               // Primitive value.
	UnserializableValue UnserializableValue `json:"unserializableValue"` // Primitive value which can not be JSON-stringified.
	ObjectId            RemoteObjectId      `json:"objectId"`            // Remote object handle.
}

// Id of an execution context.
type ExecutionContextId int

// Description of an isolated world.
type ExecutionContextDescription struct {
	Id      ExecutionContextId `json:"id"`      // Unique id of the execution context. It can be used to specify in which execution context script evaluation should be performed.
	Origin  string             `json:"origin"`  // Execution context origin.
	Name    string             `json:"name"`    // Human readable name describing given context.
	AuxData map[string]string  `json:"auxData"` // Embedder-specific auxiliary data.
}

// Detailed information about exception (or error) that was thrown during script compilation or execution.
type ExceptionDetails struct {
	ExceptionId        int                `json:"exceptionId"`        // Exception id.
	Text               string             `json:"text"`               // Exception text, which should be used together with exception object when available.
	LineNumber         int                `json:"lineNumber"`         // Line number of the exception location (0-based).
	ColumnNumber       int                `json:"columnNumber"`       // Column number of the exception location (0-based).
	ScriptId           ScriptId           `json:"scriptId"`           // Script ID of the exception location.
	Url                string             `json:"url"`                // URL of the exception location, to be used when the script was not reported.
	StackTrace         *StackTrace        `json:"stackTrace"`         // JavaScript stack trace if available.
	Exception          *RemoteObject      `json:"exception"`          // Exception object if available.
	ExecutionContextId ExecutionContextId `json:"executionContextId"` // Identifier of the context where exception happened.
}

// Number of milliseconds since epoch.
type RuntimeTimestamp int

// Stack entry for runtime errors and assertions.
type RuntimeCallFrame struct {
	FunctionName string   `json:"functionName"` // JavaScript function name.
	ScriptId     ScriptId `json:"scriptId"`     // JavaScript script id.
	Url          string   `json:"url"`          // JavaScript script name or url.
	LineNumber   int      `json:"lineNumber"`   // JavaScript script line number (0-based).
	ColumnNumber int      `json:"columnNumber"` // JavaScript script column number (0-based).
}

// Call frames for assertions or error messages.
type StackTrace struct {
	Description string              `json:"description"` // String label of this stack trace. For async traces this may be a name of the function that initiated the async call.
	CallFrames  []*RuntimeCallFrame `json:"callFrames"`  // JavaScript function name.
	Parent      *StackTrace         `json:"parent"`      // Asynchronous JavaScript stack trace that preceded this stack, if available.
}

type EvaluateParams struct {
	Expression            string             `json:"expression"`            // Expression to evaluate.
	ObjectGroup           string             `json:"objectGroup"`           // Symbolic group name that can be used to release multiple objects.
	IncludeCommandLineAPI bool               `json:"includeCommandLineAPI"` // Determines whether Command Line API should be available during the evaluation.
	Silent                bool               `json:"silent"`                // In silent mode exceptions thrown during evaluation are not reported and do not pause execution. Overrides setPauseOnException state.
	ContextId             ExecutionContextId `json:"contextId"`             // Specifies in which execution context to perform evaluation. If the parameter is omitted the evaluation will be performed in the context of the inspected page.
	ReturnByValue         bool               `json:"returnByValue"`         // Whether the result is expected to be a JSON object that should be sent by value.
	GeneratePreview       bool               `json:"generatePreview"`       // Whether preview should be generated for the result.
	UserGesture           bool               `json:"userGesture"`           // Whether execution should be treated as initiated by user in the UI.
	AwaitPromise          bool               `json:"awaitPromise"`          // Whether execution should wait for promise to be resolved. If the result of evaluation is not a Promise, it's considered to be an error.
}

type EvaluateResult struct {
	Result           *RemoteObject     `json:"result"`           // Evaluation result.
	ExceptionDetails *ExceptionDetails `json:"exceptionDetails"` // Exception details.
}

type EvaluateCB func(result *EvaluateResult, err error)

// Evaluates expression on global object.
type EvaluateCommand struct {
	params *EvaluateParams
	cb     EvaluateCB
}

func NewEvaluateCommand(params *EvaluateParams, cb EvaluateCB) *EvaluateCommand {
	return &EvaluateCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *EvaluateCommand) Name() string {
	return "Runtime.evaluate"
}

func (cmd *EvaluateCommand) Params() interface{} {
	return cmd.params
}

func (cmd *EvaluateCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj EvaluateResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type AwaitPromiseParams struct {
	PromiseObjectId RemoteObjectId `json:"promiseObjectId"` // Identifier of the promise.
	ReturnByValue   bool           `json:"returnByValue"`   // Whether the result is expected to be a JSON object that should be sent by value.
	GeneratePreview bool           `json:"generatePreview"` // Whether preview should be generated for the result.
}

type AwaitPromiseResult struct {
	Result           *RemoteObject     `json:"result"`           // Promise result. Will contain rejected value if promise was rejected.
	ExceptionDetails *ExceptionDetails `json:"exceptionDetails"` // Exception details if stack strace is available.
}

type AwaitPromiseCB func(result *AwaitPromiseResult, err error)

// Add handler to promise with given promise object id.
type AwaitPromiseCommand struct {
	params *AwaitPromiseParams
	cb     AwaitPromiseCB
}

func NewAwaitPromiseCommand(params *AwaitPromiseParams, cb AwaitPromiseCB) *AwaitPromiseCommand {
	return &AwaitPromiseCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AwaitPromiseCommand) Name() string {
	return "Runtime.awaitPromise"
}

func (cmd *AwaitPromiseCommand) Params() interface{} {
	return cmd.params
}

func (cmd *AwaitPromiseCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj AwaitPromiseResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type CallFunctionOnParams struct {
	ObjectId            RemoteObjectId  `json:"objectId"`            // Identifier of the object to call function on.
	FunctionDeclaration string          `json:"functionDeclaration"` // Declaration of the function to call.
	Arguments           []*CallArgument `json:"arguments"`           // Call arguments. All call arguments must belong to the same JavaScript world as the target object.
	Silent              bool            `json:"silent"`              // In silent mode exceptions thrown during evaluation are not reported and do not pause execution. Overrides setPauseOnException state.
	ReturnByValue       bool            `json:"returnByValue"`       // Whether the result is expected to be a JSON object which should be sent by value.
	GeneratePreview     bool            `json:"generatePreview"`     // Whether preview should be generated for the result.
	UserGesture         bool            `json:"userGesture"`         // Whether execution should be treated as initiated by user in the UI.
	AwaitPromise        bool            `json:"awaitPromise"`        // Whether execution should wait for promise to be resolved. If the result of evaluation is not a Promise, it's considered to be an error.
}

type CallFunctionOnResult struct {
	Result           *RemoteObject     `json:"result"`           // Call result.
	ExceptionDetails *ExceptionDetails `json:"exceptionDetails"` // Exception details.
}

type CallFunctionOnCB func(result *CallFunctionOnResult, err error)

// Calls function with given declaration on the given object. Object group of the result is inherited from the target object.
type CallFunctionOnCommand struct {
	params *CallFunctionOnParams
	cb     CallFunctionOnCB
}

func NewCallFunctionOnCommand(params *CallFunctionOnParams, cb CallFunctionOnCB) *CallFunctionOnCommand {
	return &CallFunctionOnCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *CallFunctionOnCommand) Name() string {
	return "Runtime.callFunctionOn"
}

func (cmd *CallFunctionOnCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CallFunctionOnCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj CallFunctionOnResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type GetPropertiesParams struct {
	ObjectId               RemoteObjectId `json:"objectId"`               // Identifier of the object to return properties for.
	OwnProperties          bool           `json:"ownProperties"`          // If true, returns properties belonging only to the element itself, not to its prototype chain.
	AccessorPropertiesOnly bool           `json:"accessorPropertiesOnly"` // If true, returns accessor properties (with getter/setter) only; internal properties are not returned either.
	GeneratePreview        bool           `json:"generatePreview"`        // Whether preview should be generated for the results.
}

type GetPropertiesResult struct {
	Result             []*PropertyDescriptor         `json:"result"`             // Object properties.
	InternalProperties []*InternalPropertyDescriptor `json:"internalProperties"` // Internal object properties (only of the element itself).
	ExceptionDetails   *ExceptionDetails             `json:"exceptionDetails"`   // Exception details.
}

type GetPropertiesCB func(result *GetPropertiesResult, err error)

// Returns properties of a given object. Object group of the result is inherited from the target object.
type GetPropertiesCommand struct {
	params *GetPropertiesParams
	cb     GetPropertiesCB
}

func NewGetPropertiesCommand(params *GetPropertiesParams, cb GetPropertiesCB) *GetPropertiesCommand {
	return &GetPropertiesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetPropertiesCommand) Name() string {
	return "Runtime.getProperties"
}

func (cmd *GetPropertiesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetPropertiesCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetPropertiesResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type ReleaseObjectParams struct {
	ObjectId RemoteObjectId `json:"objectId"` // Identifier of the object to release.
}

type ReleaseObjectCB func(err error)

// Releases remote object with given id.
type ReleaseObjectCommand struct {
	params *ReleaseObjectParams
	cb     ReleaseObjectCB
}

func NewReleaseObjectCommand(params *ReleaseObjectParams, cb ReleaseObjectCB) *ReleaseObjectCommand {
	return &ReleaseObjectCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *ReleaseObjectCommand) Name() string {
	return "Runtime.releaseObject"
}

func (cmd *ReleaseObjectCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ReleaseObjectCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type ReleaseObjectGroupParams struct {
	ObjectGroup string `json:"objectGroup"` // Symbolic object group name.
}

type ReleaseObjectGroupCB func(err error)

// Releases all remote objects that belong to a given group.
type ReleaseObjectGroupCommand struct {
	params *ReleaseObjectGroupParams
	cb     ReleaseObjectGroupCB
}

func NewReleaseObjectGroupCommand(params *ReleaseObjectGroupParams, cb ReleaseObjectGroupCB) *ReleaseObjectGroupCommand {
	return &ReleaseObjectGroupCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *ReleaseObjectGroupCommand) Name() string {
	return "Runtime.releaseObjectGroup"
}

func (cmd *ReleaseObjectGroupCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ReleaseObjectGroupCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type RunIfWaitingForDebuggerCB func(err error)

// Tells inspected instance to run if it was waiting for debugger to attach.
type RunIfWaitingForDebuggerCommand struct {
	cb RunIfWaitingForDebuggerCB
}

func NewRunIfWaitingForDebuggerCommand(cb RunIfWaitingForDebuggerCB) *RunIfWaitingForDebuggerCommand {
	return &RunIfWaitingForDebuggerCommand{
		cb: cb,
	}
}

func (cmd *RunIfWaitingForDebuggerCommand) Name() string {
	return "Runtime.runIfWaitingForDebugger"
}

func (cmd *RunIfWaitingForDebuggerCommand) Params() interface{} {
	return nil
}

func (cmd *RunIfWaitingForDebuggerCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type RuntimeEnableCB func(err error)

// Enables reporting of execution contexts creation by means of executionContextCreated event. When the reporting gets enabled the event will be sent immediately for each existing execution context.
type RuntimeEnableCommand struct {
	cb RuntimeEnableCB
}

func NewRuntimeEnableCommand(cb RuntimeEnableCB) *RuntimeEnableCommand {
	return &RuntimeEnableCommand{
		cb: cb,
	}
}

func (cmd *RuntimeEnableCommand) Name() string {
	return "Runtime.enable"
}

func (cmd *RuntimeEnableCommand) Params() interface{} {
	return nil
}

func (cmd *RuntimeEnableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type RuntimeDisableCB func(err error)

// Disables reporting of execution contexts creation.
type RuntimeDisableCommand struct {
	cb RuntimeDisableCB
}

func NewRuntimeDisableCommand(cb RuntimeDisableCB) *RuntimeDisableCommand {
	return &RuntimeDisableCommand{
		cb: cb,
	}
}

func (cmd *RuntimeDisableCommand) Name() string {
	return "Runtime.disable"
}

func (cmd *RuntimeDisableCommand) Params() interface{} {
	return nil
}

func (cmd *RuntimeDisableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type DiscardConsoleEntriesCB func(err error)

// Discards collected exceptions and console API calls.
type DiscardConsoleEntriesCommand struct {
	cb DiscardConsoleEntriesCB
}

func NewDiscardConsoleEntriesCommand(cb DiscardConsoleEntriesCB) *DiscardConsoleEntriesCommand {
	return &DiscardConsoleEntriesCommand{
		cb: cb,
	}
}

func (cmd *DiscardConsoleEntriesCommand) Name() string {
	return "Runtime.discardConsoleEntries"
}

func (cmd *DiscardConsoleEntriesCommand) Params() interface{} {
	return nil
}

func (cmd *DiscardConsoleEntriesCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetCustomObjectFormatterEnabledParams struct {
	Enabled bool `json:"enabled"`
}

type SetCustomObjectFormatterEnabledCB func(err error)

type SetCustomObjectFormatterEnabledCommand struct {
	params *SetCustomObjectFormatterEnabledParams
	cb     SetCustomObjectFormatterEnabledCB
}

func NewSetCustomObjectFormatterEnabledCommand(params *SetCustomObjectFormatterEnabledParams, cb SetCustomObjectFormatterEnabledCB) *SetCustomObjectFormatterEnabledCommand {
	return &SetCustomObjectFormatterEnabledCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetCustomObjectFormatterEnabledCommand) Name() string {
	return "Runtime.setCustomObjectFormatterEnabled"
}

func (cmd *SetCustomObjectFormatterEnabledCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetCustomObjectFormatterEnabledCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type CompileScriptParams struct {
	Expression         string             `json:"expression"`         // Expression to compile.
	SourceURL          string             `json:"sourceURL"`          // Source url to be set for the script.
	PersistScript      bool               `json:"persistScript"`      // Specifies whether the compiled script should be persisted.
	ExecutionContextId ExecutionContextId `json:"executionContextId"` // Specifies in which execution context to perform script run. If the parameter is omitted the evaluation will be performed in the context of the inspected page.
}

type CompileScriptResult struct {
	ScriptId         ScriptId          `json:"scriptId"`         // Id of the script.
	ExceptionDetails *ExceptionDetails `json:"exceptionDetails"` // Exception details.
}

type CompileScriptCB func(result *CompileScriptResult, err error)

// Compiles expression.
type CompileScriptCommand struct {
	params *CompileScriptParams
	cb     CompileScriptCB
}

func NewCompileScriptCommand(params *CompileScriptParams, cb CompileScriptCB) *CompileScriptCommand {
	return &CompileScriptCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *CompileScriptCommand) Name() string {
	return "Runtime.compileScript"
}

func (cmd *CompileScriptCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CompileScriptCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj CompileScriptResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type RunScriptParams struct {
	ScriptId              ScriptId           `json:"scriptId"`              // Id of the script to run.
	ExecutionContextId    ExecutionContextId `json:"executionContextId"`    // Specifies in which execution context to perform script run. If the parameter is omitted the evaluation will be performed in the context of the inspected page.
	ObjectGroup           string             `json:"objectGroup"`           // Symbolic group name that can be used to release multiple objects.
	Silent                bool               `json:"silent"`                // In silent mode exceptions thrown during evaluation are not reported and do not pause execution. Overrides setPauseOnException state.
	IncludeCommandLineAPI bool               `json:"includeCommandLineAPI"` // Determines whether Command Line API should be available during the evaluation.
	ReturnByValue         bool               `json:"returnByValue"`         // Whether the result is expected to be a JSON object which should be sent by value.
	GeneratePreview       bool               `json:"generatePreview"`       // Whether preview should be generated for the result.
	AwaitPromise          bool               `json:"awaitPromise"`          // Whether execution should wait for promise to be resolved. If the result of evaluation is not a Promise, it's considered to be an error.
}

type RunScriptResult struct {
	Result           *RemoteObject     `json:"result"`           // Run result.
	ExceptionDetails *ExceptionDetails `json:"exceptionDetails"` // Exception details.
}

type RunScriptCB func(result *RunScriptResult, err error)

// Runs script with given id in a given context.
type RunScriptCommand struct {
	params *RunScriptParams
	cb     RunScriptCB
}

func NewRunScriptCommand(params *RunScriptParams, cb RunScriptCB) *RunScriptCommand {
	return &RunScriptCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *RunScriptCommand) Name() string {
	return "Runtime.runScript"
}

func (cmd *RunScriptCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RunScriptCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj RunScriptResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type ExecutionContextCreatedEvent struct {
	Context *ExecutionContextDescription `json:"context"` // A newly created execution contex.
}

// Issued when new execution context is created.
type ExecutionContextCreatedEventSink struct {
	events chan *ExecutionContextCreatedEvent
}

func NewExecutionContextCreatedEventSink(bufSize int) *ExecutionContextCreatedEventSink {
	return &ExecutionContextCreatedEventSink{
		events: make(chan *ExecutionContextCreatedEvent, bufSize),
	}
}

func (s *ExecutionContextCreatedEventSink) Name() string {
	return "Runtime.executionContextCreated"
}

func (s *ExecutionContextCreatedEventSink) OnEvent(params []byte) {
	evt := &ExecutionContextCreatedEvent{}
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

type ExecutionContextDestroyedEvent struct {
	ExecutionContextId ExecutionContextId `json:"executionContextId"` // Id of the destroyed context
}

// Issued when execution context is destroyed.
type ExecutionContextDestroyedEventSink struct {
	events chan *ExecutionContextDestroyedEvent
}

func NewExecutionContextDestroyedEventSink(bufSize int) *ExecutionContextDestroyedEventSink {
	return &ExecutionContextDestroyedEventSink{
		events: make(chan *ExecutionContextDestroyedEvent, bufSize),
	}
}

func (s *ExecutionContextDestroyedEventSink) Name() string {
	return "Runtime.executionContextDestroyed"
}

func (s *ExecutionContextDestroyedEventSink) OnEvent(params []byte) {
	evt := &ExecutionContextDestroyedEvent{}
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

type ExecutionContextsClearedEvent struct {
}

// Issued when all executionContexts were cleared in browser
type ExecutionContextsClearedEventSink struct {
	events chan *ExecutionContextsClearedEvent
}

func NewExecutionContextsClearedEventSink(bufSize int) *ExecutionContextsClearedEventSink {
	return &ExecutionContextsClearedEventSink{
		events: make(chan *ExecutionContextsClearedEvent, bufSize),
	}
}

func (s *ExecutionContextsClearedEventSink) Name() string {
	return "Runtime.executionContextsCleared"
}

func (s *ExecutionContextsClearedEventSink) OnEvent(params []byte) {
	evt := &ExecutionContextsClearedEvent{}
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

type ExceptionThrownEvent struct {
	Timestamp        RuntimeTimestamp  `json:"timestamp"` // Timestamp of the exception.
	ExceptionDetails *ExceptionDetails `json:"exceptionDetails"`
}

// Issued when exception was thrown and unhandled.
type ExceptionThrownEventSink struct {
	events chan *ExceptionThrownEvent
}

func NewExceptionThrownEventSink(bufSize int) *ExceptionThrownEventSink {
	return &ExceptionThrownEventSink{
		events: make(chan *ExceptionThrownEvent, bufSize),
	}
}

func (s *ExceptionThrownEventSink) Name() string {
	return "Runtime.exceptionThrown"
}

func (s *ExceptionThrownEventSink) OnEvent(params []byte) {
	evt := &ExceptionThrownEvent{}
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

type ExceptionRevokedEvent struct {
	Reason      string `json:"reason"`      // Reason describing why exception was revoked.
	ExceptionId int    `json:"exceptionId"` // The id of revoked exception, as reported in exceptionUnhandled.
}

// Issued when unhandled exception was revoked.
type ExceptionRevokedEventSink struct {
	events chan *ExceptionRevokedEvent
}

func NewExceptionRevokedEventSink(bufSize int) *ExceptionRevokedEventSink {
	return &ExceptionRevokedEventSink{
		events: make(chan *ExceptionRevokedEvent, bufSize),
	}
}

func (s *ExceptionRevokedEventSink) Name() string {
	return "Runtime.exceptionRevoked"
}

func (s *ExceptionRevokedEventSink) OnEvent(params []byte) {
	evt := &ExceptionRevokedEvent{}
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

type ConsoleAPICalledEvent struct {
	Type               string             `json:"type"`               // Type of the call.
	Args               []*RemoteObject    `json:"args"`               // Call arguments.
	ExecutionContextId ExecutionContextId `json:"executionContextId"` // Identifier of the context where the call was made.
	Timestamp          RuntimeTimestamp   `json:"timestamp"`          // Call timestamp.
	StackTrace         *StackTrace        `json:"stackTrace"`         // Stack trace captured when the call was made.
}

// Issued when console API was called.
type ConsoleAPICalledEventSink struct {
	events chan *ConsoleAPICalledEvent
}

func NewConsoleAPICalledEventSink(bufSize int) *ConsoleAPICalledEventSink {
	return &ConsoleAPICalledEventSink{
		events: make(chan *ConsoleAPICalledEvent, bufSize),
	}
}

func (s *ConsoleAPICalledEventSink) Name() string {
	return "Runtime.consoleAPICalled"
}

func (s *ConsoleAPICalledEventSink) OnEvent(params []byte) {
	evt := &ConsoleAPICalledEvent{}
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

type InspectRequestedEvent struct {
	Object *RemoteObject     `json:"object"`
	Hints  map[string]string `json:"hints"`
}

// Issued when object should be inspected (for example, as a result of inspect() command line API call).
type InspectRequestedEventSink struct {
	events chan *InspectRequestedEvent
}

func NewInspectRequestedEventSink(bufSize int) *InspectRequestedEventSink {
	return &InspectRequestedEventSink{
		events: make(chan *InspectRequestedEvent, bufSize),
	}
}

func (s *InspectRequestedEventSink) Name() string {
	return "Runtime.inspectRequested"
}

func (s *InspectRequestedEventSink) OnEvent(params []byte) {
	evt := &InspectRequestedEvent{}
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
