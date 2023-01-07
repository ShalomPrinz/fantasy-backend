package lib

import (
	"fantasy/database/utils"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

type Query struct {
	Field string
	Limit int
	Term  string
}

func QueryTermInField[T any](ctx *gin.Context, collection string, q Query) []T {
	result := make([]T, 0)
	iter := Client.Collection(collection).
		Where(q.Field, ">=", q.Term).
		Where(q.Field, "<=", q.Term+"\uf8ff").
		OrderBy(q.Field, firestore.Asc).
		Limit(q.Limit).
		Documents(ctx)

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
