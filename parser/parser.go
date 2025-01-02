package parser

import (
	"fmt"
	"github.com/structo/types"
	"github.com/structo/utils"
	"regexp"
	"strings"
)

func ParseTypeOrInterface(input, lang string) ([]types.Field, error) {
	// Clean up input - normalize whitespace while preserving necessary spaces
	input = strings.TrimSpace(input)
	input = regexp.MustCompile(`[^\S\n]+`).ReplaceAllString(input, " ")
	isGoStruct := lang == "golang"
	var typeMatches []string
	if isGoStruct {
		// Go struct regex
		structRegex := regexp.MustCompile(`type\s+(\w+)\s+struct\s*{([^}]+)}`)
		typeMatches = structRegex.FindStringSubmatch(input)
	} else {
		// TypeScript interface regex
		interfaceRegex := regexp.MustCompile(`(?:interface|type)\s+(\w+)\s*(?:=)?\s*{([^}]+)}`)
		typeMatches = interfaceRegex.FindStringSubmatch(input)
	}
	if len(typeMatches) < 3 {
		return nil, fmt.Errorf("invalid format: couldn't parse structure")
	}
	fieldSection := typeMatches[2]
	var fieldStrings []string
	if isGoStruct {
		 // Split by newline and filter empty lines
		 lines := strings.Split(fieldSection, "\n")
		 for _, line := range lines {
			 line = strings.TrimSpace(line)
			 if line != "" {
				 fieldStrings = append(fieldStrings, line)
			 }
		 }
	} else {
		// Split TypeScript interface fields by semicolon
		fieldStrings = strings.Split(strings.TrimSpace(fieldSection), ";")
	}

	var fields []types.Field
	// Different regex patterns for Go and TypeScript
	var fieldRegex *regexp.Regexp
	if isGoStruct {
		// Go regex to handle all go types
		fieldRegex = regexp.MustCompile(`^\s*(\w+)\s+(?:\[\])?(\w+(?:\.\w+)?)\s*(?:\x60[^\x60]*\x60)?`)
	} else {
		// Typescript regex to handle all go types
		fieldRegex = regexp.MustCompile(`(\w+)(\?)?:\s*((?:\[\])?\w+(?:\[\])?(?:<\w+>)?)`)
	}
	for i, fieldStr := range fieldStrings {
		fieldStr = strings.TrimSpace(fieldStr)
		if fieldStr == "" {
			continue
		}
		matches := fieldRegex.FindStringSubmatch(fieldStr)
		if len(matches) < 3 {
			continue
		}
		var fieldName, fieldType string
		var isOptional bool
		if isGoStruct {
			fieldName = matches[1]
			fieldType = matches[2]
			// In Go, pointer types are considered optional
			isOptional = strings.HasPrefix(fieldType, "*")
			if isOptional {
				fieldType = strings.TrimPrefix(fieldType, "*")
			}
			if strings.Contains(fieldStr, "[]") {
				fieldType = "[]" + fieldType
			}
			fmt.Printf("Field: %s, Name: %s, Type: %s\n", fieldStr, fieldName, fieldType)
		} else {
			fieldName = matches[1]
			isOptional = matches[2] == "?"
			fieldType = matches[3]
		}
		ok := utils.IsSupportedLanguage(lang)
		if !ok {
			return nil, fmt.Errorf("unsupported language %s", lang)
		}
		validTypes := utils.HandleLangType(lang)
		ok = utils.IsValidType(fieldType, validTypes)
		if !ok {
			return nil, fmt.Errorf("unsupported type %s", fieldType)
		}
		fields = append(fields, types.Field{
			Index:      i,
			Name:       fieldName,
			Type:       fieldType,
			IsOptional: isOptional,
		})
	}
	if len(fields) == 0 {
		return nil, fmt.Errorf("no valid fields found in interface or type")
	}

	return fields, nil
}
