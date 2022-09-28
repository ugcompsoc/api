package server

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	h "github.com/nuigcompsoc/api/internal/helpers"
)

func (s *Server) RootGet(c *gin.Context) {
	h.RespondWithString(c, 200, "What version do you plan on using, huh?")
}

/***************************
 *
 * === AUTH V1 ENDPOINTS ===
 *
 ***************************/

func (s *Server) AuthV1RegisterPost(c *gin.Context) {
	// Check our ldap to see if their preferred username is already taken

	// Check socs portal to see if they are in the society

	// If they're a member, sign a token and send it off to them in an email

	// If not, tell them to go register
}

func (s *Server) AuthV1RegisterVerifyGet(c *gin.Context) {
	// Extract the claims from the token received (student ID & preferred username)

	// Do another check to verify them and check the username hasn't been taken

	// Register them in LDAP

	// Send out an email on account info
}

func (s *Server) AuthV1OpenIDGet(c *gin.Context) {
	// Redirect to CompSoc SSO
}

func (s *Server) AuthV1GoogleGet(c *gin.Context) {
	// Redirect to Society SSO
}

func (s *Server) AuthV1OpenIDCallbackGet(c *gin.Context) {
	// Get code from query params

	// Contact CompSoc SSO to get token

	// Set token received from CompSoc SSO
}

func (s *Server) AuthV1GoogleCallbackGet(c *gin.Context) {
	// Get code from query params

	// Contact Google SSO to get token

	// Extract token payload to get Society information (email)

	// Make a new LDAP society account if none exists

	// Set token received from Google SSO
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
