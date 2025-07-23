package consumers

import (
	"gallery-service/config"
	"gallery-service/pkg/asyncjob"
	"gallery-service/pkg/kafka"
	"gallery-service/pkg/zap"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type consumerJob struct {
	Title     string
	Group     string
	Partition *int32
	Hld       func(ctx context.Context, message kafka.Message) error
}

type consumerEngine struct {
	cfg *config.Config
	log zap.Logger
	db  *mongo.Client
}

func NewEngine(cfg *config.Config, log zap.Logger, db *mongo.Client) *consumerEngine {
	return &consumerEngine{cfg: cfg, log: log, db: db}
}

func (engine *consumerEngine) Start(ctx context.Context) error {
	//productHistoryUsageCommands := NewCreateProductHistoryConsumerJob(engine.cfg, engine.log, engine.db)
	//
	//if err := engine.startSubTopic(
	//	ctx,
	//	engine.cfg.Kafka.Topics.ProductCreated,
	//	true,
	//	productHistoryUsageCommands.CreateProductUsageHistory(),
	//); err != nil {
	//	return err
	//}

	return nil
}

type GroupJob interface {
	Run(ctx context.Context) error
}

func (engine *consumerEngine) startSubTopic(ctx context.Context, topic kafka.TopicConfig, isParallel bool, hdls ...consumerJob) error {
	k := kafka.NewKafkaConnection(engine.cfg.Kafka, engine.log, topic)

	for _, item := range hdls {
		engine.log.Infof("Setup consumer for: %s\n", item.Title)
	}

	getHld := func(job *consumerJob) func(ctx context.Context, message kafka.Message) error {
		return func(ctx context.Context, message kafka.Message) error {
			engine.log.Infof("running job for %s. Value: %v", job.Title, message)

			group := asyncjob.NewGroup(isParallel, asyncjob.NewJob(func(ctx context.Context) error {
				return job.Hld(ctx, message)
			}))

			if err := group.Run(ctx); err != nil {
				engine.log.Errorf("error running job for %s. Value: %v and cannot retry", job.Title, message)
				// handle dead letter queues / retry / etc
				// if err := deadLetterQueue.Send(ctx, message); err != nil {
				//     engine.log.Errorf("failed to send message to dead letter queue: %v", err)
				//     return err
				// }

				// // retry after a certain delay
				// time.Sleep(time.Duration(rand.Intn(5)) * time.Second)

				// retry
				// if err := retry(ctx, message); err != nil {
				//     engine.log.Errorf("failed to retry job for %s. Value: %v", job.Title, message)
				//     return err
				// }

				// // if retry limit is reached, mark message as failed
				// if err := markAsFailed(ctx, message); err != nil {
				//     engine.log.Errorf("failed to mark message as failed: %v", err)
				//     return err
				// }

				return err
			}

			return nil
		}
	}

	for i := range hdls {
		go func() {
			errChan := make(chan error)
			k.ConsumeMessages(ctx, kafka.Topic(topic.Name), kafka.Group(hdls[i].Group), hdls[i].Partition, getHld(&hdls[i]), errChan)
			engine.log.Infof("Consumer for %s is running...\n", hdls[i].Title)
			for {
				select {
				case <-ctx.Done():
					close(errChan)
					return
				case v := <-errChan:
					if v != nil {
						engine.log.Error(v)
						// handle error

						close(errChan)
						return
					}
				}
			}
		}()
	}

	return nil
}
