package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	GRPC struct {
		Target  string        `mapstructure:"target"`
		Timeout time.Duration `mapstructure:"timeout"`
	} `mapstructure:"grpc"`

	Postgres struct {
		Host        string        `mapstructure:"host"`
		Port        string        `mapstructure:"port"`
		User        string        `mapstructure:"user"`
		Password    string        `mapstructure:"password"`
		DBName      string        `mapstructure:"dbname"`
		SSLMode     string        `mapstructure:"sslmode"`
		MaxConns    int           `mapstructure:"max_connections"`
		MaxIdleTime time.Duration `mapstructure:"max_idle_time"`
	} `mapstructure:"postgres"`

	Redis struct {
		Host        string        `mapstructure:"host"`
		Port        string        `mapstructure:"port"`
		DB          int           `mapstructure:"db"`
		Password    string        `mapstructure:"password"`
		Timeout     time.Duration `mapstructure:"timeout"`
		CacheTTL    time.Duration `mapstructure:"cache_ttl"`
		MaxConns    int           `mapstructure:"max_connections"`
		MaxIdleTime time.Duration `mapstructure:"max_idle_time"`
	} `mapstructure:"redis"`

	Logger struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"logger"`
}

func LoadConfig() *Config {

	viper.SetConfigName("config")    // Имя файла без расширения
	viper.SetConfigType("yaml")      // Формат файла
	viper.AddConfigPath(".")         // Ищем в текущей директории
	viper.AddConfigPath("./configs") // Или в папке configs

	// Читаем конфигурационный файл
	if err := viper.ReadInConfig(); err != nil {
		panic("failed to read config: " + err.Error())
	}
	// Создаем субвипер для извлечения только db_service
	subv := viper.Sub("db_service")
	if subv == nil {
		panic("missing 'db_service' section in config")
	}
	// Распарсим конфиг в структуру
	var cfg Config
	if err := subv.Unmarshal(&cfg); err != nil {
		panic("failed to unmarshal config: " + err.Error())
	}

	return &cfg
}
