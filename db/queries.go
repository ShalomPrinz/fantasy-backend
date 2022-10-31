package db

import (
	"context"
	"fantasy/database/utils"
	"log"

	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetAll(collection string, attr []string) []interface{} {
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
		result = append(result, utils.GetDocData(doc, attr))
	}
	return result
}

func GetSingle(collection string, id string, attr []string) interface{} {
	ctx := context.Background()
	doc, err := Client.Collection(collection).Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		log.Fatalf("No such item with id %s in collection %s", id, collection)
	}
	if err != nil {
		log.Fatalf("Unexpected error trying to reach id %s in collection %s", id, collection)
	}

	return utils.GetDocData(doc, attr)
}

func InsertItem(collection string, item interface{}) {
	ctx := context.Background()
	_, _, err := Client.Collection(collection).Add(ctx, item)
	if err != nil {
		log.Fatalf("Failed adding item to %s collection: %v", collection, err)
	}
}

func InsertItemCustomID(collection string, id string, item interface{}) {
	ctx := context.Background()
	_, err := Client.Collection(collection).Doc(id).Set(ctx, item)
	if err != nil {
		log.Fatalf("Failed adding item to %s collection: %v", collection, err)
	}
}
