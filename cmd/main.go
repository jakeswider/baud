package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/jakeswider/baud/internal/list"
	"github.com/urfave/cli/v3"
)

func renderTable(portList []list.PortInfo, showNames bool, showAll bool) {
    headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Bold(true).Padding(0, 1)
    rowStyle := lipgloss.NewStyle().Padding(0, 1).Align(lipgloss.Left)

    headers := []string{"PORT", "USB", "VID", "PID", "SERIAL NUMBER"}
    if showNames {
        headers = append(headers, "DEVICE NAME")
    }

    t := table.New().
        Border(lipgloss.HiddenBorder()).
        Headers(headers...).
        StyleFunc(func(row, col int) lipgloss.Style {
            if row == table.HeaderRow {
                return headerStyle
            }
            // USB column color logic
            if col == 1 && portList[row].IsUSB && showAll {
                return rowStyle.Foreground(lipgloss.Color("87"))
            }
            return rowStyle
        })

    for _, port := range portList {
		if !showAll && !port.IsUSB {
			continue
		}

		isUSB := "No"
		if port.IsUSB {
			isUSB = "Yes"
		}
        row := []string{strings.TrimSpace(port.Name), isUSB, port.VID, port.PID, port.SerialNumber}
        
        if showNames {
            row = append(row, port.ProductName)
        }
        
        t.Row(row...)
    }

    fmt.Println(t)
}

func main () {
	cmd := &cli.Command{
        Name:  "baud",
        Usage: "the lightweight cli for discovering and communicating with serial ports",
        Action: func(context.Context, *cli.Command) error {
            fmt.Println("(run \"baud help\" for help)")
            return nil
        },
		Commands: []*cli.Command{
			{
				Name: "list",
				Aliases: []string{"l"},
				Usage: "lists devices connected to serial monitor",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name: "name",
						Aliases: []string{"n"},
						Usage: "Show names of connected USB devices (if they contain a VID & PID)",
					},
					&cli.BoolFlag{
						Name: "showall",
						Aliases: []string{"a"},
						Usage: "Show non USB-connected devices",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					portList, err := list.SerialListDetailed()

					if err != nil {
						log.Fatal(err)
					}
					
					if portList == nil && err == nil {
						fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Render("No ports found."))
						return nil
					} 
					renderTable(portList, cmd.Bool("name"), cmd.Bool("showall"))
					return nil
				},
			},
		},
    }


    if err := cmd.Run(context.Background(), os.Args); err != nil {
        log.Fatal(err)
    }
}	
