package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jenish-jain/logger"
)

type Server struct {
	engine       *gin.Engine
	routerGroups RouterGroups
}

type RouterGroups struct {
	rootRouter *gin.Engine
}

func NewServer() *Server {

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		logger.Info(fmt.Sprintf("Endpoint %s is declared via handler %s", absolutePath, handlerName),
			"event", "INIT_HTTP_SERVER",
			"method", httpMethod)
	}
	engine := gin.New()

	engine.Use(gin.Recovery())

	return &Server{
		engine: engine,
		routerGroups: RouterGroups{
			rootRouter: engine,
		},
	}
}

func (s *Server) Run() {
	s.InitRoutes()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logger.Info("Defaulting to port %s", port)
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: s.engine,
	}
	go listenServer(srv)
	waitForShutdown(srv)
}

func listenServer(server *http.Server) {
	logger.Info(fmt.Sprintf("listening server on %s", server.Addr),
		"event", "INIT_HTTP_SERVER")
	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func waitForShutdown(server *http.Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGINT,
		syscall.SIGTERM)
	_ = <-sig
	logger.Info("server shutting down", "event", "SHUTDOWN")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		logger.Error("server forced to shutdown", "event", "SHUTDOWN", "data", map[string]string{
			"error": fmt.Sprintf("%+v", err),
		})
	}
	logger.Info("server shutdown complete", "event", "SHUTDOWN")
}
