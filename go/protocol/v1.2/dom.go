package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

// Unique DOM node identifier.
type NodeId int

// Unique DOM node identifier used to reference a node that may not have been pushed to the front-end.
// @experimental
type BackendNodeId int

// Backend node with a friendly name.
// @experimental
type BackendNode struct {
	NodeType      int           `json:"nodeType"` // Node's nodeType.
	NodeName      string        `json:"nodeName"` // Node's nodeName.
	BackendNodeId BackendNodeId `json:"backendNodeId"`
}

// Pseudo element type.
type PseudoType string

const PseudoTypeFirstLine PseudoType = "first-line"
const PseudoTypeFirstLetter PseudoType = "first-letter"
const PseudoTypeBefore PseudoType = "before"
const PseudoTypeAfter PseudoType = "after"
const PseudoTypeBackdrop PseudoType = "backdrop"
const PseudoTypeSelection PseudoType = "selection"
const PseudoTypeFirstLineInherited PseudoType = "first-line-inherited"
const PseudoTypeScrollbar PseudoType = "scrollbar"
const PseudoTypeScrollbarThumb PseudoType = "scrollbar-thumb"
const PseudoTypeScrollbarButton PseudoType = "scrollbar-button"
const PseudoTypeScrollbarTrack PseudoType = "scrollbar-track"
const PseudoTypeScrollbarTrackPiece PseudoType = "scrollbar-track-piece"
const PseudoTypeScrollbarCorner PseudoType = "scrollbar-corner"
const PseudoTypeResizer PseudoType = "resizer"
const PseudoTypeInputListButton PseudoType = "input-list-button"

// Shadow root type.
type ShadowRootType string

const ShadowRootTypeUserAgent ShadowRootType = "user-agent"
const ShadowRootTypeOpen ShadowRootType = "open"
const ShadowRootTypeClosed ShadowRootType = "closed"

// DOM interaction is implemented in terms of mirror objects that represent the actual DOM nodes. DOMNode is a base node mirror type.
type Node struct {
	NodeId           NodeId         `json:"nodeId"`                     // Node identifier that is passed into the rest of the DOM messages as the nodeId. Backend will only push node with given id once. It is aware of all requested nodes and will only fire DOM events for nodes known to the client.
	BackendNodeId    BackendNodeId  `json:"backendNodeId"`              // The BackendNodeId for this node.
	NodeType         int            `json:"nodeType"`                   // Node's nodeType.
	NodeName         string         `json:"nodeName"`                   // Node's nodeName.
	LocalName        string         `json:"localName"`                  // Node's localName.
	NodeValue        string         `json:"nodeValue"`                  // Node's nodeValue.
	ChildNodeCount   int            `json:"childNodeCount,omitempty"`   // Child count for Container nodes.
	Children         []*Node        `json:"children,omitempty"`         // Child nodes of this node when requested with children.
	Attributes       []string       `json:"attributes,omitempty"`       // Attributes of the Element node in the form of flat array [name1, value1, name2, value2].
	DocumentURL      string         `json:"documentURL,omitempty"`      // Document URL that Document or FrameOwner node points to.
	BaseURL          string         `json:"baseURL,omitempty"`          // Base URL that Document or FrameOwner node uses for URL completion.
	PublicId         string         `json:"publicId,omitempty"`         // DocumentType's publicId.
	SystemId         string         `json:"systemId,omitempty"`         // DocumentType's systemId.
	InternalSubset   string         `json:"internalSubset,omitempty"`   // DocumentType's internalSubset.
	XmlVersion       string         `json:"xmlVersion,omitempty"`       // Document's XML version in case of XML documents.
	Name             string         `json:"name,omitempty"`             // Attr's name.
	Value            string         `json:"value,omitempty"`            // Attr's value.
	PseudoType       PseudoType     `json:"pseudoType,omitempty"`       // Pseudo element type for this node.
	ShadowRootType   ShadowRootType `json:"shadowRootType,omitempty"`   // Shadow root type.
	FrameId          *FrameId       `json:"frameId,omitempty"`          // Frame ID for frame owner elements.
	ContentDocument  *Node          `json:"contentDocument,omitempty"`  // Content document for frame owner elements.
	ShadowRoots      []*Node        `json:"shadowRoots,omitempty"`      // Shadow root list for given element host.
	TemplateContent  *Node          `json:"templateContent,omitempty"`  // Content document fragment for template elements.
	PseudoElements   []*Node        `json:"pseudoElements,omitempty"`   // Pseudo elements associated with this node.
	ImportedDocument *Node          `json:"importedDocument,omitempty"` // Import document for the HTMLImport links.
	DistributedNodes []*BackendNode `json:"distributedNodes,omitempty"` // Distributed nodes for given insertion point.
	IsSVG            bool           `json:"isSVG,omitempty"`            // Whether the node is SVG.
}

// A structure holding an RGBA color.
type RGBA struct {
	R int     `json:"r"`           // The red component, in the [0-255] range.
	G int     `json:"g"`           // The green component, in the [0-255] range.
	B int     `json:"b"`           // The blue component, in the [0-255] range.
	A float64 `json:"a,omitempty"` // The alpha component, in the [0-1] range (default: 1).
}

// An array of quad vertices, x immediately followed by y for each point, points clock-wise.
// @experimental
type Quad []float64

// Box model.
// @experimental
type BoxModel struct {
	Content      Quad              `json:"content"`                // Content box
	Padding      Quad              `json:"padding"`                // Padding box
	Border       Quad              `json:"border"`                 // Border box
	Margin       Quad              `json:"margin"`                 // Margin box
	Width        int               `json:"width"`                  // Node width
	Height       int               `json:"height"`                 // Node height
	ShapeOutside *ShapeOutsideInfo `json:"shapeOutside,omitempty"` // Shape outside coordinates
}

// CSS Shape Outside details.
// @experimental
type ShapeOutsideInfo struct {
	Bounds      Quad              `json:"bounds"`      // Shape bounds
	Shape       []json.RawMessage `json:"shape"`       // Shape coordinate details
	MarginShape []json.RawMessage `json:"marginShape"` // Margin shape bounds
}

// Rectangle.
// @experimental
type Rect struct {
	X      float64 `json:"x"`      // X coordinate
	Y      float64 `json:"y"`      // Y coordinate
	Width  float64 `json:"width"`  // Rectangle width
	Height float64 `json:"height"` // Rectangle height
}

// Configuration data for the highlighting of page elements.
type HighlightConfig struct {
	ShowInfo           bool   `json:"showInfo,omitempty"`           // Whether the node info tooltip should be shown (default: false).
	ShowRulers         bool   `json:"showRulers,omitempty"`         // Whether the rulers should be shown (default: false).
	ShowExtensionLines bool   `json:"showExtensionLines,omitempty"` // Whether the extension lines from node to the rulers should be shown (default: false).
	DisplayAsMaterial  bool   `json:"displayAsMaterial,omitempty"`
	ContentColor       *RGBA  `json:"contentColor,omitempty"`     // The content box highlight fill color (default: transparent).
	PaddingColor       *RGBA  `json:"paddingColor,omitempty"`     // The padding highlight fill color (default: transparent).
	BorderColor        *RGBA  `json:"borderColor,omitempty"`      // The border highlight fill color (default: transparent).
	MarginColor        *RGBA  `json:"marginColor,omitempty"`      // The margin highlight fill color (default: transparent).
	EventTargetColor   *RGBA  `json:"eventTargetColor,omitempty"` // The event target element highlight fill color (default: transparent).
	ShapeColor         *RGBA  `json:"shapeColor,omitempty"`       // The shape outside fill color (default: transparent).
	ShapeMarginColor   *RGBA  `json:"shapeMarginColor,omitempty"` // The shape margin fill color (default: transparent).
	SelectorList       string `json:"selectorList,omitempty"`     // Selectors to highlight relevant nodes.
}

// @experimental
type InspectMode string

const InspectModeSearchForNode InspectMode = "searchForNode"
const InspectModeSearchForUAShadowDOM InspectMode = "searchForUAShadowDOM"
const InspectModeNone InspectMode = "none"

// Enables DOM agent for the given page.

type DOMEnableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewDOMEnableCommand() *DOMEnableCommand {
	return &DOMEnableCommand{}
}

func (cmd *DOMEnableCommand) Name() string {
	return "DOM.enable"
}

func (cmd *DOMEnableCommand) Params() interface{} {
	return nil
}

