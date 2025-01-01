package generator

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/iancoleman/orderedmap"
	"github.com/structo/types"
)

var faker = gofakeit.New(0)

func GenerateMockObjects(fields []types.Field, count int) []orderedmap.OrderedMap {
	// Create slice to hold multiple objects
	mockObjects := make([]orderedmap.OrderedMap, count)

	// Generate 'count' number of objects
	for i := 0; i < count; i++ {
		mockDataMap := orderedmap.New()

		// sort by index
		sortedFields := make([]types.Field, len(fields))
		copy(sortedFields,fields)
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
	var fieldName = strings.Join(parts, "")
	switch field.Type {
	case "string":
		return generateMeaningfulString(fieldName)
	case "string[]":
		count := rand.Intn(3) + 2
		arr := make([]any, count)
		for i := range arr {
			arr[i] = (generateMeaningfulString(fieldName))
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
