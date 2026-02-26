package main

import (
	"encoding/json"
	"testing"
)

func TestWaybarOutputFormat(t *testing.T) {
	// Waybar requires a specific JSON format
	text := "[A: 125 GB/s | B: 154 GB/s]"
	tooltip := "Vextra Mesh Status: Synchronized"
	class := "active"
	
	output := WaybarOutput{
		Text:    text,
		Tooltip: tooltip,
		Class:   class,
	}
	
	data, err := json.Marshal(output)
	if err != nil {
		t.Fatal(err)
	}
	
	var decoded WaybarOutput
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatal(err)
	}
	
	if decoded.Text != text {
		t.Errorf("Expected text %s, got %s", text, decoded.Text)
	}
}
