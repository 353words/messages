package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
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

func MessageKind[T SubMessage]() string {
	var m T

	return m.Kind()
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

func ConsumeMessages(r io.Reader) error {
	dec := json.NewDecoder(r)

	for {
		var m Message
		err := dec.Decode(&m)
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return err
		}

		switch m.Kind {
		case MessageKind[StartMessage]():
			sm, err := GetSub[StartMessage](m)
			if err != nil {
				return err
			}
			fmt.Println("handling start:", sm)
		case MessageKind[StopMessage]():
			sm, err := GetSub[StopMessage](m)
			if err != nil {
				return err
			}
			fmt.Println("handling stop:", sm)
		default:
			return fmt.Errorf("%q: unknown message kind", m.Kind)
		}
	}

	return nil
}

func main() {
	data := `
		{"kind": "start", "payload": {"memory": 4, "num_cpu": 8}}
		{"kind": "stop",  "payload": {"id": "6870b39"}}
		{"kind": "start", "payload": {"memory": 32, "num_cpu": 4}}
	`
	r := strings.NewReader(data)
	if err := ConsumeMessages(r); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

}
