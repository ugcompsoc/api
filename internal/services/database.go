package services

import (
	"context"
	"sync"
	"time"

	"github.com/nuigcompsoc/api/internal/config"
	"github.com/nuigcompsoc/api/internal/models"

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

func NewDatastore(config *config.Config) *MongoDatastore {

	var mongoDataStore *MongoDatastore

	db, session := connect(config)
	if db != nil && session != nil {
		mongoDataStore = new(MongoDatastore)
		mongoDataStore.db = db
		mongoDataStore.Session = session
		return mongoDataStore
	}

	log.Fatal("Failed to create Datastore")
	return nil
}

func connect(config *config.Config) (a *mongo.Database, b *mongo.Client) {

	var connectOnce sync.Once
	var db *mongo.Database
	var session *mongo.Client

	connectOnce.Do(func() {
		db, session = connectToMongo(config)
	})

	return db, session
}

func connectToMongo(config *config.Config) (a *mongo.Database, b *mongo.Client) {

	credential := options.Credential{
		Username: config.Database.Username,
		Password: config.Database.Password,
	}
	session, err := mongo.NewClient(options.Client().ApplyURI(config.Database.Host).SetAuth(credential))
	if err != nil {
		log.WithField("error", err).Warn("Failed to connect to Database Host")
		return nil, nil
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = session.Connect(ctx)
	if err != nil {
		log.WithField("error", err).Warn("Failed to create session with Database Host")
		return nil, nil
	}

	err = session.Ping(ctx, nil)
	if err != nil {
		log.WithField("error", err).Warn("Failed to ping the Database Host")
		return nil, nil
	}

	var DB = session.Database(config.Database.Name)
	log.Info("Connected to Database")

	return DB, session
}

/*
 *	Society Database Helpers
 */
func (ds *MongoDatastore) UpsertSociety(society models.Society) error {

	var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.D{{Key: "name", Value: society.Name}}
	update := bson.D{{Key: "$set", Value: society}}
	opts := options.Update().SetUpsert(true)

	result, err := ds.db.Collection("societies").UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.WithField("error", err).Warn("Failed to update Society %v", society.Name)
	}

	log.Debug("Number of documents updated: %v\n", result.ModifiedCount)
	log.Debug("Number of documents upserted: %v\n", result.UpsertedCount)

	return nil
}

func (ds *MongoDatastore) GetAllSocieties() (map[string]models.Society, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := ds.db.Collection("societies").Find(ctx, bson.D{})
	if err != nil {
		log.WithField("error", err).Debug("Failed to return cursor to find all documents in societies collection")
		return nil, err
	}

	societies := []models.Society{}
	err = cursor.All(ctx, &societies)
	if err != nil {
		log.WithField("error", err).Debug("Failed to use cursor to find all documents in societies collection")
		return nil, err
	}

	societiesMap := map[string]models.Society{}
	for _, society := range societies {
		societiesMap[society.Name] = society
	}

	return societiesMap, nil
}

func (ds *MongoDatastore) GetSocietyBySocietyName(societyName string) (*models.Society, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var society models.Society
	err := ds.db.Collection("societies").FindOne(ctx, bson.D{{Key: "name", Value: societyName}}).Decode(&society)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Info("Society %v not found", societyName)
		} else {
			log.WithField("error", err).Warn("Failed to return single society from societies collection")
			return nil, err
		}
	}

	return &society, nil
}

/*
 *	Event Database Helpers
 */
