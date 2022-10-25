package services

import (
	"time"

	"github.com/nuigcompsoc/api/internal/config"
	"github.com/nuigcompsoc/api/internal/models"

	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
)

type SchedulerService struct {
	Config    *config.Config
	Datastore *MongoDatastore
	Scheduler *gocron.Scheduler
}

func (s *SchedulerService) DoGetAllEvents() {
	log.Info("Starting doGetAllEvents Task")

	societiesPortalService := NewSocietiesPortalService(s.Config, s.Datastore)

	var allEvents []models.Event
	var err error

	// Once an hour we want to update all events (past and upcoming)
	if time.Now().UTC().Minute() > 0 && time.Now().UTC().Minute() <= 5 {
		allEvents, err = societiesPortalService.GetAllEvents(true)
	} else {
		allEvents, err = societiesPortalService.GetAllEvents(false)
	}

	if err != nil {
		log.Warn("getAllEvents Function Failed")
		return
	}

	eventDetailsIDs := []int{}
	for _, event := range allEvents {
		eventDetailsIDs = append(eventDetailsIDs, event.EventDetailsID)
	}

	// Here we're getting all the event details for every event and ingoring
	// duplicate eventDetailsID, no point duplicating work.
	allEventDetails, err := societiesPortalService.GetAllEventsDetails(eventDetailsIDs)
	if err != nil {
		log.Warn("getAllEventsDetails Function Failed")
		return
	}

	allEventsWithEventDetails := []models.EventDetails{}
	for _, event := range allEvents {
		eventWithEventDetails := allEventDetails[event.EventDetailsID]
		eventWithEventDetails.EventID = event.EventID
		allEventsWithEventDetails = append(allEventsWithEventDetails, eventWithEventDetails)
	}

	// convert societies portal events to database events
	allDatabaseEvents := []models.DatabaseEvent{}
	for _, event := range allEventsWithEventDetails {
		allDatabaseEvents = append(allDatabaseEvents, event.ToDatabaseEvent())
	}

	err = s.Datastore.UpsertEvents(allDatabaseEvents)
	if err != nil {
		log.Warn("Datastore.upsertEvents Function Failed")
		return
	}
}

func NewSchedulerService(config *config.Config) *SchedulerService {
	return &SchedulerService{
		Config:    config,
		Datastore: NewDatastore(config),
		Scheduler: gocron.NewScheduler(time.UTC),
	}
}

func (s *SchedulerService) RunAllServices() {
	var doGetAllEventsTask = s.DoGetAllEvents

	log.Info("Starting Scheduler")
	s.Scheduler.Every("5m").Do(doGetAllEventsTask)
	s.Scheduler.StartAsync()
}
