package cmd

import (
	"github.com/urfave/cli/v2"
	"os"
	"refactory/notes/internal/db/migration"
	"refactory/notes/internal/server"
)

func Start() error {
	app := cli.NewApp()
	app.Name = "RSP Notes API"
	app.Description = "Notes API is backend application created for Refactory.id RSP recruitment test."

	app.Commands = []*cli.Command{
		{
			Name:        "migrations",
			Description: "migrations looks at the currently active migration version and will migrate all the way up (applying all up migrations)",
			Action: func(c *cli.Context) error {
				err := migration.Up()
				if nil != err {
					return err
				}
				return nil
			},
		},
		{
			Name:        "rollback",
			Description: "rollbacks looks at the currently active migration version and will migrate all the way down (applying all down migrations)",
			Action: func(c *cli.Context) error {
				err := migration.Down()
				if nil != err {
					return err
				}
				return nil
			},
		},
		{
			Name:        "steps",
			Description: "steps looks at the currently active migration version. It will migrate up if n > 0, and down if n < 0",
			Flags: []cli.Flag{
				&cli.IntFlag{Name: "n"},
			},
			Action: func(c *cli.Context) error {
				err := migration.Steps(c.Int("n"))
				if err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:        "drop",
			Description: "drop deletes everything in the database",
			Action: func(c *cli.Context) error {
				err := migration.Drop()
				if err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:        "start",
			Description: "start the server",
			Action: func(c *cli.Context) error {
				return server.Start()
			},
		},
		{
			Name:        "launch",
			Description: "launch migrate all the way up (applying all up migrations) and start the server",
			Action: func(c *cli.Context) error {
				err := migration.Up()
				if err != nil {
					return err
				}
				return server.Start()
			},
		},
	}

	return app.Run(os.Args)
}
