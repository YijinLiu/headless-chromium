package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

// Animation instance.
// @experimental
type Animation struct {
	Id           string           `json:"id"`              // Animation's id.
	Name         string           `json:"name"`            // Animation's name.
	PausedState  bool             `json:"pausedState"`     // Animation's internal paused state.
	PlayState    string           `json:"playState"`       // Animation's play state.
	PlaybackRate float64          `json:"playbackRate"`    // Animation's playback rate.
	StartTime    float64          `json:"startTime"`       // Animation's start time.
	CurrentTime  float64          `json:"currentTime"`     // Animation's current time.
	Source       *AnimationEffect `json:"source"`          // Animation's source animation node.
	Type         string           `json:"type"`            // Animation type of Animation.
	CssId        string           `json:"cssId,omitempty"` // A unique ID for Animation representing the sources that triggered this CSS animation/transition.
}

// AnimationEffect instance
// @experimental
type AnimationEffect struct {
	Delay          float64        `json:"delay"`                   // AnimationEffect's delay.
	EndDelay       float64        `json:"endDelay"`                // AnimationEffect's end delay.
	IterationStart float64        `json:"iterationStart"`          // AnimationEffect's iteration start.
	Iterations     float64        `json:"iterations"`              // AnimationEffect's iterations.
	Duration       float64        `json:"duration"`                // AnimationEffect's iteration duration.
	Direction      string         `json:"direction"`               // AnimationEffect's playback direction.
	Fill           string         `json:"fill"`                    // AnimationEffect's fill mode.
	BackendNodeId  *BackendNodeId `json:"backendNodeId"`           // AnimationEffect's target node.
	KeyframesRule  *KeyframesRule `json:"keyframesRule,omitempty"` // AnimationEffect's keyframes.
	Easing         string         `json:"easing"`                  // AnimationEffect's timing function.
}

// Keyframes Rule
type KeyframesRule struct {
	Name      string           `json:"name,omitempty"` // CSS keyframed animation's name.
	Keyframes []*KeyframeStyle `json:"keyframes"`      // List of animation keyframes.
}

// Keyframe Style
type KeyframeStyle struct {
	Offset string `json:"offset"` // Keyframe's time offset.
	Easing string `json:"easing"` // AnimationEffect's timing function.
}

// Enables animation domain notifications.

type AnimationEnableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewAnimationEnableCommand() *AnimationEnableCommand {
	return &AnimationEnableCommand{}
}

func (cmd *AnimationEnableCommand) Name() string {
	return "Animation.enable"
}

func (cmd *AnimationEnableCommand) Params() interface{} {
	return nil
}

