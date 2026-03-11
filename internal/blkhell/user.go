package blkhell

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func (c *Cli) newUserCmd() *cobra.Command {
	userCmd := &cobra.Command{
		Use:   "user",
		Short: "Manage users",
	}

	userCmd.AddCommand(
		c.createUserCmd(),
		c.removeUserCmd(),
	)

	return userCmd
}

func (c *Cli) createUserCmd() *cobra.Command {
	var role string
	cmd := &cobra.Command{
		Use:   "create [username]",
		Short: "Create a new user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			username := args[0]

			// check for valid role
			validRoles := map[string]bool{
				"member": true,
				"admin":  true,
			}
			if !validRoles[role] {
				return fmt.Errorf("invalid role: %s", role)
			}

			fmt.Print("Enter password: ")
			passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return err
			}
			fmt.Println()

			err = c.AuthService.CreateNewUser(
				ctx,
				username,
				string(passwordBytes),
				role,
			)
			if err != nil {
				return err
			}

			fmt.Println("User created successfully")
			return nil
		},
	}

	cmd.Flags().StringVarP(&role, "role", "r", "member", "Role for the new user.")
	return cmd
}

func (c *Cli) removeUserCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove [username]",
		Short: "Remove a user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := c.config.Database.DeleteUserByUsername(
				context.Background(),
				args[0],
			)
			if err != nil {
				return err
			}

			fmt.Println("User removed successfully")
			return nil
		},
	}
}
