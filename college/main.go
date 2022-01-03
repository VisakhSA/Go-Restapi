package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type College struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CollegeId   int64              `json:"collegeid" bson:"collegeid"`
	CollegeName string             `json:"collegename" bson:"collegename"`
}

func getCollege(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	collection := client.Database("collegedb").Collection("college")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	params := mux.Vars(r)
	if len(params) != 0 {
		var college College
		id, _ := strconv.ParseInt(params["id"], 0, 64)
		err := collection.FindOne(ctx, bson.M{"collegeid": id}).Decode(&college)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		json.NewEncoder(w).Encode(college)
		return
	}
	var colleges []College
	cursor, _ := collection.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var college College
		cursor.Decode(&college)
		colleges = append(colleges, college)
	}
	json.NewEncoder(w).Encode(colleges)
}

var client *mongo.Client

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	r := mux.NewRouter()
	/*collegeList = append(collegeList, College{CollegeId: 1, CollegeName: "R.M.K"})
	collegeList = append(collegeList, College{CollegeId: 2, CollegeName: "R.M.D"})*/
	r.HandleFunc("/api/getCollege", getCollege).Methods("GET")
	r.HandleFunc("/api/getCollege/{id}", getCollege).Methods("GET")
	/*r.HandleFunc("/api/addCollege", addCollege).Methods("POST")
	r.HandleFunc("/api/updateCollege/{id}", updateCollege).Methods("PUT")
	r.HandleFunc("/api/deleteCollege/{id}", deleteCollege).Methods("DELETE")*/
	log.Fatal(http.ListenAndServe(":8080", r))
}
