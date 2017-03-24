package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

// Unique identifier of the Cache object.
type CacheId string

// Data entry.
type CacheStorageDataEntry struct {
	Request  string `json:"request"`  // Request url spec.
	Response string `json:"response"` // Response stataus text.
}

// Cache identifier.
type Cache struct {
	CacheId        CacheId `json:"cacheId"`        // An opaque unique id of the cache.
	SecurityOrigin string  `json:"securityOrigin"` // Security origin of the cache.
	CacheName      string  `json:"cacheName"`      // The name of the cache.
}

type RequestCacheNamesParams struct {
	SecurityOrigin string `json:"securityOrigin"` // Security origin.
}

type RequestCacheNamesResult struct {
	Caches []*Cache `json:"caches"` // Caches for the security origin.
}

// Requests cache names.

type RequestCacheNamesCommand struct {
	params *RequestCacheNamesParams
	result RequestCacheNamesResult
	wg     sync.WaitGroup
	err    error
}

func NewRequestCacheNamesCommand(params *RequestCacheNamesParams) *RequestCacheNamesCommand {
	return &RequestCacheNamesCommand{
		params: params,
	}
}

func (cmd *RequestCacheNamesCommand) Name() string {
	return "CacheStorage.requestCacheNames"
}

func (cmd *RequestCacheNamesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RequestCacheNamesCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func RequestCacheNames(params *RequestCacheNamesParams, conn *hc.Conn) (result *RequestCacheNamesResult, err error) {
	cmd := NewRequestCacheNamesCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type RequestCacheNamesCB func(result *RequestCacheNamesResult, err error)

// Requests cache names.

type AsyncRequestCacheNamesCommand struct {
	params *RequestCacheNamesParams
	cb     RequestCacheNamesCB
}

func NewAsyncRequestCacheNamesCommand(params *RequestCacheNamesParams, cb RequestCacheNamesCB) *AsyncRequestCacheNamesCommand {
	return &AsyncRequestCacheNamesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncRequestCacheNamesCommand) Name() string {
	return "CacheStorage.requestCacheNames"
}

func (cmd *AsyncRequestCacheNamesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RequestCacheNamesCommand) Result() *RequestCacheNamesResult {
	return &cmd.result
}

func (cmd *RequestCacheNamesCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncRequestCacheNamesCommand) Done(data []byte, err error) {
	var result RequestCacheNamesResult
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

type RequestEntriesParams struct {
	CacheId   CacheId `json:"cacheId"`   // ID of cache to get entries from.
	SkipCount int     `json:"skipCount"` // Number of records to skip.
	PageSize  int     `json:"pageSize"`  // Number of records to fetch.
}

type RequestEntriesResult struct {
	CacheDataEntries []*CacheStorageDataEntry `json:"cacheDataEntries"` // Array of object store data entries.
	HasMore          bool                     `json:"hasMore"`          // If true, there are more entries to fetch in the given range.
}

// Requests data from cache.

type RequestEntriesCommand struct {
	params *RequestEntriesParams
	result RequestEntriesResult
	wg     sync.WaitGroup
	err    error
}

func NewRequestEntriesCommand(params *RequestEntriesParams) *RequestEntriesCommand {
	return &RequestEntriesCommand{
		params: params,
	}
}

func (cmd *RequestEntriesCommand) Name() string {
	return "CacheStorage.requestEntries"
}

func (cmd *RequestEntriesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RequestEntriesCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func RequestEntries(params *RequestEntriesParams, conn *hc.Conn) (result *RequestEntriesResult, err error) {
	cmd := NewRequestEntriesCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type RequestEntriesCB func(result *RequestEntriesResult, err error)

// Requests data from cache.

type AsyncRequestEntriesCommand struct {
	params *RequestEntriesParams
	cb     RequestEntriesCB
}

func NewAsyncRequestEntriesCommand(params *RequestEntriesParams, cb RequestEntriesCB) *AsyncRequestEntriesCommand {
	return &AsyncRequestEntriesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncRequestEntriesCommand) Name() string {
	return "CacheStorage.requestEntries"
}

func (cmd *AsyncRequestEntriesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RequestEntriesCommand) Result() *RequestEntriesResult {
	return &cmd.result
}

func (cmd *RequestEntriesCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncRequestEntriesCommand) Done(data []byte, err error) {
	var result RequestEntriesResult
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

type DeleteCacheParams struct {
	CacheId CacheId `json:"cacheId"` // Id of cache for deletion.
}

// Deletes a cache.

type DeleteCacheCommand struct {
	params *DeleteCacheParams
	wg     sync.WaitGroup
	err    error
}

func NewDeleteCacheCommand(params *DeleteCacheParams) *DeleteCacheCommand {
	return &DeleteCacheCommand{
		params: params,
	}
}

func (cmd *DeleteCacheCommand) Name() string {
	return "CacheStorage.deleteCache"
}

func (cmd *DeleteCacheCommand) Params() interface{} {
	return cmd.params
}

func (cmd *DeleteCacheCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func DeleteCache(params *DeleteCacheParams, conn *hc.Conn) (err error) {
	cmd := NewDeleteCacheCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type DeleteCacheCB func(err error)

// Deletes a cache.

type AsyncDeleteCacheCommand struct {
	params *DeleteCacheParams
	cb     DeleteCacheCB
}

func NewAsyncDeleteCacheCommand(params *DeleteCacheParams, cb DeleteCacheCB) *AsyncDeleteCacheCommand {
	return &AsyncDeleteCacheCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncDeleteCacheCommand) Name() string {
	return "CacheStorage.deleteCache"
}

func (cmd *AsyncDeleteCacheCommand) Params() interface{} {
	return cmd.params
}

func (cmd *DeleteCacheCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncDeleteCacheCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type DeleteEntryParams struct {
	CacheId CacheId `json:"cacheId"` // Id of cache where the entry will be deleted.
	Request string  `json:"request"` // URL spec of the request.
}

// Deletes a cache entry.

type DeleteEntryCommand struct {
	params *DeleteEntryParams
	wg     sync.WaitGroup
	err    error
}

func NewDeleteEntryCommand(params *DeleteEntryParams) *DeleteEntryCommand {
	return &DeleteEntryCommand{
		params: params,
	}
}

func (cmd *DeleteEntryCommand) Name() string {
	return "CacheStorage.deleteEntry"
}

func (cmd *DeleteEntryCommand) Params() interface{} {
	return cmd.params
}

func (cmd *DeleteEntryCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func DeleteEntry(params *DeleteEntryParams, conn *hc.Conn) (err error) {
	cmd := NewDeleteEntryCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type DeleteEntryCB func(err error)

// Deletes a cache entry.

type AsyncDeleteEntryCommand struct {
	params *DeleteEntryParams
	cb     DeleteEntryCB
}

func NewAsyncDeleteEntryCommand(params *DeleteEntryParams, cb DeleteEntryCB) *AsyncDeleteEntryCommand {
	return &AsyncDeleteEntryCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncDeleteEntryCommand) Name() string {
	return "CacheStorage.deleteEntry"
}

func (cmd *AsyncDeleteEntryCommand) Params() interface{} {
	return cmd.params
}

func (cmd *DeleteEntryCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncDeleteEntryCommand) Done(data []byte, err error) {
	cmd.cb(err)
}
