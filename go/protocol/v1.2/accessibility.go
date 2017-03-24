package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

// Unique accessibility node identifier.
type AXNodeId string

// Enum of possible property types.
type AXValueType string

const AXValueTypeBoolean AXValueType = "boolean"
const AXValueTypeTristate AXValueType = "tristate"
const AXValueTypeBooleanOrUndefined AXValueType = "booleanOrUndefined"
const AXValueTypeIdref AXValueType = "idref"
const AXValueTypeIdrefList AXValueType = "idrefList"
const AXValueTypeInteger AXValueType = "integer"
const AXValueTypeNode AXValueType = "node"
const AXValueTypeNodeList AXValueType = "nodeList"
const AXValueTypeNumber AXValueType = "number"
const AXValueTypeString AXValueType = "string"
const AXValueTypeComputedString AXValueType = "computedString"
const AXValueTypeToken AXValueType = "token"
const AXValueTypeTokenList AXValueType = "tokenList"
const AXValueTypeDomRelation AXValueType = "domRelation"
const AXValueTypeRole AXValueType = "role"
const AXValueTypeInternalRole AXValueType = "internalRole"
const AXValueTypeValueUndefined AXValueType = "valueUndefined"

// Enum of possible property sources.
type AXValueSourceType string

const AXValueSourceTypeAttribute AXValueSourceType = "attribute"
const AXValueSourceTypeImplicit AXValueSourceType = "implicit"
const AXValueSourceTypeStyle AXValueSourceType = "style"
const AXValueSourceTypeContents AXValueSourceType = "contents"
const AXValueSourceTypePlaceholder AXValueSourceType = "placeholder"
const AXValueSourceTypeRelatedElement AXValueSourceType = "relatedElement"

// Enum of possible native property sources (as a subtype of a particular AXValueSourceType).
type AXValueNativeSourceType string

const AXValueNativeSourceTypeFigcaption AXValueNativeSourceType = "figcaption"
const AXValueNativeSourceTypeLabel AXValueNativeSourceType = "label"
const AXValueNativeSourceTypeLabelfor AXValueNativeSourceType = "labelfor"
const AXValueNativeSourceTypeLabelwrapped AXValueNativeSourceType = "labelwrapped"
const AXValueNativeSourceTypeLegend AXValueNativeSourceType = "legend"
const AXValueNativeSourceTypeTablecaption AXValueNativeSourceType = "tablecaption"
const AXValueNativeSourceTypeTitle AXValueNativeSourceType = "title"
const AXValueNativeSourceTypeOther AXValueNativeSourceType = "other"

// A single source for a computed AX property.
type AXValueSource struct {
	Type              AXValueSourceType       `json:"type"`                        // What type of source this is.
	Value             *AXValue                `json:"value,omitempty"`             // The value of this property source.
	Attribute         string                  `json:"attribute,omitempty"`         // The name of the relevant attribute, if any.
	AttributeValue    *AXValue                `json:"attributeValue,omitempty"`    // The value of the relevant attribute, if any.
	Superseded        bool                    `json:"superseded,omitempty"`        // Whether this source is superseded by a higher priority source.
	NativeSource      AXValueNativeSourceType `json:"nativeSource,omitempty"`      // The native markup source for this value, e.g. a <label> element.
	NativeSourceValue *AXValue                `json:"nativeSourceValue,omitempty"` // The value, such as a node or node list, of the native source.
	Invalid           bool                    `json:"invalid,omitempty"`           // Whether the value for this property is invalid.
	InvalidReason     string                  `json:"invalidReason,omitempty"`     // Reason for the value being invalid, if it is.
}

type AXRelatedNode struct {
	BackendDOMNodeId *BackendNodeId `json:"backendDOMNodeId"` // The BackendNodeId of the related DOM node.
	Idref            string         `json:"idref,omitempty"`  // The IDRef value provided, if any.
	Text             string         `json:"text,omitempty"`   // The text alternative of this node in the current context.
}

type AXProperty struct {
	Name  string   `json:"name"`  // The name of this property.
	Value *AXValue `json:"value"` // The value of this property.
}

// A single computed AX property.
type AXValue struct {
	Type         AXValueType      `json:"type"`                   // The type of this value.
	Value        json.RawMessage  `json:"value,omitempty"`        // The computed value of this property.
	RelatedNodes []*AXRelatedNode `json:"relatedNodes,omitempty"` // One or more related nodes, if applicable.
	Sources      []*AXValueSource `json:"sources,omitempty"`      // The sources which contributed to the computation of this property.
}

// States which apply to every AX node.
type AXGlobalStates string

const AXGlobalStatesDisabled AXGlobalStates = "disabled"
const AXGlobalStatesHidden AXGlobalStates = "hidden"
const AXGlobalStatesHiddenRoot AXGlobalStates = "hiddenRoot"
const AXGlobalStatesInvalid AXGlobalStates = "invalid"

// Attributes which apply to nodes in live regions.
type AXLiveRegionAttributes string

const AXLiveRegionAttributesLive AXLiveRegionAttributes = "live"
const AXLiveRegionAttributesAtomic AXLiveRegionAttributes = "atomic"
const AXLiveRegionAttributesRelevant AXLiveRegionAttributes = "relevant"
const AXLiveRegionAttributesBusy AXLiveRegionAttributes = "busy"
const AXLiveRegionAttributesRoot AXLiveRegionAttributes = "root"

