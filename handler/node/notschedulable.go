package node

import (
	"github.com/andyxning/eventarbiter/cmd/eventarbiter/conf"
	"github.com/andyxning/eventarbiter/handler"
	"github.com/andyxning/eventarbiter/models"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/kubelet/events"
	"strings"
)

const (
	NodeNotSchedulableReason = events.NodeNotSchedulable
)

type notSchedulable struct {
	kind             string
	reason           string
	alertEventReason string
}

func NewNotSchedulable() models.EventHandler {
	return notReady{
		kind:             "NODE",
		reason:           NodeNotSchedulableReason,
		alertEventReason: "node_notschedulable",
	}
}

func (ns notSchedulable) HandleEvent(sinks []models.Sink, event *api.Event) {
	if strings.ToUpper(event.InvolvedObject.Kind) == ns.kind && event.Reason == ns.reason {
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
			sink.Sink(ns.kind, eventAlert)
		}
	}
}

func (ns notSchedulable) AlertEventReason() string {
	return ns.alertEventReason
}

func (ns notSchedulable) Reason() string {
	return ns.reason
}

func init() {
	ns := NewNotSchedulable()
	handler.MustRegisterEventAlertReason(ns.AlertEventReason(), ns)
	handler.RegisterEventReason(ns.Reason())
}
