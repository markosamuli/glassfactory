package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/markosamuli/glassfactory/internal/auth"
	"github.com/markosamuli/glassfactory/internal/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// LoginOptions for the login command
type LoginOptions struct{}

// NewLoginCommand creates new command
func NewLoginCommand() *cobra.Command {
	o := &LoginOptions{}
	c := &cobra.Command{
		Use:   "login",
		Short: "Configure Glass Factory login details.",
		Long: `Configure Glass Factory login details.

	This command will store your account information and email address in the
	local config file and your authentication token in the system keychain.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := o.Run(cmd)
			if err != nil {
				fmt.Print(err)
			}
		},
	}
	return c
}

// Run the command
func (o *LoginOptions) Run(cmd *cobra.Command) error {
	gfAuth, ok := auth.FromContext(cmd.Context())
	if !ok {
		return fmt.Errorf("failed to get authentication details")
	}
	storeLoginDetailsInKeyring := false
	if gfAuth.Account == "" {
		fmt.Print("Type your Glass Factory subdomain (eg. 'subdomain' if you login at https://subdomain.glassfactory.io):\n")
		reader := bufio.NewReader(os.Stdin)
		account, err := reader.ReadString('\n')
		if err != nil {
			return errors.Wrapf(err, "unable to read subdomain")
		}
		gfAuth.Account = strings.TrimSpace(account)
	}
	if gfAuth.Email == "" {
		fmt.Print("Type your Glass Factory login email:\n")
		reader := bufio.NewReader(os.Stdin)
		email, err := reader.ReadString('\n')
		if err != nil {
			return errors.Wrapf(err, "unable to read email address")
		}
		gfAuth.Email = strings.TrimSpace(email)
	}
	if gfAuth.APIKey == "" {
		fmt.Print("Type your Glass Factory API key:\n")
		byteAPIKey, err := terminal.ReadPassword(0)
		if err != nil {
			return errors.Wrapf(err, "unable to read API key")
		}
		apiKey := string(byteAPIKey)
		fmt.Printf("%s\n", strings.Repeat("*", len(apiKey)))
		gfAuth.APIKey = strings.TrimSpace(apiKey)
		storeLoginDetailsInKeyring = true
	}
	if err := gfAuth.Setup(); err != nil {
		fmt.Println("Deleting login details from keyring")
		err := gfAuth.DeleteLoginDetailsFromKeyring()
		if err != nil {
			return err
		}
	}
	if storeLoginDetailsInKeyring {
		fmt.Println("Storing password in keyring")
		err := gfAuth.StoreLoginDetailsInKeyring()
		if err != nil {
			return err
		}
	}

	err := config.SaveConfig(gfAuth)
	if err != nil {
		return errors.Wrapf(err, "couldn't save configuration")
	}
	if cfgFile := config.GetConfigFile(); cfgFile != "" {
		fmt.Println("Configuration saved in:", cfgFile)
	}

	return nil
}
