package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

// DOM breakpoint type.
type DOMBreakpointType string

const DOMBreakpointTypeSubtreeModified DOMBreakpointType = "subtree-modified"
const DOMBreakpointTypeAttributeModified DOMBreakpointType = "attribute-modified"
const DOMBreakpointTypeNodeRemoved DOMBreakpointType = "node-removed"

// Object event listener.
// @experimental
type EventListener struct {
	Type            string        `json:"type"`                      // EventListener's type.
	UseCapture      bool          `json:"useCapture"`                // EventListener's useCapture.
	Passive         bool          `json:"passive"`                   // EventListener's passive flag.
	Once            bool          `json:"once"`                      // EventListener's once flag.
	ScriptId        *ScriptId     `json:"scriptId"`                  // Script id of the handler code.
	LineNumber      int           `json:"lineNumber"`                // Line number in the script (0-based).
	ColumnNumber    int           `json:"columnNumber"`              // Column number in the script (0-based).
	Handler         *RemoteObject `json:"handler,omitempty"`         // Event handler function value.
	OriginalHandler *RemoteObject `json:"originalHandler,omitempty"` // Event original handler function value.
	RemoveFunction  *RemoteObject `json:"removeFunction,omitempty"`  // Event listener remove function.
}

type SetDOMBreakpointParams struct {
	NodeId *NodeId           `json:"nodeId"` // Identifier of the node to set breakpoint on.
	Type   DOMBreakpointType `json:"type"`   // Type of the operation to stop upon.
}

// Sets breakpoint on particular operation with DOM.

type SetDOMBreakpointCommand struct {
	params *SetDOMBreakpointParams
	wg     sync.WaitGroup
	err    error
}

func NewSetDOMBreakpointCommand(params *SetDOMBreakpointParams) *SetDOMBreakpointCommand {
	return &SetDOMBreakpointCommand{
		params: params,
	}
}

func (cmd *SetDOMBreakpointCommand) Name() string {
	return "DOMDebugger.setDOMBreakpoint"
}

