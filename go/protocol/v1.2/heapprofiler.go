package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
)

// Heap snapshot object id.
type HeapSnapshotObjectId string

// Sampling Heap Profile node. Holds callsite information, allocation statistics and child nodes.
type SamplingHeapProfileNode struct {
	CallFrame *RuntimeCallFrame          `json:"callFrame"` // Function location.
	SelfSize  int                        `json:"selfSize"`  // Allocations size in bytes for the node excluding children.
	Children  []*SamplingHeapProfileNode `json:"children"`  // Child nodes.
}

// Profile.
type SamplingHeapProfile struct {
	Head *SamplingHeapProfileNode `json:"head"`
}

type HeapProfilerEnableCB func(err error)

type HeapProfilerEnableCommand struct {
	cb HeapProfilerEnableCB
}

func NewHeapProfilerEnableCommand(cb HeapProfilerEnableCB) *HeapProfilerEnableCommand {
	return &HeapProfilerEnableCommand{
		cb: cb,
	}
}

func (cmd *HeapProfilerEnableCommand) Name() string {
	return "HeapProfiler.enable"
}

func (cmd *HeapProfilerEnableCommand) Params() interface{} {
	return nil
}

func (cmd *HeapProfilerEnableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type HeapProfilerDisableCB func(err error)

type HeapProfilerDisableCommand struct {
	cb HeapProfilerDisableCB
}

func NewHeapProfilerDisableCommand(cb HeapProfilerDisableCB) *HeapProfilerDisableCommand {
	return &HeapProfilerDisableCommand{
		cb: cb,
	}
}

func (cmd *HeapProfilerDisableCommand) Name() string {
	return "HeapProfiler.disable"
}

func (cmd *HeapProfilerDisableCommand) Params() interface{} {
	return nil
}

func (cmd *HeapProfilerDisableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type StartTrackingHeapObjectsParams struct {
	TrackAllocations bool `json:"trackAllocations"`
}

type StartTrackingHeapObjectsCB func(err error)

type StartTrackingHeapObjectsCommand struct {
	params *StartTrackingHeapObjectsParams
	cb     StartTrackingHeapObjectsCB
}

func NewStartTrackingHeapObjectsCommand(params *StartTrackingHeapObjectsParams, cb StartTrackingHeapObjectsCB) *StartTrackingHeapObjectsCommand {
	return &StartTrackingHeapObjectsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *StartTrackingHeapObjectsCommand) Name() string {
	return "HeapProfiler.startTrackingHeapObjects"
}

func (cmd *StartTrackingHeapObjectsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *StartTrackingHeapObjectsCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type StopTrackingHeapObjectsParams struct {
	ReportProgress bool `json:"reportProgress"` // If true 'reportHeapSnapshotProgress' events will be generated while snapshot is being taken when the tracking is stopped.
}

type StopTrackingHeapObjectsCB func(err error)

type StopTrackingHeapObjectsCommand struct {
	params *StopTrackingHeapObjectsParams
	cb     StopTrackingHeapObjectsCB
}

func NewStopTrackingHeapObjectsCommand(params *StopTrackingHeapObjectsParams, cb StopTrackingHeapObjectsCB) *StopTrackingHeapObjectsCommand {
	return &StopTrackingHeapObjectsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *StopTrackingHeapObjectsCommand) Name() string {
	return "HeapProfiler.stopTrackingHeapObjects"
}

func (cmd *StopTrackingHeapObjectsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *StopTrackingHeapObjectsCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type TakeHeapSnapshotParams struct {
	ReportProgress bool `json:"reportProgress"` // If true 'reportHeapSnapshotProgress' events will be generated while snapshot is being taken.
}

type TakeHeapSnapshotCB func(err error)

type TakeHeapSnapshotCommand struct {
	params *TakeHeapSnapshotParams
	cb     TakeHeapSnapshotCB
}

func NewTakeHeapSnapshotCommand(params *TakeHeapSnapshotParams, cb TakeHeapSnapshotCB) *TakeHeapSnapshotCommand {
	return &TakeHeapSnapshotCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *TakeHeapSnapshotCommand) Name() string {
	return "HeapProfiler.takeHeapSnapshot"
}

func (cmd *TakeHeapSnapshotCommand) Params() interface{} {
	return cmd.params
}

func (cmd *TakeHeapSnapshotCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type CollectGarbageCB func(err error)

type CollectGarbageCommand struct {
	cb CollectGarbageCB
}

func NewCollectGarbageCommand(cb CollectGarbageCB) *CollectGarbageCommand {
	return &CollectGarbageCommand{
		cb: cb,
	}
}

func (cmd *CollectGarbageCommand) Name() string {
	return "HeapProfiler.collectGarbage"
}

func (cmd *CollectGarbageCommand) Params() interface{} {
	return nil
}

func (cmd *CollectGarbageCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type GetObjectByHeapObjectIdParams struct {
	ObjectId    HeapSnapshotObjectId `json:"objectId"`
	ObjectGroup string               `json:"objectGroup"` // Symbolic group name that can be used to release multiple objects.
}

type GetObjectByHeapObjectIdResult struct {
	Result *RemoteObject `json:"result"` // Evaluation result.
}

type GetObjectByHeapObjectIdCB func(result *GetObjectByHeapObjectIdResult, err error)

type GetObjectByHeapObjectIdCommand struct {
	params *GetObjectByHeapObjectIdParams
	cb     GetObjectByHeapObjectIdCB
}

func NewGetObjectByHeapObjectIdCommand(params *GetObjectByHeapObjectIdParams, cb GetObjectByHeapObjectIdCB) *GetObjectByHeapObjectIdCommand {
	return &GetObjectByHeapObjectIdCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetObjectByHeapObjectIdCommand) Name() string {
	return "HeapProfiler.getObjectByHeapObjectId"
}

func (cmd *GetObjectByHeapObjectIdCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetObjectByHeapObjectIdCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetObjectByHeapObjectIdResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type AddInspectedHeapObjectParams struct {
	HeapObjectId HeapSnapshotObjectId `json:"heapObjectId"` // Heap snapshot object id to be accessible by means of $x command line API.
}

type AddInspectedHeapObjectCB func(err error)

// Enables console to refer to the node with given id via $x (see Command Line API for more details $x functions).
type AddInspectedHeapObjectCommand struct {
	params *AddInspectedHeapObjectParams
	cb     AddInspectedHeapObjectCB
}

func NewAddInspectedHeapObjectCommand(params *AddInspectedHeapObjectParams, cb AddInspectedHeapObjectCB) *AddInspectedHeapObjectCommand {
	return &AddInspectedHeapObjectCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AddInspectedHeapObjectCommand) Name() string {
	return "HeapProfiler.addInspectedHeapObject"
}

func (cmd *AddInspectedHeapObjectCommand) Params() interface{} {
	return cmd.params
}

func (cmd *AddInspectedHeapObjectCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type GetHeapObjectIdParams struct {
	ObjectId *RemoteObjectId `json:"objectId"` // Identifier of the object to get heap object id for.
}

type GetHeapObjectIdResult struct {
	HeapSnapshotObjectId HeapSnapshotObjectId `json:"heapSnapshotObjectId"` // Id of the heap snapshot object corresponding to the passed remote object id.
}

type GetHeapObjectIdCB func(result *GetHeapObjectIdResult, err error)

type GetHeapObjectIdCommand struct {
	params *GetHeapObjectIdParams
	cb     GetHeapObjectIdCB
}

func NewGetHeapObjectIdCommand(params *GetHeapObjectIdParams, cb GetHeapObjectIdCB) *GetHeapObjectIdCommand {
	return &GetHeapObjectIdCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetHeapObjectIdCommand) Name() string {
	return "HeapProfiler.getHeapObjectId"
}

func (cmd *GetHeapObjectIdCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetHeapObjectIdCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetHeapObjectIdResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type StartSamplingParams struct {
	SamplingInterval int `json:"samplingInterval"` // Average sample interval in bytes. Poisson distribution is used for the intervals. The default value is 32768 bytes.
}

type StartSamplingCB func(err error)

type StartSamplingCommand struct {
	params *StartSamplingParams
	cb     StartSamplingCB
}

func NewStartSamplingCommand(params *StartSamplingParams, cb StartSamplingCB) *StartSamplingCommand {
	return &StartSamplingCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *StartSamplingCommand) Name() string {
	return "HeapProfiler.startSampling"
}

func (cmd *StartSamplingCommand) Params() interface{} {
	return cmd.params
}

func (cmd *StartSamplingCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type StopSamplingResult struct {
	Profile *SamplingHeapProfile `json:"profile"` // Recorded sampling heap profile.
}

type StopSamplingCB func(result *StopSamplingResult, err error)

type StopSamplingCommand struct {
	cb StopSamplingCB
}

func NewStopSamplingCommand(cb StopSamplingCB) *StopSamplingCommand {
	return &StopSamplingCommand{
		cb: cb,
	}
}

func (cmd *StopSamplingCommand) Name() string {
	return "HeapProfiler.stopSampling"
}

func (cmd *StopSamplingCommand) Params() interface{} {
	return nil
}

func (cmd *StopSamplingCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj StopSamplingResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type AddHeapSnapshotChunkEvent struct {
	Chunk string `json:"chunk"`
}

type AddHeapSnapshotChunkEventSink struct {
	events chan *AddHeapSnapshotChunkEvent
}

func NewAddHeapSnapshotChunkEventSink(bufSize int) *AddHeapSnapshotChunkEventSink {
	return &AddHeapSnapshotChunkEventSink{
		events: make(chan *AddHeapSnapshotChunkEvent, bufSize),
	}
}

func (s *AddHeapSnapshotChunkEventSink) Name() string {
	return "HeapProfiler.addHeapSnapshotChunk"
}

func (s *AddHeapSnapshotChunkEventSink) OnEvent(params []byte) {
	evt := &AddHeapSnapshotChunkEvent{}
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

type ResetProfilesEvent struct {
}

type ResetProfilesEventSink struct {
	events chan *ResetProfilesEvent
}

func NewResetProfilesEventSink(bufSize int) *ResetProfilesEventSink {
	return &ResetProfilesEventSink{
		events: make(chan *ResetProfilesEvent, bufSize),
	}
}

func (s *ResetProfilesEventSink) Name() string {
	return "HeapProfiler.resetProfiles"
}

func (s *ResetProfilesEventSink) OnEvent(params []byte) {
	evt := &ResetProfilesEvent{}
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

type ReportHeapSnapshotProgressEvent struct {
	Done     int  `json:"done"`
	Total    int  `json:"total"`
	Finished bool `json:"finished"`
}

type ReportHeapSnapshotProgressEventSink struct {
	events chan *ReportHeapSnapshotProgressEvent
}

func NewReportHeapSnapshotProgressEventSink(bufSize int) *ReportHeapSnapshotProgressEventSink {
	return &ReportHeapSnapshotProgressEventSink{
		events: make(chan *ReportHeapSnapshotProgressEvent, bufSize),
	}
}

func (s *ReportHeapSnapshotProgressEventSink) Name() string {
	return "HeapProfiler.reportHeapSnapshotProgress"
}

func (s *ReportHeapSnapshotProgressEventSink) OnEvent(params []byte) {
	evt := &ReportHeapSnapshotProgressEvent{}
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

type LastSeenObjectIdEvent struct {
	LastSeenObjectId int `json:"lastSeenObjectId"`
	Timestamp        int `json:"timestamp"`
}

// If heap objects tracking has been started then backend regulary sends a current value for last seen object id and corresponding timestamp. If the were changes in the heap since last event then one or more heapStatsUpdate events will be sent before a new lastSeenObjectId event.
type LastSeenObjectIdEventSink struct {
	events chan *LastSeenObjectIdEvent
}

func NewLastSeenObjectIdEventSink(bufSize int) *LastSeenObjectIdEventSink {
	return &LastSeenObjectIdEventSink{
		events: make(chan *LastSeenObjectIdEvent, bufSize),
	}
}

func (s *LastSeenObjectIdEventSink) Name() string {
	return "HeapProfiler.lastSeenObjectId"
}

func (s *LastSeenObjectIdEventSink) OnEvent(params []byte) {
	evt := &LastSeenObjectIdEvent{}
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

type HeapStatsUpdateEvent struct {
	StatsUpdate []int `json:"statsUpdate"` // An array of triplets. Each triplet describes a fragment. The first integer is the fragment index, the second integer is a total count of objects for the fragment, the third integer is a total size of the objects for the fragment.
}

// If heap objects tracking has been started then backend may send update for one or more fragments
type HeapStatsUpdateEventSink struct {
	events chan *HeapStatsUpdateEvent
}

func NewHeapStatsUpdateEventSink(bufSize int) *HeapStatsUpdateEventSink {
	return &HeapStatsUpdateEventSink{
		events: make(chan *HeapStatsUpdateEvent, bufSize),
	}
}

func (s *HeapStatsUpdateEventSink) Name() string {
	return "HeapProfiler.heapStatsUpdate"
}

func (s *HeapStatsUpdateEventSink) OnEvent(params []byte) {
	evt := &HeapStatsUpdateEvent{}
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
