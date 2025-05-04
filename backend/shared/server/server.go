package server

import (
	"strconv"
	"to-read/router"
	"to-read/utils/logs"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Server struct {
	Hostname string `yaml:"hostname"`
	Port     int    `yaml:"port"`
}

func Run(s Server) error {
	e := echo.New()

	// Router and middleware
	router.Load(e)

	// Default Configurations
	if s.Hostname == "" {
		s.Hostname = "127.0.0.1"
	}
	if s.Port == 0 {
		s.Port = 3435
	}

	address := s.Hostname + ":" + strconv.Itoa(s.Port)
	err := e.Start(address)
	if err != nil {
		logs.Error("Server run failed at "+address+". ", zap.Error(err))
		return err
	}
	return nil
}
