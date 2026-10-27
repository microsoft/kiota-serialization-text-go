package textserialization

import (
	testing "testing"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
	assert "github.com/stretchr/testify/assert"
)

func TestSerializationWriterFactoryHonoursInterface(t *testing.T) {
	instance := NewTextSerializationWriterFactory()
	assert.Implements(t, (*absser.SerializationWriterFactory)(nil), instance)
}

func TestTextSerializationWriterFactoryGetValidContentType(t *testing.T) {
	factory := NewTextSerializationWriterFactory()
	contentType, err := factory.GetValidContentType()
	assert.NoError(t, err)
	assert.Equal(t, "text/plain", contentType)

	var nilFactory *TextSerializationWriterFactory
	_, err = nilFactory.GetValidContentType()
	assert.Error(t, err)
}

func TestTextSerializationWriterFactoryGetSerializationWriter(t *testing.T) {
	factory := NewTextSerializationWriterFactory()

	writer, err := factory.GetSerializationWriter("text/plain")
	assert.NoError(t, err)
	assert.NotNil(t, writer)

	writer, err = factory.GetSerializationWriter("text/plain; charset=utf-8")
	assert.NoError(t, err)
	assert.NotNil(t, writer)

	writer, err = factory.GetSerializationWriter("  TEXT/PLAIN  ")
	assert.NoError(t, err)
	assert.NotNil(t, writer)

	_, err = factory.GetSerializationWriter("")
	assert.Error(t, err)

	_, err = factory.GetSerializationWriter("   ")
	assert.Error(t, err)

	_, err = factory.GetSerializationWriter("application/json")
	assert.Error(t, err)

	var nilFactory *TextSerializationWriterFactory
	_, err = nilFactory.GetSerializationWriter("text/plain")
	assert.Error(t, err)
}
