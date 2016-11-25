package models

import "k8s.io/kubernetes/pkg/api"

type EventHandler interface {
	// HandleEvent sinks event.
	HandleEvent([]Sink, *api.Event)
	// AlertEventReason returns a unique string representing an alert event reason.
	AlertEventReason() string
	// Reason returns a string describing the event reason in short.
	Reason() string
}