func (cmd *DOMEnableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func DOMEnable(conn *hc.Conn) (err error) {
	cmd := NewDOMEnableCommand()
	cmd.Run(conn)
	return cmd.err
}

type DOMEnableCB func(err error)

// Enables DOM agent for the given page.

type AsyncDOMEnableCommand struct {
	cb DOMEnableCB
}

func NewAsyncDOMEnableCommand(cb DOMEnableCB) *AsyncDOMEnableCommand {
	return &AsyncDOMEnableCommand{
		cb: cb,
	}
}

func (cmd *AsyncDOMEnableCommand) Name() string {
	return "DOM.enable"
}

func (cmd *AsyncDOMEnableCommand) Params() interface{} {
	return nil
}

func (cmd *DOMEnableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncDOMEnableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Disables DOM agent for the given page.

type DOMDisableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewDOMDisableCommand() *DOMDisableCommand {
	return &DOMDisableCommand{}
}

func (cmd *DOMDisableCommand) Name() string {
	return "DOM.disable"
}

func (cmd *DOMDisableCommand) Params() interface{} {
	return nil
}

func (cmd *DOMDisableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func DOMDisable(conn *hc.Conn) (err error) {
	cmd := NewDOMDisableCommand()
	cmd.Run(conn)
	return cmd.err
}

type DOMDisableCB func(err error)

// Disables DOM agent for the given page.

type AsyncDOMDisableCommand struct {
	cb DOMDisableCB
}

func NewAsyncDOMDisableCommand(cb DOMDisableCB) *AsyncDOMDisableCommand {
	return &AsyncDOMDisableCommand{
		cb: cb,
	}
}

func (cmd *AsyncDOMDisableCommand) Name() string {
	return "DOM.disable"
}

func (cmd *AsyncDOMDisableCommand) Params() interface{} {
	return nil
}

func (cmd *DOMDisableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncDOMDisableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type GetDocumentParams struct {
	Depth  int  `json:"depth,omitempty"`  // The maximum depth at which children should be retrieved, defaults to 1. Use -1 for the entire subtree or provide an integer larger than 0.
	Pierce bool `json:"pierce,omitempty"` // Whether or not iframes and shadow roots should be traversed when returning the subtree (default is false).
}

type GetDocumentResult struct {
	Root *Node `json:"root"` // Resulting node.
}

// Returns the root DOM node (and optionally the subtree) to the caller.

type GetDocumentCommand struct {
	params *GetDocumentParams
	result GetDocumentResult
	wg     sync.WaitGroup
	err    error
}

func NewGetDocumentCommand(params *GetDocumentParams) *GetDocumentCommand {
	return &GetDocumentCommand{
		params: params,
	}
}

func (cmd *GetDocumentCommand) Name() string {
	return "DOM.getDocument"
}

func (cmd *GetDocumentCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetDocumentCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetDocument(params *GetDocumentParams, conn *hc.Conn) (result *GetDocumentResult, err error) {
	cmd := NewGetDocumentCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetDocumentCB func(result *GetDocumentResult, err error)

// Returns the root DOM node (and optionally the subtree) to the caller.

type AsyncGetDocumentCommand struct {
	params *GetDocumentParams
	cb     GetDocumentCB
}

func NewAsyncGetDocumentCommand(params *GetDocumentParams, cb GetDocumentCB) *AsyncGetDocumentCommand {
	return &AsyncGetDocumentCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetDocumentCommand) Name() string {
	return "DOM.getDocument"
}

func (cmd *AsyncGetDocumentCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetDocumentCommand) Result() *GetDocumentResult {
	return &cmd.result
}

func (cmd *GetDocumentCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetDocumentCommand) Done(data []byte, err error) {
	var result GetDocumentResult
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

type CollectClassNamesFromSubtreeParams struct {
	NodeId NodeId `json:"nodeId"` // Id of the node to collect class names.
}

type CollectClassNamesFromSubtreeResult struct {
	ClassNames []string `json:"classNames"` // Class name list.
}

// Collects class names for the node with given id and all of it's child nodes.
// @experimental
type CollectClassNamesFromSubtreeCommand struct {
	params *CollectClassNamesFromSubtreeParams
	result CollectClassNamesFromSubtreeResult
	wg     sync.WaitGroup
	err    error
}

func NewCollectClassNamesFromSubtreeCommand(params *CollectClassNamesFromSubtreeParams) *CollectClassNamesFromSubtreeCommand {
	return &CollectClassNamesFromSubtreeCommand{
		params: params,
	}
}

func (cmd *CollectClassNamesFromSubtreeCommand) Name() string {
	return "DOM.collectClassNamesFromSubtree"
}

func (cmd *CollectClassNamesFromSubtreeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CollectClassNamesFromSubtreeCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func CollectClassNamesFromSubtree(params *CollectClassNamesFromSubtreeParams, conn *hc.Conn) (result *CollectClassNamesFromSubtreeResult, err error) {
	cmd := NewCollectClassNamesFromSubtreeCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type CollectClassNamesFromSubtreeCB func(result *CollectClassNamesFromSubtreeResult, err error)

// Collects class names for the node with given id and all of it's child nodes.
// @experimental
type AsyncCollectClassNamesFromSubtreeCommand struct {
	params *CollectClassNamesFromSubtreeParams
	cb     CollectClassNamesFromSubtreeCB
}

func NewAsyncCollectClassNamesFromSubtreeCommand(params *CollectClassNamesFromSubtreeParams, cb CollectClassNamesFromSubtreeCB) *AsyncCollectClassNamesFromSubtreeCommand {
	return &AsyncCollectClassNamesFromSubtreeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncCollectClassNamesFromSubtreeCommand) Name() string {
	return "DOM.collectClassNamesFromSubtree"
}

func (cmd *AsyncCollectClassNamesFromSubtreeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CollectClassNamesFromSubtreeCommand) Result() *CollectClassNamesFromSubtreeResult {
	return &cmd.result
}

func (cmd *CollectClassNamesFromSubtreeCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncCollectClassNamesFromSubtreeCommand) Done(data []byte, err error) {
	var result CollectClassNamesFromSubtreeResult
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

type RequestChildNodesParams struct {
	NodeId NodeId `json:"nodeId"`           // Id of the node to get children for.
	Depth  int    `json:"depth,omitempty"`  // The maximum depth at which children should be retrieved, defaults to 1. Use -1 for the entire subtree or provide an integer larger than 0.
	Pierce bool   `json:"pierce,omitempty"` // Whether or not iframes and shadow roots should be traversed when returning the sub-tree (default is false).
}

// Requests that children of the node with given id are returned to the caller in form of setChildNodes events where not only immediate children are retrieved, but all children down to the specified depth.

type RequestChildNodesCommand struct {
	params *RequestChildNodesParams
	wg     sync.WaitGroup
	err    error
}

func NewRequestChildNodesCommand(params *RequestChildNodesParams) *RequestChildNodesCommand {
	return &RequestChildNodesCommand{
		params: params,
	}
}

func (cmd *RequestChildNodesCommand) Name() string {
	return "DOM.requestChildNodes"
}

func (cmd *RequestChildNodesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RequestChildNodesCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func RequestChildNodes(params *RequestChildNodesParams, conn *hc.Conn) (err error) {
	cmd := NewRequestChildNodesCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type RequestChildNodesCB func(err error)

// Requests that children of the node with given id are returned to the caller in form of setChildNodes events where not only immediate children are retrieved, but all children down to the specified depth.

type AsyncRequestChildNodesCommand struct {
	params *RequestChildNodesParams
	cb     RequestChildNodesCB
}

func NewAsyncRequestChildNodesCommand(params *RequestChildNodesParams, cb RequestChildNodesCB) *AsyncRequestChildNodesCommand {
	return &AsyncRequestChildNodesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncRequestChildNodesCommand) Name() string {
	return "DOM.requestChildNodes"
}

func (cmd *AsyncRequestChildNodesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RequestChildNodesCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncRequestChildNodesCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type QuerySelectorParams struct {
	NodeId   NodeId `json:"nodeId"`   // Id of the node to query upon.
	Selector string `json:"selector"` // Selector string.
}

type QuerySelectorResult struct {
	NodeId NodeId `json:"nodeId"` // Query selector result.
}

// Executes querySelector on a given node.

type QuerySelectorCommand struct {
	params *QuerySelectorParams
	result QuerySelectorResult
	wg     sync.WaitGroup
	err    error
}

func NewQuerySelectorCommand(params *QuerySelectorParams) *QuerySelectorCommand {
	return &QuerySelectorCommand{
		params: params,
	}
}

func (cmd *QuerySelectorCommand) Name() string {
	return "DOM.querySelector"
}

func (cmd *QuerySelectorCommand) Params() interface{} {
	return cmd.params
}

func (cmd *QuerySelectorCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func QuerySelector(params *QuerySelectorParams, conn *hc.Conn) (result *QuerySelectorResult, err error) {
	cmd := NewQuerySelectorCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type QuerySelectorCB func(result *QuerySelectorResult, err error)

// Executes querySelector on a given node.

type AsyncQuerySelectorCommand struct {
	params *QuerySelectorParams
	cb     QuerySelectorCB
}

func NewAsyncQuerySelectorCommand(params *QuerySelectorParams, cb QuerySelectorCB) *AsyncQuerySelectorCommand {
	return &AsyncQuerySelectorCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncQuerySelectorCommand) Name() string {
	return "DOM.querySelector"
}

func (cmd *AsyncQuerySelectorCommand) Params() interface{} {
	return cmd.params
}

func (cmd *QuerySelectorCommand) Result() *QuerySelectorResult {
	return &cmd.result
}

func (cmd *QuerySelectorCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncQuerySelectorCommand) Done(data []byte, err error) {
	var result QuerySelectorResult
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

type QuerySelectorAllParams struct {
	NodeId   NodeId `json:"nodeId"`   // Id of the node to query upon.
	Selector string `json:"selector"` // Selector string.
}

type QuerySelectorAllResult struct {
	NodeIds []NodeId `json:"nodeIds"` // Query selector result.
}

// Executes querySelectorAll on a given node.

type QuerySelectorAllCommand struct {
	params *QuerySelectorAllParams
	result QuerySelectorAllResult
	wg     sync.WaitGroup
	err    error
}

func NewQuerySelectorAllCommand(params *QuerySelectorAllParams) *QuerySelectorAllCommand {
	return &QuerySelectorAllCommand{
		params: params,
	}
}

func (cmd *QuerySelectorAllCommand) Name() string {
	return "DOM.querySelectorAll"
}

func (cmd *QuerySelectorAllCommand) Params() interface{} {
	return cmd.params
}

func (cmd *QuerySelectorAllCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func QuerySelectorAll(params *QuerySelectorAllParams, conn *hc.Conn) (result *QuerySelectorAllResult, err error) {
	cmd := NewQuerySelectorAllCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type QuerySelectorAllCB func(result *QuerySelectorAllResult, err error)

// Executes querySelectorAll on a given node.

type AsyncQuerySelectorAllCommand struct {
	params *QuerySelectorAllParams
	cb     QuerySelectorAllCB
}

func NewAsyncQuerySelectorAllCommand(params *QuerySelectorAllParams, cb QuerySelectorAllCB) *AsyncQuerySelectorAllCommand {
	return &AsyncQuerySelectorAllCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncQuerySelectorAllCommand) Name() string {
	return "DOM.querySelectorAll"
}

func (cmd *AsyncQuerySelectorAllCommand) Params() interface{} {
	return cmd.params
}

func (cmd *QuerySelectorAllCommand) Result() *QuerySelectorAllResult {
	return &cmd.result
}

func (cmd *QuerySelectorAllCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncQuerySelectorAllCommand) Done(data []byte, err error) {
	var result QuerySelectorAllResult
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

type SetNodeNameParams struct {
	NodeId NodeId `json:"nodeId"` // Id of the node to set name for.
	Name   string `json:"name"`   // New node's name.
}

type SetNodeNameResult struct {
	NodeId NodeId `json:"nodeId"` // New node's id.
}

// Sets node name for a node with given id.

type SetNodeNameCommand struct {
	params *SetNodeNameParams
	result SetNodeNameResult
	wg     sync.WaitGroup
	err    error
}

func NewSetNodeNameCommand(params *SetNodeNameParams) *SetNodeNameCommand {
	return &SetNodeNameCommand{
		params: params,
	}
}

func (cmd *SetNodeNameCommand) Name() string {
	return "DOM.setNodeName"
}

func (cmd *SetNodeNameCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetNodeNameCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetNodeName(params *SetNodeNameParams, conn *hc.Conn) (result *SetNodeNameResult, err error) {
	cmd := NewSetNodeNameCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type SetNodeNameCB func(result *SetNodeNameResult, err error)

// Sets node name for a node with given id.

type AsyncSetNodeNameCommand struct {
	params *SetNodeNameParams
	cb     SetNodeNameCB
}

func NewAsyncSetNodeNameCommand(params *SetNodeNameParams, cb SetNodeNameCB) *AsyncSetNodeNameCommand {
	return &AsyncSetNodeNameCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetNodeNameCommand) Name() string {
	return "DOM.setNodeName"
}

func (cmd *AsyncSetNodeNameCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetNodeNameCommand) Result() *SetNodeNameResult {
	return &cmd.result
}

func (cmd *SetNodeNameCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetNodeNameCommand) Done(data []byte, err error) {
	var result SetNodeNameResult
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

type SetNodeValueParams struct {
	NodeId NodeId `json:"nodeId"` // Id of the node to set value for.
	Value  string `json:"value"`  // New node's value.
}

// Sets node value for a node with given id.

type SetNodeValueCommand struct {
	params *SetNodeValueParams
	wg     sync.WaitGroup
	err    error
}

func NewSetNodeValueCommand(params *SetNodeValueParams) *SetNodeValueCommand {
	return &SetNodeValueCommand{
		params: params,
	}
}

func (cmd *SetNodeValueCommand) Name() string {
	return "DOM.setNodeValue"
}

func (cmd *SetNodeValueCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetNodeValueCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetNodeValue(params *SetNodeValueParams, conn *hc.Conn) (err error) {
	cmd := NewSetNodeValueCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetNodeValueCB func(err error)

// Sets node value for a node with given id.

type AsyncSetNodeValueCommand struct {
	params *SetNodeValueParams
	cb     SetNodeValueCB
}

func NewAsyncSetNodeValueCommand(params *SetNodeValueParams, cb SetNodeValueCB) *AsyncSetNodeValueCommand {
	return &AsyncSetNodeValueCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetNodeValueCommand) Name() string {
	return "DOM.setNodeValue"
}

func (cmd *AsyncSetNodeValueCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetNodeValueCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetNodeValueCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type RemoveNodeParams struct {
	NodeId NodeId `json:"nodeId"` // Id of the node to remove.
}

// Removes node with given id.

type RemoveNodeCommand struct {
	params *RemoveNodeParams
	wg     sync.WaitGroup
	err    error
}

func NewRemoveNodeCommand(params *RemoveNodeParams) *RemoveNodeCommand {
	return &RemoveNodeCommand{
		params: params,
	}
}

func (cmd *RemoveNodeCommand) Name() string {
	return "DOM.removeNode"
}

func (cmd *RemoveNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveNodeCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func RemoveNode(params *RemoveNodeParams, conn *hc.Conn) (err error) {
	cmd := NewRemoveNodeCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type RemoveNodeCB func(err error)

// Removes node with given id.

type AsyncRemoveNodeCommand struct {
	params *RemoveNodeParams
	cb     RemoveNodeCB
}

func NewAsyncRemoveNodeCommand(params *RemoveNodeParams, cb RemoveNodeCB) *AsyncRemoveNodeCommand {
	return &AsyncRemoveNodeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncRemoveNodeCommand) Name() string {
	return "DOM.removeNode"
}

func (cmd *AsyncRemoveNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveNodeCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncRemoveNodeCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetAttributeValueParams struct {
	NodeId NodeId `json:"nodeId"` // Id of the element to set attribute for.
	Name   string `json:"name"`   // Attribute name.
	Value  string `json:"value"`  // Attribute value.
}

// Sets attribute for an element with given id.

type SetAttributeValueCommand struct {
	params *SetAttributeValueParams
	wg     sync.WaitGroup
	err    error
}

func NewSetAttributeValueCommand(params *SetAttributeValueParams) *SetAttributeValueCommand {
	return &SetAttributeValueCommand{
		params: params,
	}
}

func (cmd *SetAttributeValueCommand) Name() string {
	return "DOM.setAttributeValue"
}

func (cmd *SetAttributeValueCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetAttributeValueCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetAttributeValue(params *SetAttributeValueParams, conn *hc.Conn) (err error) {
	cmd := NewSetAttributeValueCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetAttributeValueCB func(err error)

// Sets attribute for an element with given id.

type AsyncSetAttributeValueCommand struct {
	params *SetAttributeValueParams
	cb     SetAttributeValueCB
}

func NewAsyncSetAttributeValueCommand(params *SetAttributeValueParams, cb SetAttributeValueCB) *AsyncSetAttributeValueCommand {
	return &AsyncSetAttributeValueCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetAttributeValueCommand) Name() string {
	return "DOM.setAttributeValue"
}

func (cmd *AsyncSetAttributeValueCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetAttributeValueCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetAttributeValueCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetAttributesAsTextParams struct {
	NodeId NodeId `json:"nodeId"`         // Id of the element to set attributes for.
	Text   string `json:"text"`           // Text with a number of attributes. Will parse this text using HTML parser.
	Name   string `json:"name,omitempty"` // Attribute name to replace with new attributes derived from text in case text parsed successfully.
}

// Sets attributes on element with given id. This method is useful when user edits some existing attribute value and types in several attribute name/value pairs.

type SetAttributesAsTextCommand struct {
	params *SetAttributesAsTextParams
	wg     sync.WaitGroup
	err    error
}

func NewSetAttributesAsTextCommand(params *SetAttributesAsTextParams) *SetAttributesAsTextCommand {
	return &SetAttributesAsTextCommand{
		params: params,
	}
}

func (cmd *SetAttributesAsTextCommand) Name() string {
	return "DOM.setAttributesAsText"
}

func (cmd *SetAttributesAsTextCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetAttributesAsTextCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetAttributesAsText(params *SetAttributesAsTextParams, conn *hc.Conn) (err error) {
	cmd := NewSetAttributesAsTextCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetAttributesAsTextCB func(err error)

// Sets attributes on element with given id. This method is useful when user edits some existing attribute value and types in several attribute name/value pairs.

type AsyncSetAttributesAsTextCommand struct {
	params *SetAttributesAsTextParams
	cb     SetAttributesAsTextCB
}

func NewAsyncSetAttributesAsTextCommand(params *SetAttributesAsTextParams, cb SetAttributesAsTextCB) *AsyncSetAttributesAsTextCommand {
	return &AsyncSetAttributesAsTextCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetAttributesAsTextCommand) Name() string {
	return "DOM.setAttributesAsText"
}

func (cmd *AsyncSetAttributesAsTextCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetAttributesAsTextCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetAttributesAsTextCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type RemoveAttributeParams struct {
	NodeId NodeId `json:"nodeId"` // Id of the element to remove attribute from.
	Name   string `json:"name"`   // Name of the attribute to remove.
}

// Removes attribute with given name from an element with given id.

type RemoveAttributeCommand struct {
	params *RemoveAttributeParams
	wg     sync.WaitGroup
	err    error
}

func NewRemoveAttributeCommand(params *RemoveAttributeParams) *RemoveAttributeCommand {
	return &RemoveAttributeCommand{
		params: params,
	}
}

func (cmd *RemoveAttributeCommand) Name() string {
	return "DOM.removeAttribute"
}

func (cmd *RemoveAttributeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveAttributeCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func RemoveAttribute(params *RemoveAttributeParams, conn *hc.Conn) (err error) {
	cmd := NewRemoveAttributeCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type RemoveAttributeCB func(err error)

// Removes attribute with given name from an element with given id.

type AsyncRemoveAttributeCommand struct {
	params *RemoveAttributeParams
	cb     RemoveAttributeCB
}

func NewAsyncRemoveAttributeCommand(params *RemoveAttributeParams, cb RemoveAttributeCB) *AsyncRemoveAttributeCommand {
	return &AsyncRemoveAttributeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncRemoveAttributeCommand) Name() string {
	return "DOM.removeAttribute"
}

func (cmd *AsyncRemoveAttributeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveAttributeCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncRemoveAttributeCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type GetOuterHTMLParams struct {
	NodeId NodeId `json:"nodeId"` // Id of the node to get markup for.
}

type GetOuterHTMLResult struct {
	OuterHTML string `json:"outerHTML"` // Outer HTML markup.
}

// Returns node's HTML markup.

type GetOuterHTMLCommand struct {
	params *GetOuterHTMLParams
	result GetOuterHTMLResult
	wg     sync.WaitGroup
	err    error
}

func NewGetOuterHTMLCommand(params *GetOuterHTMLParams) *GetOuterHTMLCommand {
	return &GetOuterHTMLCommand{
		params: params,
	}
}

func (cmd *GetOuterHTMLCommand) Name() string {
	return "DOM.getOuterHTML"
}

func (cmd *GetOuterHTMLCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetOuterHTMLCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetOuterHTML(params *GetOuterHTMLParams, conn *hc.Conn) (result *GetOuterHTMLResult, err error) {
	cmd := NewGetOuterHTMLCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetOuterHTMLCB func(result *GetOuterHTMLResult, err error)

// Returns node's HTML markup.

type AsyncGetOuterHTMLCommand struct {
	params *GetOuterHTMLParams
	cb     GetOuterHTMLCB
}

func NewAsyncGetOuterHTMLCommand(params *GetOuterHTMLParams, cb GetOuterHTMLCB) *AsyncGetOuterHTMLCommand {
	return &AsyncGetOuterHTMLCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetOuterHTMLCommand) Name() string {
	return "DOM.getOuterHTML"
}

func (cmd *AsyncGetOuterHTMLCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetOuterHTMLCommand) Result() *GetOuterHTMLResult {
	return &cmd.result
}

func (cmd *GetOuterHTMLCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetOuterHTMLCommand) Done(data []byte, err error) {
	var result GetOuterHTMLResult
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

type SetOuterHTMLParams struct {
	NodeId    NodeId `json:"nodeId"`    // Id of the node to set markup for.
	OuterHTML string `json:"outerHTML"` // Outer HTML markup to set.
}

// Sets node HTML markup, returns new node id.

type SetOuterHTMLCommand struct {
	params *SetOuterHTMLParams
	wg     sync.WaitGroup
	err    error
}

func NewSetOuterHTMLCommand(params *SetOuterHTMLParams) *SetOuterHTMLCommand {
	return &SetOuterHTMLCommand{
		params: params,
	}
}

func (cmd *SetOuterHTMLCommand) Name() string {
	return "DOM.setOuterHTML"
}

func (cmd *SetOuterHTMLCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetOuterHTMLCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetOuterHTML(params *SetOuterHTMLParams, conn *hc.Conn) (err error) {
	cmd := NewSetOuterHTMLCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetOuterHTMLCB func(err error)

// Sets node HTML markup, returns new node id.

type AsyncSetOuterHTMLCommand struct {
	params *SetOuterHTMLParams
	cb     SetOuterHTMLCB
}

func NewAsyncSetOuterHTMLCommand(params *SetOuterHTMLParams, cb SetOuterHTMLCB) *AsyncSetOuterHTMLCommand {
	return &AsyncSetOuterHTMLCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetOuterHTMLCommand) Name() string {
	return "DOM.setOuterHTML"
}

func (cmd *AsyncSetOuterHTMLCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetOuterHTMLCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetOuterHTMLCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type PerformSearchParams struct {
	Query                     string `json:"query"`                               // Plain text or query selector or XPath search query.
	IncludeUserAgentShadowDOM bool   `json:"includeUserAgentShadowDOM,omitempty"` // True to search in user agent shadow DOM.
}

type PerformSearchResult struct {
	SearchId    string `json:"searchId"`    // Unique search session identifier.
	ResultCount int    `json:"resultCount"` // Number of search results.
}

// Searches for a given string in the DOM tree. Use getSearchResults to access search results or cancelSearch to end this search session.
// @experimental
type PerformSearchCommand struct {
	params *PerformSearchParams
	result PerformSearchResult
	wg     sync.WaitGroup
	err    error
}

func NewPerformSearchCommand(params *PerformSearchParams) *PerformSearchCommand {
	return &PerformSearchCommand{
		params: params,
	}
}

func (cmd *PerformSearchCommand) Name() string {
	return "DOM.performSearch"
}

func (cmd *PerformSearchCommand) Params() interface{} {
	return cmd.params
}

func (cmd *PerformSearchCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func PerformSearch(params *PerformSearchParams, conn *hc.Conn) (result *PerformSearchResult, err error) {
	cmd := NewPerformSearchCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type PerformSearchCB func(result *PerformSearchResult, err error)

// Searches for a given string in the DOM tree. Use getSearchResults to access search results or cancelSearch to end this search session.
// @experimental
type AsyncPerformSearchCommand struct {
	params *PerformSearchParams
	cb     PerformSearchCB
}

func NewAsyncPerformSearchCommand(params *PerformSearchParams, cb PerformSearchCB) *AsyncPerformSearchCommand {
	return &AsyncPerformSearchCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncPerformSearchCommand) Name() string {
	return "DOM.performSearch"
}

func (cmd *AsyncPerformSearchCommand) Params() interface{} {
	return cmd.params
}

func (cmd *PerformSearchCommand) Result() *PerformSearchResult {
	return &cmd.result
}

func (cmd *PerformSearchCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncPerformSearchCommand) Done(data []byte, err error) {
	var result PerformSearchResult
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

type GetSearchResultsParams struct {
	SearchId  string `json:"searchId"`  // Unique search session identifier.
	FromIndex int    `json:"fromIndex"` // Start index of the search result to be returned.
	ToIndex   int    `json:"toIndex"`   // End index of the search result to be returned.
}

type GetSearchResultsResult struct {
	NodeIds []NodeId `json:"nodeIds"` // Ids of the search result nodes.
}

// Returns search results from given fromIndex to given toIndex from the sarch with the given identifier.
// @experimental
type GetSearchResultsCommand struct {
	params *GetSearchResultsParams
	result GetSearchResultsResult
	wg     sync.WaitGroup
	err    error
}

func NewGetSearchResultsCommand(params *GetSearchResultsParams) *GetSearchResultsCommand {
	return &GetSearchResultsCommand{
		params: params,
	}
}

func (cmd *GetSearchResultsCommand) Name() string {
	return "DOM.getSearchResults"
}

func (cmd *GetSearchResultsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetSearchResultsCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetSearchResults(params *GetSearchResultsParams, conn *hc.Conn) (result *GetSearchResultsResult, err error) {
	cmd := NewGetSearchResultsCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetSearchResultsCB func(result *GetSearchResultsResult, err error)

// Returns search results from given fromIndex to given toIndex from the sarch with the given identifier.
// @experimental
type AsyncGetSearchResultsCommand struct {
	params *GetSearchResultsParams
	cb     GetSearchResultsCB
}

func NewAsyncGetSearchResultsCommand(params *GetSearchResultsParams, cb GetSearchResultsCB) *AsyncGetSearchResultsCommand {
	return &AsyncGetSearchResultsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetSearchResultsCommand) Name() string {
	return "DOM.getSearchResults"
}

func (cmd *AsyncGetSearchResultsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetSearchResultsCommand) Result() *GetSearchResultsResult {
	return &cmd.result
}

func (cmd *GetSearchResultsCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetSearchResultsCommand) Done(data []byte, err error) {
	var result GetSearchResultsResult
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

type DiscardSearchResultsParams struct {
	SearchId string `json:"searchId"` // Unique search session identifier.
}

// Discards search results from the session with the given id. getSearchResults should no longer be called for that search.
// @experimental
type DiscardSearchResultsCommand struct {
	params *DiscardSearchResultsParams
	wg     sync.WaitGroup
	err    error
}

func NewDiscardSearchResultsCommand(params *DiscardSearchResultsParams) *DiscardSearchResultsCommand {
	return &DiscardSearchResultsCommand{
		params: params,
	}
}

func (cmd *DiscardSearchResultsCommand) Name() string {
	return "DOM.discardSearchResults"
}

func (cmd *DiscardSearchResultsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *DiscardSearchResultsCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func DiscardSearchResults(params *DiscardSearchResultsParams, conn *hc.Conn) (err error) {
	cmd := NewDiscardSearchResultsCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type DiscardSearchResultsCB func(err error)

// Discards search results from the session with the given id. getSearchResults should no longer be called for that search.
// @experimental
type AsyncDiscardSearchResultsCommand struct {
	params *DiscardSearchResultsParams
	cb     DiscardSearchResultsCB
}

func NewAsyncDiscardSearchResultsCommand(params *DiscardSearchResultsParams, cb DiscardSearchResultsCB) *AsyncDiscardSearchResultsCommand {
	return &AsyncDiscardSearchResultsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncDiscardSearchResultsCommand) Name() string {
	return "DOM.discardSearchResults"
}

func (cmd *AsyncDiscardSearchResultsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *DiscardSearchResultsCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncDiscardSearchResultsCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type RequestNodeParams struct {
	ObjectId *RemoteObjectId `json:"objectId"` // JavaScript object id to convert into node.
}

type RequestNodeResult struct {
	NodeId NodeId `json:"nodeId"` // Node id for given object.
}

// Requests that the node is sent to the caller given the JavaScript node object reference. All nodes that form the path from the node to the root are also sent to the client as a series of setChildNodes notifications.

type RequestNodeCommand struct {
	params *RequestNodeParams
	result RequestNodeResult
	wg     sync.WaitGroup
	err    error
}

func NewRequestNodeCommand(params *RequestNodeParams) *RequestNodeCommand {
	return &RequestNodeCommand{
		params: params,
	}
}

func (cmd *RequestNodeCommand) Name() string {
	return "DOM.requestNode"
}

func (cmd *RequestNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RequestNodeCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func RequestNode(params *RequestNodeParams, conn *hc.Conn) (result *RequestNodeResult, err error) {
	cmd := NewRequestNodeCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type RequestNodeCB func(result *RequestNodeResult, err error)

// Requests that the node is sent to the caller given the JavaScript node object reference. All nodes that form the path from the node to the root are also sent to the client as a series of setChildNodes notifications.

type AsyncRequestNodeCommand struct {
	params *RequestNodeParams
	cb     RequestNodeCB
}

func NewAsyncRequestNodeCommand(params *RequestNodeParams, cb RequestNodeCB) *AsyncRequestNodeCommand {
	return &AsyncRequestNodeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncRequestNodeCommand) Name() string {
	return "DOM.requestNode"
}

func (cmd *AsyncRequestNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RequestNodeCommand) Result() *RequestNodeResult {
	return &cmd.result
}

func (cmd *RequestNodeCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncRequestNodeCommand) Done(data []byte, err error) {
	var result RequestNodeResult
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

type SetInspectModeParams struct {
	Mode            InspectMode      `json:"mode"`                      // Set an inspection mode.
	HighlightConfig *HighlightConfig `json:"highlightConfig,omitempty"` // A descriptor for the highlight appearance of hovered-over nodes. May be omitted if enabled == false.
}

// Enters the 'inspect' mode. In this mode, elements that user is hovering over are highlighted. Backend then generates 'inspectNodeRequested' event upon element selection.
// @experimental
type SetInspectModeCommand struct {
	params *SetInspectModeParams
	wg     sync.WaitGroup
	err    error
}

func NewSetInspectModeCommand(params *SetInspectModeParams) *SetInspectModeCommand {
	return &SetInspectModeCommand{
		params: params,
	}
}

func (cmd *SetInspectModeCommand) Name() string {
	return "DOM.setInspectMode"
}

func (cmd *SetInspectModeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetInspectModeCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetInspectMode(params *SetInspectModeParams, conn *hc.Conn) (err error) {
	cmd := NewSetInspectModeCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetInspectModeCB func(err error)

// Enters the 'inspect' mode. In this mode, elements that user is hovering over are highlighted. Backend then generates 'inspectNodeRequested' event upon element selection.
// @experimental
type AsyncSetInspectModeCommand struct {
	params *SetInspectModeParams
	cb     SetInspectModeCB
}

func NewAsyncSetInspectModeCommand(params *SetInspectModeParams, cb SetInspectModeCB) *AsyncSetInspectModeCommand {
	return &AsyncSetInspectModeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetInspectModeCommand) Name() string {
	return "DOM.setInspectMode"
}

func (cmd *AsyncSetInspectModeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetInspectModeCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetInspectModeCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type HighlightRectParams struct {
	X            int   `json:"x"`                      // X coordinate
	Y            int   `json:"y"`                      // Y coordinate
	Width        int   `json:"width"`                  // Rectangle width
	Height       int   `json:"height"`                 // Rectangle height
	Color        *RGBA `json:"color,omitempty"`        // The highlight fill color (default: transparent).
	OutlineColor *RGBA `json:"outlineColor,omitempty"` // The highlight outline color (default: transparent).
}

// Highlights given rectangle. Coordinates are absolute with respect to the main frame viewport.

type HighlightRectCommand struct {
	params *HighlightRectParams
	wg     sync.WaitGroup
	err    error
}

func NewHighlightRectCommand(params *HighlightRectParams) *HighlightRectCommand {
	return &HighlightRectCommand{
		params: params,
	}
}

func (cmd *HighlightRectCommand) Name() string {
	return "DOM.highlightRect"
}

func (cmd *HighlightRectCommand) Params() interface{} {
	return cmd.params
}

func (cmd *HighlightRectCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func HighlightRect(params *HighlightRectParams, conn *hc.Conn) (err error) {
	cmd := NewHighlightRectCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type HighlightRectCB func(err error)

// Highlights given rectangle. Coordinates are absolute with respect to the main frame viewport.

type AsyncHighlightRectCommand struct {
	params *HighlightRectParams
	cb     HighlightRectCB
}

func NewAsyncHighlightRectCommand(params *HighlightRectParams, cb HighlightRectCB) *AsyncHighlightRectCommand {
	return &AsyncHighlightRectCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncHighlightRectCommand) Name() string {
	return "DOM.highlightRect"
}

func (cmd *AsyncHighlightRectCommand) Params() interface{} {
	return cmd.params
}

func (cmd *HighlightRectCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncHighlightRectCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type HighlightQuadParams struct {
	Quad         Quad  `json:"quad"`                   // Quad to highlight
	Color        *RGBA `json:"color,omitempty"`        // The highlight fill color (default: transparent).
	OutlineColor *RGBA `json:"outlineColor,omitempty"` // The highlight outline color (default: transparent).
}

// Highlights given quad. Coordinates are absolute with respect to the main frame viewport.
// @experimental
type HighlightQuadCommand struct {
	params *HighlightQuadParams
	wg     sync.WaitGroup
	err    error
}

func NewHighlightQuadCommand(params *HighlightQuadParams) *HighlightQuadCommand {
	return &HighlightQuadCommand{
		params: params,
	}
}

func (cmd *HighlightQuadCommand) Name() string {
	return "DOM.highlightQuad"
}

func (cmd *HighlightQuadCommand) Params() interface{} {
	return cmd.params
}

func (cmd *HighlightQuadCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func HighlightQuad(params *HighlightQuadParams, conn *hc.Conn) (err error) {
	cmd := NewHighlightQuadCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type HighlightQuadCB func(err error)

// Highlights given quad. Coordinates are absolute with respect to the main frame viewport.
// @experimental
type AsyncHighlightQuadCommand struct {
	params *HighlightQuadParams
	cb     HighlightQuadCB
}

func NewAsyncHighlightQuadCommand(params *HighlightQuadParams, cb HighlightQuadCB) *AsyncHighlightQuadCommand {
	return &AsyncHighlightQuadCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncHighlightQuadCommand) Name() string {
	return "DOM.highlightQuad"
}

func (cmd *AsyncHighlightQuadCommand) Params() interface{} {
	return cmd.params
}

func (cmd *HighlightQuadCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncHighlightQuadCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type HighlightNodeParams struct {
	HighlightConfig *HighlightConfig `json:"highlightConfig"`         // A descriptor for the highlight appearance.
	NodeId          NodeId           `json:"nodeId,omitempty"`        // Identifier of the node to highlight.
	BackendNodeId   BackendNodeId    `json:"backendNodeId,omitempty"` // Identifier of the backend node to highlight.
	ObjectId        *RemoteObjectId  `json:"objectId,omitempty"`      // JavaScript object id of the node to be highlighted.
}

// Highlights DOM node with given id or with the given JavaScript object wrapper. Either nodeId or objectId must be specified.

type HighlightNodeCommand struct {
	params *HighlightNodeParams
	wg     sync.WaitGroup
	err    error
}

func NewHighlightNodeCommand(params *HighlightNodeParams) *HighlightNodeCommand {
	return &HighlightNodeCommand{
		params: params,
	}
}

func (cmd *HighlightNodeCommand) Name() string {
	return "DOM.highlightNode"
}

func (cmd *HighlightNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *HighlightNodeCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func HighlightNode(params *HighlightNodeParams, conn *hc.Conn) (err error) {
	cmd := NewHighlightNodeCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type HighlightNodeCB func(err error)

// Highlights DOM node with given id or with the given JavaScript object wrapper. Either nodeId or objectId must be specified.

type AsyncHighlightNodeCommand struct {
	params *HighlightNodeParams
	cb     HighlightNodeCB
}

func NewAsyncHighlightNodeCommand(params *HighlightNodeParams, cb HighlightNodeCB) *AsyncHighlightNodeCommand {
	return &AsyncHighlightNodeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncHighlightNodeCommand) Name() string {
	return "DOM.highlightNode"
}

func (cmd *AsyncHighlightNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *HighlightNodeCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncHighlightNodeCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Hides DOM node highlight.

type HideHighlightCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewHideHighlightCommand() *HideHighlightCommand {
	return &HideHighlightCommand{}
}

func (cmd *HideHighlightCommand) Name() string {
	return "DOM.hideHighlight"
}

func (cmd *HideHighlightCommand) Params() interface{} {
	return nil
}

func (cmd *HideHighlightCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func HideHighlight(conn *hc.Conn) (err error) {
	cmd := NewHideHighlightCommand()
	cmd.Run(conn)
	return cmd.err
}

type HideHighlightCB func(err error)

// Hides DOM node highlight.

type AsyncHideHighlightCommand struct {
	cb HideHighlightCB
}

func NewAsyncHideHighlightCommand(cb HideHighlightCB) *AsyncHideHighlightCommand {
	return &AsyncHideHighlightCommand{
		cb: cb,
	}
}

func (cmd *AsyncHideHighlightCommand) Name() string {
	return "DOM.hideHighlight"
}

func (cmd *AsyncHideHighlightCommand) Params() interface{} {
	return nil
}

func (cmd *HideHighlightCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncHideHighlightCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type HighlightFrameParams struct {
	FrameId             *FrameId `json:"frameId"`                       // Identifier of the frame to highlight.
	ContentColor        *RGBA    `json:"contentColor,omitempty"`        // The content box highlight fill color (default: transparent).
	ContentOutlineColor *RGBA    `json:"contentOutlineColor,omitempty"` // The content box highlight outline color (default: transparent).
}

// Highlights owner element of the frame with given id.
// @experimental
type HighlightFrameCommand struct {
	params *HighlightFrameParams
	wg     sync.WaitGroup
	err    error
}

func NewHighlightFrameCommand(params *HighlightFrameParams) *HighlightFrameCommand {
	return &HighlightFrameCommand{
		params: params,
	}
}

func (cmd *HighlightFrameCommand) Name() string {
	return "DOM.highlightFrame"
}

func (cmd *HighlightFrameCommand) Params() interface{} {
	return cmd.params
}

func (cmd *HighlightFrameCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func HighlightFrame(params *HighlightFrameParams, conn *hc.Conn) (err error) {
	cmd := NewHighlightFrameCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type HighlightFrameCB func(err error)

// Highlights owner element of the frame with given id.
// @experimental
type AsyncHighlightFrameCommand struct {
	params *HighlightFrameParams
	cb     HighlightFrameCB
}

func NewAsyncHighlightFrameCommand(params *HighlightFrameParams, cb HighlightFrameCB) *AsyncHighlightFrameCommand {
	return &AsyncHighlightFrameCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncHighlightFrameCommand) Name() string {
	return "DOM.highlightFrame"
}

func (cmd *AsyncHighlightFrameCommand) Params() interface{} {
	return cmd.params
}

func (cmd *HighlightFrameCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncHighlightFrameCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type PushNodeByPathToFrontendParams struct {
	Path string `json:"path"` // Path to node in the proprietary format.
}

type PushNodeByPathToFrontendResult struct {
	NodeId NodeId `json:"nodeId"` // Id of the node for given path.
}

// Requests that the node is sent to the caller given its path. // FIXME, use XPath
// @experimental
type PushNodeByPathToFrontendCommand struct {
	params *PushNodeByPathToFrontendParams
	result PushNodeByPathToFrontendResult
	wg     sync.WaitGroup
	err    error
}

func NewPushNodeByPathToFrontendCommand(params *PushNodeByPathToFrontendParams) *PushNodeByPathToFrontendCommand {
	return &PushNodeByPathToFrontendCommand{
		params: params,
	}
}

func (cmd *PushNodeByPathToFrontendCommand) Name() string {
	return "DOM.pushNodeByPathToFrontend"
}

func (cmd *PushNodeByPathToFrontendCommand) Params() interface{} {
	return cmd.params
}

func (cmd *PushNodeByPathToFrontendCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func PushNodeByPathToFrontend(params *PushNodeByPathToFrontendParams, conn *hc.Conn) (result *PushNodeByPathToFrontendResult, err error) {
	cmd := NewPushNodeByPathToFrontendCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type PushNodeByPathToFrontendCB func(result *PushNodeByPathToFrontendResult, err error)

// Requests that the node is sent to the caller given its path. // FIXME, use XPath
// @experimental
type AsyncPushNodeByPathToFrontendCommand struct {
	params *PushNodeByPathToFrontendParams
	cb     PushNodeByPathToFrontendCB
}

func NewAsyncPushNodeByPathToFrontendCommand(params *PushNodeByPathToFrontendParams, cb PushNodeByPathToFrontendCB) *AsyncPushNodeByPathToFrontendCommand {
	return &AsyncPushNodeByPathToFrontendCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncPushNodeByPathToFrontendCommand) Name() string {
	return "DOM.pushNodeByPathToFrontend"
}

func (cmd *AsyncPushNodeByPathToFrontendCommand) Params() interface{} {
	return cmd.params
}

func (cmd *PushNodeByPathToFrontendCommand) Result() *PushNodeByPathToFrontendResult {
	return &cmd.result
}

func (cmd *PushNodeByPathToFrontendCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncPushNodeByPathToFrontendCommand) Done(data []byte, err error) {
	var result PushNodeByPathToFrontendResult
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

type PushNodesByBackendIdsToFrontendParams struct {
	BackendNodeIds []BackendNodeId `json:"backendNodeIds"` // The array of backend node ids.
}

type PushNodesByBackendIdsToFrontendResult struct {
	NodeIds []NodeId `json:"nodeIds"` // The array of ids of pushed nodes that correspond to the backend ids specified in backendNodeIds.
}

// Requests that a batch of nodes is sent to the caller given their backend node ids.
// @experimental
type PushNodesByBackendIdsToFrontendCommand struct {
	params *PushNodesByBackendIdsToFrontendParams
	result PushNodesByBackendIdsToFrontendResult
	wg     sync.WaitGroup
	err    error
}

func NewPushNodesByBackendIdsToFrontendCommand(params *PushNodesByBackendIdsToFrontendParams) *PushNodesByBackendIdsToFrontendCommand {
	return &PushNodesByBackendIdsToFrontendCommand{
		params: params,
	}
}

func (cmd *PushNodesByBackendIdsToFrontendCommand) Name() string {
	return "DOM.pushNodesByBackendIdsToFrontend"
}

func (cmd *PushNodesByBackendIdsToFrontendCommand) Params() interface{} {
	return cmd.params
}

func (cmd *PushNodesByBackendIdsToFrontendCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func PushNodesByBackendIdsToFrontend(params *PushNodesByBackendIdsToFrontendParams, conn *hc.Conn) (result *PushNodesByBackendIdsToFrontendResult, err error) {
	cmd := NewPushNodesByBackendIdsToFrontendCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type PushNodesByBackendIdsToFrontendCB func(result *PushNodesByBackendIdsToFrontendResult, err error)

// Requests that a batch of nodes is sent to the caller given their backend node ids.
// @experimental
type AsyncPushNodesByBackendIdsToFrontendCommand struct {
	params *PushNodesByBackendIdsToFrontendParams
	cb     PushNodesByBackendIdsToFrontendCB
}

func NewAsyncPushNodesByBackendIdsToFrontendCommand(params *PushNodesByBackendIdsToFrontendParams, cb PushNodesByBackendIdsToFrontendCB) *AsyncPushNodesByBackendIdsToFrontendCommand {
	return &AsyncPushNodesByBackendIdsToFrontendCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncPushNodesByBackendIdsToFrontendCommand) Name() string {
	return "DOM.pushNodesByBackendIdsToFrontend"
}

func (cmd *AsyncPushNodesByBackendIdsToFrontendCommand) Params() interface{} {
	return cmd.params
}

func (cmd *PushNodesByBackendIdsToFrontendCommand) Result() *PushNodesByBackendIdsToFrontendResult {
	return &cmd.result
}

func (cmd *PushNodesByBackendIdsToFrontendCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncPushNodesByBackendIdsToFrontendCommand) Done(data []byte, err error) {
	var result PushNodesByBackendIdsToFrontendResult
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

type SetInspectedNodeParams struct {
	NodeId NodeId `json:"nodeId"` // DOM node id to be accessible by means of $x command line API.
}

// Enables console to refer to the node with given id via $x (see Command Line API for more details $x functions).
// @experimental
type SetInspectedNodeCommand struct {
	params *SetInspectedNodeParams
	wg     sync.WaitGroup
	err    error
}

func NewSetInspectedNodeCommand(params *SetInspectedNodeParams) *SetInspectedNodeCommand {
	return &SetInspectedNodeCommand{
		params: params,
	}
}

func (cmd *SetInspectedNodeCommand) Name() string {
	return "DOM.setInspectedNode"
}

func (cmd *SetInspectedNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetInspectedNodeCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetInspectedNode(params *SetInspectedNodeParams, conn *hc.Conn) (err error) {
	cmd := NewSetInspectedNodeCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetInspectedNodeCB func(err error)

// Enables console to refer to the node with given id via $x (see Command Line API for more details $x functions).
// @experimental
type AsyncSetInspectedNodeCommand struct {
	params *SetInspectedNodeParams
	cb     SetInspectedNodeCB
}

func NewAsyncSetInspectedNodeCommand(params *SetInspectedNodeParams, cb SetInspectedNodeCB) *AsyncSetInspectedNodeCommand {
	return &AsyncSetInspectedNodeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetInspectedNodeCommand) Name() string {
	return "DOM.setInspectedNode"
}

func (cmd *AsyncSetInspectedNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetInspectedNodeCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetInspectedNodeCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type ResolveNodeParams struct {
	NodeId      NodeId `json:"nodeId"`                // Id of the node to resolve.
	ObjectGroup string `json:"objectGroup,omitempty"` // Symbolic group name that can be used to release multiple objects.
}

type ResolveNodeResult struct {
	Object *RemoteObject `json:"object"` // JavaScript object wrapper for given node.
}

// Resolves JavaScript node object for given node id.

type ResolveNodeCommand struct {
	params *ResolveNodeParams
	result ResolveNodeResult
	wg     sync.WaitGroup
	err    error
}

func NewResolveNodeCommand(params *ResolveNodeParams) *ResolveNodeCommand {
	return &ResolveNodeCommand{
		params: params,
	}
}

func (cmd *ResolveNodeCommand) Name() string {
	return "DOM.resolveNode"
}

func (cmd *ResolveNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ResolveNodeCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ResolveNode(params *ResolveNodeParams, conn *hc.Conn) (result *ResolveNodeResult, err error) {
	cmd := NewResolveNodeCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type ResolveNodeCB func(result *ResolveNodeResult, err error)

// Resolves JavaScript node object for given node id.

type AsyncResolveNodeCommand struct {
	params *ResolveNodeParams
	cb     ResolveNodeCB
}

func NewAsyncResolveNodeCommand(params *ResolveNodeParams, cb ResolveNodeCB) *AsyncResolveNodeCommand {
	return &AsyncResolveNodeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncResolveNodeCommand) Name() string {
	return "DOM.resolveNode"
}

func (cmd *AsyncResolveNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ResolveNodeCommand) Result() *ResolveNodeResult {
	return &cmd.result
}

func (cmd *ResolveNodeCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncResolveNodeCommand) Done(data []byte, err error) {
	var result ResolveNodeResult
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

type GetAttributesParams struct {
	NodeId NodeId `json:"nodeId"` // Id of the node to retrieve attibutes for.
}

type GetAttributesResult struct {
	Attributes []string `json:"attributes"` // An interleaved array of node attribute names and values.
}

// Returns attributes for the specified node.

type GetAttributesCommand struct {
	params *GetAttributesParams
	result GetAttributesResult
	wg     sync.WaitGroup
	err    error
}

func NewGetAttributesCommand(params *GetAttributesParams) *GetAttributesCommand {
	return &GetAttributesCommand{
		params: params,
	}
}

func (cmd *GetAttributesCommand) Name() string {
	return "DOM.getAttributes"
}

func (cmd *GetAttributesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetAttributesCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetAttributes(params *GetAttributesParams, conn *hc.Conn) (result *GetAttributesResult, err error) {
	cmd := NewGetAttributesCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetAttributesCB func(result *GetAttributesResult, err error)

// Returns attributes for the specified node.

type AsyncGetAttributesCommand struct {
	params *GetAttributesParams
	cb     GetAttributesCB
}

func NewAsyncGetAttributesCommand(params *GetAttributesParams, cb GetAttributesCB) *AsyncGetAttributesCommand {
	return &AsyncGetAttributesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetAttributesCommand) Name() string {
	return "DOM.getAttributes"
}

func (cmd *AsyncGetAttributesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetAttributesCommand) Result() *GetAttributesResult {
	return &cmd.result
}

func (cmd *GetAttributesCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetAttributesCommand) Done(data []byte, err error) {
	var result GetAttributesResult
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

type CopyToParams struct {
	NodeId             NodeId `json:"nodeId"`                       // Id of the node to copy.
	TargetNodeId       NodeId `json:"targetNodeId"`                 // Id of the element to drop the copy into.
	InsertBeforeNodeId NodeId `json:"insertBeforeNodeId,omitempty"` // Drop the copy before this node (if absent, the copy becomes the last child of targetNodeId).
}

type CopyToResult struct {
	NodeId NodeId `json:"nodeId"` // Id of the node clone.
}

// Creates a deep copy of the specified node and places it into the target container before the given anchor.
// @experimental
type CopyToCommand struct {
	params *CopyToParams
	result CopyToResult
	wg     sync.WaitGroup
	err    error
}

func NewCopyToCommand(params *CopyToParams) *CopyToCommand {
	return &CopyToCommand{
		params: params,
	}
}

func (cmd *CopyToCommand) Name() string {
	return "DOM.copyTo"
}

func (cmd *CopyToCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CopyToCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func CopyTo(params *CopyToParams, conn *hc.Conn) (result *CopyToResult, err error) {
	cmd := NewCopyToCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type CopyToCB func(result *CopyToResult, err error)

// Creates a deep copy of the specified node and places it into the target container before the given anchor.
// @experimental
type AsyncCopyToCommand struct {
	params *CopyToParams
	cb     CopyToCB
}

func NewAsyncCopyToCommand(params *CopyToParams, cb CopyToCB) *AsyncCopyToCommand {
	return &AsyncCopyToCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncCopyToCommand) Name() string {
	return "DOM.copyTo"
}

func (cmd *AsyncCopyToCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CopyToCommand) Result() *CopyToResult {
	return &cmd.result
}

func (cmd *CopyToCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncCopyToCommand) Done(data []byte, err error) {
	var result CopyToResult
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

type MoveToParams struct {
	NodeId             NodeId `json:"nodeId"`                       // Id of the node to move.
	TargetNodeId       NodeId `json:"targetNodeId"`                 // Id of the element to drop the moved node into.
	InsertBeforeNodeId NodeId `json:"insertBeforeNodeId,omitempty"` // Drop node before this one (if absent, the moved node becomes the last child of targetNodeId).
}

type MoveToResult struct {
	NodeId NodeId `json:"nodeId"` // New id of the moved node.
}

// Moves node into the new container, places it before the given anchor.

type MoveToCommand struct {
	params *MoveToParams
	result MoveToResult
	wg     sync.WaitGroup
	err    error
}

func NewMoveToCommand(params *MoveToParams) *MoveToCommand {
	return &MoveToCommand{
		params: params,
	}
}

func (cmd *MoveToCommand) Name() string {
	return "DOM.moveTo"
}

func (cmd *MoveToCommand) Params() interface{} {
	return cmd.params
}

func (cmd *MoveToCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func MoveTo(params *MoveToParams, conn *hc.Conn) (result *MoveToResult, err error) {
	cmd := NewMoveToCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type MoveToCB func(result *MoveToResult, err error)

// Moves node into the new container, places it before the given anchor.

type AsyncMoveToCommand struct {
	params *MoveToParams
	cb     MoveToCB
}

func NewAsyncMoveToCommand(params *MoveToParams, cb MoveToCB) *AsyncMoveToCommand {
	return &AsyncMoveToCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncMoveToCommand) Name() string {
	return "DOM.moveTo"
}

func (cmd *AsyncMoveToCommand) Params() interface{} {
	return cmd.params
}

func (cmd *MoveToCommand) Result() *MoveToResult {
	return &cmd.result
}

func (cmd *MoveToCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncMoveToCommand) Done(data []byte, err error) {
	var result MoveToResult
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

// Undoes the last performed action.
// @experimental
type UndoCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewUndoCommand() *UndoCommand {
	return &UndoCommand{}
}

func (cmd *UndoCommand) Name() string {
	return "DOM.undo"
}

func (cmd *UndoCommand) Params() interface{} {
	return nil
}

func (cmd *UndoCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func Undo(conn *hc.Conn) (err error) {
	cmd := NewUndoCommand()
	cmd.Run(conn)
	return cmd.err
}

type UndoCB func(err error)

// Undoes the last performed action.
// @experimental
type AsyncUndoCommand struct {
	cb UndoCB
}

func NewAsyncUndoCommand(cb UndoCB) *AsyncUndoCommand {
	return &AsyncUndoCommand{
		cb: cb,
	}
}

func (cmd *AsyncUndoCommand) Name() string {
	return "DOM.undo"
}

func (cmd *AsyncUndoCommand) Params() interface{} {
	return nil
}

func (cmd *UndoCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncUndoCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Re-does the last undone action.
// @experimental
type RedoCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewRedoCommand() *RedoCommand {
	return &RedoCommand{}
}

func (cmd *RedoCommand) Name() string {
	return "DOM.redo"
}

func (cmd *RedoCommand) Params() interface{} {
	return nil
}

func (cmd *RedoCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func Redo(conn *hc.Conn) (err error) {
	cmd := NewRedoCommand()
	cmd.Run(conn)
	return cmd.err
}

type RedoCB func(err error)

// Re-does the last undone action.
// @experimental
type AsyncRedoCommand struct {
	cb RedoCB
}

func NewAsyncRedoCommand(cb RedoCB) *AsyncRedoCommand {
	return &AsyncRedoCommand{
		cb: cb,
	}
}

func (cmd *AsyncRedoCommand) Name() string {
	return "DOM.redo"
}

func (cmd *AsyncRedoCommand) Params() interface{} {
	return nil
}

func (cmd *RedoCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncRedoCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Marks last undoable state.
// @experimental
type MarkUndoableStateCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewMarkUndoableStateCommand() *MarkUndoableStateCommand {
	return &MarkUndoableStateCommand{}
}

func (cmd *MarkUndoableStateCommand) Name() string {
	return "DOM.markUndoableState"
}

func (cmd *MarkUndoableStateCommand) Params() interface{} {
	return nil
}

func (cmd *MarkUndoableStateCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func MarkUndoableState(conn *hc.Conn) (err error) {
	cmd := NewMarkUndoableStateCommand()
	cmd.Run(conn)
	return cmd.err
}

type MarkUndoableStateCB func(err error)

// Marks last undoable state.
// @experimental
type AsyncMarkUndoableStateCommand struct {
	cb MarkUndoableStateCB
}

func NewAsyncMarkUndoableStateCommand(cb MarkUndoableStateCB) *AsyncMarkUndoableStateCommand {
	return &AsyncMarkUndoableStateCommand{
		cb: cb,
	}
}

func (cmd *AsyncMarkUndoableStateCommand) Name() string {
	return "DOM.markUndoableState"
}

func (cmd *AsyncMarkUndoableStateCommand) Params() interface{} {
	return nil
}

func (cmd *MarkUndoableStateCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncMarkUndoableStateCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type FocusParams struct {
	NodeId NodeId `json:"nodeId"` // Id of the node to focus.
}

// Focuses the given element.
// @experimental
type FocusCommand struct {
	params *FocusParams
	wg     sync.WaitGroup
	err    error
}

func NewFocusCommand(params *FocusParams) *FocusCommand {
	return &FocusCommand{
		params: params,
	}
}

func (cmd *FocusCommand) Name() string {
	return "DOM.focus"
}

func (cmd *FocusCommand) Params() interface{} {
	return cmd.params
}

func (cmd *FocusCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func Focus(params *FocusParams, conn *hc.Conn) (err error) {
	cmd := NewFocusCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type FocusCB func(err error)

// Focuses the given element.
// @experimental
type AsyncFocusCommand struct {
	params *FocusParams
	cb     FocusCB
}

func NewAsyncFocusCommand(params *FocusParams, cb FocusCB) *AsyncFocusCommand {
	return &AsyncFocusCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncFocusCommand) Name() string {
	return "DOM.focus"
}

func (cmd *AsyncFocusCommand) Params() interface{} {
	return cmd.params
}

func (cmd *FocusCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncFocusCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SetFileInputFilesParams struct {
	NodeId NodeId   `json:"nodeId"` // Id of the file input node to set files for.
	Files  []string `json:"files"`  // Array of file paths to set.
}

// Sets files for the given file input element.
// @experimental
type SetFileInputFilesCommand struct {
	params *SetFileInputFilesParams
	wg     sync.WaitGroup
	err    error
}

func NewSetFileInputFilesCommand(params *SetFileInputFilesParams) *SetFileInputFilesCommand {
	return &SetFileInputFilesCommand{
		params: params,
	}
}

func (cmd *SetFileInputFilesCommand) Name() string {
	return "DOM.setFileInputFiles"
}

func (cmd *SetFileInputFilesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetFileInputFilesCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetFileInputFiles(params *SetFileInputFilesParams, conn *hc.Conn) (err error) {
	cmd := NewSetFileInputFilesCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetFileInputFilesCB func(err error)

// Sets files for the given file input element.
// @experimental
type AsyncSetFileInputFilesCommand struct {
	params *SetFileInputFilesParams
	cb     SetFileInputFilesCB
}

func NewAsyncSetFileInputFilesCommand(params *SetFileInputFilesParams, cb SetFileInputFilesCB) *AsyncSetFileInputFilesCommand {
	return &AsyncSetFileInputFilesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetFileInputFilesCommand) Name() string {
	return "DOM.setFileInputFiles"
}

func (cmd *AsyncSetFileInputFilesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetFileInputFilesCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetFileInputFilesCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type GetBoxModelParams struct {
	NodeId NodeId `json:"nodeId"` // Id of the node to get box model for.
}

type GetBoxModelResult struct {
	Model *BoxModel `json:"model"` // Box model for the node.
}

// Returns boxes for the currently selected nodes.
// @experimental
type GetBoxModelCommand struct {
	params *GetBoxModelParams
	result GetBoxModelResult
	wg     sync.WaitGroup
	err    error
}

func NewGetBoxModelCommand(params *GetBoxModelParams) *GetBoxModelCommand {
	return &GetBoxModelCommand{
		params: params,
	}
}

func (cmd *GetBoxModelCommand) Name() string {
	return "DOM.getBoxModel"
}

func (cmd *GetBoxModelCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetBoxModelCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetBoxModel(params *GetBoxModelParams, conn *hc.Conn) (result *GetBoxModelResult, err error) {
	cmd := NewGetBoxModelCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetBoxModelCB func(result *GetBoxModelResult, err error)

// Returns boxes for the currently selected nodes.
// @experimental
type AsyncGetBoxModelCommand struct {
	params *GetBoxModelParams
	cb     GetBoxModelCB
}

func NewAsyncGetBoxModelCommand(params *GetBoxModelParams, cb GetBoxModelCB) *AsyncGetBoxModelCommand {
	return &AsyncGetBoxModelCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetBoxModelCommand) Name() string {
	return "DOM.getBoxModel"
}

func (cmd *AsyncGetBoxModelCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetBoxModelCommand) Result() *GetBoxModelResult {
	return &cmd.result
}

func (cmd *GetBoxModelCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetBoxModelCommand) Done(data []byte, err error) {
	var result GetBoxModelResult
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

type GetNodeForLocationParams struct {
	X int `json:"x"` // X coordinate.
	Y int `json:"y"` // Y coordinate.
}

type GetNodeForLocationResult struct {
	NodeId NodeId `json:"nodeId"` // Id of the node at given coordinates.
}

// Returns node id at given location.
// @experimental
type GetNodeForLocationCommand struct {
	params *GetNodeForLocationParams
	result GetNodeForLocationResult
	wg     sync.WaitGroup
	err    error
}

func NewGetNodeForLocationCommand(params *GetNodeForLocationParams) *GetNodeForLocationCommand {
	return &GetNodeForLocationCommand{
		params: params,
	}
}

func (cmd *GetNodeForLocationCommand) Name() string {
	return "DOM.getNodeForLocation"
}

func (cmd *GetNodeForLocationCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetNodeForLocationCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetNodeForLocation(params *GetNodeForLocationParams, conn *hc.Conn) (result *GetNodeForLocationResult, err error) {
	cmd := NewGetNodeForLocationCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetNodeForLocationCB func(result *GetNodeForLocationResult, err error)

// Returns node id at given location.
// @experimental
type AsyncGetNodeForLocationCommand struct {
	params *GetNodeForLocationParams
	cb     GetNodeForLocationCB
}

func NewAsyncGetNodeForLocationCommand(params *GetNodeForLocationParams, cb GetNodeForLocationCB) *AsyncGetNodeForLocationCommand {
	return &AsyncGetNodeForLocationCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetNodeForLocationCommand) Name() string {
	return "DOM.getNodeForLocation"
}

func (cmd *AsyncGetNodeForLocationCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetNodeForLocationCommand) Result() *GetNodeForLocationResult {
	return &cmd.result
}

func (cmd *GetNodeForLocationCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetNodeForLocationCommand) Done(data []byte, err error) {
	var result GetNodeForLocationResult
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

type GetRelayoutBoundaryParams struct {
	NodeId NodeId `json:"nodeId"` // Id of the node.
}

type GetRelayoutBoundaryResult struct {
	NodeId NodeId `json:"nodeId"` // Relayout boundary node id for the given node.
}

// Returns the id of the nearest ancestor that is a relayout boundary.
// @experimental
type GetRelayoutBoundaryCommand struct {
	params *GetRelayoutBoundaryParams
	result GetRelayoutBoundaryResult
	wg     sync.WaitGroup
	err    error
}

func NewGetRelayoutBoundaryCommand(params *GetRelayoutBoundaryParams) *GetRelayoutBoundaryCommand {
	return &GetRelayoutBoundaryCommand{
		params: params,
	}
}

func (cmd *GetRelayoutBoundaryCommand) Name() string {
	return "DOM.getRelayoutBoundary"
}

func (cmd *GetRelayoutBoundaryCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetRelayoutBoundaryCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetRelayoutBoundary(params *GetRelayoutBoundaryParams, conn *hc.Conn) (result *GetRelayoutBoundaryResult, err error) {
	cmd := NewGetRelayoutBoundaryCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetRelayoutBoundaryCB func(result *GetRelayoutBoundaryResult, err error)

// Returns the id of the nearest ancestor that is a relayout boundary.
// @experimental
type AsyncGetRelayoutBoundaryCommand struct {
	params *GetRelayoutBoundaryParams
	cb     GetRelayoutBoundaryCB
}

func NewAsyncGetRelayoutBoundaryCommand(params *GetRelayoutBoundaryParams, cb GetRelayoutBoundaryCB) *AsyncGetRelayoutBoundaryCommand {
	return &AsyncGetRelayoutBoundaryCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetRelayoutBoundaryCommand) Name() string {
	return "DOM.getRelayoutBoundary"
}

func (cmd *AsyncGetRelayoutBoundaryCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetRelayoutBoundaryCommand) Result() *GetRelayoutBoundaryResult {
	return &cmd.result
}

func (cmd *GetRelayoutBoundaryCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetRelayoutBoundaryCommand) Done(data []byte, err error) {
	var result GetRelayoutBoundaryResult
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

type GetHighlightObjectForTestParams struct {
	NodeId NodeId `json:"nodeId"` // Id of the node to get highlight object for.
}

type GetHighlightObjectForTestResult struct {
	Highlight map[string]string `json:"highlight"` // Highlight data for the node.
}

// For testing.
// @experimental
type GetHighlightObjectForTestCommand struct {
	params *GetHighlightObjectForTestParams
	result GetHighlightObjectForTestResult
	wg     sync.WaitGroup
	err    error
}

func NewGetHighlightObjectForTestCommand(params *GetHighlightObjectForTestParams) *GetHighlightObjectForTestCommand {
	return &GetHighlightObjectForTestCommand{
		params: params,
	}
}

func (cmd *GetHighlightObjectForTestCommand) Name() string {
	return "DOM.getHighlightObjectForTest"
}

func (cmd *GetHighlightObjectForTestCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetHighlightObjectForTestCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetHighlightObjectForTest(params *GetHighlightObjectForTestParams, conn *hc.Conn) (result *GetHighlightObjectForTestResult, err error) {
	cmd := NewGetHighlightObjectForTestCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetHighlightObjectForTestCB func(result *GetHighlightObjectForTestResult, err error)

// For testing.
// @experimental
type AsyncGetHighlightObjectForTestCommand struct {
	params *GetHighlightObjectForTestParams
	cb     GetHighlightObjectForTestCB
}

func NewAsyncGetHighlightObjectForTestCommand(params *GetHighlightObjectForTestParams, cb GetHighlightObjectForTestCB) *AsyncGetHighlightObjectForTestCommand {
	return &AsyncGetHighlightObjectForTestCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetHighlightObjectForTestCommand) Name() string {
	return "DOM.getHighlightObjectForTest"
}

func (cmd *AsyncGetHighlightObjectForTestCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetHighlightObjectForTestCommand) Result() *GetHighlightObjectForTestResult {
	return &cmd.result
}

func (cmd *GetHighlightObjectForTestCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetHighlightObjectForTestCommand) Done(data []byte, err error) {
	var result GetHighlightObjectForTestResult
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

// Fired when Document has been totally updated. Node ids are no longer valid.

type DocumentUpdatedEvent struct {
}

func OnDocumentUpdated(conn *hc.Conn, cb func(evt *DocumentUpdatedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &DocumentUpdatedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("DOM.documentUpdated", sink)
}

// Fired when the node should be inspected. This happens after call to setInspectMode.
// @experimental
type InspectNodeRequestedEvent struct {
	BackendNodeId BackendNodeId `json:"backendNodeId"` // Id of the node to inspect.
}

func OnInspectNodeRequested(conn *hc.Conn, cb func(evt *InspectNodeRequestedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &InspectNodeRequestedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("DOM.inspectNodeRequested", sink)
}

// Fired when backend wants to provide client with the missing DOM structure. This happens upon most of the calls requesting node ids.

type SetChildNodesEvent struct {
	ParentId NodeId  `json:"parentId"` // Parent node id to populate with children.
	Nodes    []*Node `json:"nodes"`    // Child nodes array.
}

func OnSetChildNodes(conn *hc.Conn, cb func(evt *SetChildNodesEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &SetChildNodesEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("DOM.setChildNodes", sink)
}

// Fired when Element's attribute is modified.

type AttributeModifiedEvent struct {
	NodeId NodeId `json:"nodeId"` // Id of the node that has changed.
	Name   string `json:"name"`   // Attribute name.
	Value  string `json:"value"`  // Attribute value.
}

func OnAttributeModified(conn *hc.Conn, cb func(evt *AttributeModifiedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &AttributeModifiedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("DOM.attributeModified", sink)
}

// Fired when Element's attribute is removed.

type AttributeRemovedEvent struct {
	NodeId NodeId `json:"nodeId"` // Id of the node that has changed.
	Name   string `json:"name"`   // A ttribute name.
}

func OnAttributeRemoved(conn *hc.Conn, cb func(evt *AttributeRemovedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &AttributeRemovedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("DOM.attributeRemoved", sink)
}

// Fired when Element's inline style is modified via a CSS property modification.
// @experimental
type InlineStyleInvalidatedEvent struct {
	NodeIds []NodeId `json:"nodeIds"` // Ids of the nodes for which the inline styles have been invalidated.
}

func OnInlineStyleInvalidated(conn *hc.Conn, cb func(evt *InlineStyleInvalidatedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &InlineStyleInvalidatedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("DOM.inlineStyleInvalidated", sink)
}

// Mirrors DOMCharacterDataModified event.

type CharacterDataModifiedEvent struct {
	NodeId        NodeId `json:"nodeId"`        // Id of the node that has changed.
	CharacterData string `json:"characterData"` // New text value.
}

func OnCharacterDataModified(conn *hc.Conn, cb func(evt *CharacterDataModifiedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &CharacterDataModifiedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("DOM.characterDataModified", sink)
}

// Fired when Container's child node count has changed.

type ChildNodeCountUpdatedEvent struct {
	NodeId         NodeId `json:"nodeId"`         // Id of the node that has changed.
	ChildNodeCount int    `json:"childNodeCount"` // New node count.
}

func OnChildNodeCountUpdated(conn *hc.Conn, cb func(evt *ChildNodeCountUpdatedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &ChildNodeCountUpdatedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("DOM.childNodeCountUpdated", sink)
}

// Mirrors DOMNodeInserted event.

type ChildNodeInsertedEvent struct {
	ParentNodeId   NodeId `json:"parentNodeId"`   // Id of the node that has changed.
	PreviousNodeId NodeId `json:"previousNodeId"` // If of the previous siblint.
	Node           *Node  `json:"node"`           // Inserted node data.
}

func OnChildNodeInserted(conn *hc.Conn, cb func(evt *ChildNodeInsertedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &ChildNodeInsertedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("DOM.childNodeInserted", sink)
}

// Mirrors DOMNodeRemoved event.

type ChildNodeRemovedEvent struct {
	ParentNodeId NodeId `json:"parentNodeId"` // Parent id.
	NodeId       NodeId `json:"nodeId"`       // Id of the node that has been removed.
}

func OnChildNodeRemoved(conn *hc.Conn, cb func(evt *ChildNodeRemovedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &ChildNodeRemovedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("DOM.childNodeRemoved", sink)
}

// Called when shadow root is pushed into the element.
// @experimental
type ShadowRootPushedEvent struct {
	HostId NodeId `json:"hostId"` // Host element id.
	Root   *Node  `json:"root"`   // Shadow root.
}

func OnShadowRootPushed(conn *hc.Conn, cb func(evt *ShadowRootPushedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &ShadowRootPushedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("DOM.shadowRootPushed", sink)
}

// Called when shadow root is popped from the element.
// @experimental
type ShadowRootPoppedEvent struct {
	HostId NodeId `json:"hostId"` // Host element id.
	RootId NodeId `json:"rootId"` // Shadow root id.
}

func OnShadowRootPopped(conn *hc.Conn, cb func(evt *ShadowRootPoppedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &ShadowRootPoppedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("DOM.shadowRootPopped", sink)
}

// Called when a pseudo element is added to an element.
// @experimental
type PseudoElementAddedEvent struct {
	ParentId      NodeId `json:"parentId"`      // Pseudo element's parent element id.
	PseudoElement *Node  `json:"pseudoElement"` // The added pseudo element.
}

func OnPseudoElementAdded(conn *hc.Conn, cb func(evt *PseudoElementAddedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &PseudoElementAddedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("DOM.pseudoElementAdded", sink)
}

// Called when a pseudo element is removed from an element.
// @experimental
type PseudoElementRemovedEvent struct {
	ParentId        NodeId `json:"parentId"`        // Pseudo element's parent element id.
	PseudoElementId NodeId `json:"pseudoElementId"` // The removed pseudo element id.
}

func OnPseudoElementRemoved(conn *hc.Conn, cb func(evt *PseudoElementRemovedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &PseudoElementRemovedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("DOM.pseudoElementRemoved", sink)
}

// Called when distrubution is changed.
// @experimental
type DistributedNodesUpdatedEvent struct {
	InsertionPointId NodeId         `json:"insertionPointId"` // Insertion point where distrubuted nodes were updated.
	DistributedNodes []*BackendNode `json:"distributedNodes"` // Distributed nodes for given insertion point.
}

func OnDistributedNodesUpdated(conn *hc.Conn, cb func(evt *DistributedNodesUpdatedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &DistributedNodesUpdatedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("DOM.distributedNodesUpdated", sink)
}

// @experimental
type NodeHighlightRequestedEvent struct {
	NodeId NodeId `json:"nodeId"`
}

func OnNodeHighlightRequested(conn *hc.Conn, cb func(evt *NodeHighlightRequestedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &NodeHighlightRequestedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("DOM.nodeHighlightRequested", sink)
}
