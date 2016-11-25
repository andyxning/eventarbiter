package handler

import (
	"github.com/andyxning/eventarbiter/handler"
	"github.com/andyxning/eventarbiter/models"
	"k8s.io/kubernetes/pkg/api"
	"strings"
)

type insufficientCPU struct {
	kind             string
	reason           string
	keyWord          string
	alertEventReason string
}

func NewInsufficientCPU() models.EventHandler {
	return insufficientCPU{
		kind:   "POD",
		reason: "FailedScheduling",
		// TODO(andyxning): This should be replaced with more scalable reference instead of hard
		// code.
		keyWord:          "Insufficient CPU",
		alertEventReason: "pod_insufficentcpu",
	}
}

func (ic insufficientCPU) HandleEvent(sinks []models.Sink, event *api.Event) {
	if strings.ToUpper(event.InvolvedObject.Kind) == ic.kind && event.Reason == ic.reason &&
		strings.Contains(event.Message, ic.keyWord) {
		var eventAlert = models.PodEventAlert{
			Kind:          event.InvolvedObject.Kind,
			Name:          event.InvolvedObject.Name,
			Namespace:     event.ObjectMeta.Namespace,
			Reason:        event.Reason,
			Message:       event.Message,
			LastTimestamp: event.LastTimestamp.Local().String(),
		}

		for _, sink := range sinks {
			sink.Sink(ic.kind, eventAlert)
		}
	}
}

func (ic insufficientCPU) AlertEventReason() string {
	return ic.alertEventReason
}

func (ic insufficientCPU) Reason() string {
	return ic.reason
}

func init() {
	ic := NewInsufficientCPU()
	handler.MustRegisterEventAlertReason(ic.AlertEventReason(), ic)
	handler.RegisterEventReason(ic.Reason())
}
