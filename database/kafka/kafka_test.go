// Package kafka
//
// @author: xwc1125
// @date: 2021/3/16
package kafka

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"testing"

	"github.com/IBM/sarama"
)

func TestProducer(t *testing.T) {
	k := &KafkaConfig{
		IsAsync: false,
		Addrs:   []string{"127.0.0.1:9092"},
		GroupId: "om-default",
		Topic:   []string{"test2"},
	}
	kafka, err := New(k)
	if err != nil {
		panic(err)
	}

	result := make(chan *Result)
	go func() {
		for {
			select {
			case r := <-result:
				if r.Error != nil {
					// 出错
					log.Fatal("Producer", r.Error)
				}
				fmt.Println("offsetCfg:", r.Data.Offset, " partitions:", r.Data.Partition, " metadata:", r.Data.Metadata, " value:", string(r.Data.Value))
			}
		}
	}()
	producer, err := kafka.NewProducer()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		producer.ProducerMessage(k.Topic[0], "key"+strconv.Itoa(i), "value"+strconv.Itoa(i), 0, result)
	}
}

func TestConsumer(t *testing.T) {
	k := &KafkaConfig{
		IsAsync: false,
		Addrs:   []string{"127.0.0.1:9092"},
		GroupId: "om-default",
		Topic:   []string{"test2"},
	}
	kafka, err := New(k)
	if err != nil {
		panic(err)
	}

	result2 := make(chan *Result)
	feedback := make(chan bool)
	consumer, err := kafka.NewConsumerGroup()
	err = consumer.ConsumerGroupMessage(context.Background(), k.Topic, result2, feedback)
	// partitions, err := consumer.Consumer().Partitions(k.Topic[0])
	// fmt.Println(partitions)
	// for i := 0; i < len(partitions); i++ {
	//	err = consumer.ConsumerMessage(k.Topic[0], int32(i), result2)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	// }
	for {
		select {
		case r := <-result2:
			if r.Error != nil {
				// 出错
				log.Fatal("Consumer", r.Error)
			} else {
				log.Println("Consumer", r.Data.String())
				feedback <- true
			}
		}
	}
}

// 生产者接口
func Test_Producer(t *testing.T) {
	fmt.Printf("producer_test\n")
	config := sarama.NewConfig()
	// 等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForLocal
	// 随机的分区类型：返回一个分区器，该分区器每次选择一个随机分区
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 是否等待成功和失败后的响应
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = sarama.V0_11_0_2

	// 使用给定代理地址和配置创建一个同步生产者
	// producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	producer, err := sarama.NewAsyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		fmt.Printf("producer_test create producer error :%s\n", err.Error())
		return
	}

	defer producer.AsyncClose()

	// 构建发送的消息，
	msg := &sarama.ProducerMessage{
		Topic: "cz-default-logs",
		// Partition: int32(10), // 分区
		Key: sarama.StringEncoder("go_test"),
	}

	value := "this is message"
	for {
		// fmt.Scanln(&value)
		msg.Value = sarama.ByteEncoder(value)
		fmt.Printf("input [%s]\n", value)

		// SendMessage：该方法是生产者生产给定的消息
		// 生产成功的时候返回该消息的分区和所在的偏移量
		// 生产失败的时候返回error
		// 同步：
		// partition, offset, err := producer.SendMessage(msg)

		// 异步：send to chain
		producer.Input() <- msg

		select {
		case suc := <-producer.Successes():
			fmt.Println(fmt.Sprintf("offset: %d,  timestamp: %s", suc.Offset, suc.Timestamp.String()))
			return
		case fail := <-producer.Errors():
			fmt.Println(fail.Err.Error())
			return
		}
	}
}

