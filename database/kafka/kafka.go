// Package kafka
//
// @author: xwc1125
// @date: 2021/3/19
package kafka

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/xwc1125/xwc1125-pkg/utils/stringutil"
)

// Kafka ...
type Kafka struct {
	config  *KafkaConfig
	version sarama.KafkaVersion
}

// Producer ...
type Producer struct {
	isAsync       bool // 是否为异步
	asyncProducer sarama.AsyncProducer
	syncProducer  sarama.SyncProducer
}

// Consumer ...
type Consumer struct {
	isLog         bool
	consumer      sarama.Consumer
	consumerGroup sarama.ConsumerGroup
	pcsLock       sync.RWMutex
	pcs           map[string]sarama.PartitionConsumer
}

// New ...
func New(config *KafkaConfig) (*Kafka, error) {
	if config == nil {
		return nil, errors.New("config is empty")
	}
	if config.IsLog {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}
	version := sarama.V0_11_0_2
	var err error
	if !stringutil.IsEmpty(config.KafkaVersion) {
		version, err = sarama.ParseKafkaVersion(config.KafkaVersion)
		if err != nil {
			log.Println(fmt.Sprintf("Error parsing Kafka version: %v", err))
			return nil, err
		}
	}
	return &Kafka{
		config:  config,
		version: version,
	}, nil
}

// NewProducer ...
func (k *Kafka) NewProducer() (*Producer, error) {
	config := DefaultConfig()
	if k.config.SASLEnable {
		config.Net.SASL.Enable = k.config.SASLEnable
		config.Net.SASL.User = k.config.SASLUser
		config.Net.SASL.Password = k.config.SASLPassword
	}
	return k.newProducer(k.config.IsAsync, config)
}

// NewProducerWithConfig ...
func (k *Kafka) NewProducerWithConfig(kafkaConfig *sarama.Config) (*Producer, error) {
	return k.newProducer(k.config.IsAsync, kafkaConfig)
}

// newProducer 创建生产者
func (k *Kafka) newProducer(isAsync bool, config *sarama.Config) (*Producer, error) {
	if isAsync {
		asyncProducer, err := sarama.NewAsyncProducer(k.config.Addrs, config)
		if err != nil {
			return nil, err
		}
		return &Producer{
			isAsync:       isAsync,
			asyncProducer: asyncProducer,
		}, nil
	} else {
		syncProducer, err := sarama.NewSyncProducer(k.config.Addrs, config)
		if err != nil {
			return nil, err
		}
		return &Producer{
			isAsync:      isAsync,
			syncProducer: syncProducer,
		}, nil
	}
}

// DefaultConfig ...
func DefaultConfig() *sarama.Config {
	config := sarama.NewConfig()
	// 等待服务器所有副本都保存成功后的响应
	// config.Producer.RequiredAcks = sarama.NoResponse
	config.Producer.RequiredAcks = sarama.WaitForLocal
	// 随机的分区类型：返回一个分区器，该分区器每次选择一个随机分区
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 是否等待成功和失败后的响应
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = sarama.V0_11_0_2
	config.Consumer.Offsets.Retry.Max = 3
	config.Consumer.Offsets.CommitInterval = 1 * time.Second
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1000 * time.Millisecond
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Group.Rebalance.Timeout = 10 * time.Second
	return config
}

// Close ...
func (p *Producer) Close() error {
	if p.asyncProducer != nil {
		p.asyncProducer.AsyncClose()
	}
	if p.syncProducer != nil {
		return p.syncProducer.Close()
	}
	return nil
}

// ProducerMessage ...
func (p *Producer) ProducerMessage(topic, key, value string, partition int32, result chan *Result) error {
	// 构建发送的消息，
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: partition, // 分区
		Key:       sarama.StringEncoder(key),
		Value:     sarama.ByteEncoder(value),
	}
	if p.asyncProducer != nil {
		p.asyncProducer.Input() <- msg

	} else {
		if p.syncProducer != nil {
			partition, offset, err := p.syncProducer.SendMessage(msg)
			if err != nil {
				return err
			}
			if result != nil {
				r := &Result{
					Data: Message{
						Topic: topic,
						Key:   convertData(msg.Key),
						Value: convertData(msg.Value),
						// Headers:   msg.Headers,
						Metadata:  msg.Metadata,
						Offset:    offset,
						Partition: partition,
						Timestamp: msg.Timestamp,
					},
				}
				result <- r
			}
		}
	}
	return nil
}

func (p *Producer) listen(result chan *Result) {
	if p.asyncProducer != nil {
		select {
		case suc := <-p.asyncProducer.Successes():
			if result != nil {
				result <- &Result{
					Data: Message{
						Topic: suc.Topic,
						Key:   convertData(suc.Key),
						Value: convertData(suc.Value),
						// Headers:   suc.Headers,
						Metadata:  suc.Metadata,
						Offset:    suc.Offset,
						Partition: suc.Partition,
						Timestamp: suc.Timestamp,
					},
				}
			}
		case fail := <-p.asyncProducer.Errors():
			if result != nil {
				result <- &Result{
					Error: fail.Err,
					Data: Message{
						Topic: fail.Msg.Topic,
						Key:   convertData(fail.Msg.Key),
						Value: convertData(fail.Msg.Value),
						// Headers:   fail.Msg.Headers,
						Metadata:  fail.Msg.Metadata,
						Offset:    fail.Msg.Offset,
						Partition: fail.Msg.Partition,
						Timestamp: fail.Msg.Timestamp,
					},
				}
			}
		}
	}
}

