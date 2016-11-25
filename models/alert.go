package models

import "fmt"

type EventAlert interface{}

type PodEventAlert struct {
	Kind          string `json:"kind"`
	Name          string `json:"name"`
	Namespace     string `json:"namespace"`
	Host          string `json:"host"`
	Reason        string `json:"reason"`
	LastTimestamp string `json:"last_timestamp"`
	Message       string `json:"message"`
}

func (pea PodEventAlert) String() string {
	if pea.Host == "" {
		formatter := `Kubernetes POD Event Alert
"Kind": %s
"Name": %s
"Namespace": %s
"Reason": %s
"LastTimestamp": %s
"Message": %s
	`

		return fmt.Sprintf(formatter, pea.Kind, pea.Name, pea.Namespace,
			pea.Reason, pea.LastTimestamp, pea.Message,
		)
	}

	formatter := `Kubernetes POD Event Alert
"Kind": %s
"Name": %s
"Namespace": %s
"Host": %s
"Reason": %s
"LastTimestamp": %s
"Message": %s
	`

	return fmt.Sprintf(formatter, pea.Kind, pea.Name, pea.Namespace,
		pea.Host, pea.Reason, pea.LastTimestamp, pea.Message,
	)
}

type NodeEventAlert struct {
	Kind          string `json:"kind"`
	Name          string `json:"name"`
	Namespace     string `json:"namespace"`
	Reason        string `json:"reason"`
	LastTimestamp string `json:"last_timestamp"`
	Message       string `json:"message"`
}

func (nea NodeEventAlert) String() string {
	formatter := `Kubernetes NODE Event Alert
"Kind": %s
"Name": %s
"Namespace": %s
"Reason": %s
"LastTimestamp": %s
"Message": %s
	`

	return fmt.Sprintf(formatter, nea.Kind, nea.Name,
		nea.Namespace, nea.Reason, nea.LastTimestamp, nea.Message,
	)
}
