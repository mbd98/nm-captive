package main

import (
	"github.com/Wifx/gonetworkmanager/v2"
	"log"
)

func main() {
	nm, err := gonetworkmanager.NewNetworkManager()
	if err != nil {
		log.Fatalln(err)
	}
	defer nm.Unsubscribe()
	ch := nm.Subscribe()
	for {
		sig := <-ch
		if sig.Name == "org.freedesktop.NetworkManager.Device.StateChanged" {
			state := gonetworkmanager.NmState(sig.Body[0].(uint32))
			if state == gonetworkmanager.NmStateConnectedSite {
				err := nm.CheckConnectivity()
				if err != nil {
					log.Panicln(err)
				}
				connectivity, err := nm.GetPropertyConnectivity()
				if err != nil {
					log.Panicln(err)
				}
				if connectivity == gonetworkmanager.NmConnectivityPortal {
					log.Println("Captive portal detected")
				} else {
					log.Println("Other connectivity state")
				}
			}
		} else {
			log.Println("Other state change")
		}
	}
}
