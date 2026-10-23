package textserialization

import (
	testing "testing"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
	assert "github.com/stretchr/testify/assert"
)

func TestTextParseFactoryNodeHonoursInterface(t *testing.T) {
	instance := NewTextParseNodeFactory()
	assert.Implements(t, (*absser.ParseNodeFactory)(nil), instance)
}

func TestTextParseNodeFactoryGetValidContentType(t *testing.T) {
	factory := NewTextParseNodeFactory()
	contentType, err := factory.GetValidContentType()
	assert.NoError(t, err)
	assert.Equal(t, "text/plain", contentType)

	var nilFactory *TextParseNodeFactory
	_, err = nilFactory.GetValidContentType()
	assert.Error(t, err)
}

func TestTextParseNodeFactoryGetRootParseNode(t *testing.T) {
	factory := NewTextParseNodeFactory()

	node, err := factory.GetRootParseNode("text/plain", []byte("hello"))
	assert.NoError(t, err)
	assert.NotNil(t, node)

	node, err = factory.GetRootParseNode("text/plain; charset=utf-8", []byte("hello"))
	assert.NoError(t, err)
	assert.NotNil(t, node)

	node, err = factory.GetRootParseNode("  TEXT/PLAIN  ", []byte("hello"))
	assert.NoError(t, err)
	assert.NotNil(t, node)

	_, err = factory.GetRootParseNode("", []byte("hello"))
	assert.Error(t, err)

	_, err = factory.GetRootParseNode("   ", []byte("hello"))
	assert.Error(t, err)

	_, err = factory.GetRootParseNode("application/json", []byte("hello"))
	assert.Error(t, err)

	var nilFactory *TextParseNodeFactory
	_, err = nilFactory.GetRootParseNode("text/plain", []byte("hello"))
	assert.Error(t, err)
}
