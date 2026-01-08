package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
)

// Init initializes the i18n bundle with default language
func Init(defaultLang string) error {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Load translation files
	langs := []string{"en", "es", "fr", "de", "ja", "zh"}
	for _, lang := range langs {
		filePath := fmt.Sprintf("locales/%s.json", lang)
		if _, err := os.Stat(filePath); err == nil {
			if _, err := bundle.LoadMessageFile(filePath); err != nil {
				return fmt.Errorf("failed to load locale file %s: %w", filePath, err)
			}
		}
	}

	// Set default language
	defaultTag := language.Make(defaultLang)
	if defaultTag == language.Und {
		defaultTag = language.English
	}

	localizer = i18n.NewLocalizer(bundle, defaultTag.String())
	return nil
}

// GetLocalizer returns a localizer for the given language
func GetLocalizer(lang string) *i18n.Localizer {
	if lang == "" {
		return localizer
	}

	tag := language.Make(lang)
	if tag == language.Und {
		return localizer
	}

	return i18n.NewLocalizer(bundle, tag.String())
}

// GetLanguageFromRequest extracts language from request header or query param
func GetLanguageFromRequest(c *gin.Context) string {
	// Check Accept-Language header
	acceptLang := c.GetHeader("Accept-Language")
	if acceptLang != "" {
		langs := strings.Split(acceptLang, ",")
		if len(langs) > 0 {
			lang := strings.TrimSpace(strings.Split(langs[0], ";")[0])
			return lang
		}
	}

	// Check query parameter
	lang := c.Query("lang")
	if lang != "" {
		return lang
	}

	// Default to English
	return "en"
}

// T translates a message ID with optional arguments
func T(lang, messageID string, args ...interface{}) string {
	loc := GetLocalizer(lang)
	if loc == nil {
		return messageID
	}

	msg, err := loc.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: args,
	})

	if err != nil {
		return messageID
	}

	return msg
}

// TWithDefault translates with a default message
func TWithDefault(lang, messageID, defaultMessage string, args ...interface{}) string {
	loc := GetLocalizer(lang)
	if loc == nil {
		return defaultMessage
	}

	msg, err := loc.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
		DefaultMessage: &i18n.Message{
			ID:    messageID,
			Other: defaultMessage,
		},
		TemplateData: args,
	})

	if err != nil {
		return defaultMessage
	}

	return msg
}

// GetBundle returns the i18n bundle
func GetBundle() *i18n.Bundle {
	return bundle
}
