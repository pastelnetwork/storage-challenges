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
	Version                    string          `mapstructure:"version"`
	Remoter                    *message.Config `mapstructure:"remoter"`
	Database                   *storage.Config `mapstructure:"database"`
	MasternodePastelID         string          `mapstructure:"masternode_id"`
	MasternodePastelPassphrase string          `mapstructure:"masternode_passphrase"`
	PastelClient               *pastel.Config  `mapstructure:"pastel_client,omitempty"`

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
