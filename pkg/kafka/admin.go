package kafka

import (
	"context"
	"gallery-service/pkg/zap"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Admin struct {
	client *kadm.Client
	log    zap.Logger
}

// NewAdmin creates a new Kafka admin client
func NewAdmin(log zap.Logger, brokers []string) *Admin {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		//kgo.WithLogger(kzap.New(log.GetLogger())),
	)

	if err != nil {
		log.DPanic(err)
	}

	admin := kadm.NewClient(client)
	return &Admin{client: admin, log: log}
}

func (a *Admin) TopicExists(topic Topic) bool {
	ctx := context.Background()
	topicsMetadata, err := a.client.ListTopics(ctx)
	if err != nil {
		a.log.DPanic(err)
	}
	for _, metadata := range topicsMetadata {
		if metadata.Topic == string(topic) {
			return true
		}
	}
	return false
}

func (a *Admin) CreateTopic(topic TopicConfig) error {
	ctx := context.Background()
	resp, err := a.client.CreateTopics(ctx, int32(topic.Partitions), int16(topic.ReplicationFactor), nil, string(topic.Name))

	if err != nil {
		a.log.DPanic(err)
	}

	for _, ctr := range resp {
		if ctr.Err != nil {
			a.log.Errorf("Unable to create topic '%s': %s", ctr.Topic, ctr.Err)
			return ctr.Err
		} else {
			a.log.Infof("Created topic '%s'\n", ctr.Topic)
		}
	}

	return nil
}

func (a *Admin) Close() {
	a.client.Close()
}
