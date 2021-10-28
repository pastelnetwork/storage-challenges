package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	FolderPathContainingSampleRaptorqSymbolFiles string   `yaml:"folder_path_containing_sample_raptorq_symbol_files"`
	RqsymbolFileStorageDataFolderPath            string   `yaml:"rqsymbol_file_storage_data_folder_path"`
	NewRqsymbolFileStorageDataFolderPath         string   `yaml:"new_rqsymbol_file_storage_data_folder_path"`
	MaxSecondsToRespondToStorageChallenge        int      `yaml:"max_seconds_to_respond_to_storage_challenge"`
	NumberOfChallengeReplicas                    int      `yaml:"number_of_challenge_replicas"`
	SliceOfPastelMasternodeIds                   []string `yaml:"slice_of_pastel_masternode_ids"`
	SliceOfNewPastelMasternodeIds                []string `yaml:"slice_of_new_pastel_masternode_ids"`
	isLoaded                                     bool
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
