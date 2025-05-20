package langdetect

import (
	"testing"
)

func TestDetect(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected string
		wantErr  bool
	}{
		{
			name:     "English",
			text:     "This is a sample text in English language to test the detector.",
			expected: English,
			wantErr:  false,
		},
		{
			name:     "Hindi",
			text:     "यह हिंदी में एक नमूना पाठ है जिसका उपयोग डिटेक्टर का परीक्षण करने के लिए किया जाता है।",
			expected: Hindi,
			wantErr:  false,
		},
		{
			name:     "Tamil",
			text:     "இது டெடெக்டரை சோதிக்க தமிழில் உள்ள மாதிரி உரை.",
			expected: Tamil,
			wantErr:  false,
		},
		{
			name:     "Telugu",
			text:     "ఇది డిటెక్టర్‌ను పరీక్షించడానికి తెలుగులో నమూనా వచనం.",
			expected: Telugu,
			wantErr:  false,
		},
		{
			name:     "Kannada",
			text:     "ಇದು ಡಿಟೆಕ್ಟರ್ ಅನ್ನು ಪರೀಕ್ಷಿಸಲು ಕನ್ನಡದಲ್ಲಿ ಮಾದರಿ ಪಠ್ಯ.",
			expected: Kannada,
			wantErr:  false,
		},
		{
			name:     "Malayalam",
			text:     "ഇത് ഡിറ്റക്റ്റർ പരിശോധിക്കാൻ മലയാളത്തിലുള്ള സാമ്പിൾ ടെക്സ്റ്റ് ആണ്.",
			expected: Malayalam,
			wantErr:  false,
		},
		{
			name:     "Bengali",
			text:     "এটি ডিটেক্টর পরীক্ষা করার জন্য বাংলায় একটি নমুনা পাঠ্য।",
			expected: Bengali,
			wantErr:  false,
		},
		{
			name:     "Gujarati",
			text:     "આ ડિટેક્ટરનું પરીક્ષણ કરવા માટે ગુજરાતીમાં નમૂના ટેક્સ્ટ છે.",
			expected: Gujarati,
			wantErr:  false,
		},
		{
			name:     "Marathi",
			text:     "हा डिटेक्टरची चाचणी घेण्यासाठी मराठीतील नमुना मजकूर आहे.",
			expected: Marathi,
			wantErr:  false,
		},
		{
			name:     "Text too short",
			text:     "Hello",
			expected: Unknown,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lang, err := Detect(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("Detect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && lang != tt.expected {
				t.Errorf("Detect() = %v, want %v", lang, tt.expected)
			}
		})
	}
}

func TestDetectIndic(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected bool
		wantErr  bool
	}{
		{
			name:     "English",
			text:     "This is a sample text in English language to test the detector.",
			expected: false,
			wantErr:  false,
		},
		{
			name:     "Hindi",
			text:     "यह हिंदी में एक नमूना पाठ है जिसका उपयोग डिटेक्टर का परीक्षण करने के लिए किया जाता है।",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "Tamil",
			text:     "இது டெடெக்டரை சோதிக்க தமிழில் உள்ள மாதிரி உரை.",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "Text too short",
			text:     "Hello",
			expected: false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isIndic, err := DetectIndic(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectIndic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && isIndic != tt.expected {
				t.Errorf("DetectIndic() = %v, want %v", isIndic, tt.expected)
			}
		})
	}
}

func TestIsIndic(t *testing.T) {
	tests := []struct {
		name     string
		langCode string
		expected bool
	}{
		{name: "English", langCode: English, expected: false},
		{name: "Hindi", langCode: Hindi, expected: true},
		{name: "Tamil", langCode: Tamil, expected: true},
		{name: "Telugu", langCode: Telugu, expected: true},
		{name: "Kannada", langCode: Kannada, expected: true},
		{name: "Malayalam", langCode: Malayalam, expected: true},
		{name: "Bengali", langCode: Bengali, expected: true},
		{name: "Gujarati", langCode: Gujarati, expected: true},
		{name: "Marathi", langCode: Marathi, expected: true},
		{name: "Unknown", langCode: Unknown, expected: false},
		{name: "Invalid", langCode: "xx", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsIndic(tt.langCode)
			if result != tt.expected {
				t.Errorf("IsIndic(%s) = %v, want %v", tt.langCode, result, tt.expected)
			}
		})
	}
}

func TestSupportedLanguages(t *testing.T) {
	languages := SupportedLanguages()
	
	// Check that all language constants are included
	expectedLangs := map[string]bool{
		English:   true,
		Hindi:     true,
		Tamil:     true,
		Telugu:    true,
		Kannada:   true,
		Malayalam: true,
		Bengali:   true,
		Gujarati:  true,
		Marathi:   true,
	}
	
	if len(languages) != len(expectedLangs) {
		t.Errorf("SupportedLanguages() returned %d languages, expected %d", len(languages), len(expectedLangs))
	}
	
	for _, lang := range languages {
		if !expectedLangs[lang] {
			t.Errorf("Unexpected language in SupportedLanguages(): %s", lang)
		}
	}
}

// Let's create a benchmark to measure performance
func BenchmarkDetect(b *testing.B) {
	englishText := "This is a sample text in English language to test the detector."
	hindiText := "यह हिंदी में एक नमूना पाठ है जिसका उपयोग डिटेक्टर का परीक्षण करने के लिए किया जाता है।"
	
	b.Run("English", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Detect(englishText)
		}
	})
	
	b.Run("Hindi", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Detect(hindiText)
		}
	})
} 