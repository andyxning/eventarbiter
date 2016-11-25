package handler

import (
	"github.com/andyxning/eventarbiter/handler"
	"github.com/andyxning/eventarbiter/models"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/kubelet/events"
	"strings"
)

type backOff struct {
	kind             string
	reason           string
	alertEventReason string
}

func NewBackOff() models.EventHandler {
	return backOff{
		kind:             "POD",
		reason:           events.BackOffStartContainer,
		alertEventReason: "pod_backoff",
	}
}

func (bf backOff) HandleEvent(sinks []models.Sink, event *api.Event) {
	if strings.ToUpper(event.InvolvedObject.Kind) == bf.kind && event.Reason == bf.reason {
		var eventAlert = models.PodEventAlert{
			Kind:          event.InvolvedObject.Kind,
			Name:          event.InvolvedObject.Name,
			Namespace:     event.ObjectMeta.Namespace,
			Host:          event.Source.Host,
			Reason:        event.Reason,
			Message:       event.Message,
			LastTimestamp: event.LastTimestamp.Local().String(),
		}

		for _, sink := range sinks {
			sink.Sink(bf.kind, eventAlert)
		}
	}
}

func (bf backOff) AlertEventReason() string {
	return bf.alertEventReason
}

func (bf backOff) Reason() string {
	return bf.reason
}

func init() {
	bf := NewBackOff()
	handler.MustRegisterEventAlertReason(bf.AlertEventReason(), bf)
	handler.RegisterEventReason(bf.Reason())
}
