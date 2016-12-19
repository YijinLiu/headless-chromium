package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
)

// Unique DOM node identifier.
type NodeId int

// Unique DOM node identifier used to reference a node that may not have been pushed to the front-end.
type BackendNodeId int

// Backend node with a friendly name.
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
	NodeId           NodeId         `json:"nodeId"`           // Node identifier that is passed into the rest of the DOM messages as the nodeId. Backend will only push node with given id once. It is aware of all requested nodes and will only fire DOM events for nodes known to the client.
	NodeType         int            `json:"nodeType"`         // Node's nodeType.
	NodeName         string         `json:"nodeName"`         // Node's nodeName.
	LocalName        string         `json:"localName"`        // Node's localName.
	NodeValue        string         `json:"nodeValue"`        // Node's nodeValue.
	ChildNodeCount   int            `json:"childNodeCount"`   // Child count for Container nodes.
	Children         []*Node        `json:"children"`         // Child nodes of this node when requested with children.
	Attributes       []string       `json:"attributes"`       // Attributes of the Element node in the form of flat array [name1, value1, name2, value2].
	DocumentURL      string         `json:"documentURL"`      // Document URL that Document or FrameOwner node points to.
	BaseURL          string         `json:"baseURL"`          // Base URL that Document or FrameOwner node uses for URL completion.
	PublicId         string         `json:"publicId"`         // DocumentType's publicId.
	SystemId         string         `json:"systemId"`         // DocumentType's systemId.
	InternalSubset   string         `json:"internalSubset"`   // DocumentType's internalSubset.
	XmlVersion       string         `json:"xmlVersion"`       // Document's XML version in case of XML documents.
	Name             string         `json:"name"`             // Attr's name.
	Value            string         `json:"value"`            // Attr's value.
	PseudoType       PseudoType     `json:"pseudoType"`       // Pseudo element type for this node.
	ShadowRootType   ShadowRootType `json:"shadowRootType"`   // Shadow root type.
	FrameId          *FrameId       `json:"frameId"`          // Frame ID for frame owner elements.
	ContentDocument  *Node          `json:"contentDocument"`  // Content document for frame owner elements.
	ShadowRoots      []*Node        `json:"shadowRoots"`      // Shadow root list for given element host.
	TemplateContent  *Node          `json:"templateContent"`  // Content document fragment for template elements.
	PseudoElements   []*Node        `json:"pseudoElements"`   // Pseudo elements associated with this node.
	ImportedDocument *Node          `json:"importedDocument"` // Import document for the HTMLImport links.
	DistributedNodes []*BackendNode `json:"distributedNodes"` // Distributed nodes for given insertion point.
}

// Details of post layout rendered text positions. The exact layout should not be regarded as stable and may change between versions.
type InlineTextBox struct {
	BoundingBox         *Rect `json:"boundingBox"`         // The absolute position bounding box.
	StartCharacterIndex int   `json:"startCharacterIndex"` // The starting index in characters, for this post layout textbox substring.
	NumCharacters       int   `json:"numCharacters"`       // The number of characters in this post layout textbox substring.
}

// Details of an element in the DOM tree with a LayoutObject.
type LayoutTreeNode struct {
	BackendNodeId   BackendNodeId    `json:"backendNodeId"`   // The BackendNodeId of the related DOM node.
	BoundingBox     *Rect            `json:"boundingBox"`     // The absolute position bounding box.
	LayoutText      string           `json:"layoutText"`      // Contents of the LayoutText if any
	InlineTextNodes []*InlineTextBox `json:"inlineTextNodes"` // The post layout inline text nodes, if any.
}

// A structure holding an RGBA color.
type RGBA struct {
	R int `json:"r"` // The red component, in the [0-255] range.
	G int `json:"g"` // The green component, in the [0-255] range.
	B int `json:"b"` // The blue component, in the [0-255] range.
	A int `json:"a"` // The alpha component, in the [0-1] range (default: 1).
}

// An array of quad vertices, x immediately followed by y for each point, points clock-wise.
type Quad []int

// Box model.
type BoxModel struct {
	Content      Quad              `json:"content"`      // Content box
	Padding      Quad              `json:"padding"`      // Padding box
	Border       Quad              `json:"border"`       // Border box
	Margin       Quad              `json:"margin"`       // Margin box
	Width        int               `json:"width"`        // Node width
	Height       int               `json:"height"`       // Node height
	ShapeOutside *ShapeOutsideInfo `json:"shapeOutside"` // Shape outside coordinates
}

// CSS Shape Outside details.
type ShapeOutsideInfo struct {
	Bounds      Quad     `json:"bounds"`      // Shape bounds
	Shape       []string `json:"shape"`       // Shape coordinate details
	MarginShape []string `json:"marginShape"` // Margin shape bounds
}

// Rectangle.
type Rect struct {
	X      int `json:"x"`      // X coordinate
	Y      int `json:"y"`      // Y coordinate
	Width  int `json:"width"`  // Rectangle width
	Height int `json:"height"` // Rectangle height
}

// Configuration data for the highlighting of page elements.
type HighlightConfig struct {
	ShowInfo           bool   `json:"showInfo"`           // Whether the node info tooltip should be shown (default: false).
	ShowRulers         bool   `json:"showRulers"`         // Whether the rulers should be shown (default: false).
	ShowExtensionLines bool   `json:"showExtensionLines"` // Whether the extension lines from node to the rulers should be shown (default: false).
	DisplayAsMaterial  bool   `json:"displayAsMaterial"`
	ContentColor       *RGBA  `json:"contentColor"`     // The content box highlight fill color (default: transparent).
	PaddingColor       *RGBA  `json:"paddingColor"`     // The padding highlight fill color (default: transparent).
	BorderColor        *RGBA  `json:"borderColor"`      // The border highlight fill color (default: transparent).
	MarginColor        *RGBA  `json:"marginColor"`      // The margin highlight fill color (default: transparent).
	EventTargetColor   *RGBA  `json:"eventTargetColor"` // The event target element highlight fill color (default: transparent).
	ShapeColor         *RGBA  `json:"shapeColor"`       // The shape outside fill color (default: transparent).
	ShapeMarginColor   *RGBA  `json:"shapeMarginColor"` // The shape margin fill color (default: transparent).
	SelectorList       string `json:"selectorList"`     // Selectors to highlight relevant nodes.
}

type InspectMode string

const InspectModeSearchForNode InspectMode = "searchForNode"
const InspectModeSearchForUAShadowDOM InspectMode = "searchForUAShadowDOM"
const InspectModeShowLayoutEditor InspectMode = "showLayoutEditor"
const InspectModeNone InspectMode = "none"

type DOMEnableCB func(err error)

