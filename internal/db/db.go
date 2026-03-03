package db

import (
	"errors"
	"fmt"
	"os"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/appwrite/sdk-for-go/client"
	"github.com/appwrite/sdk-for-go/databases"
	"github.com/appwrite/sdk-for-go/users"
)

var Client client.Client
var Databases *databases.Databases
var Users *users.Users

// Collection config loaded from env
var DatabaseID string
var ProfilesCollectionID string
var ContactsCollectionID string

func Init() error {
	endpoint := os.Getenv("APPWRITE_ENDPOINT")
	projectID := os.Getenv("APPWRITE_PROJECT_ID")
	apiKey := os.Getenv("APPWRITE_API_KEY")

	if endpoint == "" || projectID == "" || apiKey == "" {
		return errors.New("APPWRITE_ENDPOINT, APPWRITE_PROJECT_ID, and APPWRITE_API_KEY must be set")
	}

	DatabaseID = os.Getenv("APPWRITE_DATABASE_ID")
	ProfilesCollectionID = os.Getenv("APPWRITE_PROFILES_COLLECTION_ID")
	ContactsCollectionID = os.Getenv("APPWRITE_CONTACTS_COLLECTION_ID")

	if DatabaseID == "" || ProfilesCollectionID == "" || ContactsCollectionID == "" {
		return errors.New("APPWRITE_DATABASE_ID, APPWRITE_PROFILES_COLLECTION_ID, and APPWRITE_CONTACTS_COLLECTION_ID must be set")
	}

	Client = appwrite.NewClient(
		appwrite.WithEndpoint(endpoint),
		appwrite.WithProject(projectID),
		appwrite.WithKey(apiKey),
	)

	Databases = appwrite.NewDatabases(Client)
	Users = appwrite.NewUsers(Client)

	fmt.Println("\nAppwrite client initialized successfully")
	return nil
}
