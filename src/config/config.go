package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type KafkaConfig struct {
	Addr        string `yaml:"address"`
	Topic       string `yaml:"topic"`
	GroupId     string `yaml:"group_id"`
	OffsetReset string `yaml:"auto_offset_reset"`
}

type WsConfig struct {
	Addr            string `yaml:"address"`
	ReadBufferSize  int    `yaml:"read-buffer-size"`
	WriteBufferSize int    `yaml:"write-buffer-size"`
}

type RedisConfig struct {
	Addr     string `yaml:"address"`
	Password string `yaml:"password"`
	Db       int    `yaml:"Db"`

	TrackingKey            string `yaml:"tracking-key"`
	ActiveDriversKeyFormat string `yaml:"drivers-key-format"`
	TripsKeyFormat         string `yaml:"trips-key-format"`
}

type Config struct {
	WsConfig    WsConfig    `yaml:"ws-server"`
	RedisConfig RedisConfig `yaml:"redis-cache"`

	SystemBrokerConfig struct {
		KafkaBroker KafkaConfig `yaml:"kafka"`
	} `yaml:"system-broker"`

	TimingBrokerConfig struct {
		KafkaBroker KafkaConfig `yaml:"kafka"`
	} `yaml:"timing-broker"`

	TrackingBrokerConfig struct {
		KafkaBroker KafkaConfig `yaml:"kafka"`
	} `yaml:"tracking-broker"`
}

type ConfigErr struct {
	Err error
}

func (err *ConfigErr) Error() string {
	return err.Err.Error()
}

func GetConfig() (*Config, *ConfigErr) {
	config_content, err := os.ReadFile("config.yaml")
	content := &Config{}

	if err != nil {
		return nil, &ConfigErr{
			Err: err,
		}
	}

	if err := yaml.Unmarshal(config_content, content); err != nil {
		return nil, &ConfigErr{
			Err: err,
		}
	}
	return content, nil
}
