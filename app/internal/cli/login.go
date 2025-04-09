package cli

import (
	auth "duckdb-test/app/internal/auth/cli"
	"duckdb-test/app/internal/config"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var ()

// loginCmd represents the login command

func newLoginCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
		and usage of using your command. For example:

		Cobra is a CLI library for Go that empowers applications.
		This application is a tool to generate the needed files
		to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("🔐 Logging in...")
			err := auth.AuthenticateAndSave(cfg)
			if err != nil {
				fmt.Println("❌ Failed to log in:", err)
				os.Exit(1)
			}
		},
	}

}
