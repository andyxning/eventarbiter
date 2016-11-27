package pod

import (
	"github.com/andyxning/eventarbiter/cmd/eventarbiter/conf"
	"github.com/andyxning/eventarbiter/handler"
	"github.com/andyxning/eventarbiter/models"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/kubelet/events"
	"strings"
)

const (
	PodFailedReason = events.FailedToStartContainer
)

type failed struct {
	kind             string
	reason           string
	alertEventReason string
}

func NewFailed() models.EventHandler {
	return failed{
		kind:             "POD",
		reason:           PodFailedReason,
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
			Environment:   conf.Conf.Environment.Value,
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
