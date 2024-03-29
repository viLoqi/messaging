package core

import (
	"context"
	"log"
	"loqi/messaging/structs"

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

func CreateNewDocument(client *firestore.Client, ctx context.Context, requestBody structs.WriteMessageRequestBody) (*firestore.DocumentRef, error) {
	ref := client.Collection(requestBody.CollectionPath).NewDoc()

	_, err := ref.Set(ctx, gin.H{
		"author":         requestBody.Author,
		"authorPhotoURL": requestBody.AuthorPhotoURL,
		"content":        requestBody.Content,
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
