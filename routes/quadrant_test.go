package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PacoDw/maze_challenge/repository"
	"github.com/stretchr/testify/assert"
)

func Test_QuadrantCreateDelete(t *testing.T) {
	Test_EnvMongoDBConnectionString(t)

	var quadrantID string

	// Create a quadrant
	t.Run("create a quadrant", func(t *testing.T) {
		blob := testExpectedBody(t, repository.Quadrant{
			Type: repository.BottomLeft,
			StartPoint: &repository.Coordinate{
				X: 0,
				Y: 0,
			},
			LimitPoint: &repository.Coordinate{
				X: 25,
				Y: 25,
			},
		})

		var expected repository.Quadrant
		err := json.Unmarshal(blob.Bytes(), &expected)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/create", blob)
		res := makeRequest(req, CreateQuadrant)

		assert.Equal(t, 200, res.Code)

		var got repository.Quadrant
		err = json.Unmarshal(res.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}

		expected.ID = got.ID

		assert.Equal(t, expected, got)

		quadrantID = got.ID
	})

	// Remove the quadrant
	t.Run("remove a quadrant", func(t *testing.T) {
		url := fmt.Sprintf("/delete/%s", quadrantID)

		req := httptest.NewRequest(http.MethodDelete, url, bytes.NewBuffer([]byte{}))
		res := makeRequest(req, DeleteQuadrant)

		assert.Equal(t, 200, res.Code)
	})
}

func Test_QuadrantGet(t *testing.T) {
	Test_EnvMongoDBConnectionString(t)

	var (
		expectedBody []byte
		quadrantID   string
	)

	// Create a quadrant
	t.Run("create a quadrant", func(t *testing.T) {
		blob := testExpectedBody(t, repository.Quadrant{
			Type: repository.BottomLeft,
			StartPoint: &repository.Coordinate{
				X: 0,
				Y: 0,
			},
			LimitPoint: &repository.Coordinate{
				X: 25,
				Y: 25,
			},
		})

		req := httptest.NewRequest(http.MethodPost, "/create", blob)
		res := makeRequest(req, CreateQuadrant)

		assert.Equal(t, 200, res.Code)

		expectedBody = res.Body.Bytes()
	})

	// read the quadrant
	t.Run("read a quadrant", func(t *testing.T) {
		var expected repository.Quadrant

		err := json.Unmarshal(expectedBody, &expected)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/read/%s", expected.ID)

		req := httptest.NewRequest(http.MethodGet, url, bytes.NewBuffer([]byte{}))
		res := makeRequest(req, GetQuadrant)

		assert.Equal(t, 200, res.Code)

		var got repository.Quadrant
		err = json.Unmarshal(res.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}

		quadrantID = expected.ID

		assert.Equal(t, expected, got)
	})

	// Remove the quadrant
	t.Run("remove a quadrant", func(t *testing.T) {
		url := fmt.Sprintf("/delete/%s", quadrantID)

		req := httptest.NewRequest(http.MethodDelete, url, bytes.NewBuffer([]byte{}))
		res := makeRequest(req, DeleteQuadrant)

		assert.Equal(t, 200, res.Code)
	})
}

func Test_QuadrantUpdate(t *testing.T) {
	Test_EnvMongoDBConnectionString(t)

	var quadrantID string

	// Create a quadrant
	t.Run("create a quadrant", func(t *testing.T) {
		blob := testExpectedBody(t, repository.Quadrant{
			Type: repository.BottomLeft,
			StartPoint: &repository.Coordinate{
				X: 0,
				Y: 0,
			},
			LimitPoint: &repository.Coordinate{
				X: 25,
				Y: 25,
			},
		})

		req := httptest.NewRequest(http.MethodPost, "/create", blob)
		res := makeRequest(req, CreateQuadrant)

		assert.Equal(t, 200, res.Code)

		var got repository.Quadrant
		err := json.Unmarshal(res.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}

		quadrantID = got.ID
	})

	// update the quadrant
	t.Run("update a quadrant", func(t *testing.T) {
		blob := testExpectedBody(t, repository.Quadrant{
			ID:   quadrantID,
			Type: repository.BottomLeft,
			StartPoint: &repository.Coordinate{
				X: 3,
				Y: 2,
			},
			LimitPoint: &repository.Coordinate{
				X: 25,
				Y: 25,
			},
		})

		var expected repository.Quadrant
		err := json.Unmarshal(blob.Bytes(), &expected)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodPatch, "/update", blob)
		res := makeRequest(req, UpdateQuadrant)

		assert.Equal(t, 200, res.Code)

		var got repository.Quadrant
		err = json.Unmarshal(res.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expected, got)
	})

	// Remove the quadrant
	t.Run("remove a quadrant", func(t *testing.T) {
		url := fmt.Sprintf("/delete/%s", quadrantID)

		req := httptest.NewRequest(http.MethodDelete, url, bytes.NewBuffer([]byte{}))
		res := makeRequest(req, DeleteQuadrant)

		assert.Equal(t, 200, res.Code)
	})
}
