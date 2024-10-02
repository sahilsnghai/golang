package database

import (
	"context"
	"log"
	"time"

	"github.com/sahilsnghai/Project5/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var cString = "mongodb://localhost:27017"

type DB struct {
	client *mongo.Client
}

func Connect() *DB {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cString))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	return &DB{
		client: client,
	}
}

func (db *DB) GetJob(id string) *model.JobListing {

	jobCollec := db.client.Database("go-mongo").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_id, _ := primitive.ObjectIDFromHex(id)

	var jobListing model.JobListing

	err := jobCollec.FindOne(ctx, bson.M{"_id": _id}).Decode(&jobListing)
	if err != nil {
		log.Fatal(err)

	}
	return &jobListing
}

func (db *DB) GetJobs() []*model.JobListing {
	var JobListings []*model.JobListing
	jobCollec := db.client.Database("go-mongo").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := jobCollec.Find(ctx, bson.D{})

	if err != nil {
		log.Fatal(err)
	}

	if err = cursor.All(context.TODO(), &JobListings); err != nil {
		panic(err)
	}

	return JobListings
}

func (db *DB) CreateJobListing(input model.CreateJobListing) *model.JobListing {
	jobCollec := db.client.Database("go-mongo").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	inserted, err := jobCollec.InsertOne(ctx, bson.M{
		"title":       input.Title,
		"description": input.Description,
		"url":         input.URL,
		"company":     input.Company,
	})

	if err != nil {
		panic(err)
	}
	insertedId := inserted.InsertedID.(primitive.ObjectID).Hex()
	JobListing := model.JobListing{
		ID:          insertedId,
		Title:       input.Title,
		URL:         input.URL,
		Description: input.Description,
		Company:     input.Company,
	}

	return &JobListing

}

func (db *DB) UpdateJobListing(jobId string, input model.UpdateJobListing) *model.JobListing {
	jobCollec := db.client.Database("go-mongo").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	updateJobInfo := bson.M{}

	if input.Title != nil {
		updateJobInfo["title"] = input.Title
	}
	if input.Description != nil {
		updateJobInfo["description"] = input.Description
	}
	if input.URL != nil {
		updateJobInfo["url"] = input.URL
	}

	_id, _ := primitive.ObjectIDFromHex(jobId)
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateJobInfo}

	results := jobCollec.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var jobListing model.JobListing

	if err := results.Decode(&jobListing); err != nil {
		log.Fatal(err)
	}
	return &jobListing

}

func (db *DB) DeleteJobListing(jobId string) *model.DeleteJobResponse {
	jobCollec := db.client.Database("go-mongo").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(jobId)
	filter := bson.M{"_id": _id}
	_, err := jobCollec.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	var DeletedJobId model.DeleteJobResponse
	return &DeletedJobId
}