func (cmd *AnimationEnableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func AnimationEnable(conn *hc.Conn) (err error) {
	cmd := NewAnimationEnableCommand()
	cmd.Run(conn)
	return cmd.err
}

type AnimationEnableCB func(err error)

// Enables animation domain notifications.

type AsyncAnimationEnableCommand struct {
	cb AnimationEnableCB
}

func NewAsyncAnimationEnableCommand(cb AnimationEnableCB) *AsyncAnimationEnableCommand {
	return &AsyncAnimationEnableCommand{
		cb: cb,
	}
}

func (cmd *AsyncAnimationEnableCommand) Name() string {
	return "Animation.enable"
}

func (cmd *AsyncAnimationEnableCommand) Params() interface{} {
	return nil
}

func (cmd *AnimationEnableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncAnimationEnableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Disables animation domain notifications.

type AnimationDisableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewAnimationDisableCommand() *AnimationDisableCommand {
	return &AnimationDisableCommand{}
}

func (cmd *AnimationDisableCommand) Name() string {
	return "Animation.disable"
}

func (cmd *AnimationDisableCommand) Params() interface{} {
	return nil
}

func (cmd *AnimationDisableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func AnimationDisable(conn *hc.Conn) (err error) {
	cmd := NewAnimationDisableCommand()
	cmd.Run(conn)
	return cmd.err
}

type AnimationDisableCB func(err error)

// Disables animation domain notifications.

type AsyncAnimationDisableCommand struct {
	cb AnimationDisableCB
}

func NewAsyncAnimationDisableCommand(cb AnimationDisableCB) *AsyncAnimationDisableCommand {
	return &AsyncAnimationDisableCommand{
		cb: cb,
	}
}

func (cmd *AsyncAnimationDisableCommand) Name() string {
	return "Animation.disable"
}

func (cmd *AsyncAnimationDisableCommand) Params() interface{} {
	return nil
}

func (cmd *AnimationDisableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncAnimationDisableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type GetPlaybackRateResult struct {
	PlaybackRate float64 `json:"playbackRate"` // Playback rate for animations on page.
}

// Gets the playback rate of the document timeline.

type GetPlaybackRateCommand struct {
	result GetPlaybackRateResult
	wg     sync.WaitGroup
	err    error
}

func NewGetPlaybackRateCommand() *GetPlaybackRateCommand {
	return &GetPlaybackRateCommand{}
}

func (cmd *GetPlaybackRateCommand) Name() string {
	return "Animation.getPlaybackRate"
}

func (cmd *GetPlaybackRateCommand) Params() interface{} {
	return nil
}

func (cmd *GetPlaybackRateCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetPlaybackRate(conn *hc.Conn) (result *GetPlaybackRateResult, err error) {
	cmd := NewGetPlaybackRateCommand()
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetPlaybackRateCB func(result *GetPlaybackRateResult, err error)

// Gets the playback rate of the document timeline.

type AsyncGetPlaybackRateCommand struct {
	cb GetPlaybackRateCB
}

func NewAsyncGetPlaybackRateCommand(cb GetPlaybackRateCB) *AsyncGetPlaybackRateCommand {
	return &AsyncGetPlaybackRateCommand{
		cb: cb,
	}
}

func (cmd *AsyncGetPlaybackRateCommand) Name() string {
	return "Animation.getPlaybackRate"
}

func (cmd *AsyncGetPlaybackRateCommand) Params() interface{} {
	return nil
}

func (cmd *GetPlaybackRateCommand) Result() *GetPlaybackRateResult {
	return &cmd.result
}

func (cmd *GetPlaybackRateCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetPlaybackRateCommand) Done(data []byte, err error) {
	var result GetPlaybackRateResult
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

type SetPlaybackRateParams struct {
	PlaybackRate float64 `json:"playbackRate"` // Playback rate for animations on page
}

// Sets the playback rate of the document timeline.

type SetPlaybackRateCommand struct {
	params *SetPlaybackRateParams
	wg     sync.WaitGroup
	err    error
}

func NewSetPlaybackRateCommand(params *SetPlaybackRateParams) *SetPlaybackRateCommand {
	return &SetPlaybackRateCommand{
		params: params,
	}
}

func (cmd *SetPlaybackRateCommand) Name() string {
	return "Animation.setPlaybackRate"
}

func (cmd *SetPlaybackRateCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetPlaybackRateCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetPlaybackRate(params *SetPlaybackRateParams, conn *hc.Conn) (err error) {
	cmd := NewSetPlaybackRateCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetPlaybackRateCB func(err error)

// Sets the playback rate of the document timeline.

type AsyncSetPlaybackRateCommand struct {
	params *SetPlaybackRateParams
	cb     SetPlaybackRateCB
}

func NewAsyncSetPlaybackRateCommand(params *SetPlaybackRateParams, cb SetPlaybackRateCB) *AsyncSetPlaybackRateCommand {
	return &AsyncSetPlaybackRateCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetPlaybackRateCommand) Name() string {
	return "Animation.setPlaybackRate"
}

func (cmd *AsyncSetPlaybackRateCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetPlaybackRateCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetPlaybackRateCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type GetCurrentTimeParams struct {
	Id string `json:"id"` // Id of animation.
}

type GetCurrentTimeResult struct {
	CurrentTime float64 `json:"currentTime"` // Current time of the page.
}

// Returns the current time of the an animation.

type GetCurrentTimeCommand struct {
	params *GetCurrentTimeParams
	result GetCurrentTimeResult
	wg     sync.WaitGroup
	err    error
}

func NewGetCurrentTimeCommand(params *GetCurrentTimeParams) *GetCurrentTimeCommand {
	return &GetCurrentTimeCommand{
		params: params,
	}
}

func (cmd *GetCurrentTimeCommand) Name() string {
	return "Animation.getCurrentTime"
}

func (cmd *GetCurrentTimeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetCurrentTimeCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetCurrentTime(params *GetCurrentTimeParams, conn *hc.Conn) (result *GetCurrentTimeResult, err error) {
	cmd := NewGetCurrentTimeCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetCurrentTimeCB func(result *GetCurrentTimeResult, err error)

// Returns the current time of the an animation.

type AsyncGetCurrentTimeCommand struct {
	params *GetCurrentTimeParams
	cb     GetCurrentTimeCB
}

func NewAsyncGetCurrentTimeCommand(params *GetCurrentTimeParams, cb GetCurrentTimeCB) *AsyncGetCurrentTimeCommand {
	return &AsyncGetCurrentTimeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetCurrentTimeCommand) Name() string {
	return "Animation.getCurrentTime"
}

func (cmd *AsyncGetCurrentTimeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetCurrentTimeCommand) Result() *GetCurrentTimeResult {
	return &cmd.result
}

func (cmd *GetCurrentTimeCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetCurrentTimeCommand) Done(data []byte, err error) {
	var result GetCurrentTimeResult
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

type SetPausedParams struct {
	Animations []string `json:"animations"` // Animations to set the pause state of.
	Paused     bool     `json:"paused"`     // Paused state to set to.
}

// Sets the paused state of a set of animations.

type SetPausedCommand struct {
	params *SetPausedParams
	wg     sync.WaitGroup
	err    error
}

func NewSetPausedCommand(params *SetPausedParams) *SetPausedCommand {
	return &SetPausedCommand{
		params: params,
	}
}

func (cmd *SetPausedCommand) Name() string {
	return "Animation.setPaused"
}

func (cmd *SetPausedCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetPausedCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetPaused(params *SetPausedParams, conn *hc.Conn) (err error) {
	cmd := NewSetPausedCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetPausedCB func(err error)

// Sets the paused state of a set of animations.

type AsyncSetPausedCommand struct {
	params *SetPausedParams
	cb     SetPausedCB
}

func NewAsyncSetPausedCommand(params *SetPausedParams, cb SetPausedCB) *AsyncSetPausedCommand {
	return &AsyncSetPausedCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetPausedCommand) Name() string {
	return "Animation.setPaused"
}

func (cmd *AsyncSetPausedCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetPausedCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetPausedCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetTimingParams struct {
	AnimationId string  `json:"animationId"` // Animation id.
	Duration    float64 `json:"duration"`    // Duration of the animation.
	Delay       float64 `json:"delay"`       // Delay of the animation.
}

// Sets the timing of an animation node.

type SetTimingCommand struct {
	params *SetTimingParams
	wg     sync.WaitGroup
	err    error
}

func NewSetTimingCommand(params *SetTimingParams) *SetTimingCommand {
	return &SetTimingCommand{
		params: params,
	}
}

func (cmd *SetTimingCommand) Name() string {
	return "Animation.setTiming"
}

func (cmd *SetTimingCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetTimingCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetTiming(params *SetTimingParams, conn *hc.Conn) (err error) {
	cmd := NewSetTimingCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetTimingCB func(err error)

// Sets the timing of an animation node.

type AsyncSetTimingCommand struct {
	params *SetTimingParams
	cb     SetTimingCB
}

func NewAsyncSetTimingCommand(params *SetTimingParams, cb SetTimingCB) *AsyncSetTimingCommand {
	return &AsyncSetTimingCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetTimingCommand) Name() string {
	return "Animation.setTiming"
}

func (cmd *AsyncSetTimingCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetTimingCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetTimingCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SeekAnimationsParams struct {
	Animations  []string `json:"animations"`  // List of animation ids to seek.
	CurrentTime float64  `json:"currentTime"` // Set the current time of each animation.
}

// Seek a set of animations to a particular time within each animation.

type SeekAnimationsCommand struct {
	params *SeekAnimationsParams
	wg     sync.WaitGroup
	err    error
}

func NewSeekAnimationsCommand(params *SeekAnimationsParams) *SeekAnimationsCommand {
	return &SeekAnimationsCommand{
		params: params,
	}
}

func (cmd *SeekAnimationsCommand) Name() string {
	return "Animation.seekAnimations"
}

func (cmd *SeekAnimationsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SeekAnimationsCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SeekAnimations(params *SeekAnimationsParams, conn *hc.Conn) (err error) {
	cmd := NewSeekAnimationsCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SeekAnimationsCB func(err error)

// Seek a set of animations to a particular time within each animation.

type AsyncSeekAnimationsCommand struct {
	params *SeekAnimationsParams
	cb     SeekAnimationsCB
}

func NewAsyncSeekAnimationsCommand(params *SeekAnimationsParams, cb SeekAnimationsCB) *AsyncSeekAnimationsCommand {
	return &AsyncSeekAnimationsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSeekAnimationsCommand) Name() string {
	return "Animation.seekAnimations"
}

func (cmd *AsyncSeekAnimationsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SeekAnimationsCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSeekAnimationsCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type ReleaseAnimationsParams struct {
	Animations []string `json:"animations"` // List of animation ids to seek.
}

// Releases a set of animations to no longer be manipulated.

type ReleaseAnimationsCommand struct {
	params *ReleaseAnimationsParams
	wg     sync.WaitGroup
	err    error
}

func NewReleaseAnimationsCommand(params *ReleaseAnimationsParams) *ReleaseAnimationsCommand {
	return &ReleaseAnimationsCommand{
		params: params,
	}
}

func (cmd *ReleaseAnimationsCommand) Name() string {
	return "Animation.releaseAnimations"
}

func (cmd *ReleaseAnimationsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ReleaseAnimationsCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ReleaseAnimations(params *ReleaseAnimationsParams, conn *hc.Conn) (err error) {
	cmd := NewReleaseAnimationsCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type ReleaseAnimationsCB func(err error)

// Releases a set of animations to no longer be manipulated.

type AsyncReleaseAnimationsCommand struct {
	params *ReleaseAnimationsParams
	cb     ReleaseAnimationsCB
}

func NewAsyncReleaseAnimationsCommand(params *ReleaseAnimationsParams, cb ReleaseAnimationsCB) *AsyncReleaseAnimationsCommand {
	return &AsyncReleaseAnimationsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncReleaseAnimationsCommand) Name() string {
	return "Animation.releaseAnimations"
}

func (cmd *AsyncReleaseAnimationsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ReleaseAnimationsCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncReleaseAnimationsCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type ResolveAnimationParams struct {
	AnimationId string `json:"animationId"` // Animation id.
}

type ResolveAnimationResult struct {
	RemoteObject *RemoteObject `json:"remoteObject"` // Corresponding remote object.
}

// Gets the remote object of the Animation.

type ResolveAnimationCommand struct {
	params *ResolveAnimationParams
	result ResolveAnimationResult
	wg     sync.WaitGroup
	err    error
}

func NewResolveAnimationCommand(params *ResolveAnimationParams) *ResolveAnimationCommand {
	return &ResolveAnimationCommand{
		params: params,
	}
}

func (cmd *ResolveAnimationCommand) Name() string {
	return "Animation.resolveAnimation"
}

func (cmd *ResolveAnimationCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ResolveAnimationCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ResolveAnimation(params *ResolveAnimationParams, conn *hc.Conn) (result *ResolveAnimationResult, err error) {
	cmd := NewResolveAnimationCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type ResolveAnimationCB func(result *ResolveAnimationResult, err error)

// Gets the remote object of the Animation.

type AsyncResolveAnimationCommand struct {
	params *ResolveAnimationParams
	cb     ResolveAnimationCB
}

func NewAsyncResolveAnimationCommand(params *ResolveAnimationParams, cb ResolveAnimationCB) *AsyncResolveAnimationCommand {
	return &AsyncResolveAnimationCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncResolveAnimationCommand) Name() string {
	return "Animation.resolveAnimation"
}

func (cmd *AsyncResolveAnimationCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ResolveAnimationCommand) Result() *ResolveAnimationResult {
	return &cmd.result
}

func (cmd *ResolveAnimationCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncResolveAnimationCommand) Done(data []byte, err error) {
	var result ResolveAnimationResult
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

// Event for each animation that has been created.

type AnimationCreatedEvent struct {
	Id string `json:"id"` // Id of the animation that was created.
}

func OnAnimationCreated(conn *hc.Conn, cb func(evt *AnimationCreatedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &AnimationCreatedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Animation.animationCreated", sink)
}

// Event for animation that has been started.

type AnimationStartedEvent struct {
	Animation *Animation `json:"animation"` // Animation that was started.
}

func OnAnimationStarted(conn *hc.Conn, cb func(evt *AnimationStartedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &AnimationStartedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Animation.animationStarted", sink)
}

// Event for when an animation has been cancelled.

type AnimationCanceledEvent struct {
	Id string `json:"id"` // Id of the animation that was cancelled.
}

func OnAnimationCanceled(conn *hc.Conn, cb func(evt *AnimationCanceledEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &AnimationCanceledEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Animation.animationCanceled", sink)
}
