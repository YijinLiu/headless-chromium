package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
)

// DOM Storage identifier.
type StorageId struct {
	SecurityOrigin string `json:"securityOrigin"` // Security origin for the storage.
	IsLocalStorage bool   `json:"isLocalStorage"` // Whether the storage is local storage (not session storage).
}

// DOM Storage item.
type Item []string

type DOMStorageEnableCB func(err error)

// Enables storage tracking, storage events will now be delivered to the client.
type DOMStorageEnableCommand struct {
	cb DOMStorageEnableCB
}

func NewDOMStorageEnableCommand(cb DOMStorageEnableCB) *DOMStorageEnableCommand {
	return &DOMStorageEnableCommand{
		cb: cb,
	}
}

func (cmd *DOMStorageEnableCommand) Name() string {
	return "DOMStorage.enable"
}

func (cmd *DOMStorageEnableCommand) Params() interface{} {
	return nil
}

func (cmd *DOMStorageEnableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type DOMStorageDisableCB func(err error)

// Disables storage tracking, prevents storage events from being sent to the client.
type DOMStorageDisableCommand struct {
	cb DOMStorageDisableCB
}

func NewDOMStorageDisableCommand(cb DOMStorageDisableCB) *DOMStorageDisableCommand {
	return &DOMStorageDisableCommand{
		cb: cb,
	}
}

func (cmd *DOMStorageDisableCommand) Name() string {
	return "DOMStorage.disable"
}

func (cmd *DOMStorageDisableCommand) Params() interface{} {
	return nil
}

func (cmd *DOMStorageDisableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type GetDOMStorageItemsParams struct {
	StorageId *StorageId `json:"storageId"`
}

type GetDOMStorageItemsResult struct {
	Entries []Item `json:"entries"`
}

type GetDOMStorageItemsCB func(result *GetDOMStorageItemsResult, err error)

type GetDOMStorageItemsCommand struct {
	params *GetDOMStorageItemsParams
	cb     GetDOMStorageItemsCB
}

func NewGetDOMStorageItemsCommand(params *GetDOMStorageItemsParams, cb GetDOMStorageItemsCB) *GetDOMStorageItemsCommand {
	return &GetDOMStorageItemsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetDOMStorageItemsCommand) Name() string {
	return "DOMStorage.getDOMStorageItems"
}

func (cmd *GetDOMStorageItemsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetDOMStorageItemsCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetDOMStorageItemsResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type SetDOMStorageItemParams struct {
	StorageId *StorageId `json:"storageId"`
	Key       string     `json:"key"`
	Value     string     `json:"value"`
}

type SetDOMStorageItemCB func(err error)

type SetDOMStorageItemCommand struct {
	params *SetDOMStorageItemParams
	cb     SetDOMStorageItemCB
}

func NewSetDOMStorageItemCommand(params *SetDOMStorageItemParams, cb SetDOMStorageItemCB) *SetDOMStorageItemCommand {
	return &SetDOMStorageItemCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetDOMStorageItemCommand) Name() string {
	return "DOMStorage.setDOMStorageItem"
}

func (cmd *SetDOMStorageItemCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetDOMStorageItemCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type RemoveDOMStorageItemParams struct {
	StorageId *StorageId `json:"storageId"`
	Key       string     `json:"key"`
}

type RemoveDOMStorageItemCB func(err error)

type RemoveDOMStorageItemCommand struct {
	params *RemoveDOMStorageItemParams
	cb     RemoveDOMStorageItemCB
}

func NewRemoveDOMStorageItemCommand(params *RemoveDOMStorageItemParams, cb RemoveDOMStorageItemCB) *RemoveDOMStorageItemCommand {
	return &RemoveDOMStorageItemCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *RemoveDOMStorageItemCommand) Name() string {
	return "DOMStorage.removeDOMStorageItem"
}

func (cmd *RemoveDOMStorageItemCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveDOMStorageItemCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type DomStorageItemsClearedEvent struct {
	StorageId *StorageId `json:"storageId"`
}

type DomStorageItemsClearedEventSink struct {
	events chan *DomStorageItemsClearedEvent
}

func NewDomStorageItemsClearedEventSink(bufSize int) *DomStorageItemsClearedEventSink {
	return &DomStorageItemsClearedEventSink{
		events: make(chan *DomStorageItemsClearedEvent, bufSize),
	}
}

func (s *DomStorageItemsClearedEventSink) Name() string {
	return "DOMStorage.domStorageItemsCleared"
}

func (s *DomStorageItemsClearedEventSink) OnEvent(params []byte) {
	evt := &DomStorageItemsClearedEvent{}
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

type DomStorageItemRemovedEvent struct {
	StorageId *StorageId `json:"storageId"`
	Key       string     `json:"key"`
}

type DomStorageItemRemovedEventSink struct {
	events chan *DomStorageItemRemovedEvent
}

func NewDomStorageItemRemovedEventSink(bufSize int) *DomStorageItemRemovedEventSink {
	return &DomStorageItemRemovedEventSink{
		events: make(chan *DomStorageItemRemovedEvent, bufSize),
	}
}

func (s *DomStorageItemRemovedEventSink) Name() string {
	return "DOMStorage.domStorageItemRemoved"
}

func (s *DomStorageItemRemovedEventSink) OnEvent(params []byte) {
	evt := &DomStorageItemRemovedEvent{}
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

type DomStorageItemAddedEvent struct {
	StorageId *StorageId `json:"storageId"`
	Key       string     `json:"key"`
	NewValue  string     `json:"newValue"`
}

type DomStorageItemAddedEventSink struct {
	events chan *DomStorageItemAddedEvent
}

func NewDomStorageItemAddedEventSink(bufSize int) *DomStorageItemAddedEventSink {
	return &DomStorageItemAddedEventSink{
		events: make(chan *DomStorageItemAddedEvent, bufSize),
	}
}

func (s *DomStorageItemAddedEventSink) Name() string {
	return "DOMStorage.domStorageItemAdded"
}

func (s *DomStorageItemAddedEventSink) OnEvent(params []byte) {
	evt := &DomStorageItemAddedEvent{}
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

type DomStorageItemUpdatedEvent struct {
	StorageId *StorageId `json:"storageId"`
	Key       string     `json:"key"`
	OldValue  string     `json:"oldValue"`
	NewValue  string     `json:"newValue"`
}

type DomStorageItemUpdatedEventSink struct {
	events chan *DomStorageItemUpdatedEvent
}

func NewDomStorageItemUpdatedEventSink(bufSize int) *DomStorageItemUpdatedEventSink {
	return &DomStorageItemUpdatedEventSink{
		events: make(chan *DomStorageItemUpdatedEvent, bufSize),
	}
}

func (s *DomStorageItemUpdatedEventSink) Name() string {
	return "DOMStorage.domStorageItemUpdated"
}

func (s *DomStorageItemUpdatedEventSink) OnEvent(params []byte) {
	evt := &DomStorageItemUpdatedEvent{}
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
