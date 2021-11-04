package config

import (
	"fmt"
	"os"

	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/storage-challenges/external/message"
	"github.com/pastelnetwork/storage-challenges/external/storage"
	"github.com/spf13/viper"
)

type Config struct {
	Version                    string          `yaml:"version"`
	Remoter                    *message.Config `yaml:"remoter,omitempty"`
	MasternodePastelID         string          `yaml:"masternode_id"`
	MasternodePastelPassphrase string          `yaml:"masternode_passphrase"`
	Database                   *storage.Config `yaml:"database,omitempty"`
	PastelClient               *pastel.Config  `yaml:"pastel_client,omitempty"`

	isLoaded bool
}

var configPath string = "./config"
var configENVPrefix string = "STORAGE_CHALLENGE_CONFIG"

func init() {
	val := os.Getenv(configENVPrefix)
	fmt.Println("STORAGE_CHALLENGE_CONFIG", val)
	if val != "" {
		configPath = val
	}
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	// if env is setted, prefer to use env config data other than config file
	// eg. if env variable STORAGE_CHALLENGE_CONFIG_DATA_PARENT_DATA_CHILDREN is setted, use that value instea of yaml config data_parent.data_children
	viper.SetEnvPrefix(configENVPrefix)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("could not read storage challenge config file: %v", err))
	}
}

func (c *Config) Load() error {
	if c.isLoaded {
		return nil
	}

	var err error
	if err = viper.Unmarshal(c); err == nil {
		c.isLoaded = true
	}
	return err
}
