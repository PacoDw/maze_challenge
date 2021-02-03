package repository

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuadrant_CreateAndDelete(t *testing.T) {
	Test_EnvMongoDBConnectionString(t)

	var (
		ctx  = DBNameSet(context.Background(), os.Getenv("DB_NAME"))
		conn = NewByConnString(os.Getenv("MOGODB_CONN"))
	)

	id, err := conn.Quadrant.Create(ctx, &Quadrant{
		Type:       TopRight,
		StartPoint: &Coordinate{X: 0, Y: 0},
		LimitPoint: &Coordinate{X: 25, Y: 25},
	})

	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, id)
	assert.IsType(t, id, "")

	isRemoved, err := conn.Quadrant.Delete(ctx, &QuadrantFilter{ID: id})
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, isRemoved)
	assert.EqualValues(t, true, isRemoved)
}

func TestQuadrant_Get(t *testing.T) {
	Test_EnvMongoDBConnectionString(t)

	var (
		ctx  = DBNameSet(context.Background(), os.Getenv("DB_NAME"))
		conn = NewByConnString(os.Getenv("MOGODB_CONN"))
	)

	expected := &Quadrant{
		Type:       TopRight,
		StartPoint: &Coordinate{X: 0, Y: 0},
		LimitPoint: &Coordinate{X: 25, Y: 25},
	}

	id, err := conn.Quadrant.Create(ctx, expected)

	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, id)
	assert.IsType(t, id, "")

	got, err := conn.Quadrant.Get(ctx, &QuadrantFilter{ID: id})

	if err != nil {
		t.Fatal(err)
	}

	assert.IsType(t, &Quadrant{}, got)
	assert.NotNil(t, got)
	assert.EqualValues(t, expected, got)

	isRemoved, err := conn.Quadrant.Delete(ctx, &QuadrantFilter{ID: id})

	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, isRemoved)
	assert.EqualValues(t, true, isRemoved)
}

func TestQuadrant_Update(t *testing.T) {
	Test_EnvMongoDBConnectionString(t)

	var (
		ctx  = DBNameSet(context.Background(), os.Getenv("DB_NAME"))
		conn = NewByConnString(os.Getenv("MOGODB_CONN"))
	)

	id, err := conn.Quadrant.Create(ctx, &Quadrant{
		Type:       TopRight,
		StartPoint: &Coordinate{X: 0, Y: 0},
		LimitPoint: &Coordinate{X: 25, Y: 25},
	})

	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, id)
	assert.IsType(t, id, "")

	expected := &Quadrant{
		Type:       TopRight,
		StartPoint: &Coordinate{X: 0, Y: 4},
		LimitPoint: &Coordinate{X: 25, Y: 30},
		SpotIDs:    []string{},
	}

	got, err := conn.Quadrant.Update(ctx, expected)

	if err != nil {
		t.Fatal(err)
	}

	assert.IsType(t, &Quadrant{}, got)
	assert.NotNil(t, got)
	assert.EqualValues(t, expected, got)

	isRemoved, err := conn.Quadrant.Delete(ctx, &QuadrantFilter{ID: id})

	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, isRemoved)
	assert.EqualValues(t, true, isRemoved)
}
