package textserialization

import (
	"errors"
	"strings"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

// TextParseNodeFactory is a ParseNodeFactory implementation for text
type TextParseNodeFactory struct {
}

// Creates a new TextParseNodeFactory
func NewTextParseNodeFactory() *TextParseNodeFactory {
	return &TextParseNodeFactory{}
}

// GetValidContentType returns the content type this factory's parse nodes can deserialize.
func (f *TextParseNodeFactory) GetValidContentType() (string, error) {
	if f == nil {
		return "", errors.New("factory cannot be nil")
	}
	return "text/plain", nil
}

// GetRootParseNode return a new ParseNode instance that is the root of the content
func (f *TextParseNodeFactory) GetRootParseNode(contentType string, content []byte) (absser.ParseNode, error) {
	if f == nil {
		return nil, errors.New("factory cannot be nil")
	}
	validType, err := f.GetValidContentType()
	if err != nil {
		return nil, err
	}
	trimmed := strings.TrimSpace(contentType)
	if trimmed == "" {
		return nil, errors.New("contentType is empty")
	}

	mediaType, _, _ := strings.Cut(trimmed, ";")
	if strings.ToLower(strings.TrimSpace(mediaType)) != validType {
		return nil, errors.New("contentType is not valid")
	}

	return NewTextParseNode(content)
}
