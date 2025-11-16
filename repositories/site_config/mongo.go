package site_config

import (
	"context"

	"github.com/imcrazytwkr/formdrain/models/site_config"
	"github.com/imcrazytwkr/formdrain/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type siteFormConfigRepository struct {
	collection *mongo.Collection
}

func NewMongoSiteConfigRepository(db *mongo.Database) repositories.SiteConfigRepository {
	return &siteFormConfigRepository{
		collection: db.Collection(siteConfigColectionName),
	}
}

func (r *siteFormConfigRepository) GetSiteConfigById(ctx context.Context, id string) (*site_config.SiteConfig, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil || objectId.IsZero() {
		return nil, nil
	}

	var config site_config.SiteConfig
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
