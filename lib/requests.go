package lib

import (
	"fantasy/database/utils"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

func IsExists(ctx *gin.Context, ref string) bool {
	snaps, err := Client.GetAll(ctx, []*firestore.DocumentRef{
		Client.Doc(ref),
	})
	if err != nil {
		if isStatusNotFound(err) {
			return false
		} else {
			log.Printf("Error getting document ref %v: %v", ref, err)
			return false
		}
	}
	if !snaps[0].Exists() {
		return false
	}
	return true
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

func GetAll[T any](ctx *gin.Context, colRef *firestore.CollectionRef) ([]T, AppError) {
	var result []T
	iter := colRef.Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Failed to iterate over %v collection: %v", colRef, err)
			return nil, GetDocumentError(err)
		}
		result = append(result, utils.GetDocData[T](doc))
	}
	return result, EmptyError
}

func GetSingleRef[T any](ctx *gin.Context, ref string) (T, AppError) {
	result, appError := GetByIds[T](ctx, []string{ref})
	return result[0], appError
}

func GetByIds[T any](ctx *gin.Context, ids []string) ([]T, AppError) {
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
	colRef := Client.Collection(collection)
	return InsertItemToCollection(ctx, colRef, item)
}

func InsertItemToCollection(ctx *gin.Context, colRef *firestore.CollectionRef, item any) (string, AppError) {
	docRef, _, err := colRef.Add(ctx, item)
	if err != nil {
		log.Printf("Failed adding item to %v collection: %v", colRef, err)
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

func RemoveItemFromCollection(ctx *gin.Context, documentPath string) AppError {
	docRef := Client.Doc(documentPath)
	if _, err := docRef.Delete(ctx); err != nil {
		log.Printf("Failed remove document %v from collection: %v", docRef, err)
		return RemoveItemError(err)
	}
	return EmptyError
}

func SubCollectionRef(collection string, docId string, subcollection string) *firestore.CollectionRef {
	return Client.Collection(collection).Doc(docId).Collection(subcollection)
}
