package server

import (
	"context"
	"errors"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine            *gin.Engine
	port              int
	ServerStoppedChan chan struct{}
}

func NewHttpServer(engine *gin.Engine, port int) *Server {
	return &Server{
		engine:            engine,
		port:              port,
		ServerStoppedChan: make(chan struct{}),
	}
}

func (s *Server) Run(shutdownChan <-chan struct{}) error {
	defer func() {
		s.ServerStoppedChan <- struct{}{}
		close(s.ServerStoppedChan)
	}()
	srv := &http.Server{
		Addr:        fmt.Sprintf(":%d", s.port),
		Handler:     s.engine,
		ReadTimeout: 10 * time.Second,
	}

	errCh := make(chan error)
	defer close(errCh)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-shutdownChan:
		err := srv.Shutdown(context.Background())
		if err != nil {
			logger.WithError(err).Errorln("error stop http server")
			return err
		}

		return nil
	case err := <-errCh:
		return err
	}
}
