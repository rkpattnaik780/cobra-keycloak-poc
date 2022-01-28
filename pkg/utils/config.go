package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type IConfig struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func getConfigLocation() (string, error) {
	cfgPath, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cfgPath, "ckp"), nil
}

func InitConfig() error {
	ckpCfgDir, err := getConfigLocation()
	if err != nil {
		return err
	}

	// create ckp config directory
	if _, err = os.Stat(ckpCfgDir); os.IsNotExist(err) {
		err = os.MkdirAll(ckpCfgDir, 0o700)
		if err != nil {
			return err
		}
	}

	return err
}

func SaveConfig(cfg *IConfig) error {
	cfgDir, err := getConfigLocation()
	if err != nil {
		return err
	}

	file := fmt.Sprintf("%s/%s", cfgDir, "config.json")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("%v: %w", "unable to marshal config", err)
	}

	if _, err = os.Stat(file); os.IsNotExist(err) {
		err = os.Mkdir(file, 0o700)
		if err != nil {
			return err
		}
	}
	err = ioutil.WriteFile(file, data, 0o600)
	if err != nil {
		return fmt.Errorf("%v: %w", "unable to save config", err)
	}
	return nil
}

func SaveToken(token string) error {

	cfgDir, err := getConfigLocation()
	if err != nil {
		return err
	}

	file := fmt.Sprintf("%s/%s", cfgDir, "config.json")

	var cfg *IConfig = &IConfig{
		AccessToken: token,
	}

	data, err := json.MarshalIndent(&cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("%v: %w", "unable to marshal config", err)
	}

	if _, err = os.Stat(cfgDir); os.IsNotExist(err) {
		err = os.Mkdir(cfgDir, 0o700)
		if err != nil {
			return err
		}
	}
	err = ioutil.WriteFile(file, data, 0o600)
	if err != nil {
		return fmt.Errorf("%v: %w", "unable to save config", err)
	}
	return nil

}
