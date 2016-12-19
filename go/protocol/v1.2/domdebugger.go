package protocol

import (
	"encoding/json"
)

// DOM breakpoint type.
type DOMBreakpointType string

const DOMBreakpointTypeSubtreeModified DOMBreakpointType = "subtree-modified"
const DOMBreakpointTypeAttributeModified DOMBreakpointType = "attribute-modified"
const DOMBreakpointTypeNodeRemoved DOMBreakpointType = "node-removed"

// Object event listener.
type EventListener struct {
	Type            string        `json:"type"`            // EventListener's type.
	UseCapture      bool          `json:"useCapture"`      // EventListener's useCapture.
	Passive         bool          `json:"passive"`         // EventListener's passive flag.
	ScriptId        *ScriptId     `json:"scriptId"`        // Script id of the handler code.
	LineNumber      int           `json:"lineNumber"`      // Line number in the script (0-based).
	ColumnNumber    int           `json:"columnNumber"`    // Column number in the script (0-based).
	Handler         *RemoteObject `json:"handler"`         // Event handler function value.
	OriginalHandler *RemoteObject `json:"originalHandler"` // Event original handler function value.
	RemoveFunction  *RemoteObject `json:"removeFunction"`  // Event listener remove function.
}

type SetDOMBreakpointParams struct {
	NodeId *NodeId           `json:"nodeId"` // Identifier of the node to set breakpoint on.
	Type   DOMBreakpointType `json:"type"`   // Type of the operation to stop upon.
}

type SetDOMBreakpointCB func(err error)

// Sets breakpoint on particular operation with DOM.
type SetDOMBreakpointCommand struct {
	params *SetDOMBreakpointParams
	cb     SetDOMBreakpointCB
}

func NewSetDOMBreakpointCommand(params *SetDOMBreakpointParams, cb SetDOMBreakpointCB) *SetDOMBreakpointCommand {
	return &SetDOMBreakpointCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetDOMBreakpointCommand) Name() string {
	return "DOMDebugger.setDOMBreakpoint"
}

