package handler

import (
	"github.com/andyxning/eventarbiter/handler"
	"github.com/andyxning/eventarbiter/models"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/kubelet/events"
	"strings"
)

type failed struct {
	kind             string
	reason           string
	alertEventReason string
}

func NewFailed() models.EventHandler {
	return failed{
		kind:             "POD",
		reason:           events.FailedToStartContainer,
		alertEventReason: "pod_failed",
	}
}

func (fd failed) HandleEvent(sinks []models.Sink, event *api.Event) {
	if strings.ToUpper(event.InvolvedObject.Kind) == fd.kind && event.Reason == fd.reason {
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
			sink.Sink(fd.kind, eventAlert)
		}
	}
}

func (fd failed) AlertEventReason() string {
	return fd.alertEventReason
}

func (fd failed) Reason() string {
	return fd.reason
}

func init() {
	fd := NewFailed()
	handler.MustRegisterEventAlertReason(fd.AlertEventReason(), fd)
	handler.RegisterEventReason(fd.Reason())
}
