package blkhell

import (
	"github.com/spf13/cobra"
	"github.com/vague2k/blkhell/config"
	"github.com/vague2k/blkhell/server/database"
	"github.com/vague2k/blkhell/server/services"
)

type Cli struct {
	DB          *database.Queries
	AuthService *services.AuthService
	cmd         *cobra.Command
}

func NewCli(cfg *config.Config) *Cli {
	cmd := &cobra.Command{
		Use:   "blkhell",
		Short: "Blkhell CLI",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	return &Cli{
		DB:  cfg.Database,
		cmd: cmd,
	}
}

func (c *Cli) Run() error {
	c.AuthService = services.NewAuthService(c.DB)

	c.cmd.AddCommand(c.newUserCmd())
	c.cmd.AddCommand(c.newDatabaseCmd())
	return c.cmd.Execute()
}
