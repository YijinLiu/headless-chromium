package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
)

// ServiceWorker registration.
type ServiceWorkerRegistration struct {
	RegistrationId string `json:"registrationId"`
	ScopeURL       string `json:"scopeURL"`
	IsDeleted      bool   `json:"isDeleted"`
}

type ServiceWorkerVersionRunningStatus string

const ServiceWorkerVersionRunningStatusStopped ServiceWorkerVersionRunningStatus = "stopped"
const ServiceWorkerVersionRunningStatusStarting ServiceWorkerVersionRunningStatus = "starting"
const ServiceWorkerVersionRunningStatusRunning ServiceWorkerVersionRunningStatus = "running"
const ServiceWorkerVersionRunningStatusStopping ServiceWorkerVersionRunningStatus = "stopping"

type ServiceWorkerVersionStatus string

const ServiceWorkerVersionStatusNew ServiceWorkerVersionStatus = "new"
const ServiceWorkerVersionStatusInstalling ServiceWorkerVersionStatus = "installing"
const ServiceWorkerVersionStatusInstalled ServiceWorkerVersionStatus = "installed"
const ServiceWorkerVersionStatusActivating ServiceWorkerVersionStatus = "activating"
const ServiceWorkerVersionStatusActivated ServiceWorkerVersionStatus = "activated"
const ServiceWorkerVersionStatusRedundant ServiceWorkerVersionStatus = "redundant"

// ServiceWorker version.
type ServiceWorkerVersion struct {
	VersionId          string                            `json:"versionId"`
	RegistrationId     string                            `json:"registrationId"`
	ScriptURL          string                            `json:"scriptURL"`
	RunningStatus      ServiceWorkerVersionRunningStatus `json:"runningStatus"`
	Status             ServiceWorkerVersionStatus        `json:"status"`
	ScriptLastModified int                               `json:"scriptLastModified"` // The Last-Modified header value of the main script.
	ScriptResponseTime int                               `json:"scriptResponseTime"` // The time at which the response headers of the main script were received from the server.  For cached script it is the last time the cache entry was validated.
	ControlledClients  []*TargetTargetID                 `json:"controlledClients"`
	TargetId           *TargetTargetID                   `json:"targetId"`
}

// ServiceWorker error message.
type ServiceWorkerErrorMessage struct {
	ErrorMessage   string `json:"errorMessage"`
	RegistrationId string `json:"registrationId"`
	VersionId      string `json:"versionId"`
	SourceURL      string `json:"sourceURL"`
	LineNumber     int    `json:"lineNumber"`
	ColumnNumber   int    `json:"columnNumber"`
}

type ServiceWorkerEnableCB func(err error)

type ServiceWorkerEnableCommand struct {
	cb ServiceWorkerEnableCB
}

func NewServiceWorkerEnableCommand(cb ServiceWorkerEnableCB) *ServiceWorkerEnableCommand {
	return &ServiceWorkerEnableCommand{
		cb: cb,
	}
}

func (cmd *ServiceWorkerEnableCommand) Name() string {
	return "ServiceWorker.enable"
}

func (cmd *ServiceWorkerEnableCommand) Params() interface{} {
	return nil
}

