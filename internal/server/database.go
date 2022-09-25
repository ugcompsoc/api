package server

import (
	"context"
	"html"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatastore struct {
	db      *mongo.Database
	Session *mongo.Client
}

/*
 *	Database Setup
 */

func (s *Server) NewDatastore() *MongoDatastore {

	var mongoDataStore *MongoDatastore

	db, session := s.connect()
	if db != nil && session != nil {
		mongoDataStore = new(MongoDatastore)
		mongoDataStore.db = db
		mongoDataStore.Session = session
		return mongoDataStore
	}

	log.Fatal("Failed to create Datastore")
	return nil
}

func (s *Server) connect() (a *mongo.Database, b *mongo.Client) {

	var connectOnce sync.Once
	var db *mongo.Database
	var session *mongo.Client

	connectOnce.Do(func() {
		db, session = s.connectToMongo()
	})

	return db, session
}

func (s *Server) connectToMongo() (a *mongo.Database, b *mongo.Client) {

	credential := options.Credential{
		Username: s.Config.Database.Username,
		Password: s.Config.Database.Password,
	}
	session, err := mongo.NewClient(options.Client().ApplyURI(s.Config.Database.Host).SetAuth(credential))
	if err != nil {
		log.Warn("Failed to connect to Database Host")
		return nil, nil
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = session.Connect(ctx)
	if err != nil {
		log.Warn("Failed to create session with Database Host")
		return nil, nil
	}

	err = session.Ping(ctx, nil)
	if err != nil {
		log.Warn("Failed to ping the Database Host")
		return nil, nil
	}

	var DB = session.Database(s.Config.Database.Name)
	log.Info("Connected to Database")

	return DB, session
}

/*
 *	Database Helpers
 */
func (ds *MongoDatastore) upsertEvents(events []EventDetails) error {

	models := []mongo.WriteModel{}
	for _, event := range events {
		// TODO should verify if html is already escaped
		// already escaped event.DescriptionHTML = html.EscapeString(event.DescriptionHTML)
		event.StartDateTimeFormatted = html.EscapeString(event.StartDateTimeFormatted)
		models = append(models,
			mongo.NewUpdateOneModel().SetFilter(
				bson.M{"eventID": event.EventID, "eventDetailsID": event.EventDetailsID}).SetUpdate(
				bson.D{{Key: "$set", Value: event}}).SetUpsert(true))
	}
	opts := options.BulkWrite().SetOrdered(false)

	var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	results, err := ds.db.Collection("events").BulkWrite(ctx, models, opts)
	if err != nil {
		log.WithField("error", err).Warn("Failed to BulkWrite events to collection")
		return err
	}

	log.Info("Inserted/Updated events to events collection")
	log.Info("Number of documents upserted: ", results.UpsertedCount)
	log.Info("Number of documents inserted: ", results.InsertedCount)
	log.Info("Number of documents modified: ", results.ModifiedCount)
	log.Info("Number of documents matched: ", results.MatchedCount)

	return nil
}

func (ds *MongoDatastore) getAllEvents() ([]EventDetails, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := ds.db.Collection("events").Find(ctx, bson.D{})
	if err != nil {
		log.WithField("error", err).Warn("Failed to return cursor to find all documents in events collection")
		return nil, err
	}

	events := []EventDetails{}
	err = cursor.All(ctx, &events)
	if err != nil {
		log.WithField("error", err).Warn("Failed to use cursor to find all documents in events collection")
		return nil, err
	}

	return events, nil
}

func (ds *MongoDatastore) getAllUpcomingEvents() ([]EventDetails, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Find upcoming events, but also include events that ended at most an hour ago
	cursor, err := ds.db.Collection("events").Find(ctx, bson.M{"end": bson.M{
		"$gte": time.Now().UTC().Add(time.Hour).Format(time.RFC3339),
	}})
	if err != nil {
		log.WithField("error", err).Warn("Failed to return cursor to find all documents that are upcoming in events collection")
		return nil, err
	}

	events := []EventDetails{}
	err = cursor.All(ctx, &events)
	if err != nil {
		log.WithField("error", err).Warn("Failed to use cursor to find all documents that are upcoming in events collection")
		return nil, err
	}

	return events, nil
}

func (ds *MongoDatastore) getAllPastEvents() ([]EventDetails, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Find past events, but also not including events that ended at most an hour ago
	cursor, err := ds.db.Collection("events").Find(ctx, bson.M{"end": bson.M{
		"$lt": time.Now().UTC().Add(-time.Hour).Format(time.RFC3339),
	}})
	if err != nil {
		log.WithField("error", err).Warn("Failed to return cursor to find all documents that are past in events collection")
		return nil, err
	}

	events := []EventDetails{}
	err = cursor.All(ctx, &events)
	if err != nil {
		log.WithField("error", err).Warn("Failed to use cursor to find all documents that are past in events collection")
		return nil, err
	}

	return events, nil
}

func (ds *MongoDatastore) getAllUpcomingEventsForSocID(socID int) ([]EventDetails, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Find upcoming events for society, but also include events that ended at most an hour ago
	cursor, err := ds.db.Collection("events").Find(ctx,
		bson.M{
			"end": bson.M{
				"$gte": time.Now().UTC().Add(time.Hour).Format(time.RFC3339),
			},
			"ownerID": socID,
		})
	if err != nil {
		log.WithFields(log.Fields{"error": err, "socID": socID}).Warn("Failed to return cursor to find all documents for society that are upcoming in events collection")
		return nil, err
	}

	events := []EventDetails{}
	err = cursor.All(ctx, &events)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "socID": socID}).Warn("Failed to use cursor to find all documents for society that are upcoming in events collection")
		return nil, err
	}

	return events, nil
}

func (ds *MongoDatastore) getAllPastEventsForSocID(socID int) ([]EventDetails, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Find past events for society, but also not including events that ended at most an hour ago
	cursor, err := ds.db.Collection("events").Find(ctx,
		bson.M{
			"end": bson.M{
				"$lt": time.Now().UTC().Add(-time.Hour).Format(time.RFC3339),
			},
			"ownerID": socID,
		})
	if err != nil {
		log.WithFields(log.Fields{"error": err, "socID": socID}).Warn("Failed to return cursor to find all documents for society that are passed in events collection")
		return nil, err
	}

	events := []EventDetails{}
	err = cursor.All(ctx, &events)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "socID": socID}).Warn("Failed to use cursor to find all documents for society that are passed in events collection")
		return nil, err
	}

	return events, nil
}
