package sarvam

import "testing"

func TestMapLanguageCodeToLanguage(t *testing.T) {
	tests := []struct {
		code string
		want Language
	}{
		{"as-IN", LanguageAssamese},
		{"bn-IN", LanguageBengali},
		{"brx-IN", LanguageBodo},
		{"doi-IN", LanguageDogri},
		{"en-IN", LanguageEnglish},
		{"gu-IN", LanguageGujarati},
		{"hi-IN", LanguageHindi},
		{"kn-IN", LanguageKannada},
		{"ks-IN", LanguageKashmiri},
		{"kok-IN", LanguageKonkani},
		{"mai-IN", LanguageMaithili},
		{"ml-IN", LanguageMalayalam},
		{"mni-IN", LanguageManipuri},
		{"mr-IN", LanguageMarathi},
		{"ne-IN", LanguageNepali},
		{"od-IN", LanguageOdia},
		{"pa-IN", LanguagePunjabi},
		{"sa-IN", LanguageSanskrit},
		{"sat-IN", LanguageSantali},
		{"sd-IN", LanguageSindhi},
		{"ta-IN", LanguageTamil},
		{"te-IN", LanguageTelugu},
		{"ur-IN", LanguageUrdu},
	}

	for _, test := range tests {
		got := mapLanguageCodeToLanguage(test.code)
		if got != test.want {
			t.Errorf("mapLanguageCodeToLanguage(%q) = %v, want %v", test.code, got, test.want)
		}
	}

	unknownCode := "random-code"
	got := mapLanguageCodeToLanguage(unknownCode)
	if got != Language(unknownCode) {
		t.Errorf("mapLanguageCodeToLanguage(%q) = %v, want %v", unknownCode, got, LanguageAuto)
	}
}
