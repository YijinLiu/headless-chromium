package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
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
	X       float64 `json:"x"`       // Offset from owning layer left boundary
	Y       float64 `json:"y"`       // Offset from owning layer top boundary
	Picture string  `json:"picture"` // Base64-encoded snapshot data.
}

// Information about a compositing layer.
type Layer struct {
	LayerId       LayerId        `json:"layerId"`                 // The unique id for this layer.
	ParentLayerId LayerId        `json:"parentLayerId,omitempty"` // The id of parent (not present for root).
	BackendNodeId *BackendNodeId `json:"backendNodeId,omitempty"` // The backend id for the node associated with this layer.
	OffsetX       float64        `json:"offsetX"`                 // Offset from parent layer, X coordinate.
	OffsetY       float64        `json:"offsetY"`                 // Offset from parent layer, Y coordinate.
	Width         float64        `json:"width"`                   // Layer width.
	Height        float64        `json:"height"`                  // Layer height.
	Transform     []float64      `json:"transform,omitempty"`     // Transformation matrix for layer, default is identity matrix
	AnchorX       float64        `json:"anchorX,omitempty"`       // Transform anchor point X, absent if no transform specified
	AnchorY       float64        `json:"anchorY,omitempty"`       // Transform anchor point Y, absent if no transform specified
	AnchorZ       float64        `json:"anchorZ,omitempty"`       // Transform anchor point Z, absent if no transform specified
	PaintCount    int            `json:"paintCount"`              // Indicates how many time this layer has painted.
	DrawsContent  bool           `json:"drawsContent"`            // Indicates whether this layer hosts any content, rather than being used for transform/scrolling purposes only.
	Invisible     bool           `json:"invisible,omitempty"`     // Set if layer is not visible.
	ScrollRects   []*ScrollRect  `json:"scrollRects,omitempty"`   // Rectangles scrolling on main thread only.
}

// Array of timings, one per paint step.
type PaintProfile []float64

// Enables compositing tree inspection.

type LayerTreeEnableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewLayerTreeEnableCommand() *LayerTreeEnableCommand {
	return &LayerTreeEnableCommand{}
}

func (cmd *LayerTreeEnableCommand) Name() string {
	return "LayerTree.enable"
}

func (cmd *LayerTreeEnableCommand) Params() interface{} {
	return nil
}

