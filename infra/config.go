package infra

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	Server     Server     `mapstructure:"server"`
	Postgres   Postgres   `mapstructure:"postgres"`
	Redis      Redis      `mapstructure:"redis"`
	Shortener  Shortener  `mapstructure:"shortener"`
	WorkerPool WorkerPool `mapstructure:"worker_pool"`
	Telemetry  Telemetry  `mapstructure:"telemetry"`
}

type Server struct {
	AppVersion string `mapstructure:"app_version"`
	Address    string `mapstructure:"address"`
	Port       string `mapstructure:"port"`
	LogLevel   string `mapstructure:"log_level"`
}

type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db_name"`
	LogLevel string `mapstructure:"log_level"`
}

type Redis struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
}

type Shortener struct {
	CodeLength int `mapstructure:"code_length"`
}

type WorkerPool struct {
	WorkerCount int `mapstructure:"worker_count"`
	QueueSize   int `mapstructure:"queue_size"`
}

type Telemetry struct {
	ServiceNamespaceKey string `mapstructure:"service_namespace_key"`
	ServiceNameKey      string `mapstructure:"service_name_key"`
	Trace               Trace  `mapstructure:"trace"`
	Metric              Meter  `mapstructure:"metric"`
}

type Trace struct {
	Enabled    bool   `mapstructure:"enabled"`
	JaegerHost string `mapstructure:"jaeger_host"`
	JaegerPort string `mapstructure:"jaeger_port"`
}

type Meter struct {
	Enabled bool   `mapstructure:"enabled"`
	Address string `mapstructure:"address"`
}

func Load() (config *Config, err error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	if err = v.ReadInConfig(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return nil, errors.New("config file not found")
		}

		return nil, err
	}

	if err = v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return config, nil
}
