package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
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
	Size         int                         `json:"size"`         // Application cache size.
	CreationTime int                         `json:"creationTime"` // Application cache creation time.
	UpdateTime   int                         `json:"updateTime"`   // Application cache update time.
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

type GetFramesWithManifestsCB func(result *GetFramesWithManifestsResult, err error)

// Returns array of frame identifiers with manifest urls for each frame containing a document associated with some application cache.
type GetFramesWithManifestsCommand struct {
	cb GetFramesWithManifestsCB
}

func NewGetFramesWithManifestsCommand(cb GetFramesWithManifestsCB) *GetFramesWithManifestsCommand {
	return &GetFramesWithManifestsCommand{
		cb: cb,
	}
}

func (cmd *GetFramesWithManifestsCommand) Name() string {
	return "ApplicationCache.getFramesWithManifests"
}

func (cmd *GetFramesWithManifestsCommand) Params() interface{} {
	return nil
}

func (cmd *GetFramesWithManifestsCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetFramesWithManifestsResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type ApplicationCacheEnableCB func(err error)

// Enables application cache domain notifications.
type ApplicationCacheEnableCommand struct {
	cb ApplicationCacheEnableCB
}

func NewApplicationCacheEnableCommand(cb ApplicationCacheEnableCB) *ApplicationCacheEnableCommand {
	return &ApplicationCacheEnableCommand{
		cb: cb,
	}
}

func (cmd *ApplicationCacheEnableCommand) Name() string {
	return "ApplicationCache.enable"
}

func (cmd *ApplicationCacheEnableCommand) Params() interface{} {
	return nil
}

func (cmd *ApplicationCacheEnableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type GetManifestForFrameParams struct {
	FrameId *FrameId `json:"frameId"` // Identifier of the frame containing document whose manifest is retrieved.
}

type GetManifestForFrameResult struct {
	ManifestURL string `json:"manifestURL"` // Manifest URL for document in the given frame.
}

type GetManifestForFrameCB func(result *GetManifestForFrameResult, err error)

// Returns manifest URL for document in the given frame.
type GetManifestForFrameCommand struct {
	params *GetManifestForFrameParams
	cb     GetManifestForFrameCB
}

func NewGetManifestForFrameCommand(params *GetManifestForFrameParams, cb GetManifestForFrameCB) *GetManifestForFrameCommand {
	return &GetManifestForFrameCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetManifestForFrameCommand) Name() string {
	return "ApplicationCache.getManifestForFrame"
}

func (cmd *GetManifestForFrameCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetManifestForFrameCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetManifestForFrameResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type GetApplicationCacheForFrameParams struct {
	FrameId *FrameId `json:"frameId"` // Identifier of the frame containing document whose application cache is retrieved.
}

type GetApplicationCacheForFrameResult struct {
	ApplicationCache *ApplicationCache `json:"applicationCache"` // Relevant application cache data for the document in given frame.
}

type GetApplicationCacheForFrameCB func(result *GetApplicationCacheForFrameResult, err error)

// Returns relevant application cache data for the document in given frame.
type GetApplicationCacheForFrameCommand struct {
	params *GetApplicationCacheForFrameParams
	cb     GetApplicationCacheForFrameCB
}

func NewGetApplicationCacheForFrameCommand(params *GetApplicationCacheForFrameParams, cb GetApplicationCacheForFrameCB) *GetApplicationCacheForFrameCommand {
	return &GetApplicationCacheForFrameCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetApplicationCacheForFrameCommand) Name() string {
	return "ApplicationCache.getApplicationCacheForFrame"
}

func (cmd *GetApplicationCacheForFrameCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetApplicationCacheForFrameCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetApplicationCacheForFrameResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type ApplicationCacheStatusUpdatedEvent struct {
	FrameId     *FrameId `json:"frameId"`     // Identifier of the frame containing document whose application cache updated status.
	ManifestURL string   `json:"manifestURL"` // Manifest URL.
	Status      int      `json:"status"`      // Updated application cache status.
}

type ApplicationCacheStatusUpdatedEventSink struct {
	events chan *ApplicationCacheStatusUpdatedEvent
}

func NewApplicationCacheStatusUpdatedEventSink(bufSize int) *ApplicationCacheStatusUpdatedEventSink {
	return &ApplicationCacheStatusUpdatedEventSink{
		events: make(chan *ApplicationCacheStatusUpdatedEvent, bufSize),
	}
}

func (s *ApplicationCacheStatusUpdatedEventSink) Name() string {
	return "ApplicationCache.applicationCacheStatusUpdated"
}

func (s *ApplicationCacheStatusUpdatedEventSink) OnEvent(params []byte) {
	evt := &ApplicationCacheStatusUpdatedEvent{}
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

type NetworkStateUpdatedEvent struct {
	IsNowOnline bool `json:"isNowOnline"`
}

type NetworkStateUpdatedEventSink struct {
	events chan *NetworkStateUpdatedEvent
}

func NewNetworkStateUpdatedEventSink(bufSize int) *NetworkStateUpdatedEventSink {
	return &NetworkStateUpdatedEventSink{
		events: make(chan *NetworkStateUpdatedEvent, bufSize),
	}
}

func (s *NetworkStateUpdatedEventSink) Name() string {
	return "ApplicationCache.networkStateUpdated"
}

func (s *NetworkStateUpdatedEventSink) OnEvent(params []byte) {
	evt := &NetworkStateUpdatedEvent{}
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