func (ds *MongoDatastore) UpsertEvents(events []models.DatabaseEvent) error {

	writeModels := []mongo.WriteModel{}
	for _, event := range events {
		writeModels = append(writeModels,
			mongo.NewUpdateOneModel().SetFilter(
				bson.M{"event_id": event.EventID, "event_details_id": event.EventDetailsID}).SetUpdate(
				bson.D{{Key: "$set", Value: event}}).SetUpsert(true))
	}
	opts := options.BulkWrite().SetOrdered(false)

	var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	results, err := ds.db.Collection("events").BulkWrite(ctx, writeModels, opts)
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

func (ds *MongoDatastore) GetAllEvents() ([]models.EventDetails, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := ds.db.Collection("events").Find(ctx, bson.D{})
	if err != nil {
		log.WithField("error", err).Warn("Failed to return cursor to find all documents in events collection")
		return nil, err
	}

	events := []models.EventDetails{}
	err = cursor.All(ctx, &events)
	if err != nil {
		log.WithField("error", err).Warn("Failed to use cursor to find all documents in events collection")
		return nil, err
	}

	return events, nil
}

func (ds *MongoDatastore) GetAllUpcomingEvents() ([]models.DatabaseEvent, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Find upcoming events, but also include events that ended at most an hour ago
	cursor, err := ds.db.Collection("events").Find(ctx, bson.M{"end_datetime": bson.M{
		"$gte": time.Now().UTC().Add(time.Hour).Format(time.RFC3339),
	}})
	if err != nil {
		log.WithField("error", err).Warn("Failed to return cursor to find all documents that are upcoming in events collection")
		return nil, err
	}

	events := []models.DatabaseEvent{}
	err = cursor.All(ctx, &events)
	if err != nil {
		log.WithField("error", err).Warn("Failed to use cursor to find all documents that are upcoming in events collection")
		return nil, err
	}

	return events, nil
}

func (ds *MongoDatastore) GetAllPastEvents() ([]models.DatabaseEvent, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Find past events, but also not including events that ended at most an hour ago
	cursor, err := ds.db.Collection("events").Find(ctx, bson.M{"end_datetime": bson.M{
		"$lt": time.Now().UTC().Add(-time.Hour).Format(time.RFC3339),
	}})
	if err != nil {
		log.WithField("error", err).Warn("Failed to return cursor to find all documents that are past in events collection")
		return nil, err
	}

	events := []models.DatabaseEvent{}
	err = cursor.All(ctx, &events)
	if err != nil {
		log.WithField("error", err).Warn("Failed to use cursor to find all documents that are past in events collection")
		return nil, err
	}

	return events, nil
}

// TODO Theres a bug in here somewhere relating to recent (within an hour) events not showing up
func (ds *MongoDatastore) GetAllUpcomingEventsForSocID(socID int) ([]models.DatabaseEvent, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Find upcoming events for society, but also include events that ended at most an hour ago
	cursor, err := ds.db.Collection("events").Find(ctx,
		bson.M{
			"end_datetime": bson.M{
				"$gte": time.Now().UTC().Add(time.Hour).Format(time.RFC3339),
			},
			"society_id": socID,
		})
	if err != nil {
		log.WithFields(log.Fields{"error": err, "socID": socID}).Warn("Failed to return cursor to find all documents for society that are upcoming in events collection")
		return nil, err
	}

	events := []models.DatabaseEvent{}
	err = cursor.All(ctx, &events)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "socID": socID}).Warn("Failed to use cursor to find all documents for society that are upcoming in events collection")
		return nil, err
	}

	return events, nil
}

func (ds *MongoDatastore) GetAllPastEventsForSocID(socID int) ([]models.DatabaseEvent, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Find past events for society, but also not including events that ended at most an hour ago
	cursor, err := ds.db.Collection("events").Find(ctx,
		bson.M{
			"end_datetime": bson.M{
				"$lt": time.Now().UTC().Add(-time.Hour).Format(time.RFC3339),
			},
			"society_id": socID,
		})
	if err != nil {
		log.WithFields(log.Fields{"error": err, "socID": socID}).Warn("Failed to return cursor to find all documents for society that are passed in events collection")
		return nil, err
	}

	events := []models.DatabaseEvent{}
	err = cursor.All(ctx, &events)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "socID": socID}).Warn("Failed to use cursor to find all documents for society that are passed in events collection")
		return nil, err
	}

	return events, nil
}
