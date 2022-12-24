package lib

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var Client *firestore.Client
var Auth *auth.Client

func InitClient() {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: os.Getenv("PROJID")}
	sa := option.WithCredentialsFile(os.Getenv("CREDPATH"))
	app, err := firebase.NewApp(ctx, conf, sa)
	if err != nil {
		log.Fatalln(err)
	}

	Client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	Auth, err = app.Auth(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}

func InitTestClient() {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "demo-test-fantasy"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	Client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	Auth, err = app.Auth(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}
