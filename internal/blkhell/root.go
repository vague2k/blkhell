package blkhell

import (
	"github.com/spf13/cobra"
	"github.com/vague2k/blkhell/server/database"
	"github.com/vague2k/blkhell/server/services"
)

type Cli struct {
	DB          *database.Queries
	AuthService *services.AuthService
	cmd         *cobra.Command
}

func NewCli() *Cli {
	cmd := &cobra.Command{
		Use:   "blkhell",
		Short: "Blkhell CLI",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	return &Cli{
		cmd: cmd,
	}
}

func (c *Cli) Run() error {
	queries, err := database.Init()
	if err != nil {
		return err
	}
	auth := services.NewAuthService(queries)

	c.DB = queries
	c.AuthService = auth

	c.cmd.AddCommand(c.newUserCmd())
	return c.cmd.Execute()
}
