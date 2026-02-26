package controller

import (
	"fmt"
	"log"
	"github.com/nats-io/nats.go"
)

// KVCacheController handles high-speed conversation state sync across the mesh
type KVCacheController struct {
	js nats.JetStreamContext
}

func NewKVCacheController(js nats.JetStreamContext) *KVCacheController {
	// Ensure the stream exists for KV sync
	_, err := js.AddStream(&nats.StreamConfig{
		Name:     "MESH_STATE",
		Subjects: []string{"mesh.kv_cache.>"},
		Storage:  nats.MemoryStorage, // High-speed, transient state
	})
	if err != nil {
		log.Printf("[KV] Warning: Could not create/verify MESH_STATE stream: %v", err)
	}

	return &KVCacheController{js: js}
}

// BroadcastDelta sends a state fragment to all nodes in the mesh
func (k *KVCacheController) BroadcastDelta(agentID string, delta []byte) error {
	subject := fmt.Sprintf("mesh.kv_cache.%s", agentID)
	_, err := k.js.Publish(subject, delta)
	if err != nil {
		return fmt.Errorf("failed to broadcast KV delta: %w", err)
	}
	
	log.Printf("[KV] 📡 Broadcasted %d bytes for agent %s", len(delta), agentID)
	return nil
}

// SubscribeToDeltas allows a node to listen for conversation state updates
func (k *KVCacheController) SubscribeToDeltas(agentID string, handler func([]byte)) (*nats.Subscription, error) {
	subject := fmt.Sprintf("mesh.kv_cache.%s", agentID)
	sub, err := k.js.Subscribe(subject, func(m *nats.Msg) {
		handler(m.Data)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to KV deltas: %w", err)
	}
	return sub, nil
}
