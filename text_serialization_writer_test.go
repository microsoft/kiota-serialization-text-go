package textserialization

import (
	"testing"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
	assert "github.com/stretchr/testify/assert"
)

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
