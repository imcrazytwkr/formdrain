package form_response

import (
	"context"
	"fmt"

	"github.com/imcrazytwkr/formdrain/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoFormResponseRepository struct {
	collection *mongo.Collection
}

func NewMongoFormResponseRepository(db *mongo.Database) repositories.FormResponseRepository {
	return &mongoFormResponseRepository{db.Collection(formResponseCollectionName)}
}

func (r *mongoFormResponseRepository) SaveFormResponse(ctx context.Context, form map[string]any) (string, error) {
	result, err := r.collection.InsertOne(ctx, form)
	if err != nil {
		return "", err
	}

	rawId := result.InsertedID
	objectID, ok := rawId.(primitive.ObjectID)
	if ok {
		return objectID.String(), nil
	}

	stringId, ok := rawId.(string)
	if ok {
		return stringId, nil
	}

	return "", fmt.Errorf("mongo returned unexpected _id key: %q", rawId)
}
