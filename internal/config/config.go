package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Api       ApiConfig       `mapstructure:"api"`
	DBService DBServiceConfig `mapstructure:"db_service"`
	Kafka     KafkaConfig     `mapstructure:"kafka"`
	Redis     RedisConfig     `mapstructure:"redis"`
}

type ApiConfig struct {
	HTTP HTTPConfig       `mapstructure:"http"`
	GRPC GRPCClientConfig `mapstructure:"grpc"`
}

type DBServiceConfig struct {
	GRPC     GRPCServerConfig `mapstructure:"grpc"`
	Postgres PostgresConfig   `mapstructure:"postgres"`
	Redis    RedisConfig      `mapstructure:"redis"`
}

type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}
type HTTPConfig struct {
	Host    string `mapstructure:"host"`
	Port    string `mapstructure:"port"`
	Timeout string `mapstructure:"timeout"`
}

type GRPCClientConfig struct {
	Host    string `mapstructure:"host"`
	Port    string `mapstructure:"port"`
	Timeout string `mapstructure:"timeout"`
}

type GRPCServerConfig struct {
	Host    string `mapstructure:"host"`
	Port    string `mapstructure:"port"`
	Timeout string `mapstructure:"timeout"`
}

func Load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	return &cfg
}