func (cmd *LayerTreeEnableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func LayerTreeEnable(conn *hc.Conn) (err error) {
	cmd := NewLayerTreeEnableCommand()
	cmd.Run(conn)
	return cmd.err
}

type LayerTreeEnableCB func(err error)

// Enables compositing tree inspection.

type AsyncLayerTreeEnableCommand struct {
	cb LayerTreeEnableCB
}

func NewAsyncLayerTreeEnableCommand(cb LayerTreeEnableCB) *AsyncLayerTreeEnableCommand {
	return &AsyncLayerTreeEnableCommand{
		cb: cb,
	}
}

func (cmd *AsyncLayerTreeEnableCommand) Name() string {
	return "LayerTree.enable"
}

func (cmd *AsyncLayerTreeEnableCommand) Params() interface{} {
	return nil
}

func (cmd *LayerTreeEnableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncLayerTreeEnableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Disables compositing tree inspection.

type LayerTreeDisableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewLayerTreeDisableCommand() *LayerTreeDisableCommand {
	return &LayerTreeDisableCommand{}
}

func (cmd *LayerTreeDisableCommand) Name() string {
	return "LayerTree.disable"
}

func (cmd *LayerTreeDisableCommand) Params() interface{} {
	return nil
}

func (cmd *LayerTreeDisableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func LayerTreeDisable(conn *hc.Conn) (err error) {
	cmd := NewLayerTreeDisableCommand()
	cmd.Run(conn)
	return cmd.err
}

type LayerTreeDisableCB func(err error)

// Disables compositing tree inspection.

type AsyncLayerTreeDisableCommand struct {
	cb LayerTreeDisableCB
}

func NewAsyncLayerTreeDisableCommand(cb LayerTreeDisableCB) *AsyncLayerTreeDisableCommand {
	return &AsyncLayerTreeDisableCommand{
		cb: cb,
	}
}

func (cmd *AsyncLayerTreeDisableCommand) Name() string {
	return "LayerTree.disable"
}

func (cmd *AsyncLayerTreeDisableCommand) Params() interface{} {
	return nil
}

func (cmd *LayerTreeDisableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncLayerTreeDisableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type CompositingReasonsParams struct {
	LayerId LayerId `json:"layerId"` // The id of the layer for which we want to get the reasons it was composited.
}

type CompositingReasonsResult struct {
	CompositingReasons []string `json:"compositingReasons"` // A list of strings specifying reasons for the given layer to become composited.
}

// Provides the reasons why the given layer was composited.

type CompositingReasonsCommand struct {
	params *CompositingReasonsParams
	result CompositingReasonsResult
	wg     sync.WaitGroup
	err    error
}

func NewCompositingReasonsCommand(params *CompositingReasonsParams) *CompositingReasonsCommand {
	return &CompositingReasonsCommand{
		params: params,
	}
}

func (cmd *CompositingReasonsCommand) Name() string {
	return "LayerTree.compositingReasons"
}

func (cmd *CompositingReasonsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CompositingReasonsCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func CompositingReasons(params *CompositingReasonsParams, conn *hc.Conn) (result *CompositingReasonsResult, err error) {
	cmd := NewCompositingReasonsCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type CompositingReasonsCB func(result *CompositingReasonsResult, err error)

// Provides the reasons why the given layer was composited.

type AsyncCompositingReasonsCommand struct {
	params *CompositingReasonsParams
	cb     CompositingReasonsCB
}

func NewAsyncCompositingReasonsCommand(params *CompositingReasonsParams, cb CompositingReasonsCB) *AsyncCompositingReasonsCommand {
	return &AsyncCompositingReasonsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncCompositingReasonsCommand) Name() string {
	return "LayerTree.compositingReasons"
}

func (cmd *AsyncCompositingReasonsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CompositingReasonsCommand) Result() *CompositingReasonsResult {
	return &cmd.result
}

func (cmd *CompositingReasonsCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncCompositingReasonsCommand) Done(data []byte, err error) {
	var result CompositingReasonsResult
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

type MakeSnapshotParams struct {
	LayerId LayerId `json:"layerId"` // The id of the layer.
}

type MakeSnapshotResult struct {
	SnapshotId SnapshotId `json:"snapshotId"` // The id of the layer snapshot.
}

// Returns the layer snapshot identifier.

type MakeSnapshotCommand struct {
	params *MakeSnapshotParams
	result MakeSnapshotResult
	wg     sync.WaitGroup
	err    error
}

func NewMakeSnapshotCommand(params *MakeSnapshotParams) *MakeSnapshotCommand {
	return &MakeSnapshotCommand{
		params: params,
	}
}

func (cmd *MakeSnapshotCommand) Name() string {
	return "LayerTree.makeSnapshot"
}

func (cmd *MakeSnapshotCommand) Params() interface{} {
	return cmd.params
}

func (cmd *MakeSnapshotCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func MakeSnapshot(params *MakeSnapshotParams, conn *hc.Conn) (result *MakeSnapshotResult, err error) {
	cmd := NewMakeSnapshotCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type MakeSnapshotCB func(result *MakeSnapshotResult, err error)

// Returns the layer snapshot identifier.

type AsyncMakeSnapshotCommand struct {
	params *MakeSnapshotParams
	cb     MakeSnapshotCB
}

func NewAsyncMakeSnapshotCommand(params *MakeSnapshotParams, cb MakeSnapshotCB) *AsyncMakeSnapshotCommand {
	return &AsyncMakeSnapshotCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncMakeSnapshotCommand) Name() string {
	return "LayerTree.makeSnapshot"
}

func (cmd *AsyncMakeSnapshotCommand) Params() interface{} {
	return cmd.params
}

func (cmd *MakeSnapshotCommand) Result() *MakeSnapshotResult {
	return &cmd.result
}

func (cmd *MakeSnapshotCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncMakeSnapshotCommand) Done(data []byte, err error) {
	var result MakeSnapshotResult
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

type LoadSnapshotParams struct {
	Tiles []*PictureTile `json:"tiles"` // An array of tiles composing the snapshot.
}

type LoadSnapshotResult struct {
	SnapshotId SnapshotId `json:"snapshotId"` // The id of the snapshot.
}

// Returns the snapshot identifier.

type LoadSnapshotCommand struct {
	params *LoadSnapshotParams
	result LoadSnapshotResult
	wg     sync.WaitGroup
	err    error
}

func NewLoadSnapshotCommand(params *LoadSnapshotParams) *LoadSnapshotCommand {
	return &LoadSnapshotCommand{
		params: params,
	}
}

func (cmd *LoadSnapshotCommand) Name() string {
	return "LayerTree.loadSnapshot"
}

func (cmd *LoadSnapshotCommand) Params() interface{} {
	return cmd.params
}

func (cmd *LoadSnapshotCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func LoadSnapshot(params *LoadSnapshotParams, conn *hc.Conn) (result *LoadSnapshotResult, err error) {
	cmd := NewLoadSnapshotCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type LoadSnapshotCB func(result *LoadSnapshotResult, err error)

// Returns the snapshot identifier.

type AsyncLoadSnapshotCommand struct {
	params *LoadSnapshotParams
	cb     LoadSnapshotCB
}

func NewAsyncLoadSnapshotCommand(params *LoadSnapshotParams, cb LoadSnapshotCB) *AsyncLoadSnapshotCommand {
	return &AsyncLoadSnapshotCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncLoadSnapshotCommand) Name() string {
	return "LayerTree.loadSnapshot"
}

func (cmd *AsyncLoadSnapshotCommand) Params() interface{} {
	return cmd.params
}

func (cmd *LoadSnapshotCommand) Result() *LoadSnapshotResult {
	return &cmd.result
}

func (cmd *LoadSnapshotCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncLoadSnapshotCommand) Done(data []byte, err error) {
	var result LoadSnapshotResult
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

type ReleaseSnapshotParams struct {
	SnapshotId SnapshotId `json:"snapshotId"` // The id of the layer snapshot.
}

// Releases layer snapshot captured by the back-end.

type ReleaseSnapshotCommand struct {
	params *ReleaseSnapshotParams
	wg     sync.WaitGroup
	err    error
}

func NewReleaseSnapshotCommand(params *ReleaseSnapshotParams) *ReleaseSnapshotCommand {
	return &ReleaseSnapshotCommand{
		params: params,
	}
}

func (cmd *ReleaseSnapshotCommand) Name() string {
	return "LayerTree.releaseSnapshot"
}

func (cmd *ReleaseSnapshotCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ReleaseSnapshotCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ReleaseSnapshot(params *ReleaseSnapshotParams, conn *hc.Conn) (err error) {
	cmd := NewReleaseSnapshotCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type ReleaseSnapshotCB func(err error)

// Releases layer snapshot captured by the back-end.

type AsyncReleaseSnapshotCommand struct {
	params *ReleaseSnapshotParams
	cb     ReleaseSnapshotCB
}

func NewAsyncReleaseSnapshotCommand(params *ReleaseSnapshotParams, cb ReleaseSnapshotCB) *AsyncReleaseSnapshotCommand {
	return &AsyncReleaseSnapshotCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncReleaseSnapshotCommand) Name() string {
	return "LayerTree.releaseSnapshot"
}

func (cmd *AsyncReleaseSnapshotCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ReleaseSnapshotCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncReleaseSnapshotCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type ProfileSnapshotParams struct {
	SnapshotId     SnapshotId `json:"snapshotId"`               // The id of the layer snapshot.
	MinRepeatCount int        `json:"minRepeatCount,omitempty"` // The maximum number of times to replay the snapshot (1, if not specified).
	MinDuration    float64    `json:"minDuration,omitempty"`    // The minimum duration (in seconds) to replay the snapshot.
	ClipRect       *Rect      `json:"clipRect,omitempty"`       // The clip rectangle to apply when replaying the snapshot.
}

type ProfileSnapshotResult struct {
	Timings []PaintProfile `json:"timings"` // The array of paint profiles, one per run.
}

type ProfileSnapshotCommand struct {
	params *ProfileSnapshotParams
	result ProfileSnapshotResult
	wg     sync.WaitGroup
	err    error
}

func NewProfileSnapshotCommand(params *ProfileSnapshotParams) *ProfileSnapshotCommand {
	return &ProfileSnapshotCommand{
		params: params,
	}
}

func (cmd *ProfileSnapshotCommand) Name() string {
	return "LayerTree.profileSnapshot"
}

func (cmd *ProfileSnapshotCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ProfileSnapshotCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ProfileSnapshot(params *ProfileSnapshotParams, conn *hc.Conn) (result *ProfileSnapshotResult, err error) {
	cmd := NewProfileSnapshotCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type ProfileSnapshotCB func(result *ProfileSnapshotResult, err error)

type AsyncProfileSnapshotCommand struct {
	params *ProfileSnapshotParams
	cb     ProfileSnapshotCB
}

func NewAsyncProfileSnapshotCommand(params *ProfileSnapshotParams, cb ProfileSnapshotCB) *AsyncProfileSnapshotCommand {
	return &AsyncProfileSnapshotCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncProfileSnapshotCommand) Name() string {
	return "LayerTree.profileSnapshot"
}

func (cmd *AsyncProfileSnapshotCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ProfileSnapshotCommand) Result() *ProfileSnapshotResult {
	return &cmd.result
}

func (cmd *ProfileSnapshotCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncProfileSnapshotCommand) Done(data []byte, err error) {
	var result ProfileSnapshotResult
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

type ReplaySnapshotParams struct {
	SnapshotId SnapshotId `json:"snapshotId"`         // The id of the layer snapshot.
	FromStep   int        `json:"fromStep,omitempty"` // The first step to replay from (replay from the very start if not specified).
	ToStep     int        `json:"toStep,omitempty"`   // The last step to replay to (replay till the end if not specified).
	Scale      float64    `json:"scale,omitempty"`    // The scale to apply while replaying (defaults to 1).
}

type ReplaySnapshotResult struct {
	DataURL string `json:"dataURL"` // A data: URL for resulting image.
}

// Replays the layer snapshot and returns the resulting bitmap.

type ReplaySnapshotCommand struct {
	params *ReplaySnapshotParams
	result ReplaySnapshotResult
	wg     sync.WaitGroup
	err    error
}

func NewReplaySnapshotCommand(params *ReplaySnapshotParams) *ReplaySnapshotCommand {
	return &ReplaySnapshotCommand{
		params: params,
	}
}

func (cmd *ReplaySnapshotCommand) Name() string {
	return "LayerTree.replaySnapshot"
}

func (cmd *ReplaySnapshotCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ReplaySnapshotCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ReplaySnapshot(params *ReplaySnapshotParams, conn *hc.Conn) (result *ReplaySnapshotResult, err error) {
	cmd := NewReplaySnapshotCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type ReplaySnapshotCB func(result *ReplaySnapshotResult, err error)

// Replays the layer snapshot and returns the resulting bitmap.

type AsyncReplaySnapshotCommand struct {
	params *ReplaySnapshotParams
	cb     ReplaySnapshotCB
}

func NewAsyncReplaySnapshotCommand(params *ReplaySnapshotParams, cb ReplaySnapshotCB) *AsyncReplaySnapshotCommand {
	return &AsyncReplaySnapshotCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncReplaySnapshotCommand) Name() string {
	return "LayerTree.replaySnapshot"
}

func (cmd *AsyncReplaySnapshotCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ReplaySnapshotCommand) Result() *ReplaySnapshotResult {
	return &cmd.result
}

func (cmd *ReplaySnapshotCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncReplaySnapshotCommand) Done(data []byte, err error) {
	var result ReplaySnapshotResult
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

type SnapshotCommandLogParams struct {
	SnapshotId SnapshotId `json:"snapshotId"` // The id of the layer snapshot.
}

type SnapshotCommandLogResult struct {
	CommandLog []map[string]string `json:"commandLog"` // The array of canvas function calls.
}

// Replays the layer snapshot and returns canvas log.

type SnapshotCommandLogCommand struct {
	params *SnapshotCommandLogParams
	result SnapshotCommandLogResult
	wg     sync.WaitGroup
	err    error
}

func NewSnapshotCommandLogCommand(params *SnapshotCommandLogParams) *SnapshotCommandLogCommand {
	return &SnapshotCommandLogCommand{
		params: params,
	}
}

func (cmd *SnapshotCommandLogCommand) Name() string {
	return "LayerTree.snapshotCommandLog"
}

func (cmd *SnapshotCommandLogCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SnapshotCommandLogCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SnapshotCommandLog(params *SnapshotCommandLogParams, conn *hc.Conn) (result *SnapshotCommandLogResult, err error) {
	cmd := NewSnapshotCommandLogCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type SnapshotCommandLogCB func(result *SnapshotCommandLogResult, err error)

// Replays the layer snapshot and returns canvas log.

type AsyncSnapshotCommandLogCommand struct {
	params *SnapshotCommandLogParams
	cb     SnapshotCommandLogCB
}

func NewAsyncSnapshotCommandLogCommand(params *SnapshotCommandLogParams, cb SnapshotCommandLogCB) *AsyncSnapshotCommandLogCommand {
	return &AsyncSnapshotCommandLogCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSnapshotCommandLogCommand) Name() string {
	return "LayerTree.snapshotCommandLog"
}

func (cmd *AsyncSnapshotCommandLogCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SnapshotCommandLogCommand) Result() *SnapshotCommandLogResult {
	return &cmd.result
}

func (cmd *SnapshotCommandLogCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSnapshotCommandLogCommand) Done(data []byte, err error) {
	var result SnapshotCommandLogResult
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

type LayerTreeDidChangeEvent struct {
	Layers []*Layer `json:"layers"` // Layer tree, absent if not in the comspositing mode.
}

func OnLayerTreeDidChange(conn *hc.Conn, cb func(evt *LayerTreeDidChangeEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &LayerTreeDidChangeEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("LayerTree.layerTreeDidChange", sink)
}

type LayerPaintedEvent struct {
	LayerId LayerId `json:"layerId"` // The id of the painted layer.
	Clip    *Rect   `json:"clip"`    // Clip rectangle.
}

func OnLayerPainted(conn *hc.Conn, cb func(evt *LayerPaintedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &LayerPaintedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("LayerTree.layerPainted", sink)
}
