package textserialization

import (
	"errors"
	testing "testing"

	"github.com/google/uuid"
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
	// Quoted boolean is rejected as invalid syntax
	boolNode, err := NewTextParseNode([]byte("\"true\""))
	assert.NoError(t, err)
	boolVal, err := boolNode.GetBoolValue()
	assert.Error(t, err)
	assert.Nil(t, boolVal)

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

	boolVal, err := nilNode.GetBoolValue()
	assert.NoError(t, err)
	assert.Nil(t, boolVal)

	durVal, err := nilNode.GetISODurationValue()
	assert.NoError(t, err)
	assert.Nil(t, durVal)

	toVal, err := nilNode.GetTimeOnlyValue()
	assert.NoError(t, err)
	assert.Nil(t, toVal)

	doVal, err := nilNode.GetDateOnlyValue()
	assert.NoError(t, err)
	assert.Nil(t, doVal)

	enumVal, err := nilNode.GetEnumValue(testStatusFactory)
	assert.NoError(t, err)
	assert.Nil(t, enumVal)

	bytesVal, err := nilNode.GetByteArrayValue()
	assert.NoError(t, err)
	assert.Nil(t, bytesVal)

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

type testStatus int

const (
	testStatusActive testStatus = iota
	testStatusInactive
)

func testStatusFactory(value string) (interface{}, error) {
	switch value {
	case "active":
		return testStatusActive, nil
	case "inactive":
		return testStatusInactive, nil
	default:
		return nil, errors.New("unknown status: " + value)
	}
}

func TestGetBoolValue_WhitespaceAndNewlines(t *testing.T) {
	for _, tc := range []struct {
		input    string
		expected bool
	}{
		{"true\n", true},
		{"false\r\n", false},
		{" true ", true},
		{"\tfalse\t", false},
		{"1\n", true},
		{"0\r\n", false},
	} {
		node, err := NewTextParseNode([]byte(tc.input))
		assert.NoError(t, err)
		val, err := node.GetBoolValue()
		assert.NoError(t, err)
		assert.NotNil(t, val)
		assert.Equal(t, tc.expected, *val)
	}

	// Whitespace only returns (nil, nil)
	emptyNode, err := NewTextParseNode([]byte("   \n\t  "))
	assert.NoError(t, err)
	val, err := emptyNode.GetBoolValue()
	assert.NoError(t, err)
	assert.Nil(t, val)

	// Values with quotes are rejected
	for _, invalid := range []string{
		"\"true\"",
		"\"false\"",
		"\" true \"",
		"\" false \"",
		"\"   \"",
		"\"\"",
		"\"true\"\n",
		" \" true \"\n",
		" \" false \"\r\n",
		" \" 1 \"\n",
		" \" 0 \"\r\n",
		" \"   \n\t  \" \n",
		"\"true",
		"true\"",
		"\"false",
		"false\"",
		"\"",
	} {
		node, err := NewTextParseNode([]byte(invalid))
		assert.NoError(t, err)
		v, err := node.GetBoolValue()
		assert.Error(t, err)
		assert.Nil(t, v)
	}
}

func TestGetDateOnlyValue_WhitespaceAndNewlines(t *testing.T) {
	for _, input := range []string{
		"2026-10-21\n",
		" 2026-10-21\r\n",
		"\t2026-10-21\t",
	} {
		node, err := NewTextParseNode([]byte(input))
		assert.NoError(t, err)
		dateVal, err := node.GetDateOnlyValue()
		assert.NoError(t, err)
		assert.NotNil(t, dateVal)
		assert.Equal(t, "2026-10-21", dateVal.String())
	}

	emptyNode, err := NewTextParseNode([]byte("   "))
	assert.NoError(t, err)
	val, err := emptyNode.GetDateOnlyValue()
	assert.NoError(t, err)
	assert.Nil(t, val)

	// Quoted values are rejected
	for _, invalid := range []string{
		"\"2026-10-21\"",
		"\" 2026-10-21 \"",
		"\"   \"",
		"\"\"",
		" \" 2026-10-21 \"\n",
		" \"   \" ",
		"\"2026-10-21\n",
		"\"2026-10-21",
		"2026-10-21\"",
	} {
		node, err := NewTextParseNode([]byte(invalid))
		assert.NoError(t, err)
		v, err := node.GetDateOnlyValue()
		assert.Error(t, err)
		assert.Nil(t, v)
	}
}

func TestGetTimeOnlyAndISODuration_WhitespaceAndNewlines(t *testing.T) {
	for _, input := range []string{
		" 15:30:00\n ",
		"15:30:00\r\n",
		"\t15:30:00\t",
	} {
		nodeTime, err := NewTextParseNode([]byte(input))
		assert.NoError(t, err)
		toVal, err := nodeTime.GetTimeOnlyValue()
		assert.NoError(t, err)
		assert.NotNil(t, toVal)
		assert.Equal(t, "15:30:00", toVal.String())
	}

	emptyTime, err := NewTextParseNode([]byte("   \n "))
	assert.NoError(t, err)
	emptyTo, err := emptyTime.GetTimeOnlyValue()
	assert.NoError(t, err)
	assert.Nil(t, emptyTo)

	for _, invalid := range []string{
		"\"15:30:00\"",
		"\" 15:30:00 \"",
		"\"   \"",
		"\"\"",
		" \" 15:30:00 \"\n ",
		" \"   \" ",
		"\"15:30:00\n",
		"\"15:30:00",
		"15:30:00\"",
	} {
		node, err := NewTextParseNode([]byte(invalid))
		assert.NoError(t, err)
		v, err := node.GetTimeOnlyValue()
		assert.Error(t, err)
		assert.Nil(t, v)
	}

	for _, input := range []string{
		" PT1H\r\n ",
		"PT1H\n",
		"\tPT1H\t",
	} {
		nodeDur, err := NewTextParseNode([]byte(input))
		assert.NoError(t, err)
		durVal, err := nodeDur.GetISODurationValue()
		assert.NoError(t, err)
		assert.NotNil(t, durVal)
		assert.Equal(t, "PT1H", durVal.String())
	}

	emptyDur, err := NewTextParseNode([]byte("   "))
	assert.NoError(t, err)
	emptyDurVal, err := emptyDur.GetISODurationValue()
	assert.NoError(t, err)
	assert.Nil(t, emptyDurVal)

	for _, invalid := range []string{
		"invalid",
		"\"PT1H\"",
		"\" PT1H \"",
		"\"   \"",
		"\"\"",
		" \" PT1H \"\n ",
		" \"   \" ",
		"\"PT1H\n",
		"\"PT1H",
		"PT1H\"",
	} {
		invalidDur, err := NewTextParseNode([]byte(invalid))
		assert.NoError(t, err)
		v, err := invalidDur.GetISODurationValue()
		assert.Error(t, err)
		assert.Nil(t, v)
	}
}

func TestGetEnumValue_WhitespaceAndNewlines(t *testing.T) {
	nodeActive, err := NewTextParseNode([]byte(" active\n "))
	assert.NoError(t, err)
	enumVal, err := nodeActive.GetEnumValue(testStatusFactory)
	assert.NoError(t, err)
	assert.Equal(t, testStatusActive, enumVal)

	nodeInactive, err := NewTextParseNode([]byte("inactive\r\n"))
	assert.NoError(t, err)
	enumVal, err = nodeInactive.GetEnumValue(testStatusFactory)
	assert.NoError(t, err)
	assert.Equal(t, testStatusInactive, enumVal)

	emptyNode, err := NewTextParseNode([]byte("   "))
	assert.NoError(t, err)
	val, err := emptyNode.GetEnumValue(testStatusFactory)
	assert.NoError(t, err)
	assert.Nil(t, val)

	for _, invalid := range []string{
		"\"active\"",
		"\"inactive\"",
		"\" active \"",
		"\" inactive \"",
		"\"   \"",
		"\"\"",
		" \" active \"\n ",
		" \" inactive \"\r\n",
		" \"   \" ",
		"\"active\n",
		"\"active",
		"active\"",
	} {
		node, err := NewTextParseNode([]byte(invalid))
		assert.NoError(t, err)
		v, err := node.GetEnumValue(testStatusFactory)
		assert.Error(t, err)
		assert.Nil(t, v)
	}
}

func TestGetByteArrayValue_WhitespaceAndEncodings(t *testing.T) {
	expected := []byte("Hello world")

	for _, input := range []string{
		"SGVsbG8gd29ybGQ=",
		"  SGVsbG8gd29ybGQ= \n",
		"SGVsbG8gd29ybGQ",
	} {
		node, err := NewTextParseNode([]byte(input))
		assert.NoError(t, err)
		data, err := node.GetByteArrayValue()
		assert.NoError(t, err)
		assert.Equal(t, expected, data)
	}

	// URL-safe encoded binary data
	urlPaddedNode, err := NewTextParseNode([]byte("-_=="))
	assert.NoError(t, err)
	data, err := urlPaddedNode.GetByteArrayValue()
	assert.NoError(t, err)
	assert.Equal(t, []byte{251}, data)

	urlRawNode, err := NewTextParseNode([]byte("-_"))
	assert.NoError(t, err)
	data, err = urlRawNode.GetByteArrayValue()
	assert.NoError(t, err)
	assert.Equal(t, []byte{251}, data)

	// Whitespace only returns (nil, nil)
	emptyNode, err := NewTextParseNode([]byte("   \n "))
	assert.NoError(t, err)
	emptyData, err := emptyNode.GetByteArrayValue()
	assert.NoError(t, err)
	assert.Nil(t, emptyData)

	// Invalid base64 and quoted inputs return error
	for _, invalid := range []string{
		"!!invalid!!",
		"\"SGVsbG8gd29ybGQ=\"",
		"\" SGVsbG8gd29ybGQ= \"",
		"\"   \"",
		"\"\"",
		"\"SGVsbG8gd29ybGQ=\"\r\n",
		" \" SGVsbG8gd29ybGQ= \"\n",
		" \"   \" ",
		"\"SGVsbG8=",
		"SGVsbG8=\"",
	} {
		node, err := NewTextParseNode([]byte(invalid))
		assert.NoError(t, err)
		_, err = node.GetByteArrayValue()
		assert.Error(t, err)
	}
}

func TestScalarParsers_RejectQuotedValues(t *testing.T) {
	// Plain text serialization format does not strip enclosing quotes per maintainer guidance.
	// Quoted strings are treated as invalid syntax for scalar types (bool, temporal, enum, etc.).
	for _, invalid := range []string{"\"15:30:00\"", "\" 15:30:00 \"", "\"   \"", "\"\"", "\"15:30:00", "15:30:00\"", " \"15:30:00\" ", "\"15:30:00\"\n"} {
		node, err := NewTextParseNode([]byte(invalid))
		assert.NoError(t, err)
		_, err = node.GetTimeOnlyValue()
		assert.Error(t, err)
	}

	for _, invalid := range []string{"\"2026-10-21\"", "\" 2026-10-21 \"", "\"   \"", "\"\"", "\"2026-10-21", "2026-10-21\"", " \"2026-10-21\" ", "\"2026-10-21\"\r\n"} {
		node, err := NewTextParseNode([]byte(invalid))
		assert.NoError(t, err)
		_, err = node.GetDateOnlyValue()
		assert.Error(t, err)
	}

	for _, invalid := range []string{"\"PT1H\"", "\" PT1H \"", "\"   \"", "\"\"", "\"PT1H", "PT1H\"", " \"PT1H\" ", "\"PT1H\"\n"} {
		node, err := NewTextParseNode([]byte(invalid))
		assert.NoError(t, err)
		_, err = node.GetISODurationValue()
		assert.Error(t, err)
	}

	for _, invalid := range []string{"\"active\"", "\"inactive\"", "\" active \"", "\" inactive \"", "\"   \"", "\"\"", "\"active", "active\"", " \"active\" ", "\"active\"\n"} {
		node, err := NewTextParseNode([]byte(invalid))
		assert.NoError(t, err)
		_, err = node.GetEnumValue(testStatusFactory)
		assert.Error(t, err)
	}

	for _, invalid := range []string{"\"true\"", "\"false\"", "\" true \"", "\" false \"", "\"   \"", "\"\"", "\"true", "true\"", " \"true\" ", "\"true\"\n", "\"false\"\r\n", " \" true \"\n"} {
		node, err := NewTextParseNode([]byte(invalid))
		assert.NoError(t, err)
		_, err = node.GetBoolValue()
		assert.Error(t, err)
	}

	for _, invalid := range []string{"\"SGVsbG8gd29ybGQ=\"", "\" SGVsbG8gd29ybGQ= \"", "\"   \"", "\"\"", "\"SGVsbG8=", "SGVsbG8=\""} {
		node, err := NewTextParseNode([]byte(invalid))
		assert.NoError(t, err)
		_, err = node.GetByteArrayValue()
		assert.Error(t, err)
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

func TestGetUUIDValue_WithUTF8BOM(t *testing.T) {
	node, err := NewTextParseNode([]byte("\xef\xbb\xbf6ec0bd7f-11c0-43da-975e-2a8ad9ebae0b"))
	assert.NoError(t, err)
	val, err := node.GetUUIDValue()
	assert.NoError(t, err)
	assert.NotNil(t, val)
	expected := uuid.MustParse("6ec0bd7f-11c0-43da-975e-2a8ad9ebae0b")
	assert.Equal(t, expected, *val)
}

func TestGetTemporalValues_WithUTF8BOM(t *testing.T) {
	timeNode, err := NewTextParseNode([]byte("\xef\xbb\xbf2026-10-21T15:30:00Z"))
	assert.NoError(t, err)
	timeVal, err := timeNode.GetTimeValue()
	assert.NoError(t, err)
	assert.NotNil(t, timeVal)
	assert.Equal(t, 2026, timeVal.Year())

	dateNode, err := NewTextParseNode([]byte("\xef\xbb\xbf2026-10-21"))
	assert.NoError(t, err)
	dateVal, err := dateNode.GetDateOnlyValue()
	assert.NoError(t, err)
	assert.NotNil(t, dateVal)
	assert.Equal(t, "2026-10-21", dateVal.String())

	timeOnlyNode, err := NewTextParseNode([]byte("\xef\xbb\xbf15:30:00"))
	assert.NoError(t, err)
	timeOnlyVal, err := timeOnlyNode.GetTimeOnlyValue()
	assert.NoError(t, err)
	assert.NotNil(t, timeOnlyVal)
	assert.Equal(t, "15:30:00", timeOnlyVal.String())

	durationNode, err := NewTextParseNode([]byte("\xef\xbb\xbfP1D"))
	assert.NoError(t, err)
	durationVal, err := durationNode.GetISODurationValue()
	assert.NoError(t, err)
	assert.NotNil(t, durationVal)
	assert.Equal(t, "P1D", durationVal.String())
}
