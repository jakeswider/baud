package watch

import (
	"fmt"
	"time"

	"github.com/jakeswider/baud/internal/list"
)

func buildPortMap(portList []list.PortInfo) map[list.PortInfo]bool {
	portMap := make(map[list.PortInfo]bool)
	for _, port := range portList {
		portMap[port] = true
	}
	return portMap
}

func WatchSerialPorts(timeToWatch uint) error {
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
	for {
		select {
		case <-timeout:
			fmt.Println("Watch time expired.")
			return nil
		case <-ticker.C:
			newDevicesList, err := list.GetSerialListDetailed()
			if err != nil {
				return fmt.Errorf("failed to get serial port list: %w", err)
			}
			newDevicesMap := buildPortMap(newDevicesList)

			//connections
			for port := range newDevicesMap {
				if _, found := oldDevicesMap[port]; !found {
					if port.ProductName != "-" {
						fmt.Printf("✅ CONNECTED:    %s Product Name: (%s)\n", port.Name, port.ProductName)
					} else {
						fmt.Printf("✅ CONNECTED:    %s \n", port.Name)
					}
				}
			}

			//disconnections
			for port := range oldDevicesMap {
				if _, found := newDevicesMap[port]; !found {
					if port.ProductName != "-" {
						fmt.Printf("❌ DISCONNECTED: %s Product Name: (%s)\n", port.Name, port.ProductName)
					} else {
						fmt.Printf("❌ DISCONNECTED: %s\n", port.Name)
					}
				}
			}
			oldDevicesMap = newDevicesMap
		}
	}
}
