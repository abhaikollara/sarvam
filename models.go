package sarvam

// ChatCompletionModel specifies the model to use for chat completions.
type ChatCompletionModel string

var (
	ChatCompletionModelBulbulV2 ChatCompletionModel = "bulbul:v2"
	ChatCompletionModelSarvamM  ChatCompletionModel = "sarvam-m"
)

// TextToSpeechModel specifies the model to use for text-to-speech conversion.
type TextToSpeechModel string

var (
	TextToSpeechModelBulbulV2 TextToSpeechModel = "bulbul:v2"
)

// SpeechToTextModel specifies the model to use for speech-to-text conversion.
type SpeechToTextModel string

var (
	SpeechToTextModelSaarikaV1     SpeechToTextModel = "saarika:v1"
	SpeechToTextModelSaarikaV2     SpeechToTextModel = "saarika:v2"
	SpeechToTextModelSaarikaV2dot5 SpeechToTextModel = "saarika:v2.5"
	SpeechToTextModelSaarikaFlash  SpeechToTextModel = "saarika:flash"
)

// SpeechToTextTranslateModel specifies the model to use for speech-to-text with translation.
type SpeechToTextTranslateModel string

var (
	SpeechToTextTranslateModelSaarasV1     SpeechToTextTranslateModel = "saaras:v1"
	SpeechToTextTranslateModelSaarasV2     SpeechToTextTranslateModel = "saaras:v2"
	SpeechToTextTranslateModelSaarasV2dot5 SpeechToTextTranslateModel = "saaras:v2.5"
	SpeechToTextTranslateModelSaarasFlash  SpeechToTextTranslateModel = "saaras:flash"
)

// TranslationModel specifies the translation model to use.
type TranslationModel string

var (
	TranslationModelMayuraV1        TranslationModel = "mayura:v1"
	TranslationModelSarvamTranslate TranslationModel = "sarvam-translate:v1"
)
