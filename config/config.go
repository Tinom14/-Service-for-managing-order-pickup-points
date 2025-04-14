package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

const (
	Secret = "secret-key"
)

type HTTPConfig struct {
	Address string `yaml:"address"`
}

type PrometheusConfig struct {
	Port uint16 `yaml:"port"`
}

type Postgres struct {
	Host          string `yaml:"host"`
	Port          uint   `yaml:"port"`
	User          string `yaml:"user"`
	Password      string `yaml:"password"`
	DBName        string `yaml:"dbname"`
	SSLMode       string `yaml:"sslmode"`
	MigrationPath string `yaml:"migrationPath"`
}

type AppConfig struct {
	HTTPConfig       `yaml:"http"`
	Postgres         `yaml:"postgres"`
	PrometheusConfig `yaml:"prometheus"`
}

type AppFlags struct {
	ConfigPath string `yaml:"config_path"`
}

func ParseFlags() AppFlags {
	configPath := flag.String("config", "", "Path to config")
	flag.Parse()
	return AppFlags{
		ConfigPath: *configPath,
	}
}
func MustLoad(cfgPath string, cfg any) {
	if cfgPath == "" {
		log.Fatal("Config path is not set")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist by this path: %s", cfgPath)
	}

	if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
		log.Fatalf("error reading config: %s", err)
	}
}
