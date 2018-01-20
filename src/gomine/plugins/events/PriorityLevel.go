package events

import (
	"reflect"
	"gomine/plugins"
	"gomine/interfaces"
)

const (
	PriorityFirst = 0
	PriorityLast = 10
)

type PriorityLevel struct {
	listeners map[string]reflect.Value
}

/**
 * Adds a listener of the given plugin to this priority level.
 */
func (level *PriorityLevel) addListener(plugin plugins.IPlugin, method reflect.Value) {
	level.listeners[plugin.GetName()] = method
}

/**
 * Removes a listener from the given plugin at this priority level.
 */
func (level *PriorityLevel) removeListener(plugin plugins.IPlugin) {
	delete(level.listeners, plugin.GetName())
}

/**
 * Calls all listeners with the given event.
 */
func (level *PriorityLevel) callListeners(event interfaces.IEvent) {
	var values = []reflect.Value{reflect.ValueOf(event)}
	for _, listener := range level.listeners {
		listener.Call(values)
	}
}