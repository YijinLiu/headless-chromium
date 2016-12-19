package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
)

// Configuration for memory dump. Used only when "memory-infra" category is enabled.
type MemoryDumpConfig struct {
}

type TraceConfig struct {
	RecordMode           string            `json:"recordMode"`           // Controls how the trace buffer stores data.
	EnableSampling       bool              `json:"enableSampling"`       // Turns on JavaScript stack sampling.
	EnableSystrace       bool              `json:"enableSystrace"`       // Turns on system tracing.
	EnableArgumentFilter bool              `json:"enableArgumentFilter"` // Turns on argument filter.
	IncludedCategories   []string          `json:"includedCategories"`   // Included category filters.
	ExcludedCategories   []string          `json:"excludedCategories"`   // Excluded category filters.
	SyntheticDelays      []string          `json:"syntheticDelays"`      // Configuration to synthesize the delays in tracing.
	MemoryDumpConfig     *MemoryDumpConfig `json:"memoryDumpConfig"`     // Configuration for memory dump triggers. Used only when "memory-infra" category is enabled.
}

type TracingStartParams struct {
	Categories                   string       `json:"categories"`                   // Category/tag filter
	Options                      string       `json:"options"`                      // Tracing options
	BufferUsageReportingInterval int          `json:"bufferUsageReportingInterval"` // If set, the agent will issue bufferUsage events at this interval, specified in milliseconds
	TransferMode                 string       `json:"transferMode"`                 // Whether to report trace events as series of dataCollected events or to save trace to a stream (defaults to ReportEvents).
	TraceConfig                  *TraceConfig `json:"traceConfig"`
}

type TracingStartCB func(err error)

// Start trace events collection.
type TracingStartCommand struct {
	params *TracingStartParams
	cb     TracingStartCB
}

func NewTracingStartCommand(params *TracingStartParams, cb TracingStartCB) *TracingStartCommand {
	return &TracingStartCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *TracingStartCommand) Name() string {
	return "Tracing.start"
}

func (cmd *TracingStartCommand) Params() interface{} {
	return cmd.params
}

func (cmd *TracingStartCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type EndCB func(err error)

// Stop trace events collection.
type EndCommand struct {
	cb EndCB
}

func NewEndCommand(cb EndCB) *EndCommand {
	return &EndCommand{
		cb: cb,
	}
}

func (cmd *EndCommand) Name() string {
	return "Tracing.end"
}

func (cmd *EndCommand) Params() interface{} {
	return nil
}

func (cmd *EndCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type GetCategoriesResult struct {
	Categories []string `json:"categories"` // A list of supported tracing categories.
}

type GetCategoriesCB func(result *GetCategoriesResult, err error)

// Gets supported tracing categories.
type GetCategoriesCommand struct {
	cb GetCategoriesCB
}

func NewGetCategoriesCommand(cb GetCategoriesCB) *GetCategoriesCommand {
	return &GetCategoriesCommand{
		cb: cb,
	}
}

func (cmd *GetCategoriesCommand) Name() string {
	return "Tracing.getCategories"
}

func (cmd *GetCategoriesCommand) Params() interface{} {
	return nil
}

func (cmd *GetCategoriesCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetCategoriesResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type RequestMemoryDumpResult struct {
	DumpGuid string `json:"dumpGuid"` // GUID of the resulting global memory dump.
	Success  bool   `json:"success"`  // True iff the global memory dump succeeded.
}

type RequestMemoryDumpCB func(result *RequestMemoryDumpResult, err error)

// Request a global memory dump.
type RequestMemoryDumpCommand struct {
	cb RequestMemoryDumpCB
}

func NewRequestMemoryDumpCommand(cb RequestMemoryDumpCB) *RequestMemoryDumpCommand {
	return &RequestMemoryDumpCommand{
		cb: cb,
	}
}

func (cmd *RequestMemoryDumpCommand) Name() string {
	return "Tracing.requestMemoryDump"
}

func (cmd *RequestMemoryDumpCommand) Params() interface{} {
	return nil
}

func (cmd *RequestMemoryDumpCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj RequestMemoryDumpResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type RecordClockSyncMarkerParams struct {
	SyncId string `json:"syncId"` // The ID of this clock sync marker
}

type RecordClockSyncMarkerCB func(err error)

// Record a clock sync marker in the trace.
type RecordClockSyncMarkerCommand struct {
	params *RecordClockSyncMarkerParams
	cb     RecordClockSyncMarkerCB
}

func NewRecordClockSyncMarkerCommand(params *RecordClockSyncMarkerParams, cb RecordClockSyncMarkerCB) *RecordClockSyncMarkerCommand {
	return &RecordClockSyncMarkerCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *RecordClockSyncMarkerCommand) Name() string {
	return "Tracing.recordClockSyncMarker"
}

func (cmd *RecordClockSyncMarkerCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RecordClockSyncMarkerCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type DataCollectedEvent struct {
	Value []map[string]string `json:"value"`
}

// Contains an bucket of collected trace events. When tracing is stopped collected events will be send as a sequence of dataCollected events followed by tracingComplete event.
type DataCollectedEventSink struct {
	events chan *DataCollectedEvent
}

func NewDataCollectedEventSink(bufSize int) *DataCollectedEventSink {
	return &DataCollectedEventSink{
		events: make(chan *DataCollectedEvent, bufSize),
	}
}

func (s *DataCollectedEventSink) Name() string {
	return "Tracing.dataCollected"
}

func (s *DataCollectedEventSink) OnEvent(params []byte) {
	evt := &DataCollectedEvent{}
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

type TracingCompleteEvent struct {
	Stream *StreamHandle `json:"stream"` // A handle of the stream that holds resulting trace data.
}

// Signals that tracing is stopped and there is no trace buffers pending flush, all data were delivered via dataCollected events.
type TracingCompleteEventSink struct {
	events chan *TracingCompleteEvent
}

func NewTracingCompleteEventSink(bufSize int) *TracingCompleteEventSink {
	return &TracingCompleteEventSink{
		events: make(chan *TracingCompleteEvent, bufSize),
	}
}

func (s *TracingCompleteEventSink) Name() string {
	return "Tracing.tracingComplete"
}

func (s *TracingCompleteEventSink) OnEvent(params []byte) {
	evt := &TracingCompleteEvent{}
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

type BufferUsageEvent struct {
	PercentFull int `json:"percentFull"` // A number in range [0..1] that indicates the used size of event buffer as a fraction of its total size.
	EventCount  int `json:"eventCount"`  // An approximate number of events in the trace log.
	Value       int `json:"value"`       // A number in range [0..1] that indicates the used size of event buffer as a fraction of its total size.
}

type BufferUsageEventSink struct {
	events chan *BufferUsageEvent
}

func NewBufferUsageEventSink(bufSize int) *BufferUsageEventSink {
	return &BufferUsageEventSink{
		events: make(chan *BufferUsageEvent, bufSize),
	}
}

func (s *BufferUsageEventSink) Name() string {
	return "Tracing.bufferUsage"
}

func (s *BufferUsageEventSink) OnEvent(params []byte) {
	evt := &BufferUsageEvent{}
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
