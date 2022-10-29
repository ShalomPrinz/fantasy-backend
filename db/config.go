package db

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var Client *firestore.Client

func InitClient() {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: os.Getenv("PROJID")}
	sa := option.WithCredentialsFile(os.Getenv("CREDPATH"))
	app, err := firebase.NewApp(ctx, conf, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	Client = client
}
