package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

// Detailed application cache resource information.
type ApplicationCacheResource struct {
	Url  string `json:"url"`  // Resource url.
	Size int    `json:"size"` // Resource size.
	Type string `json:"type"` // Resource type.
}

// Detailed application cache information.
type ApplicationCache struct {
	ManifestURL  string                      `json:"manifestURL"`  // Manifest URL.
	Size         float64                     `json:"size"`         // Application cache size.
	CreationTime float64                     `json:"creationTime"` // Application cache creation time.
	UpdateTime   float64                     `json:"updateTime"`   // Application cache update time.
	Resources    []*ApplicationCacheResource `json:"resources"`    // Application cache resources.
}

// Frame identifier - manifest URL pair.
type FrameWithManifest struct {
	FrameId     *FrameId `json:"frameId"`     // Frame identifier.
	ManifestURL string   `json:"manifestURL"` // Manifest URL.
	Status      int      `json:"status"`      // Application cache status.
}

type GetFramesWithManifestsResult struct {
	FrameIds []*FrameWithManifest `json:"frameIds"` // Array of frame identifiers with manifest urls for each frame containing a document associated with some application cache.
}

// Returns array of frame identifiers with manifest urls for each frame containing a document associated with some application cache.

type GetFramesWithManifestsCommand struct {
	result GetFramesWithManifestsResult
	wg     sync.WaitGroup
	err    error
}

func NewGetFramesWithManifestsCommand() *GetFramesWithManifestsCommand {
	return &GetFramesWithManifestsCommand{}
}

func (cmd *GetFramesWithManifestsCommand) Name() string {
	return "ApplicationCache.getFramesWithManifests"
}

func (cmd *GetFramesWithManifestsCommand) Params() interface{} {
	return nil
}

