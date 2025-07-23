package consumers

//
//import (
//	"context"
//	"gallery-service/config"
//	"gallery-service/internal/domain/service"
//	"gallery-service/internal/infrastructure/database/mongo/repository"
//	"gallery-service/pkg/kafka"
//	"gallery-service/pkg/zap"
//
//	"go.mongodb.org/mongo-driver/mongo"
//)
//
//type productUsageHistoryConsumerJob struct {
//	cfg                        *config.Config
//	log                        zap.Logger
//	productUsageHistoryService *service.HistoryService
//}
//
//func NewCreateProductHistoryConsumerJob(cfg *config.Config, log zap.Logger, db *mongo.Client) productUsageHistoryConsumerJob {
//	productUsageHistoryRepository := repository.NewProductUsageHistoryRepository(log, cfg, db)
//	return productUsageHistoryConsumerJob{
//		cfg:                        cfg,
//		log:                        log,
//		productUsageHistoryService: service.NewHistoryService(log, productUsageHistoryRepository),
//	}
//}
//
//func (p *productUsageHistoryConsumerJob) CreateProductUsageHistory() consumerJob {
//	//partition := int32(0)
//
//	return consumerJob{
//		Title:     "Create history after someone using cluster",
//		Group:     p.cfg.Kafka.Groups.ProductCreated,
//		Partition: nil,
//		Hld: func(ctx context.Context, message kafka.Message) error {
//			// Unmarshal the message
//			event := message.Value
//			//var eventData v1.ProductCreatedEvent
//			var eventData interface{}
//
//			// Get cluster usage history data from event data
//			if err := event.GetJsonData(&eventData); err != nil {
//				p.log.Errorf("failed to unmarshal data from event ProductUsageHistoryEvent consumeJob: %v", err)
//				return err
//			}
//
//			//// Create cluster usage history in repository
//			//createCommand := &productHistoryUsageCommands.CreateProductUsageHistoryCommand{
//			//	ProductID:    eventData.ProductID,
//			//	UserID:       eventData.UserID,
//			//	Username:     eventData.Username,
//			//	FullName:     eventData.FullName,
//			//	Organization: eventData.Organization,
//			//}
//			//if err := p.productUsageHistoryService.Commands.CreateProductUsageHistory.Handle(ctx, createCommand); err != nil {
//			//	p.log.Errorf("failed to create cluster usage history: %v", err)
//			//	return err
//			//}
//
//			// And emit event to all interested socket
//
//			p.log.Infof("Received Event: %v", eventData)
//
//			// Use your event bus or another eventing mechanism to emit event to interested consumers
//
//			// Example:
//			// emit.Emit("product_usage_history_created", event)
//
//			// Log the message in the logs for debugging and auditing purposes
//			// p.log.Infof("ProductUsageHistoryEvent: %+v", event)
//
//			// Log the message in the logs for debugging and auditing purposes
//			//p.log.Infof("CreateProductUsageHistory consumerJob message: %v", message.Value)
//
//			// Find all socket of userId and emit
//
//			return nil
//		},
//	}
//}
