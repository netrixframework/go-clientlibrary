package replicaclient

import (
	"time"

	"github.com/netrixframework/netrix/types"
)

var (
	MessageSendEventType    = "MessageSend"
	MessageReceiveEventType = "MessageReceive"
	TimeoutStartEventType   = "TimeoutStart"
	TimeoutEndEventType     = "TimeoutEnd"
)

type event struct {
	Type      string            `json:"type"`
	Replica   types.ReplicaID   `json:"replica"`
	Params    map[string]string `json:"params"`
	Timestamp int64             `json:"timestamp"`
}

func (c *ReplicaClient) PublishEvent(t string, params map[string]string) {
	c.sendMasterMessage(&masterRequest{
		Type: "Event",
		Event: &event{
			Type:      t,
			Replica:   types.ReplicaID(c.config.ReplicaID),
			Params:    params,
			Timestamp: time.Now().Unix(),
		},
	})
}
func (c *ReplicaClient) PublishEventAsync(t string, params map[string]string) {
	go c.PublishEvent(t, params)
}
