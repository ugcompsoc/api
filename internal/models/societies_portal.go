package models

import (
	"html"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/microcosm-cc/bluemonday"
	log "github.com/sirupsen/logrus"
)

type SocietyMember struct {
	MemberTypeTitle string `json:"MemberTypeTitle"`
	MemberID        string `json:"MemberID"`
	FirstName       string `json:"FirstName"`
	LastName        string `json:"LastName"`
	Email           string `json:"Email"`
	PhoneNumber     string `json:"PhoneNumber"`
}

type Event struct {
	EventDetailsID    int    `json:"eventDetailsID"`
	EventID           int    `json:"eventID"`
	Title             string `json:"title"`
	DescriptionAbbrev string `json:"descriptionAbbrev"`
	OwnerTitle        string `json:"ownerTitle"`
	Start             string `json:"start"`
	End               string `json:"end"`
	LocationDetails   string `json:"locationDetails"`
	OwnerID           int    `json:"ownerID"`
	AllDay            bool   `json:"allDay"`
	Icon              string `json:"icon"`
	EventLocationType string `json:"eventLocationType"`
	ClassName         string `json:"className"`
}

type EventDetails struct {
	EventDetailsID         int    `json:"eventDetailsID"`
	EventID                int    `json:"eventID"`
	Title                  string `json:"title"`
	DescriptionHTML        string `json:"descriptionHTML"`
	Description            string `json:"description"`
	EventTypeTitle         string `json:"eventTypeTitle"`
	Start                  string `json:"start"`
	End                    string `json:"end"`
	LocationDetails        string `json:"locationDetails"`
	StartDateTimeFormatted string `json:"startDateTimeFormatted"`
	OwnerID                int    `json:"ownerID"`
	OwnerTitle             string `json:"ownerTitle"`
	AllDay                 bool   `json:"allDay"`
	EventLocationGroupID   int    `json:"eventLocationGroupID"`
	Tags                   string `json:"tags"`
	LocationTypeTitle      string `json:"locationTypeTitle"`
	StatusTypeTitle        string `json:"statusTypeTitle"`
	SignUpUrl              string `json:"signUpUrl"`
	Icon                   string `json:"icon"`
	EventLocationType      string `json:"eventLocationType"`
	ClassName              string `json:"className"`
	EventUrl               string `json:"eventUrl"`
	EventReadUrl           string `json:"eventReadUrl"`
	EventICalUrl           string `json:"eventICalUrl"`
}

func (e EventDetails) ToDatabaseEvent() DatabaseEvent {

	pEasy := bluemonday.UGCPolicy()
	pStrict := bluemonday.StrictPolicy()
	mdConverter := md.NewConverter("", true, nil)
	descriptionMarkdown, err := mdConverter.ConvertString(html.UnescapeString(e.DescriptionHTML))
	if err != nil {
		log.Info(err)
	}

	return DatabaseEvent{
		EventID:                  e.EventID,
		EventDetailsID:           e.EventDetailsID,
		Title:                    pStrict.Sanitize(e.Title),
		SocietyID:                e.OwnerID,
		SocietyName:              pStrict.Sanitize(e.OwnerTitle),
		Location:                 pStrict.Sanitize(e.LocationDetails),
		Description:              pStrict.Sanitize(html.UnescapeString(e.DescriptionHTML)),
		DescriptionMarkdown:      descriptionMarkdown,
		DangerousDescriptionHTML: pEasy.Sanitize(html.UnescapeString(e.DescriptionHTML)),
		StartDatetime:            pStrict.Sanitize(e.Start),
		EndDatetime:              pStrict.Sanitize(e.End),
		DatetimeFormatted:        pStrict.Sanitize(e.StartDateTimeFormatted),
		EventURL:                 pStrict.Sanitize("https://socs.nuigalway.ie/" + e.EventReadUrl),
		EventICalURL:             pStrict.Sanitize(e.EventICalUrl),
	}
}
