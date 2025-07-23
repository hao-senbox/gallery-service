package folder

import (
	"gallery-service/internal/domain/service"
	"gallery-service/internal/infrastructure/database/mongo/repository"
	"github.com/gofiber/fiber/v2"
)

func (p *folderHandlers) MapRoutes() func(router fiber.Router) {
	return func(router fiber.Router) {
		folderRepository := repository.NewFolderRepository(p.log, p.cfg, p.mongoClient)

		p.ps = service.NewFolderService(p.cfg.Kafka, p.log, folderRepository)
		router.Get("/", p.GetAllFolder)
		router.Get("/search", p.SearchFolder)
		router.Get("/:id", p.GetFolderByID)

		router.Post("/", p.CreateFolder)
		router.Put("/", p.UpdateFolder)
		router.Delete("/:id", p.DeleteFolder)
	}
}
