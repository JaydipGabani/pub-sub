package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
)

type PubsubMsg struct {
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
	PublishTime time.Time `json:"publishTime,omitempty"`
}

var sub = &common.Subscription{
	PubsubName: "pubsub",
	Topic:      "audit",
	Route:      "/checkout",
}

func main() {
	s := daprd.NewService(":6002")
	log.Printf("Listening 6...")
	startid = ""
	auditmap = sync.Map{}
	if err := s.AddTopicEventHandler(sub, eventHandler); err != nil {
		log.Fatalf("error adding topic subscription: %v", err)
	}
	if err := s.Start(); err != nil {
		log.Fatalf("error listening: %v", err)
	}
}

type Detail struct {
	published string
	received int
	time time.Duration
	st time.Time
	latency time.Duration
}

var startid string
var eventCount int
var startTime time.Time
var flag bool
var latency time.Duration
var auditmap sync.Map

func eventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	var msg PubsubMsg
	jsonInput, err := strconv.Unquote(string(e.RawData))
	if err != nil {
		log.Fatalf("error unquoting %v", err)
		return true, err
	}
	if err := json.Unmarshal([]byte(jsonInput), &msg); err != nil {
		log.Fatalf("error %v", err)
		return true, err
	}
	// if msg.Message != "audit_finished" {
	// 	tmp, _ := auditmap.LoadOrStore(msg.ID, Detail{received: 0, st: time.Now(), latency: 0})
	// 	tt := tmp.(Detail)
	// 	tt.published = msg.ResourceName
	// 	tt.received = tt.received + 1
	// 	tt.time = time.Since(tt.st)
	// 	tt.latency = tt.latency + time.Since(msg.PublishTime)

	// 	auditmap.Store(msg.ID, tt)

	// 	auditmap.Range(func(key, value interface{}) bool {
	// 		v := value.(Detail)
	// 		t, err := strconv.Atoi(v.published)
	// 		if err != nil {
	// 			log.Fatalf("error %v", err)
	// 		}
	// 		log.Println("number of violations published", t, "number of violation received", v.received, "audit id", key, "time", v.time, "average latency per msg", v.latency/time.Duration(v.received))
	// 		return true
	// 	})

	// 	// for k, v := range auditmap {
	// 	// 	t, err := strconv.Atoi(v.published)
	// 	// 	if err != nil {
	// 	// 		log.Fatalf("error %v", err)
	// 	// 	}
	// 	// 	log.Println("number of violations published", t, "number of violation received", v.received, "audit id", k, "time", v.time, "average latency per msg", v.latency/time.Duration(v.received))
	// 	// }
	// }
	
	if startid == "" {
		startid = msg.ID
		startTime = msg.PublishTime
		eventCount = 0
	} // else if msg.ID != startid {
	// 	startid = msg.ID
	// 	startTime = msg.PublishTime
	// 	log.Println("message from new audit", startid)
	// } 
	// if msg.Message == "audit_finished" {
		// if eventCount != 0{
		// 	log.Println("number of violations published", msg.ResourceName, "number of violation received", eventCount, "audit id", startid, "time", time.Since(startTime), "average latency per msg", latency/time.Duration(eventCount), "audit_finished")
		// } else {
		// 	log.Println("number of violations published", msg.ResourceName, "number of violation received", eventCount, "audit id", startid, "time", time.Since(startTime), "audit_finished")	
		// }
		// latency = 0
		// eventCount = 0
		// return false, nil
	// }
	eventCount++
	latency += time.Since(msg.PublishTime)
	
	log.Println("number of violations published", msg.ResourceName, "number of violation received", eventCount, "audit id", startid, "time", time.Since(startTime), "average latency per msg", latency/time.Duration(eventCount))
	return false, nil
}
