package models

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
	EventDetailsID         int    `json:"eventDetailsID" bson:"eventDetailsID"`
	EventID                int    `json:"eventID" bson:"eventID"`
	Title                  string `json:"title" bson:"title"`
	DescriptionHTML        string `json:"descriptionHTML" bson:"descriptionHTML"`
	Description            string `json:"description" bson:"description"`
	EventTypeTitle         string `json:"eventTypeTitle" bson:"eventTypeTitle"`
	Start                  string `json:"start" bson:"start"`
	End                    string `json:"end" bson:"end"`
	LocationDetails        string `json:"locationDetails" bson:"locationDetails"`
	StartDateTimeFormatted string `json:"startDateTimeFormatted" bson:"startDateTimeFormatted"`
	OwnerID                int    `json:"ownerID" bson:"ownerID"`
	OwnerTitle             string `json:"ownerTitle" bson:"ownerTitle"`
	AllDay                 bool   `json:"allDay" bson:"allDay"`
	EventLocationGroupID   int    `json:"eventLocationGroupID" bson:"eventLocationGroupID"`
	Tags                   string `json:"tags" bson:"tags"`
	LocationTypeTitle      string `json:"locationTypeTitle" bson:"locationTypeTitle"`
	StatusTypeTitle        string `json:"statusTypeTitle" bson:"statusTypeTitle"`
	SignUpUrl              string `json:"signUpUrl" bson:"signUpUrl"`
	Icon                   string `json:"icon" bson:"icon"`
	EventLocationType      string `json:"eventLocationType" bson:"eventLocationType"`
	ClassName              string `json:"className" bson:"className"`
	EventUrl               string `json:"eventUrl" bson:"eventUrl"`
	EventReadUrl           string `json:"eventReadUrl" bson:"eventReadUrl"`
	EventICalUrl           string `json:"eventICalUrl" bson:"eventICalUrl"`
}
