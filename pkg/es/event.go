package es

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

// Event is an interface that represents an event.
type Event interface {
	GetID() string
	GetEventType() string
	GetOccurredOn() time.Time
	GetData() []byte
	SetData(data []byte) *BaseEvent
	GetJsonData(data interface{}) error
	SetJsonData(data interface{}) error
}

// BaseEvent is an internal representation of an event.
type BaseEvent struct {
	EventID    string    `json:"event_id"`
	EventType  string    `json:"event_type"`
	Data       []byte    `json:"data"`
	OccurredOn time.Time `json:"occurred_on"`
}

// NewBaseEvent new base Event constructor with configured EventID, Aggregate properties and Timestamp.
func NewBaseEvent(name string) BaseEvent {
	return BaseEvent{
		EventID:    uuid.NewString(),
		EventType:  name,
		OccurredOn: time.Now(),
	}
}

// GetID returns the ID of the event.
func (e *BaseEvent) GetID() string {
	return e.EventID
}

// GetEventType returns the name of the event.
func (e *BaseEvent) GetEventType() string {
	return e.EventType
}

// GetOccurredOn returns the timestamp when the event occurred.
func (e *BaseEvent) GetOccurredOn() time.Time {
	return e.OccurredOn
}

// GetData The data attached to the Event serialized to bytes.
func (e *BaseEvent) GetData() []byte {
	return e.Data
}

// SetData add the data attached to the Event serialized to bytes.
func (e *BaseEvent) SetData(data []byte) *BaseEvent {
	e.Data = data
	return e
}

// GetJsonData json unmarshal data attached to the Event.
func (e *BaseEvent) GetJsonData(data interface{}) error {
	return json.Unmarshal(e.GetData(), data)
}

// SetJsonData serialize to json and set data attached to the Event.
func (e *BaseEvent) SetJsonData(data interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	e.Data = dataBytes
	return nil
}
