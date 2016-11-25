package models

import "k8s.io/kubernetes/pkg/api"

type Source interface {
	// Start starts to collect event and send it to the eventChan.
	Start(eventChan chan<- *api.Event)
	// Stop stops Source gracefully.
	Stop()
}
