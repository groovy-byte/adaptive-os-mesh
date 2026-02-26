package controller

import (
	"testing"
	"time"
	"github.com/nats-io/nats.go"
)

func TestKVCacheSync(t *testing.T) {
	// 1. Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Skip("NATS not running, skipping TDD test")
		return
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		t.Fatal(err)
	}

	// 2. Initialize KV Sync
	k := NewKVCacheController(js)
	
	agentID := "test-agent"
	delta := []byte("The conversation context is building...")

	// 3. Broadcast Delta
	err = k.BroadcastDelta(agentID, delta)
	if err != nil {
		t.Fatal(err)
	}

	// 4. Verification: Subscribe to the mesh channel
	sub, err := js.SubscribeSync("mesh.kv_cache.>", nats.DeliverLast())
	if err != nil {
		t.Fatal(err)
	}
	
	msg, err := sub.NextMsg(1 * time.Second)
	if err != nil {
		t.Fatal("Did not receive KV cache delta broadcast")
	}
	
	if string(msg.Data) != string(delta) {
		t.Errorf("Received data mismatch: %s", msg.Data)
	}
}
