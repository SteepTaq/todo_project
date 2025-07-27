package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HTTP struct {
		Port    string        `mapstructure:"port"`
		Timeout time.Duration `mapstructure:"timeout"`
	} `mapstructure:"http"`

	GRPC struct {
		Target  string        `mapstructure:"target"`
		Timeout time.Duration `mapstructure:"timeout"`
	} `mapstructure:"grpc_db_service"`

	Kafka struct {
		Brokers []string `mapstructure:"brokers"`
		Topic   string   `mapstructure:"topic"`
	} `mapstructure:"kafka"`

	Logger struct {
		Level string `mapstructure:"level"` 
	} `mapstructure:"logger"`
}

func LoadConfig() *Config {

	viper.SetConfigName("config")   
	viper.SetConfigType("yaml")      
	viper.AddConfigPath(".")         
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		panic("failed to read config: " + err.Error())
	}
	subv := viper.Sub("api_service")
	if subv == nil {
		panic("missing 'api_service' section in config")
	}
	var cfg Config
	if err := subv.Unmarshal(&cfg); err != nil {
		panic("failed to unmarshal config: " + err.Error())
	}

	return &cfg
}