func (cmd *ServiceWorkerEnableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type ServiceWorkerDisableCB func(err error)

type ServiceWorkerDisableCommand struct {
	cb ServiceWorkerDisableCB
}

func NewServiceWorkerDisableCommand(cb ServiceWorkerDisableCB) *ServiceWorkerDisableCommand {
	return &ServiceWorkerDisableCommand{
		cb: cb,
	}
}

func (cmd *ServiceWorkerDisableCommand) Name() string {
	return "ServiceWorker.disable"
}

func (cmd *ServiceWorkerDisableCommand) Params() interface{} {
	return nil
}

func (cmd *ServiceWorkerDisableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type UnregisterParams struct {
	ScopeURL string `json:"scopeURL"`
}

type UnregisterCB func(err error)

type UnregisterCommand struct {
	params *UnregisterParams
	cb     UnregisterCB
}

func NewUnregisterCommand(params *UnregisterParams, cb UnregisterCB) *UnregisterCommand {
	return &UnregisterCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *UnregisterCommand) Name() string {
	return "ServiceWorker.unregister"
}

func (cmd *UnregisterCommand) Params() interface{} {
	return cmd.params
}

func (cmd *UnregisterCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type UpdateRegistrationParams struct {
	ScopeURL string `json:"scopeURL"`
}

type UpdateRegistrationCB func(err error)

type UpdateRegistrationCommand struct {
	params *UpdateRegistrationParams
	cb     UpdateRegistrationCB
}

func NewUpdateRegistrationCommand(params *UpdateRegistrationParams, cb UpdateRegistrationCB) *UpdateRegistrationCommand {
	return &UpdateRegistrationCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *UpdateRegistrationCommand) Name() string {
	return "ServiceWorker.updateRegistration"
}

func (cmd *UpdateRegistrationCommand) Params() interface{} {
	return cmd.params
}

func (cmd *UpdateRegistrationCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type StartWorkerParams struct {
	ScopeURL string `json:"scopeURL"`
}

type StartWorkerCB func(err error)

type StartWorkerCommand struct {
	params *StartWorkerParams
	cb     StartWorkerCB
}

func NewStartWorkerCommand(params *StartWorkerParams, cb StartWorkerCB) *StartWorkerCommand {
	return &StartWorkerCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *StartWorkerCommand) Name() string {
	return "ServiceWorker.startWorker"
}

func (cmd *StartWorkerCommand) Params() interface{} {
	return cmd.params
}

func (cmd *StartWorkerCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SkipWaitingParams struct {
	ScopeURL string `json:"scopeURL"`
}

type SkipWaitingCB func(err error)

type SkipWaitingCommand struct {
	params *SkipWaitingParams
	cb     SkipWaitingCB
}

func NewSkipWaitingCommand(params *SkipWaitingParams, cb SkipWaitingCB) *SkipWaitingCommand {
	return &SkipWaitingCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SkipWaitingCommand) Name() string {
	return "ServiceWorker.skipWaiting"
}

func (cmd *SkipWaitingCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SkipWaitingCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type StopWorkerParams struct {
	VersionId string `json:"versionId"`
}

type StopWorkerCB func(err error)

type StopWorkerCommand struct {
	params *StopWorkerParams
	cb     StopWorkerCB
}

func NewStopWorkerCommand(params *StopWorkerParams, cb StopWorkerCB) *StopWorkerCommand {
	return &StopWorkerCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *StopWorkerCommand) Name() string {
	return "ServiceWorker.stopWorker"
}

func (cmd *StopWorkerCommand) Params() interface{} {
	return cmd.params
}

func (cmd *StopWorkerCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type InspectWorkerParams struct {
	VersionId string `json:"versionId"`
}

type InspectWorkerCB func(err error)

type InspectWorkerCommand struct {
	params *InspectWorkerParams
	cb     InspectWorkerCB
}

func NewInspectWorkerCommand(params *InspectWorkerParams, cb InspectWorkerCB) *InspectWorkerCommand {
	return &InspectWorkerCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *InspectWorkerCommand) Name() string {
	return "ServiceWorker.inspectWorker"
}

func (cmd *InspectWorkerCommand) Params() interface{} {
	return cmd.params
}

func (cmd *InspectWorkerCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetForceUpdateOnPageLoadParams struct {
	ForceUpdateOnPageLoad bool `json:"forceUpdateOnPageLoad"`
}

type SetForceUpdateOnPageLoadCB func(err error)

type SetForceUpdateOnPageLoadCommand struct {
	params *SetForceUpdateOnPageLoadParams
	cb     SetForceUpdateOnPageLoadCB
}

func NewSetForceUpdateOnPageLoadCommand(params *SetForceUpdateOnPageLoadParams, cb SetForceUpdateOnPageLoadCB) *SetForceUpdateOnPageLoadCommand {
	return &SetForceUpdateOnPageLoadCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetForceUpdateOnPageLoadCommand) Name() string {
	return "ServiceWorker.setForceUpdateOnPageLoad"
}

func (cmd *SetForceUpdateOnPageLoadCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetForceUpdateOnPageLoadCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type DeliverPushMessageParams struct {
	Origin         string `json:"origin"`
	RegistrationId string `json:"registrationId"`
	Data           string `json:"data"`
}

type DeliverPushMessageCB func(err error)

type DeliverPushMessageCommand struct {
	params *DeliverPushMessageParams
	cb     DeliverPushMessageCB
}

func NewDeliverPushMessageCommand(params *DeliverPushMessageParams, cb DeliverPushMessageCB) *DeliverPushMessageCommand {
	return &DeliverPushMessageCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *DeliverPushMessageCommand) Name() string {
	return "ServiceWorker.deliverPushMessage"
}

func (cmd *DeliverPushMessageCommand) Params() interface{} {
	return cmd.params
}

func (cmd *DeliverPushMessageCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type DispatchSyncEventParams struct {
	Origin         string `json:"origin"`
	RegistrationId string `json:"registrationId"`
	Tag            string `json:"tag"`
	LastChance     bool   `json:"lastChance"`
}

type DispatchSyncEventCB func(err error)

type DispatchSyncEventCommand struct {
	params *DispatchSyncEventParams
	cb     DispatchSyncEventCB
}

func NewDispatchSyncEventCommand(params *DispatchSyncEventParams, cb DispatchSyncEventCB) *DispatchSyncEventCommand {
	return &DispatchSyncEventCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *DispatchSyncEventCommand) Name() string {
	return "ServiceWorker.dispatchSyncEvent"
}

func (cmd *DispatchSyncEventCommand) Params() interface{} {
	return cmd.params
}

func (cmd *DispatchSyncEventCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type WorkerRegistrationUpdatedEvent struct {
	Registrations []*ServiceWorkerRegistration `json:"registrations"`
}

type WorkerRegistrationUpdatedEventSink struct {
	events chan *WorkerRegistrationUpdatedEvent
}

func NewWorkerRegistrationUpdatedEventSink(bufSize int) *WorkerRegistrationUpdatedEventSink {
	return &WorkerRegistrationUpdatedEventSink{
		events: make(chan *WorkerRegistrationUpdatedEvent, bufSize),
	}
}

func (s *WorkerRegistrationUpdatedEventSink) Name() string {
	return "ServiceWorker.workerRegistrationUpdated"
}

func (s *WorkerRegistrationUpdatedEventSink) OnEvent(params []byte) {
	evt := &WorkerRegistrationUpdatedEvent{}
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

type WorkerVersionUpdatedEvent struct {
	Versions []*ServiceWorkerVersion `json:"versions"`
}

type WorkerVersionUpdatedEventSink struct {
	events chan *WorkerVersionUpdatedEvent
}

func NewWorkerVersionUpdatedEventSink(bufSize int) *WorkerVersionUpdatedEventSink {
	return &WorkerVersionUpdatedEventSink{
		events: make(chan *WorkerVersionUpdatedEvent, bufSize),
	}
}

func (s *WorkerVersionUpdatedEventSink) Name() string {
	return "ServiceWorker.workerVersionUpdated"
}

func (s *WorkerVersionUpdatedEventSink) OnEvent(params []byte) {
	evt := &WorkerVersionUpdatedEvent{}
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

type WorkerErrorReportedEvent struct {
	ErrorMessage *ServiceWorkerErrorMessage `json:"errorMessage"`
}

type WorkerErrorReportedEventSink struct {
	events chan *WorkerErrorReportedEvent
}

func NewWorkerErrorReportedEventSink(bufSize int) *WorkerErrorReportedEventSink {
	return &WorkerErrorReportedEventSink{
		events: make(chan *WorkerErrorReportedEvent, bufSize),
	}
}

func (s *WorkerErrorReportedEventSink) Name() string {
	return "ServiceWorker.workerErrorReported"
}

func (s *WorkerErrorReportedEventSink) OnEvent(params []byte) {
	evt := &WorkerErrorReportedEvent{}
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