// 消费者接口
func Test_Consumer(t *testing.T) {
	fmt.Printf("consumer_test")

	config := sarama.NewConfig()
	// config.Consumer.Return.Errors = true
	config.Version = sarama.V0_11_0_2
	config.Net.SASL.Enable = true
	config.Net.SASL.User = "ckafka-7a758rn2#changyu_test"
	config.Net.SASL.Password = "cytest123"

	// consumer
	consumer, err := sarama.NewConsumer([]string{"xxxx.ckafka.tencentcloudmq.com:6020"}, config)
	if err != nil {
		fmt.Printf("consumer_test create consumer error %s\n", err.Error())
		return
	}

	defer consumer.Close()

	// ConsumePartition方法根据主题，分区和给定的偏移量创建创建了相应的分区消费者
	// 如果该分区消费者已经消费了该信息将会返回error
	// sarama.OffsetNewest:表明了为最新消息
	// pc, err := consumer.ConsumePartition("test", int32(partition), sarama.OffsetNewest)

	pc, err := consumer.ConsumePartition("test_changyu_ts_send", 0, sarama.OffsetOldest)
	if err != nil {
		fmt.Printf("try create partition_consumer error %s\n", err.Error())
		return
	}
	defer pc.Close()

	for {
		select {
		case msg := <-pc.Messages():
			fmt.Println(fmt.Sprintf("msg offset: %d, partition: %d, timestamp: %s, value: %s\n",
				msg.Offset, msg.Partition, msg.Timestamp.String(), string(msg.Value)))
		case err := <-pc.Errors():
			fmt.Println(err.Error())
		}
	}
}

func Test_GroupConsumer(t *testing.T) {
	log.Println("Starting a new Sarama consumer")
	assignor := "range"
	oldest := true
	verbose := true

	brokers := []string{"xxxx:6020"}
	topics := []string{"test_changyu_ts_send"}
	group := "btoe_screenshot_task"

	if verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	config := sarama.NewConfig()
	config.Version = sarama.V0_11_0_2
	config.Net.SASL.Enable = true
	config.Net.SASL.User = "ckafka-7a758rn2#changyu_test"
	config.Net.SASL.Password = "cytest123"

	switch assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case "roundrobin":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case "range":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", assignor)
	}

	if oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	/**
	 * Setup a new Sarama consumer group
	 */
	consumer := ConsumerHandler{
		Ready: make(chan bool),
		IsLog: true,
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, topics, &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.Ready = make(chan bool)
		}
	}()

	<-consumer.Ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Println("terminating: context cancelled")
	case <-sigterm:
		log.Println("terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}

func Test_Metadata(t *testing.T) {
	fmt.Printf("metadata test\n")

	config := sarama.NewConfig()
	config.Version = sarama.V0_11_0_2

	client, err := sarama.NewClient([]string{"localhost:9092"}, config)
	if err != nil {
		fmt.Printf("metadata_test try create client err :%s\n", err.Error())
		return
	}

	defer client.Close()

	// get topic set
	topics, err := client.Topics()
	if err != nil {
		fmt.Printf("try get topics err %s\n", err.Error())
		return
	}

	fmt.Printf("topics(%d):\n", len(topics))

	for _, topic := range topics {
		fmt.Println(topic)
	}

	// get broker set
	brokers := client.Brokers()
	fmt.Printf("broker set(%d):\n", len(brokers))
	for _, broker := range brokers {
		fmt.Println(broker.Addr())
	}
}

func Test_GroupConsumerOrigin(t *testing.T) {
	log.Println("Starting a new Sarama consumer")
	assignor := "range"
	oldest := true
	verbose := true

	brokers := []string{"127.0.0.1:9092"}
	topics := []string{"btoe_web_evidence"}
	group := "btoe_screenshot_task"

	if verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	config := sarama.NewConfig()
	config.Version = sarama.V0_11_0_2
	config.Net.SASL.Enable = true
	config.Net.SASL.User = "ckafka-7a758rn2#changyu_test"
	config.Net.SASL.Password = "cytest123"

	switch assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case "roundrobin":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case "range":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", assignor)
	}

	if oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	/**
	 * Setup a new Sarama consumer group
	 */
	consumer := ConsumerHandler{
		Ready: make(chan bool),
		IsLog: true,
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, topics, &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.Ready = make(chan bool)
		}
	}()

	<-consumer.Ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Println("terminating: context cancelled")
	case <-sigterm:
		log.Println("terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}
