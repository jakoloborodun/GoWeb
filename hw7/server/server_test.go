package server

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var (
	ctx = context.TODO()
	db  *mongo.Database
	lg  *logrus.Logger
)

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

func TestNew(t *testing.T) {
	setupDB()
	srv := New(lg, "./www", db)

	if reflect.TypeOf(srv) != reflect.TypeOf(&Server{}) {
		t.Errorf("wrong type of server: got %v, expected %v", reflect.TypeOf(srv), reflect.TypeOf(&Server{}))
	}
	disconnectDB()
}

func TestServer_getTemplateHandler(t *testing.T) {
	setupDB()
	req, err := http.NewRequest("GET", "/blog", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	srv := New(lg, "../www", db)
	srv.ParseTemplates()
	handler := http.HandlerFunc(srv.getTemplateHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v expected %v",
			status, http.StatusOK)
	}

	disconnectDB()
}
