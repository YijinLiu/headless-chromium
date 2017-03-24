package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

// Profile node. Holds callsite information, execution statistics and child nodes.
type ProfileNode struct {
	Id            int                 `json:"id"`                      // Unique id of the node.
	CallFrame     *RuntimeCallFrame   `json:"callFrame"`               // Function location.
	HitCount      int                 `json:"hitCount,omitempty"`      // Number of samples where this node was on top of the call stack.
	Children      []int               `json:"children,omitempty"`      // Child node ids.
	DeoptReason   string              `json:"deoptReason,omitempty"`   // The reason of being not optimized. The function may be deoptimized or marked as don't optimize.
	PositionTicks []*PositionTickInfo `json:"positionTicks,omitempty"` // An array of source position ticks.
}

// Profile.
type Profile struct {
	Nodes      []*ProfileNode `json:"nodes"`                // The list of profile nodes. First item is the root node.
	StartTime  float64        `json:"startTime"`            // Profiling start timestamp in microseconds.
	EndTime    float64        `json:"endTime"`              // Profiling end timestamp in microseconds.
	Samples    []int          `json:"samples,omitempty"`    // Ids of samples top nodes.
	TimeDeltas []int          `json:"timeDeltas,omitempty"` // Time intervals between adjacent samples in microseconds. The first delta is relative to the profile startTime.
}

// Specifies a number of samples attributed to a certain source position.
// @experimental
type PositionTickInfo struct {
	Line  int `json:"line"`  // Source line number (1-based).
	Ticks int `json:"ticks"` // Number of samples attributed to the source line.
}

type ProfilerEnableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewProfilerEnableCommand() *ProfilerEnableCommand {
	return &ProfilerEnableCommand{}
}

func (cmd *ProfilerEnableCommand) Name() string {
	return "Profiler.enable"
}

func (cmd *ProfilerEnableCommand) Params() interface{} {
	return nil
}

