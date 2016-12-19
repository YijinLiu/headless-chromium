package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
)

// Animation instance.
type Animation struct {
	Id           string           `json:"id"`           // Animation's id.
	Name         string           `json:"name"`         // Animation's name.
	PausedState  bool             `json:"pausedState"`  // Animation's internal paused state.
	PlayState    string           `json:"playState"`    // Animation's play state.
	PlaybackRate int              `json:"playbackRate"` // Animation's playback rate.
	StartTime    int              `json:"startTime"`    // Animation's start time.
	CurrentTime  int              `json:"currentTime"`  // Animation's current time.
	Source       *AnimationEffect `json:"source"`       // Animation's source animation node.
	Type         string           `json:"type"`         // Animation type of Animation.
	CssId        string           `json:"cssId"`        // A unique ID for Animation representing the sources that triggered this CSS animation/transition.
}

// AnimationEffect instance
type AnimationEffect struct {
	Delay          int            `json:"delay"`          // AnimationEffect's delay.
	EndDelay       int            `json:"endDelay"`       // AnimationEffect's end delay.
	IterationStart int            `json:"iterationStart"` // AnimationEffect's iteration start.
	Iterations     int            `json:"iterations"`     // AnimationEffect's iterations.
	Duration       int            `json:"duration"`       // AnimationEffect's iteration duration.
	Direction      string         `json:"direction"`      // AnimationEffect's playback direction.
	Fill           string         `json:"fill"`           // AnimationEffect's fill mode.
	BackendNodeId  *BackendNodeId `json:"backendNodeId"`  // AnimationEffect's target node.
	KeyframesRule  *KeyframesRule `json:"keyframesRule"`  // AnimationEffect's keyframes.
	Easing         string         `json:"easing"`         // AnimationEffect's timing function.
}

// Keyframes Rule
type KeyframesRule struct {
	Name      string           `json:"name"`      // CSS keyframed animation's name.
	Keyframes []*KeyframeStyle `json:"keyframes"` // List of animation keyframes.
}

// Keyframe Style
type KeyframeStyle struct {
	Offset string `json:"offset"` // Keyframe's time offset.
	Easing string `json:"easing"` // AnimationEffect's timing function.
}

type AnimationEnableCB func(err error)

// Enables animation domain notifications.
type AnimationEnableCommand struct {
	cb AnimationEnableCB
}

func NewAnimationEnableCommand(cb AnimationEnableCB) *AnimationEnableCommand {
	return &AnimationEnableCommand{
		cb: cb,
	}
}

func (cmd *AnimationEnableCommand) Name() string {
	return "Animation.enable"
}

func (cmd *AnimationEnableCommand) Params() interface{} {
	return nil
}

