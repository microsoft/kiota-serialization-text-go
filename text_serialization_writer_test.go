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
