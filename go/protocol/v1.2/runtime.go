package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
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
	Type                string              `json:"type"`                          // Object type.
	Subtype             string              `json:"subtype,omitempty"`             // Object subtype hint. Specified for object type values only.
	ClassName           string              `json:"className,omitempty"`           // Object class (constructor) name. Specified for object type values only.
	Value               json.RawMessage     `json:"value,omitempty"`               // Remote object value in case of primitive values or JSON values (if it was requested).
	UnserializableValue UnserializableValue `json:"unserializableValue,omitempty"` // Primitive value which can not be JSON-stringified does not have value, but gets this property.
	Description         string              `json:"description,omitempty"`         // String representation of the object.
	ObjectId            RemoteObjectId      `json:"objectId,omitempty"`            // Unique object identifier (for non-primitive values).
	Preview             *ObjectPreview      `json:"preview,omitempty"`             // Preview containing abbreviated property values. Specified for object type values only.
	CustomPreview       *CustomPreview      `json:"customPreview,omitempty"`
}

// @experimental
type CustomPreview struct {
	Header                     string         `json:"header"`
	HasBody                    bool           `json:"hasBody"`
	FormatterObjectId          RemoteObjectId `json:"formatterObjectId"`
	BindRemoteObjectFunctionId RemoteObjectId `json:"bindRemoteObjectFunctionId"`
	ConfigObjectId             RemoteObjectId `json:"configObjectId,omitempty"`
}

// Object containing abbreviated remote object value.
// @experimental
type ObjectPreview struct {
	Type        string             `json:"type"`                  // Object type.
	Subtype     string             `json:"subtype,omitempty"`     // Object subtype hint. Specified for object type values only.
	Description string             `json:"description,omitempty"` // String representation of the object.
	Overflow    bool               `json:"overflow"`              // True iff some of the properties or entries of the original object did not fit.
	Properties  []*PropertyPreview `json:"properties"`            // List of the properties.
	Entries     []*EntryPreview    `json:"entries,omitempty"`     // List of the entries. Specified for map and set subtype values only.
}

// @experimental
type PropertyPreview struct {
	Name         string         `json:"name"`                   // Property name.
	Type         string         `json:"type"`                   // Object type. Accessor means that the property itself is an accessor property.
	Value        string         `json:"value,omitempty"`        // User-friendly property value string.
	ValuePreview *ObjectPreview `json:"valuePreview,omitempty"` // Nested value preview.
	Subtype      string         `json:"subtype,omitempty"`      // Object subtype hint. Specified for object type values only.
}

// @experimental
type EntryPreview struct {
	Key   *ObjectPreview `json:"key,omitempty"` // Preview of the key. Specified for map-like collection entries.
	Value *ObjectPreview `json:"value"`         // Preview of the value.
}

// Object property descriptor.
type PropertyDescriptor struct {
	Name         string        `json:"name"`                // Property name or symbol description.
	Value        *RemoteObject `json:"value,omitempty"`     // The value associated with the property.
	Writable     bool          `json:"writable,omitempty"`  // True if the value associated with the property may be changed (data descriptors only).
	Get          *RemoteObject `json:"get,omitempty"`       // A function which serves as a getter for the property, or undefined if there is no getter (accessor descriptors only).
	Set          *RemoteObject `json:"set,omitempty"`       // A function which serves as a setter for the property, or undefined if there is no setter (accessor descriptors only).
	Configurable bool          `json:"configurable"`        // True if the type of this property descriptor may be changed and if the property may be deleted from the corresponding object.
	Enumerable   bool          `json:"enumerable"`          // True if this property shows up during enumeration of the properties on the corresponding object.
	WasThrown    bool          `json:"wasThrown,omitempty"` // True if the result was thrown during the evaluation.
	IsOwn        bool          `json:"isOwn,omitempty"`     // True if the property is owned for the object.
	Symbol       *RemoteObject `json:"symbol,omitempty"`    // Property symbol object, if the property is of the symbol type.
}

// Object internal property descriptor. This property isn't normally visible in JavaScript code.
type InternalPropertyDescriptor struct {
	Name  string        `json:"name"`            // Conventional property name.
	Value *RemoteObject `json:"value,omitempty"` // The value associated with the property.
}

