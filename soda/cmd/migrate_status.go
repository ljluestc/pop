package cmd

import (
	"os"

	"github.com/gobuffalo/pop/v6"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Displays the status of all migrations.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check for build errors before attempting migration
		if err := checkBuildErrors(); err != nil {
			return errors.WithStack(err)
		}

		mig, err := pop.NewFileMigrator(migrationPath, getConn())
		if err != nil {
			return err
		}
		return mig.Status(os.Stdout)
	},
}

func init() {
	migrateCmd.AddCommand(migrateStatusCmd)
}
