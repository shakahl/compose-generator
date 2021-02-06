package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"compose-generator/commands"
)

const VERSION = "1.0.0"

func main() {
	// Version flag
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Prints the version of compose-generator",
	}

	// Main cli configuration
	app := &cli.App{
		Name:    "compose-generator",
		Version: VERSION,
		Authors: []*cli.Author{
			{Name: "Marc Auberer", Email: "marc.auberer@chillibits.com"},
		},
		Copyright: "© 2021 Marc Auberer",
		Usage:     "Generate and manage docker compose configuration files for your projects.",
		Action: func(c *cli.Context) error {
			commands.Generate(c.Bool("advanced"), c.Bool("run"), c.Bool("demonized"), c.Bool("force"))
			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "advanced", Aliases: []string{"a"}, Usage: "Generate compose file in advanced mode"},
			&cli.BoolFlag{Name: "run", Aliases: []string{"r"}, Usage: "Run docker-compose after creating the compose file"},
			&cli.BoolFlag{Name: "demonized", Aliases: []string{"d"}, Usage: "Run docker-compose demonized after creating the compose file"},
			&cli.BoolFlag{Name: "force", Aliases: []string{"f"}, Usage: "No safety checks"},
		},
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "Adds a service to an existing compose file",
				Action: func(c *cli.Context) error {
					commands.Add()
					return nil
				},
			},
			{
				Name:    "remove",
				Aliases: []string{"r", "rm"},
				Usage:   "Removes a service from an existing compose file",
				Action: func(c *cli.Context) error {
					commands.Remove()
					return nil
				},
			},
			{
				Name:    "template",
				Aliases: []string{"t"},
				Usage:   "Saves / loads snapshots of your compose configuration for later use",
				Subcommands: []*cli.Command{
					{
						Name:    "save",
						Aliases: []string{"s"},
						Usage:   "Save a custom template.",
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "stash", Aliases: []string{"s"}, Usage: "Move the regarding files, instead of copying them"},
							&cli.BoolFlag{Name: "force", Aliases: []string{"f"}, Usage: "No safety checks"},
						},
						Action: func(c *cli.Context) error {
							name := c.Args().Get(0)
							commands.SaveTemplate(name, c.Bool("stash"), c.Bool("force"))
							return nil
						},
					},
					{
						Name:    "load",
						Aliases: []string{"l"},
						Usage:   "Load a custom template.",
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "show-predefined", Aliases: []string{"p"}, Usage: "Show predefined templates in addition to the custom ones"},
							&cli.BoolFlag{Name: "force", Aliases: []string{"f"}, Usage: "No safety checks"},
						},
						Action: func(c *cli.Context) error {
							name := c.Args().Get(0)
							commands.LoadTemplate(name, c.Bool("force"))
							return nil
						},
					},
				},
			},
			{
				Name:    "install",
				Aliases: []string{"i"},
				Usage:   "Installs Docker and Docker Compose with a single command",
				Action: func(c *cli.Context) error {
					commands.Install()
					return nil
				},
			},
		},
		UseShortOptionHandling: true,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
