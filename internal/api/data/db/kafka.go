/*
 * Kafka initialization using IBM/sarama
 */

package db

import (
	"github.com/IBM/sarama"
	"github.com/NSObjects/go-template/internal/configs"
)

func NewKafkaProducer(cfg configs.KafkaConfig) (sarama.SyncProducer, error) {
	sc := sarama.NewConfig()
	sc.ClientID = cfg.ClientID
	sc.Producer.RequiredAcks = sarama.WaitForAll
	sc.Producer.Retry.Max = 3
	sc.Producer.Return.Successes = true
	return sarama.NewSyncProducer(cfg.Brokers, sc)
}
