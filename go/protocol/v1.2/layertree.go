package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
)

// Unique Layer identifier.
type LayerId string

// Unique snapshot identifier.
type SnapshotId string

// Rectangle where scrolling happens on the main thread.
type ScrollRect struct {
	Rect *Rect  `json:"rect"` // Rectangle itself.
	Type string `json:"type"` // Reason for rectangle to force scrolling on the main thread
}

// Serialized fragment of layer picture along with its offset within the layer.
type PictureTile struct {
	X       int    `json:"x"`       // Offset from owning layer left boundary
	Y       int    `json:"y"`       // Offset from owning layer top boundary
	Picture string `json:"picture"` // Base64-encoded snapshot data.
}

// Information about a compositing layer.
type Layer struct {
	LayerId       LayerId        `json:"layerId"`       // The unique id for this layer.
	ParentLayerId LayerId        `json:"parentLayerId"` // The id of parent (not present for root).
	BackendNodeId *BackendNodeId `json:"backendNodeId"` // The backend id for the node associated with this layer.
	OffsetX       int            `json:"offsetX"`       // Offset from parent layer, X coordinate.
	OffsetY       int            `json:"offsetY"`       // Offset from parent layer, Y coordinate.
	Width         int            `json:"width"`         // Layer width.
	Height        int            `json:"height"`        // Layer height.
	Transform     []int          `json:"transform"`     // Transformation matrix for layer, default is identity matrix
	AnchorX       int            `json:"anchorX"`       // Transform anchor point X, absent if no transform specified
	AnchorY       int            `json:"anchorY"`       // Transform anchor point Y, absent if no transform specified
	AnchorZ       int            `json:"anchorZ"`       // Transform anchor point Z, absent if no transform specified
	PaintCount    int            `json:"paintCount"`    // Indicates how many time this layer has painted.
	DrawsContent  bool           `json:"drawsContent"`  // Indicates whether this layer hosts any content, rather than being used for transform/scrolling purposes only.
	Invisible     bool           `json:"invisible"`     // Set if layer is not visible.
	ScrollRects   []*ScrollRect  `json:"scrollRects"`   // Rectangles scrolling on main thread only.
}

// Array of timings, one per paint step.
type PaintProfile []int

type LayerTreeEnableCB func(err error)

// Enables compositing tree inspection.
type LayerTreeEnableCommand struct {
	cb LayerTreeEnableCB
}

func NewLayerTreeEnableCommand(cb LayerTreeEnableCB) *LayerTreeEnableCommand {
	return &LayerTreeEnableCommand{
		cb: cb,
	}
}

func (cmd *LayerTreeEnableCommand) Name() string {
	return "LayerTree.enable"
}

func (cmd *LayerTreeEnableCommand) Params() interface{} {
	return nil
}