func (cmd *ProfilerEnableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ProfilerEnable(conn *hc.Conn) (err error) {
	cmd := NewProfilerEnableCommand()
	cmd.Run(conn)
	return cmd.err
}

type ProfilerEnableCB func(err error)

type AsyncProfilerEnableCommand struct {
	cb ProfilerEnableCB
}

func NewAsyncProfilerEnableCommand(cb ProfilerEnableCB) *AsyncProfilerEnableCommand {
	return &AsyncProfilerEnableCommand{
		cb: cb,
	}
}

func (cmd *AsyncProfilerEnableCommand) Name() string {
	return "Profiler.enable"
}

func (cmd *AsyncProfilerEnableCommand) Params() interface{} {
	return nil
}

func (cmd *ProfilerEnableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncProfilerEnableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type ProfilerDisableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewProfilerDisableCommand() *ProfilerDisableCommand {
	return &ProfilerDisableCommand{}
}

func (cmd *ProfilerDisableCommand) Name() string {
	return "Profiler.disable"
}

func (cmd *ProfilerDisableCommand) Params() interface{} {
	return nil
}

func (cmd *ProfilerDisableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ProfilerDisable(conn *hc.Conn) (err error) {
	cmd := NewProfilerDisableCommand()
	cmd.Run(conn)
	return cmd.err
}

type ProfilerDisableCB func(err error)

type AsyncProfilerDisableCommand struct {
	cb ProfilerDisableCB
}

func NewAsyncProfilerDisableCommand(cb ProfilerDisableCB) *AsyncProfilerDisableCommand {
	return &AsyncProfilerDisableCommand{
		cb: cb,
	}
}

func (cmd *AsyncProfilerDisableCommand) Name() string {
	return "Profiler.disable"
}

func (cmd *AsyncProfilerDisableCommand) Params() interface{} {
	return nil
}

func (cmd *ProfilerDisableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncProfilerDisableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetSamplingIntervalParams struct {
	Interval int `json:"interval"` // New sampling interval in microseconds.
}

// Changes CPU profiler sampling interval. Must be called before CPU profiles recording started.

type SetSamplingIntervalCommand struct {
	params *SetSamplingIntervalParams
	wg     sync.WaitGroup
	err    error
}

func NewSetSamplingIntervalCommand(params *SetSamplingIntervalParams) *SetSamplingIntervalCommand {
	return &SetSamplingIntervalCommand{
		params: params,
	}
}

func (cmd *SetSamplingIntervalCommand) Name() string {
	return "Profiler.setSamplingInterval"
}

func (cmd *SetSamplingIntervalCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetSamplingIntervalCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetSamplingInterval(params *SetSamplingIntervalParams, conn *hc.Conn) (err error) {
	cmd := NewSetSamplingIntervalCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetSamplingIntervalCB func(err error)

// Changes CPU profiler sampling interval. Must be called before CPU profiles recording started.

type AsyncSetSamplingIntervalCommand struct {
	params *SetSamplingIntervalParams
	cb     SetSamplingIntervalCB
}

func NewAsyncSetSamplingIntervalCommand(params *SetSamplingIntervalParams, cb SetSamplingIntervalCB) *AsyncSetSamplingIntervalCommand {
	return &AsyncSetSamplingIntervalCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetSamplingIntervalCommand) Name() string {
	return "Profiler.setSamplingInterval"
}

func (cmd *AsyncSetSamplingIntervalCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetSamplingIntervalCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetSamplingIntervalCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type ProfilerStartCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewProfilerStartCommand() *ProfilerStartCommand {
	return &ProfilerStartCommand{}
}

func (cmd *ProfilerStartCommand) Name() string {
	return "Profiler.start"
}

func (cmd *ProfilerStartCommand) Params() interface{} {
	return nil
}

func (cmd *ProfilerStartCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ProfilerStart(conn *hc.Conn) (err error) {
	cmd := NewProfilerStartCommand()
	cmd.Run(conn)
	return cmd.err
}

type ProfilerStartCB func(err error)

type AsyncProfilerStartCommand struct {
	cb ProfilerStartCB
}

func NewAsyncProfilerStartCommand(cb ProfilerStartCB) *AsyncProfilerStartCommand {
	return &AsyncProfilerStartCommand{
		cb: cb,
	}
}

func (cmd *AsyncProfilerStartCommand) Name() string {
	return "Profiler.start"
}

func (cmd *AsyncProfilerStartCommand) Params() interface{} {
	return nil
}

func (cmd *ProfilerStartCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncProfilerStartCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type StopResult struct {
	Profile *Profile `json:"profile"` // Recorded profile.
}

type StopCommand struct {
	result StopResult
	wg     sync.WaitGroup
	err    error
}

func NewStopCommand() *StopCommand {
	return &StopCommand{}
}

func (cmd *StopCommand) Name() string {
	return "Profiler.stop"
}

func (cmd *StopCommand) Params() interface{} {
	return nil
}

func (cmd *StopCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func Stop(conn *hc.Conn) (result *StopResult, err error) {
	cmd := NewStopCommand()
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type StopCB func(result *StopResult, err error)

type AsyncStopCommand struct {
	cb StopCB
}

func NewAsyncStopCommand(cb StopCB) *AsyncStopCommand {
	return &AsyncStopCommand{
		cb: cb,
	}
}

func (cmd *AsyncStopCommand) Name() string {
	return "Profiler.stop"
}

func (cmd *AsyncStopCommand) Params() interface{} {
	return nil
}

func (cmd *StopCommand) Result() *StopResult {
	return &cmd.result
}

func (cmd *StopCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncStopCommand) Done(data []byte, err error) {
	var result StopResult
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

// Sent when new profile recodring is started using console.profile() call.

type ConsoleProfileStartedEvent struct {
	Id       string    `json:"id"`
	Location *Location `json:"location"` // Location of console.profile().
	Title    string    `json:"title"`    // Profile title passed as an argument to console.profile().
}

func OnConsoleProfileStarted(conn *hc.Conn, cb func(evt *ConsoleProfileStartedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &ConsoleProfileStartedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Profiler.consoleProfileStarted", sink)
}

type ConsoleProfileFinishedEvent struct {
	Id       string    `json:"id"`
	Location *Location `json:"location"` // Location of console.profileEnd().
	Profile  *Profile  `json:"profile"`
	Title    string    `json:"title"` // Profile title passed as an argument to console.profile().
}

func OnConsoleProfileFinished(conn *hc.Conn, cb func(evt *ConsoleProfileFinishedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &ConsoleProfileFinishedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Profiler.consoleProfileFinished", sink)
}
