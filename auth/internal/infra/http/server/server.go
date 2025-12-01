package server

import (
	"go.uber.org/zap"

	"github.com/DangeL187/erax"
	"github.com/gin-gonic/gin"

	"auth/internal/app"
	"auth/internal/infra/http/routes"
)

type Server struct {
	engine *gin.Engine
}

func (s *Server) Run(addr string) error {
	zap.S().Infof("HTTP server launched on http://%s", addr)

	err := s.engine.Run(addr)
	if err != nil {
		return erax.Wrap(err, "failed to start HTTP server")
	}

	return nil
}

func NewServer(app *app.App) *Server {
	engine := gin.New()
	engine.Use(gin.Recovery(), gin.Logger())

	routes.SetupRoutes(engine, app)

	return &Server{engine: engine}
}