// Represents function call argument. Either remote object id objectId, primitive value, unserializable primitive value or neither of (for undefined) them should be specified.
type CallArgument struct {
	Value               json.RawMessage     `json:"value,omitempty"`               // Primitive value.
	UnserializableValue UnserializableValue `json:"unserializableValue,omitempty"` // Primitive value which can not be JSON-stringified.
	ObjectId            RemoteObjectId      `json:"objectId,omitempty"`            // Remote object handle.
}

// Id of an execution context.
type ExecutionContextId int

// Description of an isolated world.
type ExecutionContextDescription struct {
	Id      ExecutionContextId `json:"id"`                // Unique id of the execution context. It can be used to specify in which execution context script evaluation should be performed.
	Origin  string             `json:"origin"`            // Execution context origin.
	Name    string             `json:"name"`              // Human readable name describing given context.
	AuxData map[string]string  `json:"auxData,omitempty"` // Embedder-specific auxiliary data.
}

// Detailed information about exception (or error) that was thrown during script compilation or execution.
type ExceptionDetails struct {
	ExceptionId        int                `json:"exceptionId"`                  // Exception id.
	Text               string             `json:"text"`                         // Exception text, which should be used together with exception object when available.
	LineNumber         int                `json:"lineNumber"`                   // Line number of the exception location (0-based).
	ColumnNumber       int                `json:"columnNumber"`                 // Column number of the exception location (0-based).
	ScriptId           ScriptId           `json:"scriptId,omitempty"`           // Script ID of the exception location.
	Url                string             `json:"url,omitempty"`                // URL of the exception location, to be used when the script was not reported.
	StackTrace         *StackTrace        `json:"stackTrace,omitempty"`         // JavaScript stack trace if available.
	Exception          *RemoteObject      `json:"exception,omitempty"`          // Exception object if available.
	ExecutionContextId ExecutionContextId `json:"executionContextId,omitempty"` // Identifier of the context where exception happened.
}

// Number of milliseconds since epoch.
type RuntimeTimestamp float64

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
	Description string              `json:"description,omitempty"` // String label of this stack trace. For async traces this may be a name of the function that initiated the async call.
	CallFrames  []*RuntimeCallFrame `json:"callFrames"`            // JavaScript function name.
	Parent      *StackTrace         `json:"parent,omitempty"`      // Asynchronous JavaScript stack trace that preceded this stack, if available.
}

type EvaluateParams struct {
	Expression            string             `json:"expression"`                      // Expression to evaluate.
	ObjectGroup           string             `json:"objectGroup,omitempty"`           // Symbolic group name that can be used to release multiple objects.
	IncludeCommandLineAPI bool               `json:"includeCommandLineAPI,omitempty"` // Determines whether Command Line API should be available during the evaluation.
	Silent                bool               `json:"silent,omitempty"`                // In silent mode exceptions thrown during evaluation are not reported and do not pause execution. Overrides setPauseOnException state.
	ContextId             ExecutionContextId `json:"contextId,omitempty"`             // Specifies in which execution context to perform evaluation. If the parameter is omitted the evaluation will be performed in the context of the inspected page.
	ReturnByValue         bool               `json:"returnByValue,omitempty"`         // Whether the result is expected to be a JSON object that should be sent by value.
	GeneratePreview       bool               `json:"generatePreview,omitempty"`       // Whether preview should be generated for the result.
	UserGesture           bool               `json:"userGesture,omitempty"`           // Whether execution should be treated as initiated by user in the UI.
	AwaitPromise          bool               `json:"awaitPromise,omitempty"`          // Whether execution should wait for promise to be resolved. If the result of evaluation is not a Promise, it's considered to be an error.
}

type EvaluateResult struct {
	Result           *RemoteObject     `json:"result"`           // Evaluation result.
	ExceptionDetails *ExceptionDetails `json:"exceptionDetails"` // Exception details.
}

// Evaluates expression on global object.

type EvaluateCommand struct {
	params *EvaluateParams
	result EvaluateResult
	wg     sync.WaitGroup
	err    error
}

func NewEvaluateCommand(params *EvaluateParams) *EvaluateCommand {
	return &EvaluateCommand{
		params: params,
	}
}

func (cmd *EvaluateCommand) Name() string {
	return "Runtime.evaluate"
}

func (cmd *EvaluateCommand) Params() interface{} {
	return cmd.params
}

