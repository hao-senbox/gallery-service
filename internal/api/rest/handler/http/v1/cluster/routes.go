package cluster

import (
	"gallery-service/internal/domain/service"
	"gallery-service/internal/infrastructure/database/mongo/repository"
	"github.com/gofiber/fiber/v2"
)

func (p *clusterHandlers) MapRoutes() func(router fiber.Router) {
	return func(router fiber.Router) {
		clusterRepository := repository.NewClusterRepository(p.log, p.cfg, p.mongoClient)
		folderRepository := repository.NewFolderRepository(p.log, p.cfg, p.mongoClient)

		p.ps = service.NewClusterService(p.cfg.Kafka, p.log, clusterRepository, folderRepository)
		router.Get("/", p.GetAllCluster)
		router.Get("/search", p.SearchCluster)
		router.Get("/components", p.GetClusterComponents)
		router.Get("/languages", p.GetClusterLanguages)
		router.Get("/:id", p.GetClusterByID)
		router.Get("/:id/folders", p.GetClusterFolders)

		router.Post("/", p.CreateCluster)
		router.Put("/", p.UpdateCluster)
		router.Delete("/:id", p.DeleteCluster)
	}
}
