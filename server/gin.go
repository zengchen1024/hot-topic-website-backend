/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package server provides functionality for setting up and configuring a server for handling code repo operations.
package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/server-common-lib/interrupts"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/hot-topic-website-backend/config"
)

const (
	waitServerStart = 3 // 3s
)

var httpClient *http.Client

func RequestFilter(r *http.Request) bool {
	if strings.Contains(r.RequestURI, "/swagger/") ||
		strings.Contains(r.RequestURI, "/internal/heartbeat") {
		return false
	}

	return true
}

// StartWebServer starts a web server with the given configuration.
// It initializes the services, sets up the routers for different APIs, and starts the server.
// If TLS key and certificate are provided, it will use HTTPS.
// If removeCfg is true, it will remove the key and certificate files after starting the server.
func StartWebServer(removeCfg bool, port int, timeout time.Duration, cfg *config.Config) {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(logRequest())
	engine.UseRawPath = true

	// init services
	services, err := initServices(cfg)
	if err != nil {
		logrus.Error(err)

		return
	}

	// internal service api
	setInternalRouter("/internal", engine, cfg, &services)

	// start server
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           engine,
		ReadHeaderTimeout: time.Duration(cfg.ReadHeaderTimeout) * time.Second,
	}

	defer interrupts.WaitForGracefulShutdown()

	interrupts.ListenAndServe(srv, timeout)
}

func logRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("x-request-id", c.GetHeader("x-request-id"))

		startTime := time.Now()

		c.Next()

		endTime := time.Now()

		for _, ginErr := range c.Errors {
			logrus.Errorf("error on %s %s:\n%+v", c.Request.Method, c.Request.RequestURI, ginErr.Unwrap())
		}

		if !RequestFilter(c.Request) {
			return
		}

		log := fmt.Sprintf(
			"request_id: %s | %d | %d | %s | %s ",
			c.GetHeader("X-Request-Id"),
			c.Writer.Status(),
			endTime.Sub(startTime),
			c.Request.Method,
			c.Request.RequestURI,
		)

		logrus.Info(log)
	}
}
