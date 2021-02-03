package repository

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// QuadrantMongoDBService defines the interface that device must satisfy.
type QuadrantMongoDBService interface {
	Create(ctx context.Context, s *Quadrant) (id string, err error)
	Get(ctx context.Context, qf *QuadrantFilter) (q *Quadrant, err error)
	Update(ctx context.Context, s *Quadrant) (q *Quadrant, err error)
	Delete(ctx context.Context, qf *QuadrantFilter) (isRemoved bool, err error)
}

// QuadrantService represents a mongoServie that contains the MongoDB client.
type QuadrantService mongoService

// QuadrantService validate if it satisfy the own interface, that means
// that all mongoService can be implement its own interface but it must be
// a mongoService type.
var _ QuadrantMongoDBService = &QuadrantService{}

// Quadrant represents a quadrant information that is in the maze.
type Quadrant struct {
	ID         string       `json:"id,omitempty" bson:"_id,omitempty"`
	Type       QuadrantType `json:"type,omitempty" bson:"type"`
	Spots      []Spot       `json:"spots,omitempty"`
	SpotIDs    []string     `json:"spot_ids,omitempty" bson:"spot_ids"`
	StartPoint *Coordinate  `json:"start_point,omitempty" bson:"start_point"`
	LimitPoint *Coordinate  `json:"limit_point,omitempty" bson:"limit_point"`
}

// Coordinate represents a specific point/location in the maze.
type Coordinate struct {
	X uint `json:"x"`
	Y uint `json:"y"`
}

// QuadrantType defines the four existing quadrants.
type QuadrantType string

const (
	// TopLeft represents a quadrant location.
	TopLeft = QuadrantType("TOP_LEFT")

	// TopRight represents a quadrant location.
	TopRight = QuadrantType("TOP_RIGHT")

	// BottomLeft represents a quadrant location.
	BottomLeft = QuadrantType("BOTTOM_LEFT")

	// BottomRight represents a quadrant location.
	BottomRight = QuadrantType("BOTTOM_RIGTH")
)

// QuadrantFilter represents the filter that can be used to create a mongo query.
type QuadrantFilter struct {
	ID   string       `json:"id,omitempty"`
	Type QuadrantType `json:"type,omitempty"`
}

func (qf *QuadrantFilter) toMongoFilter() (bson.M, error) {
	if qf == nil {
		return nil, errors.New("quadrant filter must not be nil")
	}

	if qf.ID == "" && qf.Type == "" {
		return nil, errors.New("at least one of QuadrantFilter.ID and QuadrantFilter.Type must be specified")
	}

	filter := bson.M{}

	if qf.ID != "" {
		id, err := primitive.ObjectIDFromHex(qf.ID)
		if err != nil {
			return nil, fmt.Errorf("wrong id%s", qf.ID)
		}

		filter["_id"] = id
	}

	if qf.Type != "" {
		filter["type"] = qf.Type
	}

	return filter, nil
}

// Create creates a new quadrant in a maze.
func (qs *QuadrantService) Create(ctx context.Context, q *Quadrant) (string, error) {
	dbName := DBName(ctx)

	if len(q.SpotIDs) == 0 {
		q.SpotIDs = []string{}
	}

	res, err := qs.db.Database(dbName).Collection("quadrants").InsertOne(ctx, q)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

// Get gets a specific quadrant by its type.
func (qs *QuadrantService) Get(ctx context.Context, qf *QuadrantFilter) (*Quadrant, error) {
	filter, err := qf.toMongoFilter()
	if err != nil {
		return nil, err
	}

	var (
		quadrant = &Quadrant{}
	)

	err = qs.db.Database(DBName(ctx)).Collection("quadrants").FindOne(ctx, filter).Decode(quadrant)
	if err != nil {
		return nil, err
	}

	if len(quadrant.SpotIDs) > 0 {
		spots, err := New(qs.db).Spot.List(ctx, &SpotFilter{
			SpotsIDs: quadrant.SpotIDs,
		})

		if err != nil {
			return nil, err
		}

		quadrant.Spots = spots
	}

	return quadrant, nil
}

// Update creates a new quadrant in a maze.
func (qs *QuadrantService) Update(ctx context.Context, uq *Quadrant) (*Quadrant, error) {
	if uq == nil {
		return nil, errors.New("quadrant parameter must be specified")
	}

	if uq.Type == "" && uq.ID == "" {
		return nil, errors.New("at least one of Quadrant.ID and Quadrant.Type must be specified")
	}

	filter := &QuadrantFilter{
		ID:   uq.ID,
		Type: uq.Type,
	}

	cq, err := qs.Get(ctx, filter)

	if err != nil {
		return nil, fmt.Errorf("this is the get error %s", err)
	}

	if len(uq.SpotIDs) != 0 {
		idsToRemove := CompareStringSlices(cq.SpotIDs, uq.SpotIDs)

		if len(idsToRemove) > 0 {
			spotFilter := SpotFilter{SpotsIDs: idsToRemove}
			f, err := spotFilter.toMongoFilter()

			if err != nil {
				return nil, err
			}

			_, err = qs.db.Database(DBName(ctx)).Collection("spots").DeleteMany(ctx, f)

			if err != nil {
				return nil, err
			}
		}

		cq.SpotIDs = uq.SpotIDs
	}

	if uq.StartPoint != nil {
		cq.StartPoint = uq.StartPoint
	}

	if uq.LimitPoint != nil {
		cq.LimitPoint = uq.LimitPoint
	}

	update := bson.M{
		"type":        cq.Type,
		"spot_ids":    cq.SpotIDs,
		"start_point": cq.StartPoint,
		"limit_point": cq.LimitPoint,
	}

	qf, err := filter.toMongoFilter()
	if err != nil {
		return nil, err
	}

	_, err = qs.db.Database(DBName(ctx)).Collection("quadrants").ReplaceOne(ctx, qf, update)
	if err != nil {
		return nil, err
	}

	cq, err = qs.Get(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("this is the end error %s", err)
	}

	return cq, nil
}

// Delete deletes a quadrant by quadrant type.
func (qs *QuadrantService) Delete(ctx context.Context, qf *QuadrantFilter) (bool, error) {
	filter, err := qf.toMongoFilter()
	if err != nil {
		return false, err
	}

	var (
		dbName = DBName(ctx)
	)

	cq, err := qs.Get(ctx, qf)
	if err != nil {
		return false, err
	}

	if len(cq.SpotIDs) > 0 {
		if _, err := New(qs.db).Spot.Delete(ctx, &SpotFilter{SpotsIDs: cq.SpotIDs}); err != nil {
			return false, err
		}
	}

	_, err = qs.db.Database(dbName).Collection("quadrants").DeleteOne(ctx, filter)
	if err != nil {
		return false, err
	}

	return true, nil
}

// CompareStringSlices compare two slices that which will return the difference values of
// the first values that don't exists on the second slice.
func CompareStringSlices(x, y []string) []string {
	fn := func(smallSlice, maxSlice []string) []string {
		for i := range smallSlice {
			for k := 0; k < len(maxSlice); k++ {
				if smallSlice[i] == maxSlice[k] {
					maxSlice[k] = maxSlice[len(maxSlice)-1]

					maxSlice = maxSlice[:len(maxSlice)-1]

					k = -1
				}
			}
		}

		return maxSlice
	}

	return fn(y, x)
}
