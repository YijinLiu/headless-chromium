package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

// DOM Storage identifier.
// @experimental
type StorageId struct {
	SecurityOrigin string `json:"securityOrigin"` // Security origin for the storage.
	IsLocalStorage bool   `json:"isLocalStorage"` // Whether the storage is local storage (not session storage).
}

// DOM Storage item.
// @experimental
type Item []string

// Enables storage tracking, storage events will now be delivered to the client.

type DOMStorageEnableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewDOMStorageEnableCommand() *DOMStorageEnableCommand {
	return &DOMStorageEnableCommand{}
}

func (cmd *DOMStorageEnableCommand) Name() string {
	return "DOMStorage.enable"
}

func (cmd *DOMStorageEnableCommand) Params() interface{} {
	return nil
}

func (cmd *DOMStorageEnableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func DOMStorageEnable(conn *hc.Conn) (err error) {
	cmd := NewDOMStorageEnableCommand()
	cmd.Run(conn)
	return cmd.err
}

type DOMStorageEnableCB func(err error)

// Enables storage tracking, storage events will now be delivered to the client.

type AsyncDOMStorageEnableCommand struct {
	cb DOMStorageEnableCB
}

func NewAsyncDOMStorageEnableCommand(cb DOMStorageEnableCB) *AsyncDOMStorageEnableCommand {
	return &AsyncDOMStorageEnableCommand{
		cb: cb,
	}
}

func (cmd *AsyncDOMStorageEnableCommand) Name() string {
	return "DOMStorage.enable"
}

func (cmd *AsyncDOMStorageEnableCommand) Params() interface{} {
	return nil
}

func (cmd *DOMStorageEnableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncDOMStorageEnableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Disables storage tracking, prevents storage events from being sent to the client.

type DOMStorageDisableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewDOMStorageDisableCommand() *DOMStorageDisableCommand {
	return &DOMStorageDisableCommand{}
}

func (cmd *DOMStorageDisableCommand) Name() string {
	return "DOMStorage.disable"
}

func (cmd *DOMStorageDisableCommand) Params() interface{} {
	return nil
}

func (cmd *DOMStorageDisableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func DOMStorageDisable(conn *hc.Conn) (err error) {
	cmd := NewDOMStorageDisableCommand()
	cmd.Run(conn)
	return cmd.err
}

type DOMStorageDisableCB func(err error)

// Disables storage tracking, prevents storage events from being sent to the client.

type AsyncDOMStorageDisableCommand struct {
	cb DOMStorageDisableCB
}

func NewAsyncDOMStorageDisableCommand(cb DOMStorageDisableCB) *AsyncDOMStorageDisableCommand {
	return &AsyncDOMStorageDisableCommand{
		cb: cb,
	}
}

func (cmd *AsyncDOMStorageDisableCommand) Name() string {
	return "DOMStorage.disable"
}

func (cmd *AsyncDOMStorageDisableCommand) Params() interface{} {
	return nil
}

func (cmd *DOMStorageDisableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncDOMStorageDisableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type GetDOMStorageItemsParams struct {
	StorageId *StorageId `json:"storageId"`
}

type GetDOMStorageItemsResult struct {
	Entries []Item `json:"entries"`
}

type GetDOMStorageItemsCommand struct {
	params *GetDOMStorageItemsParams
	result GetDOMStorageItemsResult
	wg     sync.WaitGroup
	err    error
}

func NewGetDOMStorageItemsCommand(params *GetDOMStorageItemsParams) *GetDOMStorageItemsCommand {
	return &GetDOMStorageItemsCommand{
		params: params,
	}
}

func (cmd *GetDOMStorageItemsCommand) Name() string {
	return "DOMStorage.getDOMStorageItems"
}

func (cmd *GetDOMStorageItemsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetDOMStorageItemsCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetDOMStorageItems(params *GetDOMStorageItemsParams, conn *hc.Conn) (result *GetDOMStorageItemsResult, err error) {
	cmd := NewGetDOMStorageItemsCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetDOMStorageItemsCB func(result *GetDOMStorageItemsResult, err error)

type AsyncGetDOMStorageItemsCommand struct {
	params *GetDOMStorageItemsParams
	cb     GetDOMStorageItemsCB
}

func NewAsyncGetDOMStorageItemsCommand(params *GetDOMStorageItemsParams, cb GetDOMStorageItemsCB) *AsyncGetDOMStorageItemsCommand {
	return &AsyncGetDOMStorageItemsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetDOMStorageItemsCommand) Name() string {
	return "DOMStorage.getDOMStorageItems"
}

func (cmd *AsyncGetDOMStorageItemsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetDOMStorageItemsCommand) Result() *GetDOMStorageItemsResult {
	return &cmd.result
}

func (cmd *GetDOMStorageItemsCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetDOMStorageItemsCommand) Done(data []byte, err error) {
	var result GetDOMStorageItemsResult
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

type SetDOMStorageItemParams struct {
	StorageId *StorageId `json:"storageId"`
	Key       string     `json:"key"`
	Value     string     `json:"value"`
}

type SetDOMStorageItemCommand struct {
	params *SetDOMStorageItemParams
	wg     sync.WaitGroup
	err    error
}

func NewSetDOMStorageItemCommand(params *SetDOMStorageItemParams) *SetDOMStorageItemCommand {
	return &SetDOMStorageItemCommand{
		params: params,
	}
}

func (cmd *SetDOMStorageItemCommand) Name() string {
	return "DOMStorage.setDOMStorageItem"
}

func (cmd *SetDOMStorageItemCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetDOMStorageItemCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetDOMStorageItem(params *SetDOMStorageItemParams, conn *hc.Conn) (err error) {
	cmd := NewSetDOMStorageItemCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetDOMStorageItemCB func(err error)

type AsyncSetDOMStorageItemCommand struct {
	params *SetDOMStorageItemParams
	cb     SetDOMStorageItemCB
}

func NewAsyncSetDOMStorageItemCommand(params *SetDOMStorageItemParams, cb SetDOMStorageItemCB) *AsyncSetDOMStorageItemCommand {
	return &AsyncSetDOMStorageItemCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetDOMStorageItemCommand) Name() string {
	return "DOMStorage.setDOMStorageItem"
}

func (cmd *AsyncSetDOMStorageItemCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetDOMStorageItemCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetDOMStorageItemCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type RemoveDOMStorageItemParams struct {
	StorageId *StorageId `json:"storageId"`
	Key       string     `json:"key"`
}

type RemoveDOMStorageItemCommand struct {
	params *RemoveDOMStorageItemParams
	wg     sync.WaitGroup
	err    error
}

func NewRemoveDOMStorageItemCommand(params *RemoveDOMStorageItemParams) *RemoveDOMStorageItemCommand {
	return &RemoveDOMStorageItemCommand{
		params: params,
	}
}

func (cmd *RemoveDOMStorageItemCommand) Name() string {
	return "DOMStorage.removeDOMStorageItem"
}

func (cmd *RemoveDOMStorageItemCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveDOMStorageItemCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func RemoveDOMStorageItem(params *RemoveDOMStorageItemParams, conn *hc.Conn) (err error) {
	cmd := NewRemoveDOMStorageItemCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type RemoveDOMStorageItemCB func(err error)

type AsyncRemoveDOMStorageItemCommand struct {
	params *RemoveDOMStorageItemParams
	cb     RemoveDOMStorageItemCB
}

func NewAsyncRemoveDOMStorageItemCommand(params *RemoveDOMStorageItemParams, cb RemoveDOMStorageItemCB) *AsyncRemoveDOMStorageItemCommand {
	return &AsyncRemoveDOMStorageItemCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncRemoveDOMStorageItemCommand) Name() string {
	return "DOMStorage.removeDOMStorageItem"
}

func (cmd *AsyncRemoveDOMStorageItemCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveDOMStorageItemCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncRemoveDOMStorageItemCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type DomStorageItemsClearedEvent struct {
	StorageId *StorageId `json:"storageId"`
}

func OnDomStorageItemsCleared(conn *hc.Conn, cb func(evt *DomStorageItemsClearedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &DomStorageItemsClearedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("DOMStorage.domStorageItemsCleared", sink)
}

type DomStorageItemRemovedEvent struct {
	StorageId *StorageId `json:"storageId"`
	Key       string     `json:"key"`
}

func OnDomStorageItemRemoved(conn *hc.Conn, cb func(evt *DomStorageItemRemovedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &DomStorageItemRemovedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("DOMStorage.domStorageItemRemoved", sink)
}

type DomStorageItemAddedEvent struct {
	StorageId *StorageId `json:"storageId"`
	Key       string     `json:"key"`
	NewValue  string     `json:"newValue"`
}

func OnDomStorageItemAdded(conn *hc.Conn, cb func(evt *DomStorageItemAddedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &DomStorageItemAddedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("DOMStorage.domStorageItemAdded", sink)
}

type DomStorageItemUpdatedEvent struct {
	StorageId *StorageId `json:"storageId"`
	Key       string     `json:"key"`
	OldValue  string     `json:"oldValue"`
	NewValue  string     `json:"newValue"`
}

func OnDomStorageItemUpdated(conn *hc.Conn, cb func(evt *DomStorageItemUpdatedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &DomStorageItemUpdatedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("DOMStorage.domStorageItemUpdated", sink)
}
