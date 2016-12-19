package protocol

type SetShowPaintRectsParams struct {
	Result bool `json:"result"` // True for showing paint rectangles
}

type SetShowPaintRectsCB func(err error)

// Requests that backend shows paint rectangles
type SetShowPaintRectsCommand struct {
	params *SetShowPaintRectsParams
	cb     SetShowPaintRectsCB
}

func NewSetShowPaintRectsCommand(params *SetShowPaintRectsParams, cb SetShowPaintRectsCB) *SetShowPaintRectsCommand {
	return &SetShowPaintRectsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetShowPaintRectsCommand) Name() string {
	return "Rendering.setShowPaintRects"
}

func (cmd *SetShowPaintRectsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetShowPaintRectsCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetShowDebugBordersParams struct {
	Show bool `json:"show"` // True for showing debug borders
}

type SetShowDebugBordersCB func(err error)

// Requests that backend shows debug borders on layers
type SetShowDebugBordersCommand struct {
	params *SetShowDebugBordersParams
	cb     SetShowDebugBordersCB
}

func NewSetShowDebugBordersCommand(params *SetShowDebugBordersParams, cb SetShowDebugBordersCB) *SetShowDebugBordersCommand {
	return &SetShowDebugBordersCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetShowDebugBordersCommand) Name() string {
	return "Rendering.setShowDebugBorders"
}

func (cmd *SetShowDebugBordersCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetShowDebugBordersCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetShowFPSCounterParams struct {
	Show bool `json:"show"` // True for showing the FPS counter
}

type SetShowFPSCounterCB func(err error)

// Requests that backend shows the FPS counter
type SetShowFPSCounterCommand struct {
	params *SetShowFPSCounterParams
	cb     SetShowFPSCounterCB
}

func NewSetShowFPSCounterCommand(params *SetShowFPSCounterParams, cb SetShowFPSCounterCB) *SetShowFPSCounterCommand {
	return &SetShowFPSCounterCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetShowFPSCounterCommand) Name() string {
	return "Rendering.setShowFPSCounter"
}

func (cmd *SetShowFPSCounterCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetShowFPSCounterCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetShowScrollBottleneckRectsParams struct {
	Show bool `json:"show"` // True for showing scroll bottleneck rects
}

type SetShowScrollBottleneckRectsCB func(err error)

// Requests that backend shows scroll bottleneck rects
type SetShowScrollBottleneckRectsCommand struct {
	params *SetShowScrollBottleneckRectsParams
	cb     SetShowScrollBottleneckRectsCB
}

func NewSetShowScrollBottleneckRectsCommand(params *SetShowScrollBottleneckRectsParams, cb SetShowScrollBottleneckRectsCB) *SetShowScrollBottleneckRectsCommand {
	return &SetShowScrollBottleneckRectsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetShowScrollBottleneckRectsCommand) Name() string {
	return "Rendering.setShowScrollBottleneckRects"
}

func (cmd *SetShowScrollBottleneckRectsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetShowScrollBottleneckRectsCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetShowViewportSizeOnResizeParams struct {
	Show bool `json:"show"` // Whether to paint size or not.
}

type SetShowViewportSizeOnResizeCB func(err error)

// Paints viewport size upon main frame resize.
type SetShowViewportSizeOnResizeCommand struct {
	params *SetShowViewportSizeOnResizeParams
	cb     SetShowViewportSizeOnResizeCB
}

func NewSetShowViewportSizeOnResizeCommand(params *SetShowViewportSizeOnResizeParams, cb SetShowViewportSizeOnResizeCB) *SetShowViewportSizeOnResizeCommand {
	return &SetShowViewportSizeOnResizeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetShowViewportSizeOnResizeCommand) Name() string {
	return "Rendering.setShowViewportSizeOnResize"
}

func (cmd *SetShowViewportSizeOnResizeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetShowViewportSizeOnResizeCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}
