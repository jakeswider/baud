package list

import (
	_ "embed"
	"bufio"
	"strings"

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
	
	// FIX 1: Increase buffer to 1MB to handle long lines in usb.ids
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	currentVendor := ""

	for scanner.Scan() {
		line := scanner.Text()

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// FIX 2: Check tab depth to ignore Interface/Class lines
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
	Name 	     string
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
			Name: port.Name,
			IsUSB: port.IsUSB,
			VID: port.VID,
			PID: port.PID,
			SerialNumber: port.SerialNumber,
			ProductName: IdentifyDevice(port.VID, port.PID),
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

// Simple lookup logic (you can expand this or move it to internal/list)
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