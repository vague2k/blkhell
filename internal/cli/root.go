package cli

import (
	"github.com/spf13/cobra"
	"github.com/vague2k/blkhell/server/auth"
	"github.com/vague2k/blkhell/server/database"
)

type App struct {
	DB   *database.Queries
	Auth *auth.Service
}

func NewRootCmd() *cobra.Command {
	app := &App{}
	cmd := &cobra.Command{
		Use:   "blkhell",
		Short: "Blkhell CLI",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			queries, err := database.Init()
			if err != nil {
				return err
			}

			app.DB = queries
			app.Auth = auth.New(queries)

			return nil
		},
	}

	cmd.AddCommand(newUserCmd(app))

	return cmd
}
