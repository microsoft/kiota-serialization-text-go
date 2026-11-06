package textserialization

import (
	"testing"
	"time"

	"github.com/google/uuid"
	absser "github.com/microsoft/kiota-abstractions-go/serialization"
	assert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteTimeValueUsesRFC3339(t *testing.T) {
	tests := []struct {
		name  string
		value time.Time
		want  string
	}{
		{"UTC", time.Date(2026, 10, 16, 12, 34, 56, 0, time.UTC), "2026-10-16T12:34:56Z"},
		{"nanoseconds", time.Date(2026, 10, 16, 12, 34, 56, 123456789, time.UTC), "2026-10-16T12:34:56.123456789Z"},
		{"fractional seconds", time.Date(2026, 10, 16, 12, 34, 56, 120000000, time.UTC), "2026-10-16T12:34:56.12Z"},
		{"positive offset", time.Date(2026, 10, 16, 12, 34, 56, 0, time.FixedZone("test", 19800)), "2026-10-16T12:34:56+05:30"},
		{"negative offset", time.Date(2026, 10, 16, 12, 34, 56, 0, time.FixedZone("test", -14400)), "2026-10-16T12:34:56-04:00"},
		{"zero time", time.Time{}, "0001-01-01T00:00:00Z"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer, err := NewTextSerializationWriterFactory().GetSerializationWriter("text/plain")
			require.NoError(t, err)
			require.NoError(t, writer.WriteTimeValue("", &tt.value))
			content, err := writer.GetSerializedContent()
			require.NoError(t, err)
			assert.Equal(t, tt.want, string(content))

			node, err := NewTextParseNodeFactory().GetRootParseNode("text/plain", content)
			require.NoError(t, err)
			parsed, err := node.GetTimeValue()
			require.NoError(t, err)
			require.NotNil(t, parsed)
			assert.True(t, tt.value.Equal(*parsed))
			_, wantOffset := tt.value.Zone()
			_, gotOffset := parsed.Zone()
			assert.Equal(t, wantOffset, gotOffset)
		})
	}
}

func TestWriteTimeValueWithMonotonicClock(t *testing.T) {
	value := time.Now()
	writer := NewTextSerializationWriter()
	require.NoError(t, writer.WriteTimeValue("", &value))
	content, err := writer.GetSerializedContent()
	require.NoError(t, err)
	node, err := NewTextParseNode(content)
	require.NoError(t, err)
	parsed, err := node.GetTimeValue()
	require.NoError(t, err)
	require.NotNil(t, parsed)
	assert.True(t, value.Equal(*parsed))
}

func TestWriteTimeValuePreservesWriterGuards(t *testing.T) {
	value := time.Date(2026, 10, 16, 12, 34, 56, 0, time.UTC)
	writer := NewTextSerializationWriter()
	require.NoError(t, writer.WriteTimeValue("", nil))
	content, err := writer.GetSerializedContent()
	require.NoError(t, err)
	assert.Empty(t, content)
	assert.ErrorIs(t, writer.WriteTimeValue("nested", &value), NoStructuredDataError)
	require.NoError(t, writer.WriteTimeValue("", &value))
	assert.ErrorIs(t, writer.WriteTimeValue("", &value), OnlyOneValue)
}

func TestItDoesntWriteAnythingForAdditionalData(t *testing.T) {
	serializer := NewTextSerializationWriter()
	err := serializer.WriteAdditionalData(nil)
	assert.NotNil(t, err)
}

func TestSerializationWriterHonoursInterface(t *testing.T) {
	instance := NewTextSerializationWriter()
	assert.Implements(t, (*absser.SerializationWriter)(nil), instance)
}

func TestTextSerializationWriter(t *testing.T) {
	serializer := NewTextSerializationWriter()
	countBefore := 0
	onBefore := func(parsable absser.Parsable) error {
		countBefore++
		return nil
	}
	err := serializer.SetOnBeforeSerialization(onBefore)
	assert.NoError(t, err)
	assert.NotNil(t, serializer.GetOnBeforeSerialization())

	countAfter := 0
	onAfter := func(parsable absser.Parsable) error {
		countAfter++
		return nil
	}
	err = serializer.SetOnAfterObjectSerialization(onAfter)
	assert.NoError(t, err)
	assert.NotNil(t, serializer.GetOnAfterObjectSerialization())

	countStart := 0
	onStart := func(absser.Parsable, absser.SerializationWriter) error {
		countStart++
		return nil
	}

	err = serializer.SetOnStartObjectSerialization(onStart)
	assert.NoError(t, err)
	assert.NotNil(t, serializer.GetOnStartObjectSerialization())
}

func TestNilSerializationWriterMethods(t *testing.T) {
	var nilWriter *TextSerializationWriter
	content, err := nilWriter.GetSerializedContent()
	assert.Error(t, err)
	assert.Nil(t, content)

	testStr := "hello"
	assert.Error(t, nilWriter.WriteStringValue("", &testStr))

	assert.Nil(t, nilWriter.GetOnBeforeSerialization())
	assert.Error(t, nilWriter.SetOnBeforeSerialization(func(p absser.Parsable) error { return nil }))
	assert.Nil(t, nilWriter.GetOnAfterObjectSerialization())
	assert.Error(t, nilWriter.SetOnAfterObjectSerialization(func(p absser.Parsable) error { return nil }))
	assert.Nil(t, nilWriter.GetOnStartObjectSerialization())
	assert.Error(t, nilWriter.SetOnStartObjectSerialization(func(p absser.Parsable, w absser.SerializationWriter) error { return nil }))
	assert.NoError(t, nilWriter.Close())
}

