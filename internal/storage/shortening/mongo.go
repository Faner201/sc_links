package shortening

import (
	"context"
	"fmt"
	"time"

	"github.com/Faner201/sc_links/internal/db"
	errors "github.com/Faner201/sc_links/internal/error"
	"github.com/Faner201/sc_links/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mgo struct {
	db *mongo.Database
}

func NewMongoDB(client *mongo.Database) *mgo {
	return &mgo{db: client}
}

func (m *mgo) col() *mongo.Collection {
	return m.db.Collection("shortening")
}

func (m *mgo) Put(ctx context.Context, shortening model.Shortering) (*model.Shortering, error) {
	const operation = "shortening.mgo.Put"

	shortening.CreatedAt = time.Now().UTC()

	count, err := m.col().CountDocuments(ctx, bson.M{"_id": shortening.Identifier})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	if count > 0 {
		return nil, fmt.Errorf("%s: %w", operation, errors.ErrIdentifiExists)
	}

	_, err = m.col().InsertOne(ctx, mgoShoreningFromModel(shortening))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &shortening, nil
}

func (m *mgo) Get(ctx context.Context, shorteningID string) (*model.Shortering, error) {
	const operation = "shortening.mgo.Get"

	var shortening db.MgoShortening
	if err := m.col().FindOne(ctx, bson.M{"_id": shorteningID}).Decode(&shortening); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%s: %w", operation, errors.ErrNotFound)
		}

		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return modelShorteningFromMgo(shortening), nil
}

func (m *mgo) IncrementVisits(ctx context.Context, shorteningID string) error {
	const operation = "shortening.mgo.IncrementVisits"

	var (
		filter = bson.M{"_id": shorteningID}
		update = bson.M{
			"$inc": bson.M{"visits": 1},
			"$set": bson.M{"updated_at": time.Now().UTC()},
		}
	)

	_, err := m.col().UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	return nil
}

func mgoShoreningFromModel(shortening model.Shortering) db.MgoShortening {
	return db.MgoShortening{
		Identifier:  shortening.Identifier,
		CreatedBy:   shortening.CreatedBy,
		OriginalURL: shortening.OriginalURL,
		Visits:      shortening.Visits,
		CreatedAt:   shortening.CreatedAt,
		UpdatedAt:   shortening.UpdatedAt,
	}
}

func modelShorteningFromMgo(shortening db.MgoShortening) *model.Shortering {
	return &model.Shortering{
		Identifier:  shortening.Identifier,
		CreatedBy:   shortening.CreatedBy,
		OriginalURL: shortening.OriginalURL,
		Visits:      shortening.Visits,
		CreatedAt:   shortening.CreatedAt,
		UpdatedAt:   shortening.UpdatedAt,
	}
}
