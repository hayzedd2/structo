package utils

import(
	"strings"
)

func IsValidType(fieldType string, validTypes []string) bool {
	for _, v := range validTypes {
		if v == fieldType {
			return true
		}
	}
	return false
}


func HandleLangType(lang string) []string {
	// Valid TypeScript types
	validTypescriptTypes := []string{
		"string",
		"string | null",
		"string | undefined",
		"number | null",
		"number | undefined",
		"number",
		"boolean",
		"string[]",
		"number[]",
		"boolean[]",
		"Date",
		"any",
	}

	// Valid Go types
	validGoTypes := []string{
		"string",
		"bool",
		"int",
		"int8",
		"int16",
		"int32",
		"int64",
		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"float32",
		"float64",
		"[]string",
		"[]int",
		"[]int8",
		"[]int16",
		"[]int32",
		"[]int64",
		"[]uint",
		"[]uint8",
		"[]uint16",
		"[]uint32",
		"[]uint64",
		"[]bool",
		"[]float32",
		"[]float64",
	}

	lowerCaseLang := strings.ToLower(lang)

	// Choose the correct type slice
	switch lowerCaseLang {
	case "typescript":
		return validTypescriptTypes
	case "golang", "go":
		return validGoTypes
	default:
		return nil // Unsupported language
	}
}

func IsSupportedLanguage(lang string) bool {
	switch strings.ToLower(lang) {
	case "typescript", "golang", "go":
		return true
	default:
		return false
	}
}