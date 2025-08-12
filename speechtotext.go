package sarvam

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Timestamps represents word-level timing information for speech-to-text results.
type Timestamps struct {
	Words            []string  `json:"words"`
	StartTimeSeconds []float64 `json:"start_time_seconds"`
	EndTimeSeconds   []float64 `json:"end_time_seconds"`
}

// DiarizedEntry represents a single speaker's transcript segment.
type DiarizedEntry struct {
	Transcript       string  `json:"transcript"`
	StartTimeSeconds float64 `json:"start_time_seconds"`
	EndTimeSeconds   float64 `json:"end_time_seconds"`
	SpeakerID        string  `json:"speaker_id"`
}

// DiarizedTranscript represents the complete diarized transcript.
type DiarizedTranscript struct {
	Entries []DiarizedEntry `json:"entries"`
}

// SpeechToTextResponse represents the result of a speech-to-text operation.
type SpeechToTextResponse struct {
	RequestId          string              `json:"request_id"`
	Transcript         string              `json:"transcript"`
	Timestamps         *Timestamps         `json:"timestamps,omitempty"`
	DiarizedTranscript *DiarizedTranscript `json:"diarized_transcript,omitempty"`
	LanguageCode       Language            `json:"language_code"`
}

// String returns the transcribed text.
func (s *SpeechToTextResponse) String() string {
	return s.Transcript
}

// SpeechToTextParams contains parameters for speech-to-text conversion.
type SpeechToTextParams struct {
	Model          *SpeechToTextModel // Optional: Model to use (default: saarika:v2.5)
	LanguageCode   *Language          // Optional: Language code for the input audio
	WithTimestamps *bool              // Optional: Whether to include timestamps in response
}

// SpeechToText converts speech from an audio file to text.
func (c *Client) SpeechToText(speech io.Reader, params SpeechToTextParams) (*SpeechToTextResponse, error) {
	resp, err := c.buildSpeechToTextRequest("/speech-to-text", speech, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, parseAPIError(resp)
	}

	// Parse the response
	type speechToTextResponse struct {
		RequestId          string              `json:"request_id"`
		Transcript         string              `json:"transcript"`
		Timestamps         *Timestamps         `json:"timestamps,omitempty"`
		DiarizedTranscript *DiarizedTranscript `json:"diarized_transcript,omitempty"`
		LanguageCode       string              `json:"language_code"`
	}

	var response speechToTextResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &SpeechToTextResponse{
		RequestId:          response.RequestId,
		Transcript:         response.Transcript,
		Timestamps:         response.Timestamps,
		DiarizedTranscript: response.DiarizedTranscript,
		LanguageCode:       mapLanguageCodeToLanguage(response.LanguageCode),
	}, nil
}

// SpeechToTextTranslateResponse represents the result of a speech-to-text-translate operation.
type SpeechToTextTranslateResponse struct {
	RequestId          string              `json:"request_id"`
	Transcript         string              `json:"transcript"`
	LanguageCode       Language            `json:"language_code"`
	DiarizedTranscript *DiarizedTranscript `json:"diarized_transcript,omitempty"`
}

// String returns the transcribed and translated text.
func (s *SpeechToTextTranslateResponse) String() string {
	return s.Transcript
}

// SpeechToTextTranslateParams contains parameters for speech-to-text-translate conversion.
type SpeechToTextTranslateParams struct {
	Prompt     *string                     // Optional: Conversation context to boost model accuracy
	Model      *SpeechToTextTranslateModel // Optional: Model to use for speech-to-text conversion
	AudioCodec *AudioCodec                 // Optional: Audio codec to use for speech-to-text conversion
}

type AudioCodec string

var (
	AudioCodecWav      AudioCodec = "wav"
	AudioCodecXWav     AudioCodec = "x-wav"
	AudioCodecWave     AudioCodec = "wave"
	AudioCodecMp3      AudioCodec = "mp3"
	AudioCodecMpeg     AudioCodec = "mpeg"
	AudioCodecMpeg3    AudioCodec = "mpeg3"
	AudioCodecXMp3     AudioCodec = "x-mp3"
	AudioCodecXMpeg3   AudioCodec = "x-mpeg-3"
	AudioCodecAac      AudioCodec = "aac"
	AudioCodecXAac     AudioCodec = "x-aac"
	AudioCodecAiff     AudioCodec = "aiff"
	AudioCodecXAiff    AudioCodec = "x-aiff"
	AudioCodecOgg      AudioCodec = "ogg"
	AudioCodecOpus     AudioCodec = "opus"
	AudioCodecFlac     AudioCodec = "flac"
	AudioCodecXFlac    AudioCodec = "x-flac"
	AudioCodecMp4      AudioCodec = "mp4"
	AudioCodecXM4a     AudioCodec = "x-m4a"
	AudioCodecAmr      AudioCodec = "amr"
	AudioCodecXMsWma   AudioCodec = "x-ms-wma"
	AudioCodecWebm     AudioCodec = "webm"
	AudioCodecPcmS16le AudioCodec = "pcm_s16le"
	AudioCodecPcmL16   AudioCodec = "pcm_l16"
	AudioCodecPcmRaw   AudioCodec = "pcm_raw"
)

// SpeechToTextTranslate automatically detects the input language, transcribes the speech, and translates the text to English.
func (c *Client) SpeechToTextTranslate(speech io.Reader, params SpeechToTextTranslateParams) (*SpeechToTextTranslateResponse, error) {
	resp, err := c.buildSpeechToTextTranslateRequest("/speech-to-text-translate", speech, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, parseAPIError(resp)
	}

	// Parse the response
	type speechToTextTranslateResponse struct {
		RequestId          string              `json:"request_id"`
		Transcript         string              `json:"transcript"`
		LanguageCode       string              `json:"language_code"`
		DiarizedTranscript *DiarizedTranscript `json:"diarized_transcript,omitempty"`
	}

	var response speechToTextTranslateResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &SpeechToTextTranslateResponse{
		RequestId:          response.RequestId,
		Transcript:         response.Transcript,
		LanguageCode:       mapLanguageCodeToLanguage(response.LanguageCode),
		DiarizedTranscript: response.DiarizedTranscript,
	}, nil
}
