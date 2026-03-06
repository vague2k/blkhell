package blkhell

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/vague2k/blkhell/server/database"
)

func (c *Cli) newDatabaseCmd() *cobra.Command {
	databaseCmd := &cobra.Command{
		Use:   "db",
		Short: "Manage database operations",
	}

	databaseCmd.AddCommand(
		c.newMigrateCmd(),
	)

	return databaseCmd
}

func (c *Cli) newMigrateCmd() *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Apply database migrations",
	}

	migrateCmd.AddCommand(
		c.migrateUpCmd(),
		c.migrateDownCmd(),
		c.migrateForceCmd(),
	)

	return migrateCmd
}

func (c *Cli) migrateUpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Apply all pending up migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := database.Open()
			if err != nil {
				return fmt.Errorf("failed to open database: %w", err)
			}
			defer db.Close()

			if err := database.MigrateUp(db); err != nil {
				return fmt.Errorf("migration up failed: %w", err)
			}

			return nil
		},
	}
}

func (c *Cli) migrateDownCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "Roll back one migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := database.Open()
			if err != nil {
				return fmt.Errorf("failed to open database: %w", err)
			}
			defer db.Close()

			if err := database.MigrateDown(db); err != nil {
				return fmt.Errorf("migration down failed: %w", err)
			}

			return nil
		},
	}
}

func (c *Cli) migrateForceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "force [version]",
		Short: "Force the database migration version",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			version, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid version: %w", err)
			}

			db, err := database.Open()
			if err != nil {
				return fmt.Errorf("failed to open database: %w", err)
			}
			defer db.Close()

			if err := database.MigrateForce(db, version); err != nil {
				return fmt.Errorf("migration force failed: %w", err)
			}

			return nil
		},
	}
}
