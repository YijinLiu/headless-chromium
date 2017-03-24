package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

// Heap snapshot object id.
type HeapSnapshotObjectId string

// Sampling Heap Profile node. Holds callsite information, allocation statistics and child nodes.
type SamplingHeapProfileNode struct {
	CallFrame *RuntimeCallFrame          `json:"callFrame"` // Function location.
	SelfSize  float64                    `json:"selfSize"`  // Allocations size in bytes for the node excluding children.
	Children  []*SamplingHeapProfileNode `json:"children"`  // Child nodes.
}

// Profile.
type SamplingHeapProfile struct {
	Head *SamplingHeapProfileNode `json:"head"`
}

type HeapProfilerEnableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewHeapProfilerEnableCommand() *HeapProfilerEnableCommand {
	return &HeapProfilerEnableCommand{}
}

func (cmd *HeapProfilerEnableCommand) Name() string {
	return "HeapProfiler.enable"
}

func (cmd *HeapProfilerEnableCommand) Params() interface{} {
	return nil
}

func (cmd *HeapProfilerEnableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func HeapProfilerEnable(conn *hc.Conn) (err error) {
	cmd := NewHeapProfilerEnableCommand()
	cmd.Run(conn)
	return cmd.err
}

type HeapProfilerEnableCB func(err error)

type AsyncHeapProfilerEnableCommand struct {
	cb HeapProfilerEnableCB
}

func NewAsyncHeapProfilerEnableCommand(cb HeapProfilerEnableCB) *AsyncHeapProfilerEnableCommand {
	return &AsyncHeapProfilerEnableCommand{
		cb: cb,
	}
}

func (cmd *AsyncHeapProfilerEnableCommand) Name() string {
	return "HeapProfiler.enable"
}

func (cmd *AsyncHeapProfilerEnableCommand) Params() interface{} {
	return nil
}

func (cmd *HeapProfilerEnableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncHeapProfilerEnableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type HeapProfilerDisableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewHeapProfilerDisableCommand() *HeapProfilerDisableCommand {
	return &HeapProfilerDisableCommand{}
}

func (cmd *HeapProfilerDisableCommand) Name() string {
	return "HeapProfiler.disable"
}

func (cmd *HeapProfilerDisableCommand) Params() interface{} {
	return nil
}

func (cmd *HeapProfilerDisableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func HeapProfilerDisable(conn *hc.Conn) (err error) {
	cmd := NewHeapProfilerDisableCommand()
	cmd.Run(conn)
	return cmd.err
}

type HeapProfilerDisableCB func(err error)

type AsyncHeapProfilerDisableCommand struct {
	cb HeapProfilerDisableCB
}

func NewAsyncHeapProfilerDisableCommand(cb HeapProfilerDisableCB) *AsyncHeapProfilerDisableCommand {
	return &AsyncHeapProfilerDisableCommand{
		cb: cb,
	}
}

func (cmd *AsyncHeapProfilerDisableCommand) Name() string {
	return "HeapProfiler.disable"
}

func (cmd *AsyncHeapProfilerDisableCommand) Params() interface{} {
	return nil
}

func (cmd *HeapProfilerDisableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncHeapProfilerDisableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type StartTrackingHeapObjectsParams struct {
	TrackAllocations bool `json:"trackAllocations,omitempty"`
}

type StartTrackingHeapObjectsCommand struct {
	params *StartTrackingHeapObjectsParams
	wg     sync.WaitGroup
	err    error
}

func NewStartTrackingHeapObjectsCommand(params *StartTrackingHeapObjectsParams) *StartTrackingHeapObjectsCommand {
	return &StartTrackingHeapObjectsCommand{
		params: params,
	}
}

func (cmd *StartTrackingHeapObjectsCommand) Name() string {
	return "HeapProfiler.startTrackingHeapObjects"
}

func (cmd *StartTrackingHeapObjectsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *StartTrackingHeapObjectsCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func StartTrackingHeapObjects(params *StartTrackingHeapObjectsParams, conn *hc.Conn) (err error) {
	cmd := NewStartTrackingHeapObjectsCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type StartTrackingHeapObjectsCB func(err error)

type AsyncStartTrackingHeapObjectsCommand struct {
	params *StartTrackingHeapObjectsParams
	cb     StartTrackingHeapObjectsCB
}

func NewAsyncStartTrackingHeapObjectsCommand(params *StartTrackingHeapObjectsParams, cb StartTrackingHeapObjectsCB) *AsyncStartTrackingHeapObjectsCommand {
	return &AsyncStartTrackingHeapObjectsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncStartTrackingHeapObjectsCommand) Name() string {
	return "HeapProfiler.startTrackingHeapObjects"
}

func (cmd *AsyncStartTrackingHeapObjectsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *StartTrackingHeapObjectsCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncStartTrackingHeapObjectsCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type StopTrackingHeapObjectsParams struct {
	ReportProgress bool `json:"reportProgress,omitempty"` // If true 'reportHeapSnapshotProgress' events will be generated while snapshot is being taken when the tracking is stopped.
}

type StopTrackingHeapObjectsCommand struct {
	params *StopTrackingHeapObjectsParams
	wg     sync.WaitGroup
	err    error
}

func NewStopTrackingHeapObjectsCommand(params *StopTrackingHeapObjectsParams) *StopTrackingHeapObjectsCommand {
	return &StopTrackingHeapObjectsCommand{
		params: params,
	}
}

func (cmd *StopTrackingHeapObjectsCommand) Name() string {
	return "HeapProfiler.stopTrackingHeapObjects"
}

func (cmd *StopTrackingHeapObjectsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *StopTrackingHeapObjectsCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func StopTrackingHeapObjects(params *StopTrackingHeapObjectsParams, conn *hc.Conn) (err error) {
	cmd := NewStopTrackingHeapObjectsCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type StopTrackingHeapObjectsCB func(err error)

type AsyncStopTrackingHeapObjectsCommand struct {
	params *StopTrackingHeapObjectsParams
	cb     StopTrackingHeapObjectsCB
}

func NewAsyncStopTrackingHeapObjectsCommand(params *StopTrackingHeapObjectsParams, cb StopTrackingHeapObjectsCB) *AsyncStopTrackingHeapObjectsCommand {
	return &AsyncStopTrackingHeapObjectsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncStopTrackingHeapObjectsCommand) Name() string {
	return "HeapProfiler.stopTrackingHeapObjects"
}

func (cmd *AsyncStopTrackingHeapObjectsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *StopTrackingHeapObjectsCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncStopTrackingHeapObjectsCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type TakeHeapSnapshotParams struct {
	ReportProgress bool `json:"reportProgress,omitempty"` // If true 'reportHeapSnapshotProgress' events will be generated while snapshot is being taken.
}

type TakeHeapSnapshotCommand struct {
	params *TakeHeapSnapshotParams
	wg     sync.WaitGroup
	err    error
}

func NewTakeHeapSnapshotCommand(params *TakeHeapSnapshotParams) *TakeHeapSnapshotCommand {
	return &TakeHeapSnapshotCommand{
		params: params,
	}
}

func (cmd *TakeHeapSnapshotCommand) Name() string {
	return "HeapProfiler.takeHeapSnapshot"
}

func (cmd *TakeHeapSnapshotCommand) Params() interface{} {
	return cmd.params
}

func (cmd *TakeHeapSnapshotCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func TakeHeapSnapshot(params *TakeHeapSnapshotParams, conn *hc.Conn) (err error) {
	cmd := NewTakeHeapSnapshotCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type TakeHeapSnapshotCB func(err error)

type AsyncTakeHeapSnapshotCommand struct {
	params *TakeHeapSnapshotParams
	cb     TakeHeapSnapshotCB
}

func NewAsyncTakeHeapSnapshotCommand(params *TakeHeapSnapshotParams, cb TakeHeapSnapshotCB) *AsyncTakeHeapSnapshotCommand {
	return &AsyncTakeHeapSnapshotCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncTakeHeapSnapshotCommand) Name() string {
	return "HeapProfiler.takeHeapSnapshot"
}

func (cmd *AsyncTakeHeapSnapshotCommand) Params() interface{} {
	return cmd.params
}

func (cmd *TakeHeapSnapshotCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncTakeHeapSnapshotCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type CollectGarbageCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewCollectGarbageCommand() *CollectGarbageCommand {
	return &CollectGarbageCommand{}
}

func (cmd *CollectGarbageCommand) Name() string {
	return "HeapProfiler.collectGarbage"
}

func (cmd *CollectGarbageCommand) Params() interface{} {
	return nil
}

func (cmd *CollectGarbageCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func CollectGarbage(conn *hc.Conn) (err error) {
	cmd := NewCollectGarbageCommand()
	cmd.Run(conn)
	return cmd.err
}

type CollectGarbageCB func(err error)

type AsyncCollectGarbageCommand struct {
	cb CollectGarbageCB
}

func NewAsyncCollectGarbageCommand(cb CollectGarbageCB) *AsyncCollectGarbageCommand {
	return &AsyncCollectGarbageCommand{
		cb: cb,
	}
}

func (cmd *AsyncCollectGarbageCommand) Name() string {
	return "HeapProfiler.collectGarbage"
}

func (cmd *AsyncCollectGarbageCommand) Params() interface{} {
	return nil
}

func (cmd *CollectGarbageCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncCollectGarbageCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type GetObjectByHeapObjectIdParams struct {
	ObjectId    HeapSnapshotObjectId `json:"objectId"`
	ObjectGroup string               `json:"objectGroup,omitempty"` // Symbolic group name that can be used to release multiple objects.
}

type GetObjectByHeapObjectIdResult struct {
	Result *RemoteObject `json:"result"` // Evaluation result.
}

type GetObjectByHeapObjectIdCommand struct {
	params *GetObjectByHeapObjectIdParams
	result GetObjectByHeapObjectIdResult
	wg     sync.WaitGroup
	err    error
}

func NewGetObjectByHeapObjectIdCommand(params *GetObjectByHeapObjectIdParams) *GetObjectByHeapObjectIdCommand {
	return &GetObjectByHeapObjectIdCommand{
		params: params,
	}
}

func (cmd *GetObjectByHeapObjectIdCommand) Name() string {
	return "HeapProfiler.getObjectByHeapObjectId"
}

func (cmd *GetObjectByHeapObjectIdCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetObjectByHeapObjectIdCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetObjectByHeapObjectId(params *GetObjectByHeapObjectIdParams, conn *hc.Conn) (result *GetObjectByHeapObjectIdResult, err error) {
	cmd := NewGetObjectByHeapObjectIdCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetObjectByHeapObjectIdCB func(result *GetObjectByHeapObjectIdResult, err error)

type AsyncGetObjectByHeapObjectIdCommand struct {
	params *GetObjectByHeapObjectIdParams
	cb     GetObjectByHeapObjectIdCB
}

func NewAsyncGetObjectByHeapObjectIdCommand(params *GetObjectByHeapObjectIdParams, cb GetObjectByHeapObjectIdCB) *AsyncGetObjectByHeapObjectIdCommand {
	return &AsyncGetObjectByHeapObjectIdCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetObjectByHeapObjectIdCommand) Name() string {
	return "HeapProfiler.getObjectByHeapObjectId"
}

func (cmd *AsyncGetObjectByHeapObjectIdCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetObjectByHeapObjectIdCommand) Result() *GetObjectByHeapObjectIdResult {
	return &cmd.result
}

func (cmd *GetObjectByHeapObjectIdCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetObjectByHeapObjectIdCommand) Done(data []byte, err error) {
	var result GetObjectByHeapObjectIdResult
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

type AddInspectedHeapObjectParams struct {
	HeapObjectId HeapSnapshotObjectId `json:"heapObjectId"` // Heap snapshot object id to be accessible by means of $x command line API.
}

// Enables console to refer to the node with given id via $x (see Command Line API for more details $x functions).

type AddInspectedHeapObjectCommand struct {
	params *AddInspectedHeapObjectParams
	wg     sync.WaitGroup
	err    error
}

func NewAddInspectedHeapObjectCommand(params *AddInspectedHeapObjectParams) *AddInspectedHeapObjectCommand {
	return &AddInspectedHeapObjectCommand{
		params: params,
	}
}

func (cmd *AddInspectedHeapObjectCommand) Name() string {
	return "HeapProfiler.addInspectedHeapObject"
}

func (cmd *AddInspectedHeapObjectCommand) Params() interface{} {
	return cmd.params
}

func (cmd *AddInspectedHeapObjectCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func AddInspectedHeapObject(params *AddInspectedHeapObjectParams, conn *hc.Conn) (err error) {
	cmd := NewAddInspectedHeapObjectCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type AddInspectedHeapObjectCB func(err error)

// Enables console to refer to the node with given id via $x (see Command Line API for more details $x functions).

type AsyncAddInspectedHeapObjectCommand struct {
	params *AddInspectedHeapObjectParams
	cb     AddInspectedHeapObjectCB
}

func NewAsyncAddInspectedHeapObjectCommand(params *AddInspectedHeapObjectParams, cb AddInspectedHeapObjectCB) *AsyncAddInspectedHeapObjectCommand {
	return &AsyncAddInspectedHeapObjectCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncAddInspectedHeapObjectCommand) Name() string {
	return "HeapProfiler.addInspectedHeapObject"
}

func (cmd *AsyncAddInspectedHeapObjectCommand) Params() interface{} {
	return cmd.params
}

func (cmd *AddInspectedHeapObjectCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncAddInspectedHeapObjectCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type GetHeapObjectIdParams struct {
	ObjectId *RemoteObjectId `json:"objectId"` // Identifier of the object to get heap object id for.
}

type GetHeapObjectIdResult struct {
	HeapSnapshotObjectId HeapSnapshotObjectId `json:"heapSnapshotObjectId"` // Id of the heap snapshot object corresponding to the passed remote object id.
}

type GetHeapObjectIdCommand struct {
	params *GetHeapObjectIdParams
	result GetHeapObjectIdResult
	wg     sync.WaitGroup
	err    error
}

func NewGetHeapObjectIdCommand(params *GetHeapObjectIdParams) *GetHeapObjectIdCommand {
	return &GetHeapObjectIdCommand{
		params: params,
	}
}

func (cmd *GetHeapObjectIdCommand) Name() string {
	return "HeapProfiler.getHeapObjectId"
}

func (cmd *GetHeapObjectIdCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetHeapObjectIdCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetHeapObjectId(params *GetHeapObjectIdParams, conn *hc.Conn) (result *GetHeapObjectIdResult, err error) {
	cmd := NewGetHeapObjectIdCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetHeapObjectIdCB func(result *GetHeapObjectIdResult, err error)

type AsyncGetHeapObjectIdCommand struct {
	params *GetHeapObjectIdParams
	cb     GetHeapObjectIdCB
}

func NewAsyncGetHeapObjectIdCommand(params *GetHeapObjectIdParams, cb GetHeapObjectIdCB) *AsyncGetHeapObjectIdCommand {
	return &AsyncGetHeapObjectIdCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetHeapObjectIdCommand) Name() string {
	return "HeapProfiler.getHeapObjectId"
}

func (cmd *AsyncGetHeapObjectIdCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetHeapObjectIdCommand) Result() *GetHeapObjectIdResult {
	return &cmd.result
}

func (cmd *GetHeapObjectIdCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetHeapObjectIdCommand) Done(data []byte, err error) {
	var result GetHeapObjectIdResult
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

type StartSamplingParams struct {
	SamplingInterval float64 `json:"samplingInterval,omitempty"` // Average sample interval in bytes. Poisson distribution is used for the intervals. The default value is 32768 bytes.
}

type StartSamplingCommand struct {
	params *StartSamplingParams
	wg     sync.WaitGroup
	err    error
}

func NewStartSamplingCommand(params *StartSamplingParams) *StartSamplingCommand {
	return &StartSamplingCommand{
		params: params,
	}
}

func (cmd *StartSamplingCommand) Name() string {
	return "HeapProfiler.startSampling"
}

func (cmd *StartSamplingCommand) Params() interface{} {
	return cmd.params
}

func (cmd *StartSamplingCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func StartSampling(params *StartSamplingParams, conn *hc.Conn) (err error) {
	cmd := NewStartSamplingCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type StartSamplingCB func(err error)

type AsyncStartSamplingCommand struct {
	params *StartSamplingParams
	cb     StartSamplingCB
}

func NewAsyncStartSamplingCommand(params *StartSamplingParams, cb StartSamplingCB) *AsyncStartSamplingCommand {
	return &AsyncStartSamplingCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncStartSamplingCommand) Name() string {
	return "HeapProfiler.startSampling"
}

func (cmd *AsyncStartSamplingCommand) Params() interface{} {
	return cmd.params
}

func (cmd *StartSamplingCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncStartSamplingCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type StopSamplingResult struct {
	Profile *SamplingHeapProfile `json:"profile"` // Recorded sampling heap profile.
}

type StopSamplingCommand struct {
	result StopSamplingResult
	wg     sync.WaitGroup
	err    error
}

func NewStopSamplingCommand() *StopSamplingCommand {
	return &StopSamplingCommand{}
}

func (cmd *StopSamplingCommand) Name() string {
	return "HeapProfiler.stopSampling"
}

func (cmd *StopSamplingCommand) Params() interface{} {
	return nil
}

func (cmd *StopSamplingCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func StopSampling(conn *hc.Conn) (result *StopSamplingResult, err error) {
	cmd := NewStopSamplingCommand()
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type StopSamplingCB func(result *StopSamplingResult, err error)

type AsyncStopSamplingCommand struct {
	cb StopSamplingCB
}

func NewAsyncStopSamplingCommand(cb StopSamplingCB) *AsyncStopSamplingCommand {
	return &AsyncStopSamplingCommand{
		cb: cb,
	}
}

func (cmd *AsyncStopSamplingCommand) Name() string {
	return "HeapProfiler.stopSampling"
}

func (cmd *AsyncStopSamplingCommand) Params() interface{} {
	return nil
}

func (cmd *StopSamplingCommand) Result() *StopSamplingResult {
	return &cmd.result
}

func (cmd *StopSamplingCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncStopSamplingCommand) Done(data []byte, err error) {
	var result StopSamplingResult
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

type AddHeapSnapshotChunkEvent struct {
	Chunk string `json:"chunk"`
}

func OnAddHeapSnapshotChunk(conn *hc.Conn, cb func(evt *AddHeapSnapshotChunkEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &AddHeapSnapshotChunkEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("HeapProfiler.addHeapSnapshotChunk", sink)
}

type ResetProfilesEvent struct {
}

func OnResetProfiles(conn *hc.Conn, cb func(evt *ResetProfilesEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &ResetProfilesEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("HeapProfiler.resetProfiles", sink)
}

type ReportHeapSnapshotProgressEvent struct {
	Done     int  `json:"done"`
	Total    int  `json:"total"`
	Finished bool `json:"finished"`
}

func OnReportHeapSnapshotProgress(conn *hc.Conn, cb func(evt *ReportHeapSnapshotProgressEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &ReportHeapSnapshotProgressEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("HeapProfiler.reportHeapSnapshotProgress", sink)
}

// If heap objects tracking has been started then backend regulary sends a current value for last seen object id and corresponding timestamp. If the were changes in the heap since last event then one or more heapStatsUpdate events will be sent before a new lastSeenObjectId event.

type LastSeenObjectIdEvent struct {
	LastSeenObjectId int     `json:"lastSeenObjectId"`
	Timestamp        float64 `json:"timestamp"`
}

func OnLastSeenObjectId(conn *hc.Conn, cb func(evt *LastSeenObjectIdEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &LastSeenObjectIdEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("HeapProfiler.lastSeenObjectId", sink)
}

// If heap objects tracking has been started then backend may send update for one or more fragments

type HeapStatsUpdateEvent struct {
	StatsUpdate []int `json:"statsUpdate"` // An array of triplets. Each triplet describes a fragment. The first integer is the fragment index, the second integer is a total count of objects for the fragment, the third integer is a total size of the objects for the fragment.
}

func OnHeapStatsUpdate(conn *hc.Conn, cb func(evt *HeapStatsUpdateEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &HeapStatsUpdateEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("HeapProfiler.heapStatsUpdate", sink)
}
