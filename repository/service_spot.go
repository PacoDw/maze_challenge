package repository

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SpotMongoDBService defines the interface that device must satisfy.
type SpotMongoDBService interface {
	Create(ctx context.Context, s *Spot) (id string, err error)
	Update(ctx context.Context, su *Spot) (s *Spot, err error)
	Get(ctx context.Context, sf *SpotFilter) (s *Spot, err error)
	List(ctx context.Context, sf *SpotFilter) (spots []Spot, err error)
	Delete(ctx context.Context, sf *SpotFilter) (isRemoved bool, err error)
}

// SpotService represents a mongoServie that contains the MongoDB client.
type SpotService mongoService

// SpotService validate if it satisfy the own interface, that means
// that all mongoService can be implement its own interface but it must be
// a mongoService type.
var _ SpotMongoDBService = &SpotService{}

// Spot represents a spot information that is in the maze.
type Spot struct {
	ID         string      `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string      `json:"name,omitempty" bson:"name,omitempty" binding:"required"`
	GoldAmount string      `json:"gold_mount,omitempty" bson:"gold_mount,omitempty" binding:"required"`
	Coordinate *Coordinate `json:"Coordinate,omitempty" bson:"Coordinate,omitempty" binding:"required"`
	QuadrantID string      `json:"quadrant_id,omitempty" bson:"quadrant_id,omitempty" binding:"required"`
}

// SpotFilter represents the filter that can be used to create a mongo query.
type SpotFilter struct {
	ID         string   `json:"id,omitempty" bson:"id,omitempty"`
	QuadrantID string   `json:"quadrant_id,omitempty" bson:"quadrant_id,omitempty"`
	SpotsIDs   []string `json:"spot_ids,omitempty"`
}

func (qf *SpotFilter) toMongoFilter() (bson.M, error) {
	if qf == nil {
		return nil, errors.New("quadrant filter must not be nil")
	}

	if qf.ID == "" && qf.QuadrantID == "" && len(qf.SpotsIDs) == 0 {
		return nil, errors.New("at least one of SpotFilter.ID, SpotFilter.QuadrantID and SpotFilter.SpotIDs must be specified")
	}

	filter := bson.M{}

	if qf.ID != "" {
		id, err := primitive.ObjectIDFromHex(qf.ID)
		if err != nil {
			return nil, fmt.Errorf("wrong id%s", qf.ID)
		}

		filter["_id"] = id
	}

	if qf.QuadrantID != "" {
		filter["quadrant_id"] = qf.QuadrantID
	}

	if len(qf.SpotsIDs) > 0 {
		ids := make([]primitive.ObjectID, len(qf.SpotsIDs))

		for i := range qf.SpotsIDs {
			id, err := primitive.ObjectIDFromHex(qf.SpotsIDs[i])
			if err == nil {
				ids[i] = id
			}
		}

		filter["_id"] = bson.M{"$in": ids}
	}

	return filter, nil
}

// Create creates a new spot in a maze.
func (ss *SpotService) Create(ctx context.Context, s *Spot) (string, error) {
	if s.QuadrantID == "" {
		return "", errors.New("the Spot.QuadrantID attribute must be specified")
	}

	if s.Name == "" {
		return "", errors.New("the Spot.Name attribute must be specified")
	}

	if s.Coordinate == nil {
		return "", errors.New("the Spot.Coordinate attribute must be specified")
	}

	if _, err := New(ss.db).Quadrant.Get(ctx, &QuadrantFilter{ID: s.QuadrantID}); err != nil {
		return "", fmt.Errorf("can't reach quadrant: %s", err)
	}

	res, err := ss.db.Database(DBName(ctx)).Collection("spots").InsertOne(ctx, s)
	if err != nil {
		return "", fmt.Errorf("inserting a spot: %s", err)
	}

	q, err := New(ss.db).Quadrant.Get(ctx, &QuadrantFilter{ID: s.QuadrantID})
	if err != nil {
		return "", fmt.Errorf("getting a quadrant: %s", err)
	}

	q.SpotIDs = append(q.SpotIDs, res.InsertedID.(primitive.ObjectID).Hex())

	if _, err := New(ss.db).Quadrant.Update(ctx, q); err != nil {
		return "", fmt.Errorf("updating a quadrant: %s", err)
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

// Get gets a spot in a maze.
func (ss *SpotService) Get(ctx context.Context, sf *SpotFilter) (*Spot, error) {
	if sf.ID == "" {
		return nil, errors.New("the SpotFilter.ID attribute must be specified")
	}

	filter, err := sf.toMongoFilter()
	if err != nil {
		return nil, err
	}

	var (
		dbName = DBName(ctx)
		spot   = &Spot{}
	)

	err = ss.db.Database(dbName).Collection("spots").FindOne(ctx, filter).Decode(spot)
	if err != nil {
		return nil, err
	}

	return spot, nil
}

// Update creates a spot in a maze.
func (ss *SpotService) Update(ctx context.Context, su *Spot) (*Spot, error) {
	if su.ID == "" {
		return nil, errors.New("the Spot.ID attribute must be specified")
	}

	if su.QuadrantID == "" {
		return nil, errors.New("the Spot.QuadrantID attribute must be specified")
	}

	filter := &SpotFilter{ID: su.ID}

	spot, err := ss.Get(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("can't find the spot: %s", err)
	}

	if su.Name != "" {
		spot.Name = su.Name
	}

	if su.Coordinate != nil {
		spot.Coordinate = su.Coordinate
	}

	if su.GoldAmount != "" {
		spot.GoldAmount = su.GoldAmount
	}

	update := bson.M{
		"name":        spot.Name,
		"gold_mount":  spot.GoldAmount,
		"Coordinate":  spot.Coordinate,
		"quadrant_id": su.QuadrantID,
	}

	sf, err := filter.toMongoFilter()
	if err != nil {
		return nil, fmt.Errorf("parsing internal filter: %s", err)
	}

	_, err = ss.db.Database(DBName(ctx)).Collection("spots").ReplaceOne(ctx, sf, update)
	if err != nil {
		return nil, err
	}

	if spot.QuadrantID != su.QuadrantID {
		q, err := New(ss.db).Quadrant.Get(ctx, &QuadrantFilter{ID: spot.QuadrantID})
		if err != nil {
			return nil, err
		}

		newSpotIDs := make([]string, 0)

		for i := range q.SpotIDs {
			if spot.ID != q.SpotIDs[i] {
				newSpotIDs = append(newSpotIDs, spot.ID)
			}
		}

		if _, err := New(ss.db).Quadrant.Update(ctx, &Quadrant{ID: spot.ID, SpotIDs: newSpotIDs}); err != nil {
			return nil, err
		}

		q, err = New(ss.db).Quadrant.Get(ctx, &QuadrantFilter{ID: su.QuadrantID})
		if err != nil {
			return nil, err
		}

		q.SpotIDs = append(q.SpotIDs, spot.ID)

		if _, err := New(ss.db).Quadrant.Update(ctx, q); err != nil {
			return nil, err
		}
	}

	return spot, nil
}

// List gets all spots by quadrant id in a maze.
func (ss *SpotService) List(ctx context.Context, sf *SpotFilter) ([]Spot, error) {
	if sf.ID != "" {
		return nil, errors.New("the SpotFilter.ID attribute must not be specified")
	}

	if sf.QuadrantID == "" && len(sf.SpotsIDs) == 0 {
		return nil, errors.New("at least one of SpotFilter.QuadrantID and SpotFilter.SpotIDs attributes must be specified")
	}

	filter, err := sf.toMongoFilter()
	if err != nil {
		return nil, fmt.Errorf("decoding filter for Spot.List: %s", err)
	}

	cursor, err := ss.db.Database(DBName(ctx)).Collection("spots").Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("finding all spots: %s", err)
	}

	defer cursor.Close(ctx)

	spots := make([]Spot, 0)

	for cursor.Next(context.TODO()) {
		var s Spot
		err := cursor.Decode(&s)

		if err == nil {
			spots = append(spots, s)
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("can't decode spots: %s", err)
	}

	return spots, nil
}

// Delete deletes a spot.
func (ss *SpotService) Delete(ctx context.Context, sf *SpotFilter) (bool, error) {
	spotIDs := make([]string, 0)

	switch {
	case sf.ID != "":
		spotIDs = append(spotIDs, sf.ID)
	case sf.QuadrantID != "":
		spots, err := ss.List(ctx, sf)
		if err != nil {
			return false, err
		}

		for i := range spots {
			spotIDs = append(spotIDs, spots[i].ID)
		}

	case len(sf.SpotsIDs) > 0:
		spotIDs = sf.SpotsIDs
	}

	if len(spotIDs) > 0 {
		filter := SpotFilter{
			SpotsIDs: spotIDs,
		}

		f, err := filter.toMongoFilter()
		if err != nil {
			return false, err
		}

		_, err = ss.db.Database(DBName(ctx)).Collection("spots").DeleteMany(ctx, f)
		if err != nil {
			return false, err
		}

		_, err = ss.db.Database(DBName(ctx)).Collection("quadrants").UpdateMany(ctx,
			bson.M{},
			bson.M{"$pullAll": bson.M{
				"spot_ids": spotIDs,
			}},
		)

		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, errors.New("filter is wrong")
}
