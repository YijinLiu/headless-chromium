package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

type TargetID string

type BrowserContextID string

type TargetInfo struct {
	TargetId TargetID `json:"targetId"`
	Type     string   `json:"type"`
	Title    string   `json:"title"`
	Url      string   `json:"url"`
}

type RemoteLocation struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type SetDiscoverTargetsParams struct {
	Discover bool `json:"discover"` // Whether to discover available targets.
}

// Controls whether to discover available targets and notify via targetCreated/targetDestroyed events.

type SetDiscoverTargetsCommand struct {
	params *SetDiscoverTargetsParams
	wg     sync.WaitGroup
	err    error
}

func NewSetDiscoverTargetsCommand(params *SetDiscoverTargetsParams) *SetDiscoverTargetsCommand {
	return &SetDiscoverTargetsCommand{
		params: params,
	}
}

func (cmd *SetDiscoverTargetsCommand) Name() string {
	return "Target.setDiscoverTargets"
}

func (cmd *SetDiscoverTargetsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetDiscoverTargetsCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetDiscoverTargets(params *SetDiscoverTargetsParams, conn *hc.Conn) (err error) {
	cmd := NewSetDiscoverTargetsCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetDiscoverTargetsCB func(err error)

// Controls whether to discover available targets and notify via targetCreated/targetDestroyed events.

type AsyncSetDiscoverTargetsCommand struct {
	params *SetDiscoverTargetsParams
	cb     SetDiscoverTargetsCB
}

func NewAsyncSetDiscoverTargetsCommand(params *SetDiscoverTargetsParams, cb SetDiscoverTargetsCB) *AsyncSetDiscoverTargetsCommand {
	return &AsyncSetDiscoverTargetsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetDiscoverTargetsCommand) Name() string {
	return "Target.setDiscoverTargets"
}

func (cmd *AsyncSetDiscoverTargetsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetDiscoverTargetsCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetDiscoverTargetsCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetAutoAttachParams struct {
	AutoAttach             bool `json:"autoAttach"`             // Whether to auto-attach to related targets.
	WaitForDebuggerOnStart bool `json:"waitForDebuggerOnStart"` // Whether to pause new targets when attaching to them. Use Runtime.runIfWaitingForDebugger to run paused targets.
}

// Controls whether to automatically attach to new targets which are considered to be related to this one. When turned on, attaches to all existing related targets as well. When turned off, automatically detaches from all currently attached targets.

type SetAutoAttachCommand struct {
	params *SetAutoAttachParams
	wg     sync.WaitGroup
	err    error
}

func NewSetAutoAttachCommand(params *SetAutoAttachParams) *SetAutoAttachCommand {
	return &SetAutoAttachCommand{
		params: params,
	}
}

func (cmd *SetAutoAttachCommand) Name() string {
	return "Target.setAutoAttach"
}

func (cmd *SetAutoAttachCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetAutoAttachCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetAutoAttach(params *SetAutoAttachParams, conn *hc.Conn) (err error) {
	cmd := NewSetAutoAttachCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetAutoAttachCB func(err error)

// Controls whether to automatically attach to new targets which are considered to be related to this one. When turned on, attaches to all existing related targets as well. When turned off, automatically detaches from all currently attached targets.

type AsyncSetAutoAttachCommand struct {
	params *SetAutoAttachParams
	cb     SetAutoAttachCB
}

func NewAsyncSetAutoAttachCommand(params *SetAutoAttachParams, cb SetAutoAttachCB) *AsyncSetAutoAttachCommand {
	return &AsyncSetAutoAttachCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetAutoAttachCommand) Name() string {
	return "Target.setAutoAttach"
}

func (cmd *AsyncSetAutoAttachCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetAutoAttachCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetAutoAttachCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetAttachToFramesParams struct {
	Value bool `json:"value"` // Whether to attach to frames.
}

type SetAttachToFramesCommand struct {
	params *SetAttachToFramesParams
	wg     sync.WaitGroup
	err    error
}

func NewSetAttachToFramesCommand(params *SetAttachToFramesParams) *SetAttachToFramesCommand {
	return &SetAttachToFramesCommand{
		params: params,
	}
}

func (cmd *SetAttachToFramesCommand) Name() string {
	return "Target.setAttachToFrames"
}

func (cmd *SetAttachToFramesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetAttachToFramesCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetAttachToFrames(params *SetAttachToFramesParams, conn *hc.Conn) (err error) {
	cmd := NewSetAttachToFramesCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetAttachToFramesCB func(err error)

type AsyncSetAttachToFramesCommand struct {
	params *SetAttachToFramesParams
	cb     SetAttachToFramesCB
}

func NewAsyncSetAttachToFramesCommand(params *SetAttachToFramesParams, cb SetAttachToFramesCB) *AsyncSetAttachToFramesCommand {
	return &AsyncSetAttachToFramesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetAttachToFramesCommand) Name() string {
	return "Target.setAttachToFrames"
}

func (cmd *AsyncSetAttachToFramesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetAttachToFramesCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetAttachToFramesCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetRemoteLocationsParams struct {
	Locations []*RemoteLocation `json:"locations"` // List of remote locations.
}

// Enables target discovery for the specified locations, when setDiscoverTargets was set to true.

type SetRemoteLocationsCommand struct {
	params *SetRemoteLocationsParams
	wg     sync.WaitGroup
	err    error
}

func NewSetRemoteLocationsCommand(params *SetRemoteLocationsParams) *SetRemoteLocationsCommand {
	return &SetRemoteLocationsCommand{
		params: params,
	}
}

func (cmd *SetRemoteLocationsCommand) Name() string {
	return "Target.setRemoteLocations"
}

func (cmd *SetRemoteLocationsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetRemoteLocationsCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetRemoteLocations(params *SetRemoteLocationsParams, conn *hc.Conn) (err error) {
	cmd := NewSetRemoteLocationsCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetRemoteLocationsCB func(err error)

// Enables target discovery for the specified locations, when setDiscoverTargets was set to true.

type AsyncSetRemoteLocationsCommand struct {
	params *SetRemoteLocationsParams
	cb     SetRemoteLocationsCB
}

func NewAsyncSetRemoteLocationsCommand(params *SetRemoteLocationsParams, cb SetRemoteLocationsCB) *AsyncSetRemoteLocationsCommand {
	return &AsyncSetRemoteLocationsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetRemoteLocationsCommand) Name() string {
	return "Target.setRemoteLocations"
}

func (cmd *AsyncSetRemoteLocationsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetRemoteLocationsCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetRemoteLocationsCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SendMessageToTargetParams struct {
	TargetId string `json:"targetId"`
	Message  string `json:"message"`
}

// Sends protocol message to the target with given id.

type SendMessageToTargetCommand struct {
	params *SendMessageToTargetParams
	wg     sync.WaitGroup
	err    error
}

func NewSendMessageToTargetCommand(params *SendMessageToTargetParams) *SendMessageToTargetCommand {
	return &SendMessageToTargetCommand{
		params: params,
	}
}

func (cmd *SendMessageToTargetCommand) Name() string {
	return "Target.sendMessageToTarget"
}

func (cmd *SendMessageToTargetCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SendMessageToTargetCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SendMessageToTarget(params *SendMessageToTargetParams, conn *hc.Conn) (err error) {
	cmd := NewSendMessageToTargetCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SendMessageToTargetCB func(err error)

// Sends protocol message to the target with given id.

type AsyncSendMessageToTargetCommand struct {
	params *SendMessageToTargetParams
	cb     SendMessageToTargetCB
}

func NewAsyncSendMessageToTargetCommand(params *SendMessageToTargetParams, cb SendMessageToTargetCB) *AsyncSendMessageToTargetCommand {
	return &AsyncSendMessageToTargetCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSendMessageToTargetCommand) Name() string {
	return "Target.sendMessageToTarget"
}

func (cmd *AsyncSendMessageToTargetCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SendMessageToTargetCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSendMessageToTargetCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type GetTargetInfoParams struct {
	TargetId TargetID `json:"targetId"`
}

type GetTargetInfoResult struct {
	TargetInfo *TargetInfo `json:"targetInfo"`
}

// Returns information about a target.

type GetTargetInfoCommand struct {
	params *GetTargetInfoParams
	result GetTargetInfoResult
	wg     sync.WaitGroup
	err    error
}

func NewGetTargetInfoCommand(params *GetTargetInfoParams) *GetTargetInfoCommand {
	return &GetTargetInfoCommand{
		params: params,
	}
}

func (cmd *GetTargetInfoCommand) Name() string {
	return "Target.getTargetInfo"
}

func (cmd *GetTargetInfoCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetTargetInfoCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetTargetInfo(params *GetTargetInfoParams, conn *hc.Conn) (result *GetTargetInfoResult, err error) {
	cmd := NewGetTargetInfoCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetTargetInfoCB func(result *GetTargetInfoResult, err error)

// Returns information about a target.

type AsyncGetTargetInfoCommand struct {
	params *GetTargetInfoParams
	cb     GetTargetInfoCB
}

func NewAsyncGetTargetInfoCommand(params *GetTargetInfoParams, cb GetTargetInfoCB) *AsyncGetTargetInfoCommand {
	return &AsyncGetTargetInfoCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetTargetInfoCommand) Name() string {
	return "Target.getTargetInfo"
}

func (cmd *AsyncGetTargetInfoCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetTargetInfoCommand) Result() *GetTargetInfoResult {
	return &cmd.result
}

func (cmd *GetTargetInfoCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetTargetInfoCommand) Done(data []byte, err error) {
	var result GetTargetInfoResult
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

type ActivateTargetParams struct {
	TargetId TargetID `json:"targetId"`
}

// Activates (focuses) the target.

type ActivateTargetCommand struct {
	params *ActivateTargetParams
	wg     sync.WaitGroup
	err    error
}

func NewActivateTargetCommand(params *ActivateTargetParams) *ActivateTargetCommand {
	return &ActivateTargetCommand{
		params: params,
	}
}

func (cmd *ActivateTargetCommand) Name() string {
	return "Target.activateTarget"
}

func (cmd *ActivateTargetCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ActivateTargetCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ActivateTarget(params *ActivateTargetParams, conn *hc.Conn) (err error) {
	cmd := NewActivateTargetCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type ActivateTargetCB func(err error)

// Activates (focuses) the target.

type AsyncActivateTargetCommand struct {
	params *ActivateTargetParams
	cb     ActivateTargetCB
}

func NewAsyncActivateTargetCommand(params *ActivateTargetParams, cb ActivateTargetCB) *AsyncActivateTargetCommand {
	return &AsyncActivateTargetCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncActivateTargetCommand) Name() string {
	return "Target.activateTarget"
}

func (cmd *AsyncActivateTargetCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ActivateTargetCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncActivateTargetCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type CloseTargetParams struct {
	TargetId TargetID `json:"targetId"`
}

type CloseTargetResult struct {
	Success bool `json:"success"`
}

// Closes the target. If the target is a page that gets closed too.

type CloseTargetCommand struct {
	params *CloseTargetParams
	result CloseTargetResult
	wg     sync.WaitGroup
	err    error
}

func NewCloseTargetCommand(params *CloseTargetParams) *CloseTargetCommand {
	return &CloseTargetCommand{
		params: params,
	}
}

func (cmd *CloseTargetCommand) Name() string {
	return "Target.closeTarget"
}

func (cmd *CloseTargetCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CloseTargetCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func CloseTarget(params *CloseTargetParams, conn *hc.Conn) (result *CloseTargetResult, err error) {
	cmd := NewCloseTargetCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type CloseTargetCB func(result *CloseTargetResult, err error)

// Closes the target. If the target is a page that gets closed too.

type AsyncCloseTargetCommand struct {
	params *CloseTargetParams
	cb     CloseTargetCB
}

func NewAsyncCloseTargetCommand(params *CloseTargetParams, cb CloseTargetCB) *AsyncCloseTargetCommand {
	return &AsyncCloseTargetCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncCloseTargetCommand) Name() string {
	return "Target.closeTarget"
}

func (cmd *AsyncCloseTargetCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CloseTargetCommand) Result() *CloseTargetResult {
	return &cmd.result
}

func (cmd *CloseTargetCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncCloseTargetCommand) Done(data []byte, err error) {
	var result CloseTargetResult
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

type AttachToTargetParams struct {
	TargetId TargetID `json:"targetId"`
}

type AttachToTargetResult struct {
	Success bool `json:"success"` // Whether attach succeeded.
}

// Attaches to the target with given id.

type AttachToTargetCommand struct {
	params *AttachToTargetParams
	result AttachToTargetResult
	wg     sync.WaitGroup
	err    error
}

func NewAttachToTargetCommand(params *AttachToTargetParams) *AttachToTargetCommand {
	return &AttachToTargetCommand{
		params: params,
	}
}

func (cmd *AttachToTargetCommand) Name() string {
	return "Target.attachToTarget"
}

func (cmd *AttachToTargetCommand) Params() interface{} {
	return cmd.params
}

func (cmd *AttachToTargetCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func AttachToTarget(params *AttachToTargetParams, conn *hc.Conn) (result *AttachToTargetResult, err error) {
	cmd := NewAttachToTargetCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type AttachToTargetCB func(result *AttachToTargetResult, err error)

// Attaches to the target with given id.

type AsyncAttachToTargetCommand struct {
	params *AttachToTargetParams
	cb     AttachToTargetCB
}

func NewAsyncAttachToTargetCommand(params *AttachToTargetParams, cb AttachToTargetCB) *AsyncAttachToTargetCommand {
	return &AsyncAttachToTargetCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncAttachToTargetCommand) Name() string {
	return "Target.attachToTarget"
}

func (cmd *AsyncAttachToTargetCommand) Params() interface{} {
	return cmd.params
}

func (cmd *AttachToTargetCommand) Result() *AttachToTargetResult {
	return &cmd.result
}

func (cmd *AttachToTargetCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncAttachToTargetCommand) Done(data []byte, err error) {
	var result AttachToTargetResult
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

type DetachFromTargetParams struct {
	TargetId TargetID `json:"targetId"`
}

// Detaches from the target with given id.

type DetachFromTargetCommand struct {
	params *DetachFromTargetParams
	wg     sync.WaitGroup
	err    error
}

func NewDetachFromTargetCommand(params *DetachFromTargetParams) *DetachFromTargetCommand {
	return &DetachFromTargetCommand{
		params: params,
	}
}

func (cmd *DetachFromTargetCommand) Name() string {
	return "Target.detachFromTarget"
}

func (cmd *DetachFromTargetCommand) Params() interface{} {
	return cmd.params
}

func (cmd *DetachFromTargetCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func DetachFromTarget(params *DetachFromTargetParams, conn *hc.Conn) (err error) {
	cmd := NewDetachFromTargetCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type DetachFromTargetCB func(err error)

// Detaches from the target with given id.

type AsyncDetachFromTargetCommand struct {
	params *DetachFromTargetParams
	cb     DetachFromTargetCB
}

func NewAsyncDetachFromTargetCommand(params *DetachFromTargetParams, cb DetachFromTargetCB) *AsyncDetachFromTargetCommand {
	return &AsyncDetachFromTargetCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncDetachFromTargetCommand) Name() string {
	return "Target.detachFromTarget"
}

func (cmd *AsyncDetachFromTargetCommand) Params() interface{} {
	return cmd.params
}

func (cmd *DetachFromTargetCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncDetachFromTargetCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type CreateBrowserContextResult struct {
	BrowserContextId BrowserContextID `json:"browserContextId"` // The id of the context created.
}

// Creates a new empty BrowserContext. Similar to an incognito profile but you can have more than one.

type CreateBrowserContextCommand struct {
	result CreateBrowserContextResult
	wg     sync.WaitGroup
	err    error
}

func NewCreateBrowserContextCommand() *CreateBrowserContextCommand {
	return &CreateBrowserContextCommand{}
}

func (cmd *CreateBrowserContextCommand) Name() string {
	return "Target.createBrowserContext"
}

func (cmd *CreateBrowserContextCommand) Params() interface{} {
	return nil
}

func (cmd *CreateBrowserContextCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func CreateBrowserContext(conn *hc.Conn) (result *CreateBrowserContextResult, err error) {
	cmd := NewCreateBrowserContextCommand()
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type CreateBrowserContextCB func(result *CreateBrowserContextResult, err error)

// Creates a new empty BrowserContext. Similar to an incognito profile but you can have more than one.

type AsyncCreateBrowserContextCommand struct {
	cb CreateBrowserContextCB
}

func NewAsyncCreateBrowserContextCommand(cb CreateBrowserContextCB) *AsyncCreateBrowserContextCommand {
	return &AsyncCreateBrowserContextCommand{
		cb: cb,
	}
}

func (cmd *AsyncCreateBrowserContextCommand) Name() string {
	return "Target.createBrowserContext"
}

func (cmd *AsyncCreateBrowserContextCommand) Params() interface{} {
	return nil
}

func (cmd *CreateBrowserContextCommand) Result() *CreateBrowserContextResult {
	return &cmd.result
}

func (cmd *CreateBrowserContextCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncCreateBrowserContextCommand) Done(data []byte, err error) {
	var result CreateBrowserContextResult
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

type DisposeBrowserContextParams struct {
	BrowserContextId BrowserContextID `json:"browserContextId"`
}

type DisposeBrowserContextResult struct {
	Success bool `json:"success"`
}

// Deletes a BrowserContext, will fail of any open page uses it.

type DisposeBrowserContextCommand struct {
	params *DisposeBrowserContextParams
	result DisposeBrowserContextResult
	wg     sync.WaitGroup
	err    error
}

func NewDisposeBrowserContextCommand(params *DisposeBrowserContextParams) *DisposeBrowserContextCommand {
	return &DisposeBrowserContextCommand{
		params: params,
	}
}

func (cmd *DisposeBrowserContextCommand) Name() string {
	return "Target.disposeBrowserContext"
}

func (cmd *DisposeBrowserContextCommand) Params() interface{} {
	return cmd.params
}

func (cmd *DisposeBrowserContextCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func DisposeBrowserContext(params *DisposeBrowserContextParams, conn *hc.Conn) (result *DisposeBrowserContextResult, err error) {
	cmd := NewDisposeBrowserContextCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type DisposeBrowserContextCB func(result *DisposeBrowserContextResult, err error)

// Deletes a BrowserContext, will fail of any open page uses it.

type AsyncDisposeBrowserContextCommand struct {
	params *DisposeBrowserContextParams
	cb     DisposeBrowserContextCB
}

func NewAsyncDisposeBrowserContextCommand(params *DisposeBrowserContextParams, cb DisposeBrowserContextCB) *AsyncDisposeBrowserContextCommand {
	return &AsyncDisposeBrowserContextCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncDisposeBrowserContextCommand) Name() string {
	return "Target.disposeBrowserContext"
}

func (cmd *AsyncDisposeBrowserContextCommand) Params() interface{} {
	return cmd.params
}

func (cmd *DisposeBrowserContextCommand) Result() *DisposeBrowserContextResult {
	return &cmd.result
}

func (cmd *DisposeBrowserContextCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncDisposeBrowserContextCommand) Done(data []byte, err error) {
	var result DisposeBrowserContextResult
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

type CreateTargetParams struct {
	Url              string           `json:"url"`                        // The initial URL the page will be navigated to.
	Width            int              `json:"width,omitempty"`            // Frame width in DIP (headless chrome only).
	Height           int              `json:"height,omitempty"`           // Frame height in DIP (headless chrome only).
	BrowserContextId BrowserContextID `json:"browserContextId,omitempty"` // The browser context to create the page in (headless chrome only).
}

type CreateTargetResult struct {
	TargetId TargetID `json:"targetId"` // The id of the page opened.
}

// Creates a new page.

type CreateTargetCommand struct {
	params *CreateTargetParams
	result CreateTargetResult
	wg     sync.WaitGroup
	err    error
}

func NewCreateTargetCommand(params *CreateTargetParams) *CreateTargetCommand {
	return &CreateTargetCommand{
		params: params,
	}
}

func (cmd *CreateTargetCommand) Name() string {
	return "Target.createTarget"
}

func (cmd *CreateTargetCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CreateTargetCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func CreateTarget(params *CreateTargetParams, conn *hc.Conn) (result *CreateTargetResult, err error) {
	cmd := NewCreateTargetCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type CreateTargetCB func(result *CreateTargetResult, err error)

// Creates a new page.

type AsyncCreateTargetCommand struct {
	params *CreateTargetParams
	cb     CreateTargetCB
}

func NewAsyncCreateTargetCommand(params *CreateTargetParams, cb CreateTargetCB) *AsyncCreateTargetCommand {
	return &AsyncCreateTargetCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncCreateTargetCommand) Name() string {
	return "Target.createTarget"
}

func (cmd *AsyncCreateTargetCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CreateTargetCommand) Result() *CreateTargetResult {
	return &cmd.result
}

func (cmd *CreateTargetCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncCreateTargetCommand) Done(data []byte, err error) {
	var result CreateTargetResult
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

type GetTargetsResult struct {
	TargetInfos []*TargetInfo `json:"targetInfos"` // The list of targets.
}

// Retrieves a list of available targets.

type GetTargetsCommand struct {
	result GetTargetsResult
	wg     sync.WaitGroup
	err    error
}

func NewGetTargetsCommand() *GetTargetsCommand {
	return &GetTargetsCommand{}
}

func (cmd *GetTargetsCommand) Name() string {
	return "Target.getTargets"
}

func (cmd *GetTargetsCommand) Params() interface{} {
	return nil
}

func (cmd *GetTargetsCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetTargets(conn *hc.Conn) (result *GetTargetsResult, err error) {
	cmd := NewGetTargetsCommand()
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetTargetsCB func(result *GetTargetsResult, err error)

// Retrieves a list of available targets.

type AsyncGetTargetsCommand struct {
	cb GetTargetsCB
}

func NewAsyncGetTargetsCommand(cb GetTargetsCB) *AsyncGetTargetsCommand {
	return &AsyncGetTargetsCommand{
		cb: cb,
	}
}

func (cmd *AsyncGetTargetsCommand) Name() string {
	return "Target.getTargets"
}

func (cmd *AsyncGetTargetsCommand) Params() interface{} {
	return nil
}

func (cmd *GetTargetsCommand) Result() *GetTargetsResult {
	return &cmd.result
}

func (cmd *GetTargetsCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetTargetsCommand) Done(data []byte, err error) {
	var result GetTargetsResult
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

// Issued when a possible inspection target is created.

type TargetCreatedEvent struct {
	TargetInfo *TargetInfo `json:"targetInfo"`
}

func OnTargetCreated(conn *hc.Conn, cb func(evt *TargetCreatedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &TargetCreatedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Target.targetCreated", sink)
}

// Issued when a target is destroyed.

type TargetDestroyedEvent struct {
	TargetId TargetID `json:"targetId"`
}

func OnTargetDestroyed(conn *hc.Conn, cb func(evt *TargetDestroyedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &TargetDestroyedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Target.targetDestroyed", sink)
}

// Issued when attached to target because of auto-attach or attachToTarget command.

type AttachedToTargetEvent struct {
	TargetInfo         *TargetInfo `json:"targetInfo"`
	WaitingForDebugger bool        `json:"waitingForDebugger"`
}

func OnAttachedToTarget(conn *hc.Conn, cb func(evt *AttachedToTargetEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &AttachedToTargetEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Target.attachedToTarget", sink)
}

// Issued when detached from target for any reason (including detachFromTarget command).

type DetachedFromTargetEvent struct {
	TargetId TargetID `json:"targetId"`
}

func OnDetachedFromTarget(conn *hc.Conn, cb func(evt *DetachedFromTargetEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &DetachedFromTargetEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Target.detachedFromTarget", sink)
}

// Notifies about new protocol message from attached target.

type ReceivedMessageFromTargetEvent struct {
	TargetId TargetID `json:"targetId"`
	Message  string   `json:"message"`
}

func OnReceivedMessageFromTarget(conn *hc.Conn, cb func(evt *ReceivedMessageFromTargetEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &ReceivedMessageFromTargetEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Target.receivedMessageFromTarget", sink)
}
