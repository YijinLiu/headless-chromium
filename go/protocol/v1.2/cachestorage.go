package protocol

import (
	"encoding/json"
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

type RequestCacheNamesCB func(result *RequestCacheNamesResult, err error)

// Requests cache names.
type RequestCacheNamesCommand struct {
	params *RequestCacheNamesParams
	cb     RequestCacheNamesCB
}

func NewRequestCacheNamesCommand(params *RequestCacheNamesParams, cb RequestCacheNamesCB) *RequestCacheNamesCommand {
	return &RequestCacheNamesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *RequestCacheNamesCommand) Name() string {
	return "CacheStorage.requestCacheNames"
}

func (cmd *RequestCacheNamesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RequestCacheNamesCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj RequestCacheNamesResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
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

type RequestEntriesCB func(result *RequestEntriesResult, err error)

// Requests data from cache.
type RequestEntriesCommand struct {
	params *RequestEntriesParams
	cb     RequestEntriesCB
}

func NewRequestEntriesCommand(params *RequestEntriesParams, cb RequestEntriesCB) *RequestEntriesCommand {
	return &RequestEntriesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *RequestEntriesCommand) Name() string {
	return "CacheStorage.requestEntries"
}

func (cmd *RequestEntriesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RequestEntriesCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj RequestEntriesResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type DeleteCacheParams struct {
	CacheId CacheId `json:"cacheId"` // Id of cache for deletion.
}

type DeleteCacheCB func(err error)

// Deletes a cache.
type DeleteCacheCommand struct {
	params *DeleteCacheParams
	cb     DeleteCacheCB
}

func NewDeleteCacheCommand(params *DeleteCacheParams, cb DeleteCacheCB) *DeleteCacheCommand {
	return &DeleteCacheCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *DeleteCacheCommand) Name() string {
	return "CacheStorage.deleteCache"
}

func (cmd *DeleteCacheCommand) Params() interface{} {
	return cmd.params
}

func (cmd *DeleteCacheCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type DeleteEntryParams struct {
	CacheId CacheId `json:"cacheId"` // Id of cache where the entry will be deleted.
	Request string  `json:"request"` // URL spec of the request.
}

type DeleteEntryCB func(err error)

// Deletes a cache entry.
type DeleteEntryCommand struct {
	params *DeleteEntryParams
	cb     DeleteEntryCB
}

func NewDeleteEntryCommand(params *DeleteEntryParams, cb DeleteEntryCB) *DeleteEntryCommand {
	return &DeleteEntryCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *DeleteEntryCommand) Name() string {
	return "CacheStorage.deleteEntry"
}

func (cmd *DeleteEntryCommand) Params() interface{} {
	return cmd.params
}

func (cmd *DeleteEntryCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}
