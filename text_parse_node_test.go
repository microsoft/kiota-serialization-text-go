package textserialization

import (
	testing "testing"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
	assert "github.com/stretchr/testify/assert"
)

func TestTree(t *testing.T) {
	source := "\"stringValue\""
	sourceArray := []byte(source)
	parseNode, err := NewTextParseNode(sourceArray)
	if err != nil {
		t.Errorf("Error creating parse node: %s", err.Error())
	}
	someProp, err := parseNode.GetChildNode("someProp")
	assert.NotNil(t, err)
	assert.Nil(t, someProp)

	stringValue, err := parseNode.GetStringValue()
	assert.Nil(t, err)
	assert.NotNil(t, stringValue)
	assert.Equal(t, "stringValue", *stringValue)
}

func TestTextParseNodeHonoursInterface(t *testing.T) {
	source := "\"stringValue\""
	sourceArray := []byte(source)
	instance, err := NewTextParseNode(sourceArray)
	assert.Nil(t, err)
	assert.Implements(t, (*absser.ParseNode)(nil), instance)

	action := func(parsable absser.Parsable) error {
		return nil
	}

	err = instance.SetOnBeforeAssignFieldValues(action)
	assert.NoError(t, err)
	assert.NotNil(t, instance.GetOnBeforeAssignFieldValues())

	err = instance.SetOnAfterAssignFieldValues(action)
	assert.NoError(t, err)
	assert.NotNil(t, instance.GetOnAfterAssignFieldValues())
}

func TestQuotedScalarsAndSingleQuote(t *testing.T) {
	// Quoted boolean
	boolNode, err := NewTextParseNode([]byte("\"true\""))
	assert.NoError(t, err)
	boolVal, err := boolNode.GetBoolValue()
	assert.NoError(t, err)
	assert.NotNil(t, boolVal)
	assert.True(t, *boolVal)

	// Quoted integer
	intNode, err := NewTextParseNode([]byte("\"42\""))
	assert.NoError(t, err)
	intVal, err := intNode.GetInt32Value()
	assert.NoError(t, err)
	assert.NotNil(t, intVal)
	assert.Equal(t, int32(42), *intVal)

	// Single quote
	singleQuoteNode, err := NewTextParseNode([]byte("\""))
	assert.NoError(t, err)
	singleQuoteVal, err := singleQuoteNode.GetStringValue()
	assert.NoError(t, err)
	assert.NotNil(t, singleQuoteVal)
	assert.Equal(t, "\"", *singleQuoteVal)

	// Empty quotes
	emptyQuoteNode, err := NewTextParseNode([]byte("\"\""))
	assert.NoError(t, err)
	emptyQuoteVal, err := emptyQuoteNode.GetStringValue()
	assert.NoError(t, err)
	assert.NotNil(t, emptyQuoteVal)
	assert.Equal(t, "", *emptyQuoteVal)
}

func TestUnbalancedQuotesAndNumbers(t *testing.T) {
	// Unbalanced leading quote should not be stripped
	leadingNode, err := NewTextParseNode([]byte("\"hello"))
	assert.NoError(t, err)
	leadingVal, err := leadingNode.GetStringValue()
	assert.NoError(t, err)
	assert.NotNil(t, leadingVal)
	assert.Equal(t, "\"hello", *leadingVal)

	// Unbalanced trailing quote should not be stripped
	trailingNode, err := NewTextParseNode([]byte("hello\""))
	assert.NoError(t, err)
	trailingVal, err := trailingNode.GetStringValue()
	assert.NoError(t, err)
	assert.NotNil(t, trailingVal)
	assert.Equal(t, "hello\"", *trailingVal)

	// Quoted float64
	floatNode, err := NewTextParseNode([]byte("\"3.14159\""))
	assert.NoError(t, err)
	floatVal, err := floatNode.GetFloat64Value()
	assert.NoError(t, err)
	assert.NotNil(t, floatVal)
	assert.Equal(t, 3.14159, *floatVal)
}


