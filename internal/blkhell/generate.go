package blkhell

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func (c *Cli) newGenerateCmd() *cobra.Command {
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate fake data",
		Long:  "Generate fake data for bands, releases, and projects (development only)",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if environment is development
			env := os.Getenv("GO_ENV")
			if env != "development" {
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
			err = c.DB.GenerateFakeData(ctx)
			if err != nil {
				return fmt.Errorf("failed to generate fake data: %w", err)
			}

			fmt.Println("Fake data generated successfully")
			return nil
		},
	}

	return generateCmd
}