func (cmd *GetFramesWithManifestsCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetFramesWithManifests(conn *hc.Conn) (result *GetFramesWithManifestsResult, err error) {
	cmd := NewGetFramesWithManifestsCommand()
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetFramesWithManifestsCB func(result *GetFramesWithManifestsResult, err error)

// Returns array of frame identifiers with manifest urls for each frame containing a document associated with some application cache.

type AsyncGetFramesWithManifestsCommand struct {
	cb GetFramesWithManifestsCB
}

func NewAsyncGetFramesWithManifestsCommand(cb GetFramesWithManifestsCB) *AsyncGetFramesWithManifestsCommand {
	return &AsyncGetFramesWithManifestsCommand{
		cb: cb,
	}
}

func (cmd *AsyncGetFramesWithManifestsCommand) Name() string {
	return "ApplicationCache.getFramesWithManifests"
}

func (cmd *AsyncGetFramesWithManifestsCommand) Params() interface{} {
	return nil
}

func (cmd *GetFramesWithManifestsCommand) Result() *GetFramesWithManifestsResult {
	return &cmd.result
}

func (cmd *GetFramesWithManifestsCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetFramesWithManifestsCommand) Done(data []byte, err error) {
	var result GetFramesWithManifestsResult
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

// Enables application cache domain notifications.

type ApplicationCacheEnableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewApplicationCacheEnableCommand() *ApplicationCacheEnableCommand {
	return &ApplicationCacheEnableCommand{}
}

func (cmd *ApplicationCacheEnableCommand) Name() string {
	return "ApplicationCache.enable"
}

func (cmd *ApplicationCacheEnableCommand) Params() interface{} {
	return nil
}

func (cmd *ApplicationCacheEnableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ApplicationCacheEnable(conn *hc.Conn) (err error) {
	cmd := NewApplicationCacheEnableCommand()
	cmd.Run(conn)
	return cmd.err
}

type ApplicationCacheEnableCB func(err error)

// Enables application cache domain notifications.

type AsyncApplicationCacheEnableCommand struct {
	cb ApplicationCacheEnableCB
}

func NewAsyncApplicationCacheEnableCommand(cb ApplicationCacheEnableCB) *AsyncApplicationCacheEnableCommand {
	return &AsyncApplicationCacheEnableCommand{
		cb: cb,
	}
}

func (cmd *AsyncApplicationCacheEnableCommand) Name() string {
	return "ApplicationCache.enable"
}

func (cmd *AsyncApplicationCacheEnableCommand) Params() interface{} {
	return nil
}

func (cmd *ApplicationCacheEnableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncApplicationCacheEnableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type GetManifestForFrameParams struct {
	FrameId *FrameId `json:"frameId"` // Identifier of the frame containing document whose manifest is retrieved.
}

type GetManifestForFrameResult struct {
	ManifestURL string `json:"manifestURL"` // Manifest URL for document in the given frame.
}

// Returns manifest URL for document in the given frame.

type GetManifestForFrameCommand struct {
	params *GetManifestForFrameParams
	result GetManifestForFrameResult
	wg     sync.WaitGroup
	err    error
}

func NewGetManifestForFrameCommand(params *GetManifestForFrameParams) *GetManifestForFrameCommand {
	return &GetManifestForFrameCommand{
		params: params,
	}
}

func (cmd *GetManifestForFrameCommand) Name() string {
	return "ApplicationCache.getManifestForFrame"
}

func (cmd *GetManifestForFrameCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetManifestForFrameCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetManifestForFrame(params *GetManifestForFrameParams, conn *hc.Conn) (result *GetManifestForFrameResult, err error) {
	cmd := NewGetManifestForFrameCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetManifestForFrameCB func(result *GetManifestForFrameResult, err error)

// Returns manifest URL for document in the given frame.

type AsyncGetManifestForFrameCommand struct {
	params *GetManifestForFrameParams
	cb     GetManifestForFrameCB
}

func NewAsyncGetManifestForFrameCommand(params *GetManifestForFrameParams, cb GetManifestForFrameCB) *AsyncGetManifestForFrameCommand {
	return &AsyncGetManifestForFrameCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetManifestForFrameCommand) Name() string {
	return "ApplicationCache.getManifestForFrame"
}

func (cmd *AsyncGetManifestForFrameCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetManifestForFrameCommand) Result() *GetManifestForFrameResult {
	return &cmd.result
}

func (cmd *GetManifestForFrameCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetManifestForFrameCommand) Done(data []byte, err error) {
	var result GetManifestForFrameResult
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

type GetApplicationCacheForFrameParams struct {
	FrameId *FrameId `json:"frameId"` // Identifier of the frame containing document whose application cache is retrieved.
}

type GetApplicationCacheForFrameResult struct {
	ApplicationCache *ApplicationCache `json:"applicationCache"` // Relevant application cache data for the document in given frame.
}

// Returns relevant application cache data for the document in given frame.

type GetApplicationCacheForFrameCommand struct {
	params *GetApplicationCacheForFrameParams
	result GetApplicationCacheForFrameResult
	wg     sync.WaitGroup
	err    error
}

func NewGetApplicationCacheForFrameCommand(params *GetApplicationCacheForFrameParams) *GetApplicationCacheForFrameCommand {
	return &GetApplicationCacheForFrameCommand{
		params: params,
	}
}

func (cmd *GetApplicationCacheForFrameCommand) Name() string {
	return "ApplicationCache.getApplicationCacheForFrame"
}

func (cmd *GetApplicationCacheForFrameCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetApplicationCacheForFrameCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetApplicationCacheForFrame(params *GetApplicationCacheForFrameParams, conn *hc.Conn) (result *GetApplicationCacheForFrameResult, err error) {
	cmd := NewGetApplicationCacheForFrameCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetApplicationCacheForFrameCB func(result *GetApplicationCacheForFrameResult, err error)

// Returns relevant application cache data for the document in given frame.

type AsyncGetApplicationCacheForFrameCommand struct {
	params *GetApplicationCacheForFrameParams
	cb     GetApplicationCacheForFrameCB
}

func NewAsyncGetApplicationCacheForFrameCommand(params *GetApplicationCacheForFrameParams, cb GetApplicationCacheForFrameCB) *AsyncGetApplicationCacheForFrameCommand {
	return &AsyncGetApplicationCacheForFrameCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetApplicationCacheForFrameCommand) Name() string {
	return "ApplicationCache.getApplicationCacheForFrame"
}

func (cmd *AsyncGetApplicationCacheForFrameCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetApplicationCacheForFrameCommand) Result() *GetApplicationCacheForFrameResult {
	return &cmd.result
}

func (cmd *GetApplicationCacheForFrameCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetApplicationCacheForFrameCommand) Done(data []byte, err error) {
	var result GetApplicationCacheForFrameResult
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

type ApplicationCacheStatusUpdatedEvent struct {
	FrameId     *FrameId `json:"frameId"`     // Identifier of the frame containing document whose application cache updated status.
	ManifestURL string   `json:"manifestURL"` // Manifest URL.
	Status      int      `json:"status"`      // Updated application cache status.
}

func OnApplicationCacheStatusUpdated(conn *hc.Conn, cb func(evt *ApplicationCacheStatusUpdatedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &ApplicationCacheStatusUpdatedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("ApplicationCache.applicationCacheStatusUpdated", sink)
}

type NetworkStateUpdatedEvent struct {
	IsNowOnline bool `json:"isNowOnline"`
}

func OnNetworkStateUpdated(conn *hc.Conn, cb func(evt *NetworkStateUpdatedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &NetworkStateUpdatedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("ApplicationCache.networkStateUpdated", sink)
}
