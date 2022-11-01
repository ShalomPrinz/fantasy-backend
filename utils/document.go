package utils

import (
	"log"
	"strings"

	"cloud.google.com/go/firestore"
)

func GetDocData(doc *firestore.DocumentSnapshot, attr []string) map[string]string {
	m := map[string]string{
		"id": doc.Ref.ID,
	}

	for _, value := range attr {
		m[strings.ToLower(value)] = GetDocString(doc, value)
	}

	return m
}

func GetDocString(doc *firestore.DocumentSnapshot, field string) string {
	value, exist := doc.Data()[field]
	if !exist {
		log.Fatalf("No value found in doc: %v", field)
	}

	stringValue, ok := value.(string)
	if !ok {
		log.Fatalf("Not a string in doc field %v: %v", field, value)
	}
	return stringValue
}