func (cmd *LayerTreeEnableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type LayerTreeDisableCB func(err error)

// Disables compositing tree inspection.
type LayerTreeDisableCommand struct {
	cb LayerTreeDisableCB
}

func NewLayerTreeDisableCommand(cb LayerTreeDisableCB) *LayerTreeDisableCommand {
	return &LayerTreeDisableCommand{
		cb: cb,
	}
}

func (cmd *LayerTreeDisableCommand) Name() string {
	return "LayerTree.disable"
}

func (cmd *LayerTreeDisableCommand) Params() interface{} {
	return nil
}

func (cmd *LayerTreeDisableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type CompositingReasonsParams struct {
	LayerId LayerId `json:"layerId"` // The id of the layer for which we want to get the reasons it was composited.
}

type CompositingReasonsResult struct {
	CompositingReasons []string `json:"compositingReasons"` // A list of strings specifying reasons for the given layer to become composited.
}

type CompositingReasonsCB func(result *CompositingReasonsResult, err error)

// Provides the reasons why the given layer was composited.
type CompositingReasonsCommand struct {
	params *CompositingReasonsParams
	cb     CompositingReasonsCB
}

func NewCompositingReasonsCommand(params *CompositingReasonsParams, cb CompositingReasonsCB) *CompositingReasonsCommand {
	return &CompositingReasonsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *CompositingReasonsCommand) Name() string {
	return "LayerTree.compositingReasons"
}

func (cmd *CompositingReasonsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CompositingReasonsCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj CompositingReasonsResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type MakeSnapshotParams struct {
	LayerId LayerId `json:"layerId"` // The id of the layer.
}

type MakeSnapshotResult struct {
	SnapshotId SnapshotId `json:"snapshotId"` // The id of the layer snapshot.
}

type MakeSnapshotCB func(result *MakeSnapshotResult, err error)

// Returns the layer snapshot identifier.
type MakeSnapshotCommand struct {
	params *MakeSnapshotParams
	cb     MakeSnapshotCB
}

func NewMakeSnapshotCommand(params *MakeSnapshotParams, cb MakeSnapshotCB) *MakeSnapshotCommand {
	return &MakeSnapshotCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *MakeSnapshotCommand) Name() string {
	return "LayerTree.makeSnapshot"
}

func (cmd *MakeSnapshotCommand) Params() interface{} {
	return cmd.params
}

func (cmd *MakeSnapshotCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj MakeSnapshotResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type LoadSnapshotParams struct {
	Tiles []*PictureTile `json:"tiles"` // An array of tiles composing the snapshot.
}

type LoadSnapshotResult struct {
	SnapshotId SnapshotId `json:"snapshotId"` // The id of the snapshot.
}

type LoadSnapshotCB func(result *LoadSnapshotResult, err error)

// Returns the snapshot identifier.
type LoadSnapshotCommand struct {
	params *LoadSnapshotParams
	cb     LoadSnapshotCB
}

func NewLoadSnapshotCommand(params *LoadSnapshotParams, cb LoadSnapshotCB) *LoadSnapshotCommand {
	return &LoadSnapshotCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *LoadSnapshotCommand) Name() string {
	return "LayerTree.loadSnapshot"
}

func (cmd *LoadSnapshotCommand) Params() interface{} {
	return cmd.params
}

func (cmd *LoadSnapshotCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj LoadSnapshotResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type ReleaseSnapshotParams struct {
	SnapshotId SnapshotId `json:"snapshotId"` // The id of the layer snapshot.
}

type ReleaseSnapshotCB func(err error)

// Releases layer snapshot captured by the back-end.
type ReleaseSnapshotCommand struct {
	params *ReleaseSnapshotParams
	cb     ReleaseSnapshotCB
}

func NewReleaseSnapshotCommand(params *ReleaseSnapshotParams, cb ReleaseSnapshotCB) *ReleaseSnapshotCommand {
	return &ReleaseSnapshotCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *ReleaseSnapshotCommand) Name() string {
	return "LayerTree.releaseSnapshot"
}

func (cmd *ReleaseSnapshotCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ReleaseSnapshotCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type ProfileSnapshotParams struct {
	SnapshotId     SnapshotId `json:"snapshotId"`     // The id of the layer snapshot.
	MinRepeatCount int        `json:"minRepeatCount"` // The maximum number of times to replay the snapshot (1, if not specified).
	MinDuration    int        `json:"minDuration"`    // The minimum duration (in seconds) to replay the snapshot.
	ClipRect       *Rect      `json:"clipRect"`       // The clip rectangle to apply when replaying the snapshot.
}

type ProfileSnapshotResult struct {
	Timings []PaintProfile `json:"timings"` // The array of paint profiles, one per run.
}

type ProfileSnapshotCB func(result *ProfileSnapshotResult, err error)

type ProfileSnapshotCommand struct {
	params *ProfileSnapshotParams
	cb     ProfileSnapshotCB
}

func NewProfileSnapshotCommand(params *ProfileSnapshotParams, cb ProfileSnapshotCB) *ProfileSnapshotCommand {
	return &ProfileSnapshotCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *ProfileSnapshotCommand) Name() string {
	return "LayerTree.profileSnapshot"
}

func (cmd *ProfileSnapshotCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ProfileSnapshotCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj ProfileSnapshotResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type ReplaySnapshotParams struct {
	SnapshotId SnapshotId `json:"snapshotId"` // The id of the layer snapshot.
	FromStep   int        `json:"fromStep"`   // The first step to replay from (replay from the very start if not specified).
	ToStep     int        `json:"toStep"`     // The last step to replay to (replay till the end if not specified).
	Scale      int        `json:"scale"`      // The scale to apply while replaying (defaults to 1).
}

type ReplaySnapshotResult struct {
	DataURL string `json:"dataURL"` // A data: URL for resulting image.
}

type ReplaySnapshotCB func(result *ReplaySnapshotResult, err error)

// Replays the layer snapshot and returns the resulting bitmap.
type ReplaySnapshotCommand struct {
	params *ReplaySnapshotParams
	cb     ReplaySnapshotCB
}

func NewReplaySnapshotCommand(params *ReplaySnapshotParams, cb ReplaySnapshotCB) *ReplaySnapshotCommand {
	return &ReplaySnapshotCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *ReplaySnapshotCommand) Name() string {
	return "LayerTree.replaySnapshot"
}

func (cmd *ReplaySnapshotCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ReplaySnapshotCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj ReplaySnapshotResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type SnapshotCommandLogParams struct {
	SnapshotId SnapshotId `json:"snapshotId"` // The id of the layer snapshot.
}

type SnapshotCommandLogResult struct {
	CommandLog []map[string]string `json:"commandLog"` // The array of canvas function calls.
}

type SnapshotCommandLogCB func(result *SnapshotCommandLogResult, err error)

// Replays the layer snapshot and returns canvas log.
type SnapshotCommandLogCommand struct {
	params *SnapshotCommandLogParams
	cb     SnapshotCommandLogCB
}

func NewSnapshotCommandLogCommand(params *SnapshotCommandLogParams, cb SnapshotCommandLogCB) *SnapshotCommandLogCommand {
	return &SnapshotCommandLogCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SnapshotCommandLogCommand) Name() string {
	return "LayerTree.snapshotCommandLog"
}

func (cmd *SnapshotCommandLogCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SnapshotCommandLogCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj SnapshotCommandLogResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type LayerTreeDidChangeEvent struct {
	Layers []*Layer `json:"layers"` // Layer tree, absent if not in the comspositing mode.
}

type LayerTreeDidChangeEventSink struct {
	events chan *LayerTreeDidChangeEvent
}

func NewLayerTreeDidChangeEventSink(bufSize int) *LayerTreeDidChangeEventSink {
	return &LayerTreeDidChangeEventSink{
		events: make(chan *LayerTreeDidChangeEvent, bufSize),
	}
}

func (s *LayerTreeDidChangeEventSink) Name() string {
	return "LayerTree.layerTreeDidChange"
}

func (s *LayerTreeDidChangeEventSink) OnEvent(params []byte) {
	evt := &LayerTreeDidChangeEvent{}
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

type LayerPaintedEvent struct {
	LayerId LayerId `json:"layerId"` // The id of the painted layer.
	Clip    *Rect   `json:"clip"`    // Clip rectangle.
}

type LayerPaintedEventSink struct {
	events chan *LayerPaintedEvent
}

func NewLayerPaintedEventSink(bufSize int) *LayerPaintedEventSink {
	return &LayerPaintedEventSink{
		events: make(chan *LayerPaintedEvent, bufSize),
	}
}

func (s *LayerPaintedEventSink) Name() string {
	return "LayerTree.layerPainted"
}

func (s *LayerPaintedEventSink) OnEvent(params []byte) {
	evt := &LayerPaintedEvent{}
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
