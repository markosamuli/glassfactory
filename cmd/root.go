package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/markosamuli/glassfactory"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const cfgName = ".glassfactory"
const envPrefix = "gf"

var ( // Used for flags.
	cfgFile     string
	accountSubdomain string
	userEmail string
	userToken string
	rootCmd = &cobra.Command{
		Use:   "glassfactory",
		Short: "Glass Factory reports tool",
		Long: `CLI reports tool for Glass Factory.`,
	}
)

func createApiService() (*glassfactory.Service, error) {
	ctx := context.Background()
	gfs := glassfactory.Settings{}
	gfs.AccountSubdomain = viper.GetString("account_subdomain")
	gfs.UserEmail = viper.GetString("user_email")
	gfs.UserToken = viper.GetString("user_token")
	s, err := glassfactory.NewService(ctx, &gfs)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.glassfactory.yaml)")
	rootCmd.PersistentFlags().StringVarP(&accountSubdomain, "subdomain", "d", "", "account subdomain")
	rootCmd.PersistentFlags().StringVarP(&userEmail, "email", "e", "", "user email address")
	rootCmd.PersistentFlags().StringVarP(&userToken, "token", "t", "", "user API key")
	viper.BindPFlag("account_subdomain", rootCmd.PersistentFlags().Lookup("subdomain"))
	viper.BindPFlag("user_email", rootCmd.PersistentFlags().Lookup("email"))
	viper.BindPFlag("user_token", rootCmd.PersistentFlags().Lookup("token"))
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		// Search config in home directory with name ".glassfactory" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(cfgName)
	}

	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

}