package main

import (
	"encoding/json"
	"fmt"
)

type StartMessage struct {
	Memory int `json:"memory"` // GB
	NumCPU int `json:"num_cpu"`
}

func (StartMessage) Kind() string {
	return "start"
}

type StopMessage struct {
	ID string `json:"id"`
}

func (StopMessage) Kind() string {
	return "stop"
}

type Message struct {
	Kind    string
	Payload json.RawMessage
}

type SubMessage interface {
	StartMessage | StopMessage

	Kind() string
}

func GetSub[T SubMessage](m Message) (T, error) {
	var sm T
	if m.Kind != sm.Kind() {
		var zero T
		return zero, fmt.Errorf("expected kind %q, got %q", sm.Kind(), m.Kind)
	}

	if err := json.Unmarshal(m.Payload, &sm); err != nil {
		var zero T
		return zero, fmt.Errorf("unmarshal: %w", err)
	}

	return sm, nil
}

var stream = `
{"kind": "start", "payload": {"memory": 4, "num_cpu": 8}}`
`


func main() {
	data := []byte(

}
