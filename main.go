package main

import (
	"github.com/Wifx/gonetworkmanager/v2"
	"log"
)

func checkConnectivity(nm gonetworkmanager.NetworkManager) {
	err := nm.CheckConnectivity()
	if err != nil {
		log.Panicln(err)
	}
	connectivity, err := nm.GetPropertyConnectivity()
	if err != nil {
		log.Panicln(err)
	}
	if connectivity == gonetworkmanager.NmConnectivityPortal {
		log.Println("Captive portal detected!  You probably want to open a web browser now.")
	}
}

func handleStateChange(nm gonetworkmanager.NetworkManager, state gonetworkmanager.NmState) {
	switch state {
	case gonetworkmanager.NmStateConnecting:
	case gonetworkmanager.NmStateConnectedLocal:
	case gonetworkmanager.NmStateConnectedSite:
		checkConnectivity(nm)
	}
}

func handleDeviceStateChange(nm gonetworkmanager.NetworkManager, state gonetworkmanager.NmDeviceState) {
	switch state {
	case gonetworkmanager.NmDeviceStateNeedAuth:
	case gonetworkmanager.NmDeviceStateIpCheck:
		checkConnectivity(nm)
	}
}

func main() {
	nm, err := gonetworkmanager.NewNetworkManager()
	if err != nil {
		log.Fatalln(err)
	}
	defer nm.Unsubscribe()
	ch := nm.Subscribe()
	for {
		sig := <-ch
		switch sig.Name {
		case "org.freedesktop.NetworkManager.StateChanged":
			handleStateChange(nm, sig.Body[0].(gonetworkmanager.NmState))
			break
		case "org.freedesktop.NetworkManager.Device.StateChanged":
			handleDeviceStateChange(nm, sig.Body[0].(gonetworkmanager.NmDeviceState))
			break
		}
	}
}
