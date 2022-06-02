package replicaclient

import (
	"sync"
	"time"

	"github.com/netrixframework/netrix/types"
)

// TimeoutInfo encapsulates the timeout information that needs to be scheduled
type TimeoutInfo interface {
	// Key returns a unique key for the given timeout. Only one timeout for a specific key can be running at any given time.
	Key() string
	// Duration of the timeout
	Duration() time.Duration
}

type timer struct {
	timeouts map[string]TimeoutInfo
	outChan  chan TimeoutInfo
	lock     *sync.Mutex
}

type timeout struct {
	Type     string          `json:"type"`
	Duration string          `json:"duration"`
	Replica  types.ReplicaID `json:"replica"`
}

func newTimer() *timer {
	return &timer{
		timeouts: make(map[string]TimeoutInfo),
		outChan:  make(chan TimeoutInfo, 10),
		lock:     new(sync.Mutex),
	}
}

func (t *timer) AddTimeout(info TimeoutInfo) bool {
	t.lock.Lock()
	defer t.lock.Unlock()
	key := info.Key()
	_, ok := t.timeouts[key]
	if !ok {
		t.timeouts[key] = info
		return true
	}
	return false
}

func (t *timer) FireTimeout(key string) {
	t.lock.Lock()
	info, ok := t.timeouts[key]
	t.lock.Unlock()
	if ok {
		t.lock.Lock()
		delete(t.timeouts, key)
		t.lock.Unlock()
		go func(i TimeoutInfo) { t.outChan <- i }(info)
	}
}

// StartTimer schedules a timer for the given TimerInfo
// Note: The timers are implemented as message sends and receives that are to be scheduled by the
// testing strategy. If you do not want to instrument timers as message send/receives then do not use this function.
func (c *ReplicaClient) StartTimer(i TimeoutInfo) {
	c.logger.Info("Starting timer", "type", i.Key(), "duration", i.Duration().String())
	if c.timer.AddTimeout(i) {
		duration := i.Duration()
		if duration < 0 {
			duration = time.Duration(0)
		}
		c.PublishEvent(TimeoutStartEventType, map[string]string{
			"type":     i.Key(),
			"duration": i.Duration().String(),
		})
	}
}

// TimeoutChan returns the channel on which timeouts are delivered.
func (c *ReplicaClient) TimeoutChan() chan TimeoutInfo {
	return c.timer.outChan
}
