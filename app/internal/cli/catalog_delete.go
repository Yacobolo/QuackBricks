package cli

import (
	auth "duckdb-test/app/internal/auth/cli"
	"duckdb-test/app/internal/client"
	"duckdb-test/app/internal/config"
	"fmt"

	"github.com/spf13/cobra"
)

var ()

func newCatalogDeleteCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <catalog_name>",
		Short: "Delete a catalog",
		Long:  "Delete a catalog in the DuckDB service.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			catalogName := args[0]

			tokenStr, err := auth.GetAuthToken(cfg)

			if err != nil {
				return err
			}

			path := fmt.Sprintf("/catalogs/%s", catalogName)

			err = client.DoAndPrintRequest(client.RequestParams{
				Cfg:        cfg,
				Token:      *tokenStr,
				Path:       path,
				HttpMethod: client.MethodDelete,
			},
			)

			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("Catalog '%s' deletion initiated.\n", catalogName)
			return nil
		},
	}

	return cmd
}
