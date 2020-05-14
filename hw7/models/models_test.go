package models

import (
	"bytes"
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"testing"
)

var (
	ctx = context.TODO()
	db  *mongo.Database
	lg  *logrus.Logger
)

func areEqual(t *testing.T, a interface{}, b interface{}) bool {
	exp, ok := a.([]byte)
	if !ok {
		return reflect.DeepEqual(a, b)
	}

	act, ok := b.([]byte)
	if !ok {
		return false
	}

	return bytes.Equal(exp, act)
}

func setupDB() {
	uri := "mongodb://root:root@localhost:27017"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		lg.WithError(err).Fatal("can't connect to db")
	}

	db = client.Database("hw7")
}

func disconnectDB() {
	db.Client().Disconnect(ctx)
}

func TestBlogPost_Create(t *testing.T) {
	var blogPosts []BlogPost

	setupDB()

	post1 := NewBlogPost("Hello world", "Lorem ipsum dolor sit amet", true, &Category{}, "")
	post2 := NewBlogPost("Second", "Dura lex sed lex", false, &Category{}, "<p>Dura lex sed lex</p>")
	blogPosts = append(blogPosts, *post1)
	blogPosts = append(blogPosts, *post2)

	for _, post := range blogPosts {
		if err := post.Create(ctx, db); err != nil {
			t.Error(err)
		}
	}

	disconnectDB()
}

func TestGetAllPosts(t *testing.T) {
	var blogPosts []BlogPost

	setupDB()
	p := BlogPost{}
	if err := db.Collection(p.GetMongoCollectionName()).Drop(ctx); err != nil {
		t.Fatal(err)
	}

	post1 := NewBlogPost("Hello world", "Lorem ipsum dolor sit amet", true, &Category{}, "")
	post2 := NewBlogPost("Second", "Dura lex sed lex", false, &Category{}, "<p>Dura lex sed lex</p>")
	blogPosts = append(blogPosts, *post1)
	blogPosts = append(blogPosts, *post2)

	for _, post := range blogPosts {
		if err := post.Create(ctx, db); err != nil {
			t.Error(err)
		}
	}

	posts, err := GetAllPosts(ctx, db)
	if err != nil {
		t.Error(err)
	}
	disconnectDB()

	if !areEqual(t, posts, blogPosts) {
		t.Errorf("Received %v (type %v), expected %v (type %v)", posts, reflect.TypeOf(posts), blogPosts, reflect.TypeOf(blogPosts))
	}
}

func TestBlogPost_Update(t *testing.T) {
	setupDB()

	post := NewBlogPost("Test Update", "Lorem ipsum dolor sit amet", true, &Category{}, "")
	if err := post.Create(ctx, db); err != nil {
		t.Error(err)
	}

	post.Title = "Update successful"
	if err := post.Update(ctx, db); err != nil {
		t.Error(err)
	}

	afterUpdate, err := GetPost(post.ID, ctx, db)
	if err != nil {
		t.Error(err)
	}

	if afterUpdate.Title != "Update successful" {
		t.Error(
			"expected", "Update successful",
			"got", afterUpdate.Title,
		)
	}

	disconnectDB()
}
