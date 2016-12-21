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
		keyWord:          "insufficient memory",
		alertEventReason: "pod_insufficientmemory",
	}
}

func (im insufficientMemory) HandleEvent(sinks []models.Sink, event *api.Event) {
	// Currently in Kubernetes 1.4.0, message for insufficient cpu is proportional with
	// the number of minion machines. This may be very large and longer information is useless.
	// So, we can just truncate it to maxMessageLength.
	// See https://gist.github.com/andyxning/8065dc35889f07073e129bb75a6e57fe for an example.
	if len(event.Message) >= maxMessageLength {
		event.Message = event.Message[:maxMessageLength]
	}

	if strings.ToUpper(event.InvolvedObject.Kind) == im.kind && event.Reason == im.reason &&
		strings.Contains(strings.ToLower(event.Message), im.keyWord) {
		var eventAlert = models.PodEventAlert{
			Kind:          strings.ToUpper(event.InvolvedObject.Kind),
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
