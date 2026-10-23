package textserialization

import (
	"errors"
	"strings"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

// TextSerializationWriterFactory implements SerializationWriterFactory for text.
type TextSerializationWriterFactory struct {
}

// NewTextSerializationWriterFactory creates a new instance of the TextSerializationWriterFactory.
func NewTextSerializationWriterFactory() *TextSerializationWriterFactory {
	return &TextSerializationWriterFactory{}
}

// GetValidContentType returns the valid content type for the SerializationWriterFactoryRegistry
func (f *TextSerializationWriterFactory) GetValidContentType() (string, error) {
	if f == nil {
		return "", errors.New("factory cannot be nil")
	}
	return "text/plain", nil
}

// GetSerializationWriter returns the relevant SerializationWriter instance for the given content type
func (f *TextSerializationWriterFactory) GetSerializationWriter(contentType string) (absser.SerializationWriter, error) {
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

	return NewTextSerializationWriter(), nil
}
