package pod

import (
	"github.com/andyxning/eventarbiter/cmd/eventarbiter/conf"
	"github.com/andyxning/eventarbiter/handler"
	"github.com/andyxning/eventarbiter/models"
	"k8s.io/kubernetes/pkg/api"
	"strings"
)

const (
	PodInsufficientCPUReason = "FailedScheduling"
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
		reason: PodInsufficientCPUReason,
		// TODO(andyxning): This should be replaced with more scalable reference instead of hard
		// code.
		keyWord:          "insufficient cpu",
		alertEventReason: "pod_insufficentcpu",
	}
}

func (ic insufficientCPU) HandleEvent(sinks []models.Sink, event *api.Event) {
	if strings.ToUpper(event.InvolvedObject.Kind) == ic.kind && event.Reason == ic.reason &&
		strings.Contains(strings.ToLower(event.Message), ic.keyWord) {
		// Currently in Kubernetes 1.4.0, message for insufficient cpu is proportional with
		// the number of minion machines. This may be very large and longer information is useless.
		// So, we can just truncate it to maxMessageLength.
		// See https://gist.github.com/andyxning/8065dc35889f07073e129bb75a6e57fe for an example.
		if len(event.Message) >= maxMessageLength {
			event.Message = event.Message[:maxMessageLength]
		}
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
