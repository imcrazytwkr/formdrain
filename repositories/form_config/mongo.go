package form_config

import (
	"context"

	"github.com/imcrazytwkr/formdrain/models/form_config"
	"github.com/imcrazytwkr/formdrain/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoFormConfigRepository struct {
	collection *mongo.Collection
}

func NewMongoFormConfigRepository(db *mongo.Database) repositories.FormConfigRepository {
	return &mongoFormConfigRepository{
		collection: db.Collection(formConfigColectionName),
	}
}

func (r *mongoFormConfigRepository) GetFormConfigById(ctx context.Context, id string) (*form_config.FormConfig, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil || objectId.IsZero() {
		return nil, nil
	}

	var config form_config.FormConfig
	err = r.collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&config)
	switch err {
	case mongo.ErrNoDocuments:
		return nil, nil
	case nil:
		return &config, nil
	default:
		return nil, err
	}
}
