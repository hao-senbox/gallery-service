package topic

import (
	"gallery-service/internal/domain/service"
	"gallery-service/internal/infrastructure/database/mongo/repository"

	"github.com/gofiber/fiber/v2"
)

func (p *topicHandlers) MapRoutesAdmin() func(router fiber.Router) {
	return func(router fiber.Router) {
		topicRepository := repository.NewTopicRepository(p.log, p.cfg, p.mongoClient)
		folderRepository := repository.NewFolderRepository(p.log, p.cfg, p.mongoClient)

		p.ps = service.NewTopicService(p.cfg.Kafka, p.log, topicRepository, folderRepository)
		router.Get("", p.GetAllTopic)
		router.Get("/search", p.SearchTopic)
		router.Get("/components", p.GetTopicComponents)
		router.Get("/languages", p.GetTopicLanguages)
		router.Get("/:id", p.GetTopicByID)

		router.Post("", p.CreateTopic)
		router.Put("", p.UpdateTopic)
		router.Delete("/:id", p.DeleteTopic)
	}
}

func (p *topicHandlers) MapRoutesUser() func(router fiber.Router) {
	return func(router fiber.Router) {
		topicRepository := repository.NewTopicRepository(p.log, p.cfg, p.mongoClient)
		folderRepository := repository.NewFolderRepository(p.log, p.cfg, p.mongoClient)

		p.ps = service.NewTopicService(p.cfg.Kafka, p.log, topicRepository, folderRepository)
		router.Get("", p.GetAllTopic4App)
		router.Get("/:id", p.GetTopicByID)
	}
}

func (p *topicHandlers) MapRoutesGateway() func(router fiber.Router) {
	return func(router fiber.Router) {
		topicRepository := repository.NewTopicRepository(p.log, p.cfg, p.mongoClient)
		folderRepository := repository.NewFolderRepository(p.log, p.cfg, p.mongoClient)

		p.ps = service.NewTopicService(p.cfg.Kafka, p.log, topicRepository, folderRepository)
		router.Get("", p.GetAllTopic4Gateway)
		router.Get("/:id", p.GetTopicByID4Gateway)
	}
}
