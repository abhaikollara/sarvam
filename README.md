# SarvamAI Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/abhaikollara/sarvam.svg)](https://pkg.go.dev/github.com/abhaikollara/sarvam)
[![Go Report Card](https://goreportcard.com/badge/abhaikollara/sarvam)](https://goreportcard.com/report/abhaikollara/sarvam)
[![Go Version](https://img.shields.io/github/go-mod/go-version/abhaikollara/sarvam)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

An unofficial Go SDK for the [Sarvam AI](https://sarvam.ai) APIs, providing easy access to Indian language AI services including translation, text-to-speech, chat completions, and language identification.

> **⚠️ Breaking Changes Notice**: This SDK is currently in pre-v1 development. Breaking changes may occur until v1.0.0 is released. Please pin your dependency version if you need stability in production environments.

## 🌟 API Parity (wip)

- [x] Text Translation
- [x] Language Identification
- [x] Text-to-Speech
- [x] Chat Completions
- [x] Transliteration
- [x] Speech to text
- [x] Speech to text translation


## 🚀 Quick Start

### Installation

```bash
go get code.abhai.dev/sarvam
```

### Basic Usage

The SDK provides both instance-based and package-level APIs for convenience.

#### Using Package-Level Functions (Recommended for simple use cases)

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "code.abhai.dev/sarvam"
)

func main() {
    // Set API key (or set SARVAM_API_KEY environment variable)
    sarvam.SetAPIKey("your-api-key-here")
    
    // Use package-level functions directly
    result, err := sarvam.SpeechToText(sarvam.SpeechToTextParams{
        FilePath: "audio.wav",
        Model:    &sarvam.SpeechToTextModelSaarikaV2dot5,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Transcript:", result.Transcript)
}
```

#### Using Instance-Based Client (Recommended for advanced use cases)

```go
package main

import (
    "fmt"
    "log"
    
    "code.abhai.dev/sarvam"
)

func main() {
    // Create a client instance
    client := sarvam.NewClient("your-api-key-here")
    
    // Use the client instance
    result, err := client.SpeechToText(sarvam.SpeechToTextParams{
        FilePath: "audio.wav",
        Model:    &sarvam.SpeechToTextModelSaarikaV2dot5,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Transcript:", result.Transcript)
}
```

### Environment Variable

You can set the `SARVAM_API_KEY` environment variable instead of calling `SetAPIKey()`:

```bash
export SARVAM_API_KEY="your-api-key-here"
```

The SDK will automatically pick up this environment variable on initialization.

### Available Package-Level Functions

The SDK provides the following package-level functions that use the default client:

#### Speech & Audio
- `sarvam.SpeechToText(params)` - Convert speech to text
- `sarvam.SpeechToTextTranslate(params)` - Convert speech to text with translation to English
- `sarvam.TextToSpeech(params)` - Convert text to speech

#### Chat & AI
- `sarvam.ChatCompletion(request)` - Generate chat completions
- `sarvam.SimpleChatCompletion(messages, model)` - Generate chat completions with simplified parameters
- `sarvam.ChatCompletionWithParams(params)` - Generate chat completions with custom parameters

#### Translation & Language
- `sarvam.Translate(input, sourceLang, targetLang)` - Translate text between languages
- `sarvam.TranslateWithParams(params)` - Translate text with advanced parameters
- `sarvam.IdentifyLanguage(input)` - Identify the language of input text
- `sarvam.Transliterate(input, sourceLang, targetLang)` - Convert text between scripts

#### Utility Functions
- `sarvam.SetAPIKey(key)` - Set the API key for the default client
- `sarvam.GetDefaultClient()` - Get the current default client instance

### Error Handling

All package-level functions return appropriate errors if:
- The default client is not initialized (call `SetAPIKey()` first)
- The API key is invalid or expired
- The request parameters are invalid
- The API returns an error response

## 📖 Examples

Check out the [examples](./examples) directory for complete working examples:

- [Text Translation](./examples/text/translate.go)
- [Text-to-Speech](./examples/texttospeech/main.go)
- [Chat Completions](./examples/chatcompletions/chatcompletion.go)
- [Speech-to-Text](./examples/speechtotext/main.go)


## 🤝 Contributing

We welcome contributions! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development Setup

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- [Sarvam AI Platform](https://sarvam.ai)
- [API Documentation](https://docs.sarvam.ai)
- [Go Documentation](https://pkg.go.dev/code.abhai.dev/sarvam)

## 🆘 Support

If you encounter any issues or have questions:

1. Check the [examples](./examples) directory
2. Review the [API documentation](https://docs.sarvam.ai)
3. Open an [issue](https://github.com/abhaikollara/sarvam-go/issues) on GitHub

## Disclaimer

This SDK is an **unofficial** client for the Sarvam API and is not affiliated with, endorsed by, or maintained by Sarvam.

All trademarks, service marks, and copyrights associated with Sarvam belong to their respective owners.

Use this SDK at your own risk. Please review and comply with Sarvam’s terms of service and API usage policies.

---

**Made with ❤️ for the Indian AI community** 
