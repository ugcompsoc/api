package server

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	h "github.com/nuigcompsoc/api/internal/helpers"
)

func (s *Server) RootGet(c *gin.Context) {
	h.RespondWithString(c, 200, "Homepage that won't be hosted by api")
}

/***************************
 *
 * === MISC V1 ENDPOINTS ===
 *
 ***************************/

func (s *Server) MiscV1BrewGet(c *gin.Context) {
	h.RespondWithError(c, 418, errors.New("I refuse to brew coffee because I am, permanently, a teapot."))
	return
}

func (s *Server) MiscV1PingGet(c *gin.Context) {
	h.RespondWithString(c, 200, "Pong!")
	return
}

/***************************
 *
 * == EVENTS V1 ENDPOINTS ==
 *
 ***************************/

func (s *Server) EventsV1Get(c *gin.Context) {
	events, err := s.Datastore.GetAllEvents()
	if err != nil {
		h.RespondWithError(c, 500, errors.New("failed to query database for events"))
		return
	}

	h.RespondWithJSON(c, 200, events)
	return
}

func (s *Server) EventsV1UpcomingGet(c *gin.Context) {
	events, err := s.Datastore.GetAllUpcomingEvents()
	if err != nil {
		h.RespondWithError(c, 500, errors.New("failed to query database for events"))
		return
	}

	h.RespondWithJSON(c, 200, events)
	return
}

func (s *Server) EventsV1UpcomingSocIDGet(c *gin.Context) {
	socID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.RespondWithError(c, 400, errors.New("could not convert socID into integer"))
		return
	}

	events, err := s.Datastore.GetAllUpcomingEventsForSocID(socID)
	if err != nil {
		h.RespondWithError(c, 500, errors.New("failed to query database for events"))
		return
	}

	h.RespondWithJSON(c, 200, events)
	return
}

func (s *Server) EventsV1PastGet(c *gin.Context) {
	events, err := s.Datastore.GetAllPastEvents()
	if err != nil {
		h.RespondWithError(c, 500, errors.New("failed to query database for events"))
		return
	}

	h.RespondWithJSON(c, 200, events)
	return
}

func (s *Server) EventsV1PastSocIDGet(c *gin.Context) {
	socID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.RespondWithError(c, 400, errors.New("could not convert socID into integer"))
		return
	}

	events, err := s.Datastore.GetAllPastEventsForSocID(socID)
	if err != nil {
		h.RespondWithError(c, 500, errors.New("failed to query database for events"))
		return
	}

	h.RespondWithJSON(c, 200, events)
	return
}
