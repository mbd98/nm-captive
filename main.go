package main

import (
	"fmt"
	"github.com/Wifx/gonetworkmanager/v2"
	"log"
)

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
		sig := <-ch
		if sig.Name == "org.freedesktop.NetworkManager.StateChanged" {
			state := gonetworkmanager.NmState(sig.Body[0].(uint32))
			if state == gonetworkmanager.NmStateConnectedSite {
				c, err := nm.GetPropertyConnectivity()
				if err != nil {
					log.Panicln(err)
				}
				if c == gonetworkmanager.NmConnectivityPortal {
					log.Println("Captive portal detected")
					fmt.Println("You probably want to open a web browser now")
				}
			}
		}
	}
}
