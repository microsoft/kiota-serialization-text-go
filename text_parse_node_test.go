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

func TestGetByteValue(t *testing.T) {
	t.Run("parses valid unsigned byte values", func(t *testing.T) {
		for _, tc := range []struct {
			input    string
			expected byte
		}{
			{"0", 0},
			{"127", 127},
			{"128", 128},
			{"200", 200},
			{"255", 255},
		} {
			node, err := NewTextParseNode([]byte(tc.input))
			assert.NoError(t, err)
			val, err := node.GetByteValue()
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, *val)
		}
	})

	t.Run("rejects out of range byte values", func(t *testing.T) {
		for _, input := range []string{"256", "-1", "invalid"} {
			node, err := NewTextParseNode([]byte(input))
			assert.NoError(t, err)
			val, err := node.GetByteValue()
			assert.Error(t, err)
			assert.Nil(t, val)
		}
	})
}

func TestGetInt32Value(t *testing.T) {
	t.Run("parses valid int32", func(t *testing.T) {
		node, err := NewTextParseNode([]byte("2147483647"))
		assert.NoError(t, err)
		val, err := node.GetInt32Value()
		assert.NoError(t, err)
		assert.Equal(t, int32(2147483647), *val)
	})

	t.Run("parses negative int32", func(t *testing.T) {
		node, err := NewTextParseNode([]byte("-2147483648"))
		assert.NoError(t, err)
		val, err := node.GetInt32Value()
		assert.NoError(t, err)
		assert.Equal(t, int32(-2147483648), *val)
	})

	t.Run("parses leading-zero decimal int32", func(t *testing.T) {
		node, err := NewTextParseNode([]byte("010"))
		assert.NoError(t, err)
		val, err := node.GetInt32Value()
		assert.NoError(t, err)
		assert.Equal(t, int32(10), *val)
	})

	t.Run("coerces exponent int32 after integer overflow", func(t *testing.T) {
		node, err := NewTextParseNode([]byte("100000000000000000000e-20"))
		assert.NoError(t, err)
		val, err := node.GetInt32Value()
		assert.NoError(t, err)
		assert.Equal(t, int32(1), *val)
	})

	t.Run("rejects overflow int32", func(t *testing.T) {
		node, err := NewTextParseNode([]byte("2147483648"))
		assert.NoError(t, err)
		val, err := node.GetInt32Value()
		assert.Error(t, err)
		assert.Nil(t, val)
	})

	t.Run("coerces floating point integer format", func(t *testing.T) {
		node, err := NewTextParseNode([]byte("42.0"))
		assert.NoError(t, err)
		val, err := node.GetInt32Value()
		assert.NoError(t, err)
		assert.Equal(t, int32(42), *val)
	})
}

func TestGetInt64Value(t *testing.T) {
	t.Run("parses MaxInt64 without precision loss", func(t *testing.T) {
		node, err := NewTextParseNode([]byte("9223372036854775807"))
		assert.NoError(t, err)
		val, err := node.GetInt64Value()
		assert.NoError(t, err)
		assert.Equal(t, int64(9223372036854775807), *val)
	})

	t.Run("parses MinInt64", func(t *testing.T) {
		node, err := NewTextParseNode([]byte("-9223372036854775808"))
		assert.NoError(t, err)
		val, err := node.GetInt64Value()
		assert.NoError(t, err)
		assert.Equal(t, int64(-9223372036854775808), *val)
	})

	t.Run("parses leading-zero decimal int64", func(t *testing.T) {
		node, err := NewTextParseNode([]byte("010"))
		assert.NoError(t, err)
		val, err := node.GetInt64Value()
		assert.NoError(t, err)
		assert.Equal(t, int64(10), *val)
	})

	t.Run("coerces exponent int64 after integer overflow", func(t *testing.T) {
		node, err := NewTextParseNode([]byte("100000000000000000000e-20"))
		assert.NoError(t, err)
		val, err := node.GetInt64Value()
		assert.NoError(t, err)
		assert.Equal(t, int64(1), *val)
	})

	t.Run("rejects overflow int64", func(t *testing.T) {
		node, err := NewTextParseNode([]byte("9223372036854775808"))
		assert.NoError(t, err)
		val, err := node.GetInt64Value()
		assert.Error(t, err)
		assert.Nil(t, val)
	})
}
