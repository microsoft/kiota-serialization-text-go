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

func TestGetFloat32Value(t *testing.T) {
	source := "3.14"
	parseNode, err := NewTextParseNode([]byte(source))
	assert.Nil(t, err)
	val, err := parseNode.GetFloat32Value()
	assert.Nil(t, err)
	assert.NotNil(t, val)
	assert.InDelta(t, float32(3.14), *val, 0.001)
}

func TestGetFloat64Value(t *testing.T) {
	source := "3.141592653589793"
	parseNode, err := NewTextParseNode([]byte(source))
	assert.Nil(t, err)
	val, err := parseNode.GetFloat64Value()
	assert.Nil(t, err)
	assert.NotNil(t, val)
	assert.InDelta(t, 3.141592653589793, *val, 0.000000001)
}

func TestGetInt32Value(t *testing.T) {
	source := "42"
	parseNode, err := NewTextParseNode([]byte(source))
	assert.Nil(t, err)
	val, err := parseNode.GetInt32Value()
	assert.Nil(t, err)
	assert.NotNil(t, val)
	assert.Equal(t, int32(42), *val)
}

func TestGetInt64Value(t *testing.T) {
	source := "9876543210"
	parseNode, err := NewTextParseNode([]byte(source))
	assert.Nil(t, err)
	val, err := parseNode.GetInt64Value()
	assert.Nil(t, err)
	assert.NotNil(t, val)
	assert.Equal(t, int64(9876543210), *val)
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
