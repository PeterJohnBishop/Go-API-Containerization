package db

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Message struct {
	ID     string   `json:"id"`
	Sender string   `json:"sender"`
	Text   string   `json:"text"`
	Media  []string `json:"media"`
	Date   int64    `json:"date"` // func (t time.Time) UnixMilli() int64
}

type Chat struct {
	ID       string   `json:"id"`
	Users    []string `json:"users"`
	Messages []string `json:"messages"`
	Active   int64    `json:"active"`
}
type Alert struct {
	ActivePeriods   []ActivePeriod   `json:"active_period,omitempty"`
	InformedEntity  []EntitySelector `json:"informed_entity,omitempty"`
	Cause           string           `json:"cause,omitempty"`
	Effect          string           `json:"effect,omitempty"`
	HeaderText      Translation      `json:"header_text,omitempty"`
	DescriptionText Translation      `json:"description_text,omitempty"`
}
type ActivePeriod struct {
	Start int64 `json:"start,omitempty"`
	End   int64 `json:"end,omitempty"`
}
type EntitySelector struct {
	AgencyID  string `json:"agency_id,omitempty"`
	RouteID   string `json:"route_id,omitempty"`
	RouteType int    `json:"route_type,omitempty"`
	StopID    string `json:"stop_id,omitempty"`
}
type Translation struct {
	Text     string `json:"text"`
	Language string `json:"language,omitempty"`
}
type TripUpdate struct {
	Trip           TripDescriptor    `json:"trip"`
	Vehicle        VehicleDescriptor `json:"vehicle,omitempty"`
	StopTimeUpdate []StopTimeUpdate  `json:"stop_time_update,omitempty"`
	Timestamp      int64             `json:"timestamp,omitempty"`
}
type TripDescriptor struct {
	TripID               string `json:"trip_id,omitempty"`
	RouteID              string `json:"route_id,omitempty"`
	DirectionID          int    `json:"direction_id,omitempty"`
	ScheduleRelationship string `json:"schedule_relationship,omitempty"` // SCHEDULED, ADDED, CANCELED
}
type StopTimeUpdate struct {
	StopSequence         int           `json:"stop_sequence"`
	StopID               string        `json:"stop_id"`
	Arrival              StopTimeEvent `json:"arrival,omitempty"`
	Departure            StopTimeEvent `json:"departure,omitempty"`
	ScheduleRelationship string        `json:"schedule_relationship,omitempty"` // SCHEDULED, SKIPPED
}
type StopTimeEvent struct {
	Time int64 `json:"time,omitempty"`
}
type VehiclePosition struct {
	Trip          TripDescriptor    `json:"trip,omitempty"`
	Vehicle       VehicleDescriptor `json:"vehicle,omitempty"`
	Position      Position          `json:"position,omitempty"`
	StopID        string            `json:"stop_id,omitempty"`
	CurrentStatus int               `json:"current_status,omitempty"` // 0 = INCOMING_AT, 1 = STOPPED_AT, 2 = IN_TRANSIT_TO
	Timestamp     int64             `json:"timestamp,omitempty"`
}
type VehicleDescriptor struct {
	ID    string `json:"id,omitempty"`
	Label string `json:"label,omitempty"`
}
type Position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Bearing   float64 `json:"bearing,omitempty"`
}
