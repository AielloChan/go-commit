package config

import (
	"errors"
	"os"
	"path/filepath"
)

const CONFIG_FILE_NAME = "commit.config.json"

func FindStageIndexByName(stages *[]Stage, name string) (int, bool) {
	if name == "" {
		return -1, false
	}
	for i, item := range *stages {
		if item.Name == name {
			return i, true
		}
	}
	return -1, false
}

func InitConfig() (Configuration, error) {
	curPath, _ := os.Getwd()
	_, err := os.Stat(filepath.Join(curPath, CONFIG_FILE_NAME))
	if err != nil {
		return Configuration{}, errors.New("Config 'commit.config.json' not exist at current path '" + curPath + "' or do not have access permission")
	}

	cfg, err := GetConfig(CONFIG_FILE_NAME)
	if err != nil {
		return Configuration{}, errors.New("Error while read config: " + err.Error())
	}
	if len(cfg.Stages) == 0 {
		return cfg, errors.New("Does not have 'stages' config")
	}
	return cfg, nil
}