func (cmd *EvaluateCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func Evaluate(params *EvaluateParams, conn *hc.Conn) (result *EvaluateResult, err error) {
	cmd := NewEvaluateCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type EvaluateCB func(result *EvaluateResult, err error)

// Evaluates expression on global object.

type AsyncEvaluateCommand struct {
	params *EvaluateParams
	cb     EvaluateCB
}

func NewAsyncEvaluateCommand(params *EvaluateParams, cb EvaluateCB) *AsyncEvaluateCommand {
	return &AsyncEvaluateCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncEvaluateCommand) Name() string {
	return "Runtime.evaluate"
}

func (cmd *AsyncEvaluateCommand) Params() interface{} {
	return cmd.params
}

func (cmd *EvaluateCommand) Result() *EvaluateResult {
	return &cmd.result
}

func (cmd *EvaluateCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncEvaluateCommand) Done(data []byte, err error) {
	var result EvaluateResult
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

type AwaitPromiseParams struct {
	PromiseObjectId RemoteObjectId `json:"promiseObjectId"`           // Identifier of the promise.
	ReturnByValue   bool           `json:"returnByValue,omitempty"`   // Whether the result is expected to be a JSON object that should be sent by value.
	GeneratePreview bool           `json:"generatePreview,omitempty"` // Whether preview should be generated for the result.
}

type AwaitPromiseResult struct {
	Result           *RemoteObject     `json:"result"`           // Promise result. Will contain rejected value if promise was rejected.
	ExceptionDetails *ExceptionDetails `json:"exceptionDetails"` // Exception details if stack strace is available.
}

// Add handler to promise with given promise object id.

type AwaitPromiseCommand struct {
	params *AwaitPromiseParams
	result AwaitPromiseResult
	wg     sync.WaitGroup
	err    error
}

func NewAwaitPromiseCommand(params *AwaitPromiseParams) *AwaitPromiseCommand {
	return &AwaitPromiseCommand{
		params: params,
	}
}

func (cmd *AwaitPromiseCommand) Name() string {
	return "Runtime.awaitPromise"
}

func (cmd *AwaitPromiseCommand) Params() interface{} {
	return cmd.params
}

func (cmd *AwaitPromiseCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func AwaitPromise(params *AwaitPromiseParams, conn *hc.Conn) (result *AwaitPromiseResult, err error) {
	cmd := NewAwaitPromiseCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type AwaitPromiseCB func(result *AwaitPromiseResult, err error)

// Add handler to promise with given promise object id.

type AsyncAwaitPromiseCommand struct {
	params *AwaitPromiseParams
	cb     AwaitPromiseCB
}

func NewAsyncAwaitPromiseCommand(params *AwaitPromiseParams, cb AwaitPromiseCB) *AsyncAwaitPromiseCommand {
	return &AsyncAwaitPromiseCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncAwaitPromiseCommand) Name() string {
	return "Runtime.awaitPromise"
}

func (cmd *AsyncAwaitPromiseCommand) Params() interface{} {
	return cmd.params
}

func (cmd *AwaitPromiseCommand) Result() *AwaitPromiseResult {
	return &cmd.result
}

func (cmd *AwaitPromiseCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncAwaitPromiseCommand) Done(data []byte, err error) {
	var result AwaitPromiseResult
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

type CallFunctionOnParams struct {
	ObjectId            RemoteObjectId  `json:"objectId"`                  // Identifier of the object to call function on.
	FunctionDeclaration string          `json:"functionDeclaration"`       // Declaration of the function to call.
	Arguments           []*CallArgument `json:"arguments,omitempty"`       // Call arguments. All call arguments must belong to the same JavaScript world as the target object.
	Silent              bool            `json:"silent,omitempty"`          // In silent mode exceptions thrown during evaluation are not reported and do not pause execution. Overrides setPauseOnException state.
	ReturnByValue       bool            `json:"returnByValue,omitempty"`   // Whether the result is expected to be a JSON object which should be sent by value.
	GeneratePreview     bool            `json:"generatePreview,omitempty"` // Whether preview should be generated for the result.
	UserGesture         bool            `json:"userGesture,omitempty"`     // Whether execution should be treated as initiated by user in the UI.
	AwaitPromise        bool            `json:"awaitPromise,omitempty"`    // Whether execution should wait for promise to be resolved. If the result of evaluation is not a Promise, it's considered to be an error.
}

type CallFunctionOnResult struct {
	Result           *RemoteObject     `json:"result"`           // Call result.
	ExceptionDetails *ExceptionDetails `json:"exceptionDetails"` // Exception details.
}

// Calls function with given declaration on the given object. Object group of the result is inherited from the target object.

type CallFunctionOnCommand struct {
	params *CallFunctionOnParams
	result CallFunctionOnResult
	wg     sync.WaitGroup
	err    error
}

func NewCallFunctionOnCommand(params *CallFunctionOnParams) *CallFunctionOnCommand {
	return &CallFunctionOnCommand{
		params: params,
	}
}

func (cmd *CallFunctionOnCommand) Name() string {
	return "Runtime.callFunctionOn"
}

func (cmd *CallFunctionOnCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CallFunctionOnCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func CallFunctionOn(params *CallFunctionOnParams, conn *hc.Conn) (result *CallFunctionOnResult, err error) {
	cmd := NewCallFunctionOnCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type CallFunctionOnCB func(result *CallFunctionOnResult, err error)

// Calls function with given declaration on the given object. Object group of the result is inherited from the target object.

type AsyncCallFunctionOnCommand struct {
	params *CallFunctionOnParams
	cb     CallFunctionOnCB
}

func NewAsyncCallFunctionOnCommand(params *CallFunctionOnParams, cb CallFunctionOnCB) *AsyncCallFunctionOnCommand {
	return &AsyncCallFunctionOnCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncCallFunctionOnCommand) Name() string {
	return "Runtime.callFunctionOn"
}

func (cmd *AsyncCallFunctionOnCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CallFunctionOnCommand) Result() *CallFunctionOnResult {
	return &cmd.result
}

func (cmd *CallFunctionOnCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncCallFunctionOnCommand) Done(data []byte, err error) {
	var result CallFunctionOnResult
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

type GetPropertiesParams struct {
	ObjectId               RemoteObjectId `json:"objectId"`                         // Identifier of the object to return properties for.
	OwnProperties          bool           `json:"ownProperties,omitempty"`          // If true, returns properties belonging only to the element itself, not to its prototype chain.
	AccessorPropertiesOnly bool           `json:"accessorPropertiesOnly,omitempty"` // If true, returns accessor properties (with getter/setter) only; internal properties are not returned either.
	GeneratePreview        bool           `json:"generatePreview,omitempty"`        // Whether preview should be generated for the results.
}

type GetPropertiesResult struct {
	Result             []*PropertyDescriptor         `json:"result"`             // Object properties.
	InternalProperties []*InternalPropertyDescriptor `json:"internalProperties"` // Internal object properties (only of the element itself).
	ExceptionDetails   *ExceptionDetails             `json:"exceptionDetails"`   // Exception details.
}

// Returns properties of a given object. Object group of the result is inherited from the target object.

type GetPropertiesCommand struct {
	params *GetPropertiesParams
	result GetPropertiesResult
	wg     sync.WaitGroup
	err    error
}

func NewGetPropertiesCommand(params *GetPropertiesParams) *GetPropertiesCommand {
	return &GetPropertiesCommand{
		params: params,
	}
}

func (cmd *GetPropertiesCommand) Name() string {
	return "Runtime.getProperties"
}

func (cmd *GetPropertiesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetPropertiesCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetProperties(params *GetPropertiesParams, conn *hc.Conn) (result *GetPropertiesResult, err error) {
	cmd := NewGetPropertiesCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetPropertiesCB func(result *GetPropertiesResult, err error)

// Returns properties of a given object. Object group of the result is inherited from the target object.

type AsyncGetPropertiesCommand struct {
	params *GetPropertiesParams
	cb     GetPropertiesCB
}

func NewAsyncGetPropertiesCommand(params *GetPropertiesParams, cb GetPropertiesCB) *AsyncGetPropertiesCommand {
	return &AsyncGetPropertiesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetPropertiesCommand) Name() string {
	return "Runtime.getProperties"
}

func (cmd *AsyncGetPropertiesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetPropertiesCommand) Result() *GetPropertiesResult {
	return &cmd.result
}

func (cmd *GetPropertiesCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetPropertiesCommand) Done(data []byte, err error) {
	var result GetPropertiesResult
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

type ReleaseObjectParams struct {
	ObjectId RemoteObjectId `json:"objectId"` // Identifier of the object to release.
}

// Releases remote object with given id.

type ReleaseObjectCommand struct {
	params *ReleaseObjectParams
	wg     sync.WaitGroup
	err    error
}

func NewReleaseObjectCommand(params *ReleaseObjectParams) *ReleaseObjectCommand {
	return &ReleaseObjectCommand{
		params: params,
	}
}

func (cmd *ReleaseObjectCommand) Name() string {
	return "Runtime.releaseObject"
}

func (cmd *ReleaseObjectCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ReleaseObjectCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ReleaseObject(params *ReleaseObjectParams, conn *hc.Conn) (err error) {
	cmd := NewReleaseObjectCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type ReleaseObjectCB func(err error)

// Releases remote object with given id.

type AsyncReleaseObjectCommand struct {
	params *ReleaseObjectParams
	cb     ReleaseObjectCB
}

func NewAsyncReleaseObjectCommand(params *ReleaseObjectParams, cb ReleaseObjectCB) *AsyncReleaseObjectCommand {
	return &AsyncReleaseObjectCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncReleaseObjectCommand) Name() string {
	return "Runtime.releaseObject"
}

func (cmd *AsyncReleaseObjectCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ReleaseObjectCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncReleaseObjectCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type ReleaseObjectGroupParams struct {
	ObjectGroup string `json:"objectGroup"` // Symbolic object group name.
}

// Releases all remote objects that belong to a given group.

type ReleaseObjectGroupCommand struct {
	params *ReleaseObjectGroupParams
	wg     sync.WaitGroup
	err    error
}

func NewReleaseObjectGroupCommand(params *ReleaseObjectGroupParams) *ReleaseObjectGroupCommand {
	return &ReleaseObjectGroupCommand{
		params: params,
	}
}

func (cmd *ReleaseObjectGroupCommand) Name() string {
	return "Runtime.releaseObjectGroup"
}

func (cmd *ReleaseObjectGroupCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ReleaseObjectGroupCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ReleaseObjectGroup(params *ReleaseObjectGroupParams, conn *hc.Conn) (err error) {
	cmd := NewReleaseObjectGroupCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type ReleaseObjectGroupCB func(err error)

// Releases all remote objects that belong to a given group.

type AsyncReleaseObjectGroupCommand struct {
	params *ReleaseObjectGroupParams
	cb     ReleaseObjectGroupCB
}

func NewAsyncReleaseObjectGroupCommand(params *ReleaseObjectGroupParams, cb ReleaseObjectGroupCB) *AsyncReleaseObjectGroupCommand {
	return &AsyncReleaseObjectGroupCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncReleaseObjectGroupCommand) Name() string {
	return "Runtime.releaseObjectGroup"
}

func (cmd *AsyncReleaseObjectGroupCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ReleaseObjectGroupCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncReleaseObjectGroupCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Tells inspected instance to run if it was waiting for debugger to attach.

type RunIfWaitingForDebuggerCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewRunIfWaitingForDebuggerCommand() *RunIfWaitingForDebuggerCommand {
	return &RunIfWaitingForDebuggerCommand{}
}

func (cmd *RunIfWaitingForDebuggerCommand) Name() string {
	return "Runtime.runIfWaitingForDebugger"
}

func (cmd *RunIfWaitingForDebuggerCommand) Params() interface{} {
	return nil
}

func (cmd *RunIfWaitingForDebuggerCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func RunIfWaitingForDebugger(conn *hc.Conn) (err error) {
	cmd := NewRunIfWaitingForDebuggerCommand()
	cmd.Run(conn)
	return cmd.err
}

type RunIfWaitingForDebuggerCB func(err error)

// Tells inspected instance to run if it was waiting for debugger to attach.

type AsyncRunIfWaitingForDebuggerCommand struct {
	cb RunIfWaitingForDebuggerCB
}

func NewAsyncRunIfWaitingForDebuggerCommand(cb RunIfWaitingForDebuggerCB) *AsyncRunIfWaitingForDebuggerCommand {
	return &AsyncRunIfWaitingForDebuggerCommand{
		cb: cb,
	}
}

func (cmd *AsyncRunIfWaitingForDebuggerCommand) Name() string {
	return "Runtime.runIfWaitingForDebugger"
}

func (cmd *AsyncRunIfWaitingForDebuggerCommand) Params() interface{} {
	return nil
}

func (cmd *RunIfWaitingForDebuggerCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncRunIfWaitingForDebuggerCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Enables reporting of execution contexts creation by means of executionContextCreated event. When the reporting gets enabled the event will be sent immediately for each existing execution context.

type RuntimeEnableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewRuntimeEnableCommand() *RuntimeEnableCommand {
	return &RuntimeEnableCommand{}
}

func (cmd *RuntimeEnableCommand) Name() string {
	return "Runtime.enable"
}

func (cmd *RuntimeEnableCommand) Params() interface{} {
	return nil
}

func (cmd *RuntimeEnableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func RuntimeEnable(conn *hc.Conn) (err error) {
	cmd := NewRuntimeEnableCommand()
	cmd.Run(conn)
	return cmd.err
}

type RuntimeEnableCB func(err error)

// Enables reporting of execution contexts creation by means of executionContextCreated event. When the reporting gets enabled the event will be sent immediately for each existing execution context.

type AsyncRuntimeEnableCommand struct {
	cb RuntimeEnableCB
}

func NewAsyncRuntimeEnableCommand(cb RuntimeEnableCB) *AsyncRuntimeEnableCommand {
	return &AsyncRuntimeEnableCommand{
		cb: cb,
	}
}

func (cmd *AsyncRuntimeEnableCommand) Name() string {
	return "Runtime.enable"
}

func (cmd *AsyncRuntimeEnableCommand) Params() interface{} {
	return nil
}

func (cmd *RuntimeEnableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncRuntimeEnableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Disables reporting of execution contexts creation.

type RuntimeDisableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewRuntimeDisableCommand() *RuntimeDisableCommand {
	return &RuntimeDisableCommand{}
}

func (cmd *RuntimeDisableCommand) Name() string {
	return "Runtime.disable"
}

func (cmd *RuntimeDisableCommand) Params() interface{} {
	return nil
}

func (cmd *RuntimeDisableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func RuntimeDisable(conn *hc.Conn) (err error) {
	cmd := NewRuntimeDisableCommand()
	cmd.Run(conn)
	return cmd.err
}

type RuntimeDisableCB func(err error)

// Disables reporting of execution contexts creation.

type AsyncRuntimeDisableCommand struct {
	cb RuntimeDisableCB
}

func NewAsyncRuntimeDisableCommand(cb RuntimeDisableCB) *AsyncRuntimeDisableCommand {
	return &AsyncRuntimeDisableCommand{
		cb: cb,
	}
}

func (cmd *AsyncRuntimeDisableCommand) Name() string {
	return "Runtime.disable"
}

func (cmd *AsyncRuntimeDisableCommand) Params() interface{} {
	return nil
}

func (cmd *RuntimeDisableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncRuntimeDisableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Discards collected exceptions and console API calls.

type DiscardConsoleEntriesCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewDiscardConsoleEntriesCommand() *DiscardConsoleEntriesCommand {
	return &DiscardConsoleEntriesCommand{}
}

func (cmd *DiscardConsoleEntriesCommand) Name() string {
	return "Runtime.discardConsoleEntries"
}

func (cmd *DiscardConsoleEntriesCommand) Params() interface{} {
	return nil
}

func (cmd *DiscardConsoleEntriesCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func DiscardConsoleEntries(conn *hc.Conn) (err error) {
	cmd := NewDiscardConsoleEntriesCommand()
	cmd.Run(conn)
	return cmd.err
}

type DiscardConsoleEntriesCB func(err error)

// Discards collected exceptions and console API calls.

type AsyncDiscardConsoleEntriesCommand struct {
	cb DiscardConsoleEntriesCB
}

func NewAsyncDiscardConsoleEntriesCommand(cb DiscardConsoleEntriesCB) *AsyncDiscardConsoleEntriesCommand {
	return &AsyncDiscardConsoleEntriesCommand{
		cb: cb,
	}
}

func (cmd *AsyncDiscardConsoleEntriesCommand) Name() string {
	return "Runtime.discardConsoleEntries"
}

func (cmd *AsyncDiscardConsoleEntriesCommand) Params() interface{} {
	return nil
}

func (cmd *DiscardConsoleEntriesCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncDiscardConsoleEntriesCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetCustomObjectFormatterEnabledParams struct {
	Enabled bool `json:"enabled"`
}

// @experimental
type SetCustomObjectFormatterEnabledCommand struct {
	params *SetCustomObjectFormatterEnabledParams
	wg     sync.WaitGroup
	err    error
}

func NewSetCustomObjectFormatterEnabledCommand(params *SetCustomObjectFormatterEnabledParams) *SetCustomObjectFormatterEnabledCommand {
	return &SetCustomObjectFormatterEnabledCommand{
		params: params,
	}
}

func (cmd *SetCustomObjectFormatterEnabledCommand) Name() string {
	return "Runtime.setCustomObjectFormatterEnabled"
}

func (cmd *SetCustomObjectFormatterEnabledCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetCustomObjectFormatterEnabledCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetCustomObjectFormatterEnabled(params *SetCustomObjectFormatterEnabledParams, conn *hc.Conn) (err error) {
	cmd := NewSetCustomObjectFormatterEnabledCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetCustomObjectFormatterEnabledCB func(err error)

// @experimental
type AsyncSetCustomObjectFormatterEnabledCommand struct {
	params *SetCustomObjectFormatterEnabledParams
	cb     SetCustomObjectFormatterEnabledCB
}

func NewAsyncSetCustomObjectFormatterEnabledCommand(params *SetCustomObjectFormatterEnabledParams, cb SetCustomObjectFormatterEnabledCB) *AsyncSetCustomObjectFormatterEnabledCommand {
	return &AsyncSetCustomObjectFormatterEnabledCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetCustomObjectFormatterEnabledCommand) Name() string {
	return "Runtime.setCustomObjectFormatterEnabled"
}

func (cmd *AsyncSetCustomObjectFormatterEnabledCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetCustomObjectFormatterEnabledCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetCustomObjectFormatterEnabledCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type CompileScriptParams struct {
	Expression         string             `json:"expression"`                   // Expression to compile.
	SourceURL          string             `json:"sourceURL"`                    // Source url to be set for the script.
	PersistScript      bool               `json:"persistScript"`                // Specifies whether the compiled script should be persisted.
	ExecutionContextId ExecutionContextId `json:"executionContextId,omitempty"` // Specifies in which execution context to perform script run. If the parameter is omitted the evaluation will be performed in the context of the inspected page.
}

type CompileScriptResult struct {
	ScriptId         ScriptId          `json:"scriptId"`         // Id of the script.
	ExceptionDetails *ExceptionDetails `json:"exceptionDetails"` // Exception details.
}

// Compiles expression.

type CompileScriptCommand struct {
	params *CompileScriptParams
	result CompileScriptResult
	wg     sync.WaitGroup
	err    error
}

func NewCompileScriptCommand(params *CompileScriptParams) *CompileScriptCommand {
	return &CompileScriptCommand{
		params: params,
	}
}

func (cmd *CompileScriptCommand) Name() string {
	return "Runtime.compileScript"
}

func (cmd *CompileScriptCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CompileScriptCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func CompileScript(params *CompileScriptParams, conn *hc.Conn) (result *CompileScriptResult, err error) {
	cmd := NewCompileScriptCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type CompileScriptCB func(result *CompileScriptResult, err error)

// Compiles expression.

type AsyncCompileScriptCommand struct {
	params *CompileScriptParams
	cb     CompileScriptCB
}

func NewAsyncCompileScriptCommand(params *CompileScriptParams, cb CompileScriptCB) *AsyncCompileScriptCommand {
	return &AsyncCompileScriptCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncCompileScriptCommand) Name() string {
	return "Runtime.compileScript"
}

func (cmd *AsyncCompileScriptCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CompileScriptCommand) Result() *CompileScriptResult {
	return &cmd.result
}

func (cmd *CompileScriptCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncCompileScriptCommand) Done(data []byte, err error) {
	var result CompileScriptResult
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

type RunScriptParams struct {
	ScriptId              ScriptId           `json:"scriptId"`                        // Id of the script to run.
	ExecutionContextId    ExecutionContextId `json:"executionContextId,omitempty"`    // Specifies in which execution context to perform script run. If the parameter is omitted the evaluation will be performed in the context of the inspected page.
	ObjectGroup           string             `json:"objectGroup,omitempty"`           // Symbolic group name that can be used to release multiple objects.
	Silent                bool               `json:"silent,omitempty"`                // In silent mode exceptions thrown during evaluation are not reported and do not pause execution. Overrides setPauseOnException state.
	IncludeCommandLineAPI bool               `json:"includeCommandLineAPI,omitempty"` // Determines whether Command Line API should be available during the evaluation.
	ReturnByValue         bool               `json:"returnByValue,omitempty"`         // Whether the result is expected to be a JSON object which should be sent by value.
	GeneratePreview       bool               `json:"generatePreview,omitempty"`       // Whether preview should be generated for the result.
	AwaitPromise          bool               `json:"awaitPromise,omitempty"`          // Whether execution should wait for promise to be resolved. If the result of evaluation is not a Promise, it's considered to be an error.
}

type RunScriptResult struct {
	Result           *RemoteObject     `json:"result"`           // Run result.
	ExceptionDetails *ExceptionDetails `json:"exceptionDetails"` // Exception details.
}

// Runs script with given id in a given context.

type RunScriptCommand struct {
	params *RunScriptParams
	result RunScriptResult
	wg     sync.WaitGroup
	err    error
}

func NewRunScriptCommand(params *RunScriptParams) *RunScriptCommand {
	return &RunScriptCommand{
		params: params,
	}
}

func (cmd *RunScriptCommand) Name() string {
	return "Runtime.runScript"
}

func (cmd *RunScriptCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RunScriptCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func RunScript(params *RunScriptParams, conn *hc.Conn) (result *RunScriptResult, err error) {
	cmd := NewRunScriptCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type RunScriptCB func(result *RunScriptResult, err error)

// Runs script with given id in a given context.

type AsyncRunScriptCommand struct {
	params *RunScriptParams
	cb     RunScriptCB
}

func NewAsyncRunScriptCommand(params *RunScriptParams, cb RunScriptCB) *AsyncRunScriptCommand {
	return &AsyncRunScriptCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncRunScriptCommand) Name() string {
	return "Runtime.runScript"
}

func (cmd *AsyncRunScriptCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RunScriptCommand) Result() *RunScriptResult {
	return &cmd.result
}

func (cmd *RunScriptCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncRunScriptCommand) Done(data []byte, err error) {
	var result RunScriptResult
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

// Issued when new execution context is created.

type ExecutionContextCreatedEvent struct {
	Context *ExecutionContextDescription `json:"context"` // A newly created execution contex.
}

func OnExecutionContextCreated(conn *hc.Conn, cb func(evt *ExecutionContextCreatedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &ExecutionContextCreatedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Runtime.executionContextCreated", sink)
}

// Issued when execution context is destroyed.

type ExecutionContextDestroyedEvent struct {
	ExecutionContextId ExecutionContextId `json:"executionContextId"` // Id of the destroyed context
}

func OnExecutionContextDestroyed(conn *hc.Conn, cb func(evt *ExecutionContextDestroyedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &ExecutionContextDestroyedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Runtime.executionContextDestroyed", sink)
}

// Issued when all executionContexts were cleared in browser

type ExecutionContextsClearedEvent struct {
}

func OnExecutionContextsCleared(conn *hc.Conn, cb func(evt *ExecutionContextsClearedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &ExecutionContextsClearedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Runtime.executionContextsCleared", sink)
}

// Issued when exception was thrown and unhandled.

type ExceptionThrownEvent struct {
	Timestamp        RuntimeTimestamp  `json:"timestamp"` // Timestamp of the exception.
	ExceptionDetails *ExceptionDetails `json:"exceptionDetails"`
}

func OnExceptionThrown(conn *hc.Conn, cb func(evt *ExceptionThrownEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &ExceptionThrownEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Runtime.exceptionThrown", sink)
}

// Issued when unhandled exception was revoked.

type ExceptionRevokedEvent struct {
	Reason      string `json:"reason"`      // Reason describing why exception was revoked.
	ExceptionId int    `json:"exceptionId"` // The id of revoked exception, as reported in exceptionUnhandled.
}

func OnExceptionRevoked(conn *hc.Conn, cb func(evt *ExceptionRevokedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &ExceptionRevokedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Runtime.exceptionRevoked", sink)
}

// Issued when console API was called.

type ConsoleAPICalledEvent struct {
	Type               string             `json:"type"`               // Type of the call.
	Args               []*RemoteObject    `json:"args"`               // Call arguments.
	ExecutionContextId ExecutionContextId `json:"executionContextId"` // Identifier of the context where the call was made.
	Timestamp          RuntimeTimestamp   `json:"timestamp"`          // Call timestamp.
	StackTrace         *StackTrace        `json:"stackTrace"`         // Stack trace captured when the call was made.
}

func OnConsoleAPICalled(conn *hc.Conn, cb func(evt *ConsoleAPICalledEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &ConsoleAPICalledEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Runtime.consoleAPICalled", sink)
}

// Issued when object should be inspected (for example, as a result of inspect() command line API call).

type InspectRequestedEvent struct {
	Object *RemoteObject     `json:"object"`
	Hints  map[string]string `json:"hints"`
}

func OnInspectRequested(conn *hc.Conn, cb func(evt *InspectRequestedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &InspectRequestedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Runtime.inspectRequested", sink)
}
