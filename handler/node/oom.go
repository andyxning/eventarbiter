package node

import (
	"github.com/andyxning/eventarbiter/cmd/eventarbiter/conf"
	"github.com/andyxning/eventarbiter/handler"
	"github.com/andyxning/eventarbiter/models"
	"k8s.io/kubernetes/pkg/api"
	"strings"
)

const (
	NodeOOMReason = "SystemOOM"
)

type oom struct {
	kind             string
	reason           string
	alertEventReason string
}

func NewOOM() models.EventHandler {
	return oom{
		kind:             "NODE",
		reason:           NodeOOMReason,
		alertEventReason: "node_systemoom",
	}
}

func (om oom) HandleEvent(sinks []models.Sink, event *api.Event) {
	if strings.ToUpper(event.InvolvedObject.Kind) == om.kind && event.Reason == om.reason {
		var eventAlert = models.NodeEventAlert{
			Kind:          event.InvolvedObject.Kind,
			Name:          event.InvolvedObject.Name,
			Namespace:     event.ObjectMeta.Namespace,
			Reason:        event.Reason,
			Message:       event.Message,
			LastTimestamp: event.LastTimestamp.Local().String(),
			Environment:   conf.Conf.Environment.Value,
		}

		for _, sink := range sinks {
			sink.Sink(om.kind, eventAlert)
		}
	}
}

func (om oom) AlertEventReason() string {
	return om.alertEventReason
}

func (om oom) Reason() string {
	return om.reason
}

func init() {
	om := NewOOM()
	handler.MustRegisterEventAlertReason(om.AlertEventReason(), om)
	handler.RegisterEventReason(om.Reason())
}
