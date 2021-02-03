package repository

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpot_CreateListDelete(t *testing.T) {
	Test_EnvMongoDBConnectionString(t)

	var (
		ctx  = DBNameSet(context.Background(), os.Getenv("DB_NAME"))
		conn = NewByConnString(os.Getenv("MOGODB_CONN"))
	)

	qID, err := conn.Quadrant.Create(ctx, &Quadrant{
		Type:       TopRight,
		StartPoint: &Coordinate{X: 0, Y: 0},
		LimitPoint: &Coordinate{X: 25, Y: 25},
	})

	if err != nil {
		t.Fatal(err)
	}

	sID, err := conn.Spot.Create(ctx, &Spot{
		Name:       "exit",
		GoldAmount: "4000",
		Coordinate: &Coordinate{
			X: 9,
			Y: 0,
		},
		QuadrantID: qID,
	})

	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, sID)
	assert.IsType(t, sID, "")

	sID, err = conn.Spot.Create(ctx, &Spot{
		Name:       "entrace",
		GoldAmount: "9000",
		Coordinate: &Coordinate{
			X: 0,
			Y: 10,
		},
		QuadrantID: qID,
	})

	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, sID)
	assert.IsType(t, sID, "")

	spots, err := conn.Spot.List(ctx, &SpotFilter{QuadrantID: qID})
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, spots)
	assert.Len(t, spots, 2)

	isRemoved, err := conn.Spot.Delete(ctx, &SpotFilter{QuadrantID: qID})
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, isRemoved)
	assert.EqualValues(t, true, isRemoved)

	if _, err := conn.Quadrant.Delete(ctx, &QuadrantFilter{ID: qID}); err != nil {
		t.Fatal(err)
	}
}

func TestSpot_GetDelete(t *testing.T) {
	Test_EnvMongoDBConnectionString(t)

	var (
		ctx  = DBNameSet(context.Background(), os.Getenv("DB_NAME"))
		conn = NewByConnString(os.Getenv("MOGODB_CONN"))
	)

	qID, err := conn.Quadrant.Create(ctx, &Quadrant{
		Type:       TopRight,
		StartPoint: &Coordinate{X: 0, Y: 0},
		LimitPoint: &Coordinate{X: 25, Y: 25},
	})

	if err != nil {
		t.Fatal(err)
	}

	expectedValue := &Spot{
		Name:       "exit",
		GoldAmount: "4000",
		Coordinate: &Coordinate{
			X: 9,
			Y: 0,
		},
		QuadrantID: qID,
	}

	sID, err := conn.Spot.Create(ctx, expectedValue)

	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, sID)
	assert.IsType(t, sID, "")

	got, err := conn.Spot.Get(ctx, &SpotFilter{ID: sID})
	if err != nil {
		t.Fatal(err)
	}

	expectedValue.ID = got.ID

	assert.NotEmpty(t, got)
	assert.IsType(t, &Spot{}, got)
	assert.EqualValues(t, expectedValue, got)

	isRemoved, err := conn.Spot.Delete(ctx, &SpotFilter{ID: got.ID})
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, isRemoved)
	assert.EqualValues(t, true, isRemoved)

	if _, err := conn.Quadrant.Delete(ctx, &QuadrantFilter{ID: qID}); err != nil {
		t.Fatal(err)
	}
}

func TestSpot_filterSpotList(t *testing.T) {
	Test_EnvMongoDBConnectionString(t)

	var (
		ctx  = DBNameSet(context.Background(), os.Getenv("DB_NAME"))
		conn = NewByConnString(os.Getenv("MOGODB_CONN"))
	)

	qID, err := conn.Quadrant.Create(ctx, &Quadrant{
		Type:       TopRight,
		StartPoint: &Coordinate{X: 0, Y: 0},
		LimitPoint: &Coordinate{X: 25, Y: 25},
	})

	if err != nil {
		t.Fatal(err)
	}

	sID1, err := conn.Spot.Create(ctx, &Spot{
		Name:       "exit",
		GoldAmount: "4000",
		Coordinate: &Coordinate{
			X: 9,
			Y: 0,
		},
		QuadrantID: qID,
	})

	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, sID1)
	assert.IsType(t, sID1, "")

	sID2, err := conn.Spot.Create(ctx, &Spot{
		Name:       "entrace",
		GoldAmount: "9000",
		Coordinate: &Coordinate{
			X: 0,
			Y: 10,
		},
		QuadrantID: qID,
	})

	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, sID2)
	assert.IsType(t, sID2, "")

	spots, err := conn.Spot.List(ctx, &SpotFilter{SpotsIDs: []string{sID1, sID2}})
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, spots)
	assert.Len(t, spots, 2)

	isRemoved, err := conn.Spot.Delete(ctx, &SpotFilter{QuadrantID: qID})
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, isRemoved)
	assert.EqualValues(t, true, isRemoved)

	if _, err := conn.Quadrant.Delete(ctx, &QuadrantFilter{ID: qID}); err != nil {
		t.Fatal(err)
	}
}
