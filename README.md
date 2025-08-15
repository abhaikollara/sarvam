# SarvamAI Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/abhaikollara/sarvam.svg)](https://pkg.go.dev/github.com/abhaikollara/sarvam)
[![Go Report Card](https://goreportcard.com/badge/abhaikollara/sarvam)](https://goreportcard.com/report/abhaikollara/sarvam)
[![Go Version](https://img.shields.io/github/go-mod/go-version/abhaikollara/sarvam)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

An unofficial Go SDK for the [Sarvam AI](https://sarvam.ai) APIs, providing easy access to Indian language AI services including translation, text-to-speech, chat completions, and language identification.

> **⚠️ Breaking Changes Notice**: This SDK is currently in pre-v1 development. Breaking changes may occur until v1.0.0 is released. Please pin your dependency version if you need stability in production environments.

## 🌟 Features

- **Text Translation** - Translate text between 22+ Indian languages
- **Language Identification** - Automatically detect the language of input text
- **Text-to-Speech** - Convert text to natural-sounding speech
- **Speech-to-Text** - Transcribe audio to text with support for Indian languages
- **Speech-to-Text Translation** - Transcribe and translate audio in one request
- **Chat Completions** - Generate AI responses using Sarvam's language models
- **Transliteration** - Convert text between different writing scripts

## 🚀 Quick Start

### Installation

```bash
go get code.abhai.dev/sarvam
```

### Basic Usage

The SDK provides both instance-based and package-level APIs for convenience.


### Environment Variable

You can set the `SARVAM_API_KEY` environment variable instead of calling `SetAPIKey()`:

```bash
export SARVAM_API_KEY="your-api-key-here"
```

The SDK will automatically pick up this environment variable on initialization.

## 📖 Examples

Check out the [examples](./examples) directory for complete working examples:

- [Text Translation](./examples/text/translate.go)
- [Text-to-Speech](./examples/texttospeech/main.go)
- [Chat Completions](./examples/chatcompletions/chatcompletion.go)
- [Speech-to-Text](./examples/speechtotext/main.go)
- [Speech-to-Text Translation](./examples/speechtotexttranslate/main.go)
- [Language Identification](./examples/languageidentification/main.go)
- [Transliteration](./examples/transliteratetext/main.go)


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

Use this SDK at your own risk. Please review and comply with Sarvam's terms of service and API usage policies.

---

**Made with ❤️ for the Indian AI community** 
