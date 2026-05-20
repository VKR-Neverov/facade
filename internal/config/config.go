package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

const (
	PgDsn          = "pg-dsn"
	JaegerEndpoint = "jaeger-endpoint"
	ServiceName    = "service-name"
)

var conf *Config

type Config struct {
	Secrets struct {
		PgDsn          string `yaml:"pg-dsn"`
		JaegerEndpoint string `yaml:"jaeger-endpoint"`
		ServiceName    string `yaml:"service-name"`
	} `yaml:"secrets"`
}

func Init() error {
	body, err := os.ReadFile("./configs/values.yaml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(body, &conf)
	return err
}

func GetValue(key string) (string, error) {
	switch key {
	case PgDsn:
		return conf.Secrets.PgDsn, nil
	case JaegerEndpoint:
		return conf.Secrets.JaegerEndpoint, nil
	case ServiceName:
		return conf.Secrets.ServiceName, nil
	default:
		return "", ErrConfigNotFoundByKey(key)
	}
	//return getValue(reflect.ValueOf(conf).Elem(), key)
}

func InitTestConfigs() error {
	body, err := os.ReadFile("./../../../../configs/values_test.yaml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(body, &conf)
	return err
}
