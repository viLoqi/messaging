package core

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

func CreateNewDocument(client *firestore.Client, ctx context.Context, path, author, content, authorPhotoURL string) (*firestore.DocumentRef, error) {
	ref := client.Collection(path).NewDoc()

	_, err := ref.Set(ctx, gin.H{
		"author":         author,
		"authorPhotoURL": authorPhotoURL,
		"content":        content,
		"firstCreated":   firestore.ServerTimestamp,
		"lastUpdated":    firestore.ServerTimestamp,
	})

	return ref, err
}

func DeleteDocument(client *firestore.Client, ctx context.Context, path string) (*firestore.WriteResult, error) {
	res, err := client.Doc(path).Delete(ctx)
	return res, err
}

func UpdateDocument(client *firestore.Client, ctx context.Context, path, content string) (*firestore.WriteResult, error) {
	res, err := client.Doc(path).Update(ctx, []firestore.Update{
		{
			Path:  "content",
			Value: content,
		},
		{
			Path:  "lastUpdated",
			Value: firestore.ServerTimestamp,
		},
	})
	return res, err
}
