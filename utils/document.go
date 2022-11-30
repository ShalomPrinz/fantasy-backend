package utils

import (
	"log"

	"cloud.google.com/go/firestore"
)

func GetDocData[T any](doc *firestore.DocumentSnapshot) T {
	if !doc.Exists() {
		log.Fatalf("Given document %v doesn't exist, couldn't return its data", doc.Ref.Path)
	}

	entity := map[string]any{
		"ID": doc.Ref.ID,
	}

	if err := doc.DataTo(&entity); err != nil {
		log.Fatalf("Couldn't copy data from doc into the given struct. %v", err)
	}

	return MapToStruct[T](entity)
}

func GetDocArrayData[T any](docArray []*firestore.DocumentSnapshot) []T {
	var result []T
	for _, item := range docArray {
		result = append(result, GetDocData[T](item))
	}
	return result
}
