package server

import (
	clusterV1 "gallery-service/internal/api/rest/handler/http/v1/cluster"
	folderV1 "gallery-service/internal/api/rest/handler/http/v1/folder"
	"github.com/gofiber/fiber/v2"
)

func (s *server) routes() {
	// Define routes
	api := s.fiber.Group("/api/v1/gallery") // Root api

	//api.Use(s.mw.Auth()) // Middleware Auth

	clusterHandlers := clusterV1.NewClusterHandlers(s.log, s.cfg, s.mongoClient)
	clusterGroup := api.Group("/clusters", s.mw.Auth(s.consulClient)) // Group for cluster
	clusterGroup.Route("", clusterHandlers.MapRoutes())

	folderHandlers := folderV1.NewFolderHandlers(s.log, s.cfg, s.mongoClient)
	folderGroup := api.Group("/folders", s.mw.Auth(s.consulClient)) // Group for folder
	folderGroup.Route("", folderHandlers.MapRoutes())

	s.fiber.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).JSON(nil)
	})
}
