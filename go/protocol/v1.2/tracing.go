package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

// Configuration for memory dump. Used only when "memory-infra" category is enabled.
type MemoryDumpConfig struct {
}

type TraceConfig struct {
	RecordMode           string            `json:"recordMode,omitempty"`           // Controls how the trace buffer stores data.
	EnableSampling       bool              `json:"enableSampling,omitempty"`       // Turns on JavaScript stack sampling.
	EnableSystrace       bool              `json:"enableSystrace,omitempty"`       // Turns on system tracing.
	EnableArgumentFilter bool              `json:"enableArgumentFilter,omitempty"` // Turns on argument filter.
	IncludedCategories   []string          `json:"includedCategories,omitempty"`   // Included category filters.
	ExcludedCategories   []string          `json:"excludedCategories,omitempty"`   // Excluded category filters.
	SyntheticDelays      []string          `json:"syntheticDelays,omitempty"`      // Configuration to synthesize the delays in tracing.
	MemoryDumpConfig     *MemoryDumpConfig `json:"memoryDumpConfig,omitempty"`     // Configuration for memory dump triggers. Used only when "memory-infra" category is enabled.
}

type TracingStartParams struct {
	Categories                   string       `json:"categories,omitempty"`                   // Category/tag filter
	Options                      string       `json:"options,omitempty"`                      // Tracing options
	BufferUsageReportingInterval float64      `json:"bufferUsageReportingInterval,omitempty"` // If set, the agent will issue bufferUsage events at this interval, specified in milliseconds
	TransferMode                 string       `json:"transferMode,omitempty"`                 // Whether to report trace events as series of dataCollected events or to save trace to a stream (defaults to ReportEvents).
	TraceConfig                  *TraceConfig `json:"traceConfig,omitempty"`
}

// Start trace events collection.

type TracingStartCommand struct {
	params *TracingStartParams
	wg     sync.WaitGroup
	err    error
}

func NewTracingStartCommand(params *TracingStartParams) *TracingStartCommand {
	return &TracingStartCommand{
		params: params,
	}
}

func (cmd *TracingStartCommand) Name() string {
	return "Tracing.start"
}

func (cmd *TracingStartCommand) Params() interface{} {
	return cmd.params
}

