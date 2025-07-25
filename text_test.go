package sarvam

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranslateString(t *testing.T) {
	translation := Translation{
		RequestId:      "123",
		TranslatedText: "Hello, world!",
		SourceLanguage: LanguageEnglish,
	}

	assert.Equal(t, translation.String(), "Hello, world!")
}
