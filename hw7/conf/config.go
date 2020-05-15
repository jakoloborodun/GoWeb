package conf

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// Config struct for webapp config
type Config struct {
	Server struct {
		Name    string `yaml:"name"`
		RootDir string `yaml:"root_dir"`
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Name     string `yaml:"db_name"`
		Username string `yaml:"db_user"`
		Password string `yaml:"db_pass"`
		Host     string `yaml:"db_host"`
		Port     string `yaml:"db_port"`
	} `yaml:"database"`
	Logger struct {
		Level        string `yaml:"level"`
		ReportCaller bool   `yaml:"report_caller"`
		Output       string `yaml:"output"`
		Syslog       bool   `yaml:"syslog"`
	} `yaml:"logger"`
}

// Will parse config flag and returns the path to config file
// or error if file doesn't exist.
func ParseConfigFlag() (string, error) {
	var configPath string

	flag.StringVar(&configPath, "c", "./config.yml", "path to yaml config file")
	flag.Parse()

	_, err := os.Stat(configPath)
	if err != nil {
		return "", err
	}

	return configPath, nil
}

// NewConfig returns a new decoded Config struct.
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, err
}
