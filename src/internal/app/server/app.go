package serverapp

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

var errServerStopped = errors.New("server has stopped")

type Server struct {
	port   int
	router *gin.Engine
}

func New(
	port int,
	router *gin.Engine,
) *Server {

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	return &Server{
		port:   port,
		router: router,
	}
}

func (s *Server) Serve() error {
	addr := fmt.Sprintf(":%s", s.port)
	err := s.router.Run(addr)

	if err == nil {
		return nil
	}

	return fmt.Errorf("%w: %s", errServerStopped, err.Error())
}
