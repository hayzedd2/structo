package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/structo/types"
	"github.com/structo/utils"
)



func ParseTypeOrInterface(input string) ([]types.Field, error) {
	// Clean up input - normalize whitespace while preserving necessary spaces
	input = strings.TrimSpace(input)
	input = regexp.MustCompile(`\s+`).ReplaceAllString(input, " ")

	typeOrInterfaceRegex := regexp.MustCompile(`(?:interface|type)\s+(\w+)\s*(?:=)?\s*{([^}]+)}`)

	matches := typeOrInterfaceRegex.FindStringSubmatch(input)
	if len(matches) < 3 {
		return nil, fmt.Errorf("invalid type or interface format: couldn't parse structure")
	}
	fieldSection := matches[2]

	// Split fields by semicolon, handling possible trailing semicolon
	fieldStrings := strings.Split(strings.TrimSpace(fieldSection), ";")
	var fields []types.Field

	fieldRegex := regexp.MustCompile(`(\w+)(\?)?:\s*((?:\[\])?\w+(?:\[\])?(?:<\w+>)?)`)
	for _, fieldStr := range fieldStrings {
		fieldStr = strings.TrimSpace(fieldStr)
		if fieldStr == "" {
			continue
		}
		matches := fieldRegex.FindStringSubmatch(fieldStr)
		if len(matches) < 4 {
			continue
		}
		fieldName := matches[1]
		isOptional := matches[2] == "?"
		fieldType := matches[3]
		validTypes := utils.HandleLangType("typescript")
		ok := utils.IsValidType(fieldType, validTypes)
		if !ok {
			return nil, fmt.Errorf("unsupported type %s", fieldType)
		}
		fields = append(fields, types.Field{
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
