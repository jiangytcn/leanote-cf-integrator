package mci

import (
	"os"
	"errors"
	toml "github.com/robfig/config"
	"github.com/cloudfoundry-community/go-cfenv"
)

func PushConfig(config *toml.Config, configFilePath string) error {
	config.WriteFile(configFilePath, 0644, "cf integrator generated this file")
	return nil
}
func CFConfig(config *toml.Config) error {
	var err error
	if !IsInCloudFoundry() {
		return errors.New("Not in Cloud Foundry environment.")
	}
	c := toml.NewDefault()
	c.AddOption("", "http.port", os.Getenv("PORT"))
	c.AddOption("", "site.url", "http://0.0.0.0:" + os.Getenv("PORT"))
	config.Merge(c)

	appEnv, err := cfenv.Current()
	if err != nil {
		return err
	}

	err = cfDatabase(appEnv, config)
	if err != nil {
		return err
	}

	return nil
}

func IsInCloudFoundry() bool {
	d := os.Getenv("VCAP_APPLICATION")
	if d != "" {
		return true
	}
	return false
}
func ExtractConfig(configFilePath string) (*toml.Config, error) {
	_, err := os.Stat(configFilePath)
	if err != nil {
		return nil, err
	}

	config, err := toml.ReadDefault(configFilePath)

	return config, err
}