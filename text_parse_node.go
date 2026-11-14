// Package textserialization is the default Kiota serialization implementation for text.
package textserialization

import (
	"bytes"
	"encoding/base64"
	"errors"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

var (
	minInt64Rat = new(big.Rat).SetInt64(math.MinInt64)
	maxInt64Rat = new(big.Rat).SetInt64(math.MaxInt64)
)

// TextParseNode is a ParseNode implementation for text.
type TextParseNode struct {
	value                     string
	onBeforeAssignFieldValues absser.ParsableAction
	onAfterAssignFieldValues  absser.ParsableAction
}

// NewTextParseNode creates a new TextParseNode. If present, a leading UTF-8 byte order mark (BOM)
// is stripped prior to initializing the node.
func NewTextParseNode(content []byte) (*TextParseNode, error) {
	content = bytes.TrimPrefix(content, []byte("\xef\xbb\xbf"))
	if len(content) == 0 {
		return nil, errors.New("content is empty")
	}
	value, err := loadTextTree(content)
	return value, err
}

// loadTextTree initializes a TextParseNode, stripping enclosing double quotes when the content
// is wrapped with matching quotes and has length of at least 2.
func loadTextTree(content []byte) (*TextParseNode, error) {
	val := string(content)
	if len(val) >= 2 && strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"") {
		val = val[1 : len(val)-1]
	}
	return &TextParseNode{
		value: val,
	}, nil
}

// GetChildNode returns a new parse node for the given identifier.
func (n *TextParseNode) GetChildNode(index string) (absser.ParseNode, error) {
	if n == nil {
		return nil, nil
	}
	return nil, NoStructuredDataError
}

// GetObjectValue returns the Parsable value from the node.
func (n *TextParseNode) GetObjectValue(ctor absser.ParsableFactory) (absser.Parsable, error) {
	if n == nil {
		return nil, nil
	}
	return nil, NoStructuredDataError
}

// GetCollectionOfObjectValues returns the collection of Parsable values from the node.
func (n *TextParseNode) GetCollectionOfObjectValues(ctor absser.ParsableFactory) ([]absser.Parsable, error) {
	if n == nil {
		return nil, nil
	}
	return nil, NoStructuredDataError
}

// GetCollectionOfPrimitiveValues returns the collection of primitive values from the node.
func (n *TextParseNode) GetCollectionOfPrimitiveValues(targetType string) ([]interface{}, error) {
	if n == nil {
		return nil, nil
	}
	return nil, NoStructuredDataError
}

// GetCollectionOfEnumValues returns the collection of Enum values from the node.
func (n *TextParseNode) GetCollectionOfEnumValues(parser absser.EnumFactory) ([]interface{}, error) {
	if n == nil {
		return nil, nil
	}
	return nil, NoStructuredDataError
}

// GetStringValue returns a String value from the nodes.
func (n *TextParseNode) GetStringValue() (*string, error) {
	if n == nil {
		return nil, nil
	}
	val := n.value
	return &val, nil
}

// GetBoolValue returns a Bool value from the nodes. Leading and trailing whitespace
// and newlines are trimmed. Quoted values are preserved as-is and rejected as invalid syntax.
func (n *TextParseNode) GetBoolValue() (*bool, error) {
	if n == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(n.value)
	if trimmed == "" {
		return nil, nil
	}
	val, err := strconv.ParseBool(trimmed)
	if err != nil {
		return nil, err
	}
	return &val, nil
}

// GetInt8Value returns a int8 value from the nodes. Leading and trailing whitespace
// and newlines are trimmed.
func (n *TextParseNode) GetInt8Value() (*int8, error) {
	if n == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(n.value)
	if trimmed == "" {
		return nil, nil
	}
	if parsed, err := strconv.ParseInt(trimmed, 10, 8); err == nil {
		cast := int8(parsed)
		return &cast, nil
	} else if !errors.Is(err, strconv.ErrSyntax) && !strings.ContainsAny(trimmed, ".eE") {
		return nil, err
	}
	v, err := n.GetFloat64Value()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	if !(*v >= -128 && *v <= 127) {
		return nil, strconv.ErrRange
	}
	cast := int8(*v)
	return &cast, nil
}

