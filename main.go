package main

import (
	"github.com/gosidekick/goconfig"
	_ "github.com/gosidekick/goconfig/toml"
	"github.com/sirupsen/logrus"
	"mqtt-go-playground/mqtt"
	"mqtt-go-playground/serverMgmt"
	"mqtt-go-playground/service_cover"
)

var version string

// Config the large config struct that contains the whole app config
type Config struct {
	LogLevel string `cfgDefault:"INFO"`
	Mqtt     mqtt.Config
}

func main() {
	logrus.Infof("version=%s", version)
	cfg := &Config{}
	goconfig.Path = "./config"
	goconfig.File = "config.toml"
	err := goconfig.Parse(cfg)
	if err != nil {
		logrus.Error("config error")
		logrus.Fatal(err)
	}
	serverMgmt.Log(cfg)
	serverMgmt.CustomLogging(cfg.LogLevel)

	//init mqtt service
	mqtt.Init(&cfg.Mqtt)

	defer mqtt.Disconnect()

	// Start the services
	go service_cover.Start()

	// Block until we receive a shutdown signal
	<-serverMgmt.GracefulStop

	// Stop all background work, previously defined with defers
	logrus.Info("Received a quit signal. Stopping background work now...")

}
