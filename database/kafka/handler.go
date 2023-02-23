// Package kafka
//
// @author: xwc1125
// @date: 2021/3/24
package kafka

import (
	"fmt"
	"log"

	"github.com/Shopify/sarama"
)

// ConsumerHandler represents a Sarama consumer group consumer
type ConsumerHandler struct {
	IsLog    bool
	Ready    chan bool
	Result   chan *Result
	Feedback <-chan bool // 成功还是失败
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(h.Ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for msg := range claim.Messages() {
		if h.IsLog {
			log.Println(fmt.Sprintf("Message claimed: value = %s, timestamp = %v, topic = %s", string(msg.Value), msg.Timestamp, msg.Topic))
		}
		if h.Result != nil {
			h.Result <- &Result{
				Data: Message{
					Topic: msg.Topic,
					Key:   msg.Key,
					Value: msg.Value,
					// Headers:   msg.Headers,
					Offset:    msg.Offset,
					Partition: msg.Partition,
					Timestamp: msg.Timestamp,
				},
			}
		}
		if h.Feedback != nil {
			select {
			case f := <-h.Feedback:
				if f {
					session.MarkMessage(msg, "")
				}
			}
		} else {
			session.MarkMessage(msg, "")
		}
	}

	return nil
}
