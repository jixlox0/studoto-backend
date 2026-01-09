package uuid

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Generate generates a UUID with a custom prefix
func Generate(prefix string) string {
	id := uuid.New().String()

	// Remove hyphens from UUID
	id = strings.ReplaceAll(id, "-", "")

	if prefix == "" {
		return id
	}

	// Return prefixed UUID: PREFIX-UUID
	return fmt.Sprintf("%s-%s", strings.ToLower(prefix), id)
}

// GenerateShort generates a short UUID with prefix (first 8 chars)
func GenerateShort(prefix string) string {
	id := uuid.New().String()
	id = strings.ReplaceAll(id, "-", "")

	if prefix == "" {
		return id[:8]
	}

	return fmt.Sprintf("%s-%s", strings.ToLower(prefix), id[:8])
}

// GenerateWithSeparator generates UUID with custom prefix and separator
func GenerateWithSeparator(prefix, separator string) string {
	id := uuid.New().String()
	id = strings.ReplaceAll(id, "-", "")

	if separator == "" {
		separator = "-"
	}

	if prefix == "" {
		return id
	}

	return fmt.Sprintf("%s%s%s", strings.ToLower(prefix), separator, id)
}

// Validate validates if a string is a valid prefixed UUID
func Validate(id, prefix string) bool {
	if prefix == "" {
		return uuid.Validate(id) == nil
	}

	expectedPrefix := strings.ToLower(prefix) + "-"
	if !strings.HasPrefix(id, expectedPrefix) {
		return false
	}

	uuidPart := strings.TrimPrefix(id, expectedPrefix)
	return len(uuidPart) == 32 // UUID without hyphens is 32 chars
}

// ExtractUUID extracts the UUID part from a prefixed UUID
func ExtractUUID(prefixedID, prefix string) (string, error) {
	if prefix == "" {
		return prefixedID, nil
	}

	expectedPrefix := strings.ToLower(prefix) + "-"
	if !strings.HasPrefix(prefixedID, expectedPrefix) {
		return "", fmt.Errorf("invalid prefix: expected %s", expectedPrefix)
	}

	uuidPart := strings.TrimPrefix(prefixedID, expectedPrefix)
	if len(uuidPart) != 32 {
		return "", fmt.Errorf("invalid UUID length")
	}

	// Add hyphens back for standard UUID format
	standardUUID := fmt.Sprintf("%s-%s-%s-%s-%s",
		uuidPart[0:8],
		uuidPart[8:12],
		uuidPart[12:16],
		uuidPart[16:20],
		uuidPart[20:32],
	)

	return standardUUID, nil
}