// Enables DOM agent for the given page.
type DOMEnableCommand struct {
	cb DOMEnableCB
}

func NewDOMEnableCommand(cb DOMEnableCB) *DOMEnableCommand {
	return &DOMEnableCommand{
		cb: cb,
	}
}

func (cmd *DOMEnableCommand) Name() string {
	return "DOM.enable"
}

func (cmd *DOMEnableCommand) Params() interface{} {
	return nil
}

func (cmd *DOMEnableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type DOMDisableCB func(err error)

// Disables DOM agent for the given page.
type DOMDisableCommand struct {
	cb DOMDisableCB
}

func NewDOMDisableCommand(cb DOMDisableCB) *DOMDisableCommand {
	return &DOMDisableCommand{
		cb: cb,
	}
}

func (cmd *DOMDisableCommand) Name() string {
	return "DOM.disable"
}

func (cmd *DOMDisableCommand) Params() interface{} {
	return nil
}

func (cmd *DOMDisableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type GetDocumentResult struct {
	Root *Node `json:"root"` // Resulting node.
}

type GetDocumentCB func(result *GetDocumentResult, err error)

// Returns the root DOM node to the caller.
type GetDocumentCommand struct {
	cb GetDocumentCB
}

func NewGetDocumentCommand(cb GetDocumentCB) *GetDocumentCommand {
	return &GetDocumentCommand{
		cb: cb,
	}
}

func (cmd *GetDocumentCommand) Name() string {
	return "DOM.getDocument"
}

func (cmd *GetDocumentCommand) Params() interface{} {
	return nil
}

func (cmd *GetDocumentCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetDocumentResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type GetLayoutTreeNodesResult struct {
	LayoutTreeNodes []*LayoutTreeNode `json:"layoutTreeNodes"`
}

type GetLayoutTreeNodesCB func(result *GetLayoutTreeNodesResult, err error)

// Returns the document's LayoutTreeNodes to the caller, and those of any iframes too.
type GetLayoutTreeNodesCommand struct {
	cb GetLayoutTreeNodesCB
}

func NewGetLayoutTreeNodesCommand(cb GetLayoutTreeNodesCB) *GetLayoutTreeNodesCommand {
	return &GetLayoutTreeNodesCommand{
		cb: cb,
	}
}

func (cmd *GetLayoutTreeNodesCommand) Name() string {
	return "DOM.getLayoutTreeNodes"
}

func (cmd *GetLayoutTreeNodesCommand) Params() interface{} {
	return nil
}

func (cmd *GetLayoutTreeNodesCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetLayoutTreeNodesResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type CollectClassNamesFromSubtreeParams struct {
	NodeId NodeId `json:"nodeId"` // Id of the node to collect class names.
}

type CollectClassNamesFromSubtreeResult struct {
	ClassNames []string `json:"classNames"` // Class name list.
}

type CollectClassNamesFromSubtreeCB func(result *CollectClassNamesFromSubtreeResult, err error)

// Collects class names for the node with given id and all of it's child nodes.
type CollectClassNamesFromSubtreeCommand struct {
	params *CollectClassNamesFromSubtreeParams
	cb     CollectClassNamesFromSubtreeCB
}

func NewCollectClassNamesFromSubtreeCommand(params *CollectClassNamesFromSubtreeParams, cb CollectClassNamesFromSubtreeCB) *CollectClassNamesFromSubtreeCommand {
	return &CollectClassNamesFromSubtreeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *CollectClassNamesFromSubtreeCommand) Name() string {
	return "DOM.collectClassNamesFromSubtree"
}

func (cmd *CollectClassNamesFromSubtreeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CollectClassNamesFromSubtreeCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj CollectClassNamesFromSubtreeResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type RequestChildNodesParams struct {
	NodeId NodeId `json:"nodeId"` // Id of the node to get children for.
	Depth  int    `json:"depth"`  // The maximum depth at which children should be retrieved, defaults to 1. Use -1 for the entire subtree or provide an integer larger than 0.
}

type RequestChildNodesCB func(err error)

// Requests that children of the node with given id are returned to the caller in form of setChildNodes events where not only immediate children are retrieved, but all children down to the specified depth.
type RequestChildNodesCommand struct {
	params *RequestChildNodesParams
	cb     RequestChildNodesCB
}

func NewRequestChildNodesCommand(params *RequestChildNodesParams, cb RequestChildNodesCB) *RequestChildNodesCommand {
	return &RequestChildNodesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *RequestChildNodesCommand) Name() string {
	return "DOM.requestChildNodes"
}

func (cmd *RequestChildNodesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RequestChildNodesCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type QuerySelectorParams struct {
	NodeId   NodeId `json:"nodeId"`   // Id of the node to query upon.
	Selector string `json:"selector"` // Selector string.
}

type QuerySelectorResult struct {
	NodeId NodeId `json:"nodeId"` // Query selector result.
}

type QuerySelectorCB func(result *QuerySelectorResult, err error)

// Executes querySelector on a given node.
type QuerySelectorCommand struct {
	params *QuerySelectorParams
	cb     QuerySelectorCB
}

func NewQuerySelectorCommand(params *QuerySelectorParams, cb QuerySelectorCB) *QuerySelectorCommand {
	return &QuerySelectorCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *QuerySelectorCommand) Name() string {
	return "DOM.querySelector"
}

func (cmd *QuerySelectorCommand) Params() interface{} {
	return cmd.params
}

func (cmd *QuerySelectorCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj QuerySelectorResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type QuerySelectorAllParams struct {
	NodeId   NodeId `json:"nodeId"`   // Id of the node to query upon.
	Selector string `json:"selector"` // Selector string.
}

type QuerySelectorAllResult struct {
	NodeIds []NodeId `json:"nodeIds"` // Query selector result.
}

type QuerySelectorAllCB func(result *QuerySelectorAllResult, err error)

// Executes querySelectorAll on a given node.
type QuerySelectorAllCommand struct {
	params *QuerySelectorAllParams
	cb     QuerySelectorAllCB
}

func NewQuerySelectorAllCommand(params *QuerySelectorAllParams, cb QuerySelectorAllCB) *QuerySelectorAllCommand {
	return &QuerySelectorAllCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *QuerySelectorAllCommand) Name() string {
	return "DOM.querySelectorAll"
}

func (cmd *QuerySelectorAllCommand) Params() interface{} {
	return cmd.params
}

func (cmd *QuerySelectorAllCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj QuerySelectorAllResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type SetNodeNameParams struct {
	NodeId NodeId `json:"nodeId"` // Id of the node to set name for.
	Name   string `json:"name"`   // New node's name.
}

type SetNodeNameResult struct {
	NodeId NodeId `json:"nodeId"` // New node's id.
}

type SetNodeNameCB func(result *SetNodeNameResult, err error)

// Sets node name for a node with given id.
type SetNodeNameCommand struct {
	params *SetNodeNameParams
	cb     SetNodeNameCB
}

func NewSetNodeNameCommand(params *SetNodeNameParams, cb SetNodeNameCB) *SetNodeNameCommand {
	return &SetNodeNameCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetNodeNameCommand) Name() string {
	return "DOM.setNodeName"
}

func (cmd *SetNodeNameCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetNodeNameCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj SetNodeNameResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type SetNodeValueParams struct {
	NodeId NodeId `json:"nodeId"` // Id of the node to set value for.
	Value  string `json:"value"`  // New node's value.
}

type SetNodeValueCB func(err error)

// Sets node value for a node with given id.
type SetNodeValueCommand struct {
	params *SetNodeValueParams
	cb     SetNodeValueCB
}

func NewSetNodeValueCommand(params *SetNodeValueParams, cb SetNodeValueCB) *SetNodeValueCommand {
	return &SetNodeValueCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetNodeValueCommand) Name() string {
	return "DOM.setNodeValue"
}

func (cmd *SetNodeValueCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetNodeValueCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type RemoveNodeParams struct {
	NodeId NodeId `json:"nodeId"` // Id of the node to remove.
}

type RemoveNodeCB func(err error)

// Removes node with given id.
type RemoveNodeCommand struct {
	params *RemoveNodeParams
	cb     RemoveNodeCB
}

func NewRemoveNodeCommand(params *RemoveNodeParams, cb RemoveNodeCB) *RemoveNodeCommand {
	return &RemoveNodeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *RemoveNodeCommand) Name() string {
	return "DOM.removeNode"
}

func (cmd *RemoveNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveNodeCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetAttributeValueParams struct {
	NodeId NodeId `json:"nodeId"` // Id of the element to set attribute for.
	Name   string `json:"name"`   // Attribute name.
	Value  string `json:"value"`  // Attribute value.
}

type SetAttributeValueCB func(err error)

// Sets attribute for an element with given id.
type SetAttributeValueCommand struct {
	params *SetAttributeValueParams
	cb     SetAttributeValueCB
}

func NewSetAttributeValueCommand(params *SetAttributeValueParams, cb SetAttributeValueCB) *SetAttributeValueCommand {
	return &SetAttributeValueCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetAttributeValueCommand) Name() string {
	return "DOM.setAttributeValue"
}

func (cmd *SetAttributeValueCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetAttributeValueCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetAttributesAsTextParams struct {
	NodeId NodeId `json:"nodeId"` // Id of the element to set attributes for.
	Text   string `json:"text"`   // Text with a number of attributes. Will parse this text using HTML parser.
	Name   string `json:"name"`   // Attribute name to replace with new attributes derived from text in case text parsed successfully.
}

type SetAttributesAsTextCB func(err error)

// Sets attributes on element with given id. This method is useful when user edits some existing attribute value and types in several attribute name/value pairs.
type SetAttributesAsTextCommand struct {
	params *SetAttributesAsTextParams
	cb     SetAttributesAsTextCB
}

func NewSetAttributesAsTextCommand(params *SetAttributesAsTextParams, cb SetAttributesAsTextCB) *SetAttributesAsTextCommand {
	return &SetAttributesAsTextCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetAttributesAsTextCommand) Name() string {
	return "DOM.setAttributesAsText"
}

func (cmd *SetAttributesAsTextCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetAttributesAsTextCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type RemoveAttributeParams struct {
	NodeId NodeId `json:"nodeId"` // Id of the element to remove attribute from.
	Name   string `json:"name"`   // Name of the attribute to remove.
}

type RemoveAttributeCB func(err error)

// Removes attribute with given name from an element with given id.
type RemoveAttributeCommand struct {
	params *RemoveAttributeParams
	cb     RemoveAttributeCB
}

func NewRemoveAttributeCommand(params *RemoveAttributeParams, cb RemoveAttributeCB) *RemoveAttributeCommand {
	return &RemoveAttributeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *RemoveAttributeCommand) Name() string {
	return "DOM.removeAttribute"
}

func (cmd *RemoveAttributeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RemoveAttributeCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type GetOuterHTMLParams struct {
	NodeId NodeId `json:"nodeId"` // Id of the node to get markup for.
}

type GetOuterHTMLResult struct {
	OuterHTML string `json:"outerHTML"` // Outer HTML markup.
}

type GetOuterHTMLCB func(result *GetOuterHTMLResult, err error)

// Returns node's HTML markup.
type GetOuterHTMLCommand struct {
	params *GetOuterHTMLParams
	cb     GetOuterHTMLCB
}

func NewGetOuterHTMLCommand(params *GetOuterHTMLParams, cb GetOuterHTMLCB) *GetOuterHTMLCommand {
	return &GetOuterHTMLCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetOuterHTMLCommand) Name() string {
	return "DOM.getOuterHTML"
}

func (cmd *GetOuterHTMLCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetOuterHTMLCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetOuterHTMLResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type SetOuterHTMLParams struct {
	NodeId    NodeId `json:"nodeId"`    // Id of the node to set markup for.
	OuterHTML string `json:"outerHTML"` // Outer HTML markup to set.
}

type SetOuterHTMLCB func(err error)

// Sets node HTML markup, returns new node id.
type SetOuterHTMLCommand struct {
	params *SetOuterHTMLParams
	cb     SetOuterHTMLCB
}

func NewSetOuterHTMLCommand(params *SetOuterHTMLParams, cb SetOuterHTMLCB) *SetOuterHTMLCommand {
	return &SetOuterHTMLCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetOuterHTMLCommand) Name() string {
	return "DOM.setOuterHTML"
}

func (cmd *SetOuterHTMLCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetOuterHTMLCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type PerformSearchParams struct {
	Query                     string `json:"query"`                     // Plain text or query selector or XPath search query.
	IncludeUserAgentShadowDOM bool   `json:"includeUserAgentShadowDOM"` // True to search in user agent shadow DOM.
}

type PerformSearchResult struct {
	SearchId    string `json:"searchId"`    // Unique search session identifier.
	ResultCount int    `json:"resultCount"` // Number of search results.
}

type PerformSearchCB func(result *PerformSearchResult, err error)

// Searches for a given string in the DOM tree. Use getSearchResults to access search results or cancelSearch to end this search session.
type PerformSearchCommand struct {
	params *PerformSearchParams
	cb     PerformSearchCB
}

func NewPerformSearchCommand(params *PerformSearchParams, cb PerformSearchCB) *PerformSearchCommand {
	return &PerformSearchCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *PerformSearchCommand) Name() string {
	return "DOM.performSearch"
}

func (cmd *PerformSearchCommand) Params() interface{} {
	return cmd.params
}

func (cmd *PerformSearchCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj PerformSearchResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
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

type GetSearchResultsCB func(result *GetSearchResultsResult, err error)

// Returns search results from given fromIndex to given toIndex from the sarch with the given identifier.
type GetSearchResultsCommand struct {
	params *GetSearchResultsParams
	cb     GetSearchResultsCB
}

func NewGetSearchResultsCommand(params *GetSearchResultsParams, cb GetSearchResultsCB) *GetSearchResultsCommand {
	return &GetSearchResultsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetSearchResultsCommand) Name() string {
	return "DOM.getSearchResults"
}

func (cmd *GetSearchResultsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetSearchResultsCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetSearchResultsResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type DiscardSearchResultsParams struct {
	SearchId string `json:"searchId"` // Unique search session identifier.
}

type DiscardSearchResultsCB func(err error)

// Discards search results from the session with the given id. getSearchResults should no longer be called for that search.
type DiscardSearchResultsCommand struct {
	params *DiscardSearchResultsParams
	cb     DiscardSearchResultsCB
}

func NewDiscardSearchResultsCommand(params *DiscardSearchResultsParams, cb DiscardSearchResultsCB) *DiscardSearchResultsCommand {
	return &DiscardSearchResultsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *DiscardSearchResultsCommand) Name() string {
	return "DOM.discardSearchResults"
}

func (cmd *DiscardSearchResultsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *DiscardSearchResultsCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type RequestNodeParams struct {
	ObjectId *RemoteObjectId `json:"objectId"` // JavaScript object id to convert into node.
}

type RequestNodeResult struct {
	NodeId NodeId `json:"nodeId"` // Node id for given object.
}

type RequestNodeCB func(result *RequestNodeResult, err error)

// Requests that the node is sent to the caller given the JavaScript node object reference. All nodes that form the path from the node to the root are also sent to the client as a series of setChildNodes notifications.
type RequestNodeCommand struct {
	params *RequestNodeParams
	cb     RequestNodeCB
}

func NewRequestNodeCommand(params *RequestNodeParams, cb RequestNodeCB) *RequestNodeCommand {
	return &RequestNodeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *RequestNodeCommand) Name() string {
	return "DOM.requestNode"
}

func (cmd *RequestNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RequestNodeCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj RequestNodeResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type SetInspectModeParams struct {
	Mode            InspectMode      `json:"mode"`            // Set an inspection mode.
	HighlightConfig *HighlightConfig `json:"highlightConfig"` // A descriptor for the highlight appearance of hovered-over nodes. May be omitted if enabled == false.
}

type SetInspectModeCB func(err error)

// Enters the 'inspect' mode. In this mode, elements that user is hovering over are highlighted. Backend then generates 'inspectNodeRequested' event upon element selection.
type SetInspectModeCommand struct {
	params *SetInspectModeParams
	cb     SetInspectModeCB
}

func NewSetInspectModeCommand(params *SetInspectModeParams, cb SetInspectModeCB) *SetInspectModeCommand {
	return &SetInspectModeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetInspectModeCommand) Name() string {
	return "DOM.setInspectMode"
}

func (cmd *SetInspectModeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetInspectModeCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type HighlightRectParams struct {
	X            int   `json:"x"`            // X coordinate
	Y            int   `json:"y"`            // Y coordinate
	Width        int   `json:"width"`        // Rectangle width
	Height       int   `json:"height"`       // Rectangle height
	Color        *RGBA `json:"color"`        // The highlight fill color (default: transparent).
	OutlineColor *RGBA `json:"outlineColor"` // The highlight outline color (default: transparent).
}

type HighlightRectCB func(err error)

// Highlights given rectangle. Coordinates are absolute with respect to the main frame viewport.
type HighlightRectCommand struct {
	params *HighlightRectParams
	cb     HighlightRectCB
}

func NewHighlightRectCommand(params *HighlightRectParams, cb HighlightRectCB) *HighlightRectCommand {
	return &HighlightRectCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *HighlightRectCommand) Name() string {
	return "DOM.highlightRect"
}

func (cmd *HighlightRectCommand) Params() interface{} {
	return cmd.params
}

func (cmd *HighlightRectCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type HighlightQuadParams struct {
	Quad         Quad  `json:"quad"`         // Quad to highlight
	Color        *RGBA `json:"color"`        // The highlight fill color (default: transparent).
	OutlineColor *RGBA `json:"outlineColor"` // The highlight outline color (default: transparent).
}

type HighlightQuadCB func(err error)

// Highlights given quad. Coordinates are absolute with respect to the main frame viewport.
type HighlightQuadCommand struct {
	params *HighlightQuadParams
	cb     HighlightQuadCB
}

func NewHighlightQuadCommand(params *HighlightQuadParams, cb HighlightQuadCB) *HighlightQuadCommand {
	return &HighlightQuadCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *HighlightQuadCommand) Name() string {
	return "DOM.highlightQuad"
}

func (cmd *HighlightQuadCommand) Params() interface{} {
	return cmd.params
}

func (cmd *HighlightQuadCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type HighlightNodeParams struct {
	HighlightConfig *HighlightConfig `json:"highlightConfig"` // A descriptor for the highlight appearance.
	NodeId          NodeId           `json:"nodeId"`          // Identifier of the node to highlight.
	BackendNodeId   BackendNodeId    `json:"backendNodeId"`   // Identifier of the backend node to highlight.
	ObjectId        *RemoteObjectId  `json:"objectId"`        // JavaScript object id of the node to be highlighted.
}

type HighlightNodeCB func(err error)

// Highlights DOM node with given id or with the given JavaScript object wrapper. Either nodeId or objectId must be specified.
type HighlightNodeCommand struct {
	params *HighlightNodeParams
	cb     HighlightNodeCB
}

func NewHighlightNodeCommand(params *HighlightNodeParams, cb HighlightNodeCB) *HighlightNodeCommand {
	return &HighlightNodeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *HighlightNodeCommand) Name() string {
	return "DOM.highlightNode"
}

func (cmd *HighlightNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *HighlightNodeCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type HideHighlightCB func(err error)

// Hides DOM node highlight.
type HideHighlightCommand struct {
	cb HideHighlightCB
}

func NewHideHighlightCommand(cb HideHighlightCB) *HideHighlightCommand {
	return &HideHighlightCommand{
		cb: cb,
	}
}

func (cmd *HideHighlightCommand) Name() string {
	return "DOM.hideHighlight"
}

func (cmd *HideHighlightCommand) Params() interface{} {
	return nil
}

func (cmd *HideHighlightCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type HighlightFrameParams struct {
	FrameId             *FrameId `json:"frameId"`             // Identifier of the frame to highlight.
	ContentColor        *RGBA    `json:"contentColor"`        // The content box highlight fill color (default: transparent).
	ContentOutlineColor *RGBA    `json:"contentOutlineColor"` // The content box highlight outline color (default: transparent).
}

type HighlightFrameCB func(err error)

// Highlights owner element of the frame with given id.
type HighlightFrameCommand struct {
	params *HighlightFrameParams
	cb     HighlightFrameCB
}

func NewHighlightFrameCommand(params *HighlightFrameParams, cb HighlightFrameCB) *HighlightFrameCommand {
	return &HighlightFrameCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *HighlightFrameCommand) Name() string {
	return "DOM.highlightFrame"
}

func (cmd *HighlightFrameCommand) Params() interface{} {
	return cmd.params
}

func (cmd *HighlightFrameCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type PushNodeByPathToFrontendParams struct {
	Path string `json:"path"` // Path to node in the proprietary format.
}

type PushNodeByPathToFrontendResult struct {
	NodeId NodeId `json:"nodeId"` // Id of the node for given path.
}

type PushNodeByPathToFrontendCB func(result *PushNodeByPathToFrontendResult, err error)

// Requests that the node is sent to the caller given its path. // FIXME, use XPath
type PushNodeByPathToFrontendCommand struct {
	params *PushNodeByPathToFrontendParams
	cb     PushNodeByPathToFrontendCB
}

func NewPushNodeByPathToFrontendCommand(params *PushNodeByPathToFrontendParams, cb PushNodeByPathToFrontendCB) *PushNodeByPathToFrontendCommand {
	return &PushNodeByPathToFrontendCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *PushNodeByPathToFrontendCommand) Name() string {
	return "DOM.pushNodeByPathToFrontend"
}

func (cmd *PushNodeByPathToFrontendCommand) Params() interface{} {
	return cmd.params
}

func (cmd *PushNodeByPathToFrontendCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj PushNodeByPathToFrontendResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type PushNodesByBackendIdsToFrontendParams struct {
	BackendNodeIds []BackendNodeId `json:"backendNodeIds"` // The array of backend node ids.
}

type PushNodesByBackendIdsToFrontendResult struct {
	NodeIds []NodeId `json:"nodeIds"` // The array of ids of pushed nodes that correspond to the backend ids specified in backendNodeIds.
}

type PushNodesByBackendIdsToFrontendCB func(result *PushNodesByBackendIdsToFrontendResult, err error)

// Requests that a batch of nodes is sent to the caller given their backend node ids.
type PushNodesByBackendIdsToFrontendCommand struct {
	params *PushNodesByBackendIdsToFrontendParams
	cb     PushNodesByBackendIdsToFrontendCB
}

func NewPushNodesByBackendIdsToFrontendCommand(params *PushNodesByBackendIdsToFrontendParams, cb PushNodesByBackendIdsToFrontendCB) *PushNodesByBackendIdsToFrontendCommand {
	return &PushNodesByBackendIdsToFrontendCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *PushNodesByBackendIdsToFrontendCommand) Name() string {
	return "DOM.pushNodesByBackendIdsToFrontend"
}

func (cmd *PushNodesByBackendIdsToFrontendCommand) Params() interface{} {
	return cmd.params
}

func (cmd *PushNodesByBackendIdsToFrontendCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj PushNodesByBackendIdsToFrontendResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type SetInspectedNodeParams struct {
	NodeId NodeId `json:"nodeId"` // DOM node id to be accessible by means of $x command line API.
}

type SetInspectedNodeCB func(err error)

// Enables console to refer to the node with given id via $x (see Command Line API for more details $x functions).
type SetInspectedNodeCommand struct {
	params *SetInspectedNodeParams
	cb     SetInspectedNodeCB
}

func NewSetInspectedNodeCommand(params *SetInspectedNodeParams, cb SetInspectedNodeCB) *SetInspectedNodeCommand {
	return &SetInspectedNodeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetInspectedNodeCommand) Name() string {
	return "DOM.setInspectedNode"
}

func (cmd *SetInspectedNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetInspectedNodeCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type ResolveNodeParams struct {
	NodeId      NodeId `json:"nodeId"`      // Id of the node to resolve.
	ObjectGroup string `json:"objectGroup"` // Symbolic group name that can be used to release multiple objects.
}

type ResolveNodeResult struct {
	Object *RemoteObject `json:"object"` // JavaScript object wrapper for given node.
}

type ResolveNodeCB func(result *ResolveNodeResult, err error)

// Resolves JavaScript node object for given node id.
type ResolveNodeCommand struct {
	params *ResolveNodeParams
	cb     ResolveNodeCB
}

func NewResolveNodeCommand(params *ResolveNodeParams, cb ResolveNodeCB) *ResolveNodeCommand {
	return &ResolveNodeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *ResolveNodeCommand) Name() string {
	return "DOM.resolveNode"
}

func (cmd *ResolveNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ResolveNodeCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj ResolveNodeResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type GetAttributesParams struct {
	NodeId NodeId `json:"nodeId"` // Id of the node to retrieve attibutes for.
}

type GetAttributesResult struct {
	Attributes []string `json:"attributes"` // An interleaved array of node attribute names and values.
}

type GetAttributesCB func(result *GetAttributesResult, err error)

// Returns attributes for the specified node.
type GetAttributesCommand struct {
	params *GetAttributesParams
	cb     GetAttributesCB
}

func NewGetAttributesCommand(params *GetAttributesParams, cb GetAttributesCB) *GetAttributesCommand {
	return &GetAttributesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetAttributesCommand) Name() string {
	return "DOM.getAttributes"
}

func (cmd *GetAttributesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetAttributesCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetAttributesResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type CopyToParams struct {
	NodeId             NodeId `json:"nodeId"`             // Id of the node to copy.
	TargetNodeId       NodeId `json:"targetNodeId"`       // Id of the element to drop the copy into.
	InsertBeforeNodeId NodeId `json:"insertBeforeNodeId"` // Drop the copy before this node (if absent, the copy becomes the last child of targetNodeId).
}

type CopyToResult struct {
	NodeId NodeId `json:"nodeId"` // Id of the node clone.
}

type CopyToCB func(result *CopyToResult, err error)

// Creates a deep copy of the specified node and places it into the target container before the given anchor.
type CopyToCommand struct {
	params *CopyToParams
	cb     CopyToCB
}

func NewCopyToCommand(params *CopyToParams, cb CopyToCB) *CopyToCommand {
	return &CopyToCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *CopyToCommand) Name() string {
	return "DOM.copyTo"
}

func (cmd *CopyToCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CopyToCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj CopyToResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type MoveToParams struct {
	NodeId             NodeId `json:"nodeId"`             // Id of the node to move.
	TargetNodeId       NodeId `json:"targetNodeId"`       // Id of the element to drop the moved node into.
	InsertBeforeNodeId NodeId `json:"insertBeforeNodeId"` // Drop node before this one (if absent, the moved node becomes the last child of targetNodeId).
}

type MoveToResult struct {
	NodeId NodeId `json:"nodeId"` // New id of the moved node.
}

type MoveToCB func(result *MoveToResult, err error)

// Moves node into the new container, places it before the given anchor.
type MoveToCommand struct {
	params *MoveToParams
	cb     MoveToCB
}

func NewMoveToCommand(params *MoveToParams, cb MoveToCB) *MoveToCommand {
	return &MoveToCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *MoveToCommand) Name() string {
	return "DOM.moveTo"
}

func (cmd *MoveToCommand) Params() interface{} {
	return cmd.params
}

func (cmd *MoveToCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj MoveToResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type UndoCB func(err error)

// Undoes the last performed action.
type UndoCommand struct {
	cb UndoCB
}

func NewUndoCommand(cb UndoCB) *UndoCommand {
	return &UndoCommand{
		cb: cb,
	}
}

func (cmd *UndoCommand) Name() string {
	return "DOM.undo"
}

func (cmd *UndoCommand) Params() interface{} {
	return nil
}

func (cmd *UndoCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type RedoCB func(err error)

// Re-does the last undone action.
type RedoCommand struct {
	cb RedoCB
}

func NewRedoCommand(cb RedoCB) *RedoCommand {
	return &RedoCommand{
		cb: cb,
	}
}

func (cmd *RedoCommand) Name() string {
	return "DOM.redo"
}

func (cmd *RedoCommand) Params() interface{} {
	return nil
}

func (cmd *RedoCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type MarkUndoableStateCB func(err error)

// Marks last undoable state.
type MarkUndoableStateCommand struct {
	cb MarkUndoableStateCB
}

func NewMarkUndoableStateCommand(cb MarkUndoableStateCB) *MarkUndoableStateCommand {
	return &MarkUndoableStateCommand{
		cb: cb,
	}
}

func (cmd *MarkUndoableStateCommand) Name() string {
	return "DOM.markUndoableState"
}

func (cmd *MarkUndoableStateCommand) Params() interface{} {
	return nil
}

func (cmd *MarkUndoableStateCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type FocusParams struct {
	NodeId NodeId `json:"nodeId"` // Id of the node to focus.
}

type FocusCB func(err error)

// Focuses the given element.
type FocusCommand struct {
	params *FocusParams
	cb     FocusCB
}

func NewFocusCommand(params *FocusParams, cb FocusCB) *FocusCommand {
	return &FocusCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *FocusCommand) Name() string {
	return "DOM.focus"
}

func (cmd *FocusCommand) Params() interface{} {
	return cmd.params
}

func (cmd *FocusCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetFileInputFilesParams struct {
	NodeId NodeId   `json:"nodeId"` // Id of the file input node to set files for.
	Files  []string `json:"files"`  // Array of file paths to set.
}

type SetFileInputFilesCB func(err error)

// Sets files for the given file input element.
type SetFileInputFilesCommand struct {
	params *SetFileInputFilesParams
	cb     SetFileInputFilesCB
}

func NewSetFileInputFilesCommand(params *SetFileInputFilesParams, cb SetFileInputFilesCB) *SetFileInputFilesCommand {
	return &SetFileInputFilesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetFileInputFilesCommand) Name() string {
	return "DOM.setFileInputFiles"
}

func (cmd *SetFileInputFilesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetFileInputFilesCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type GetBoxModelParams struct {
	NodeId NodeId `json:"nodeId"` // Id of the node to get box model for.
}

type GetBoxModelResult struct {
	Model *BoxModel `json:"model"` // Box model for the node.
}

type GetBoxModelCB func(result *GetBoxModelResult, err error)

// Returns boxes for the currently selected nodes.
type GetBoxModelCommand struct {
	params *GetBoxModelParams
	cb     GetBoxModelCB
}

func NewGetBoxModelCommand(params *GetBoxModelParams, cb GetBoxModelCB) *GetBoxModelCommand {
	return &GetBoxModelCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetBoxModelCommand) Name() string {
	return "DOM.getBoxModel"
}

func (cmd *GetBoxModelCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetBoxModelCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetBoxModelResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type GetNodeForLocationParams struct {
	X int `json:"x"` // X coordinate.
	Y int `json:"y"` // Y coordinate.
}

type GetNodeForLocationResult struct {
	NodeId NodeId `json:"nodeId"` // Id of the node at given coordinates.
}

type GetNodeForLocationCB func(result *GetNodeForLocationResult, err error)

// Returns node id at given location.
type GetNodeForLocationCommand struct {
	params *GetNodeForLocationParams
	cb     GetNodeForLocationCB
}

func NewGetNodeForLocationCommand(params *GetNodeForLocationParams, cb GetNodeForLocationCB) *GetNodeForLocationCommand {
	return &GetNodeForLocationCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetNodeForLocationCommand) Name() string {
	return "DOM.getNodeForLocation"
}

func (cmd *GetNodeForLocationCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetNodeForLocationCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetNodeForLocationResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type GetRelayoutBoundaryParams struct {
	NodeId NodeId `json:"nodeId"` // Id of the node.
}

type GetRelayoutBoundaryResult struct {
	NodeId NodeId `json:"nodeId"` // Relayout boundary node id for the given node.
}

type GetRelayoutBoundaryCB func(result *GetRelayoutBoundaryResult, err error)

// Returns the id of the nearest ancestor that is a relayout boundary.
type GetRelayoutBoundaryCommand struct {
	params *GetRelayoutBoundaryParams
	cb     GetRelayoutBoundaryCB
}

func NewGetRelayoutBoundaryCommand(params *GetRelayoutBoundaryParams, cb GetRelayoutBoundaryCB) *GetRelayoutBoundaryCommand {
	return &GetRelayoutBoundaryCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetRelayoutBoundaryCommand) Name() string {
	return "DOM.getRelayoutBoundary"
}

func (cmd *GetRelayoutBoundaryCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetRelayoutBoundaryCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetRelayoutBoundaryResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type GetHighlightObjectForTestParams struct {
	NodeId NodeId `json:"nodeId"` // Id of the node to get highlight object for.
}

type GetHighlightObjectForTestResult struct {
	Highlight map[string]string `json:"highlight"` // Highlight data for the node.
}

type GetHighlightObjectForTestCB func(result *GetHighlightObjectForTestResult, err error)

// For testing.
type GetHighlightObjectForTestCommand struct {
	params *GetHighlightObjectForTestParams
	cb     GetHighlightObjectForTestCB
}

func NewGetHighlightObjectForTestCommand(params *GetHighlightObjectForTestParams, cb GetHighlightObjectForTestCB) *GetHighlightObjectForTestCommand {
	return &GetHighlightObjectForTestCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetHighlightObjectForTestCommand) Name() string {
	return "DOM.getHighlightObjectForTest"
}

func (cmd *GetHighlightObjectForTestCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetHighlightObjectForTestCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetHighlightObjectForTestResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type DocumentUpdatedEvent struct {
}

// Fired when Document has been totally updated. Node ids are no longer valid.
type DocumentUpdatedEventSink struct {
	events chan *DocumentUpdatedEvent
}

func NewDocumentUpdatedEventSink(bufSize int) *DocumentUpdatedEventSink {
	return &DocumentUpdatedEventSink{
		events: make(chan *DocumentUpdatedEvent, bufSize),
	}
}

func (s *DocumentUpdatedEventSink) Name() string {
	return "DOM.documentUpdated"
}

func (s *DocumentUpdatedEventSink) OnEvent(params []byte) {
	evt := &DocumentUpdatedEvent{}
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

type InspectNodeRequestedEvent struct {
	BackendNodeId BackendNodeId `json:"backendNodeId"` // Id of the node to inspect.
}

// Fired when the node should be inspected. This happens after call to setInspectMode.
type InspectNodeRequestedEventSink struct {
	events chan *InspectNodeRequestedEvent
}

func NewInspectNodeRequestedEventSink(bufSize int) *InspectNodeRequestedEventSink {
	return &InspectNodeRequestedEventSink{
		events: make(chan *InspectNodeRequestedEvent, bufSize),
	}
}

func (s *InspectNodeRequestedEventSink) Name() string {
	return "DOM.inspectNodeRequested"
}

func (s *InspectNodeRequestedEventSink) OnEvent(params []byte) {
	evt := &InspectNodeRequestedEvent{}
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

type SetChildNodesEvent struct {
	ParentId NodeId  `json:"parentId"` // Parent node id to populate with children.
	Nodes    []*Node `json:"nodes"`    // Child nodes array.
}

// Fired when backend wants to provide client with the missing DOM structure. This happens upon most of the calls requesting node ids.
type SetChildNodesEventSink struct {
	events chan *SetChildNodesEvent
}

func NewSetChildNodesEventSink(bufSize int) *SetChildNodesEventSink {
	return &SetChildNodesEventSink{
		events: make(chan *SetChildNodesEvent, bufSize),
	}
}

func (s *SetChildNodesEventSink) Name() string {
	return "DOM.setChildNodes"
}

func (s *SetChildNodesEventSink) OnEvent(params []byte) {
	evt := &SetChildNodesEvent{}
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

type AttributeModifiedEvent struct {
	NodeId NodeId `json:"nodeId"` // Id of the node that has changed.
	Name   string `json:"name"`   // Attribute name.
	Value  string `json:"value"`  // Attribute value.
}

// Fired when Element's attribute is modified.
type AttributeModifiedEventSink struct {
	events chan *AttributeModifiedEvent
}

func NewAttributeModifiedEventSink(bufSize int) *AttributeModifiedEventSink {
	return &AttributeModifiedEventSink{
		events: make(chan *AttributeModifiedEvent, bufSize),
	}
}

func (s *AttributeModifiedEventSink) Name() string {
	return "DOM.attributeModified"
}

func (s *AttributeModifiedEventSink) OnEvent(params []byte) {
	evt := &AttributeModifiedEvent{}
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

type AttributeRemovedEvent struct {
	NodeId NodeId `json:"nodeId"` // Id of the node that has changed.
	Name   string `json:"name"`   // A ttribute name.
}

// Fired when Element's attribute is removed.
type AttributeRemovedEventSink struct {
	events chan *AttributeRemovedEvent
}

func NewAttributeRemovedEventSink(bufSize int) *AttributeRemovedEventSink {
	return &AttributeRemovedEventSink{
		events: make(chan *AttributeRemovedEvent, bufSize),
	}
}

func (s *AttributeRemovedEventSink) Name() string {
	return "DOM.attributeRemoved"
}

func (s *AttributeRemovedEventSink) OnEvent(params []byte) {
	evt := &AttributeRemovedEvent{}
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

type InlineStyleInvalidatedEvent struct {
	NodeIds []NodeId `json:"nodeIds"` // Ids of the nodes for which the inline styles have been invalidated.
}

// Fired when Element's inline style is modified via a CSS property modification.
type InlineStyleInvalidatedEventSink struct {
	events chan *InlineStyleInvalidatedEvent
}

func NewInlineStyleInvalidatedEventSink(bufSize int) *InlineStyleInvalidatedEventSink {
	return &InlineStyleInvalidatedEventSink{
		events: make(chan *InlineStyleInvalidatedEvent, bufSize),
	}
}

func (s *InlineStyleInvalidatedEventSink) Name() string {
	return "DOM.inlineStyleInvalidated"
}

func (s *InlineStyleInvalidatedEventSink) OnEvent(params []byte) {
	evt := &InlineStyleInvalidatedEvent{}
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

type CharacterDataModifiedEvent struct {
	NodeId        NodeId `json:"nodeId"`        // Id of the node that has changed.
	CharacterData string `json:"characterData"` // New text value.
}

// Mirrors DOMCharacterDataModified event.
type CharacterDataModifiedEventSink struct {
	events chan *CharacterDataModifiedEvent
}

func NewCharacterDataModifiedEventSink(bufSize int) *CharacterDataModifiedEventSink {
	return &CharacterDataModifiedEventSink{
		events: make(chan *CharacterDataModifiedEvent, bufSize),
	}
}

func (s *CharacterDataModifiedEventSink) Name() string {
	return "DOM.characterDataModified"
}

func (s *CharacterDataModifiedEventSink) OnEvent(params []byte) {
	evt := &CharacterDataModifiedEvent{}
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

type ChildNodeCountUpdatedEvent struct {
	NodeId         NodeId `json:"nodeId"`         // Id of the node that has changed.
	ChildNodeCount int    `json:"childNodeCount"` // New node count.
}

// Fired when Container's child node count has changed.
type ChildNodeCountUpdatedEventSink struct {
	events chan *ChildNodeCountUpdatedEvent
}

func NewChildNodeCountUpdatedEventSink(bufSize int) *ChildNodeCountUpdatedEventSink {
	return &ChildNodeCountUpdatedEventSink{
		events: make(chan *ChildNodeCountUpdatedEvent, bufSize),
	}
}

func (s *ChildNodeCountUpdatedEventSink) Name() string {
	return "DOM.childNodeCountUpdated"
}

func (s *ChildNodeCountUpdatedEventSink) OnEvent(params []byte) {
	evt := &ChildNodeCountUpdatedEvent{}
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

type ChildNodeInsertedEvent struct {
	ParentNodeId   NodeId `json:"parentNodeId"`   // Id of the node that has changed.
	PreviousNodeId NodeId `json:"previousNodeId"` // If of the previous siblint.
	Node           *Node  `json:"node"`           // Inserted node data.
}

// Mirrors DOMNodeInserted event.
type ChildNodeInsertedEventSink struct {
	events chan *ChildNodeInsertedEvent
}

func NewChildNodeInsertedEventSink(bufSize int) *ChildNodeInsertedEventSink {
	return &ChildNodeInsertedEventSink{
		events: make(chan *ChildNodeInsertedEvent, bufSize),
	}
}

func (s *ChildNodeInsertedEventSink) Name() string {
	return "DOM.childNodeInserted"
}

func (s *ChildNodeInsertedEventSink) OnEvent(params []byte) {
	evt := &ChildNodeInsertedEvent{}
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

type ChildNodeRemovedEvent struct {
	ParentNodeId NodeId `json:"parentNodeId"` // Parent id.
	NodeId       NodeId `json:"nodeId"`       // Id of the node that has been removed.
}

// Mirrors DOMNodeRemoved event.
type ChildNodeRemovedEventSink struct {
	events chan *ChildNodeRemovedEvent
}

func NewChildNodeRemovedEventSink(bufSize int) *ChildNodeRemovedEventSink {
	return &ChildNodeRemovedEventSink{
		events: make(chan *ChildNodeRemovedEvent, bufSize),
	}
}

func (s *ChildNodeRemovedEventSink) Name() string {
	return "DOM.childNodeRemoved"
}

func (s *ChildNodeRemovedEventSink) OnEvent(params []byte) {
	evt := &ChildNodeRemovedEvent{}
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

type ShadowRootPushedEvent struct {
	HostId NodeId `json:"hostId"` // Host element id.
	Root   *Node  `json:"root"`   // Shadow root.
}

// Called when shadow root is pushed into the element.
type ShadowRootPushedEventSink struct {
	events chan *ShadowRootPushedEvent
}

func NewShadowRootPushedEventSink(bufSize int) *ShadowRootPushedEventSink {
	return &ShadowRootPushedEventSink{
		events: make(chan *ShadowRootPushedEvent, bufSize),
	}
}

func (s *ShadowRootPushedEventSink) Name() string {
	return "DOM.shadowRootPushed"
}

func (s *ShadowRootPushedEventSink) OnEvent(params []byte) {
	evt := &ShadowRootPushedEvent{}
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

type ShadowRootPoppedEvent struct {
	HostId NodeId `json:"hostId"` // Host element id.
	RootId NodeId `json:"rootId"` // Shadow root id.
}

// Called when shadow root is popped from the element.
type ShadowRootPoppedEventSink struct {
	events chan *ShadowRootPoppedEvent
}

func NewShadowRootPoppedEventSink(bufSize int) *ShadowRootPoppedEventSink {
	return &ShadowRootPoppedEventSink{
		events: make(chan *ShadowRootPoppedEvent, bufSize),
	}
}

func (s *ShadowRootPoppedEventSink) Name() string {
	return "DOM.shadowRootPopped"
}

func (s *ShadowRootPoppedEventSink) OnEvent(params []byte) {
	evt := &ShadowRootPoppedEvent{}
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

type PseudoElementAddedEvent struct {
	ParentId      NodeId `json:"parentId"`      // Pseudo element's parent element id.
	PseudoElement *Node  `json:"pseudoElement"` // The added pseudo element.
}

// Called when a pseudo element is added to an element.
type PseudoElementAddedEventSink struct {
	events chan *PseudoElementAddedEvent
}

func NewPseudoElementAddedEventSink(bufSize int) *PseudoElementAddedEventSink {
	return &PseudoElementAddedEventSink{
		events: make(chan *PseudoElementAddedEvent, bufSize),
	}
}

func (s *PseudoElementAddedEventSink) Name() string {
	return "DOM.pseudoElementAdded"
}

func (s *PseudoElementAddedEventSink) OnEvent(params []byte) {
	evt := &PseudoElementAddedEvent{}
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

type PseudoElementRemovedEvent struct {
	ParentId        NodeId `json:"parentId"`        // Pseudo element's parent element id.
	PseudoElementId NodeId `json:"pseudoElementId"` // The removed pseudo element id.
}

// Called when a pseudo element is removed from an element.
type PseudoElementRemovedEventSink struct {
	events chan *PseudoElementRemovedEvent
}

func NewPseudoElementRemovedEventSink(bufSize int) *PseudoElementRemovedEventSink {
	return &PseudoElementRemovedEventSink{
		events: make(chan *PseudoElementRemovedEvent, bufSize),
	}
}

func (s *PseudoElementRemovedEventSink) Name() string {
	return "DOM.pseudoElementRemoved"
}

func (s *PseudoElementRemovedEventSink) OnEvent(params []byte) {
	evt := &PseudoElementRemovedEvent{}
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

type DistributedNodesUpdatedEvent struct {
	InsertionPointId NodeId         `json:"insertionPointId"` // Insertion point where distrubuted nodes were updated.
	DistributedNodes []*BackendNode `json:"distributedNodes"` // Distributed nodes for given insertion point.
}

// Called when distrubution is changed.
type DistributedNodesUpdatedEventSink struct {
	events chan *DistributedNodesUpdatedEvent
}

func NewDistributedNodesUpdatedEventSink(bufSize int) *DistributedNodesUpdatedEventSink {
	return &DistributedNodesUpdatedEventSink{
		events: make(chan *DistributedNodesUpdatedEvent, bufSize),
	}
}

func (s *DistributedNodesUpdatedEventSink) Name() string {
	return "DOM.distributedNodesUpdated"
}

func (s *DistributedNodesUpdatedEventSink) OnEvent(params []byte) {
	evt := &DistributedNodesUpdatedEvent{}
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

type NodeHighlightRequestedEvent struct {
	NodeId NodeId `json:"nodeId"`
}

type NodeHighlightRequestedEventSink struct {
	events chan *NodeHighlightRequestedEvent
}

func NewNodeHighlightRequestedEventSink(bufSize int) *NodeHighlightRequestedEventSink {
	return &NodeHighlightRequestedEventSink{
		events: make(chan *NodeHighlightRequestedEvent, bufSize),
	}
}

func (s *NodeHighlightRequestedEventSink) Name() string {
	return "DOM.nodeHighlightRequested"
}

func (s *NodeHighlightRequestedEventSink) OnEvent(params []byte) {
	evt := &NodeHighlightRequestedEvent{}
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
