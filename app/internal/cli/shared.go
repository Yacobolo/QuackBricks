package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// Helper function to get a required string flag and collect errors
func collectRequiredStringFlagError(cmd *cobra.Command, name string, errs *[]error) string {
	value, err := cmd.Flags().GetString(name)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("error getting '%s' flag: %w", name, err))
		return ""
	}
	if strings.TrimSpace(value) == "" {
		*errs = append(*errs, fmt.Errorf("the '%s' flag is required", name))
		return ""
	}
	return value
}

// Helper function to get an optional string flag value
func getOptionalStringFlag(cmd *cobra.Command, name string) string {
	value, _ := cmd.Flags().GetString(name) // Error is ignored as the flag is optional
	return value
}
