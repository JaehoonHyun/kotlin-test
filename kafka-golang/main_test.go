package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func TestSimplePubSub(t *testing.T) {
	// to create topics when auto.create.topics.enable='false'
	topic := "my-topic"
	conn, controllerConn, err := Dial("tcp", "kafka-headless.banzaic-kafka:29092")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	defer controllerConn.Close()

	err = TopicDefault(controllerConn, topic, 1, 1)
	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {

		w := &kafka.Writer{
			Addr:     kafka.TCP("kafka-headless.banzaic-kafka:29092"),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		}

		for {

			err = w.WriteMessages(context.Background(),
				kafka.Message{
					Key:   []byte("500KB"),
					Value: []byte(RandStringRunes(1024 * 500)),
				},
			)

			if err != nil {
				log.Fatal(err)
			}

			log.Println("published")
		}

		wg.Done()
	}()

	wg.Add(1)
	go func() {
		log.Println("reader start")
		// make a new reader that consumes from topic-A, partition 0, at offset 0
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{"kafka-headless.banzaic-kafka:29092"},
			Topic:          topic,
			Partition:      0,
			MinBytes:       0,    // 10KB
			MaxBytes:       10e3, // 10MB
			CommitInterval: time.Second,
		})
		err := r.SetOffset(kafka.LastOffset)
		if err != nil {
			log.Fatal(err)
		}

		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				break
			}
			fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		}

		if err := r.Close(); err != nil {
			log.Fatal("failed to close reader:", err)
		}

		wg.Done()
	}()

	time.Sleep(10 * time.Second)
}

func TestPubConsumerGroup(t *testing.T) {

	// to create topics when auto.create.topics.enable='false'
	topic := "consumergroup2"
	partition := 4
	conn, controllerConn, err := Dial("tcp", "kafka-headless.banzaic-kafka:29092")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	defer controllerConn.Close()

	err = TopicDefault(controllerConn, topic, partition, 1)
	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {

		w := &kafka.Writer{
			Addr:     kafka.TCP("kafka-headless.banzaic-kafka:29092"),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		}

		for {

			//메시지 키의 개수에 따라 파티셔닝 된다.
			err = w.WriteMessages(context.Background(),
				kafka.Message{
					Key:   []byte("100B"),
					Value: []byte(RandStringRunes(100 * 1)),
				},
				kafka.Message{
					Key:   []byte("50B"),
					Value: []byte(RandStringRunes(50 * 1)),
				},
				kafka.Message{
					Key:   []byte("200B"),
					Value: []byte(RandStringRunes(200 * 1)),
				},
				kafka.Message{
					Key:   []byte("400B"),
					Value: []byte(RandStringRunes(400 * 1)),
				},
			)

			if err != nil {
				log.Fatal(err)
			}

			log.Println("published")
		}

		wg.Done()
	}()

	wg.Add(1)
	go func() {
		// make a new reader that consumes from topic-A
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{"kafka-headless.banzaic-kafka:29092"},
			GroupID:  "consumer-group-id",
			Topic:    topic,
			MinBytes: 10,   // 10B
			MaxBytes: 10e6, // 10MB
		})

		for i := 0; i < partition; i++ {
			go func() {
				for {
					m, err := r.ReadMessage(context.Background())
					if err != nil {
						break
					}
					fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
				}

				if err := r.Close(); err != nil {
					log.Fatal("failed to close reader:", err)
				}
			}()
		}
		wg.Done()
	}()

	time.Sleep(25 * time.Second)
}
