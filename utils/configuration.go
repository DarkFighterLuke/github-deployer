package utils

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Server struct {
		Host            string `yaml:"host"`
		Port            string `yaml:"port"`
		PayloadEndpoint string `yaml:"payload_endpoint"`
	} `yaml:"server"`
	Repository struct {
		Branch string `yaml:"branch"`
	} `yaml:"repository"`
	Script struct {
		Path string `yaml:"path"`
	} `yaml:"script"`
}

//GetConfig Given a configuration path validate and return new Config
func GetConfig(configPath string) (*Config, error) {
	err := ValidateConfigPath(configPath)
	if err != nil {
		log.Fatalln(err)
	}
	conf, err := NewConfig(configPath)
	return conf, err
}

//NewConfig Given a configuration path read a new Config
func NewConfig(congigPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(congigPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}

//ValidateConfigPath Check if path is valid or not
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}
