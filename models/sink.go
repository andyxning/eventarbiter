package models

type Sink interface {
	// Sink sends alter to specific destination.
	Sink(kind string, eventAlert EventAlert)
	// Stop stops sinking gracefully.
	Stop()
	// Name returns the sink name.
	Name() string
}
