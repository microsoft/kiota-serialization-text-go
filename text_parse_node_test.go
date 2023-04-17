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

	err = instance.SetOnBeforeAssignFieldValues(nil)
	assert.NoError(t, err)

	err = instance.SetOnAfterAssignFieldValues(nil)
	assert.NoError(t, err)
}
