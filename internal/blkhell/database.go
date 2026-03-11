package blkhell

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

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
		c.newGenerateCmd(),
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

func (c *Cli) newGenerateCmd() *cobra.Command {
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate fake data",
		Long:  "Generate fake data for bands, releases, and projects (development only)",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if environment is development
			if c.config.Environment != "development" {
				return fmt.Errorf("this command can only be run in development environment")
			}

			reader := bufio.NewReader(os.Stdin)

			// First confirmation
			fmt.Print("This will generate fake data in your database. Are you sure? (yes/no): ")
			response1, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			response1 = strings.TrimSpace(strings.ToLower(response1))
			if response1 != "yes" && response1 != "y" {
				fmt.Println("Operation cancelled.")
				return nil
			}

			// Second confirmation
			fmt.Print("\033[1;33mWARNING\033[0m: This action cannot be undone. (yes/no): ")
			response2, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			response2 = strings.TrimSpace(strings.ToLower(response2))
			if response2 != "yes" && response2 != "y" {
				fmt.Println("Operation cancelled.")
				return nil
			}

			fmt.Println("Generating fake data...")
			ctx := context.Background()
			err = c.config.Database.GenerateFakeData(ctx)
			if err != nil {
				return fmt.Errorf("failed to generate fake data: %w", err)
			}

			fmt.Println("Fake data generated successfully")
			return nil
		},
	}

	return generateCmd
}

func (c *Cli) migrateUpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Apply all pending up migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			db := c.config.SqlDB
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
			db := c.config.SqlDB
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

			db := c.config.SqlDB
			defer db.Close()
			if err := database.MigrateForce(db, version); err != nil {
				return fmt.Errorf("migration force failed: %w", err)
			}

			return nil
		},
	}
}
