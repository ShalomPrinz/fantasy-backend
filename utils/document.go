package utils

import (
	"log"

	"cloud.google.com/go/firestore"
	"github.com/mitchellh/mapstructure"
)

func GetDocData[T any](doc *firestore.DocumentSnapshot) T {
	entity := map[string]any{
		"ID": doc.Ref.ID,
	}

	if err := doc.DataTo(&entity); err != nil {
		log.Fatalf("Couldn't copy data from doc into the given struct. %v", err)
	}

	var result T
	err := mapstructure.Decode(entity, &result)
	if err != nil {
		log.Fatalf("Couldn't convert data given struct. %v", err)
	}
	return result
}
