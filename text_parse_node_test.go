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
		} {
			node, err := NewTextParseNode([]byte(tc.input))
			assert.NoError(t, err)
			val, err := node.GetByteValue()
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, *val)
		}
	})

	t.Run("rejects out of range byte values", func(t *testing.T) {
		for _, input := range []string{"256", "+256", "-1", "invalid"} {
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

func TestGetStringValue_WithUTF8BOM(t *testing.T) {
	node, err := NewTextParseNode([]byte("\xef\xbb\xbfHello World"))
	assert.NoError(t, err)
	assert.NotNil(t, node)

	val, err := node.GetStringValue()
	assert.NoError(t, err)
	assert.NotNil(t, val)
	assert.Equal(t, "Hello World", *val)
}

func TestQuotedString_WithUTF8BOM(t *testing.T) {
	node, err := NewTextParseNode([]byte("\xef\xbb\xbf\"Hello World\""))
	assert.NoError(t, err)
	assert.NotNil(t, node)

	val, err := node.GetStringValue()
	assert.NoError(t, err)
	assert.NotNil(t, val)
	assert.Equal(t, "Hello World", *val)

	emptyQuotedNode, err := NewTextParseNode([]byte("\xef\xbb\xbf\"\""))
	assert.NoError(t, err)
	assert.NotNil(t, emptyQuotedNode)

	emptyVal, err := emptyQuotedNode.GetStringValue()
	assert.NoError(t, err)
	assert.NotNil(t, emptyVal)
	assert.Equal(t, "", *emptyVal)
}

func TestEmptyContent_OnlyUTF8BOM(t *testing.T) {
	node, err := NewTextParseNode([]byte{0xef, 0xbb, 0xbf})
	assert.Error(t, err)
	assert.Nil(t, node)
	assert.EqualError(t, err, "content is empty")
}

func TestGetBoolValue_WithUTF8BOM(t *testing.T) {
	for _, tc := range []struct {
		input    string
		expected bool
	}{
		{"\xef\xbb\xbftrue", true},
		{"\xef\xbb\xbffalse", false},
		{"\xef\xbb\xbf1", true},
		{"\xef\xbb\xbf0", false},
	} {
		node, err := NewTextParseNode([]byte(tc.input))
		assert.NoError(t, err)
		val, err := node.GetBoolValue()
		assert.NoError(t, err)
		assert.NotNil(t, val)
		assert.Equal(t, tc.expected, *val)
	}
}

func TestGetIntegerValues_WithUTF8BOM(t *testing.T) {
	node32, err := NewTextParseNode([]byte("\xef\xbb\xbf42"))
	assert.NoError(t, err)
	val32, err := node32.GetInt32Value()
	assert.NoError(t, err)
	assert.NotNil(t, val32)
	assert.Equal(t, int32(42), *val32)

	node64, err := NewTextParseNode([]byte("\xef\xbb\xbf9223372036854775807"))
	assert.NoError(t, err)
	val64, err := node64.GetInt64Value()
	assert.NoError(t, err)
	assert.NotNil(t, val64)
	assert.Equal(t, int64(9223372036854775807), *val64)

	node8, err := NewTextParseNode([]byte("\xef\xbb\xbf-128"))
	assert.NoError(t, err)
	val8, err := node8.GetInt8Value()
	assert.NoError(t, err)
	assert.NotNil(t, val8)
	assert.Equal(t, int8(-128), *val8)

	nodeByte, err := NewTextParseNode([]byte("\xef\xbb\xbf255"))
	assert.NoError(t, err)
	valByte, err := nodeByte.GetByteValue()
	assert.NoError(t, err)
	assert.NotNil(t, valByte)
	assert.Equal(t, byte(255), *valByte)
}

func TestGetFloatValues_WithUTF8BOM(t *testing.T) {
	node32, err := NewTextParseNode([]byte("\xef\xbb\xbf3.14"))
	assert.NoError(t, err)
	val32, err := node32.GetFloat32Value()
	assert.NoError(t, err)
	assert.NotNil(t, val32)
	assert.Equal(t, float32(3.14), *val32)

	node64, err := NewTextParseNode([]byte("\xef\xbb\xbf2.718281828459"))
	assert.NoError(t, err)
	val64, err := node64.GetFloat64Value()
	assert.NoError(t, err)
	assert.NotNil(t, val64)
	assert.Equal(t, 2.718281828459, *val64)
}






