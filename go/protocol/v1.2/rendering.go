package protocol

import (
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

type SetShowPaintRectsParams struct {
	Result bool `json:"result"` // True for showing paint rectangles
}

// Requests that backend shows paint rectangles

type SetShowPaintRectsCommand struct {
	params *SetShowPaintRectsParams
	wg     sync.WaitGroup
	err    error
}

func NewSetShowPaintRectsCommand(params *SetShowPaintRectsParams) *SetShowPaintRectsCommand {
	return &SetShowPaintRectsCommand{
		params: params,
	}
}

func (cmd *SetShowPaintRectsCommand) Name() string {
	return "Rendering.setShowPaintRects"
}

func (cmd *SetShowPaintRectsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetShowPaintRectsCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetShowPaintRects(params *SetShowPaintRectsParams, conn *hc.Conn) (err error) {
	cmd := NewSetShowPaintRectsCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetShowPaintRectsCB func(err error)

// Requests that backend shows paint rectangles

type AsyncSetShowPaintRectsCommand struct {
	params *SetShowPaintRectsParams
	cb     SetShowPaintRectsCB
}

func NewAsyncSetShowPaintRectsCommand(params *SetShowPaintRectsParams, cb SetShowPaintRectsCB) *AsyncSetShowPaintRectsCommand {
	return &AsyncSetShowPaintRectsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetShowPaintRectsCommand) Name() string {
	return "Rendering.setShowPaintRects"
}

func (cmd *AsyncSetShowPaintRectsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetShowPaintRectsCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetShowPaintRectsCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetShowDebugBordersParams struct {
	Show bool `json:"show"` // True for showing debug borders
}

// Requests that backend shows debug borders on layers

type SetShowDebugBordersCommand struct {
	params *SetShowDebugBordersParams
	wg     sync.WaitGroup
	err    error
}

func NewSetShowDebugBordersCommand(params *SetShowDebugBordersParams) *SetShowDebugBordersCommand {
	return &SetShowDebugBordersCommand{
		params: params,
	}
}

func (cmd *SetShowDebugBordersCommand) Name() string {
	return "Rendering.setShowDebugBorders"
}

func (cmd *SetShowDebugBordersCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetShowDebugBordersCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetShowDebugBorders(params *SetShowDebugBordersParams, conn *hc.Conn) (err error) {
	cmd := NewSetShowDebugBordersCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetShowDebugBordersCB func(err error)

// Requests that backend shows debug borders on layers

type AsyncSetShowDebugBordersCommand struct {
	params *SetShowDebugBordersParams
	cb     SetShowDebugBordersCB
}

func NewAsyncSetShowDebugBordersCommand(params *SetShowDebugBordersParams, cb SetShowDebugBordersCB) *AsyncSetShowDebugBordersCommand {
	return &AsyncSetShowDebugBordersCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetShowDebugBordersCommand) Name() string {
	return "Rendering.setShowDebugBorders"
}

func (cmd *AsyncSetShowDebugBordersCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetShowDebugBordersCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetShowDebugBordersCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetShowFPSCounterParams struct {
	Show bool `json:"show"` // True for showing the FPS counter
}

// Requests that backend shows the FPS counter

type SetShowFPSCounterCommand struct {
	params *SetShowFPSCounterParams
	wg     sync.WaitGroup
	err    error
}

func NewSetShowFPSCounterCommand(params *SetShowFPSCounterParams) *SetShowFPSCounterCommand {
	return &SetShowFPSCounterCommand{
		params: params,
	}
}

func (cmd *SetShowFPSCounterCommand) Name() string {
	return "Rendering.setShowFPSCounter"
}

func (cmd *SetShowFPSCounterCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetShowFPSCounterCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetShowFPSCounter(params *SetShowFPSCounterParams, conn *hc.Conn) (err error) {
	cmd := NewSetShowFPSCounterCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetShowFPSCounterCB func(err error)

// Requests that backend shows the FPS counter

type AsyncSetShowFPSCounterCommand struct {
	params *SetShowFPSCounterParams
	cb     SetShowFPSCounterCB
}

func NewAsyncSetShowFPSCounterCommand(params *SetShowFPSCounterParams, cb SetShowFPSCounterCB) *AsyncSetShowFPSCounterCommand {
	return &AsyncSetShowFPSCounterCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetShowFPSCounterCommand) Name() string {
	return "Rendering.setShowFPSCounter"
}

func (cmd *AsyncSetShowFPSCounterCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetShowFPSCounterCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetShowFPSCounterCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetShowScrollBottleneckRectsParams struct {
	Show bool `json:"show"` // True for showing scroll bottleneck rects
}

// Requests that backend shows scroll bottleneck rects

type SetShowScrollBottleneckRectsCommand struct {
	params *SetShowScrollBottleneckRectsParams
	wg     sync.WaitGroup
	err    error
}

func NewSetShowScrollBottleneckRectsCommand(params *SetShowScrollBottleneckRectsParams) *SetShowScrollBottleneckRectsCommand {
	return &SetShowScrollBottleneckRectsCommand{
		params: params,
	}
}

func (cmd *SetShowScrollBottleneckRectsCommand) Name() string {
	return "Rendering.setShowScrollBottleneckRects"
}

func (cmd *SetShowScrollBottleneckRectsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetShowScrollBottleneckRectsCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetShowScrollBottleneckRects(params *SetShowScrollBottleneckRectsParams, conn *hc.Conn) (err error) {
	cmd := NewSetShowScrollBottleneckRectsCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetShowScrollBottleneckRectsCB func(err error)

// Requests that backend shows scroll bottleneck rects

type AsyncSetShowScrollBottleneckRectsCommand struct {
	params *SetShowScrollBottleneckRectsParams
	cb     SetShowScrollBottleneckRectsCB
}

func NewAsyncSetShowScrollBottleneckRectsCommand(params *SetShowScrollBottleneckRectsParams, cb SetShowScrollBottleneckRectsCB) *AsyncSetShowScrollBottleneckRectsCommand {
	return &AsyncSetShowScrollBottleneckRectsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetShowScrollBottleneckRectsCommand) Name() string {
	return "Rendering.setShowScrollBottleneckRects"
}

func (cmd *AsyncSetShowScrollBottleneckRectsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetShowScrollBottleneckRectsCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetShowScrollBottleneckRectsCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetShowViewportSizeOnResizeParams struct {
	Show bool `json:"show"` // Whether to paint size or not.
}

// Paints viewport size upon main frame resize.

type SetShowViewportSizeOnResizeCommand struct {
	params *SetShowViewportSizeOnResizeParams
	wg     sync.WaitGroup
	err    error
}

func NewSetShowViewportSizeOnResizeCommand(params *SetShowViewportSizeOnResizeParams) *SetShowViewportSizeOnResizeCommand {
	return &SetShowViewportSizeOnResizeCommand{
		params: params,
	}
}

func (cmd *SetShowViewportSizeOnResizeCommand) Name() string {
	return "Rendering.setShowViewportSizeOnResize"
}

func (cmd *SetShowViewportSizeOnResizeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetShowViewportSizeOnResizeCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetShowViewportSizeOnResize(params *SetShowViewportSizeOnResizeParams, conn *hc.Conn) (err error) {
	cmd := NewSetShowViewportSizeOnResizeCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetShowViewportSizeOnResizeCB func(err error)

// Paints viewport size upon main frame resize.

type AsyncSetShowViewportSizeOnResizeCommand struct {
	params *SetShowViewportSizeOnResizeParams
	cb     SetShowViewportSizeOnResizeCB
}

func NewAsyncSetShowViewportSizeOnResizeCommand(params *SetShowViewportSizeOnResizeParams, cb SetShowViewportSizeOnResizeCB) *AsyncSetShowViewportSizeOnResizeCommand {
	return &AsyncSetShowViewportSizeOnResizeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetShowViewportSizeOnResizeCommand) Name() string {
	return "Rendering.setShowViewportSizeOnResize"
}

func (cmd *AsyncSetShowViewportSizeOnResizeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetShowViewportSizeOnResizeCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetShowViewportSizeOnResizeCommand) Done(data []byte, err error) {
	cmd.cb(err)
}
