package protocol

import (
	"encoding/json"
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
	Type              AXValueSourceType       `json:"type"`              // What type of source this is.
	Value             *AXValue                `json:"value"`             // The value of this property source.
	Attribute         string                  `json:"attribute"`         // The name of the relevant attribute, if any.
	AttributeValue    *AXValue                `json:"attributeValue"`    // The value of the relevant attribute, if any.
	Superseded        bool                    `json:"superseded"`        // Whether this source is superseded by a higher priority source.
	NativeSource      AXValueNativeSourceType `json:"nativeSource"`      // The native markup source for this value, e.g. a <label> element.
	NativeSourceValue *AXValue                `json:"nativeSourceValue"` // The value, such as a node or node list, of the native source.
	Invalid           bool                    `json:"invalid"`           // Whether the value for this property is invalid.
	InvalidReason     string                  `json:"invalidReason"`     // Reason for the value being invalid, if it is.
}

type AXRelatedNode struct {
	BackendNodeId *BackendNodeId `json:"backendNodeId"` // The BackendNodeId of the related node.
	Idref         string         `json:"idref"`         // The IDRef value provided, if any.
	Text          string         `json:"text"`          // The text alternative of this node in the current context.
}

type AXProperty struct {
	Name  string   `json:"name"`  // The name of this property.
	Value *AXValue `json:"value"` // The value of this property.
}

// A single computed AX property.
type AXValue struct {
	Type         AXValueType      `json:"type"`         // The type of this value.
	Value        string           `json:"value"`        // The computed value of this property.
	RelatedNodes []*AXRelatedNode `json:"relatedNodes"` // One or more related nodes, if applicable.
	Sources      []*AXValueSource `json:"sources"`      // The sources which contributed to the computation of this property.
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
	NodeId         AXNodeId      `json:"nodeId"`         // Unique identifier for this node.
	Ignored        bool          `json:"ignored"`        // Whether this node is ignored for accessibility
	IgnoredReasons []*AXProperty `json:"ignoredReasons"` // Collection of reasons why this node is hidden.
	Role           *AXValue      `json:"role"`           // This Node's role, whether explicit or implicit.
	Name           *AXValue      `json:"name"`           // The accessible name for this Node.
	Description    *AXValue      `json:"description"`    // The accessible description for this Node.
	Value          *AXValue      `json:"value"`          // The value for this Node.
	Properties     []*AXProperty `json:"properties"`     // All other properties
}

type GetAXNodeChainParams struct {
	NodeId         *NodeId `json:"nodeId"`         // ID of node to get accessibility node for.
	FetchAncestors bool    `json:"fetchAncestors"` // Whether to also push down a partial tree (parent chain).
}

type GetAXNodeChainResult struct {
	Nodes []*AXNode `json:"nodes"` // The Accessibility.AXNode for this DOM node, if it exists, plus ancestors if requested.
}

type GetAXNodeChainCB func(result *GetAXNodeChainResult, err error)

// Fetches the accessibility node for this DOM node, if it exists.
type GetAXNodeChainCommand struct {
	params *GetAXNodeChainParams
	cb     GetAXNodeChainCB
}

func NewGetAXNodeChainCommand(params *GetAXNodeChainParams, cb GetAXNodeChainCB) *GetAXNodeChainCommand {
	return &GetAXNodeChainCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetAXNodeChainCommand) Name() string {
	return "Accessibility.getAXNodeChain"
}

func (cmd *GetAXNodeChainCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetAXNodeChainCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetAXNodeChainResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}
