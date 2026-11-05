package textserialization

import (
	"errors"
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
	nodeTime, err := NewTextParseNode([]byte(" 15:30:00\n "))
	assert.NoError(t, err)
	toVal, err := nodeTime.GetTimeOnlyValue()
	assert.NoError(t, err)
	assert.NotNil(t, toVal)
	assert.Equal(t, "15:30:00", toVal.String())

	emptyTime, err := NewTextParseNode([]byte("   \n "))
	assert.NoError(t, err)
	emptyTo, err := emptyTime.GetTimeOnlyValue()
	assert.NoError(t, err)
	assert.Nil(t, emptyTo)

	for _, invalid := range []string{
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

	nodeDur, err := NewTextParseNode([]byte(" PT1H\r\n "))
	assert.NoError(t, err)
	durVal, err := nodeDur.GetISODurationValue()
	assert.NoError(t, err)
	assert.NotNil(t, durVal)
	assert.Equal(t, "PT1H", durVal.String())

	emptyDur, err := NewTextParseNode([]byte("   "))
	assert.NoError(t, err)
	emptyDurVal, err := emptyDur.GetISODurationValue()
	assert.NoError(t, err)
	assert.Nil(t, emptyDurVal)

	invalidDur, err := NewTextParseNode([]byte("invalid"))
	assert.NoError(t, err)
	_, err = invalidDur.GetISODurationValue()
	assert.Error(t, err)
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
	for _, invalid := range []string{"\"15:30:00", "15:30:00\"", " \"15:30:00\" "} {
		node, err := NewTextParseNode([]byte(invalid))
		assert.NoError(t, err)
		_, err = node.GetTimeOnlyValue()
		assert.Error(t, err)
	}

	for _, invalid := range []string{"\"2026-10-21", "2026-10-21\"", " \"2026-10-21\" "} {
		node, err := NewTextParseNode([]byte(invalid))
		assert.NoError(t, err)
		_, err = node.GetDateOnlyValue()
		assert.Error(t, err)
	}

	for _, invalid := range []string{"\"active", "active\"", " \"active\" "} {
		node, err := NewTextParseNode([]byte(invalid))
		assert.NoError(t, err)
		_, err = node.GetEnumValue(testStatusFactory)
		assert.Error(t, err)
	}

	for _, invalid := range []string{"\"true", "true\"", " \"true\" "} {
		node, err := NewTextParseNode([]byte(invalid))
		assert.NoError(t, err)
		_, err = node.GetBoolValue()
		assert.Error(t, err)
	}
}
