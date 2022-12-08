package lib

import (
	"fantasy/database/utils"
	"log"

	"cloud.google.com/go/firestore"
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

func GetByIds[T any](ctx *gin.Context, collection string, ids []string) []T {
	var refs []*firestore.DocumentRef
	for _, item := range ids {
		refs = append(refs, Client.Doc(item))
	}
	snaps, err := Client.GetAll(ctx, refs)
	if err != nil {
		log.Fatalf("Unexpected error getting documents by ids. %v", err)
	}
	return utils.GetDocArrayData[T](snaps)
}

func InsertItem(ctx *gin.Context, collection string, item any) string {
	docRef, _, err := Client.Collection(collection).Add(ctx, item)
	if err != nil {
		log.Fatalf("Failed adding item to %s collection: %v", collection, err)
	}
	return docRef.ID
}

func InsertItemCustomID(ctx *gin.Context, collection string, id string, item any) {
	_, err := Client.Collection(collection).Doc(id).Set(ctx, item)
	if err != nil {
		log.Fatalf("Failed adding item to %s collection: %v", collection, err)
	}
}

func InsertItemIntoArray(ctx *gin.Context, collection string, doc string, path string, item any) {
	docs := Client.Collection(collection).Doc(doc)
	_, err := docs.Update(ctx, []firestore.Update{
		{Path: path, Value: firestore.ArrayUnion(item)},
	})
	if err != nil {
		log.Fatalf(
			"Failed adding item to array %v in doc %v in %v collection. %v",
			path, doc, collection, err)
	}
}

func RemoveItemFromArray(ctx *gin.Context, collection string, doc string, path string, item any) {
	docs := Client.Collection(collection).Doc(doc)
	_, err := docs.Update(ctx, []firestore.Update{
		{Path: path, Value: firestore.ArrayRemove(item)},
	})
	if err != nil {
		log.Fatalf(
			"Failed removing item from array %v in doc %v in %v collection. %v",
			path, doc, collection, err)
	}
}
