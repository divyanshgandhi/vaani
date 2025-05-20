package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/NavoDayAI/vaani/backend/services/api-gateway/utils/langdetect"
	"go.uber.org/zap"
)

// LanguageDetectionRequest represents a request to detect language
type LanguageDetectionRequest struct {
	Text string `json:"text"`
}

// LanguageDetectionResponse represents a language detection response
type LanguageDetectionResponse struct {
	Language      string `json:"language"`
	LanguageName  string `json:"language_name"`
	IsIndic       bool   `json:"is_indic"`
	TextLength    int    `json:"text_length"`
	SupportedLangs []string `json:"supported_languages,omitempty"`
}

// LanguageMap maps language codes to human-readable language names
var LanguageMap = map[string]string{
	langdetect.English:   "English",
	langdetect.Hindi:     "Hindi",
	langdetect.Tamil:     "Tamil",
	langdetect.Telugu:    "Telugu",
	langdetect.Kannada:   "Kannada",
	langdetect.Malayalam: "Malayalam",
	langdetect.Bengali:   "Bengali",
	langdetect.Gujarati:  "Gujarati",
	langdetect.Marathi:   "Marathi",
	langdetect.Unknown:   "Unknown",
}

// DetectLanguageHandler handles requests to detect the language of text
func DetectLanguageHandler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			logger.Error("Error reading request body", zap.Error(err))
			return
		}
		defer r.Body.Close()

		// Parse request
		var req LanguageDetectionRequest
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			logger.Error("Error parsing JSON", zap.Error(err))
			return
		}

		// Validate request
		if req.Text == "" {
			http.Error(w, "Text is required", http.StatusBadRequest)
			return
		}

		// Detect language
		language, err := langdetect.Detect(req.Text)
		
		// Prepare response
		response := LanguageDetectionResponse{
			Language:     language,
			LanguageName: LanguageMap[language],
			TextLength:   len(req.Text),
			SupportedLangs: langdetect.SupportedLanguages(),
		}
		
		// Add error info if needed
		if err != nil {
			if err == langdetect.ErrTextTooShort {
				// Still return the response with the error, but add a note
				response.Language = langdetect.Unknown
				response.LanguageName = "Text too short for reliable detection"
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				logger.Error("Language detection error", zap.Error(err))
				return
			}
		} else {
			// Determine if the detected language is Indic
			response.IsIndic = langdetect.IsIndic(language)
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		
		// Log success
		logger.Info("Language detected", 
			zap.String("language", language),
			zap.Int("text_length", len(req.Text)),
			zap.Bool("is_indic", response.IsIndic))
	}
} 