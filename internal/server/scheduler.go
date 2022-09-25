package server

import (
	"time"

	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
)

func (s *Server) doGetAllEvents() {
	log.Info("Starting doGetAllEvents Task")

	var allEvents []Event
	var err error

	// Once an hour we want to update all events (past and upcoming)
	if time.Now().UTC().Minute() > 0 && time.Now().UTC().Minute() <= 5 {
		allEvents, err = s.getAllEvents(true)
	} else {
		allEvents, err = s.getAllEvents(false)
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
	allEventDetails, err := s.getAllEventsDetails(eventDetailsIDs)
	if err != nil {
		log.Warn("getAllEventsDetails Function Failed")
		return
	}

	allEventsWithEventDetails := []EventDetails{}
	for _, event := range allEvents {
		eventWithEventDetails := allEventDetails[event.EventDetailsID]
		eventWithEventDetails.EventID = event.EventID
		allEventsWithEventDetails = append(allEventsWithEventDetails, eventWithEventDetails)
	}

	err = s.Datastore.upsertEvents(allEventsWithEventDetails)
	if err != nil {
		log.Warn("Datastore.upsertEvents Function Failed")
		return
	}
}

// TODO add a pause scheduler until timestap
func (s *Server) RunAllServices() *gocron.Scheduler {
	var doGetAllEventsTask = s.doGetAllEvents

	log.Info("Starting Scheduler")

	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every("5m").Do(doGetAllEventsTask)

	scheduler.StartAsync()

	return scheduler
}
