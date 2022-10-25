package models

type DatabaseEvent struct {
	EventID                  int    `bson:"event_id, omitempty"`
	EventDetailsID           int    `bson:"event_details_id, omitempty"`
	Title                    string `bson:"title, omitempty"`
	SocietyID                int    `bson:"society_id, omitempty"`
	SocietyName              string `bson:"society_name, omitempty"`
	Location                 string `bson:"location, omitempty"`
	Description              string `bson:"description, omitempty"`
	DescriptionMarkdown      string `bson:"description_markdown, omitempty"`
	DangerousDescriptionHTML string `bson:"dangerous_description_html, omitempty"`
	StartDatetime            string `bson:"start_datetime, omitempty"`
	EndDatetime              string `bson:"end_datetime, omitempty"`
	DatetimeFormatted        string `bson:"datetime_formatted, omitempty"`
	EventURL                 string `bson:"event_url, omitempty"`
	EventICalURL             string `bson:"event_ical_url, omitempty"`
}

type Society struct {
	Name              string
	SocietiesPortalID int32
}
