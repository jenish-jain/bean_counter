package server

import (
	"bean_counter/internal/reporter"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jenish-jain/logger"
)

func (s *Server) InitRoutes() {
	router := s.routerGroups.rootRouter
	router.Use(requestIdMiddleWare)

	reporterHandler := reporter.NewHandler()
	reporterHandler.InitRoutes(router)
}

func requestIdMiddleWare(c *gin.Context) {
	u := uuid.New()
	c.Set("request_id", u)
	logger.DebugWithCtx(c, "request started", "method", c.Request.Method, "path", c.Request.URL.Path)
	c.Next()
}
