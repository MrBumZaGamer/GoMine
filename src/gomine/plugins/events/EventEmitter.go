package events

import (
	"reflect"
	"gomine/interfaces"
	"errors"
	"gomine/plugins"
)

var NoFunctionErr = errors.New("the function passed as argument is not a function")
var InvalidCountErr = errors.New("the function passed as argument has an invalid parameter count. Should be none or one argument")
var InvalidEventErr = errors.New("the first argument of the function passed must be a valid event")
var InvalidPriorityErr = errors.New("event listener priorities should be in a range of 0-10")

type EventEmitter struct {
	listeners map[string][]*PriorityLevel
}

func NewEventEmitter() *EventEmitter {
	return &EventEmitter{make(map[string][]*PriorityLevel)}
}

/**
 * Registers a listener function for the given plugin and with the given priority.
 */
func (emitter *EventEmitter) RegisterListener(plugin plugins.IPlugin, priority byte, function interface{}) error {
	if priority > 10 || priority < 0 {
		return InvalidPriorityErr
	}

	var funcType = reflect.TypeOf(function)
	if funcType.Kind() != reflect.Func {
		return NoFunctionErr
	}

	var method = reflect.ValueOf(function)
	if method.Type().NumIn() != 1 {
		return InvalidCountErr
	}

	var firstArgument = method.Type().In(0)

	var eventInterface = reflect.TypeOf((*interfaces.IEvent)(nil)).Elem()
	if !firstArgument.Implements(eventInterface) {
		return InvalidEventErr
	}

	var name = firstArgument.Name()
	if _, ok := emitter.listeners[name]; !ok {
		emitter.listeners[name] = make([]*PriorityLevel, 11)
	}
	emitter.listeners[name][priority].addListener(plugin, method)

	return nil
}

/**
 * Emits an event.
 */
func (emitter *EventEmitter) EmitEvent(event interfaces.IEvent) {
	emitter.callListeners(event)
}

/**
 * Calls all listeners of the given event.
 */
func (emitter *EventEmitter) callListeners(event interfaces.IEvent) {
	for _, level := range emitter.listeners[event.GetEventName()] {
		level.callListeners(event)
	}
}