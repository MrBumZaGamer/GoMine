package events

type Event struct {
	Name string
}

func NewEvent(name string) *Event {
	return &Event{name}
}

/**
 * Returns the name of the event.
 */
func (event *Event) GetEventName() string {
	return event.Name
}