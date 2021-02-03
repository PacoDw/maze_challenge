package repository

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func Test_EnvMongoDBConnectionString(t *testing.T) {
	// check connection string is not empty
	if os.Getenv("MOGODB_CONN") == "" {
		t.Logf("we can't read env variable trying find it in .env file...")

		if err := godotenv.Load("../.env"); err != nil {
			t.Fatalf("error loading .env file %+v", err)
		}
	}
}

func Test_NewMongoDBConnection_withRightCredentials(t *testing.T) {
	Test_EnvMongoDBConnectionString(t)

	client := NewMongoDBConn(os.Getenv("MOGODB_CONN"))

	defer func() {
		err := client.Disconnect(context.Background())
		assert.NoError(t, err)
	}()

	// check that the client returned by the function is not empty or nil
	assert.NotEmpty(t, client)

	// check we are returning the right value
	assert.IsType(t, client, &mongo.Client{})
}
