package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Category struct {
	Mongo       `inline`
	Title       string `bson:"title"`
	Description string `bson:"desc"`
}

func (cat *Category) GetMongoCollectionName() string {
	return "categories"
}

func (cat *Category) Create(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(cat.GetMongoCollectionName())
	_, err := coll.InsertOne(ctx, cat)
	if err != nil {
		return err
	}

	return nil
}

func (cat *Category) Update(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(cat.GetMongoCollectionName())
	_, err := coll.ReplaceOne(ctx, bson.M{"_id": cat.ID}, cat)
	return err
}

func (cat *Category) Delete(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(cat.GetMongoCollectionName())
	_, err := coll.DeleteOne(ctx, bson.M{"_id": cat.ID})
	return err
}

func GetCategory(id primitive.ObjectID, ctx context.Context, db *mongo.Database) (*Category, error) {
	cat := Category{}
	coll := db.Collection(cat.GetMongoCollectionName())

	res := coll.FindOne(ctx, bson.M{"_id": id})
	if err := res.Decode(&cat); err != nil {
		return nil, err
	}
	return &cat, nil
}

func GetAllCategories(ctx context.Context, db *mongo.Database) ([]Category, error) {
	cat := Category{}
	coll := db.Collection(cat.GetMongoCollectionName())

	cur, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var categories []Category
	if err := cur.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

func NewCategory(title, description string) *Category {
	return &Category{
		Mongo:       Mongo{ID: primitive.NewObjectID()},
		Title:       title,
		Description: description,
	}
}
