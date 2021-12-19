package main

import (
	"context"
	"encoding/json"
	"fmt"
	"golang-mongo-api/env"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func GetRestaurants(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var places []bson.M
	coll := client.Database("sample_restaurants").Collection("restaurants")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var result bson.M
		err := cursor.Decode((&result))
		if err != nil {
			fmt.Println("cursor.Next() error: ", err)
			os.Exit(1)
		} else {
			places = append(places, result)
		}
	}

	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	json.NewEncoder(w).Encode(places)
}
func GetRestaurantByName(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	params := mux.Vars(r)
	name := params["name"]
	fmt.Println(name)

	coll := client.Database("sample_restaurants").Collection("restaurants")

	var result bson.M
	s := coll.FindOne(context.TODO(), bson.D{{"name", name}}).Decode(&result)
	if s != nil {
		if s == mongo.ErrNoDocuments {
			return
		}
		panic(s)
	}
	json.NewEncoder(w).Encode(result)
}

func main() {
	fmt.Println("Starting...")
	router := mux.NewRouter()
	clientOptions := options.Client().ApplyURI(env.Connect)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(context.TODO(), clientOptions)

	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("HERE")
	fmt.Println(databases)
	router.HandleFunc("/restaurants", GetRestaurants).Methods("GET")
	router.HandleFunc("/restaurant/{name}", GetRestaurantByName).Methods("GET")
	log.Fatal(http.ListenAndServe(":8033", router))
}
