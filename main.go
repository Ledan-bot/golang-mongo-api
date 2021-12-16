package main

import (
	"context"
	"encoding/json"
	"fmt"
	"golang-mongo-api/env"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Restaurant struct {
	id            primitive.ObjectID
	borough       string
	cuisine       string
	name          string
	restaurant_id string
}

var client *mongo.Client

func GetRestaurants(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var places []Restaurant
	collection := client.Database("sample_restaurants").Collection("restaurants")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var rest Restaurant
		cursor.Decode((&rest))
		places = append(places, rest)
	}

	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	json.NewEncoder(w).Encode(places)
}

func main() {
	fmt.Println("Starting...")
	client, err := mongo.NewClient(options.Client().ApplyURI(env.Connect))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
	router := mux.NewRouter()
	router.HandleFunc("/restaurants", GetRestaurants).Methods("GET")
	http.ListenAndServe(":8033", router)
}
