package cmd

import (
	"fmt"
	"github.com/markosamuli/glassfactory/settings"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

func getSettingsFromUser() (*settings.GlassFactorySettings, error) {

	fmt.Print("Type your Glass Factory subdomain (eg. 'subdomain' if you login at https://subdomain.glassfactory.io):\n")

	var subdomain string
	if _, err := fmt.Scan(&subdomain); err != nil {
		return nil, fmt.Errorf("unable to read subdomain %v", err)
	}

	fmt.Print("Type your Glass Factory login email:\n")

	var email string
	if _, err := fmt.Scan(&email); err != nil {
		return nil, fmt.Errorf("unable to read email address %v", err)
	}

	fmt.Print("Type your Glass Factory API key:\n")

	var token string
	if _, err := fmt.Scan(&token); err != nil {
		return nil, fmt.Errorf("unable to read API key %v", err)
	}

	gfs := &settings.GlassFactorySettings{
		AccountSubdomain: subdomain,
		UserEmail: email,
		UserToken: token,
	}
	return gfs, nil
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure",
	Long:  `First time configuration`,
	Run: func(cmd *cobra.Command, args []string) {

		gfs, err := getSettingsFromUser()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		if cfgFile == "" {
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			cfgFile = filepath.Join(home, fmt.Sprintf("%s.yaml", cfgName))
			viper.SetConfigType("yaml")
			viper.SetConfigFile(cfgFile)
		}

		viper.Set("account_subdomain", gfs.AccountSubdomain)
		viper.Set("user_email", gfs.UserEmail)
		viper.Set("user_token", gfs.UserToken)
		err = viper.SafeWriteConfig()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		fmt.Println("Configuration saved in:", cfgFile)
	},
}