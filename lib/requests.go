package lib

import (
	"fantasy/database/utils"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

func IsExists(ctx *gin.Context, collection string, id string) bool {
	_, err := Client.Collection(collection).Doc(id).Get(ctx)
	if err != nil {
		if isStatusNotFound(err) {
			return false
		} else {
			log.Printf("Error checking if id %s exists in collection %s. The result might be corrupted", id, collection)
			return false
		}
	}
	return true
}

func GetAll[T any](ctx *gin.Context, collection string) ([]T, AppError) {
	var result []T
	iter := Client.Collection(collection).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Failed to iterate over %s collection: %v", collection, err)
			return nil, GetDocumentError(err)
		}
		result = append(result, utils.GetDocData[T](doc))
	}
	return result, EmptyError
}

func GetSingle[T any](ctx *gin.Context, collection string, id string) (T, AppError) {
	doc, err := Client.Collection(collection).Doc(id).Get(ctx)
	if err != nil {
		log.Printf("Error trying to reach id %s in collection %s", id, collection)
		var zeroValue T
		return zeroValue, GetDocumentError(err)
	}
	return utils.GetDocData[T](doc), EmptyError
}

func GetByIds[T any](ctx *gin.Context, collection string, ids []string) ([]T, AppError) {
	var refs []*firestore.DocumentRef
	for _, item := range ids {
		refs = append(refs, Client.Doc(item))
	}
	snaps, err := Client.GetAll(ctx, refs)
	if err != nil {
		log.Printf("Unexpected error getting documents by ids. %v", err)
		return nil, GetDocumentError(err)
	}
	return utils.GetDocArrayData[T](snaps), EmptyError
}

func InsertItem(ctx *gin.Context, collection string, item any) (string, AppError) {
	docRef, _, err := Client.Collection(collection).Add(ctx, item)
	if err != nil {
		log.Printf("Failed adding item to %s collection: %v", collection, err)
		return "", InsertItemError(err)
	}
	return docRef.ID, EmptyError
}

func InsertItemCustomID(ctx *gin.Context, collection string, id string, item any) AppError {
	_, err := Client.Collection(collection).Doc(id).Set(ctx, item)
	if err != nil {
		log.Printf("Failed adding item to %s collection: %v", collection, err)
		return InsertItemError(err)
	}
	return EmptyError
}

func InsertItemIntoArray(ctx *gin.Context, collection string, doc string, path string, item any) AppError {
	docs := Client.Collection(collection).Doc(doc)
	_, err := docs.Update(ctx, []firestore.Update{
		{Path: path, Value: firestore.ArrayUnion(item)},
	})
	if err != nil {
		log.Printf(
			"Failed adding item to array %v in doc %v in %v collection. %v",
			path, doc, collection, err)
		return InsertItemError(err)
	}
	return EmptyError
}

func RemoveItemFromArray(ctx *gin.Context, collection string, doc string, path string, item any) AppError {
	docs := Client.Collection(collection).Doc(doc)
	_, err := docs.Update(ctx, []firestore.Update{
		{Path: path, Value: firestore.ArrayRemove(item)},
	})
	if err != nil {
		log.Printf(
			"Failed removing item from array %v in doc %v in %v collection. %v",
			path, doc, collection, err)
		return RemoveItemError(err)
	}
	return EmptyError
}
