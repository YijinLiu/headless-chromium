package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
)

// Profile node. Holds callsite information, execution statistics and child nodes.
type ProfileNode struct {
	Id            int                 `json:"id"`            // Unique id of the node.
	CallFrame     *RuntimeCallFrame   `json:"callFrame"`     // Function location.
	HitCount      int                 `json:"hitCount"`      // Number of samples where this node was on top of the call stack.
	Children      []int               `json:"children"`      // Child node ids.
	DeoptReason   string              `json:"deoptReason"`   // The reason of being not optimized. The function may be deoptimized or marked as don't optimize.
	PositionTicks []*PositionTickInfo `json:"positionTicks"` // An array of source position ticks.
}

// Profile.
type Profile struct {
	Nodes      []*ProfileNode `json:"nodes"`      // The list of profile nodes. First item is the root node.
	StartTime  int            `json:"startTime"`  // Profiling start timestamp in microseconds.
	EndTime    int            `json:"endTime"`    // Profiling end timestamp in microseconds.
	Samples    []int          `json:"samples"`    // Ids of samples top nodes.
	TimeDeltas []int          `json:"timeDeltas"` // Time intervals between adjacent samples in microseconds. The first delta is relative to the profile startTime.
}

// Specifies a number of samples attributed to a certain source position.
type PositionTickInfo struct {
	Line  int `json:"line"`  // Source line number (1-based).
	Ticks int `json:"ticks"` // Number of samples attributed to the source line.
}

type ProfilerEnableCB func(err error)

type ProfilerEnableCommand struct {
	cb ProfilerEnableCB
}

func NewProfilerEnableCommand(cb ProfilerEnableCB) *ProfilerEnableCommand {
	return &ProfilerEnableCommand{
		cb: cb,
	}
}

func (cmd *ProfilerEnableCommand) Name() string {
	return "Profiler.enable"
}

func (cmd *ProfilerEnableCommand) Params() interface{} {
	return nil
}

func (cmd *ProfilerEnableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type ProfilerDisableCB func(err error)

type ProfilerDisableCommand struct {
	cb ProfilerDisableCB
}

func NewProfilerDisableCommand(cb ProfilerDisableCB) *ProfilerDisableCommand {
	return &ProfilerDisableCommand{
		cb: cb,
	}
}

func (cmd *ProfilerDisableCommand) Name() string {
	return "Profiler.disable"
}

func (cmd *ProfilerDisableCommand) Params() interface{} {
	return nil
}

func (cmd *ProfilerDisableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetSamplingIntervalParams struct {
	Interval int `json:"interval"` // New sampling interval in microseconds.
}

type SetSamplingIntervalCB func(err error)

// Changes CPU profiler sampling interval. Must be called before CPU profiles recording started.
type SetSamplingIntervalCommand struct {
	params *SetSamplingIntervalParams
	cb     SetSamplingIntervalCB
}

func NewSetSamplingIntervalCommand(params *SetSamplingIntervalParams, cb SetSamplingIntervalCB) *SetSamplingIntervalCommand {
	return &SetSamplingIntervalCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetSamplingIntervalCommand) Name() string {
	return "Profiler.setSamplingInterval"
}

func (cmd *SetSamplingIntervalCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetSamplingIntervalCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type ProfilerStartCB func(err error)

type ProfilerStartCommand struct {
	cb ProfilerStartCB
}

func NewProfilerStartCommand(cb ProfilerStartCB) *ProfilerStartCommand {
	return &ProfilerStartCommand{
		cb: cb,
	}
}

func (cmd *ProfilerStartCommand) Name() string {
	return "Profiler.start"
}

func (cmd *ProfilerStartCommand) Params() interface{} {
	return nil
}

func (cmd *ProfilerStartCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type StopResult struct {
	Profile *Profile `json:"profile"` // Recorded profile.
}

type StopCB func(result *StopResult, err error)

type StopCommand struct {
	cb StopCB
}

func NewStopCommand(cb StopCB) *StopCommand {
	return &StopCommand{
		cb: cb,
	}
}

func (cmd *StopCommand) Name() string {
	return "Profiler.stop"
}

func (cmd *StopCommand) Params() interface{} {
	return nil
}

func (cmd *StopCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj StopResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type ConsoleProfileStartedEvent struct {
	Id       string    `json:"id"`
	Location *Location `json:"location"` // Location of console.profile().
	Title    string    `json:"title"`    // Profile title passed as an argument to console.profile().
}

// Sent when new profile recodring is started using console.profile() call.
type ConsoleProfileStartedEventSink struct {
	events chan *ConsoleProfileStartedEvent
}

func NewConsoleProfileStartedEventSink(bufSize int) *ConsoleProfileStartedEventSink {
	return &ConsoleProfileStartedEventSink{
		events: make(chan *ConsoleProfileStartedEvent, bufSize),
	}
}

func (s *ConsoleProfileStartedEventSink) Name() string {
	return "Profiler.consoleProfileStarted"
}

func (s *ConsoleProfileStartedEventSink) OnEvent(params []byte) {
	evt := &ConsoleProfileStartedEvent{}
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

type ConsoleProfileFinishedEvent struct {
	Id       string    `json:"id"`
	Location *Location `json:"location"` // Location of console.profileEnd().
	Profile  *Profile  `json:"profile"`
	Title    string    `json:"title"` // Profile title passed as an argument to console.profile().
}

type ConsoleProfileFinishedEventSink struct {
	events chan *ConsoleProfileFinishedEvent
}

func NewConsoleProfileFinishedEventSink(bufSize int) *ConsoleProfileFinishedEventSink {
	return &ConsoleProfileFinishedEventSink{
		events: make(chan *ConsoleProfileFinishedEvent, bufSize),
	}
}

func (s *ConsoleProfileFinishedEventSink) Name() string {
	return "Profiler.consoleProfileFinished"
}

func (s *ConsoleProfileFinishedEventSink) OnEvent(params []byte) {
	evt := &ConsoleProfileFinishedEvent{}
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
