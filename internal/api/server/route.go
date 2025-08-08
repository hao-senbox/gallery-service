package server

import (
	clusterV1 "gallery-service/internal/api/rest/handler/http/v1/cluster"
	folderV1 "gallery-service/internal/api/rest/handler/http/v1/folder"
	topicV1 "gallery-service/internal/api/rest/handler/http/v1/topic"

	"github.com/gofiber/fiber/v2"
)

func (s *server) routes() {

	clusterHandlers := clusterV1.NewClusterHandlers(s.log, s.cfg, s.mongoClient)
	folderHandlers := folderV1.NewFolderHandlers(s.log, s.cfg, s.mongoClient)
	topicHandlers := topicV1.NewTopicHandlers(s.log, s.cfg, s.mongoClient)

	// ===== Admin Routes =====
	adminAPI := s.fiber.Group("/api/v1/admin/gallery")

	clusterGroup := adminAPI.Group("/clusters", s.mw.Auth(s.consulClient))
	clusterGroup.Route("", clusterHandlers.MapRoutes())

	folderGroup := adminAPI.Group("/folders", s.mw.Auth(s.consulClient))
	folderGroup.Route("", folderHandlers.MapRoutes())

	topicGroup := adminAPI.Group("/topics", s.mw.Auth(s.consulClient), s.mw.ValidateSuperAdminRole())
	topicGroup.Route("", topicHandlers.MapRoutes())

	// ===== User Routes =====
	userAPI := s.fiber.Group("/api/v1/user/gallery")

	userClusterGroup := userAPI.Group("/clusters", s.mw.Auth(s.consulClient))
	userClusterGroup.Route("", clusterHandlers.MapRoutes())

	userFolderGroup := userAPI.Group("/folders", s.mw.Auth(s.consulClient))
	userFolderGroup.Route("", folderHandlers.MapRoutes())

	userTopicGroup := userAPI.Group("/topics", s.mw.Auth(s.consulClient))
	userTopicGroup.Route("", topicHandlers.MapRoutes())

	// ===== Health Check =====
	s.fiber.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).JSON(nil)
	})
}
