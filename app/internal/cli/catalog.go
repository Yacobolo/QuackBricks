package cli

import (
	"duckdb-test/app/internal/config"

	"github.com/spf13/cobra"
)

// Function to create the main 'catalog' command
func newCatalogCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "catalog",
		Short: "Manage catalog entries",
		Long:  `Provides commands to register, list, update, and delete catalog entries.`,
		// Run: func(cmd *cobra.Command, args []string) {
		//   // Optional: Display help if no subcommand is given
		// 	 cmd.Help()
		// },
	}

	// Add subcommands
	cmd.AddCommand(newCatalogRegisterCmd(cfg)) // Defined in catalog_register.go
	cmd.AddCommand(newCatalogListCmd(cfg))     // Defined in catalog_list.go
	cmd.AddCommand(newCatalogDeleteCmd(cfg))   // Defined in catalog_delete.go
	cmd.AddCommand(newCatalogGetCmd(cfg))      // Defined in catalog_update.go

	return cmd
}

// You would then add this 'catalog' command to your root command:
// rootCommand.AddCommand(newCatalogCmd(cfg))
