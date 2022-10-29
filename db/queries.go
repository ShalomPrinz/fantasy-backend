package db

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type fn func(*firestore.DocumentSnapshot) interface{}

func GetAll(collection string, callback fn) []interface{} {
	var result []interface{}
	ctx := context.Background()
	iter := Client.Collection(collection).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate over %s collection: %v", collection, err)
		}
		result = append(result, callback(doc))
	}
	return result
}

func InsertItem(collection string, item interface{}) {
	ctx := context.Background()
	_, _, err := Client.Collection(collection).Add(ctx, item)
	if err != nil {
		log.Fatalf("Failed adding item to %s collection: %v", collection, err)
	}
}
