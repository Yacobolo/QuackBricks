package cli

import (
	auth "duckdb-test/app/internal/auth/cli"
	"duckdb-test/app/internal/client"
	"duckdb-test/app/internal/config"
	"fmt"

	"github.com/spf13/cobra"
)

var ()

func newCatalogListCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all catalog entries.",
		Long:  `Lists all catalog entries in the system.`,
		Run: func(cmd *cobra.Command, args []string) {

			tokenStr, err := auth.GetAuthToken(cfg)
			if err != nil {
				fmt.Println("❌ Not authenticated. Please run 'cli login'")
				return
			}

			err = client.DoAndPrintRequest(client.RequestParams{
				Cfg:        cfg,
				Token:      *tokenStr,
				Path:       "/catalog",
				HttpMethod: client.MethodGet,
			},
			)
			if err != nil {
				fmt.Println(err)
			}

		},
	}
	return cmd

}
