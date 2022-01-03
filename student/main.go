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

type Student struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CollegeId   int64              `json:"collegeid" bson:"collegeid"`
	StudentId   int64              `json:"studentid" bson:"studentid"`
	StudentName string             `json:"studentname" bson:"studentname"`
}

func getStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	collection := client.Database("studentdb").Collection("student")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	params := mux.Vars(r)
	if len(params) != 0 {
		var student Student
		id, _ := strconv.ParseInt(params["id"], 0, 64)
		err := collection.FindOne(ctx, bson.M{"studentid": id}).Decode(&student)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		json.NewEncoder(w).Encode(student)
		return
	}
	var students []Student
	cursor, _ := collection.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var student Student
		cursor.Decode(&student)
		students = append(students, student)
	}
	json.NewEncoder(w).Encode(students)
}

var client *mongo.Client

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	r := mux.NewRouter()
	/*collegeList = append(collegeList, College{CollegeId: 1, CollegeName: "R.M.K"})
	collegeList = append(collegeList, College{CollegeId: 2, CollegeName: "R.M.D"})*/
	r.HandleFunc("/api/getStudent", getStudent).Methods("GET")
	r.HandleFunc("/api/getStudent/{id}", getStudent).Methods("GET")
	/*r.HandleFunc("/api/addCollege", addCollege).Methods("POST")
	r.HandleFunc("/api/updateCollege/{id}", updateCollege).Methods("PUT")
	r.HandleFunc("/api/deleteCollege/{id}", deleteCollege).Methods("DELETE")*/
	log.Fatal(http.ListenAndServe(":8081", r))
}
