package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

var (
	cfg Config
)

// Config ...
type Config struct {
	Kafka struct {
		KafkaTopic string `yaml:"kafka.topic"`
	}
	RelativePath string `envconfig:"DATASOURCEPATH"`
}

func init() {
	readYAML(&cfg)
	readEnv(&cfg)
}

func readYAML(cfg *Config) {
	f, _ := os.Open("config/config.yml")

	defer f.Close()
	decoder := yaml.NewDecoder(f)
	decoder.Decode(cfg)
}

func readEnv(cfg *Config) {
	var fileName string
	fileName = "env/" + os.Getenv("ENV") + ".env"
	fmt.Println(fileName)
	godotenv.Load(fileName)
	envconfig.Process("", cfg)
}
