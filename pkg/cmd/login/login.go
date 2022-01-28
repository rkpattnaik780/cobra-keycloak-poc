package login

import (
	"log"

	"github.com/rkpattnaik780/cobra-keycloak-poc/pkg/utils"
	"github.com/spf13/cobra"
)

func NewLoginCmd() *cobra.Command {

	loginCmd := &cobra.Command{
		SilenceUsage:  true,
		SilenceErrors: true,
		Use:           "login",
		Short:         "Command to  login using keycloak",
		Run: func(cmd *cobra.Command, args []string) {
			handleLoginCallback()
		},
	}

	return loginCmd

}

func handleLoginCallback() {

	utils.CloseApp.Add(1)
	config := utils.Config{
		KeycloakConfig: utils.KeycloakConfig{
			KeycloakURL: "http://localhost:8080/auth",
			Realm:       "demo",
			ClientID:    "cli-client",
		},
		EmbeddedServerConfig: utils.EmbeddedServerConfig{
			Port:         3000,
			CallbackPath: "sso-callback",
		},
	}

	utils.StartServer(config)
	err := utils.OpenBrowser(utils.BuildAuthorizationRequest(config))
	if err != nil {
		log.Fatalf("Could not open the browser for url %v", utils.BuildAuthorizationRequest(config))
	}

	utils.CloseApp.Wait()
}