func (cmd *TracingStartCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func TracingStart(params *TracingStartParams, conn *hc.Conn) (err error) {
	cmd := NewTracingStartCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type TracingStartCB func(err error)

// Start trace events collection.

type AsyncTracingStartCommand struct {
	params *TracingStartParams
	cb     TracingStartCB
}

func NewAsyncTracingStartCommand(params *TracingStartParams, cb TracingStartCB) *AsyncTracingStartCommand {
	return &AsyncTracingStartCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncTracingStartCommand) Name() string {
	return "Tracing.start"
}

func (cmd *AsyncTracingStartCommand) Params() interface{} {
	return cmd.params
}

func (cmd *TracingStartCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncTracingStartCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Stop trace events collection.

type EndCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewEndCommand() *EndCommand {
	return &EndCommand{}
}

func (cmd *EndCommand) Name() string {
	return "Tracing.end"
}

func (cmd *EndCommand) Params() interface{} {
	return nil
}

func (cmd *EndCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func End(conn *hc.Conn) (err error) {
	cmd := NewEndCommand()
	cmd.Run(conn)
	return cmd.err
}

type EndCB func(err error)

// Stop trace events collection.

type AsyncEndCommand struct {
	cb EndCB
}

func NewAsyncEndCommand(cb EndCB) *AsyncEndCommand {
	return &AsyncEndCommand{
		cb: cb,
	}
}

func (cmd *AsyncEndCommand) Name() string {
	return "Tracing.end"
}

func (cmd *AsyncEndCommand) Params() interface{} {
	return nil
}

func (cmd *EndCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncEndCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type GetCategoriesResult struct {
	Categories []string `json:"categories"` // A list of supported tracing categories.
}

// Gets supported tracing categories.

type GetCategoriesCommand struct {
	result GetCategoriesResult
	wg     sync.WaitGroup
	err    error
}

func NewGetCategoriesCommand() *GetCategoriesCommand {
	return &GetCategoriesCommand{}
}

func (cmd *GetCategoriesCommand) Name() string {
	return "Tracing.getCategories"
}

func (cmd *GetCategoriesCommand) Params() interface{} {
	return nil
}

func (cmd *GetCategoriesCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetCategories(conn *hc.Conn) (result *GetCategoriesResult, err error) {
	cmd := NewGetCategoriesCommand()
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetCategoriesCB func(result *GetCategoriesResult, err error)

// Gets supported tracing categories.

type AsyncGetCategoriesCommand struct {
	cb GetCategoriesCB
}

func NewAsyncGetCategoriesCommand(cb GetCategoriesCB) *AsyncGetCategoriesCommand {
	return &AsyncGetCategoriesCommand{
		cb: cb,
	}
}

func (cmd *AsyncGetCategoriesCommand) Name() string {
	return "Tracing.getCategories"
}

func (cmd *AsyncGetCategoriesCommand) Params() interface{} {
	return nil
}

func (cmd *GetCategoriesCommand) Result() *GetCategoriesResult {
	return &cmd.result
}

func (cmd *GetCategoriesCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetCategoriesCommand) Done(data []byte, err error) {
	var result GetCategoriesResult
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

type RequestMemoryDumpResult struct {
	DumpGuid string `json:"dumpGuid"` // GUID of the resulting global memory dump.
	Success  bool   `json:"success"`  // True iff the global memory dump succeeded.
}

// Request a global memory dump.

type RequestMemoryDumpCommand struct {
	result RequestMemoryDumpResult
	wg     sync.WaitGroup
	err    error
}

func NewRequestMemoryDumpCommand() *RequestMemoryDumpCommand {
	return &RequestMemoryDumpCommand{}
}

func (cmd *RequestMemoryDumpCommand) Name() string {
	return "Tracing.requestMemoryDump"
}

func (cmd *RequestMemoryDumpCommand) Params() interface{} {
	return nil
}

func (cmd *RequestMemoryDumpCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func RequestMemoryDump(conn *hc.Conn) (result *RequestMemoryDumpResult, err error) {
	cmd := NewRequestMemoryDumpCommand()
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type RequestMemoryDumpCB func(result *RequestMemoryDumpResult, err error)

// Request a global memory dump.

type AsyncRequestMemoryDumpCommand struct {
	cb RequestMemoryDumpCB
}

func NewAsyncRequestMemoryDumpCommand(cb RequestMemoryDumpCB) *AsyncRequestMemoryDumpCommand {
	return &AsyncRequestMemoryDumpCommand{
		cb: cb,
	}
}

func (cmd *AsyncRequestMemoryDumpCommand) Name() string {
	return "Tracing.requestMemoryDump"
}

func (cmd *AsyncRequestMemoryDumpCommand) Params() interface{} {
	return nil
}

func (cmd *RequestMemoryDumpCommand) Result() *RequestMemoryDumpResult {
	return &cmd.result
}

func (cmd *RequestMemoryDumpCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncRequestMemoryDumpCommand) Done(data []byte, err error) {
	var result RequestMemoryDumpResult
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

type RecordClockSyncMarkerParams struct {
	SyncId string `json:"syncId"` // The ID of this clock sync marker
}

// Record a clock sync marker in the trace.

type RecordClockSyncMarkerCommand struct {
	params *RecordClockSyncMarkerParams
	wg     sync.WaitGroup
	err    error
}

func NewRecordClockSyncMarkerCommand(params *RecordClockSyncMarkerParams) *RecordClockSyncMarkerCommand {
	return &RecordClockSyncMarkerCommand{
		params: params,
	}
}

func (cmd *RecordClockSyncMarkerCommand) Name() string {
	return "Tracing.recordClockSyncMarker"
}

func (cmd *RecordClockSyncMarkerCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RecordClockSyncMarkerCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func RecordClockSyncMarker(params *RecordClockSyncMarkerParams, conn *hc.Conn) (err error) {
	cmd := NewRecordClockSyncMarkerCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type RecordClockSyncMarkerCB func(err error)

// Record a clock sync marker in the trace.

type AsyncRecordClockSyncMarkerCommand struct {
	params *RecordClockSyncMarkerParams
	cb     RecordClockSyncMarkerCB
}

func NewAsyncRecordClockSyncMarkerCommand(params *RecordClockSyncMarkerParams, cb RecordClockSyncMarkerCB) *AsyncRecordClockSyncMarkerCommand {
	return &AsyncRecordClockSyncMarkerCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncRecordClockSyncMarkerCommand) Name() string {
	return "Tracing.recordClockSyncMarker"
}

func (cmd *AsyncRecordClockSyncMarkerCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RecordClockSyncMarkerCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncRecordClockSyncMarkerCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Contains an bucket of collected trace events. When tracing is stopped collected events will be send as a sequence of dataCollected events followed by tracingComplete event.

type DataCollectedEvent struct {
	Value []map[string]string `json:"value"`
}

func OnDataCollected(conn *hc.Conn, cb func(evt *DataCollectedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &DataCollectedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Tracing.dataCollected", sink)
}

// Signals that tracing is stopped and there is no trace buffers pending flush, all data were delivered via dataCollected events.

type TracingCompleteEvent struct {
	Stream *StreamHandle `json:"stream"` // A handle of the stream that holds resulting trace data.
}

func OnTracingComplete(conn *hc.Conn, cb func(evt *TracingCompleteEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &TracingCompleteEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Tracing.tracingComplete", sink)
}

type BufferUsageEvent struct {
	PercentFull float64 `json:"percentFull"` // A number in range [0..1] that indicates the used size of event buffer as a fraction of its total size.
	EventCount  float64 `json:"eventCount"`  // An approximate number of events in the trace log.
	Value       float64 `json:"value"`       // A number in range [0..1] that indicates the used size of event buffer as a fraction of its total size.
}

func OnBufferUsage(conn *hc.Conn, cb func(evt *BufferUsageEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &BufferUsageEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Tracing.bufferUsage", sink)
}