// GetByteValue returns a Byte value from the nodes. Leading and trailing whitespace
// and newlines are trimmed.
func (n *TextParseNode) GetByteValue() (*byte, error) {
	if n == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(n.value)
	if trimmed == "" {
		return nil, nil
	}
	cleaned := strings.TrimPrefix(trimmed, "+")
	if parsed, err := strconv.ParseUint(cleaned, 10, 8); err == nil {
		cast := byte(parsed)
		return &cast, nil
	} else if !errors.Is(err, strconv.ErrSyntax) && !strings.ContainsAny(cleaned, ".eE") {
		return nil, err
	}
	v, err := n.GetFloat64Value()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	if !(*v >= 0 && *v <= 255) {
		return nil, strconv.ErrRange
	}
	cast := byte(*v)
	return &cast, nil
}

// GetFloat32Value returns a Float32 value from the nodes. Leading and trailing whitespace
// and newlines are trimmed.
func (n *TextParseNode) GetFloat32Value() (*float32, error) {
	if n == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(n.value)
	if trimmed == "" {
		return nil, nil
	}
	val, err := strconv.ParseFloat(trimmed, 32)
	if err != nil {
		return nil, err
	}
	cast := float32(val)
	return &cast, nil
}

// GetFloat64Value returns a Float64 value from the nodes. Leading and trailing whitespace
// and newlines are trimmed.
func (n *TextParseNode) GetFloat64Value() (*float64, error) {
	if n == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(n.value)
	if trimmed == "" {
		return nil, nil
	}
	val, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return nil, err
	}
	return &val, nil
}

// GetInt32Value returns a Int32 value from the nodes. Leading and trailing whitespace
// and newlines are trimmed.
func (n *TextParseNode) GetInt32Value() (*int32, error) {
	if n == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(n.value)
	if trimmed == "" {
		return nil, nil
	}
	if parsed, err := strconv.ParseInt(trimmed, 10, 32); err == nil {
		cast := int32(parsed)
		return &cast, nil
	} else if !errors.Is(err, strconv.ErrSyntax) && !strings.ContainsAny(trimmed, ".eE") {
		return nil, err
	}
	v, err := n.GetFloat64Value()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	if !(*v >= -2147483648 && *v <= 2147483647) {
		return nil, strconv.ErrRange
	}
	cast := int32(*v)
	return &cast, nil
}

// GetInt64Value returns a Int64 value from the nodes. Leading and trailing whitespace
// and newlines are trimmed.
func (n *TextParseNode) GetInt64Value() (*int64, error) {
	if n == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(n.value)
	if trimmed == "" {
		return nil, nil
	}
	if parsed, err := strconv.ParseInt(trimmed, 10, 64); err == nil {
		return &parsed, nil
	} else if !errors.Is(err, strconv.ErrSyntax) && !strings.ContainsAny(trimmed, ".eE") {
		return nil, err
	}
	if _, err := strconv.ParseFloat(trimmed, 64); err != nil && errors.Is(err, strconv.ErrSyntax) {
		return nil, err
	}
	rat, ok := new(big.Rat).SetString(trimmed)
	if !ok {
		return nil, strconv.ErrSyntax
	}
	if rat.Cmp(minInt64Rat) < 0 || rat.Cmp(maxInt64Rat) > 0 {
		return nil, strconv.ErrRange
	}
	var intVal big.Int
	intVal.Quo(rat.Num(), rat.Denom())
	cast := intVal.Int64()
	return &cast, nil
}

// GetTimeValue returns a Time value from the nodes. Leading and trailing whitespace
// and newlines are trimmed. Quoted values are preserved as-is and rejected as invalid syntax.
func (n *TextParseNode) GetTimeValue() (*time.Time, error) {
	if n == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(n.value)
	if trimmed == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, trimmed)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

