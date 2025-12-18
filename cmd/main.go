package main

import (
	"context"
	"fmt"
	"log"
	"os"
	
	"github.com/jakeswider/baud/internal/list"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "baud",
		Usage: "the lightweight cli for discovering and communicating with serial ports",
		Action: func(context.Context, *cli.Command) error {
			fmt.Println("(run \"baud help\" for help)")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "lists devices connected to serial monitor",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "Show names of connected USB devices (if they contain a VID & PID)",
					},
					&cli.BoolFlag{
						Name:    "showall",
						Aliases: []string{"a"},
						Usage:   "Show non USB-connected devices",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					portList, err := list.GetSerialListDetailed()

					if err != nil {
						log.Fatal(err)
					}

					list.RenderTable(portList, cmd.Bool("name"), cmd.Bool("showall"))
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
