package kafka

import (
	"context"
	"gallery-service/pkg/es"
	"gallery-service/pkg/zap"
	"sync"

	"github.com/twmb/franz-go/pkg/kgo"
)

type TopicConfig struct {
	Name              string `mapstructure:"name" validate:"required"`
	Partitions        int    `mapstructure:"partitions" validate:"required"`
	ReplicationFactor int    `mapstructure:"replication_factor" validate:"required"`
}

type Topics struct {
	ProductCreated TopicConfig `mapstructure:"product_created" validate:"required"`
	ProductUpdated TopicConfig `mapstructure:"product_updated" validate:"required"`
	ProductDeleted TopicConfig `mapstructure:"product_deleted" validate:"required"`
	ProductPicked  TopicConfig `mapstructure:"product_picked" validate:"required"`
}

type ProducerConfig struct {
	MetadataMaxAge        int `mapstructure:"metadata_max_age"`
	MaxBufferedRecords    int `mapstructure:"max_buffered_records"`
	ProducerBatchMaxBytes int `mapstructure:"producer_batch_max_bytes"`
	RecordPartitioner     int `mapstructure:"record_partitioner"`
}

type ConsumerConfig struct {
	MetadataMaxAge         int `mapstructure:"metadata_max_age"`
	FetchMaxBytes          int `mapstructure:"fetch_max_bytes"`
	FetchMaxPartitionBytes int `mapstructure:"fetch_max_partition_bytes"`
	FetchMaxWait           int `mapstructure:"fetch_max_wait"`
}

type ClientConfig struct {
	ProducerConfig ProducerConfig `mapstructure:"producer_config"`
	ConsumerConfig ConsumerConfig `mapstructure:"consumer_config"`
}

type Groups struct {
	ProductCreated string `mapstructure:"product_created" validate:"required"`
	ProductUpdated string `mapstructure:"product_updated" validate:"required"`
	ProductDeleted string `mapstructure:"product_deleted" validate:"required"`
	ProductPicked  string `mapstructure:"product_picked" validate:"required"`
}

type Config struct {
	Brokers           []string     `mapstructure:"brokers" validate:"required"`
	ReplicationFactor int          `mapstructure:"replication_factor" validate:"required"`
	Topics            Topics       `mapstructure:"topics" validate:"required"`
	Groups            Groups       `mapstructure:"groups" validate:"required"`
	ClientConfig      ClientConfig `mapstructure:"client_config" validate:"required"`
}

type Topic string
type Group string

type Message struct {
	Key   string
	Value es.BaseEvent
}

type kafkaConnection struct {
	producer  *producer
	consumers map[Topic]*consumer
	cfg       Config
	log       zap.Logger
	mu        sync.Mutex
}

var (
	kfrConnection *kafkaConnection
	consumers     = make(map[Topic]*consumer)
)

func NewKafkaConnection(cfg Config, log zap.Logger, topic TopicConfig) *kafkaConnection {
	admin := NewAdmin(log, cfg.Brokers)
	defer admin.Close()
	if !admin.TopicExists(Topic(topic.Name)) {
		if err := admin.CreateTopic(topic); err != nil {
			return nil
		}
	}

	if kfrConnection != nil {
		return kfrConnection
	}

	kfrConnection = &kafkaConnection{
		cfg:       cfg,
		log:       log,
		consumers: consumers,
	}

	return kfrConnection
}

func (k *kafkaConnection) SendMessage(ctx context.Context, data interface{}, promise *func(*kgo.Record, error)) {
	k.mu.Lock()
	defer k.mu.Unlock()

	if k.producer == nil {
		k.producer = NewProducer(k.cfg, k.log, k.cfg.Brokers, Topic(k.cfg.Topics.ProductCreated.Name))
	}

	k.producer.SendMessage(ctx, data, promise)
}

func (k *kafkaConnection) SendMessageWithKey(ctx context.Context, data interface{}, key string, promise *func(*kgo.Record, error)) {
	k.mu.Lock()
	defer k.mu.Unlock()

	if k.producer == nil {
		k.producer = NewProducer(k.cfg, k.log, k.cfg.Brokers, Topic(k.cfg.Topics.ProductCreated.Name))
	}

	k.producer.SendMessageWithKey(ctx, data, key, promise)
}

func (k *kafkaConnection) ConsumeMessages(
	ctx context.Context,
	topic Topic,
	group Group,
	partition *int32,
	handler func(ctx context.Context, message Message) error,
	errCh chan<- error,
) {
	k.mu.Lock()
	defer k.mu.Unlock()

	if _, ok := k.consumers[topic]; !ok {
		k.consumers[topic] = NewConsumer(k.cfg, k.log, k.cfg.Brokers, topic, group, partition, handler, errCh)
	}

	k.consumers[topic].StartPolling(ctx)
}

func (k *kafkaConnection) Close() {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.producer.Close()
	for _, c := range k.consumers {
		c.Close()
	}
}
