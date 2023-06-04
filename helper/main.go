package helper

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

func CreateFireStoreClient() (*firestore.Client, context.Context) {
	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsFile("./credentials.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return client, ctx
}

func GetSnapShotData(client *firestore.Client, ctx context.Context, path string) (map[string]interface{}, error) {
	dsnap, err := client.Doc(path).Get(ctx)
	return dsnap.Data(), err
}

func CreateNewDocument(client *firestore.Client, ctx context.Context, path, content string) (*firestore.WriteResult, error) {
	ref := client.Collection(path).NewDoc()

	res, err := ref.Set(ctx, gin.H{
		"author":       "Testing Script",
		"content":      content,
		"firstCreated": firestore.ServerTimestamp,
		"lastUpdated":  firestore.ServerTimestamp,
	})

	return res, err
}
