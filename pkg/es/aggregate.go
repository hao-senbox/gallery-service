package es

const (
	aggregateStartVersion                = -1 // used for EventStoreDB
	aggregateAppliedEventsInitialCap     = 10
	aggregateUncommittedEventsInitialCap = 10
)

type When interface {
	When(event Event) error
}

type when func(event Event) error

// Apply process Aggregate Event
type Apply interface {
	Apply(event Event) error
}

type Aggregate interface {
	When
	AggregateRoot
}

// AggregateRoot is an interface for models that persist state changes
type AggregateRoot interface {
	AddDomainEvent(Event)
	GetDomainEvents() []Event
	Apply
}

// AggregateBase is a struct that provides a base implementation of the AggregateRoot interface
type AggregateBase struct {
	events []Event
	when   when
}

// NewAggregateBase AggregateBase constructor, contains all main fields and methods,
// main aggregate must realize When interface and pass as argument to constructor
func NewAggregateBase(when when) *AggregateBase {
	if when == nil {
		return nil
	}

	return &AggregateBase{
		events: make([]Event, 0, aggregateAppliedEventsInitialCap),
		when:   when,
	}
}

// AddDomainEvent adds a domain event to the aggregate root
func (a *AggregateBase) AddDomainEvent(event Event) {
	a.events = append(a.events, event)
}

// GetDomainEvents returns all domain events collected by the aggregate root
func (a *AggregateBase) GetDomainEvents() []Event {
	return a.events
}

// Apply push event to aggregate uncommitted events using When method
func (a *AggregateBase) Apply(event Event) error {
	if err := a.when(event); err != nil {
		return err
	}

	return nil
}
