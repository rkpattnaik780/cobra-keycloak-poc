package main

import (
	"sync"

	"github.com/rkpattnaik780/cobra-keycloak-poc/pkg/cmd/root"
	"github.com/rkpattnaik780/cobra-keycloak-poc/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	CloseApp sync.WaitGroup
	RootCmd  *cobra.Command
)

func init() {
	utils.InitConfig()
	RootCmd = root.NewRootCmd()
}

func main() {
	RootCmd.Execute()
}