func TestWriteEmptyStringEnforcesOnlyOneValue(t *testing.T) {
	writer := NewTextSerializationWriter()
	empty := ""
	require.NoError(t, writer.WriteStringValue("", &empty))
	second := "hello"
	assert.ErrorIs(t, writer.WriteStringValue("", &second), OnlyOneValue)
	content, err := writer.GetSerializedContent()
	require.NoError(t, err)
	assert.Equal(t, "", string(content))
}

func TestCloseClearsInternalBuffer(t *testing.T) {
	writer := NewTextSerializationWriter()
	val := "hello"
	require.NoError(t, writer.WriteStringValue("", &val))
	content, err := writer.GetSerializedContent()
	require.NoError(t, err)
	assert.Equal(t, "hello", string(content))

	require.NoError(t, writer.Close())
	contentAfter, err := writer.GetSerializedContent()
	require.NoError(t, err)
	assert.Empty(t, contentAfter)

	nextVal := "world"
	require.NoError(t, writer.WriteStringValue("", &nextVal))
	contentNext, err := writer.GetSerializedContent()
	require.NoError(t, err)
	assert.Equal(t, "world", string(contentNext))
}

func TestWriteScalars(t *testing.T) {
	t.Run("WriteFloat32Value preserves 32-bit shortest representation", func(t *testing.T) {
		writer := NewTextSerializationWriter()
		val := float32(0.1)
		require.NoError(t, writer.WriteFloat32Value("", &val))
		content, err := writer.GetSerializedContent()
		require.NoError(t, err)
		assert.Equal(t, "0.1", string(content))
	})

	t.Run("WriteFloat64Value", func(t *testing.T) {
		writer := NewTextSerializationWriter()
		val := 3.14159
		require.NoError(t, writer.WriteFloat64Value("", &val))
		content, err := writer.GetSerializedContent()
		require.NoError(t, err)
		assert.Equal(t, "3.14159", string(content))
	})

	t.Run("WriteBoolValue", func(t *testing.T) {
		writer := NewTextSerializationWriter()
		val := true
		require.NoError(t, writer.WriteBoolValue("", &val))
		content, err := writer.GetSerializedContent()
		require.NoError(t, err)
		assert.Equal(t, "true", string(content))
	})

	t.Run("WriteByteValue", func(t *testing.T) {
		writer := NewTextSerializationWriter()
		val := byte(255)
		require.NoError(t, writer.WriteByteValue("", &val))
		content, err := writer.GetSerializedContent()
		require.NoError(t, err)
		assert.Equal(t, "255", string(content))
	})

	t.Run("WriteInt8Value", func(t *testing.T) {
		writer := NewTextSerializationWriter()
		val := int8(-128)
		require.NoError(t, writer.WriteInt8Value("", &val))
		content, err := writer.GetSerializedContent()
		require.NoError(t, err)
		assert.Equal(t, "-128", string(content))
	})

	t.Run("WriteInt32Value", func(t *testing.T) {
		writer := NewTextSerializationWriter()
		val := int32(42)
		require.NoError(t, writer.WriteInt32Value("", &val))
		content, err := writer.GetSerializedContent()
		require.NoError(t, err)
		assert.Equal(t, "42", string(content))
	})

	t.Run("WriteInt64Value", func(t *testing.T) {
		writer := NewTextSerializationWriter()
		val := int64(9223372036854775807)
		require.NoError(t, writer.WriteInt64Value("", &val))
		content, err := writer.GetSerializedContent()
		require.NoError(t, err)
		assert.Equal(t, "9223372036854775807", string(content))
	})

	t.Run("WriteUUIDValue", func(t *testing.T) {
		writer := NewTextSerializationWriter()
		parsed := uuid.MustParse("6ec0bd7f-11c0-43da-975e-2a8ad9ebae0b")
		require.NoError(t, writer.WriteUUIDValue("", &parsed))
		content, err := writer.GetSerializedContent()
		require.NoError(t, err)
		assert.Equal(t, "6ec0bd7f-11c0-43da-975e-2a8ad9ebae0b", string(content))
	})

	t.Run("WriteByteArrayValue", func(t *testing.T) {
		writer := NewTextSerializationWriter()
		bytes := []byte("hello world")
		require.NoError(t, writer.WriteByteArrayValue("", bytes))
		content, err := writer.GetSerializedContent()
		require.NoError(t, err)
		assert.Equal(t, "aGVsbG8gd29ybGQ=", string(content))
	})

	t.Run("nil values write nothing and return no error", func(t *testing.T) {
		writer := NewTextSerializationWriter()
		require.NoError(t, writer.WriteBoolValue("", nil))
		require.NoError(t, writer.WriteByteValue("", nil))
		require.NoError(t, writer.WriteInt8Value("", nil))
		require.NoError(t, writer.WriteInt32Value("", nil))
		require.NoError(t, writer.WriteInt64Value("", nil))
		require.NoError(t, writer.WriteFloat32Value("", nil))
		require.NoError(t, writer.WriteFloat64Value("", nil))
		require.NoError(t, writer.WriteUUIDValue("", nil))
		require.NoError(t, writer.WriteByteArrayValue("", nil))
		require.NoError(t, writer.WriteISODurationValue("", nil))
		require.NoError(t, writer.WriteTimeOnlyValue("", nil))
		require.NoError(t, writer.WriteDateOnlyValue("", nil))
		content, err := writer.GetSerializedContent()
		require.NoError(t, err)
		assert.Empty(t, content)
	})
}
