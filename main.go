package main

import (
	"fmt"
	"github.com/Wifx/gonetworkmanager/v2"
	"log"
)

func checkConnectivity(nm gonetworkmanager.NetworkManager) {
	// Ask NM to run the connectivity test: https://networkmanager.dev/docs/api/latest/gdbus-org.freedesktop.NetworkManager.html#gdbus-method-org-freedesktop-NetworkManager.CheckConnectivity
	err := nm.CheckConnectivity()
	if err != nil {
		log.Panicln(err)
	}
	// Get test results: https://networkmanager.dev/docs/api/latest/gdbus-org.freedesktop.NetworkManager.html#gdbus-property-org-freedesktop-NetworkManager.Connectivity
	connectivity, err := nm.GetPropertyConnectivity()
	if err != nil {
		log.Panicln(err)
	}
	// Portal detected? https://networkmanager.dev/docs/api/latest/nm-dbus-types.html#NMConnectivityState
	if connectivity == gonetworkmanager.NmConnectivityPortal {
		log.Println("Captive portal detected")
		fmt.Println("You probably want to open a web browser now.")
	}
}

func handleStateChange(nm gonetworkmanager.NetworkManager, state gonetworkmanager.NmState) {
	// https://networkmanager.dev/docs/api/latest/nm-dbus-types.html#NMState
	if state == gonetworkmanager.NmStateConnectedSite {
		log.Printf("state changed: %v\n", state)
		checkConnectivity(nm)
	}
}

func handleDeviceStateChange(nm gonetworkmanager.NetworkManager, state gonetworkmanager.NmDeviceState) {
	// https://networkmanager.dev/docs/api/latest/nm-dbus-types.html#NMDeviceState
	if state == gonetworkmanager.NmDeviceStateIpCheck {
		log.Printf("device state changed: %v\n", state)
		checkConnectivity(nm)
	}
}

func main() {
	nm, err := gonetworkmanager.NewNetworkManager()
	if err != nil {
		log.Fatalln(err)
	}
	// Stop listening for signals when done
	defer nm.Unsubscribe()
	// Subscribe to the bus for new signals
	ch := nm.Subscribe()
	for {
		// Get the next signal
		sig := <-ch
		// Filter signals we care about - should be able to get by with just the one
		switch sig.Name {
		// https://networkmanager.dev/docs/api/latest/gdbus-org.freedesktop.NetworkManager.html#gdbus-signal-org-freedesktop-NetworkManager.StateChanged
		case "org.freedesktop.NetworkManager.StateChanged":
			handleStateChange(nm, gonetworkmanager.NmState(sig.Body[0].(uint32)))
			break
		// https://networkmanager.dev/docs/api/latest/gdbus-org.freedesktop.NetworkManager.Device.html#gdbus-signal-org-freedesktop-NetworkManager-Device.StateChanged
		case "org.freedesktop.NetworkManager.Device.StateChanged":
			handleDeviceStateChange(nm, gonetworkmanager.NmDeviceState(sig.Body[0].(uint32)))
			break
		}
	}
}
