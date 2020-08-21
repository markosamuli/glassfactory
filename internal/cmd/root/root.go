package root

import (
	"context"
	"fmt"
	"os"

	"github.com/markosamuli/glassfactory/internal/auth"
	authCmd "github.com/markosamuli/glassfactory/internal/cmd/auth"
	"github.com/markosamuli/glassfactory/internal/cmd/report"
	"github.com/markosamuli/glassfactory/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ( // Used for flags.
	cfgFile string
	verbose bool
	gfAuth  = auth.NewAuth()
	rootCmd = &cobra.Command{
		Use:   "glassfactory",
		Short: "Glass Factory reports tool",
		Long:  `CLI reports tool for Glass Factory.`,
	}
)

// Execute executes the root command.
func Execute() {
	ctx := context.Background()
	ctx = auth.NewContext(ctx, gfAuth)
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.glassfactory.yaml)")
	rootCmd.PersistentFlags().StringVar(&gfAuth.Account, "account", "", "Glass Factory account subdomain")
	rootCmd.PersistentFlags().StringVar(&gfAuth.Email, "email", "", "Glass Factory user email address")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	viper.BindPFlag("account", rootCmd.PersistentFlags().Lookup("account"))
	viper.BindPFlag("email", rootCmd.PersistentFlags().Lookup("email"))

	rootCmd.AddCommand(authCmd.NewCommand())
	rootCmd.AddCommand(report.NewCommand())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	config.InitConfig(cfgFile, gfAuth)
	if verbose {
		if cfgFile := config.GetConfigFile(); cfgFile != "" {
			fmt.Println("Using config file:", cfgFile)
		}
	}
}
