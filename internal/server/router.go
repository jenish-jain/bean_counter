package server

import (
	"bean_counter/internal/reporter"

	"github.com/jenish-jain/logger"
)

func (s *Server) InitRoutes() {
	router := s.routerGroups.rootRouter
	router.Use(logger.AttachRequestIdToRequests)

	reporterHandler := reporter.NewHandler()
	reporterHandler.InitRoutes(router)
}
