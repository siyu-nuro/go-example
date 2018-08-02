package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
	Addr string `yaml:"addr"`
	Name string `yaml:"name"`
}

type HTTPServerConfig struct {
	Port string `yaml:"port"`
}

type Config struct {
	Database              DatabaseConfig              `yaml:"database"`
	HTTPServer            HTTPServerConfig            `yaml:"httpServer"`
}

var (
	configSingleton *Config
	mutex           sync.Mutex
)

func newConfigFromFile(file *string) (*Config, error) {
	cfg := &Config{}

	if file == nil {
		fmt.Errorf("failed to load config - file is nil")
		return nil, errors.New("config file path is nil")
	}

	data, err := ioutil.ReadFile(*file)
	if err != nil {
		fmt.Errorf("failed to load config - failed to read file. error %v", err)
		return nil, err
	}

	if err = yaml.UnmarshalStrict(data, cfg); err != nil {
		fmt.Errorf("failed to load config - failed to unmarshal yaml file. error %v", err)
		return nil, err
	}

	return cfg, nil
}

// GetConfig returns config singleton
func GetConfig(file *string) (*Config, error) {
	if configSingleton == nil && file == nil {
		fmt.Errorf("cannot load config from file %v", file)
		return nil, errors.New("config is not initialized and config file path is nil")
	}
	if configSingleton == nil {
		mutex.Lock()
		if configSingleton == nil {
			cfg, err := newConfigFromFile(file)
			if err != nil {
				fmt.Errorf("failed to load config from file %v", file)
				return nil, err
			}
			configSingleton = cfg
		}
		mutex.Unlock()
	}
	return configSingleton, nil
}
