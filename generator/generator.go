package generator

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/iancoleman/orderedmap"
	"github.com/structo/types"
	"math/rand"
	"sort"
	"strings"
)

var faker = gofakeit.New(0)
type Age struct{
	Age int;
	Addresss  string;
}

func GenerateMockObjects(fields []types.Field, count int) []orderedmap.OrderedMap {
	// Create slice to hold multiple objects
	mockObjects := make([]orderedmap.OrderedMap, count)

	// Generate 'count' number of objects
	for i := 0; i < count; i++ {
		mockDataMap := orderedmap.New()

		// sort by index
		sortedFields := make([]types.Field, len(fields))
		copy(sortedFields, fields)
		sort.Slice(sortedFields, func(i, j int) bool {
			return sortedFields[i].Index < sortedFields[j].Index
		})
		for _, field := range sortedFields {
			mockDataMap.Set(field.Name, GenerateMockData(field))
		}
		mockObjects[i] = *mockDataMap
	}
	return mockObjects
}

func GenerateMockData(field types.Field) interface{} {
    parts := strings.Split(strings.ToLower(field.Name), "_")
    fieldName := strings.Join(parts, "")
    fmt.Println("fieldtypes", field.Type)
    // Extract base type and check if it's an array
    baseType := field.Type
    isArray := strings.HasSuffix(baseType, "[]") || strings.HasPrefix(baseType, "[]")
    if isArray {
        baseType = strings.TrimSuffix(strings.TrimPrefix(baseType, "[]"), "[]")
    }
	fmt.Println("basetype", baseType)

    // Generate single value based on type
    generateValue := func() interface{} {
        switch baseType {
        case "string":
            return generateMeaningfulString(fieldName)
            
        case "number", "int", "int8", "int16", "int32", "int64", "Number":
            return faker.Number(1, 1000)
            
        case "uint", "uint8", "uint16", "uint32", "uint64":
            return faker.Uint64()
            
        case "float32", "Float32":
            return faker.Float32()
            
        case "float64", "Float64":
            return faker.Float64()
            
        case "bool", "boolean", "Boolean":
            return faker.Bool()
            
        case "Date", "date":
            return faker.Date()
            
        // TypeScript specific types
        case "string | null", "String":
            if faker.Bool() {
                return nil
            }
            return generateMeaningfulString(fieldName)
            
        case "number | null", "Number | null":
            if faker.Bool() {
                return nil
            }
            return faker.Number(1, 1000)
		case "boolean | null", "Boolean | null":
			if faker.Bool() {
                return nil
            }
            return faker.Bool()
		case "Date | null":
			if faker.Bool() {
                return nil
            }
            return faker.Date()
            
        // Complex TypeScript types
        // case "Record<string, string>":
        //     count := rand.Intn(3) + 1
        //     record := make(map[string]string)
        //     for i := 0; i < count; i++ {
        //         record[generateMeaningfulString("key")] = generateMeaningfulString("value")
        //     }
        //     return record
            
        // case "Record<string, number>":
        //     count := rand.Intn(3) + 1
        //     record := make(map[string]int)
        //     for i := 0; i < count; i++ {
        //         record[generateMeaningfulString("key")] = faker.Number(1, 1000)
        //     }
        //     return record
            
        default:
            return nil
        }
    }

    // Handle arrays
    if isArray {
        count := rand.Intn(3) + 2
        arr := make([]interface{}, count)
        for i := range arr {
            arr[i] = generateValue()
        }
        return arr
    }
    return generateValue()
}

func contains(str string, patterns ...string) bool {
	for _, pattern := range patterns {
		if strings.Contains(str, pattern) {
			return true
		}
	}
	return false
}

func generateMeaningfulString(fieldName string) interface{} {
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
	case contains(fieldName, "gender"):
		return faker.Gender()

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

	// products

	case contains(fieldName, "productname"):
		return faker.Product().Name
	case contains(fieldName, "price"):
		return fmt.Sprintf("%v", faker.Product().Price)
	default:
		return faker.LoremIpsumSentence(2)
	}

}
