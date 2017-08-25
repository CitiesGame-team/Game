package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
	DriverName  string `yaml:"driver_name"`
	ConnStr     string `yaml:"conn_string"`
	PoolMaxIdle int    `yaml:"pool_max_idle"`
	PoolMaxOpen int    `yaml:"pool_max_open"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type ProjectConfig struct {
	Name           string          `yaml:"name"`
	SplashFile     string          `yaml:"splash_file"`
	CitiesBaseFile string          `yaml:"cities_base_file"`
	Db             *DatabaseConfig `yaml:"db"`
	Server         *ServerConfig   `yaml:"server"`
}

func ReadYAML(fName string, dest interface{}) error {
	file, err := os.Open(fName)
	if err != nil {
		return fmt.Errorf("can't open yaml file %q: %s", fName, err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("can't read yaml file %q: %s", fName, err)
	}

	if err := yaml.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("can't yaml.unmarshal file %q: %s", fName, err)
	}

	return nil
}

func ReadProjectConfig(fName string) (ProjectConfig, error) {
	var dest ProjectConfig

	err := ReadYAML(fName, &dest)
	if err != nil {
		return dest, err
	}

	return dest, nil
}

type CitiesCollection struct {
	Cities []string `yaml:"cities"`
}

func ReadCitiesBase(fName string) (CitiesCollection, error) {
	var dest CitiesCollection

	err := ReadYAML(fName, &dest)
	if err != nil {
		return dest, err
	}

	return dest, nil
}
