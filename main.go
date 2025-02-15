package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	secretsmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

func main() {
	projectNumber := os.Getenv("GOOGLE_CLOUD_PROJECT_NUMBER")
	secretName := os.Getenv("SECRET_NAME")

	ctx := context.Background()
	client, err := secretsmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}
	defer client.Close()

	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: "projects/" + projectNumber + "/secrets/" + secretName + "/versions/latest",
	}

	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		log.Fatalf("failed to access secret version: %v", err)
	}

	secret := string(result.Payload.Data)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Your secret is: %s", secret)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
