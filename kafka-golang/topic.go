package main

import kafka "github.com/segmentio/kafka-go"

func TopicDefault(controllerConn *kafka.Conn, topic string, partitions int, replicas int) error {

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     partitions,
			ReplicationFactor: replicas,
		},
	}

	return controllerConn.CreateTopics(topicConfigs...)
}
