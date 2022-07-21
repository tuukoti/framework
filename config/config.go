package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host  string `yaml:"host"`
	Debug bool   `yaml:"debug"`
}

// Load simply loads the config from location "path".
func Load(path string) (*Config, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, &ErrMissingConfigFile{
			Path: path,
			Err:  err,
		}
	}

	cfg := &Config{}

	err = yaml.Unmarshal(f, cfg)
	if err != nil {
		return nil, &ErrInvalidConfig{
			Path: path,
			Err:  err,
		}
	}

	return nil, nil
}

type ErrMissingConfigFile struct {
	Path string
	Err  error
}

func (e *ErrMissingConfigFile) Error() string {
	return fmt.Sprintf("unable to location config file: %s, | err: %s", e.Path, e.Err.Error())
}

type ErrInvalidConfig struct {
	Path string
	Err  error
}

func (e *ErrInvalidConfig) Error() string {
	return fmt.Sprintf("invalid config file: %s, | err: %s", e.Path, e.Err.Error())
}