func (cmd *SetDOMBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetDOMBreakpointCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetDOMBreakpoint(params *SetDOMBreakpointParams, conn *hc.Conn) (err error) {
	cmd := NewSetDOMBreakpointCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetDOMBreakpointCB func(err error)

// Sets breakpoint on particular operation with DOM.

type AsyncSetDOMBreakpointCommand struct {
	params *SetDOMBreakpointParams
	cb     SetDOMBreakpointCB
}

func NewAsyncSetDOMBreakpointCommand(params *SetDOMBreakpointParams, cb SetDOMBreakpointCB) *AsyncSetDOMBreakpointCommand {
	return &AsyncSetDOMBreakpointCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetDOMBreakpointCommand) Name() string {
	return "DOMDebugger.setDOMBreakpoint"
}

func (cmd *AsyncSetDOMBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetDOMBreakpointCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetDOMBreakpointCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type RemoveDOMBreakpointParams struct {
	NodeId *NodeId           `json:"nodeId"` // Identifier of the node to remove breakpoint from.
	Type   DOMBreakpointType `json:"type"`   // Type of the breakpoint to remove.
}

// Removes DOM breakpoint that was set using setDOMBreakpoint.

type RemoveDOMBreakpointCommand struct {
	params *RemoveDOMBreakpointParams
	wg     sync.WaitGroup
	err    error
}

func NewRemoveDOMBreakpointCommand(params *RemoveDOMBreakpointParams) *RemoveDOMBreakpointCommand {
	return &RemoveDOMBreakpointCommand{
		params: params,
	}
}

func (cmd *RemoveDOMBreakpointCommand) Name() string {
	return "DOMDebugger.removeDOMBreakpoint"
}

func (cmd *RemoveDOMBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveDOMBreakpointCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func RemoveDOMBreakpoint(params *RemoveDOMBreakpointParams, conn *hc.Conn) (err error) {
	cmd := NewRemoveDOMBreakpointCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type RemoveDOMBreakpointCB func(err error)

// Removes DOM breakpoint that was set using setDOMBreakpoint.

type AsyncRemoveDOMBreakpointCommand struct {
	params *RemoveDOMBreakpointParams
	cb     RemoveDOMBreakpointCB
}

func NewAsyncRemoveDOMBreakpointCommand(params *RemoveDOMBreakpointParams, cb RemoveDOMBreakpointCB) *AsyncRemoveDOMBreakpointCommand {
	return &AsyncRemoveDOMBreakpointCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncRemoveDOMBreakpointCommand) Name() string {
	return "DOMDebugger.removeDOMBreakpoint"
}

func (cmd *AsyncRemoveDOMBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveDOMBreakpointCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncRemoveDOMBreakpointCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetEventListenerBreakpointParams struct {
	EventName  string `json:"eventName"`            // DOM Event name to stop on (any DOM event will do).
	TargetName string `json:"targetName,omitempty"` // EventTarget interface name to stop on. If equal to "*" or not provided, will stop on any EventTarget.
}

// Sets breakpoint on particular DOM event.

type SetEventListenerBreakpointCommand struct {
	params *SetEventListenerBreakpointParams
	wg     sync.WaitGroup
	err    error
}

func NewSetEventListenerBreakpointCommand(params *SetEventListenerBreakpointParams) *SetEventListenerBreakpointCommand {
	return &SetEventListenerBreakpointCommand{
		params: params,
	}
}

func (cmd *SetEventListenerBreakpointCommand) Name() string {
	return "DOMDebugger.setEventListenerBreakpoint"
}

func (cmd *SetEventListenerBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetEventListenerBreakpointCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetEventListenerBreakpoint(params *SetEventListenerBreakpointParams, conn *hc.Conn) (err error) {
	cmd := NewSetEventListenerBreakpointCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetEventListenerBreakpointCB func(err error)

// Sets breakpoint on particular DOM event.

type AsyncSetEventListenerBreakpointCommand struct {
	params *SetEventListenerBreakpointParams
	cb     SetEventListenerBreakpointCB
}

func NewAsyncSetEventListenerBreakpointCommand(params *SetEventListenerBreakpointParams, cb SetEventListenerBreakpointCB) *AsyncSetEventListenerBreakpointCommand {
	return &AsyncSetEventListenerBreakpointCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetEventListenerBreakpointCommand) Name() string {
	return "DOMDebugger.setEventListenerBreakpoint"
}

func (cmd *AsyncSetEventListenerBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetEventListenerBreakpointCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetEventListenerBreakpointCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type RemoveEventListenerBreakpointParams struct {
	EventName  string `json:"eventName"`            // Event name.
	TargetName string `json:"targetName,omitempty"` // EventTarget interface name.
}

// Removes breakpoint on particular DOM event.

type RemoveEventListenerBreakpointCommand struct {
	params *RemoveEventListenerBreakpointParams
	wg     sync.WaitGroup
	err    error
}

func NewRemoveEventListenerBreakpointCommand(params *RemoveEventListenerBreakpointParams) *RemoveEventListenerBreakpointCommand {
	return &RemoveEventListenerBreakpointCommand{
		params: params,
	}
}

func (cmd *RemoveEventListenerBreakpointCommand) Name() string {
	return "DOMDebugger.removeEventListenerBreakpoint"
}

func (cmd *RemoveEventListenerBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveEventListenerBreakpointCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func RemoveEventListenerBreakpoint(params *RemoveEventListenerBreakpointParams, conn *hc.Conn) (err error) {
	cmd := NewRemoveEventListenerBreakpointCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type RemoveEventListenerBreakpointCB func(err error)

// Removes breakpoint on particular DOM event.

type AsyncRemoveEventListenerBreakpointCommand struct {
	params *RemoveEventListenerBreakpointParams
	cb     RemoveEventListenerBreakpointCB
}

func NewAsyncRemoveEventListenerBreakpointCommand(params *RemoveEventListenerBreakpointParams, cb RemoveEventListenerBreakpointCB) *AsyncRemoveEventListenerBreakpointCommand {
	return &AsyncRemoveEventListenerBreakpointCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncRemoveEventListenerBreakpointCommand) Name() string {
	return "DOMDebugger.removeEventListenerBreakpoint"
}

func (cmd *AsyncRemoveEventListenerBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveEventListenerBreakpointCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncRemoveEventListenerBreakpointCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetInstrumentationBreakpointParams struct {
	EventName string `json:"eventName"` // Instrumentation name to stop on.
}

// Sets breakpoint on particular native event.
// @experimental
type SetInstrumentationBreakpointCommand struct {
	params *SetInstrumentationBreakpointParams
	wg     sync.WaitGroup
	err    error
}

func NewSetInstrumentationBreakpointCommand(params *SetInstrumentationBreakpointParams) *SetInstrumentationBreakpointCommand {
	return &SetInstrumentationBreakpointCommand{
		params: params,
	}
}

func (cmd *SetInstrumentationBreakpointCommand) Name() string {
	return "DOMDebugger.setInstrumentationBreakpoint"
}

func (cmd *SetInstrumentationBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetInstrumentationBreakpointCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetInstrumentationBreakpoint(params *SetInstrumentationBreakpointParams, conn *hc.Conn) (err error) {
	cmd := NewSetInstrumentationBreakpointCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetInstrumentationBreakpointCB func(err error)

// Sets breakpoint on particular native event.
// @experimental
type AsyncSetInstrumentationBreakpointCommand struct {
	params *SetInstrumentationBreakpointParams
	cb     SetInstrumentationBreakpointCB
}

func NewAsyncSetInstrumentationBreakpointCommand(params *SetInstrumentationBreakpointParams, cb SetInstrumentationBreakpointCB) *AsyncSetInstrumentationBreakpointCommand {
	return &AsyncSetInstrumentationBreakpointCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetInstrumentationBreakpointCommand) Name() string {
	return "DOMDebugger.setInstrumentationBreakpoint"
}

func (cmd *AsyncSetInstrumentationBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetInstrumentationBreakpointCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetInstrumentationBreakpointCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type RemoveInstrumentationBreakpointParams struct {
	EventName string `json:"eventName"` // Instrumentation name to stop on.
}

// Removes breakpoint on particular native event.
// @experimental
type RemoveInstrumentationBreakpointCommand struct {
	params *RemoveInstrumentationBreakpointParams
	wg     sync.WaitGroup
	err    error
}

func NewRemoveInstrumentationBreakpointCommand(params *RemoveInstrumentationBreakpointParams) *RemoveInstrumentationBreakpointCommand {
	return &RemoveInstrumentationBreakpointCommand{
		params: params,
	}
}

func (cmd *RemoveInstrumentationBreakpointCommand) Name() string {
	return "DOMDebugger.removeInstrumentationBreakpoint"
}

func (cmd *RemoveInstrumentationBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveInstrumentationBreakpointCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func RemoveInstrumentationBreakpoint(params *RemoveInstrumentationBreakpointParams, conn *hc.Conn) (err error) {
	cmd := NewRemoveInstrumentationBreakpointCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type RemoveInstrumentationBreakpointCB func(err error)

// Removes breakpoint on particular native event.
// @experimental
type AsyncRemoveInstrumentationBreakpointCommand struct {
	params *RemoveInstrumentationBreakpointParams
	cb     RemoveInstrumentationBreakpointCB
}

func NewAsyncRemoveInstrumentationBreakpointCommand(params *RemoveInstrumentationBreakpointParams, cb RemoveInstrumentationBreakpointCB) *AsyncRemoveInstrumentationBreakpointCommand {
	return &AsyncRemoveInstrumentationBreakpointCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncRemoveInstrumentationBreakpointCommand) Name() string {
	return "DOMDebugger.removeInstrumentationBreakpoint"
}

func (cmd *AsyncRemoveInstrumentationBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveInstrumentationBreakpointCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncRemoveInstrumentationBreakpointCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetXHRBreakpointParams struct {
	Url string `json:"url"` // Resource URL substring. All XHRs having this substring in the URL will get stopped upon.
}

// Sets breakpoint on XMLHttpRequest.

type SetXHRBreakpointCommand struct {
	params *SetXHRBreakpointParams
	wg     sync.WaitGroup
	err    error
}

func NewSetXHRBreakpointCommand(params *SetXHRBreakpointParams) *SetXHRBreakpointCommand {
	return &SetXHRBreakpointCommand{
		params: params,
	}
}

func (cmd *SetXHRBreakpointCommand) Name() string {
	return "DOMDebugger.setXHRBreakpoint"
}

func (cmd *SetXHRBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetXHRBreakpointCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetXHRBreakpoint(params *SetXHRBreakpointParams, conn *hc.Conn) (err error) {
	cmd := NewSetXHRBreakpointCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetXHRBreakpointCB func(err error)

// Sets breakpoint on XMLHttpRequest.

type AsyncSetXHRBreakpointCommand struct {
	params *SetXHRBreakpointParams
	cb     SetXHRBreakpointCB
}

func NewAsyncSetXHRBreakpointCommand(params *SetXHRBreakpointParams, cb SetXHRBreakpointCB) *AsyncSetXHRBreakpointCommand {
	return &AsyncSetXHRBreakpointCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetXHRBreakpointCommand) Name() string {
	return "DOMDebugger.setXHRBreakpoint"
}

func (cmd *AsyncSetXHRBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetXHRBreakpointCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetXHRBreakpointCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type RemoveXHRBreakpointParams struct {
	Url string `json:"url"` // Resource URL substring.
}

// Removes breakpoint from XMLHttpRequest.

type RemoveXHRBreakpointCommand struct {
	params *RemoveXHRBreakpointParams
	wg     sync.WaitGroup
	err    error
}

func NewRemoveXHRBreakpointCommand(params *RemoveXHRBreakpointParams) *RemoveXHRBreakpointCommand {
	return &RemoveXHRBreakpointCommand{
		params: params,
	}
}

func (cmd *RemoveXHRBreakpointCommand) Name() string {
	return "DOMDebugger.removeXHRBreakpoint"
}

func (cmd *RemoveXHRBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveXHRBreakpointCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func RemoveXHRBreakpoint(params *RemoveXHRBreakpointParams, conn *hc.Conn) (err error) {
	cmd := NewRemoveXHRBreakpointCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type RemoveXHRBreakpointCB func(err error)

// Removes breakpoint from XMLHttpRequest.

type AsyncRemoveXHRBreakpointCommand struct {
	params *RemoveXHRBreakpointParams
	cb     RemoveXHRBreakpointCB
}

func NewAsyncRemoveXHRBreakpointCommand(params *RemoveXHRBreakpointParams, cb RemoveXHRBreakpointCB) *AsyncRemoveXHRBreakpointCommand {
	return &AsyncRemoveXHRBreakpointCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncRemoveXHRBreakpointCommand) Name() string {
	return "DOMDebugger.removeXHRBreakpoint"
}

func (cmd *AsyncRemoveXHRBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveXHRBreakpointCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncRemoveXHRBreakpointCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type GetEventListenersParams struct {
	ObjectId *RemoteObjectId `json:"objectId"` // Identifier of the object to return listeners for.
}

type GetEventListenersResult struct {
	Listeners []*EventListener `json:"listeners"` // Array of relevant listeners.
}

// Returns event listeners of the given object.
// @experimental
type GetEventListenersCommand struct {
	params *GetEventListenersParams
	result GetEventListenersResult
	wg     sync.WaitGroup
	err    error
}

func NewGetEventListenersCommand(params *GetEventListenersParams) *GetEventListenersCommand {
	return &GetEventListenersCommand{
		params: params,
	}
}

func (cmd *GetEventListenersCommand) Name() string {
	return "DOMDebugger.getEventListeners"
}

func (cmd *GetEventListenersCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetEventListenersCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetEventListeners(params *GetEventListenersParams, conn *hc.Conn) (result *GetEventListenersResult, err error) {
	cmd := NewGetEventListenersCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetEventListenersCB func(result *GetEventListenersResult, err error)

// Returns event listeners of the given object.
// @experimental
type AsyncGetEventListenersCommand struct {
	params *GetEventListenersParams
	cb     GetEventListenersCB
}

func NewAsyncGetEventListenersCommand(params *GetEventListenersParams, cb GetEventListenersCB) *AsyncGetEventListenersCommand {
	return &AsyncGetEventListenersCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetEventListenersCommand) Name() string {
	return "DOMDebugger.getEventListeners"
}

func (cmd *AsyncGetEventListenersCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetEventListenersCommand) Result() *GetEventListenersResult {
	return &cmd.result
}

func (cmd *GetEventListenersCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetEventListenersCommand) Done(data []byte, err error) {
	var result GetEventListenersResult
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
