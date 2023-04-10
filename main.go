package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
)

type PubsubMsg struct {
	Key string `json:"key,omitempty"`
	ID                    string            `json:"id,omitempty"`
	Details               interface{}       `json:"details,omitempty"`
	EventType             string            `json:"eventType,omitempty"`
	Group                 string            `json:"group,omitempty"`
	Version               string            `json:"version,omitempty"`
	Kind                  string            `json:"kind,omitempty"`
	Name                  string            `json:"name,omitempty"`
	Namespace             string            `json:"namespace,omitempty"`
	Message               string            `json:"message,omitempty"`
	EnforcementAction     string            `json:"enforcementAction,omitempty"`
	ConstraintAnnotations map[string]string `json:"constraintAnnotations,omitempty"`
	ResourceGroup         string            `json:"resourceGroup,omitempty"`
	ResourceAPIVersion    string            `json:"resourceAPIVersion,omitempty"`
	ResourceKind          string            `json:"resourceKind,omitempty"`
	ResourceNamespace     string            `json:"resourceNamespace,omitempty"`
	ResourceName          string            `json:"resourceName,omitempty"`
	ResourceLabels        map[string]string `json:"resourceLabels,omitempty"`
}

var sub = &common.Subscription{
	PubsubName: "pubsub",
	Topic:      "audit",
	Route:      "/checkout",
}

func main() {
	s := daprd.NewService(":6002")
	log.Printf("Listening...")
	if err := s.AddTopicEventHandler(sub, eventHandler); err != nil {
		log.Fatalf("error adding topic subscription: %v", err)
	}
	flag = true

	if err := s.Start(); err != nil {
		log.Fatalf("error listenning: %v", err)
	}

	for {
		if time.Since(startTime) > time.Duration(20 * time.Second) {
			endTime = time.Now() // record end time
			duration := endTime.Sub(startTime)
			log.Printf("Total events received: %d, Time taken: %v", eventCount, duration)
			panic("test")
		} 
		time.Sleep(20 * time.Second)
	}
}

var eventCount int
var startTime time.Time
var endTime time.Time
var flag bool

func eventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	var msg PubsubMsg
	jsonInput, err := strconv.Unquote(string(e.RawData))
	if err != nil {
		log.Fatalf("error unquoting %v", err)
	}
	if err := json.Unmarshal([]byte(jsonInput), &msg); err != nil {
		log.Fatalf("error %v", err)
	}
	if flag {
		startTime = time.Now() // record start time
		flag = false
	} 
	
	eventCount++
	log.Printf("%#v", eventCount)
	
	return false, nil
}