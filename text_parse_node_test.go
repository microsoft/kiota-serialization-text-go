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

func TestGetInt8Value(t *testing.T) {
	t.Run("parses valid int8 values", func(t *testing.T) {
		for _, tc := range []struct {
			input    string
			expected int8
		}{
			{"0", 0},
			{"127", 127},
			{"-128", -128},
			{"+0", 0},
			{"+127", 127},
			{"010", 10},
			{"08", 8},
			{"09", 9},
			{"42.0", 42},
		} {
			node, err := NewTextParseNode([]byte(tc.input))
			assert.NoError(t, err)
			val, err := node.GetInt8Value()
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, *val)
		}
	})

	t.Run("rejects out of range int8 values", func(t *testing.T) {
		for _, input := range []string{"128", "-129", "256", "300.0", "-300.0", "invalid"} {
			node, err := NewTextParseNode([]byte(input))
			assert.NoError(t, err)
			val, err := node.GetInt8Value()
			assert.Error(t, err)
			assert.Nil(t, val)
		}
	})

	t.Run("nil receiver", func(t *testing.T) {
		var nilNode *TextParseNode
		val, err := nilNode.GetInt8Value()
		assert.NoError(t, err)
		assert.Nil(t, val)
	})
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
			{"+0", 0},
			{"+127", 127},
			{"+255", 255},
			{"010", 10},
			{"08", 8},
			{"09", 9},
			{"42.0", 42},
		} {
			node, err := NewTextParseNode([]byte(tc.input))
			assert.NoError(t, err)
			val, err := node.GetByteValue()
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, *val)
		}
	})

	t.Run("rejects out of range byte values", func(t *testing.T) {
		for _, input := range []string{"256", "+256", "-1", "300.0", "-10.0", "invalid"} {
			node, err := NewTextParseNode([]byte(input))
			assert.NoError(t, err)
			val, err := node.GetByteValue()
			assert.Error(t, err)
			assert.Nil(t, val)
		}
	})

	t.Run("nil receiver", func(t *testing.T) {
		var nilNode *TextParseNode
		val, err := nilNode.GetByteValue()
		assert.NoError(t, err)
		assert.Nil(t, val)
	})
}

func TestGetFloat32Value(t *testing.T) {
	t.Run("parses valid float32", func(t *testing.T) {
		node, err := NewTextParseNode([]byte("3.14159"))
		assert.NoError(t, err)
		val, err := node.GetFloat32Value()
		assert.NoError(t, err)
		assert.InDelta(t, float32(3.14159), *val, 0.0001)
	})

	t.Run("rejects overflow float32", func(t *testing.T) {
		node, err := NewTextParseNode([]byte("1e50"))
		assert.NoError(t, err)
		val, err := node.GetFloat32Value()
		assert.Error(t, err)
		assert.Nil(t, val)
	})

	t.Run("nil receiver", func(t *testing.T) {
		var nilNode *TextParseNode
		val, err := nilNode.GetFloat32Value()
		assert.NoError(t, err)
		assert.Nil(t, val)
	})
}

func TestGetFloat64Value(t *testing.T) {
	t.Run("parses valid float64", func(t *testing.T) {
		node, err := NewTextParseNode([]byte("1.7976931348623157e+308"))
		assert.NoError(t, err)
		val, err := node.GetFloat64Value()
		assert.NoError(t, err)
		assert.NotNil(t, val)
	})

	t.Run("nil receiver", func(t *testing.T) {
		var nilNode *TextParseNode
		val, err := nilNode.GetFloat64Value()
		assert.NoError(t, err)
		assert.Nil(t, val)
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

	t.Run("rejects coerced float overflow int32", func(t *testing.T) {
		node, err := NewTextParseNode([]byte("3000000000.0"))
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

	t.Run("rejects coerced float overflow int64", func(t *testing.T) {
		node, err := NewTextParseNode([]byte("1e25"))
		assert.NoError(t, err)
		val, err := node.GetInt64Value()
		assert.Error(t, err)
		assert.Nil(t, val)
	})
}

func TestNilParseNodeMethods(t *testing.T) {
	var nilNode *TextParseNode
	val, err := nilNode.GetRawValue()
	assert.NoError(t, err)
	assert.Nil(t, val)

	assert.Nil(t, nilNode.GetOnBeforeAssignFieldValues())
	assert.Error(t, nilNode.SetOnBeforeAssignFieldValues(func(p absser.Parsable) error { return nil }))
	assert.Nil(t, nilNode.GetOnAfterAssignFieldValues())
	assert.Error(t, nilNode.SetOnAfterAssignFieldValues(func(p absser.Parsable) error { return nil }))
}

func TestGetTimeValueWithSubseconds(t *testing.T) {
	inputs := []string{
		"2026-10-21T15:30:00.123456789Z",
		"2026-10-21T15:30:00.123Z",
		"2026-10-21T15:30:00.123456+05:30",
	}
	for _, input := range inputs {
		node, err := NewTextParseNode([]byte(input))
		assert.NoError(t, err)
		parsed, err := node.GetTimeValue()
		assert.NoError(t, err)
		assert.NotNil(t, parsed)
	}
}

func TestGetTimeValueInvalidReturnsNil(t *testing.T) {
	node, err := NewTextParseNode([]byte("invalid-time"))
	assert.NoError(t, err)
	parsed, err := node.GetTimeValue()
	assert.Error(t, err)
	assert.Nil(t, parsed)

	var nilNode *TextParseNode
	parsedNil, err := nilNode.GetTimeValue()
	assert.NoError(t, err)
	assert.Nil(t, parsedNil)
}

func TestGetUUIDValueInvalidReturnsNil(t *testing.T) {
	node, err := NewTextParseNode([]byte("invalid-uuid"))
	assert.NoError(t, err)
	parsed, err := node.GetUUIDValue()
	assert.Error(t, err)
	assert.Nil(t, parsed)

	var nilNode *TextParseNode
	parsedNil, err := nilNode.GetUUIDValue()
	assert.NoError(t, err)
	assert.Nil(t, parsedNil)
}
