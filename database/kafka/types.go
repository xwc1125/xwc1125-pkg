// Package kafka
//
// @author: xwc1125
// @date: 2021/3/19
package kafka

import (
	"encoding/json"
	"fmt"
	"time"
)

type Result struct {
	Error error
	Data  Message
}

type Message struct {
	Topic string
	Key   []byte
	Value []byte
	// Headers   []sarama.RecordHeader
	Metadata  interface{}
	Offset    int64
	Partition int32
	Timestamp time.Time
}

type jsonMessage struct {
	Topic string
	Key   string
	Value string
	// Headers   []sarama.RecordHeader
	Metadata  interface{}
	Offset    int64
	Partition int32
	Timestamp time.Time
}

func (m Message) toJsonMessage() jsonMessage {
	return jsonMessage{
		Topic:     m.Topic,
		Key:       string(m.Key),
		Value:     string(m.Value),
		Metadata:  m.Metadata,
		Offset:    m.Offset,
		Partition: m.Partition,
		Timestamp: m.Timestamp,
	}
}

func (m *Message) String() string {
	bytes, _ := json.Marshal(m.toJsonMessage())
	return string(bytes)
}

func (m *Message) TerminalString() string {
	bytes, _ := json.Marshal(m.toJsonMessage())
	return string(bytes)
}

func (m *Message) Format(s fmt.State, c rune) {
	bytes, _ := json.Marshal(m.toJsonMessage())
	fmt.Fprintf(s, "%"+string(c), bytes[:])
}

func (m *Message) MarshalText() ([]byte, error) {
	return m.MarshalJSON()
}

func (m *Message) UnmarshalText(input []byte) error {
	return m.UnmarshalJSON(input)
}

func (m *Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.toJsonMessage())
}

func (m *Message) UnmarshalJSON(input []byte) error {
	jsonMessage := new(jsonMessage)
	err := json.Unmarshal(input, &jsonMessage)
	if err != nil {
		return err
	}
	m.Topic = jsonMessage.Topic
	m.Key = []byte(jsonMessage.Key)
	m.Value = []byte(jsonMessage.Value)
	m.Metadata = jsonMessage.Metadata
	m.Offset = jsonMessage.Offset
	m.Partition = jsonMessage.Partition
	m.Timestamp = jsonMessage.Timestamp
	return nil
}
