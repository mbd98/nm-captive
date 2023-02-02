package main

import (
	"github.com/Wifx/gonetworkmanager/v2"
	"log"
	"net/http"
)

func determineUrl() {
	resp, err := http.Get("http://detectportal.firefox.com/canonical.html")
	if err != nil {
		log.Panicln(err)
	}
	log.Println(resp)
	loc, err := resp.Location()
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Portal URL: %v\n", loc)
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
		if sig.Name == "org.freedesktop.NetworkManager.Device.StateChanged" {
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
				go determineUrl()
			}
		}
	}
}
