package cli

import (
	auth "duckdb-test/app/internal/auth/cli"
	"duckdb-test/app/internal/client"
	"duckdb-test/app/internal/config"
	"fmt"

	"github.com/spf13/cobra"
)

func newQueryCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "query",
		Short: "Query using duckdb",
		Long:  `This command will send a query to the DuckDB server and return the results.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			tokenStr, err := auth.CheckAuth(cfg)
			if err != nil {
				fmt.Println("❌ Not authenticated. Please run 'cli login'")
				return
			}

			QueryParam := client.QueryParam{
				Key:   "q",
				Value: args[0]}

			// err = client.DoAndPrintRequest(cfg, *tokenStr, "/query", param)

			err = client.DoAndPrintRequest(client.RequestParams{
				Cfg:    cfg,
				Token:  *tokenStr,
				Path:   "/query",
				Method: "GET",
				QueryParams: []client.QueryParam{
					QueryParam,
				},
			},
			)

			if err != nil {
				fmt.Println(err)
			}

		},
	}

}
