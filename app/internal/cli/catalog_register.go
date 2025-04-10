package cli

import (
	auth "duckdb-test/app/internal/auth/cli"
	"duckdb-test/app/internal/catalog"
	"duckdb-test/app/internal/client"
	"duckdb-test/app/internal/config"
	"fmt"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command

func newCatalogRegisterCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register a new catalog entry.",
		Long: `Registers a new catalog entry with the specified details.
		This command requires the 'name', 'source_type', and 'location' flags.
		The 'schema_name' and 'description' flags are optional.

		Example:
		  cli catalog register --name my_table --source_type duckdb --location '/path/to/my.db' --schema_name public --description "My important table"
		`,
		Run: func(cmd *cobra.Command, args []string) {

			tokenStr, err := auth.GetAuthToken(cfg)
			if err != nil {
				fmt.Println("❌ Not authenticated. Please run 'cli login'")
				return
			}

			// Collect errors for required flags
			var errs []error
			name := collectRequiredStringFlagError(cmd, "name", &errs)
			schemaName := collectRequiredStringFlagError(cmd, "schema_name", &errs)
			sourceType := collectRequiredStringFlagError(cmd, "source_type", &errs)
			location := collectRequiredStringFlagError(cmd, "location", &errs)

			if len(errs) > 0 {
				for _, err := range errs {
					fmt.Println(err)
				}
				return
			}

			description := getOptionalStringFlag(cmd, "description")

			input := catalog.CreateCatalogEntryParams{
				Name:        name,
				SourceType:  sourceType,
				Location:    location,
				SchemaName:  schemaName,
				Description: &description,
			}

			err = client.DoAndPrintRequest(client.RequestParams{
				Cfg:        cfg,
				Token:      *tokenStr,
				Path:       "/catalog/register",
				HttpMethod: client.MethodPost,
				Payload:    input,
			},
			)
			if err != nil {
				fmt.Println(err)
			}

		},
	}
	// Define flags and mark required ones
	cmd.Flags().String("name", "", "The name of the catalog entry (required)")
	cmd.Flags().String("source_type", "", "The source type (e.g., duckdb, postgres) (required)")
	cmd.Flags().String("location", "", "The location of the data source (required)")
	cmd.Flags().String("schema_name", "public", "The schema name for the catalog entry (optional, default: public)")
	cmd.Flags().String("description", "", "A description of the catalog entry (optional)")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("schema_name")
	cmd.MarkFlagRequired("source_type")
	cmd.MarkFlagRequired("location")

	return cmd

}