func (cmd *SetDOMBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetDOMBreakpointCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type RemoveDOMBreakpointParams struct {
	NodeId *NodeId           `json:"nodeId"` // Identifier of the node to remove breakpoint from.
	Type   DOMBreakpointType `json:"type"`   // Type of the breakpoint to remove.
}

type RemoveDOMBreakpointCB func(err error)

// Removes DOM breakpoint that was set using setDOMBreakpoint.
type RemoveDOMBreakpointCommand struct {
	params *RemoveDOMBreakpointParams
	cb     RemoveDOMBreakpointCB
}

func NewRemoveDOMBreakpointCommand(params *RemoveDOMBreakpointParams, cb RemoveDOMBreakpointCB) *RemoveDOMBreakpointCommand {
	return &RemoveDOMBreakpointCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *RemoveDOMBreakpointCommand) Name() string {
	return "DOMDebugger.removeDOMBreakpoint"
}

func (cmd *RemoveDOMBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveDOMBreakpointCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetEventListenerBreakpointParams struct {
	EventName  string `json:"eventName"`  // DOM Event name to stop on (any DOM event will do).
	TargetName string `json:"targetName"` // EventTarget interface name to stop on. If equal to "*" or not provided, will stop on any EventTarget.
}

type SetEventListenerBreakpointCB func(err error)

// Sets breakpoint on particular DOM event.
type SetEventListenerBreakpointCommand struct {
	params *SetEventListenerBreakpointParams
	cb     SetEventListenerBreakpointCB
}

func NewSetEventListenerBreakpointCommand(params *SetEventListenerBreakpointParams, cb SetEventListenerBreakpointCB) *SetEventListenerBreakpointCommand {
	return &SetEventListenerBreakpointCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetEventListenerBreakpointCommand) Name() string {
	return "DOMDebugger.setEventListenerBreakpoint"
}

func (cmd *SetEventListenerBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetEventListenerBreakpointCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type RemoveEventListenerBreakpointParams struct {
	EventName  string `json:"eventName"`  // Event name.
	TargetName string `json:"targetName"` // EventTarget interface name.
}

type RemoveEventListenerBreakpointCB func(err error)

// Removes breakpoint on particular DOM event.
type RemoveEventListenerBreakpointCommand struct {
	params *RemoveEventListenerBreakpointParams
	cb     RemoveEventListenerBreakpointCB
}

func NewRemoveEventListenerBreakpointCommand(params *RemoveEventListenerBreakpointParams, cb RemoveEventListenerBreakpointCB) *RemoveEventListenerBreakpointCommand {
	return &RemoveEventListenerBreakpointCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *RemoveEventListenerBreakpointCommand) Name() string {
	return "DOMDebugger.removeEventListenerBreakpoint"
}

func (cmd *RemoveEventListenerBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveEventListenerBreakpointCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetInstrumentationBreakpointParams struct {
	EventName string `json:"eventName"` // Instrumentation name to stop on.
}

type SetInstrumentationBreakpointCB func(err error)

// Sets breakpoint on particular native event.
type SetInstrumentationBreakpointCommand struct {
	params *SetInstrumentationBreakpointParams
	cb     SetInstrumentationBreakpointCB
}

func NewSetInstrumentationBreakpointCommand(params *SetInstrumentationBreakpointParams, cb SetInstrumentationBreakpointCB) *SetInstrumentationBreakpointCommand {
	return &SetInstrumentationBreakpointCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetInstrumentationBreakpointCommand) Name() string {
	return "DOMDebugger.setInstrumentationBreakpoint"
}

func (cmd *SetInstrumentationBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetInstrumentationBreakpointCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type RemoveInstrumentationBreakpointParams struct {
	EventName string `json:"eventName"` // Instrumentation name to stop on.
}

type RemoveInstrumentationBreakpointCB func(err error)

// Removes breakpoint on particular native event.
type RemoveInstrumentationBreakpointCommand struct {
	params *RemoveInstrumentationBreakpointParams
	cb     RemoveInstrumentationBreakpointCB
}

func NewRemoveInstrumentationBreakpointCommand(params *RemoveInstrumentationBreakpointParams, cb RemoveInstrumentationBreakpointCB) *RemoveInstrumentationBreakpointCommand {
	return &RemoveInstrumentationBreakpointCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *RemoveInstrumentationBreakpointCommand) Name() string {
	return "DOMDebugger.removeInstrumentationBreakpoint"
}

func (cmd *RemoveInstrumentationBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveInstrumentationBreakpointCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetXHRBreakpointParams struct {
	Url string `json:"url"` // Resource URL substring. All XHRs having this substring in the URL will get stopped upon.
}

type SetXHRBreakpointCB func(err error)

// Sets breakpoint on XMLHttpRequest.
type SetXHRBreakpointCommand struct {
	params *SetXHRBreakpointParams
	cb     SetXHRBreakpointCB
}

func NewSetXHRBreakpointCommand(params *SetXHRBreakpointParams, cb SetXHRBreakpointCB) *SetXHRBreakpointCommand {
	return &SetXHRBreakpointCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetXHRBreakpointCommand) Name() string {
	return "DOMDebugger.setXHRBreakpoint"
}

func (cmd *SetXHRBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetXHRBreakpointCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type RemoveXHRBreakpointParams struct {
	Url string `json:"url"` // Resource URL substring.
}

type RemoveXHRBreakpointCB func(err error)

// Removes breakpoint from XMLHttpRequest.
type RemoveXHRBreakpointCommand struct {
	params *RemoveXHRBreakpointParams
	cb     RemoveXHRBreakpointCB
}

func NewRemoveXHRBreakpointCommand(params *RemoveXHRBreakpointParams, cb RemoveXHRBreakpointCB) *RemoveXHRBreakpointCommand {
	return &RemoveXHRBreakpointCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *RemoveXHRBreakpointCommand) Name() string {
	return "DOMDebugger.removeXHRBreakpoint"
}

func (cmd *RemoveXHRBreakpointCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveXHRBreakpointCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type GetEventListenersParams struct {
	ObjectId *RemoteObjectId `json:"objectId"` // Identifier of the object to return listeners for.
}

type GetEventListenersResult struct {
	Listeners []*EventListener `json:"listeners"` // Array of relevant listeners.
}

type GetEventListenersCB func(result *GetEventListenersResult, err error)

// Returns event listeners of the given object.
type GetEventListenersCommand struct {
	params *GetEventListenersParams
	cb     GetEventListenersCB
}

func NewGetEventListenersCommand(params *GetEventListenersParams, cb GetEventListenersCB) *GetEventListenersCommand {
	return &GetEventListenersCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetEventListenersCommand) Name() string {
	return "DOMDebugger.getEventListeners"
}

func (cmd *GetEventListenersCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetEventListenersCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetEventListenersResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}
