package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
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
	ScriptLastModified float64                           `json:"scriptLastModified,omitempty"` // The Last-Modified header value of the main script.
	ScriptResponseTime float64                           `json:"scriptResponseTime,omitempty"` // The time at which the response headers of the main script were received from the server.  For cached script it is the last time the cache entry was validated.
	ControlledClients  []*TargetID                       `json:"controlledClients,omitempty"`
	TargetId           *TargetID                         `json:"targetId,omitempty"`
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

type ServiceWorkerEnableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewServiceWorkerEnableCommand() *ServiceWorkerEnableCommand {
	return &ServiceWorkerEnableCommand{}
}

func (cmd *ServiceWorkerEnableCommand) Name() string {
	return "ServiceWorker.enable"
}

func (cmd *ServiceWorkerEnableCommand) Params() interface{} {
	return nil
}

func (cmd *ServiceWorkerEnableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ServiceWorkerEnable(conn *hc.Conn) (err error) {
	cmd := NewServiceWorkerEnableCommand()
	cmd.Run(conn)
	return cmd.err
}

type ServiceWorkerEnableCB func(err error)

type AsyncServiceWorkerEnableCommand struct {
	cb ServiceWorkerEnableCB
}

func NewAsyncServiceWorkerEnableCommand(cb ServiceWorkerEnableCB) *AsyncServiceWorkerEnableCommand {
	return &AsyncServiceWorkerEnableCommand{
		cb: cb,
	}
}

func (cmd *AsyncServiceWorkerEnableCommand) Name() string {
	return "ServiceWorker.enable"
}

func (cmd *AsyncServiceWorkerEnableCommand) Params() interface{} {
	return nil
}

func (cmd *ServiceWorkerEnableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncServiceWorkerEnableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type ServiceWorkerDisableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewServiceWorkerDisableCommand() *ServiceWorkerDisableCommand {
	return &ServiceWorkerDisableCommand{}
}

func (cmd *ServiceWorkerDisableCommand) Name() string {
	return "ServiceWorker.disable"
}

func (cmd *ServiceWorkerDisableCommand) Params() interface{} {
	return nil
}

func (cmd *ServiceWorkerDisableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ServiceWorkerDisable(conn *hc.Conn) (err error) {
	cmd := NewServiceWorkerDisableCommand()
	cmd.Run(conn)
	return cmd.err
}

type ServiceWorkerDisableCB func(err error)

type AsyncServiceWorkerDisableCommand struct {
	cb ServiceWorkerDisableCB
}

func NewAsyncServiceWorkerDisableCommand(cb ServiceWorkerDisableCB) *AsyncServiceWorkerDisableCommand {
	return &AsyncServiceWorkerDisableCommand{
		cb: cb,
	}
}

func (cmd *AsyncServiceWorkerDisableCommand) Name() string {
	return "ServiceWorker.disable"
}

func (cmd *AsyncServiceWorkerDisableCommand) Params() interface{} {
	return nil
}

func (cmd *ServiceWorkerDisableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncServiceWorkerDisableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type UnregisterParams struct {
	ScopeURL string `json:"scopeURL"`
}

type UnregisterCommand struct {
	params *UnregisterParams
	wg     sync.WaitGroup
	err    error
}

func NewUnregisterCommand(params *UnregisterParams) *UnregisterCommand {
	return &UnregisterCommand{
		params: params,
	}
}

func (cmd *UnregisterCommand) Name() string {
	return "ServiceWorker.unregister"
}

func (cmd *UnregisterCommand) Params() interface{} {
	return cmd.params
}

func (cmd *UnregisterCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func Unregister(params *UnregisterParams, conn *hc.Conn) (err error) {
	cmd := NewUnregisterCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type UnregisterCB func(err error)

type AsyncUnregisterCommand struct {
	params *UnregisterParams
	cb     UnregisterCB
}

func NewAsyncUnregisterCommand(params *UnregisterParams, cb UnregisterCB) *AsyncUnregisterCommand {
	return &AsyncUnregisterCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncUnregisterCommand) Name() string {
	return "ServiceWorker.unregister"
}

func (cmd *AsyncUnregisterCommand) Params() interface{} {
	return cmd.params
}

func (cmd *UnregisterCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncUnregisterCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type UpdateRegistrationParams struct {
	ScopeURL string `json:"scopeURL"`
}

type UpdateRegistrationCommand struct {
	params *UpdateRegistrationParams
	wg     sync.WaitGroup
	err    error
}

func NewUpdateRegistrationCommand(params *UpdateRegistrationParams) *UpdateRegistrationCommand {
	return &UpdateRegistrationCommand{
		params: params,
	}
}

func (cmd *UpdateRegistrationCommand) Name() string {
	return "ServiceWorker.updateRegistration"
}

func (cmd *UpdateRegistrationCommand) Params() interface{} {
	return cmd.params
}

func (cmd *UpdateRegistrationCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func UpdateRegistration(params *UpdateRegistrationParams, conn *hc.Conn) (err error) {
	cmd := NewUpdateRegistrationCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type UpdateRegistrationCB func(err error)

type AsyncUpdateRegistrationCommand struct {
	params *UpdateRegistrationParams
	cb     UpdateRegistrationCB
}

func NewAsyncUpdateRegistrationCommand(params *UpdateRegistrationParams, cb UpdateRegistrationCB) *AsyncUpdateRegistrationCommand {
	return &AsyncUpdateRegistrationCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncUpdateRegistrationCommand) Name() string {
	return "ServiceWorker.updateRegistration"
}

func (cmd *AsyncUpdateRegistrationCommand) Params() interface{} {
	return cmd.params
}

func (cmd *UpdateRegistrationCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncUpdateRegistrationCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type StartWorkerParams struct {
	ScopeURL string `json:"scopeURL"`
}

type StartWorkerCommand struct {
	params *StartWorkerParams
	wg     sync.WaitGroup
	err    error
}

func NewStartWorkerCommand(params *StartWorkerParams) *StartWorkerCommand {
	return &StartWorkerCommand{
		params: params,
	}
}

func (cmd *StartWorkerCommand) Name() string {
	return "ServiceWorker.startWorker"
}

func (cmd *StartWorkerCommand) Params() interface{} {
	return cmd.params
}

func (cmd *StartWorkerCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func StartWorker(params *StartWorkerParams, conn *hc.Conn) (err error) {
	cmd := NewStartWorkerCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type StartWorkerCB func(err error)

type AsyncStartWorkerCommand struct {
	params *StartWorkerParams
	cb     StartWorkerCB
}

func NewAsyncStartWorkerCommand(params *StartWorkerParams, cb StartWorkerCB) *AsyncStartWorkerCommand {
	return &AsyncStartWorkerCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncStartWorkerCommand) Name() string {
	return "ServiceWorker.startWorker"
}

func (cmd *AsyncStartWorkerCommand) Params() interface{} {
	return cmd.params
}

func (cmd *StartWorkerCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncStartWorkerCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SkipWaitingParams struct {
	ScopeURL string `json:"scopeURL"`
}

type SkipWaitingCommand struct {
	params *SkipWaitingParams
	wg     sync.WaitGroup
	err    error
}

func NewSkipWaitingCommand(params *SkipWaitingParams) *SkipWaitingCommand {
	return &SkipWaitingCommand{
		params: params,
	}
}

func (cmd *SkipWaitingCommand) Name() string {
	return "ServiceWorker.skipWaiting"
}

func (cmd *SkipWaitingCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SkipWaitingCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SkipWaiting(params *SkipWaitingParams, conn *hc.Conn) (err error) {
	cmd := NewSkipWaitingCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SkipWaitingCB func(err error)

type AsyncSkipWaitingCommand struct {
	params *SkipWaitingParams
	cb     SkipWaitingCB
}

func NewAsyncSkipWaitingCommand(params *SkipWaitingParams, cb SkipWaitingCB) *AsyncSkipWaitingCommand {
	return &AsyncSkipWaitingCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSkipWaitingCommand) Name() string {
	return "ServiceWorker.skipWaiting"
}

func (cmd *AsyncSkipWaitingCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SkipWaitingCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSkipWaitingCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type StopWorkerParams struct {
	VersionId string `json:"versionId"`
}

type StopWorkerCommand struct {
	params *StopWorkerParams
	wg     sync.WaitGroup
	err    error
}

func NewStopWorkerCommand(params *StopWorkerParams) *StopWorkerCommand {
	return &StopWorkerCommand{
		params: params,
	}
}

func (cmd *StopWorkerCommand) Name() string {
	return "ServiceWorker.stopWorker"
}

func (cmd *StopWorkerCommand) Params() interface{} {
	return cmd.params
}

func (cmd *StopWorkerCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func StopWorker(params *StopWorkerParams, conn *hc.Conn) (err error) {
	cmd := NewStopWorkerCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type StopWorkerCB func(err error)

type AsyncStopWorkerCommand struct {
	params *StopWorkerParams
	cb     StopWorkerCB
}

func NewAsyncStopWorkerCommand(params *StopWorkerParams, cb StopWorkerCB) *AsyncStopWorkerCommand {
	return &AsyncStopWorkerCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncStopWorkerCommand) Name() string {
	return "ServiceWorker.stopWorker"
}

func (cmd *AsyncStopWorkerCommand) Params() interface{} {
	return cmd.params
}

func (cmd *StopWorkerCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncStopWorkerCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type InspectWorkerParams struct {
	VersionId string `json:"versionId"`
}

type InspectWorkerCommand struct {
	params *InspectWorkerParams
	wg     sync.WaitGroup
	err    error
}

func NewInspectWorkerCommand(params *InspectWorkerParams) *InspectWorkerCommand {
	return &InspectWorkerCommand{
		params: params,
	}
}

func (cmd *InspectWorkerCommand) Name() string {
	return "ServiceWorker.inspectWorker"
}

func (cmd *InspectWorkerCommand) Params() interface{} {
	return cmd.params
}

func (cmd *InspectWorkerCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func InspectWorker(params *InspectWorkerParams, conn *hc.Conn) (err error) {
	cmd := NewInspectWorkerCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type InspectWorkerCB func(err error)

type AsyncInspectWorkerCommand struct {
	params *InspectWorkerParams
	cb     InspectWorkerCB
}

func NewAsyncInspectWorkerCommand(params *InspectWorkerParams, cb InspectWorkerCB) *AsyncInspectWorkerCommand {
	return &AsyncInspectWorkerCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncInspectWorkerCommand) Name() string {
	return "ServiceWorker.inspectWorker"
}

func (cmd *AsyncInspectWorkerCommand) Params() interface{} {
	return cmd.params
}

func (cmd *InspectWorkerCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncInspectWorkerCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetForceUpdateOnPageLoadParams struct {
	ForceUpdateOnPageLoad bool `json:"forceUpdateOnPageLoad"`
}

type SetForceUpdateOnPageLoadCommand struct {
	params *SetForceUpdateOnPageLoadParams
	wg     sync.WaitGroup
	err    error
}

func NewSetForceUpdateOnPageLoadCommand(params *SetForceUpdateOnPageLoadParams) *SetForceUpdateOnPageLoadCommand {
	return &SetForceUpdateOnPageLoadCommand{
		params: params,
	}
}

func (cmd *SetForceUpdateOnPageLoadCommand) Name() string {
	return "ServiceWorker.setForceUpdateOnPageLoad"
}

func (cmd *SetForceUpdateOnPageLoadCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetForceUpdateOnPageLoadCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetForceUpdateOnPageLoad(params *SetForceUpdateOnPageLoadParams, conn *hc.Conn) (err error) {
	cmd := NewSetForceUpdateOnPageLoadCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetForceUpdateOnPageLoadCB func(err error)

type AsyncSetForceUpdateOnPageLoadCommand struct {
	params *SetForceUpdateOnPageLoadParams
	cb     SetForceUpdateOnPageLoadCB
}

func NewAsyncSetForceUpdateOnPageLoadCommand(params *SetForceUpdateOnPageLoadParams, cb SetForceUpdateOnPageLoadCB) *AsyncSetForceUpdateOnPageLoadCommand {
	return &AsyncSetForceUpdateOnPageLoadCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetForceUpdateOnPageLoadCommand) Name() string {
	return "ServiceWorker.setForceUpdateOnPageLoad"
}

func (cmd *AsyncSetForceUpdateOnPageLoadCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetForceUpdateOnPageLoadCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetForceUpdateOnPageLoadCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type DeliverPushMessageParams struct {
	Origin         string `json:"origin"`
	RegistrationId string `json:"registrationId"`
	Data           string `json:"data"`
}

type DeliverPushMessageCommand struct {
	params *DeliverPushMessageParams
	wg     sync.WaitGroup
	err    error
}

func NewDeliverPushMessageCommand(params *DeliverPushMessageParams) *DeliverPushMessageCommand {
	return &DeliverPushMessageCommand{
		params: params,
	}
}

func (cmd *DeliverPushMessageCommand) Name() string {
	return "ServiceWorker.deliverPushMessage"
}

func (cmd *DeliverPushMessageCommand) Params() interface{} {
	return cmd.params
}

func (cmd *DeliverPushMessageCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func DeliverPushMessage(params *DeliverPushMessageParams, conn *hc.Conn) (err error) {
	cmd := NewDeliverPushMessageCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type DeliverPushMessageCB func(err error)

type AsyncDeliverPushMessageCommand struct {
	params *DeliverPushMessageParams
	cb     DeliverPushMessageCB
}

func NewAsyncDeliverPushMessageCommand(params *DeliverPushMessageParams, cb DeliverPushMessageCB) *AsyncDeliverPushMessageCommand {
	return &AsyncDeliverPushMessageCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncDeliverPushMessageCommand) Name() string {
	return "ServiceWorker.deliverPushMessage"
}

func (cmd *AsyncDeliverPushMessageCommand) Params() interface{} {
	return cmd.params
}

func (cmd *DeliverPushMessageCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncDeliverPushMessageCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type DispatchSyncEventParams struct {
	Origin         string `json:"origin"`
	RegistrationId string `json:"registrationId"`
	Tag            string `json:"tag"`
	LastChance     bool   `json:"lastChance"`
}

type DispatchSyncEventCommand struct {
	params *DispatchSyncEventParams
	wg     sync.WaitGroup
	err    error
}

func NewDispatchSyncEventCommand(params *DispatchSyncEventParams) *DispatchSyncEventCommand {
	return &DispatchSyncEventCommand{
		params: params,
	}
}

func (cmd *DispatchSyncEventCommand) Name() string {
	return "ServiceWorker.dispatchSyncEvent"
}

func (cmd *DispatchSyncEventCommand) Params() interface{} {
	return cmd.params
}

func (cmd *DispatchSyncEventCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func DispatchSyncEvent(params *DispatchSyncEventParams, conn *hc.Conn) (err error) {
	cmd := NewDispatchSyncEventCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type DispatchSyncEventCB func(err error)

type AsyncDispatchSyncEventCommand struct {
	params *DispatchSyncEventParams
	cb     DispatchSyncEventCB
}

func NewAsyncDispatchSyncEventCommand(params *DispatchSyncEventParams, cb DispatchSyncEventCB) *AsyncDispatchSyncEventCommand {
	return &AsyncDispatchSyncEventCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncDispatchSyncEventCommand) Name() string {
	return "ServiceWorker.dispatchSyncEvent"
}

func (cmd *AsyncDispatchSyncEventCommand) Params() interface{} {
	return cmd.params
}

func (cmd *DispatchSyncEventCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncDispatchSyncEventCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type WorkerRegistrationUpdatedEvent struct {
	Registrations []*ServiceWorkerRegistration `json:"registrations"`
}

func OnWorkerRegistrationUpdated(conn *hc.Conn, cb func(evt *WorkerRegistrationUpdatedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &WorkerRegistrationUpdatedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("ServiceWorker.workerRegistrationUpdated", sink)
}

type WorkerVersionUpdatedEvent struct {
	Versions []*ServiceWorkerVersion `json:"versions"`
}

func OnWorkerVersionUpdated(conn *hc.Conn, cb func(evt *WorkerVersionUpdatedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &WorkerVersionUpdatedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("ServiceWorker.workerVersionUpdated", sink)
}

type WorkerErrorReportedEvent struct {
	ErrorMessage *ServiceWorkerErrorMessage `json:"errorMessage"`
}

func OnWorkerErrorReported(conn *hc.Conn, cb func(evt *WorkerErrorReportedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &WorkerErrorReportedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("ServiceWorker.workerErrorReported", sink)
}
