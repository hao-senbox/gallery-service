package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"gallery-service/pkg/zap"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Consumer interface {
	StartPolling(ctx context.Context)
	Close()
}

type tp struct {
	t string
	p int32
}

type splitConsume struct {
	// Using BlockRebalanceOnCommit means we do not need a mu to manage
	// consumers, unlike the autocommit normal example.
	consumers map[tp]*consumer
}

func (s *splitConsume) assigned(
	log zap.Logger,
	group Group,
	partition *int32,
	handler func(ctx context.Context, message Message) error,
	chHdlErr chan<- error,
) func(ctx context.Context, cl *kgo.Client, assigned map[string][]int32) {
	return func(ctx context.Context, cl *kgo.Client, assigned map[string][]int32) {
		if partition != nil {
			if *partition < 0 || *partition > int32(len(assigned)) {
				log.DPanicf("invalid partition: %+v", partition)
			}
		}

		for topic, partitions := range assigned {
			if partition != nil {
				for _, p := range partitions {
					if p == *partition {
						if _, ok := s.consumers[tp{topic, p}]; !ok {
							c := &consumer{
								client:    cl,
								topic:     Topic(topic),
								group:     group,
								partition: &p,
								log:       log,

								quit: make(chan struct{}),
								done: make(chan struct{}),
								recs: make(chan []*kgo.Record, 5),
							}
							s.consumers[tp{topic, p}] = c
							go c.consume(ctx, handler, chHdlErr)
						}
					}
				}
			} else {
				for _, p := range partitions {
					if _, ok := s.consumers[tp{topic, p}]; !ok {
						c := &consumer{
							client:    cl,
							topic:     Topic(topic),
							group:     group,
							partition: &p,
							log:       log,

							quit: make(chan struct{}),
							done: make(chan struct{}),
							recs: make(chan []*kgo.Record, 5),
						}
						s.consumers[tp{topic, p}] = c
						go c.consume(ctx, handler, chHdlErr)
					}
				}
			}
		}
	}
}

// In this example, each partition consumer commits itself. Those commits will
// fail if partitions are lost, but will succeed if partitions are revoked. We
// only need one revoked or lost function (and we name it "lost").
func (s *splitConsume) lost(_ context.Context, _ *kgo.Client, lost map[string][]int32) {
	var wg sync.WaitGroup
	defer wg.Wait()

	for topic, partitions := range lost {
		for _, partition := range partitions {
			tp := tp{topic, partition}
			pc := s.consumers[tp]
			delete(s.consumers, tp)
			close(pc.quit)
			fmt.Printf("waiting for work to finish t %s p %d\n", topic, partition)
			wg.Add(1)
			go func() { <-pc.done; wg.Done() }()
		}
	}
}

func (s *splitConsume) poll(ctx context.Context, cl *kgo.Client) {
	for {
		// PollRecords is strongly recommended when using
		// BlockRebalanceOnPoll. You can tune how many records to
		// process at once (upper bound -- could all be on one
		// partition), ensuring that your processor loops complete fast
		// enough to not block a rebalance too long.
		fetches := cl.PollRecords(ctx, 10000)
		if fetches.IsClientClosed() {
			return
		}
		//fetches.EachError(func(_ string, _ int32, err error) {
		//	// Note: you can delete this block, which will result
		//	// in these fetchErrors being sent to the partition
		//	// consumers, and then you can handle the fetchErrors there.
		//	panic(err)
		//})
		fetches.EachPartition(func(p kgo.FetchTopicPartition) {
			tp := tp{p.Topic, p.Partition}

			consumer, exists := s.consumers[tp]
			if !exists || consumer == nil {
				// Handle the error (e.g., log it or retry)
				fmt.Printf("consumer lost or revoked t %s p %d\n", p.Topic, p.Partition)
				return
			}

			// Since we are using BlockRebalanceOnPoll, we can be
			// sure this partition consumer exists:
			//
			// * onAssigned is guaranteed to be called before we
			// fetch offsets for newly added partitions
			//
			// * onRevoked waits for partition consumers to quit
			// and be deleted before re-allowing polling.
			s.consumers[tp].recs <- p.Records
		})
		cl.AllowRebalance()
	}
}

