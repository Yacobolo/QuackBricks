package cli

import (
	auth "duckdb-test/app/internal/auth/cli"
	"duckdb-test/app/internal/client"
	"duckdb-test/app/internal/config"
	"fmt"

	"github.com/spf13/cobra"
)

func newCatalogGetCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "get",
		Args:    cobra.ExactArgs(1),
		Short:   "Get a catalog",
		Long:    `Get a catalog in the DuckDB service.`,
		Example: `cli catalog get <catalog_name> or cli catalog get <catalog_id>`,
		Run: func(cmd *cobra.Command, args []string) {

			tokenStr, err := auth.GetAuthToken(cfg)

			if err != nil {
				fmt.Println("❌ Failed to get token:", err)
				return
			}

			err = client.DoAndPrintRequest(client.RequestParams{
				Cfg:        cfg,
				Token:      *tokenStr,
				Path:       "/me",
				HttpMethod: client.MethodGet,
			},
			)

			if err != nil {
				fmt.Println(err)
			}

		},
	}

}
