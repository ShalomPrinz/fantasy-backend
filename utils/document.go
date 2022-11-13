package utils

import (
	"log"
	"strings"

	"cloud.google.com/go/firestore"
)

func GetDocData(doc *firestore.DocumentSnapshot, attr []string) map[string]any {
	m := map[string]any{
		"id": doc.Ref.ID,
	}

	for _, value := range attr {
		m[strings.ToLower(value)] = GetDocString(doc, value)
	}

	return m
}

func GetDocString(doc *firestore.DocumentSnapshot, field string) any {
	value, exist := doc.Data()[field]
	if !exist {
		log.Fatalf("No value found in doc: %v", field)
	}
	return value
}
