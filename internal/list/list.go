package list

import (
	"bufio"
	_ "embed"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

//go:embed usb.ids
var usbDB string

var vendorMap = make(map[string]string)
var productMap = make(map[string]string)

func normID(s string) string {
	return strings.ToLower(strings.TrimPrefix(s, "0x"))
}

func init() {
	scanner := bufio.NewScanner(strings.NewReader(usbDB))

	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	currentVendor := ""

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		tabs := 0
		for i := 0; i < len(line) && line[i] == '\t'; i++ {
			tabs++
		}

		trimmed := strings.TrimSpace(line)
		parts := strings.Fields(trimmed)
		if len(parts) < 2 {
			continue
		}

		id := normID(parts[0])
		name := strings.Join(parts[1:], " ")

		switch tabs {
		case 0: // Vendor line
			currentVendor = id
			vendorMap[currentVendor] = name
		case 1: // Product line
			if currentVendor != "" {
				productMap[currentVendor+":"+id] = name
			}
		default:
			// Ignore lines with 2+ tabs (Interfaces/Protocols)
			continue
		}
	}
}

type PortInfo struct {
	Name         string
	IsUSB        bool
	VID          string
	PID          string
	SerialNumber string
	ProductName  string
}

func PortDetailstoPortInfo(pList []*enumerator.PortDetails) []PortInfo {
	PortInfoList := make([]PortInfo, 0, len(pList))
	for _, port := range pList {
		PortInfoList = append(PortInfoList, PortInfo{
			Name:         port.Name,
			IsUSB:        port.IsUSB,
			VID:          port.VID,
			PID:          port.PID,
			SerialNumber: port.SerialNumber,
			ProductName:  IdentifyDevice(port.VID, port.PID),
		})
	}
	return PortInfoList
}

func SerialList() ([]string, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return []string{""}, err
	}
	if len(ports) == 0 {
		return []string{"No ports found"}, nil
	}

	return ports, nil
}

func SerialListDetailed() ([]PortInfo, error) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return nil, err
	}
	if len(ports) == 0 {
		return nil, nil
	}
	PortInfoList := PortDetailstoPortInfo(ports)
	return PortInfoList, nil
}

func IdentifyDevice(vid, pid string) string {
	v := normID(vid)
	p := normID(pid)

	if v == "" {
		return "-"
	}

	vName, vFound := vendorMap[v]
	if !vFound {
		return "-"
	}

	pName, pFound := productMap[v+":"+p]
	if !pFound {
		return vName
	}

	return vName + " " + pName

}
func usbPortList(portList []PortInfo) []PortInfo {
	usbPorts := make([]PortInfo, 0)
	for _, port := range portList {
		if port.IsUSB {
			usbPorts = append(usbPorts, port)
		}
	}
	return usbPorts
}

func RenderTable(portList []PortInfo, showNames bool, showAll bool) {
	usbPorts := usbPortList(portList)

	if len(usbPorts) == 0 && !showAll {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Render("No USB connections found. (consider running with --showall)"))
		return
	}
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
