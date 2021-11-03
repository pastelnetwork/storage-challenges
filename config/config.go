package config

import (
	"os"

	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/storage-challenges/external/message"
	"github.com/spf13/viper"
)

type Config struct {
	Remoter                   message.Config `yaml:"remoter"`
	MasterNodeID              string         `yaml:"masternode_id"`
	NumberOfChallengeReplicas int            `yaml:"number_of_challenge_replicas"`
	Pastel                    pastel.Config  `yaml:"pastel_config"`

	isLoaded bool
}

var configPath string = "./config"
var configENVPrefix string = "STORAGE_CHALLENGE_CONFIG"

func init() {
	val := os.Getenv(configENVPrefix)
	if val != "" {
		configPath = val
	}
}

func (c *Config) Load() error {
	if c.isLoaded {
		return nil
	}
	viper.AddConfigPath(configPath)
	// if env is setted, prefer to use env config data other than config file
	// eg. if env variable STORAGE_CHALLENGE_CONFIG_DATA_PARENT_DATA_CHILDREN is setted, use that value instea of yaml config data_parent.data_children
	viper.SetEnvPrefix(configENVPrefix)
	viper.AutomaticEnv()

	var err error
	if err = viper.Unmarshal(c); err == nil {
		c.isLoaded = true
	}
	return err
}
