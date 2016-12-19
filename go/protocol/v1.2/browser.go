package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
)

type BrowserContextID string

type BrowserTargetID string

type BrowserTargetInfo struct {
	TargetId BrowserTargetID `json:"targetId"`
	Type     string          `json:"type"`
	Title    string          `json:"title"`
	Url      string          `json:"url"`
}

type RemoteLocation struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type CreateBrowserContextResult struct {
	BrowserContextId BrowserContextID `json:"browserContextId"` // The id of the context created.
}

type CreateBrowserContextCB func(result *CreateBrowserContextResult, err error)

// Creates a new empty BrowserContext. Similar to an incognito profile but you can have more than one.
type CreateBrowserContextCommand struct {
	cb CreateBrowserContextCB
}

func NewCreateBrowserContextCommand(cb CreateBrowserContextCB) *CreateBrowserContextCommand {
	return &CreateBrowserContextCommand{
		cb: cb,
	}
}

func (cmd *CreateBrowserContextCommand) Name() string {
	return "Browser.createBrowserContext"
}

func (cmd *CreateBrowserContextCommand) Params() interface{} {
	return nil
}

func (cmd *CreateBrowserContextCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj CreateBrowserContextResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type DisposeBrowserContextParams struct {
	BrowserContextId BrowserContextID `json:"browserContextId"`
}

type DisposeBrowserContextResult struct {
	Success bool `json:"success"`
}

type DisposeBrowserContextCB func(result *DisposeBrowserContextResult, err error)

// Deletes a BrowserContext, will fail of any open page uses it.
type DisposeBrowserContextCommand struct {
	params *DisposeBrowserContextParams
	cb     DisposeBrowserContextCB
}

func NewDisposeBrowserContextCommand(params *DisposeBrowserContextParams, cb DisposeBrowserContextCB) *DisposeBrowserContextCommand {
	return &DisposeBrowserContextCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *DisposeBrowserContextCommand) Name() string {
	return "Browser.disposeBrowserContext"
}

func (cmd *DisposeBrowserContextCommand) Params() interface{} {
	return cmd.params
}

func (cmd *DisposeBrowserContextCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj DisposeBrowserContextResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type CreateTargetParams struct {
	Url              string           `json:"url"`              // The initial URL the page will be navigated to.
	Width            int              `json:"width"`            // Frame width in DIP (headless chrome only).
	Height           int              `json:"height"`           // Frame height in DIP (headless chrome only).
	BrowserContextId BrowserContextID `json:"browserContextId"` // The browser context to create the page in (headless chrome only).
}

type CreateTargetResult struct {
	TargetId BrowserTargetID `json:"targetId"` // The id of the page opened.
}

type CreateTargetCB func(result *CreateTargetResult, err error)

// Creates a new page.
type CreateTargetCommand struct {
	params *CreateTargetParams
	cb     CreateTargetCB
}

func NewCreateTargetCommand(params *CreateTargetParams, cb CreateTargetCB) *CreateTargetCommand {
	return &CreateTargetCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *CreateTargetCommand) Name() string {
	return "Browser.createTarget"
}

func (cmd *CreateTargetCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CreateTargetCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj CreateTargetResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type CloseTargetParams struct {
	TargetId BrowserTargetID `json:"targetId"`
}

type CloseTargetResult struct {
	Success bool `json:"success"`
}

type CloseTargetCB func(result *CloseTargetResult, err error)

// Closes the target. If the target is a page that gets closed too.
type CloseTargetCommand struct {
	params *CloseTargetParams
	cb     CloseTargetCB
}

func NewCloseTargetCommand(params *CloseTargetParams, cb CloseTargetCB) *CloseTargetCommand {
	return &CloseTargetCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *CloseTargetCommand) Name() string {
	return "Browser.closeTarget"
}

func (cmd *CloseTargetCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CloseTargetCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj CloseTargetResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type GetTargetsResult struct {
	TargetInfo []*BrowserTargetInfo `json:"targetInfo"`
}

type GetTargetsCB func(result *GetTargetsResult, err error)

// Returns target information for all potential targets.
type GetTargetsCommand struct {
	cb GetTargetsCB
}

func NewGetTargetsCommand(cb GetTargetsCB) *GetTargetsCommand {
	return &GetTargetsCommand{
		cb: cb,
	}
}

func (cmd *GetTargetsCommand) Name() string {
	return "Browser.getTargets"
}

func (cmd *GetTargetsCommand) Params() interface{} {
	return nil
}

func (cmd *GetTargetsCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetTargetsResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type SetRemoteLocationsParams struct {
	Locations []*RemoteLocation `json:"locations"` // List of remote locations
}

type SetRemoteLocationsCB func(err error)

// Enables target discovery for the specified locations.
type SetRemoteLocationsCommand struct {
	params *SetRemoteLocationsParams
	cb     SetRemoteLocationsCB
}

func NewSetRemoteLocationsCommand(params *SetRemoteLocationsParams, cb SetRemoteLocationsCB) *SetRemoteLocationsCommand {
	return &SetRemoteLocationsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetRemoteLocationsCommand) Name() string {
	return "Browser.setRemoteLocations"
}

func (cmd *SetRemoteLocationsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetRemoteLocationsCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type AttachParams struct {
	TargetId BrowserTargetID `json:"targetId"` // Target id.
}

type AttachResult struct {
	Success bool `json:"success"` // Whether attach succeeded.
}

type AttachCB func(result *AttachResult, err error)

// Attaches to the target with given id.
type AttachCommand struct {
	params *AttachParams
	cb     AttachCB
}

func NewAttachCommand(params *AttachParams, cb AttachCB) *AttachCommand {
	return &AttachCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AttachCommand) Name() string {
	return "Browser.attach"
}

func (cmd *AttachCommand) Params() interface{} {
	return cmd.params
}

func (cmd *AttachCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj AttachResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type DetachParams struct {
	TargetId BrowserTargetID `json:"targetId"`
}

type DetachResult struct {
	Success bool `json:"success"` // Whether detach succeeded.
}

type DetachCB func(result *DetachResult, err error)

// Detaches from the target with given id.
type DetachCommand struct {
	params *DetachParams
	cb     DetachCB
}

func NewDetachCommand(params *DetachParams, cb DetachCB) *DetachCommand {
	return &DetachCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *DetachCommand) Name() string {
	return "Browser.detach"
}

func (cmd *DetachCommand) Params() interface{} {
	return cmd.params
}

func (cmd *DetachCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj DetachResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type SendMessageParams struct {
	TargetId BrowserTargetID `json:"targetId"`
	Message  string          `json:"message"`
}

type SendMessageCB func(err error)

// Sends protocol message to the target with given id.
type SendMessageCommand struct {
	params *SendMessageParams
	cb     SendMessageCB
}

func NewSendMessageCommand(params *SendMessageParams, cb SendMessageCB) *SendMessageCommand {
	return &SendMessageCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SendMessageCommand) Name() string {
	return "Browser.sendMessage"
}

func (cmd *SendMessageCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SendMessageCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type DispatchMessageEvent struct {
	TargetId BrowserTargetID `json:"targetId"`
	Message  string          `json:"message"`
}

// Dispatches protocol message from the target with given id.
type DispatchMessageEventSink struct {
	events chan *DispatchMessageEvent
}

func NewDispatchMessageEventSink(bufSize int) *DispatchMessageEventSink {
	return &DispatchMessageEventSink{
		events: make(chan *DispatchMessageEvent, bufSize),
	}
}

func (s *DispatchMessageEventSink) Name() string {
	return "Browser.dispatchMessage"
}

func (s *DispatchMessageEventSink) OnEvent(params []byte) {
	evt := &DispatchMessageEvent{}
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
