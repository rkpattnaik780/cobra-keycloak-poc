package root

import (
	"github.com/rkpattnaik780/cobra-keycloak-poc/pkg/cmd/login"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {

	rootCmd := &cobra.Command{
		SilenceUsage:  true,
		SilenceErrors: true,
		Use:           "ckp",
		Short:         "Root Command",
	}

	rootCmd.AddCommand(login.NewLoginCmd())

	return rootCmd

}
