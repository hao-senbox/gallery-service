package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"gallery-service/pkg/zap"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer interface {
	SendMessage(ctx context.Context, data interface{}, promise *func(*kgo.Record, error))
	SendMessageWithKey(ctx context.Context, data interface{}, key string, promise *func(*kgo.Record, error))
	Close()
}

type producer struct {
	client *kgo.Client
	topic  Topic
}

func NewProducer(cfg Config, log zap.Logger, brokers []string, topic Topic) *producer {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ClientID("gallery-service"),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		//kgo.WithLogger(kzap.New(log.GetLogger())),

		// tuning performance
		kgo.ProducerBatchCompression(kgo.Lz4Compression()),
		kgo.MetadataMaxAge(time.Duration(cfg.ClientConfig.ProducerConfig.MetadataMaxAge)*time.Second),
		kgo.MaxBufferedRecords(cfg.ClientConfig.ProducerConfig.MaxBufferedRecords),
		kgo.ProducerBatchMaxBytes(int32(cfg.ClientConfig.ProducerConfig.ProducerBatchMaxBytes)),
		kgo.RecordPartitioner(kgo.UniformBytesPartitioner(cfg.ClientConfig.ProducerConfig.RecordPartitioner, false, false, nil)),
	)

	if err != nil {
		log.DPanic(err)
	}

	return &producer{client: client, topic: topic}
}

func (p *producer) SendMessage(ctx context.Context, data interface{}, promise *func(*kgo.Record, error)) {
	b, _ := json.Marshal(data)
	k := []byte(fmt.Sprintf("key-%d", time.Now().UnixNano()))
	p.client.Produce(ctx, &kgo.Record{Key: k, Topic: string(p.topic), Value: b}, *promise)
}

func (p *producer) SendMessageWithKey(ctx context.Context, data interface{}, key string, promise *func(*kgo.Record, error)) {
	b, _ := json.Marshal(data)
	k := []byte(key)
	p.client.Produce(ctx, &kgo.Record{Key: k, Topic: string(p.topic), Value: b}, *promise)
}

func (p *producer) Close() {
	p.client.Close()
}
