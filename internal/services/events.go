package services

import (
	b64 "encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/nuigcompsoc/api/internal/config"
	"github.com/nuigcompsoc/api/internal/models"

	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

type EventService struct {
	Datastore                    *MongoDatastore
	Endpoint                     string
	EventService                 string
	EventServiceMethodAll        string
	EventServiceMethodIndividual string
	EventServiceAction           string
}

func NewEventService(config *config.Config, datastore *MongoDatastore) *EventService {
	return &EventService{
		Datastore:                    datastore,
		Endpoint:                     config.SocsPortal.Endpoint,
		EventService:                 config.SocsPortal.EventService,
		EventServiceMethodAll:        config.SocsPortal.EventServiceMethodAll,
		EventServiceMethodIndividual: config.SocsPortal.EventServiceMethodIndividual,
		EventServiceAction:           config.SocsPortal.EventServiceAction,
	}
}

// This contacts the Socs Portal to get the list of past and upcoming events.
// We then have no use for the rest of the data as it doesn't give us enough
// details anyways. So we've to contact the Socs Portal again on a seperate
// API for every eventID we get. This function just returns an array of
// eventDetailsIDs as strings.
func (s *EventService) GetEventsForSocID(socID string) ([]models.Event, error) {
	/*
		A response from the societies portal will look like this:
		[
			{
				"eventDetailsID": 34719,
				"eventID": 16913,
				"title": "Compsoc",
				"descriptionAbbrev": "SOCs Day",
				"ownerTitle": "Compsoc",
				"start": "2022-09-07T12:00",
				"end": "2022-09-07T17:00",
				"locationDetails": " Aras an Mac Leinn",
				"ownerID": 30,
				"allDay": false,
				"icon": "fa-home",
				"eventLocationType": " home",
				"className": "ic_other On Campus "
			},
			{
				"eventDetailsID": ...
			}
		]
	*/

	req, err := http.NewRequest("GET", s.Endpoint, nil)
	if err != nil {
		log.WithField("error", err.Error()).Warn("Could not create a request to SocsPortal endpoint")
		return nil, err
	}

	q := req.URL.Query()
	q.Add("object", b64.StdEncoding.EncodeToString([]byte(s.EventService)))
	q.Add("method", b64.StdEncoding.EncodeToString([]byte(s.EventServiceMethodAll)))
	q.Add("action", b64.StdEncoding.EncodeToString([]byte(s.EventServiceAction)))
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithField("error", err.Error()).Warn("Could not make a request to SocsPortal endpoint")
		return nil, err
	}
	defer res.Body.Close()

	var data []models.Event
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		log.WithField("error", err.Error()).Warn("Could not decode JSON response from Socs Portal into interface")
		return nil, err
	}

	return data, nil
}

// Here we're gettings the eventDetailsIDs and contacting the Socs Poral to get
// details on every event. The idea then is to save them to the database so we're
// not annoying Socs Portal every time we want to find out about our events.
func (s *EventService) GetAllEventsDetails(eventDetailIDs []int) (map[int]models.EventDetails, error) {
	/*
		A response from the societies portal will look like this:
		[
			{
				"eventDetailsID": 34719,
				"eventID": 16913,
				"title": "SOCs Day",
				"descriptionHTML": "<p><br /></p>",
				"description": "",
				"eventTypeTitle": "Other",
				"start": "2022-09-07T12:00",
				"end": "2022-09-07T17:00",
				"locationDetails": " Aras an Mac Leinn",
				"startDateTimeFormatted": "Wednesday, 12:00 to 17:00 Sep 07th 2022</span>",
				"ownerID": 30,
				"ownerTitle": "Compsoc",
				"allDay": false,
				"eventLocationGroupID": 0,
				"tags": null,
				"locationTypeTitle": "On Campus",
				"statusTypeTitle": "Approved",
				"signUpUrl": "",
				"icon": "fa-bank",
				"eventLocationType": " home",
				"className": "bg-color-pinkDark Other On Campus ",
				"eventUrl": "",
				"eventReadUrl": "calendar.php?object=Q2FsZW5kYXI=&method=ZXZlbnRSZWFkVmlldw==&action=Ng==&eventDetailsID=MzQ3MTk=&view=&ownerID=MzA=",
				"eventICalUrl": "https://socs.nuigalway.ie/calendar.php?object=Q2FsZW5kYXJTaGFyaW5n&method=ZXZlbnRUb0ljYWw=&action=Ng==&eventDetailsID=MzQ3MTk="
			}
		]
	*/

	eventsDetails := map[int]models.EventDetails{}
	for _, eventDetailsID := range eventDetailIDs {
		req, err := http.NewRequest("GET", s.Endpoint, nil)
		if err != nil {
			log.WithField("error", err.Error()).Warn("Could not create a request to SocsPortal endpoint")
			return nil, err
		}

		q := req.URL.Query()
		q.Add("object", b64.StdEncoding.EncodeToString([]byte(s.EventService)))
		q.Add("method", b64.StdEncoding.EncodeToString([]byte(s.EventServiceMethodIndividual)))
		q.Add("eventDetailsID", strconv.Itoa(eventDetailsID))
		q.Add("action", b64.StdEncoding.EncodeToString([]byte(s.EventServiceAction)))
		req.URL.RawQuery = q.Encode()

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.WithField("error", err.Error()).Warn("Could not make a request to SocsPortal endpoint")
			return nil, err
		}
		defer res.Body.Close()

		data := []models.EventDetails{}
		err = json.NewDecoder(res.Body).Decode(&data)
		if err != nil {
			log.WithField("error", err.Error()).Warn("Could not decode JSON response from Socs Portal into interface")
			return nil, err
		}
		eventsDetails[data[0].EventDetailsID] = data[0]
	}

	return eventsDetails, nil
}

func (s *EventService) GetAllEvents(onlyUpcomingEvents bool) ([]models.Event, error) {
	req, err := http.NewRequest("GET", s.Endpoint, nil)
	if err != nil {
		log.WithField("error", err.Error()).Warn("Could not create a request to SocsPortal endpoint")
		return nil, err
	}

	q := req.URL.Query()
	q.Add("object", b64.StdEncoding.EncodeToString([]byte(s.EventService)))
	q.Add("method", b64.StdEncoding.EncodeToString([]byte(s.EventServiceMethodAll)))
	q.Add("action", b64.StdEncoding.EncodeToString([]byte(s.EventServiceAction)))
	if !onlyUpcomingEvents {
		log.Info("Requesting only upcoming events")
		q.Add("start", time.Now().UTC().Format(time.RFC3339))
		q.Add("end", time.Now().UTC().AddDate(1, 0, 0).Format(time.RFC3339))
	} else {
		log.Info("Requesting all events")
	}
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithField("error", err.Error()).Warn("Could not make a request to SocsPortal endpoint")
		return nil, err
	}
	defer res.Body.Close()

	events := []models.Event{}
	err = json.NewDecoder(res.Body).Decode(&events)
	if err != nil {
		log.WithField("error", err.Error()).Warn("Could not decode JSON response from Socs Portal into interface")
		return nil, err
	}

	societies, err := s.Datastore.GetAllSocieties()
	if err != nil {
		return nil, err
	}

	socIDs := make([]int, len(societies))
	index := 0
	for _, society := range societies {
		socIDs[index] = int(society.SocietiesPortalID)
		index++
	}

	eventsWeWant := []models.Event{}
	for _, event := range events {
		if slices.Contains(socIDs, event.OwnerID) {
			eventsWeWant = append(eventsWeWant, event)
		}
	}

	return eventsWeWant, nil
}