type consumer struct {
	client    *kgo.Client
	topic     Topic
	group     Group
	partition *int32
	s         *splitConsume

	log zap.Logger

	quit chan struct{}
	done chan struct{}
	recs chan []*kgo.Record
}

func NewConsumer(
	cfg Config,
	log zap.Logger,
	brokers []string,
	topic Topic,
	group Group,
	partition *int32,
	handler func(ctx context.Context, message Message) error,
	chHdlErr chan<- error,
) *consumer {
	s := &splitConsume{
		consumers: make(map[tp]*consumer),
	}

	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(string(group)),
		kgo.ConsumeTopics(string(topic)),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
		//kgo.Balancers(kgo.RangeBalancer()),
		//kgo.WithLogger(kzap.New(log.GetLogger())),
		kgo.OnPartitionsAssigned(s.assigned(log, group, partition, handler, chHdlErr)),
		kgo.OnPartitionsRevoked(s.lost),
		kgo.OnPartitionsLost(s.lost),
		kgo.DisableAutoCommit(),
		kgo.BlockRebalanceOnPoll(),
		//kgo.AdjustFetchOffsetsFn(func(ctx context.Context, m map[string]map[int32]kgo.Offset) (map[string]map[int32]kgo.Offset, error) {
		//	for k, v := range m {
		//		for i := range v {
		//			m[k][i] = kgo.NewOffset().At(-2).WithEpoch(-1)
		//		}
		//	}
		//	return m, nil
		//}),

		//tuning performance
		kgo.MetadataMaxAge(time.Duration(cfg.ClientConfig.ConsumerConfig.MetadataMaxAge)*time.Second),
		kgo.FetchMaxBytes(int32(cfg.ClientConfig.ConsumerConfig.FetchMaxBytes)),
		kgo.FetchMaxPartitionBytes(int32(cfg.ClientConfig.ConsumerConfig.FetchMaxPartitionBytes)),
		kgo.FetchMaxWait(time.Duration(cfg.ClientConfig.ConsumerConfig.FetchMaxWait)*time.Millisecond),
	)

	if err != nil {
		log.DPanic(err)
	}

	return &consumer{
		client:    client,
		topic:     topic,
		group:     group,
		partition: partition,
		s:         s,
		log:       log,
	}
}

func (c *consumer) consume(ctx context.Context, handler func(ctx context.Context, message Message) error, chHdlErr chan<- error) {
	for {
		select {
		case <-c.quit:
			return
		case <-ctx.Done():
			return
		case recs := <-c.recs:
			// simulate work
			for _, rec := range recs {
				message := &Message{
					Key: string(rec.Key),
				}
				if err := json.Unmarshal(rec.Value, &message.Value); err != nil {
					c.log.Errorf("Error unmarshalling: %v t: %s p: %+v offset %d\n", err, c.topic, c.partition, rec.Offset)
				}

				go func() {
					err := handler(ctx, *message)
					if err != nil {
						c.log.Errorf("Error processing message: %v t: %s p: %+v offset %d\n", err, c.topic, c.partition, rec.Offset)
						chHdlErr <- err
					}
				}()
			}
			c.log.Infof("Some sort of work done, about to commit t %s p %+v\n", c.topic, c.partition)
			err := c.client.CommitRecords(ctx, recs...)
			if err != nil {
				c.log.Errorf("Error when committing offsets to kafka err: %v t: %s p: %+v offset %d\n", err, c.topic, c.partition, recs[len(recs)-1].Offset+1)
			}
		}
	}
}

func (c *consumer) StartPolling(ctx context.Context) {
	go c.s.poll(ctx, c.client)
}

func (c *consumer) Close() {
	close(c.quit)
	close(c.recs)
	c.client.Close()
}
