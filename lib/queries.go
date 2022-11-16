package lib

import (
	"fantasy/database/utils"
	"log"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetAll[T any](ctx *gin.Context, collection string) []T {
	var result []T
	iter := Client.Collection(collection).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate over %s collection: %v", collection, err)
		}
		result = append(result, utils.GetDocData[T](doc))
	}
	return result
}

func GetSingle[T any](ctx *gin.Context, collection string, id string) T {
	doc, err := Client.Collection(collection).Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		log.Fatalf("No such item with id %s in collection %s", id, collection)
	}
	if err != nil {
		log.Fatalf("Unexpected error trying to reach id %s in collection %s", id, collection)
	}

	return utils.GetDocData[T](doc)
}

func InsertItem(ctx *gin.Context, collection string, item interface{}) {
	_, _, err := Client.Collection(collection).Add(ctx, item)
	if err != nil {
		log.Fatalf("Failed adding item to %s collection: %v", collection, err)
	}
}

func InsertItemCustomID(ctx *gin.Context, collection string, id string, item interface{}) {
	_, err := Client.Collection(collection).Doc(id).Set(ctx, item)
	if err != nil {
		log.Fatalf("Failed adding item to %s collection: %v", collection, err)
	}
}
