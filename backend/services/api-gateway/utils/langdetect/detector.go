package langdetect

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Supported languages
const (
	English   = "en"
	Hindi     = "hi"
	Tamil     = "ta"
	Telugu    = "te"
	Kannada   = "kn"
	Malayalam = "ml"
	Bengali   = "bn"
	Gujarati  = "gu"
	Marathi   = "mr"
	Unknown   = "unknown"
)

var (
	// ErrTextTooShort is returned when the text is too short for reliable detection
	ErrTextTooShort = errors.New("text too short for reliable language detection")

	// ErrUnsupportedLanguage is returned when the detected language is not in our supported list
	ErrUnsupportedLanguage = errors.New("detected language not supported")
)

// Marathi-specific characters and character combinations
var marathiCharacters = []string{
	"ळ", "ऌ", "ऍ", "ऎ", "ख्र", "त्र्य", "ज्ञ", "ट्र", "ड्र", "क्ष",
}

// Hindi-specific characters and character combinations
var hindiCharacters = []string{
	"ऋ", "ॠ", "ऌ", "ॡ", "ड़", "ढ़", "क़", "ख़", "ग़", "ज़", "फ़",
}

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
			// Both Hindi and Marathi use Devanagari
			// We'll differentiate them later
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

	// Additional check for Hindi/Marathi distinction
	if maxLang == Hindi {
		if isLikelyMarathi(text) {
			return Marathi
		}
	}

	return maxLang
}

// isLikelyMarathi attempts to determine if Devanagari text is more likely to be Marathi than Hindi
func isLikelyMarathi(text string) bool {
	// Count specific Marathi and Hindi character occurrences
	marathiScore := 0
	hindiScore := 0

	// Check for Marathi-specific characters
	for _, char := range marathiCharacters {
		count := strings.Count(text, char)
		marathiScore += count
	}

	// Check for Hindi-specific characters
	for _, char := range hindiCharacters {
		count := strings.Count(text, char)
		hindiScore += count
	}

	// Marathi-specific word endings
	marathiEndings := []string{"मध्ये", "साठी", "तील", "मुळे", "चा", "ची", "चे", "ला", "हे", "आहे"}
	for _, ending := range marathiEndings {
		if strings.Contains(text, ending) {
			marathiScore += 2
		}
	}

	// Hindi-specific word endings
	hindiEndings := []string{"का", "के", "की", "में", "है", "से", "कर", "को", "ने", "पर"}
	for _, ending := range hindiEndings {
		if strings.Contains(text, ending) {
			hindiScore += 2
		}
	}

	// Return true if Marathi score is higher or equal to Hindi score
	return marathiScore >= hindiScore && marathiScore > 0
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
	if utf8.RuneCountInString(text) < 10 {
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
