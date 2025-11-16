package common

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ObjectId string

func (id ObjectId) String() string {
	return string(id)
}

func (id ObjectId) MarshalBSONValue() (bsontype.Type, []byte, error) {
	p, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return bsontype.Null, nil, err
	}

	return bson.MarshalValue(p)
}

func (id *ObjectId) UnmarshalBSONValue(t bsontype.Type, b []byte) error {
	var raw primitive.ObjectID
	err := bson.Unmarshal(b, raw)
	if err != nil {
		return err
	}

	*id = ObjectId(raw.Hex())
	return nil
}
