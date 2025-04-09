package cli

import (
	auth "duckdb-test/app/internal/auth/cli"
	"fmt"

	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command

func newLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
		and usage of using your command. For example:

		Cobra is a CLI library for Go that empowers applications.
		This application is a tool to generate the needed files
		to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("logout called")

			err := auth.Logout()
			if err != nil {
				fmt.Println("❌ Failed to log out:", err)
			}

		},
	}
}
