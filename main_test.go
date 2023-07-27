package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestMongoDBConnection(t *testing.T) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	assert.Nil(t, err, "Failed to connect to MongoDB")
	defer func() {
		assert.Nil(t, client.Disconnect(context.Background()), "Failed to disconnect from MongoDB")
	}()
	err = client.Ping(context.Background(), nil)
	assert.Nil(t, err, "Failed to ping MongoDB")
	collection = client.Database("local").Collection("todoapp")
	assert.NotNil(t, collection)
}

func TestFindAllTodoHandler(t *testing.T) {
	initMongoDB()
	filter := bson.M{}
	var todos []Todo
	cur, err := collection.Find(context.Background(), filter)
	assert.Nil(t, err, err)

	for cur.Next(context.Background()) {
		var todo Todo
		if err := cur.Decode(&todo); err != nil {
			assert.Nil(t, err, err)
		}
		todos = append(todos, todo)
	}
	if err := cur.Err(); err != nil {
		assert.Nil(t, err, err)
	}
}

func TestFindByIdTodoHandler(t *testing.T) {
	initMongoDB()
	id, err := primitive.ObjectIDFromHex("64c06ec5c1d62c139842f15f")
	assert.Nil(t, err)

	filter := bson.M{"_id": id}
	result := collection.FindOne(context.Background(), filter)
	assert.NotNil(t, result)

}

func TestCreateTodoHandler(t *testing.T) {
	initMongoDB()
	todo := Todo{
		Title:     "Test Todo",
		Completed: false,
	}
	_, err := collection.InsertOne(context.Background(), todo)
	assert.Nil(t, err)
}
