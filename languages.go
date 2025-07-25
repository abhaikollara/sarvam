package sarvam

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

func mapLanguageCodeToLanguage(code string) Language {
	if language, ok := languageMap[code]; ok {
		return language
	}
	return Language(code)
}
