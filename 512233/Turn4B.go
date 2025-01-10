
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Attendee struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Age       int                `bson:"age" json:"age"`
	Gender    string             `bson:"gender" json:"gender"`
	Category  string             `bson:"category" json:"category"`
	RegisteredAt time.Time         `bson:"registeredAt" json:"registeredAt"`
}

type Event struct {
	Name      string
	Attendees []Attendee
	lock      sync.Mutex
}

var (
	events = make(map[string]Event)
	db     *mongo.Database
)

func init() {
	// Initialize MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	db = client.Database("eventregistration")
}

func (e *Event) RegisterAttendee(attendee Attendee) {
	e.lock.Lock()
	defer e.lock.Unlock()
	attendee.ID = primitive.NewObjectID()
	attendee.RegisteredAt = time.Now()
	e.Attendees = append(e.Attendees, attendee)
	// Save the attendee to MongoDB
	_, err := db.Collection("attendees").InsertOne(context.TODO(), attendee)
	if err != nil {
		log.Printf("Error saving attendee to MongoDB: %v", err)
	}
}

func (e *Event) GetAttendees() []Attendee {
	e.lock.Lock()
	defer e.lock.Unlock()
	return e.Attendees
}

func (e *Event) CountAttendees() int {
	return len(e.Attendees)
}

// SyncInMemoryData will load attendee data from MongoDB into memory
func SyncInMemoryData() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cur, err := db.Collection("attendees").Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error fetching attendees from MongoDB: %v", err)
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var attendee Attendee
		if err := cur.Decode(&attendee); err != nil {
			log.Printf("Error decoding attendee: %v", err)
			continue
		}
		event, ok := events[attendee.Category]
		if !ok {
			event = Event{Name: attendee.Category}
			events[attendee.Category] = event
		}
		event.RegisterAttendee(attendee)
	}
	if err := cur.Err(); err != nil {
		log.Printf("Error iterating attendees cursor: %v", err)