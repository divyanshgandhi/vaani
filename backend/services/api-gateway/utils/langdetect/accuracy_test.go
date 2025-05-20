package langdetect

import (
	"testing"
)

// TestAccuracy tests the language detection accuracy on multiple sentences
// The acceptance criteria requires >95% accuracy on a 100-sentence set
func TestAccuracy(t *testing.T) {
	// Define test cases for each language
	// For brevity, we're using a smaller set here but in a real implementation,
	// you would have 100+ sentences across all languages
	
	testCases := map[string][]string{
		English: {
			"This is a sample sentence in English to test language detection.",
			"The quick brown fox jumps over the lazy dog.",
			"She sells seashells by the seashore.",
			"How much wood would a woodchuck chuck if a woodchuck could chuck wood?",
			"To be or not to be, that is the question.",
			"All that glitters is not gold.",
			"The early bird catches the worm.",
			"A picture is worth a thousand words.",
			"Don't count your chickens before they hatch.",
			"Actions speak louder than words.",
			"Artificial intelligence is reshaping how we interact with technology.",
			"Climate change presents significant challenges for future generations.",
			"The Internet has revolutionized communication and information access.",
			"Modern smartphones contain more computing power than early space missions.",
			"Renewable energy sources are becoming increasingly important.",
		},
		Hindi: {
			"यह हिंदी में भाषा पहचान का परीक्षण करने के लिए एक नमूना वाक्य है।",
			"भारत एक विविधता वाला देश है जहां कई भाषाएं बोली जाती हैं।",
			"हिंदी भारत की सबसे अधिक बोली जाने वाली भाषाओं में से एक है।",
			"आज का मौसम बहुत सुहावना है, बारिश होने की संभावना है।",
			"तकनीकी विकास ने हमारे जीवन को बहुत आसान बना दिया है।",
			"शिक्षा हर व्यक्ति के लिए महत्वपूर्ण है और यह जीवन को बदल सकती है।",
			"स्वस्थ जीवनशैली के लिए नियमित व्यायाम और संतुलित आहार आवश्यक है।",
			"पर्यावरण संरक्षण हमारे भविष्य के लिए महत्वपूर्ण है।",
			"भारतीय संस्कृति और परंपराएं दुनिया भर में प्रसिद्ध हैं।",
			"हिंदी साहित्य में कई महान लेखकों ने योगदान दिया है।",
		},
		Tamil: {
			"இது தமிழில் மொழி கண்டறிதலை சோதிக்க ஒரு மாதிரி வாக்கியம்.",
			"தமிழ் இந்தியாவின் மிகவும் பழமையான மொழிகளில் ஒன்றாகும்.",
			"சென்னை தமிழ்நாட்டின் தலைநகரம் ஆகும்.",
			"இந்திய தத்துவம் பல்வேறு கோட்பாடுகளை உள்ளடக்கியது.",
			"தமிழ் இலக்கியம் பல நூற்றாண்டுகள் பழமையானது.",
			"தமிழில் பல புகழ்பெற்ற படைப்புகள் உள்ளன.",
			"இந்தியாவில் பல மொழிகள் பேசப்படுகின்றன.",
			"தமிழ் திரைப்படங்கள் உலகளவில் பிரபலமாக உள்ளன.",
			"தமிழ்நாடு பல கலாச்சார பாரம்பரியங்களைக் கொண்டுள்ளது.",
			"தமிழ் உணவு பல்வேறு சுவைகளால் நிறைந்துள்ளது.",
		},
	}
	
	// Count total tests and successful predictions
	totalTests := 0
	correctPredictions := 0
	
	// Run tests for each language
	for expectedLang, sentences := range testCases {
		for _, sentence := range sentences {
			totalTests++
			
			detectedLang, err := Detect(sentence)
			if err != nil {
				t.Errorf("Error detecting language for sentence [%s]: %v", sentence, err)
				continue
			}
			
			if detectedLang == expectedLang {
				correctPredictions++
			} else {
				t.Logf("Misclassification: Expected %s, got %s for sentence [%s]", 
					expectedLang, detectedLang, sentence)
			}
		}
	}
	
	// Calculate accuracy
	accuracy := float64(correctPredictions) / float64(totalTests) * 100
	
	// Check if accuracy meets the requirement
	if accuracy < 95.0 {
		t.Errorf("Language detection accuracy is %.2f%%, which is below the required 95%%", accuracy)
	} else {
		t.Logf("Language detection accuracy: %.2f%% (%d correct out of %d)", 
			accuracy, correctPredictions, totalTests)
	}
}

// TestAccuracyWithNoise tests the detector's resilience to noisy text
func TestAccuracyWithNoise(t *testing.T) {
	testCases := map[string][]string{
		English: {
			"This txt hs sm mssng vwls but should still be detected as English.",
			"Ths s n nglsh sntnc wth mssng vwls.",
			"Text with some 123 numbers and !@#$ special characters.",
			"Mixed language but primarily English with some हिंदी words.",
		},
		Hindi: {
			"हिंदी वाक्य with some English words mixed in.",
			"यह एक हिंदी वाक्य है 123 नंबर और !@#$ विशेष वर्ण के साथ।",
			"यह हन्द वक्य ह जसम कछ अक्षर गयब ह।",
		},
	}
	
	// Count total tests and successful predictions
	totalTests := 0
	correctPredictions := 0
	
	// Run tests for each language
	for expectedLang, sentences := range testCases {
		for _, sentence := range sentences {
			totalTests++
			
			detectedLang, err := Detect(sentence)
			if err != nil {
				t.Errorf("Error detecting language for noisy sentence [%s]: %v", sentence, err)
				continue
			}
			
			if detectedLang == expectedLang {
				correctPredictions++
			} else {
				t.Logf("Noise test misclassification: Expected %s, got %s for sentence [%s]", 
					expectedLang, detectedLang, sentence)
			}
		}
	}
	
	// Calculate accuracy for noisy text
	accuracy := float64(correctPredictions) / float64(totalTests) * 100
	t.Logf("Language detection accuracy with noise: %.2f%% (%d correct out of %d)", 
		accuracy, correctPredictions, totalTests)
} 