// Attributes which apply to widgets.
type AXWidgetAttributes string

const AXWidgetAttributesAutocomplete AXWidgetAttributes = "autocomplete"
const AXWidgetAttributesHaspopup AXWidgetAttributes = "haspopup"
const AXWidgetAttributesLevel AXWidgetAttributes = "level"
const AXWidgetAttributesMultiselectable AXWidgetAttributes = "multiselectable"
const AXWidgetAttributesOrientation AXWidgetAttributes = "orientation"
const AXWidgetAttributesMultiline AXWidgetAttributes = "multiline"
const AXWidgetAttributesReadonly AXWidgetAttributes = "readonly"
const AXWidgetAttributesRequired AXWidgetAttributes = "required"
const AXWidgetAttributesValuemin AXWidgetAttributes = "valuemin"
const AXWidgetAttributesValuemax AXWidgetAttributes = "valuemax"
const AXWidgetAttributesValuetext AXWidgetAttributes = "valuetext"

// States which apply to widgets.
type AXWidgetStates string

const AXWidgetStatesChecked AXWidgetStates = "checked"
const AXWidgetStatesExpanded AXWidgetStates = "expanded"
const AXWidgetStatesPressed AXWidgetStates = "pressed"
const AXWidgetStatesSelected AXWidgetStates = "selected"

// Relationships between elements other than parent/child/sibling.
type AXRelationshipAttributes string

const AXRelationshipAttributesActivedescendant AXRelationshipAttributes = "activedescendant"
const AXRelationshipAttributesFlowto AXRelationshipAttributes = "flowto"
const AXRelationshipAttributesControls AXRelationshipAttributes = "controls"
const AXRelationshipAttributesDescribedby AXRelationshipAttributes = "describedby"
const AXRelationshipAttributesLabelledby AXRelationshipAttributes = "labelledby"
const AXRelationshipAttributesOwns AXRelationshipAttributes = "owns"

// A node in the accessibility tree.
type AXNode struct {
	NodeId           AXNodeId       `json:"nodeId"`                     // Unique identifier for this node.
	Ignored          bool           `json:"ignored"`                    // Whether this node is ignored for accessibility
	IgnoredReasons   []*AXProperty  `json:"ignoredReasons,omitempty"`   // Collection of reasons why this node is hidden.
	Role             *AXValue       `json:"role,omitempty"`             // This Node's role, whether explicit or implicit.
	Name             *AXValue       `json:"name,omitempty"`             // The accessible name for this Node.
	Description      *AXValue       `json:"description,omitempty"`      // The accessible description for this Node.
	Value            *AXValue       `json:"value,omitempty"`            // The value for this Node.
	Properties       []*AXProperty  `json:"properties,omitempty"`       // All other properties
	ChildIds         []AXNodeId     `json:"childIds,omitempty"`         // IDs for each of this node's child nodes.
	BackendDOMNodeId *BackendNodeId `json:"backendDOMNodeId,omitempty"` // The backend ID for the associated DOM node, if any.
}

type GetPartialAXTreeParams struct {
	NodeId         *NodeId `json:"nodeId"`                   // ID of node to get the partial accessibility tree for.
	FetchRelatives bool    `json:"fetchRelatives,omitempty"` // Whether to fetch this nodes ancestors, siblings and children. Defaults to true.
}

type GetPartialAXTreeResult struct {
	Nodes []*AXNode `json:"nodes"` // The Accessibility.AXNode for this DOM node, if it exists, plus its ancestors, siblings and children, if requested.
}

// Fetches the accessibility node and partial accessibility tree for this DOM node, if it exists.
// @experimental
type GetPartialAXTreeCommand struct {
	params *GetPartialAXTreeParams
	result GetPartialAXTreeResult
	wg     sync.WaitGroup
	err    error
}

func NewGetPartialAXTreeCommand(params *GetPartialAXTreeParams) *GetPartialAXTreeCommand {
	return &GetPartialAXTreeCommand{
		params: params,
	}
}

func (cmd *GetPartialAXTreeCommand) Name() string {
	return "Accessibility.getPartialAXTree"
}

func (cmd *GetPartialAXTreeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetPartialAXTreeCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetPartialAXTree(params *GetPartialAXTreeParams, conn *hc.Conn) (result *GetPartialAXTreeResult, err error) {
	cmd := NewGetPartialAXTreeCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetPartialAXTreeCB func(result *GetPartialAXTreeResult, err error)

// Fetches the accessibility node and partial accessibility tree for this DOM node, if it exists.
// @experimental
type AsyncGetPartialAXTreeCommand struct {
	params *GetPartialAXTreeParams
	cb     GetPartialAXTreeCB
}

func NewAsyncGetPartialAXTreeCommand(params *GetPartialAXTreeParams, cb GetPartialAXTreeCB) *AsyncGetPartialAXTreeCommand {
	return &AsyncGetPartialAXTreeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetPartialAXTreeCommand) Name() string {
	return "Accessibility.getPartialAXTree"
}

func (cmd *AsyncGetPartialAXTreeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetPartialAXTreeCommand) Result() *GetPartialAXTreeResult {
	return &cmd.result
}

func (cmd *GetPartialAXTreeCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetPartialAXTreeCommand) Done(data []byte, err error) {
	var result GetPartialAXTreeResult
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
