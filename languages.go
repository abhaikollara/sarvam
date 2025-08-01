package sarvam

// Language represents a supported language code.
type Language string

const (
	LanguageAssamese  Language = "as-IN"
	LanguageBengali   Language = "bn-IN"
	LanguageBodo      Language = "brx-IN"
	LanguageDogri     Language = "doi-IN"
	LanguageEnglish   Language = "en-IN"
	LanguageGujarati  Language = "gu-IN"
	LanguageHindi     Language = "hi-IN"
	LanguageKannada   Language = "kn-IN"
	LanguageKashmiri  Language = "ks-IN"
	LanguageKonkani   Language = "kok-IN"
	LanguageMaithili  Language = "mai-IN"
	LanguageMalayalam Language = "ml-IN"
	LanguageManipuri  Language = "mni-IN"
	LanguageMarathi   Language = "mr-IN"
	LanguageNepali    Language = "ne-IN"
	LanguageOdia      Language = "od-IN"
	LanguagePunjabi   Language = "pa-IN"
	LanguageSanskrit  Language = "sa-IN"
	LanguageSantali   Language = "sat-IN"
	LanguageSindhi    Language = "sd-IN"
	LanguageTamil     Language = "ta-IN"
	LanguageTelugu    Language = "te-IN"
	LanguageUrdu      Language = "ur-IN"
)

const LanguageAuto Language = "auto"

var languageNameMap = map[Language]string{
	LanguageEnglish:   "English",
	LanguageHindi:     "Hindi",
	LanguageBengali:   "Bengali",
	LanguageGujarati:  "Gujarati",
	LanguageKannada:   "Kannada",
	LanguageMalayalam: "Malayalam",
	LanguageMarathi:   "Marathi",
	LanguageOdia:      "Odia",
	LanguagePunjabi:   "Punjabi",
	LanguageSanskrit:  "Sanskrit",
	LanguageSantali:   "Santali",
	LanguageSindhi:    "Sindhi",
	LanguageTamil:     "Tamil",
	LanguageTelugu:    "Telugu",
	LanguageUrdu:      "Urdu",
	LanguageAssamese:  "Assamese",
	LanguageBodo:      "Bodo",
	LanguageDogri:     "Dogri",
	LanguageKashmiri:  "Kashmiri",
	LanguageKonkani:   "Konkani",
	LanguageMaithili:  "Maithili",
	LanguageManipuri:  "Manipuri",
	LanguageNepali:    "Nepali",
}

// String returns the human-readable name of the language.
func (l Language) String() string {
	return languageNameMap[l]
}

var languageMap = map[string]Language{
	"en-IN":  LanguageEnglish,
	"hi-IN":  LanguageHindi,
	"bn-IN":  LanguageBengali,
	"gu-IN":  LanguageGujarati,
	"kn-IN":  LanguageKannada,
	"ml-IN":  LanguageMalayalam,
	"mr-IN":  LanguageMarathi,
	"od-IN":  LanguageOdia,
	"pa-IN":  LanguagePunjabi,
	"ta-IN":  LanguageTamil,
	"te-IN":  LanguageTelugu,
	"ur-IN":  LanguageUrdu,
	"as-IN":  LanguageAssamese,
	"brx-IN": LanguageBodo,
	"doi-IN": LanguageDogri,
	"ks-IN":  LanguageKashmiri,
	"kok-IN": LanguageKonkani,
	"mai-IN": LanguageMaithili,
	"mni-IN": LanguageManipuri,
	"ne-IN":  LanguageNepali,
	"sa-IN":  LanguageSanskrit,
	"sat-IN": LanguageSantali,
	"sd-IN":  LanguageSindhi,
}

// mapLanguageCodeToLanguage converts a language code string to a Language type.
func mapLanguageCodeToLanguage(code string) Language {
	if language, ok := languageMap[code]; ok {
		return language
	}
	return Language(code)
}

// Script represents a writing script.
type Script string

const (
	ScriptLatin      Script = "Latn"
	ScriptDevanagari Script = "Deva"
	ScriptBengali    Script = "Beng"
	ScriptGujarati   Script = "Gujr"
	ScriptKannada    Script = "Knda"
	ScriptMalayalam  Script = "Mlym"
	ScriptOdia       Script = "Orya"
	ScriptGurmukhi   Script = "Guru"
	ScriptTamil      Script = "Taml"
	ScriptTelugu     Script = "Telu"
)

var scriptNameMap = map[Script]string{
	ScriptLatin:      "Latin",
	ScriptDevanagari: "Devanagari",
	ScriptBengali:    "Bengali",
	ScriptGujarati:   "Gujarati",
	ScriptKannada:    "Kannada",
	ScriptMalayalam:  "Malayalam",
	ScriptOdia:       "Odia",
	ScriptGurmukhi:   "Gurmukhi",
	ScriptTamil:      "Tamil",
	ScriptTelugu:     "Telugu",
}

func (s Script) String() string {
	return scriptNameMap[s]
}

var scriptMap = map[string]Script{
	"Latn": ScriptLatin,
	"Deva": ScriptDevanagari,
	"Beng": ScriptBengali,
	"Gujr": ScriptGujarati,
	"Knda": ScriptKannada,
	"Mlym": ScriptMalayalam,
	"Orya": ScriptOdia,
	"Guru": ScriptGurmukhi,
	"Taml": ScriptTamil,
	"Telu": ScriptTelugu,
}

func mapScriptCodeToScript(code string) Script {
	if script, ok := scriptMap[code]; ok {
		return script
	}
	return Script(code)
}
