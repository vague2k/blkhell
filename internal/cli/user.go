package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func newUserCmd(app *App) *cobra.Command {
	userCmd := &cobra.Command{
		Use:   "user",
		Short: "Manage users",
	}

	userCmd.AddCommand(newCreateUserCmd(app))
	userCmd.AddCommand(newRemoveUserCmd(app))

	return userCmd
}

func newCreateUserCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "create [username]",
		Short: "Create a new user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			username := args[0]

			fmt.Print("Enter password: ")
			passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return err
			}
			fmt.Println()

			err = app.Auth.CreateNewUser(
				context.Background(),
				username,
				string(passwordBytes),
				"user",
			)
			if err != nil {
				return err
			}

			fmt.Println("User created successfully")
			return nil
		},
	}
}

func newRemoveUserCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "remove [username]",
		Short: "Remove a user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := app.DB.DeleteUserByUsername(
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