func convertData(encoder sarama.Encoder) []byte {
	bytes, _ := encoder.Encode()
	return bytes
}

// NewConsumer ...
func (k *Kafka) NewConsumer() (*Consumer, error) {
	config := DefaultConfig()
	if k.config.SASLEnable {
		config.Net.SASL.Enable = k.config.SASLEnable
		config.Net.SASL.User = k.config.SASLUser
		config.Net.SASL.Password = k.config.SASLPassword
	}
	return k.NewConsumerWithConfig(config)
}

// NewConsumerWithConfig ...
func (k *Kafka) NewConsumerWithConfig(kafkaConfig *sarama.Config) (*Consumer, error) {
	consumer, err := sarama.NewConsumer(k.config.Addrs, kafkaConfig)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		isLog:    k.config.IsLog,
		consumer: consumer,
		pcs:      make(map[string]sarama.PartitionConsumer),
	}, nil
}

// NewConsumerWithGroup ...
func (k *Kafka) NewConsumerGroup() (*Consumer, error) {
	config := DefaultConfig()
	if k.config.SASLEnable {
		config.Net.SASL.Enable = k.config.SASLEnable
		config.Net.SASL.User = k.config.SASLUser
		config.Net.SASL.Password = k.config.SASLPassword
	}
	k.setStrategy(config)
	return k.NewConsumerGroupConfig(config)
}

func (k *Kafka) NewConsumerGroupConfig(config *sarama.Config) (*Consumer, error) {
	consumer, err := sarama.NewConsumerGroup(k.config.Addrs, k.config.GroupId, config)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		isLog:         k.config.IsLog,
		consumerGroup: consumer,
		pcs:           make(map[string]sarama.PartitionConsumer),
	}, nil
}

func (k *Kafka) setStrategy(config *sarama.Config) {
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	switch k.config.Strategy {
	case "sticky":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case "roundrobin":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case "range":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	}
}

// ConsumerMessage ...
func (c *Consumer) ConsumerMessage(topic string, partition int32, result chan *Result) error {
	// ConsumePartition方法根据主题，分区和给定的偏移量创建创建了相应的分区消费者
	// 如果该分区消费者已经消费了该信息将会返回error
	// sarama.OffsetNewest:表明了为最新消息
	pc, err := c.consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
	if err != nil {
		return err
	}
	pcKey := strconv.Itoa(int(partition))
	key := topic + "_" + pcKey
	c.pcsLock.Lock()
	defer c.pcsLock.Unlock()
	if _, ok := c.pcs[key]; ok {
		return errors.New("the topic and partition is exist")
	}
	go c.listen(pc, result)
	return nil
}

func (c *Consumer) listen(pc sarama.PartitionConsumer, result chan *Result) {
	for {
		if pc != nil {
			select {
			case msg := <-pc.Messages():
				result <- &Result{
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
			case err := <-pc.Errors():
				result <- &Result{
					Error: err.Err,
					Data: Message{
						Topic:     err.Topic,
						Partition: err.Partition,
					},
				}
			}
		}
	}
}

func (c *Consumer) ConsumerGroupMessage(ctx context.Context, topic []string, result chan *Result, feedback <-chan bool) error {
	c.listenGroup(ctx, topic, result, feedback)
	return nil
}

func (c *Consumer) listenGroup(ctx context.Context, topic []string, result chan *Result, feedback <-chan bool) {
	/**
	 * Setup a new Sarama consumer group
	 */
	consumer := ConsumerHandler{
		Ready:    make(chan bool),
		Result:   result,
		Feedback: feedback,
	}
	go func() {
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := c.consumerGroup.Consume(ctx, topic, &consumer); err != nil {
				log.Printf("Error from consumer,topic: %s, err: %v", topic, err)
				time.Sleep(5 * time.Second)
				continue
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.Ready = make(chan bool)
		}
	}()
	<-consumer.Ready // Await till the consumer has been set up
}

// ClosePartition ...
func (c *Consumer) ClosePartition(topic string, partition int32) error {
	c.pcsLock.Lock()
	defer c.pcsLock.Unlock()
	pcKey := strconv.Itoa(int(partition))
	key := topic + "_" + pcKey
	if consumer, ok := c.pcs[key]; ok {
		return consumer.Close()
	}
	return nil
}

// Close ...
func (c *Consumer) Close() error {
	if err := c.consumer.Close(); err != nil {
		return err
	}
	c.pcsLock.Lock()
	defer c.pcsLock.Unlock()
	for _, consumer := range c.pcs {
		consumer.AsyncClose()
	}
	return nil
}

// Consumer ...
func (c *Consumer) Consumer() sarama.Consumer {
	return c.consumer
}

// ConsumerGroup ...
func (c *Consumer) ConsumerGroup() sarama.ConsumerGroup {
	return c.consumerGroup
}
