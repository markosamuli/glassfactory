package config

import (
	"fmt"
	"path/filepath"

	"github.com/markosamuli/glassfactory/internal/auth"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const cfgName = ".glassfactory"
const envPrefix = "gf"

// InitConfig reads in config file and ENV variables if set.
func InitConfig(cfgFile string, gfAuth *auth.Auth) error {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		// Search config in home directory with name ".glassfactory" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(cfgName)
	}

	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if gfAuth.Account == "" {
			if account := viper.GetString("account"); account != "" {
				gfAuth.Account = account
			}
		}
		if gfAuth.Email == "" {
			if email := viper.GetString("email"); email != "" {
				gfAuth.Email = email
			}
		}
	}
	return nil
}

// GetConfigFile returns the current config file
func GetConfigFile() string {
	return viper.ConfigFileUsed()
}

// SaveConfig writes authentication into Viper config file
func SaveConfig(gfAuth *auth.Auth) error {
	cfgFile := viper.ConfigFileUsed()
	if cfgFile == "" {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		cfgFile = filepath.Join(home, fmt.Sprintf("%s.yaml", cfgName))
		viper.SetConfigType("yaml")
		viper.SetConfigFile(cfgFile)
	}
	viper.Set("account", gfAuth.Account)
	viper.Set("email", gfAuth.Email)
	if err := viper.SafeWriteConfig(); err != nil {
		return err
	}
	return nil
}
