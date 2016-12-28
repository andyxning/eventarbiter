package pod

import (
	"github.com/andyxning/eventarbiter/cmd/eventarbiter/conf"
	"github.com/andyxning/eventarbiter/handler"
	"github.com/andyxning/eventarbiter/models"
	"k8s.io/kubernetes/pkg/api"
	"strings"
)

const (
	originalPodInsufficientCPUMemoryReason = "FailedScheduling"
	PodInsufficientCPUMemoryReason         = "InsufficientCPUMemory"
)

type insufficientCPUMemory struct {
	kind             string
	reason           string
	originalReason   string
	keyWords         []string
	alertEventReason string
}

func NewInsufficientCPUMemory() models.EventHandler {
	return insufficientCPUMemory{
		kind:           "POD",
		reason:         PodInsufficientCPUMemoryReason,
		originalReason: originalPodInsufficientCPUMemoryReason,
		// TODO(andyxning): This should be replaced with more scalable reference instead of hard
		// code.
		keyWords:         []string{"insufficient cpu", "insufficient memory"},
		alertEventReason: "pod_insufficentcpumemory",
	}
}

func (icm insufficientCPUMemory) HandleEvent(sinks []models.Sink, event *api.Event) {
	if strings.ToUpper(event.InvolvedObject.Kind) == icm.kind && event.Reason == icm.reason {
		for _, keyWord := range icm.keyWords {
			if strings.Contains(strings.ToLower(event.Message), keyWord) {
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
					Reason:        icm.reason,
					Message:       event.Message,
					LastTimestamp: event.LastTimestamp.Local().String(),
					Environment:   conf.Conf.Environment.Value,
				}

				for _, sink := range sinks {
					sink.Sink(icm.kind, eventAlert)
				}

				break
			}
		}
	}
}

func (icm insufficientCPUMemory) AlertEventReason() string {
	return icm.alertEventReason
}

func (icm insufficientCPUMemory) Reason() string {
	return icm.reason
}

func init() {
	icm := NewInsufficientCPUMemory()
	handler.MustRegisterEventAlertReason(icm.AlertEventReason(), icm)
	handler.RegisterEventReason(icm.Reason())
}
