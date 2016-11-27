package pod

import (
	"github.com/andyxning/eventarbiter/cmd/eventarbiter/conf"
	"github.com/andyxning/eventarbiter/handler"
	"github.com/andyxning/eventarbiter/models"
	"k8s.io/kubernetes/pkg/api"
	"strings"
)

const (
	PodInsufficientMemoryReason = "FailedScheduling"
)

type insufficientMemory struct {
	kind             string
	reason           string
	keyWord          string
	alertEventReason string
}

func NewInsufficientMemory() models.EventHandler {
	return insufficientMemory{
		kind:   "POD",
		reason: PodInsufficientMemoryReason,
		// TODO(andyxning): This should be replaced with more scalable reference instead of hard
		// code.
		keyWord:          "Insufficient Memory",
		alertEventReason: "pod_insufficientmemory",
	}
}

func (im insufficientMemory) HandleEvent(sinks []models.Sink, event *api.Event) {
	if strings.ToUpper(event.InvolvedObject.Kind) == im.kind && event.Reason == im.reason &&
		strings.Contains(event.Message, im.keyWord) {
		var eventAlert = models.PodEventAlert{
			Kind:          event.InvolvedObject.Kind,
			Name:          event.InvolvedObject.Name,
			Namespace:     event.ObjectMeta.Namespace,
			Reason:        event.Reason,
			Message:       event.Message,
			LastTimestamp: event.LastTimestamp.Local().String(),
			Environment:   conf.Conf.Environment.Value,
		}

		for _, sink := range sinks {
			sink.Sink(im.kind, eventAlert)
		}
	}
}

func (im insufficientMemory) AlertEventReason() string {
	return im.alertEventReason
}

func (im insufficientMemory) Reason() string {
	return im.reason
}

func init() {
	im := NewInsufficientMemory()
	handler.MustRegisterEventAlertReason(im.AlertEventReason(), im)
	handler.RegisterEventReason(im.Reason())
}
