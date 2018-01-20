package server

import "gomine/plugins/events"

type LogEvent struct {
	*events.Event
	LogText string
}

func NewLogEvent(text string) *LogEvent {
	return &LogEvent{events.NewEvent("server.LogEvent"), text}
}