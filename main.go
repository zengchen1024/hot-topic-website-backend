/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package main is the entry point for the application.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/hot-topic-website-backend/common/infrastructure/mongodb"
	"github.com/opensourceways/hot-topic-website-backend/config"
	"github.com/opensourceways/hot-topic-website-backend/server"
)

const (
	port        = 8888
	gracePeriod = 180
)

type options struct {
	service     ServiceOptions
	enableDebug bool
}

// ServiceOptions defines configuration parameters for the service.
type ServiceOptions struct {
	Port        int
	RemoveCfg   bool
	ConfigFile  string
	GracePeriod time.Duration
}

// Validate checks if the ServiceOptions are valid.
// It returns an error if the config file is missing.
func (o *ServiceOptions) Validate() error {
	if o.ConfigFile == "" {
		return fmt.Errorf("missing config-file")
	}

	return nil
}

// AddFlags adds flags for ServiceOptions to the provided FlagSet.
// It includes flags for port, remove-config, config-file, cert, key, and grace-period.
func (o *ServiceOptions) AddFlags(fs *flag.FlagSet) {
	fs.IntVar(&o.Port, "port", port, "Port to listen on.")

	fs.BoolVar(&o.RemoveCfg, "rm-cfg", false, "whether remove the cfg file after initialized .")

	fs.StringVar(&o.ConfigFile, "config-file", "", "Path to config file.")

	fs.DurationVar(&o.GracePeriod, "grace-period", time.Duration(gracePeriod)*time.Second,
		"On shutdown, try to handle remaining events for the specified duration.")
}

// Validate validates the options and returns an error if any validation fails.
func (o *options) Validate() error {
	return o.service.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) (options, error) {
	var o options

	o.service.AddFlags(fs)

	fs.BoolVar(
		&o.enableDebug, "enable_debug", false,
		"whether to enable debug model.",
	)

	err := fs.Parse(args)

	return o, err
}

// @securityDefinitions.apikey Internal
// @in header
// @name TOKEN
// @description Type "Internal" followed by a space and internal token.
func main() {
	o, err := gatherOptions(
		flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		os.Args[1:]...,
	)
	if err != nil {
		logrus.Errorf("new options failed, err:%s", err.Error())

		return
	}

	if err := o.Validate(); err != nil {
		logrus.Errorf("Invalid options, err:%s", err.Error())

		return
	}

	if o.enableDebug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug enabled.")
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		DisableColors: true,
		DisableQuote:  true,
	})

	// cfg
	cfg := new(config.Config)

	if err := config.LoadConfig(o.service.ConfigFile, cfg, o.service.RemoveCfg); err != nil {
		logrus.Errorf("main load config, err:%s", err.Error())

		return
	}

	// init cfg
	if err := cfg.Init(); err != nil {
		logrus.Errorf("init cfg failed, err:%s", err.Error())

		return
	}

	// mongodb
	if err := mongodb.Init(&cfg.MongoDB); err != nil {
		logrus.Error(err)

		return
	}

	defer exitMongoService()

	// init api doc
	//api.Init(cfg.SwaggerInfo)

	// run
	server.StartWebServer(o.service.RemoveCfg, o.service.Port, o.service.GracePeriod, cfg)
}

func exitMongoService() {
	if err := mongodb.Close(); err != nil {
		logrus.Error(err)
	}
}
