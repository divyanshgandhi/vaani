package langdetect

import (
	"errors"
	"strings"
	"unicode"
)

// Supported languages
const (
	English  = "en"
	Hindi    = "hi"
	Tamil    = "ta"
	Telugu   = "te"
	Kannada  = "kn"
	Malayalam = "ml"
	Bengali  = "bn"
	Gujarati = "gu"
	Marathi  = "mr"
	Unknown  = "unknown"
)

var (
	// ErrTextTooShort is returned when the text is too short for reliable detection
	ErrTextTooShort = errors.New("text too short for reliable language detection")
	
	// ErrUnsupportedLanguage is returned when the detected language is not in our supported list
	ErrUnsupportedLanguage = errors.New("detected language not supported")
)

// Simple language detection based on Unicode ranges
func detectLanguageByUnicodeRange(text string) string {
	if len(text) == 0 {
		return Unknown
	}

	// Count characters in different script ranges
	counts := make(map[string]int)
	
	for _, r := range text {
		switch {
		case unicode.Is(unicode.Latin, r):
			counts[English]++
		case unicode.Is(unicode.Devanagari, r):
			counts[Hindi]++
		case unicode.Is(unicode.Tamil, r):
			counts[Tamil]++
		case unicode.Is(unicode.Telugu, r):
			counts[Telugu]++
		case unicode.Is(unicode.Kannada, r):
			counts[Kannada]++
		case unicode.Is(unicode.Malayalam, r):
			counts[Malayalam]++
		case unicode.Is(unicode.Bengali, r):
			counts[Bengali]++
		case unicode.Is(unicode.Gujarati, r):
			counts[Gujarati]++
		// Marathi uses Devanagari script, so it's counted under Hindi
		}
	}
	
	// Find the language with the most characters
	maxCount := 0
	maxLang := Unknown
	
	for lang, count := range counts {
		if count > maxCount {
			maxCount = count
			maxLang = lang
		}
	}
	
	return maxLang
}

// SupportedLanguages returns a list of ISO 639-1 language codes that we support
func SupportedLanguages() []string {
	return []string{
		English,
		Hindi,
		Tamil,
		Telugu,
		Kannada,
		Malayalam,
		Bengali,
		Gujarati,
		Marathi,
	}
}

// Detect detects the language of the given text
// Returns the ISO 639-1 language code (e.g., "en", "hi", "ta")
// Minimum 10 characters are needed for reliable detection
func Detect(text string) (string, error) {
	// Clean the text
	text = strings.TrimSpace(text)
	
	// Check if the text is long enough for reliable detection
	if len(text) < 10 {
		return Unknown, ErrTextTooShort
	}
	
	lang := detectLanguageByUnicodeRange(text)
	if lang == Unknown {
		return Unknown, ErrUnsupportedLanguage
	}
	
	return lang, nil
}

// DetectIndic returns true if the text is in an Indic language
// This is useful for routing TTS requests to the appropriate backend
func DetectIndic(text string) (bool, error) {
	lang, err := Detect(text)
	if err != nil {
		return false, err
	}
	
	// Check if the language is not English
	return lang != English && lang != Unknown, nil
}

// IsIndic checks if the given language code is an Indic language
func IsIndic(langCode string) bool {
	indicLanguages := map[string]bool{
		Hindi:     true,
		Tamil:     true,
		Telugu:    true,
		Kannada:   true,
		Malayalam: true,
		Bengali:   true,
		Gujarati:  true,
		Marathi:   true,
	}
	
	return indicLanguages[langCode]
} 