package main

import (
	"bufio"
	"fmt"
	"github.com/Wifx/gonetworkmanager/v2"
	"log"
	"net/http"
)

func printPortalBody() {
	resp, err := http.Get("http://detectportal.firefox.com/canonical.html")
	if err == nil {
		scan := bufio.NewScanner(resp.Body)
		for scan.Scan() {
			fmt.Println(scan.Text())
		}
		err := scan.Err()
		if err != nil {
			log.Printf("Error reading body: %v\n", err)
		}
		err = resp.Body.Close()
		if err != nil {
			log.Printf("Error closing body: %v\n", err)
		}
	} else {
		log.Printf("Error making request: %v\n", err)
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
				go printPortalBody()
			}
		}
	}
}
