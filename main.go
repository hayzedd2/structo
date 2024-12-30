package main

import (
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"math/rand"
	"regexp"
	"strings"
)

type Field struct {
	Name       string
	Type       string
	IsOptional bool
}

func main() {
	input := `export interface EventResponse {
    ID: number;
    Name: string;
    Description: string;
    Location: string;
    StartDate: string;
    StartTime: string
    Category: string
    UserId: string;
}
`
	fields, err := parseTypeOrInterface(input)
	mockDataMap := make(map[string]interface{})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, field := range fields {
		mockData := generateMockData(field)
		mockDataMap[field.Name] = mockData
	}
	jsonData, err := json.MarshalIndent(mockDataMap, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	fmt.Println("Generated JSON Schema:")
	fmt.Println(string(jsonData))
}

func parseTypeOrInterface(input string) ([]Field, error) {

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

	var fields []Field

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
		validTypes := handleLangType("typescript")
		ok := isValidType(fieldType, validTypes)
		if !ok {
			return nil, fmt.Errorf("unsupported type %s", fieldType)
		}
		fields = append(fields, Field{
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

func generateMockData(field Field) interface{} {
	faker := gofakeit.New(0)
	parts := strings.Split(strings.ToLower(field.Name), "_")
	var fieldName = strings.Join(parts, "")

	generateString := func(fieldName string) interface{} {
		switch {
		// Personal Information
		case contains(fieldName, "firstname", "fname"):
			return faker.FirstName()
		case contains(fieldName, "lastname", "lname"):
			return faker.LastName()
		case contains(fieldName, "username"):
			return faker.Username()
		case contains(fieldName, "fullname", "name"):
			return faker.Name()
		case contains(fieldName, "email"):
			return faker.Email()
		case contains(fieldName, "username"):
			return faker.Username()
		case contains(fieldName, "password"):
			return faker.Password(true, true, true, true, false, 10)

		// Identification
		case contains(fieldName, "uuid", "id"):
			return faker.UUID()
		case contains(fieldName, "ssn"):
			return faker.SSN()

		// Contact
		case contains(fieldName, "phone", "telephone", "mobile"):
			return faker.Phone()

		// Location
		case contains(fieldName, "country", "location"):
			return faker.Country()
		case contains(fieldName, "city"):
			return faker.City()
		case contains(fieldName, "state"):
			return faker.State()
		case contains(fieldName, "street"):
			return faker.Street()
		case contains(fieldName, "address"):
			return faker.Address().Address
		case contains(fieldName, "zipcode", "postalcode"):
			return faker.Zip()
		case contains(fieldName, "latitude", "lat"):
			return fmt.Sprintf("%.6f", faker.Latitude())
		case contains(fieldName, "longitude", "lng", "long"):
			return fmt.Sprintf("%.6f", faker.Longitude())

		// Internet
		case contains(fieldName, "url", "website"):
			return faker.URL()
		case contains(fieldName, "ipv4"):
			return faker.IPv4Address()
		case contains(fieldName, "ipv6"):
			return faker.IPv6Address()
		case contains(fieldName, "useragent"):
			return faker.UserAgent()

		// Business
		case contains(fieldName, "company", "business"):
			return faker.Company()
		case contains(fieldName, "jobTitle", "job"):
			return faker.JobTitle()
		case contains(fieldName, "product"):
			return faker.Product()

		// Payment
		case contains(fieldName, "creditcard", "cc"):
			return faker.CreditCardNumber
		case contains(fieldName, "currency"):
			return faker.Currency().Short

		// Dates
		case contains(fieldName, "date", "createdat", "updatedat"):
			return faker.Date().Format("2006-01-02")
		case contains(fieldName, "time"):
			return faker.Date().Format("15:04:05")

		// Content
		case contains(fieldName, "description", "desc"):
			return faker.Sentence(10)
		case contains(fieldName, "title"):
			return faker.Sentence(3)
		case contains(fieldName, "comment"):
			return faker.Sentence(5)
		case contains(fieldName, "paragraph"):
			return faker.Paragraph(2, 4, 4, "\n")

		// Colors
		case contains(fieldName, "color"):
			return faker.Color()
		case contains(fieldName, "hexcolor"):
			return faker.HexColor()

		// File-related

		case contains(fieldName, "extension", "ext"):
			return faker.FileExtension()

		// Miscellaneous
		case contains(fieldName, "language", "lang"):
			return faker.Language()
		case contains(fieldName, "timezone", "tz"):
			return faker.TimeZone()
		default:
			return faker.LoremIpsumSentence
		}
	}
	switch field.Type {
	case "string":
		return generateString(fieldName)
	case "string[]":
		count := rand.Intn(3) + 2
		arr := make([]any, count)
		for i := range arr {
			arr[i] = generateString(fieldName)
		}
		return arr
	case "number":
		return faker.Number(1, 1000)
	case "number[]":
		count := rand.Intn(3) + 2
		arr := make([]int, count)
		for i := range arr {
			arr[i] = faker.Number(1, 1000)
		}
		return arr
	case "bool", "boolean":
		return faker.Bool()
	case "Date":
		return faker.Date()
	case "[]Date":
		count := rand.Intn(3) + 2
		arr := make([]any, count)
		for i := range arr {
			arr[i] = faker.Date()
		}
		return arr
	default:
		return nil
	}
}

func isValidType(fieldType string, validTypes []string) bool {
	for _, v := range validTypes {
		if v == fieldType {
			return true
		}
	}
	return false
}

func handleLangType(lang string) []string {
	// Valid TypeScript types
	validTypescriptTypes := []string{
		"string",
		"number",
		"boolean",
		"string[]",
		"number[]",
		"boolean[]",
		"Date",
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

	normalizedLang := strings.ToLower(lang)

	// Choose the correct type slice
	switch normalizedLang {
	case "typescript":
		return validTypescriptTypes
	case "golang", "go":
		return validGoTypes
	default:
		return nil // Unsupported language
	}
}

func contains(str string, patterns ...string) bool {
	for _, pattern := range patterns {
		if strings.Contains(str, pattern) {
			return true
		}
	}
	return false
}
