package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

type Todo struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title,omitempty" bson:"title,omitempty"`
	Completed bool               `json:"completed,omitempty" bson:"completed,omitempty"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}

var client *mongo.Client
var collection *mongo.Collection

func createTodoHandler(c *gin.Context) {
	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	todo.CreatedAt = time.Now()
	result, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create a new Todo"})
		return
	}

	todo.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, todo)
}

func listTodoHandler(c *gin.Context) {
	var todos []Todo
	filter := bson.M{}
	result, err := collection.Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get a List"})
		return
	}

	for result.Next(context.Background()) {
		var todo Todo
		if err := result.Decode(&todo); err != nil {
			continue
		}
		todos = append(todos, todo)
	}
	c.JSON(http.StatusOK, todos)
}

func getTodoHandler(c *gin.Context) {
	var todo Todo
	idStr := c.Param("id")

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	filter := bson.M{"_id": id}
	err = collection.FindOne(context.Background(), filter).Decode(&todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	} else {
		c.JSON(http.StatusOK, todo)
		return
	}
}

func updateTodoHandler(c *gin.Context) {
	var todo Todo
	var inputTodo Todo

	if err := c.ShouldBindJSON(&inputTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	idStr := c.Param("id")

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	filter := bson.M{"_id": id}
	err = collection.FindOne(context.Background(), filter).Decode(&todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"title":     inputTodo.Title,
			"completed": inputTodo.Completed,
			"createdAt": time.Now(),
		},
	}
	updateResult, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Todo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"updatedCount": updateResult.ModifiedCount, "updatedTodo": todo})
}

func deleteTodoHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	result, err := collection.DeleteOne(context.Background(), bson.D{{"_id", id}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func initMongoDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	collection = client.Database("local").Collection("todoapp")
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/todos", createTodoHandler)
	router.GET("/todos", listTodoHandler)
	router.GET("/todos/:id", getTodoHandler)
	router.PUT("/todos/:id", updateTodoHandler)
	router.DELETE("/todos/:id", deleteTodoHandler)

	return router
}

func main() {
	initMongoDB()
	gin.SetMode(gin.ReleaseMode)

	router := setupRouter()
	err := router.Run(":7070")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
