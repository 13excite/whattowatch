package conf

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

const DefaultConfigPath = "/etc/kp_config.yaml"

// C is the global configuration
var C = Config{}

// Defaults returns config's object with default values
func (conf *Config) Defaults() {
	conf.LogLevel = "info"
	conf.LogEncoding = "console"
	conf.LoggerColor = true
	conf.LoggerDisableStacktrace = true
	conf.LoggerDevMode = true
	conf.LoggerDisableCaller = false
	conf.LoggerDisabledHttp = []string{"/version"}
	conf.ServerHost = "127.0.0.1"
	conf.ServerPort = "8081"
	conf.PidFile = ""
	conf.Database = Database{
		Password:            "postgres",
		Username:            "postgres",
		Hostname:            "localhost",
		Database:            "test",
		Port:                5432,
		MaxConnections:      20,
		LogQueries:          false,
		Retries:             5,
		SleepBetweenRetries: "7s",
	}
}

// ReadConfigFile reading and parsing configuration yaml file
func (conf *Config) ReadConfigFile(configPath string) {
	if configPath == "" {
		configPath = DefaultConfigPath
	}
	yamlConfig, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yamlConfig, &conf)
	if err != nil {
		fmt.Errorf("could not unmarshal config %v", conf)
		log.Fatal(err)
	}
}
