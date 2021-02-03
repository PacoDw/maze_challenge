package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/PacoDw/maze_challenge/repository"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func makeRequest(req *http.Request, handler gin.HandlerFunc) *httptest.ResponseRecorder {
	var (
		w    = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
	)

	switch req.Method {
	case http.MethodDelete, http.MethodGet:
		splitURL := strings.Split(req.RequestURI, "/")
		param := splitURL[len(splitURL)-1]

		c.Params = append(c.Params, gin.Param{Key: "id", Value: param})
	}

	c.Request = req

	conn := repository.NewByConnString(os.Getenv("MOGODB_CONN"))

	c.Set("mongoRepoConn", conn)

	c.Set(string(repository.ContextDBName), os.Getenv("DB_NAME"))

	handler(c)

	return w
}

func Test_EnvMongoDBConnectionString(t *testing.T) {
	// check connection string is not empty
	if os.Getenv("MOGODB_CONN") == "" {
		t.Logf("we can't read env variable trying find it in .env file...")

		if err := godotenv.Load("../.env"); err != nil {
			t.Fatalf("error loading .env file %+v", err)
		}
	}
}

func testExpectedBody(t *testing.T, expectedBody interface{}) *bytes.Buffer {
	t.Helper()

	b, err := json.Marshal(expectedBody)
	if err != nil {
		t.Fatal(err)
	}

	return bytes.NewBuffer(b)
}

func Test_SpotCreateDelete(t *testing.T) {
	Test_EnvMongoDBConnectionString(t)

	var spotID string

	// Create a spot
	t.Run("create a spot", func(t *testing.T) {
		blob := testExpectedBody(t, repository.Spot{
			Name:       "exit",
			GoldAmount: "200",
			Coordinate: &repository.Coordinate{
				X: 10,
				Y: 3,
			},
			QuadrantID: "601989b15f19695a0a281bef",
		})

		var expected repository.Spot
		err := json.Unmarshal(blob.Bytes(), &expected)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/create", blob)
		res := makeRequest(req, CreateSpot)

		assert.Equal(t, 200, res.Code)

		var got repository.Spot
		err = json.Unmarshal(res.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}

		expected.ID = got.ID

		assert.Equal(t, expected, got)

		spotID = got.ID
	})

	// Remove the spot
	t.Run("remove a spot", func(t *testing.T) {
		url := fmt.Sprintf("/delete/%s", spotID)

		req := httptest.NewRequest(http.MethodDelete, url, bytes.NewBuffer([]byte{}))
		res := makeRequest(req, DeleteSpot)

		assert.Equal(t, 200, res.Code)
	})
}

func Test_SpotGet(t *testing.T) {
	Test_EnvMongoDBConnectionString(t)

	var (
		expectedBody []byte
		spotID       string
	)

	// Create a spot
	t.Run("create a spot", func(t *testing.T) {
		blob := testExpectedBody(t, repository.Spot{
			Name:       "exit",
			GoldAmount: "200",
			Coordinate: &repository.Coordinate{
				X: 10,
				Y: 3,
			},
			QuadrantID: "601989b15f19695a0a281bef",
		})

		req := httptest.NewRequest(http.MethodPost, "/create", blob)
		res := makeRequest(req, CreateSpot)

		assert.Equal(t, 200, res.Code)

		expectedBody = res.Body.Bytes()
	})

	// read the spot
	t.Run("read a spot", func(t *testing.T) {
		var expected repository.Spot

		err := json.Unmarshal(expectedBody, &expected)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/read/%s", expected.ID)

		req := httptest.NewRequest(http.MethodGet, url, bytes.NewBuffer([]byte{}))
		res := makeRequest(req, GetSpot)

		assert.Equal(t, 200, res.Code)

		var got repository.Spot
		err = json.Unmarshal(res.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}

		spotID = expected.ID

		assert.Equal(t, expected, got)
	})

	// Remove the spot
	t.Run("remove a spot", func(t *testing.T) {
		url := fmt.Sprintf("/delete/%s", spotID)

		req := httptest.NewRequest(http.MethodDelete, url, bytes.NewBuffer([]byte{}))
		res := makeRequest(req, DeleteSpot)

		assert.Equal(t, 200, res.Code)
	})
}

func Test_SpotUpdate(t *testing.T) {
	Test_EnvMongoDBConnectionString(t)

	var spotID string

	// Create a spot
	t.Run("create a spot", func(t *testing.T) {
		blob := testExpectedBody(t, repository.Spot{
			Name:       "exit",
			GoldAmount: "200",
			Coordinate: &repository.Coordinate{
				X: 10,
				Y: 3,
			},
			QuadrantID: "601989b15f19695a0a281bef",
		})

		req := httptest.NewRequest(http.MethodPost, "/create", blob)
		res := makeRequest(req, CreateSpot)

		assert.Equal(t, 200, res.Code)

		var got repository.Spot
		err := json.Unmarshal(res.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}

		spotID = got.ID
	})

	// update the spot
	t.Run("update a spot", func(t *testing.T) {
		blob := testExpectedBody(t, repository.Spot{
			ID:         spotID,
			Name:       "entrace",
			GoldAmount: "500",
			Coordinate: &repository.Coordinate{
				X: 21,
				Y: 4,
			},
			QuadrantID: "601989b15f19695a0a281bef",
		})

		var expected repository.Spot
		err := json.Unmarshal(blob.Bytes(), &expected)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodPatch, "/update", blob)
		res := makeRequest(req, UpdateSpot)

		assert.Equal(t, 200, res.Code)

		var got repository.Spot
		err = json.Unmarshal(res.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expected, got)
	})

	// Remove the spot
	t.Run("remove a spot", func(t *testing.T) {
		url := fmt.Sprintf("/delete/%s", spotID)

		req := httptest.NewRequest(http.MethodDelete, url, bytes.NewBuffer([]byte{}))
		res := makeRequest(req, DeleteSpot)

		assert.Equal(t, 200, res.Code)
	})
}
