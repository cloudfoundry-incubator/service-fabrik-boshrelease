package main

import (
	"encoding/json"
	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type ServiceInfo struct {
	// The id mentioned is the SKU name of service
	// like redis, postgresql and not uutd
	ID   string `json:"id"`
	Plan string `json:"plan"`
}

type ConsumerInfo struct {
	Environment string `json:"environment"`
	Region      string `json:"region"`
	Org         string `json:"org"`
	Space       string `json:"space"`
	Instance    string `json:"instance"`
}

type InstancesMeasure struct {
	ID    string `json:"id"`
	Value int    `json:"value"`
}

// MeteringOptions represents the options field of Metering Resource
type MeteringOptions struct {
	ID                string             `json:"id"`
	Timestamp         string             `json:"timestamp"`
	ServiceInfo       ServiceInfo        `json:"service"`
	ConsumerInfo      ConsumerInfo       `json:"consumer"`
	InstancesMeasures []InstancesMeasure `json:"measues"`
}

// MeteringSpec represents the spec field of metering resource
type MeteringSpec struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Options           string `json:"options,omitempty"`
}

// Metering structure holds all the details related to
// Metering event
type Metering struct {
	Spec MeteringSpec `json:"spec"`
}

func (m *Metering) getName() string {
	var meteringOptions MeteringOptions
	json.Unmarshal([]byte(m.Spec.Options), &meteringOptions)
	return meteringOptions.ID
}

func newMetering(opt GenericOptions, crd GenericResource, signal int) *Metering {
	si := ServiceInfo{
		ID:   opt.ServiceID,
		Plan: opt.PlanID,
	}
	ci := ConsumerInfo{
		Environment: "",
		Region:      "",
		Org:         opt.Context.OrganizationGUID,
		Space:       opt.Context.SpaceGUID,
		Instance:    crd.Name,
	}
	im := InstancesMeasure{
		ID:    "instances",
		Value: signal,
	}
	mo := MeteringOptions{
		ID: uuid.New().String(),
		// Go has wierd time formating rules !!
		// https://golang.org/src/time/format.go
		Timestamp:         time.Now().UTC().Format(time.RFC3339Nano),
		ServiceInfo:       si,
		ConsumerInfo:      ci,
		InstancesMeasures: []InstancesMeasure{im},
	}
	meteringOptions, _ := json.Marshal(mo)
	m := &Metering{
		Spec: MeteringSpec{
			Options: string(meteringOptions),
		},
	}
	return m
}