func (cmd *AnimationEnableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type AnimationDisableCB func(err error)

// Disables animation domain notifications.
type AnimationDisableCommand struct {
	cb AnimationDisableCB
}

func NewAnimationDisableCommand(cb AnimationDisableCB) *AnimationDisableCommand {
	return &AnimationDisableCommand{
		cb: cb,
	}
}

func (cmd *AnimationDisableCommand) Name() string {
	return "Animation.disable"
}

func (cmd *AnimationDisableCommand) Params() interface{} {
	return nil
}

func (cmd *AnimationDisableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type GetPlaybackRateResult struct {
	PlaybackRate int `json:"playbackRate"` // Playback rate for animations on page.
}

type GetPlaybackRateCB func(result *GetPlaybackRateResult, err error)

// Gets the playback rate of the document timeline.
type GetPlaybackRateCommand struct {
	cb GetPlaybackRateCB
}

func NewGetPlaybackRateCommand(cb GetPlaybackRateCB) *GetPlaybackRateCommand {
	return &GetPlaybackRateCommand{
		cb: cb,
	}
}

func (cmd *GetPlaybackRateCommand) Name() string {
	return "Animation.getPlaybackRate"
}

func (cmd *GetPlaybackRateCommand) Params() interface{} {
	return nil
}

func (cmd *GetPlaybackRateCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetPlaybackRateResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type SetPlaybackRateParams struct {
	PlaybackRate int `json:"playbackRate"` // Playback rate for animations on page
}

type SetPlaybackRateCB func(err error)

// Sets the playback rate of the document timeline.
type SetPlaybackRateCommand struct {
	params *SetPlaybackRateParams
	cb     SetPlaybackRateCB
}

func NewSetPlaybackRateCommand(params *SetPlaybackRateParams, cb SetPlaybackRateCB) *SetPlaybackRateCommand {
	return &SetPlaybackRateCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetPlaybackRateCommand) Name() string {
	return "Animation.setPlaybackRate"
}

func (cmd *SetPlaybackRateCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetPlaybackRateCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type GetCurrentTimeParams struct {
	Id string `json:"id"` // Id of animation.
}

type GetCurrentTimeResult struct {
	CurrentTime int `json:"currentTime"` // Current time of the page.
}

type GetCurrentTimeCB func(result *GetCurrentTimeResult, err error)

// Returns the current time of the an animation.
type GetCurrentTimeCommand struct {
	params *GetCurrentTimeParams
	cb     GetCurrentTimeCB
}

func NewGetCurrentTimeCommand(params *GetCurrentTimeParams, cb GetCurrentTimeCB) *GetCurrentTimeCommand {
	return &GetCurrentTimeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetCurrentTimeCommand) Name() string {
	return "Animation.getCurrentTime"
}

func (cmd *GetCurrentTimeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetCurrentTimeCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetCurrentTimeResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type SetPausedParams struct {
	Animations []string `json:"animations"` // Animations to set the pause state of.
	Paused     bool     `json:"paused"`     // Paused state to set to.
}

type SetPausedCB func(err error)

// Sets the paused state of a set of animations.
type SetPausedCommand struct {
	params *SetPausedParams
	cb     SetPausedCB
}

func NewSetPausedCommand(params *SetPausedParams, cb SetPausedCB) *SetPausedCommand {
	return &SetPausedCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetPausedCommand) Name() string {
	return "Animation.setPaused"
}

func (cmd *SetPausedCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetPausedCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetTimingParams struct {
	AnimationId string `json:"animationId"` // Animation id.
	Duration    int    `json:"duration"`    // Duration of the animation.
	Delay       int    `json:"delay"`       // Delay of the animation.
}

type SetTimingCB func(err error)

// Sets the timing of an animation node.
type SetTimingCommand struct {
	params *SetTimingParams
	cb     SetTimingCB
}

func NewSetTimingCommand(params *SetTimingParams, cb SetTimingCB) *SetTimingCommand {
	return &SetTimingCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetTimingCommand) Name() string {
	return "Animation.setTiming"
}

func (cmd *SetTimingCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetTimingCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SeekAnimationsParams struct {
	Animations  []string `json:"animations"`  // List of animation ids to seek.
	CurrentTime int      `json:"currentTime"` // Set the current time of each animation.
}

type SeekAnimationsCB func(err error)

// Seek a set of animations to a particular time within each animation.
type SeekAnimationsCommand struct {
	params *SeekAnimationsParams
	cb     SeekAnimationsCB
}

func NewSeekAnimationsCommand(params *SeekAnimationsParams, cb SeekAnimationsCB) *SeekAnimationsCommand {
	return &SeekAnimationsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SeekAnimationsCommand) Name() string {
	return "Animation.seekAnimations"
}

func (cmd *SeekAnimationsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SeekAnimationsCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type ReleaseAnimationsParams struct {
	Animations []string `json:"animations"` // List of animation ids to seek.
}

type ReleaseAnimationsCB func(err error)

// Releases a set of animations to no longer be manipulated.
type ReleaseAnimationsCommand struct {
	params *ReleaseAnimationsParams
	cb     ReleaseAnimationsCB
}

func NewReleaseAnimationsCommand(params *ReleaseAnimationsParams, cb ReleaseAnimationsCB) *ReleaseAnimationsCommand {
	return &ReleaseAnimationsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *ReleaseAnimationsCommand) Name() string {
	return "Animation.releaseAnimations"
}

func (cmd *ReleaseAnimationsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ReleaseAnimationsCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type ResolveAnimationParams struct {
	AnimationId string `json:"animationId"` // Animation id.
}

type ResolveAnimationResult struct {
	RemoteObject *RemoteObject `json:"remoteObject"` // Corresponding remote object.
}

type ResolveAnimationCB func(result *ResolveAnimationResult, err error)

// Gets the remote object of the Animation.
type ResolveAnimationCommand struct {
	params *ResolveAnimationParams
	cb     ResolveAnimationCB
}

func NewResolveAnimationCommand(params *ResolveAnimationParams, cb ResolveAnimationCB) *ResolveAnimationCommand {
	return &ResolveAnimationCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *ResolveAnimationCommand) Name() string {
	return "Animation.resolveAnimation"
}

func (cmd *ResolveAnimationCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ResolveAnimationCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj ResolveAnimationResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type AnimationCreatedEvent struct {
	Id string `json:"id"` // Id of the animation that was created.
}

// Event for each animation that has been created.
type AnimationCreatedEventSink struct {
	events chan *AnimationCreatedEvent
}

func NewAnimationCreatedEventSink(bufSize int) *AnimationCreatedEventSink {
	return &AnimationCreatedEventSink{
		events: make(chan *AnimationCreatedEvent, bufSize),
	}
}

func (s *AnimationCreatedEventSink) Name() string {
	return "Animation.animationCreated"
}

func (s *AnimationCreatedEventSink) OnEvent(params []byte) {
	evt := &AnimationCreatedEvent{}
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

type AnimationStartedEvent struct {
	Animation *Animation `json:"animation"` // Animation that was started.
}

// Event for animation that has been started.
type AnimationStartedEventSink struct {
	events chan *AnimationStartedEvent
}

func NewAnimationStartedEventSink(bufSize int) *AnimationStartedEventSink {
	return &AnimationStartedEventSink{
		events: make(chan *AnimationStartedEvent, bufSize),
	}
}

func (s *AnimationStartedEventSink) Name() string {
	return "Animation.animationStarted"
}

func (s *AnimationStartedEventSink) OnEvent(params []byte) {
	evt := &AnimationStartedEvent{}
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

type AnimationCanceledEvent struct {
	Id string `json:"id"` // Id of the animation that was cancelled.
}

// Event for when an animation has been cancelled.
type AnimationCanceledEventSink struct {
	events chan *AnimationCanceledEvent
}

func NewAnimationCanceledEventSink(bufSize int) *AnimationCanceledEventSink {
	return &AnimationCanceledEventSink{
		events: make(chan *AnimationCanceledEvent, bufSize),
	}
}

func (s *AnimationCanceledEventSink) Name() string {
	return "Animation.animationCanceled"
}

func (s *AnimationCanceledEventSink) OnEvent(params []byte) {
	evt := &AnimationCanceledEvent{}
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
