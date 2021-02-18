package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	kafka "github.com/segmentio/kafka-go"
)

func main() {

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
					Key:   []byte("Key-A"),
					Value: []byte("Hello World!"),
				},
				kafka.Message{
					Key:   []byte("Key-B"),
					Value: []byte("One!"),
				},
				kafka.Message{
					Key:   []byte("Key-C"),
					Value: []byte("Two!"),
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
			Brokers:   []string{"kafka-headless.banzaic-kafka:29092"},
			Topic:     topic,
			Partition: 0,
			MinBytes:  0,    // 10KB
			MaxBytes:  10e3, // 10MB
		})
		err := r.SetOffset(0)
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

	wg.Wait()
}