// GetISODurationValue returns a ISODuration value from the nodes. Leading and trailing whitespace
// and newlines are trimmed. Quoted values are preserved as-is and rejected as invalid syntax.
func (n *TextParseNode) GetISODurationValue() (*absser.ISODuration, error) {
	if n == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(n.value)
	if trimmed == "" {
		return nil, nil
	}
	return absser.ParseISODuration(trimmed)
}

// GetTimeOnlyValue returns a TimeOnly value from the nodes. Leading and trailing whitespace
// and newlines are trimmed. Quoted values are preserved as-is and rejected as invalid syntax.
func (n *TextParseNode) GetTimeOnlyValue() (*absser.TimeOnly, error) {
	if n == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(n.value)
	if trimmed == "" {
		return nil, nil
	}
	return absser.ParseTimeOnly(trimmed)
}

// GetDateOnlyValue returns a DateOnly value from the nodes. Leading and trailing whitespace
// and newlines are trimmed. Quoted values are preserved as-is and rejected as invalid syntax.
func (n *TextParseNode) GetDateOnlyValue() (*absser.DateOnly, error) {
	if n == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(n.value)
	if trimmed == "" {
		return nil, nil
	}
	return absser.ParseDateOnly(trimmed)
}

// GetUUIDValue returns a UUID value from the nodes. Leading and trailing whitespace
// and newlines are trimmed. Quoted values are preserved as-is and rejected as invalid syntax.
func (n *TextParseNode) GetUUIDValue() (*uuid.UUID, error) {
	if n == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(n.value)
	if trimmed == "" {
		return nil, nil
	}
	parsed, err := uuid.Parse(trimmed)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

// GetEnumValue returns a Enum value from the nodes. Leading and trailing whitespace
// and newlines are trimmed. Quoted values are preserved as-is and rejected as invalid syntax.
func (n *TextParseNode) GetEnumValue(parser absser.EnumFactory) (interface{}, error) {
	if parser == nil {
		return nil, errors.New("parser is nil")
	}
	if n == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(n.value)
	if trimmed == "" {
		return nil, nil
	}
	return parser(trimmed)
}

// GetByteArrayValue returns a ByteArray value from the nodes. Leading and trailing whitespace
// and newlines are trimmed. Quoted values are preserved as-is and rejected as invalid syntax.
func (n *TextParseNode) GetByteArrayValue() ([]byte, error) {
	if n == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(n.value)
	if trimmed == "" {
		return nil, nil
	}
	if data, err := base64.StdEncoding.DecodeString(trimmed); err == nil {
		return data, nil
	}
	if data, err := base64.RawStdEncoding.DecodeString(trimmed); err == nil {
		return data, nil
	}
	if data, err := base64.URLEncoding.DecodeString(trimmed); err == nil {
		return data, nil
	}
	return base64.RawURLEncoding.DecodeString(trimmed)
}

// GetRawValue returns a ByteArray value from the nodes.
func (n *TextParseNode) GetRawValue() (interface{}, error) {
	if n == nil {
		return nil, nil
	}
	return n.value, nil
}

// GetOnBeforeAssignFieldValues returns a ByteArray value from the nodes.
func (n *TextParseNode) GetOnBeforeAssignFieldValues() absser.ParsableAction {
	if n == nil {
		return nil
	}
	return n.onBeforeAssignFieldValues
}

// SetOnBeforeAssignFieldValues returns a ByteArray value from the nodes.
func (n *TextParseNode) SetOnBeforeAssignFieldValues(action absser.ParsableAction) error {
	if n == nil {
		return errors.New("parse node cannot be nil")
	}
	n.onBeforeAssignFieldValues = action
	return nil
}

// GetOnAfterAssignFieldValues returns a ByteArray value from the nodes.
func (n *TextParseNode) GetOnAfterAssignFieldValues() absser.ParsableAction {
	if n == nil {
		return nil
	}
	return n.onAfterAssignFieldValues
}

// SetOnAfterAssignFieldValues returns a ByteArray value from the nodes.
func (n *TextParseNode) SetOnAfterAssignFieldValues(action absser.ParsableAction) error {
	if n == nil {
		return errors.New("parse node cannot be nil")
	}
	n.onAfterAssignFieldValues = action
	return nil
}
