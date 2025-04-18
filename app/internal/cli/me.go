package cli

import (
	auth "duckdb-test/app/internal/auth/cli"
	"duckdb-test/app/internal/client"
	"duckdb-test/app/internal/config"
	"fmt"

	"github.com/spf13/cobra"
)

func newMeCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "me",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
		and usage of using your command. For example:

		Cobra is a CLI library for Go that empowers applications.
		This application is a tool to generate the needed files
		to quickly create a Cobra application.`,
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
