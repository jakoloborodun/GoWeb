package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"time"
)

type BlogPost struct {
	Mongo    `inline`
	Title    string
	Text     string
	Status   bool
	Created  time.Time
	Category *Category
	Content  template.HTML
}

func (p *BlogPost) GetMongoCollectionName() string {
	return "posts"
}

func (p *BlogPost) Create(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(p.GetMongoCollectionName())
	_, err := coll.InsertOne(ctx, p)
	return err
}

func (p *BlogPost) Update(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(p.GetMongoCollectionName())
	_, err := coll.ReplaceOne(ctx, bson.M{"_id": p.ID}, p)
	return err
}

func (p *BlogPost) Delete(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(p.GetMongoCollectionName())
	_, err := coll.DeleteOne(ctx, bson.M{"_id": p.ID})
	return err
}

func GetPost(id primitive.ObjectID, ctx context.Context, db *mongo.Database) (*BlogPost, error) {
	p := BlogPost{}
	coll := db.Collection(p.GetMongoCollectionName())

	res := coll.FindOne(ctx, bson.M{"_id": id})
	if err := res.Decode(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

func GetAllPosts(ctx context.Context, db *mongo.Database) ([]BlogPost, error) {
	p := BlogPost{}
	coll := db.Collection(p.GetMongoCollectionName())

	cur, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var posts []BlogPost
	if err := cur.All(ctx, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func FindPost(ctx context.Context, db *mongo.Database, field string, value interface{}) ([]BlogPost, error) {
	p := BlogPost{}
	coll := db.Collection(p.GetMongoCollectionName())

	cur, err := coll.Find(ctx, bson.M{field: value})
	if err != nil {
		return nil, err
	}

	var posts []BlogPost
	if err := cur.All(ctx, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func NewBlogPost(title string, text string, status bool, category *Category, content template.HTML) *BlogPost {
	return &BlogPost{
		Mongo:    Mongo{ID: primitive.NewObjectID()},
		Title:    title,
		Text:     text,
		Created:  time.Now(),
		Status:   status,
		Category: category,
		Content:  content,
	}
}
