package watch

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/jakeswider/baud/internal/list"
)

var (
	connectStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("42")).
		Padding(0, 1).
		Bold(true).
		MarginRight(1)

	disconnectStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("#D3494E")).
		Padding(0, 1).
		Bold(true).
		MarginRight(1)

	headerStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("99")).
		Padding(0, 1).
		Bold(true).
		Foreground(lipgloss.Color("99"))

	dimStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	portStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("87")).
		Bold(true)

	timerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("208")).
		Bold(true)
)

func buildPortMap(portList []list.PortInfo) map[list.PortInfo]bool {
	portMap := make(map[list.PortInfo]bool)
	for _, port := range portList {
		portMap[port] = true
	}
	return portMap
}

func WatchSerialPorts(timeToWatch uint) error {
	endTime := time.Now().Add(time.Duration(timeToWatch) * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	var timeout <-chan time.Time
	if timeToWatch > 0 {
		timeout = time.After(time.Duration(timeToWatch) * time.Second)
	}

	initialDeviceList, err := list.GetSerialListDetailed()
	if err != nil {
		return fmt.Errorf("failed to get initial serial port list: %w", err)
	}
	oldDevicesMap := buildPortMap(initialDeviceList)

	// Print the header once to establish its position
	fmt.Println(headerStyle.Render("Baud Watch Mode Active | Ctrl+C to Exit"))
	// Save the position where logs will start
	fmt.Print("\033[s") 

	for {
		select {
		case <-timeout:
			// Final update to show session finished
			fmt.Printf("\033[u\033[2F\033[K%s\033[u", headerStyle.Render("Baud Watch Mode | Session Finished"))
			fmt.Println("\nWatch time expired.")
			return nil

		case <-ticker.C:
			// 1. Build the Combined Header String
			headerText := "Baud Watch Mode Active | Ctrl+C to Exit"
			if timeToWatch > 0 {
				remaining := time.Until(endTime).Round(time.Second)
				if remaining < 0 {
					remaining = 0
				}
				headerText = fmt.Sprintf("Baud Watch Mode Active | Ctrl+C to Exit | %s", timerStyle.Render(remaining.String()))
			}

			// 2. INLINE HEADER UPDATE
			// \033[u: restore cursor to log start
			// \033[3F: move UP 3 lines (to catch the top of the bordered box)
			// \033[K: clear line
			// \033[u: restore cursor to log start for next scan
			fmt.Printf("\033[u\033[3F\033[K%s\033[u", headerStyle.Render(headerText))

			// 3. Scan and Compare
			newDevicesList, _ := list.GetSerialListDetailed()
			newDevicesMap := buildPortMap(newDevicesList)
			ts := dimStyle.Render(time.Now().Format("15:04:05") + " ")

			for port := range newDevicesMap {
				if _, found := oldDevicesMap[port]; !found {
					msg := fmt.Sprintf("%s%s %s", ts, connectStyle.Render("CONNECTED"), portStyle.Render(port.Name))
					if port.ProductName != "-" {
						msg += dimStyle.Render(" (" + port.ProductName + ")")
					}
					fmt.Println(msg)
					// Update the "saved position" so the header doesn't move up
					fmt.Print("\033[s") 
				}
			}

			for port := range oldDevicesMap {
				if _, found := newDevicesMap[port]; !found {
					msg := fmt.Sprintf("%s%s %s", ts, disconnectStyle.Render("DISCONNECTED"), portStyle.Render(port.Name))
					if port.ProductName != "-" {
						msg += dimStyle.Render(" (" + port.ProductName + ")")
					}
					fmt.Println(msg)
					// Update the "saved position"
					fmt.Print("\033[s")
				}
			}

			oldDevicesMap = newDevicesMap
		}
	}
}